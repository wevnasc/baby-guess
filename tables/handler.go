package tables

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
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

		uuid, err := uuid.Parse(mux.Vars(r)["id"])

		if err != nil {
			return server.NewError("error to parse account id", server.URLParse)
		}

		table := newTable(uuid, body.Name, body.Items)
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

		accountID, err := uuid.Parse(mux.Vars(r)["account_id"])

		if err != nil {
			return server.NewError("error to parse account id", server.URLParse)
		}

		tableID, err := uuid.Parse(mux.Vars(r)["table_id"])

		if err != nil {
			return server.NewError("error to parse table id", server.URLParse)
		}

		itemID, err := uuid.Parse(mux.Vars(r)["id"])

		if err != nil {
			return server.NewError("error to parse item id", server.URLParse)
		}

		item := item{
			owner: &owner{accountID},
			id:    itemID,
		}

		err = h.ctrl.selectItem(r.Context(), tableID, item)

		if err != nil {
			return err
		}

		w.WriteHeader(http.StatusNoContent)
		return nil
	})
}

func (h *Handler) SetupRoutes(r *mux.Router) {
	// TODO to use account ID inside authentication token
	r.Methods(http.MethodPost).Subrouter().HandleFunc("/accounts/{id}/tables", h.createTablesHandler())
	r.Methods(http.MethodPost).Subrouter().HandleFunc("/accounts/{account_id}/tables/{table_id}/items/{id}/selected", h.selectItemHandler())
}

func NewHandler(db *db.Store) *Handler {
	ctrl := newController(newDatabase(db))
	return &Handler{ctrl}
}
