package accounts

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/wevnasc/baby-guess/server"
)

type Handler struct {
	*server.Middleware
	ctrl *controller
}

func (h *Handler) postAccountsHandler(rw http.ResponseWriter, req *http.Request) {

	type request struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	var body request

	err := json.NewDecoder(req.Body).Decode(&body)

	if err != nil {
		server.NewError("error to parse body", http.StatusBadRequest).Json(rw)
		return
	}

	account, err := newAccount(body.Name, body.Password, body.Email)

	if err != nil {
		server.NewError("not was possible to create the account", http.StatusBadRequest).Json(rw)
		return
	}

	account, err = h.ctrl.create(req.Context(), account)

	if err != nil {
		server.NewError(err.Error(), http.StatusBadRequest).Json(rw)
		return
	}

	res := &response{
		ID:    account.id.String(),
		Name:  account.name,
		Email: account.email,
	}

	server.Json(rw, res, http.StatusCreated)
}

func (h *Handler) accountsHandler(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		h.postAccountsHandler(rw, req)
	}
}

func (h *Handler) SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/accounts", h.Resource(h.accountsHandler, []string{http.MethodPost}))
}

func NewHandler(logger *log.Logger, db *sql.DB) *Handler {
	ctrl := newController(newDatabase(db))
	return &Handler{&server.Middleware{Log: logger}, ctrl}
}
