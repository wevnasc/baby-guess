package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func New(c *Connection) (*sql.DB, error) {
	url := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.Database,
	)

	db, err := sql.Open("postgres", url)

	if err != nil {
		return nil, fmt.Errorf("not was possible to create a coonection %v", err)
	}

	err = db.Ping()

	if err != nil {
		return nil, fmt.Errorf("not was possible to connect with the database %v", err)
	}

	return db, nil
}
