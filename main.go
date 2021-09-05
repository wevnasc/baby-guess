package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/wevnasc/baby-guess/accounts"
	"github.com/wevnasc/baby-guess/db"
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
	logger := log.New(os.Stdout, "HTTP: ", log.LstdFlags|log.Lshortfile)

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

	h := accounts.NewHandler(logger, database)
	mux := http.NewServeMux()

	h.SetupRoutes(mux)
	srv := server.New(mux, ServerAddr)

	fmt.Printf("starting server on %s\n", ServerAddr)

	if err := srv.ListenAndServe(); err != nil {
		return fmt.Errorf("server failed to start: %v", err)
	}

	return nil
}
