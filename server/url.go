package server

import (
	"net/http"

	"github.com/google/uuid"
)

func PathUUID(r *http.Request, id string) uuid.UUID {
	return r.Context().Value("account_id").(uuid.UUID)
}
