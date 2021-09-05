package server

import "net/http"

type Error struct {
	Message string `json:"message"`
	Code    int    `json:"-"`
}

func NewError(message string, code int) *Error {
	return &Error{Message: message, Code: code}
}

func (e *Error) Json(rw http.ResponseWriter) {
	Json(rw, e, e.Code)
}

func (e *Error) Error() string {
	return e.Message
}
