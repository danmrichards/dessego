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
func (s *SQLiteService) EnsureCreate(id string) error {
	stmt, err := s.db.Prepare(
		`INSERT OR IGNORE INTO player (id) VALUES (?)`,
	)
	if err != nil {
		return fmt.Errorf("prepare insert: %w", err)
	}

	if _, err = stmt.Exec(id); err != nil {
		return fmt.Errorf("create player: %w", err)
	}

	return nil
}

// DesiredTendency returns the desired tendency for the player with the
// given ID.
func (s *SQLiteService) DesiredTendency(id string) (dt int, err error) {
	var stmt *sql.Stmt
	stmt, err = s.db.Prepare(
		`SELECT desired_tendency FROM player WHERE id = ?`,
	)
	if err != nil {
		return 0, fmt.Errorf("prepare select: %w", err)
	}

	if err = stmt.QueryRow(id).Scan(&dt); err != nil {
		return 0, fmt.Errorf("query row: %w", err)
	}

	return dt, nil
}

// Stats returns a map of statistics for the given player.
func (s *SQLiteService) Stats(id string) (map[string]int, error) {
	var stmt *sql.Stmt
	stmt, err := s.db.Prepare(
		`SELECT grade_s, grade_a, grade_b, grade_c, grade_d, sessions
		FROM player 
		WHERE id = ?`,
	)
	if err != nil {
		return nil, fmt.Errorf("prepare select: %w", err)
	}

	var gs, ga, gb, gc, gd, sess int
	if err = stmt.QueryRow(id).Scan(&gs, &ga, &gb, &gc, &gd, &sess); err != nil {
		return nil, fmt.Errorf("query row: %w", err)
	}

	return map[string]int{
		"grade_s":  gs,
		"grade_a":  ga,
		"grade_b":  gb,
		"grade_c":  gc,
		"grade_d":  gd,
		"sessions": sess,
	}, nil
}

// MsgRating returns the message rating for the player with the given ID.
func (s *SQLiteService) MsgRating(id string) (mr int, err error) {
	var stmt *sql.Stmt
	stmt, err = s.db.Prepare(
		`SELECT msg_rating FROM player WHERE id = ?`,
	)
	if err != nil {
		return 0, fmt.Errorf("prepare select: %w", err)
	}

	if err = stmt.QueryRow(id).Scan(&mr); err != nil {
		return 0, fmt.Errorf("query row: %w", err)
	}

	return mr, nil
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
