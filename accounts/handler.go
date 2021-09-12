package accounts

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

func (h *Handler) createAccountsHandler() http.HandlerFunc {

	type request struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		ID    uuid.UUID `json:"id"`
		Name  string    `json:"name"`
		Email string    `json:"email"`
	}

	return server.ErrorHandler(func(rw http.ResponseWriter, r *http.Request) error {

		var body request

		err := json.NewDecoder(r.Body).Decode(&body)

		if err != nil {
			return server.NewError("Error to parse resource", server.ResourceParse)
		}

		account, err := newAccount(body.Name, body.Password, body.Email)

		if err != nil {
			return server.NewError("Error to create account", server.ResourceInvalid)
		}

		account, err = h.ctrl.create(r.Context(), account)

		if err != nil {
			return err
		}

		res := &response{
			ID:    account.id,
			Name:  account.name,
			Email: account.email,
		}

		server.Json(rw, res, http.StatusCreated)
		return nil
	})
}

func (h *Handler) SetupRoutes(r *mux.Router) {
	r.Methods(http.MethodPost).Subrouter().HandleFunc("/accounts", h.createAccountsHandler())
}

func NewHandler(db *db.Store) *Handler {
	ctrl := newController(newDatabase(db))
	return &Handler{ctrl}
}
