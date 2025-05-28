package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func ConnectionDB(user, password, dbname, host string, port int) (*sql.DB, error) {

	// str for connection to db in pq driver format
	connect := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%d sslmode=disable", user, password, dbname, host, port,
	)

	// Try connection to db, not real connection
	db, err := sql.Open("postgres", connect)
	if err != nil {
		return nil, fmt.Errorf("error open db: %w", err)
	}

	//Check real connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error ping db: %w", err)
	}

	return db, nil
}
