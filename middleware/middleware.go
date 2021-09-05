package middleware

import (
	"net/http"

	"github.com/wevnasc/baby-guess/server"
)

func Headers(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-type", "application/json")
		h.ServeHTTP(rw, r)
	})
}

type ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request) error

func ErrorHandler(next ErrorHandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		err := next(rw, r)

		if err == nil {
			return
		}

		if httpError, ok := err.(*server.Error); ok {
			httpError.Json(rw)
			return
		}

		server.NewError(err.Error(), server.UnkownError).Json(rw)
	}
}
