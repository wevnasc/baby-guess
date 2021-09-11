package server

import "net/http"

type reason int

const (
	ResourceAlreadyExists reason = http.StatusConflict
	ResourceInvalid              = http.StatusBadRequest
	ResourceNotFound             = http.StatusNotFound
	ResourceParse                = http.StatusBadRequest
	OperationNotAllowed          = http.StatusForbidden
	OperationError               = http.StatusBadRequest
	URLParse                     = http.StatusBadRequest
	UnkownError                  = http.StatusBadRequest
)

type Error struct {
	Message string `json:"message"`
	Reason  reason `json:"-"`
}

func NewError(message string, reason reason) *Error {
	return &Error{Message: message, Reason: reason}
}

func (e *Error) Json(rw http.ResponseWriter) {
	Json(rw, e, int(e.Reason))
}

func (e *Error) Error() string {
	return e.Message
}
