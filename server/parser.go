package server

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func Json(rw http.ResponseWriter, body interface{}, code int) {
	rw.WriteHeader(code)
	value, err := json.Marshal(body)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = rw.Write(value)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func AccountUUID(r *http.Request) uuid.UUID {
	return r.Context().Value("account_id").(uuid.UUID)
}
