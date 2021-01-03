package database

import (
	"database/sql"
	"fmt"
	"os"

	// SQLite driver.
	_ "github.com/mattn/go-sqlite3"
)

// NewSQLite returns a new SQLite database.
func NewSQLite(path string) (*sql.DB, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if _, err = os.Create(path); err != nil {
			return nil, fmt.Errorf("create DB: %w", err)
		}
	}

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("open DB: %w", err)
	}

	return db, nil
}
