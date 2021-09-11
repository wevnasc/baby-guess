package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gorilla/mux"
	"github.com/wevnasc/baby-guess/accounts"
	"github.com/wevnasc/baby-guess/db"
	"github.com/wevnasc/baby-guess/middleware"
	"github.com/wevnasc/baby-guess/server"
	"github.com/wevnasc/baby-guess/tables"
)

var (
	ServerAddr = os.Getenv("HTTP_SERVER_ADDR")
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

	defer store.Close()

	if err != nil {
		return err
	}

	mux := mux.NewRouter()
	mux.Use(middleware.Headers)

	accounts.NewHandler(store).SetupRoutes(mux)
	tables.NewHandler(store).SetupRoutes(mux)

	srv := server.New(mux, ServerAddr)

	if err := srv.ListenAndServe(); err != nil {
		return fmt.Errorf("server failed to start: %v", err)
	}

	return nil
}
