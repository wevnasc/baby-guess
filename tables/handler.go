package tables

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/wevnasc/baby-guess/config"
	"github.com/wevnasc/baby-guess/db"
	"github.com/wevnasc/baby-guess/email"
	"github.com/wevnasc/baby-guess/server"
)

type Handler struct {
	ctrl   *controller
	config *config.Config
}

func (h *Handler) createTablesHandler() http.HandlerFunc {
	type request struct {
		Name  string `json:"name"`
		Items int    `json:"items"`
	}

	type response struct {
		ID string `json:"id"`
	}

	return server.ErrorHandler(func(rw http.ResponseWriter, r *http.Request) error {
		var body request

		err := json.NewDecoder(r.Body).Decode(&body)

		if err != nil {
			return server.NewError("error to parse body", server.ResourceParse)
		}

		table, err := newTable(server.PathUUID(r, "account_id"), body.Name, body.Items)

		if err != nil {
			return server.NewError(err.Error(), server.ResourceInvalid)
		}

		table, err = h.ctrl.create(r.Context(), table)

		if err != nil {
			return err
		}

		res := &response{
			ID: table.id.String(),
		}

		server.Json(rw, res, http.StatusCreated)
		return nil
	})
}

func (h *Handler) allTablesHandler() http.HandlerFunc {

	type tablesResponse struct {
		ID   uuid.UUID `json:"id"`
		Name string    `json:"name"`
	}

	return server.ErrorHandler(func(rw http.ResponseWriter, r *http.Request) error {
		tables, err := h.ctrl.allTables(r.Context(), server.AccountUUID(r))

		if err != nil {
			return err
		}

		tt := make([]tablesResponse, len(tables))
		for i := range tables {

			current := tables[i]

			tt[i] = tablesResponse{
				ID:   current.id,
				Name: current.name,
			}

		}

		server.Json(rw, tt, http.StatusOK)
		return nil
	})

}

func (h *Handler) oneTablesHandler() http.HandlerFunc {

	type ownerResponse struct {
		ID    uuid.UUID `json:"id"`
		Name  string    `json:"name"`
		Email string    `json:"email"`
	}

	type itemResponse struct {
		ID          uuid.UUID      `json:"id"`
		Description string         `json:"description"`
		Owner       *ownerResponse `json:"owner,omitempty"`
	}

	type tablesResponse struct {
		ID    uuid.UUID      `json:"id"`
		Name  string         `json:"name"`
		Items []itemResponse `json:"items"`
	}

	toResponse := func(table table) tablesResponse {

		var ii = []itemResponse{}

		for _, item := range table.items {

			var o *ownerResponse
			if item.owner.id.Valid {
				o = &ownerResponse{
					ID:    item.owner.id.UUID,
					Name:  item.owner.name,
					Email: item.owner.email,
				}
			}

			i := itemResponse{
				ID:          item.id,
				Description: item.description,
				Owner:       o,
			}

			ii = append(ii, i)
		}

		t := tablesResponse{
			ID:    table.id,
			Name:  table.name,
			Items: ii,
		}

		return t
	}

	return server.ErrorHandler(func(rw http.ResponseWriter, r *http.Request) error {
		table, err := h.ctrl.oneTable(r.Context(), server.PathUUID(r, "table_id"))

		if err != nil {
			return err
		}

		response := toResponse(*table)
		server.Json(rw, response, http.StatusOK)
		return nil
	})
}

func (h *Handler) selectItemHandler() http.HandlerFunc {
	return server.ErrorHandler(func(w http.ResponseWriter, r *http.Request) error {
		item := item{
			owner: newOwner(server.AccountUUID(r)),
			id:    server.PathUUID(r, "item_id"),
		}

		if err := h.ctrl.selectItem(r.Context(), server.PathUUID(r, "table_id"), item); err != nil {
			return err
		}

		w.WriteHeader(http.StatusNoContent)
		return nil
	})
}

func (h *Handler) unselectItemHandler() http.HandlerFunc {
	return server.ErrorHandler(func(w http.ResponseWriter, r *http.Request) error {
		item := item{
			owner: newOwner(server.AccountUUID(r)),
			id:    server.PathUUID(r, "item_id"),
		}

		if err := h.ctrl.unselectItem(r.Context(), server.PathUUID(r, "table_id"), item); err != nil {
			return err
		}

		w.WriteHeader(http.StatusNoContent)
		return nil
	})
}

func (h *Handler) approveItemHandler() http.HandlerFunc {
	return server.ErrorHandler(func(w http.ResponseWriter, r *http.Request) error {
		if err := h.ctrl.approveItem(
			r.Context(),
			newOwner(server.AccountUUID(r)),
			server.PathUUID(r, "table_id"),
			server.PathUUID(r, "item_id"),
		); err != nil {
			return err
		}

		w.WriteHeader(http.StatusNoContent)
		return nil
	})
}

func (h *Handler) drawHandler() http.HandlerFunc {

	type response struct {
		ID        uuid.UUID  `json:"id"`
		AccountID *uuid.UUID `json:"account_id,omitempty"`
	}

	return server.ErrorHandler(func(w http.ResponseWriter, r *http.Request) error {
		item, err := h.ctrl.draw(
			r.Context(),
			newOwner(server.AccountUUID(r)),
			server.PathUUID(r, "table_id"),
		)

		if err != nil {
			return err
		}

		server.Json(w, response{ID: item.id, AccountID: item.owner.nullableID()}, http.StatusOK)
		return nil
	})

}

func (h *Handler) SetupRoutes(r *mux.Router) {
	//private routes
	tRouter := r.PathPrefix("/tables").Subrouter()
	tRouter.Use(server.Auth(h.config))
	tRouter.HandleFunc("", h.createTablesHandler()).Methods(http.MethodPost)
	tRouter.HandleFunc("", h.allTablesHandler()).Methods(http.MethodGet)

	tdRouter := r.PathPrefix("/tables/{table_id}").Subrouter()
	tdRouter.Use(server.Auth(h.config), server.ParseUUID("table_id"))
	tdRouter.HandleFunc("/draw", h.drawHandler()).Methods(http.MethodPost)

	iRouter := r.PathPrefix("/tables/{table_id}/items/{item_id}").Subrouter()
	iRouter.Use(server.Auth(h.config), server.ParseUUID("table_id", "item_id"))

	iRouter.HandleFunc("/select", h.selectItemHandler()).Methods(http.MethodPost)
	iRouter.HandleFunc("/unselect", h.unselectItemHandler()).Methods(http.MethodPost)
	iRouter.HandleFunc("/approve", h.approveItemHandler()).Methods(http.MethodPost)

	//public routes
	tdpRouter := r.PathPrefix("/tables/{table_id}").Subrouter()
	tdpRouter.Use(server.ParseUUID("table_id"))
	tdpRouter.HandleFunc("", h.oneTablesHandler()).Methods(http.MethodGet)
}

func NewHandler(db *db.Store, config *config.Config, email email.Client) *Handler {
	ctrl := newController(newDatabase(db), email)
	return &Handler{ctrl, config}
}
