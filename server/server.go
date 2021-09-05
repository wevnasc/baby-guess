package server

import (
	"net/http"
	"time"
)

func New(mux *http.ServeMux, serverAddr string) *http.Server {
	return &http.Server{
		Addr:         serverAddr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}
}
