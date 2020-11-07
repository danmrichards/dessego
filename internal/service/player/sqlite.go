package player

import (
	"database/sql"
	"fmt"
	"io/ioutil"
)

// SQLiteService is a player service backed by a SQLite database.
type SQLiteService struct {
	db *sql.DB
}

// NewSQLiteService returns an initialised SQLite player service.
func NewSQLiteService(db *sql.DB) (*SQLiteService, error) {
	s := &SQLiteService{
		db: db,
	}

	if err := s.init(); err != nil {
		return nil, fmt.Errorf("initialise: %w", err)
	}

	return s, nil
}

// EnsureCreate creates a player with the given ID and index.
//
// If a player with the given ID and index already exists, no error will
// be returned.
func (s *SQLiteService) EnsureCreate(id string, index int) error {
	stmt, err := s.db.Prepare(
		`INSERT OR IGNORE INTO player (id, idx) VALUES (?, ?)`,
	)
	if err != nil {
		return fmt.Errorf("prepare insert: %w", err)
	}

	if _, err = stmt.Exec(id, index); err != nil {
		return fmt.Errorf("create player: %w", err)
	}

	return nil
}

// init initialises the database tables required by this service.
func (s *SQLiteService) init() error {
	ddl, err := ioutil.ReadFile("internal/service/player/ddl.sql")
	if err != nil {
		return fmt.Errorf("read DDL: %w", err)
	}

	var stmt *sql.Stmt
	stmt, err = s.db.Prepare(string(ddl))
	if err != nil {
		return fmt.Errorf("prepare DDL: %w", err)
	}

	if _, err = stmt.Exec(); err != nil {
		return fmt.Errorf("init table: %w", err)
	}

	return nil
}
