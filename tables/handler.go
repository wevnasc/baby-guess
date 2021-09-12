package tables

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/wevnasc/baby-guess/db"
	"github.com/wevnasc/baby-guess/middleware"
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

	return middleware.ErrorHandler(func(rw http.ResponseWriter, r *http.Request) error {
		var body request

		err := json.NewDecoder(r.Body).Decode(&body)

		if err != nil {
			return server.NewError("error to parse body", server.ResourceParse)
		}

		table := newTable(server.PathUUID(r, "account_id"), body.Name, body.Items)
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

func (h *Handler) selectItemHandler() http.HandlerFunc {
	return middleware.ErrorHandler(func(w http.ResponseWriter, r *http.Request) error {
		item := item{
			owner: newOwner(server.PathUUID(r, "account_id")),
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
	return middleware.ErrorHandler(func(w http.ResponseWriter, r *http.Request) error {
		item := item{
			owner: newOwner(server.PathUUID(r, "account_id")),
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
	return middleware.ErrorHandler(func(w http.ResponseWriter, r *http.Request) error {
		if err := h.ctrl.approveItem(
			r.Context(),
			newOwner(server.PathUUID(r, "account_id")),
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
	aRouter := r.PathPrefix("/accounts/{account_id}").Subrouter()
	aRouter.Use(middleware.ParseUUID("account_id"))

	tRouter := aRouter.PathPrefix("/tables").Subrouter()
	tRouter.HandleFunc("", h.createTablesHandler()).Methods(http.MethodPost)

	iRouter := aRouter.PathPrefix("/tables/{table_id}/items/{item_id}").Subrouter()
	iRouter.Use(middleware.ParseUUID("table_id", "item_id"))

	iRouter.HandleFunc("/select", h.selectItemHandler()).Methods(http.MethodPost)
	iRouter.HandleFunc("/unselect", h.unselectItemHandler()).Methods(http.MethodPost)
	iRouter.HandleFunc("/approve", h.approveItemHandler()).Methods(http.MethodPost)
}

func NewHandler(db *db.Store) *Handler {
	ctrl := newController(newDatabase(db))
	return &Handler{ctrl}
}
