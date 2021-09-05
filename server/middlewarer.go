package server

import (
	"log"
	"net/http"
	"time"
)

type Middleware struct {
	Log *log.Logger
}

func (m *Middleware) Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		next(rw, r)
		m.Log.Printf("request processed in %s\n", time.Now().Sub(startTime))
	}
}

func (m *Middleware) Headers(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-type", "application/json")
		next(rw, r)
	}
}

func (m *Middleware) Method(next http.HandlerFunc, allowed []string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		for _, method := range allowed {
			if method == r.Method {
				next(rw, r)
				return
			}
		}

		rw.WriteHeader(http.StatusMethodNotAllowed)
		return

	}
}

func (m *Middleware) Resource(next http.HandlerFunc, allowed []string) http.HandlerFunc {
	return m.Logger(m.Headers(m.Method(func(rw http.ResponseWriter, r *http.Request) {
		next(rw, r)
	}, allowed)))
}
