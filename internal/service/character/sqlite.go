package character

import (
	"database/sql"
	"fmt"
	"io/ioutil"
)

// SQLiteService is a character service backed by a SQLite database.
type SQLiteService struct {
	db *sql.DB
}

// NewSQLiteService returns an initialised SQLite character service.
func NewSQLiteService(db *sql.DB) (*SQLiteService, error) {
	s := &SQLiteService{
		db: db,
	}

	if err := s.init(); err != nil {
		return nil, fmt.Errorf("initialise: %w", err)
	}

	return s, nil
}

// EnsureCreate creates a character with the given ID and index.
//
// If a character with the given ID and index already exists, no error will
// be returned.
func (s *SQLiteService) EnsureCreate(id string) error {
	stmt, err := s.db.Prepare(
		`INSERT OR IGNORE INTO character (id) VALUES (?)`,
	)
	if err != nil {
		return fmt.Errorf("prepare insert: %w", err)
	}
	if _, err = stmt.Exec(id); err != nil {
		return fmt.Errorf("create character: %w", err)
	}

	return nil
}

// WorldTendency returns a maximum of n world tendency entries.
func (s *SQLiteService) WorldTendency(n int) (wts []WorldTendency, err error) {
	wts = make([]WorldTendency, 0, n)

	var stmt *sql.Stmt
	stmt, err = s.db.Prepare(
		`SELECT area_1, wb_1, lr_1,
       	area_2, wb_2, lr_2,
		area_3, wb_3, lr_3,
		area_4, wb_4, lr_4,
		area_5, wb_5, lr_5,
		area_6, wb_6, lr_6,
		area_7, wb_7, lr_7
		FROM world_tendency
		ORDER BY id DESC
		LIMIT ?`,
	)
	if err != nil {
		return nil, fmt.Errorf("prepare select: %w", err)
	}

	var rows *sql.Rows
	rows, err = stmt.Query(n)
	if err != nil {
		return nil, fmt.Errorf("query rows: %w", err)
	}

	for rows.Next() {
		var wt WorldTendency
		if err = rows.Scan(
			&wt.Area1, &wt.WB1, &wt.LR1,
			&wt.Area2, &wt.WB2, &wt.LR2,
			&wt.Area3, &wt.WB3, &wt.LR3,
			&wt.Area4, &wt.WB4, &wt.LR4,
			&wt.Area5, &wt.WB5, &wt.LR5,
			&wt.Area6, &wt.WB6, &wt.LR6,
			&wt.Area7, &wt.WB7, &wt.LR7,
		); err != nil {
			return nil, fmt.Errorf("query row: %w", err)
		}

		wts = append(wts, wt)
	}

	return wts, nil
}

// SetTendency sets the world tendency for the character with the given ID.
func (s *SQLiteService) SetTendency(id string, wt WorldTendency) error {
	stmt, err := s.db.Prepare(
		`INSERT INTO world_tendency (
            character_id,
			area_1, wb_1, lr_1, 
		    area_2, wb_2, lr_2,
		    area_3, wb_3, lr_3,
			area_4, wb_4, lr_4,
			area_5, wb_5, lr_5,
			area_6, wb_6, lr_6,
			area_7, wb_7, lr_7
		) VALUES (
		    ?,
			?, ?, ?,
			?, ?, ?,
			?, ?, ?,
			?, ?, ?,
			?, ?, ?,
			?, ?, ?,
			?, ?, ?
		)`,
	)
	if err != nil {
		return fmt.Errorf("prepare update: %w", err)
	}

	if _, err := stmt.Exec(
		id,
		wt.Area1, wt.WB1, wt.LR1,
		wt.Area2, wt.WB2, wt.LR2,
		wt.Area3, wt.WB3, wt.LR3,
		wt.Area4, wt.WB4, wt.LR4,
		wt.Area5, wt.WB5, wt.LR5,
		wt.Area6, wt.WB6, wt.LR6,
		wt.Area7, wt.WB7, wt.LR7,
	); err != nil {
		return fmt.Errorf("update row: %w", err)
	}

	return nil
}

// Stats returns a map of statistics for the given character.
func (s *SQLiteService) Stats(id string) (*Stats, error) {
	stmt, err := s.db.Prepare(
		`SELECT grade_s, grade_a, grade_b, grade_c, grade_d, sessions
		FROM character
		WHERE id = ?`,
	)
	if err != nil {
		return nil, fmt.Errorf("prepare select: %w", err)
	}

	st := &Stats{}
	if err = stmt.QueryRow(id).Scan(
		&st.GradeS, &st.GradeA, &st.GradeB, &st.GradeC, &st.GradeD, &st.Sessions,
	); err != nil {
		return nil, fmt.Errorf("query row: %w", err)
	}

	return st, nil
}

// MsgRating returns the message rating for the character with the given ID.
func (s *SQLiteService) MsgRating(id string) (mr int, err error) {
	var stmt *sql.Stmt
	stmt, err = s.db.Prepare(
		`SELECT msg_rating FROM character WHERE id = ?`,
	)
	if err != nil {
		return 0, fmt.Errorf("prepare select: %w", err)
	}

	if err = stmt.QueryRow(id).Scan(&mr); err != nil {
		return 0, fmt.Errorf("query row: %w", err)
	}

	return mr, nil
}

// UpdateMsgRating updates the message rating for the character with the
// given ID.
func (s *SQLiteService) UpdateMsgRating(id string) error {
	stmt, err := s.db.Prepare(
		`UPDATE character SET msg_rating = msg_rating + 1 WHERE id = ?`,
	)
	if err != nil {
		return fmt.Errorf("prepare query: %w", err)
	}

	if _, err = stmt.Exec(id); err != nil {
		return fmt.Errorf("update message rating: %w", err)
	}

	return nil
}

// InitMultiplayer initialises a multiplayer session for the given
// characterID.
func (s *SQLiteService) InitMultiplayer(id string) error {
	stmt, err := s.db.Prepare(
		`UPDATE character SET sessions = sessions + 1 WHERE id = ?`,
	)
	if err != nil {
		return fmt.Errorf("prepare query: %w", err)
	}

	if _, err = stmt.Exec(id); err != nil {
		return fmt.Errorf("update message rating: %w", err)
	}

	return nil
}

// UpdatePlayerGrade updates the given player with the given grade.
func (s *SQLiteService) UpdatePlayerGrade(id string, grade MultiplayerGrade) error {
	stmt, err := s.db.Prepare(
		`UPDATE character SET ? = ? + 1 WHERE id = ?`,
	)
	if err != nil {
		return fmt.Errorf("prepare query: %w", err)
	}

	if _, err = stmt.Exec(grade, grade, id); err != nil {
		return fmt.Errorf("update message rating: %w", err)
	}

	return nil
}

// init initialises the database tables required by this service.
func (s *SQLiteService) init() error {
	for _, t := range []string{"character", "world_tendency"} {
		if err := s.initTable(t); err != nil {
			return err
		}
	}

	return nil
}

func (s *SQLiteService) initTable(table string) error {
	ddl, err := ioutil.ReadFile("internal/service/character/" + table + ".sql")
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
