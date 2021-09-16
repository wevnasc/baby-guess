package server

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/wevnasc/baby-guess/config"
	"github.com/wevnasc/baby-guess/token"
)

type ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request) error

func Headers(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-type", "application/json")
		h.ServeHTTP(rw, r)
	})
}

func ParseUUID(keys ...string) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

			vars := mux.Vars(r)
			ctx := r.Context()

			for _, key := range keys {

				value, ok := vars[key]

				if !ok {
					NewError("not found parameter %s on the url", URLParse)
					return
				}

				id, err := uuid.Parse(value)

				if err != nil {
					NewError("error to parse account id", URLParse).Json(rw)
					return
				}

				ctx = context.WithValue(ctx, key, id)
			}

			h.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}

func ErrorHandler(next ErrorHandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		err := next(rw, r)

		if err == nil {
			return
		}

		if httpError, ok := err.(*Error); ok {
			httpError.Json(rw)
			return
		}

		NewError(err.Error(), UnkownError).Json(rw)
	}
}

func Auth(config *config.Config) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			authorization := r.Header.Get("Authorization")

			accountID, err := token.Auth(authorization, config.Secret)

			if err != nil {
				NewError("account unauthorized", AccountUnauthorized).Json(rw)
				return
			}

			ctx := context.WithValue(r.Context(), "account_id", accountID)

			h.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}
