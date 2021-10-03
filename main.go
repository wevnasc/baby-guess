package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gorilla/mux"
	"github.com/wevnasc/baby-guess/auth"
	"github.com/wevnasc/baby-guess/config"
	"github.com/wevnasc/baby-guess/db"
	"github.com/wevnasc/baby-guess/email"
	"github.com/wevnasc/baby-guess/server"
	"github.com/wevnasc/baby-guess/tables"
)

var Local = os.Getenv("LOCAL")

func main() {
	err := run()

	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	config := config.New(Local)
	store, err := db.New(config)

	if err != nil {
		return err
	}

	defer store.Close()

	err = db.RunMigrations(store)

	if err != nil {
		return err
	}

	var emailClient email.Client

	if Local == "true" {
		emailClient = &email.DebugClient{}
	} else {
		emailClient = email.NewSmtpClient(config)
	}

	mux := mux.NewRouter()
	mux.Use(server.Headers)

	auth.NewHandler(store, config, emailClient).SetupRoutes(mux)
	tables.NewHandler(store, config, emailClient).SetupRoutes(mux)

	srv := server.New(mux, config.ServerAddr)

	if err := srv.ListenAndServe(); err != nil {
		return fmt.Errorf("server failed to start: %v", err)
	}

	return nil
}
