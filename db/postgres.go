package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/wevnasc/baby-guess/config"
)

func New(c *config.Config) (*Store, error) {
	url := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost,
		c.DBPort,
		c.DBUser,
		c.DBPass,
		c.DBName,
	)

	db, err := sql.Open("postgres", url)

	if err != nil {
		return nil, fmt.Errorf("not was possible to create a coonection %v", err)
	}

	err = db.Ping()

	if err != nil {
		return nil, fmt.Errorf("not was possible to connect with the database %v", err)
	}

	return &Store{db}, nil
}
