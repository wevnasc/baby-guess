package db

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func RunMigrations(store *Store) error {
	driver, err := postgres.WithInstance(store.DB, &postgres.Config{})

	if err != nil {
		return fmt.Errorf("error to get driver %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://db/migrations", "postgres", driver)

	if err != nil {
		return fmt.Errorf("error to create migration instance %v", err)
	}

	m.Up()

	return nil
}
