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
	database, err := db.New(&db.Connection{
		Host:     "localhost",
		User:     "postgres",
		Password: "postgres",
		Port:     "5432",
		Database: "baby_guess",
	})

	defer database.Close()

	if err != nil {
		return err
	}

	h := accounts.NewHandler(database)
	mux := mux.NewRouter()

	mux.Use(middleware.Headers)
	h.SetupRoutes(mux)
	srv := server.New(mux, ServerAddr)

	fmt.Printf("starting server on localhost:%s\n", ServerAddr)

	if err := srv.ListenAndServe(); err != nil {
		return fmt.Errorf("server failed to start: %v", err)
	}

	return nil
}
