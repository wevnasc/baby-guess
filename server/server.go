package server

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func New(mux *mux.Router, port string) *http.Server {
	return &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}
}
