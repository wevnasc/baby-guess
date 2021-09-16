package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/wevnasc/baby-guess/config"
	"github.com/wevnasc/baby-guess/db"
	"github.com/wevnasc/baby-guess/server"
	"github.com/wevnasc/baby-guess/token"
)

type Handler struct {
	ctrl   *controller
	config *config.Config
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

func (h *Handler) loginHandler() http.HandlerFunc {
	type response struct {
		Token string `json:"token"`
	}

	return server.ErrorHandler(func(rw http.ResponseWriter, r *http.Request) error {
		authorization := r.Header.Get("Authorization")

		credentials, err := token.BasicAuth(authorization)

		if err != nil {
			return server.NewError(err.Error(), server.OperationError)
		}

		account, err := h.ctrl.login(r.Context(), credentials)

		if err != nil {
			return err
		}

		token, err := token.NewAuth(account.id, h.config.Secret, time.Hour*24)

		if err != nil {
			return server.NewError("not was possible to authenticate the account", server.ResourceInvalid)
		}

		server.Json(rw, &response{Token: token}, http.StatusOK)
		return nil
	})
}

func (h *Handler) SetupRoutes(r *mux.Router) {
	r.Methods(http.MethodPost).Subrouter().HandleFunc("/accounts", h.createAccountsHandler())
	r.Methods(http.MethodGet).Subrouter().HandleFunc("/login", h.loginHandler())
}

func NewHandler(db *db.Store, config *config.Config) *Handler {
	ctrl := newController(newDatabase(db))
	return &Handler{ctrl, config}
}
