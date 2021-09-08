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
			return server.NewError("error to get account id", server.URLParse)
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

func (h *Handler) SetupRoutes(r *mux.Router) {
	r.Methods(http.MethodPost).Subrouter().HandleFunc("/accounts/{id}/tables", h.createTablesHandler())
}

func NewHandler(db *db.Store) *Handler {
	ctrl := newController(newDatabase(db))
	return &Handler{ctrl}
}
