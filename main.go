package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gorilla/mux"
	"github.com/wevnasc/baby-guess/auth"
	"github.com/wevnasc/baby-guess/config"
	"github.com/wevnasc/baby-guess/db"
	"github.com/wevnasc/baby-guess/server"
	"github.com/wevnasc/baby-guess/tables"
)

var (
	ServerAddr = os.Getenv("HTTP_SERVER_ADDR")
	Secret     = os.Getenv("AUTH_SECRET")
)

func main() {
	err := run()

	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	store, err := db.New(&db.Connection{
		Host:     "localhost",
		User:     "postgres",
		Password: "postgres",
		Port:     "5432",
		Database: "baby_guess",
	})

	if err != nil {
		return err
	}

	defer store.Close()

	config := config.New(Secret)

	mux := mux.NewRouter()
	mux.Use(server.Headers)

	auth.NewHandler(store, config).SetupRoutes(mux)
	tables.NewHandler(store, config).SetupRoutes(mux)

	srv := server.New(mux, ServerAddr)

	if err := srv.ListenAndServe(); err != nil {
		return fmt.Errorf("server failed to start: %v", err)
	}

	return nil
}
