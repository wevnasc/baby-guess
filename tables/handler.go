package tables

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/wevnasc/baby-guess/db"
	"github.com/wevnasc/baby-guess/server"
)

type Handler struct {
	ctrl *controller
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

	type itemResponse struct {
		ID          uuid.UUID  `json:"id"`
		Description string     `json:"description"`
		AccountID   *uuid.UUID `json:"account_id,omitempty"`
	}

	type tablesResponse struct {
		ID    uuid.UUID      `json:"id"`
		Name  string         `json:"name"`
		Items []itemResponse `json:"items"`
	}

	return server.ErrorHandler(func(rw http.ResponseWriter, r *http.Request) error {
		tables, err := h.ctrl.all(r.Context(), server.AccountUUID(r))

		if err != nil {
			return err
		}

		var tt = []tablesResponse{}

		for _, table := range tables {

			var ii = []itemResponse{}

			for _, item := range table.items {

				i := itemResponse{
					ID:          item.id,
					Description: item.description,
					AccountID:   item.owner.nullableID(),
				}

				ii = append(ii, i)
			}

			t := tablesResponse{
				ID:    table.id,
				Name:  table.name,
				Items: ii,
			}

			tt = append(tt, t)
		}

		server.Json(rw, tt, http.StatusOK)
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

func (h *Handler) SetupRoutes(r *mux.Router) {
	tRouter := r.PathPrefix("/tables").Subrouter()
	tRouter.Use(server.Authenticated)
	tRouter.HandleFunc("", h.createTablesHandler()).Methods(http.MethodPost)
	tRouter.HandleFunc("", h.allTablesHandler()).Methods(http.MethodGet)

	iRouter := r.PathPrefix("/tables/{table_id}/items/{item_id}").Subrouter()
	iRouter.Use(server.Authenticated, server.ParseUUID("table_id", "item_id"))

	iRouter.HandleFunc("/select", h.selectItemHandler()).Methods(http.MethodPost)
	iRouter.HandleFunc("/unselect", h.unselectItemHandler()).Methods(http.MethodPost)
	iRouter.HandleFunc("/approve", h.approveItemHandler()).Methods(http.MethodPost)
}

func NewHandler(db *db.Store) *Handler {
	ctrl := newController(newDatabase(db))
	return &Handler{ctrl}
}
