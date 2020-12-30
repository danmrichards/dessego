package replay

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

type sqlPreparer interface {
	Prepare(query string) (*sql.Stmt, error)
}

// Option is a functional option that configures the SQLite service.
type Option func(*SQLiteService)

// SQLiteService is a msg service backed by a SQLite database.
type SQLiteService struct {
	db   *sql.DB
	l    zerolog.Logger
	seed bool
}

// Seed configures the service to seed the database on startup.
func Seed() Option {
	return func(s *SQLiteService) {
		s.seed = true
	}
}

// NewSQLiteService returns an initialised SQLite replays service.
func NewSQLiteService(db *sql.DB, l zerolog.Logger, opts ...Option) (*SQLiteService, error) {
	s := &SQLiteService{
		db: db,
		l:  l,
	}

	for _, o := range opts {
		o(s)
	}

	if err := s.init(); err != nil {
		return nil, fmt.Errorf("initialise: %w", err)
	}

	return s, nil
}

// List returns n replays for the given block ID and legacy type.
func (s *SQLiteService) List(blockID int32, n int, legacy LegacyType) (rs []Replay, err error) {
	rs = make([]Replay, 0, n)

	var stmt *sql.Stmt
	stmt, err = s.db.Prepare(
		`SELECT *
		FROM replay 
		WHERE block_id = ?
		AND legacy = ?
		ORDER BY random()
		LIMIT ?`,
	)
	if err != nil {
		return nil, fmt.Errorf("prepare select: %w", err)
	}

	var rows *sql.Rows
	rows, err = stmt.Query(blockID, legacy, n)
	if err != nil {
		return nil, fmt.Errorf("query rows: %w", err)
	}

	for rows.Next() {
		var r Replay
		if err = rows.Scan(
			&r.ID,
			&r.CharacterID,
			&r.BlockID,
			&r.PosX,
			&r.PosY,
			&r.PosZ,
			&r.AngX,
			&r.AngY,
			&r.AngZ,
			&r.MsgID,
			&r.MainMsgID,
			&r.AddMsgCateID,
			&r.Data,
			&r.Legacy,
		); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}

		rs = append(rs, r)
	}

	return rs, nil
}

// Get returns a given replay.
func (s *SQLiteService) Get(id uint32) (r *Replay, err error) {
	var stmt *sql.Stmt
	stmt, err = s.db.Prepare(
		`SELECT *
		FROM replay 
		WHERE id = ?`,
	)
	if err != nil {
		return nil, fmt.Errorf("prepare select: %w", err)
	}

	r = &Replay{}
	if err = stmt.QueryRow(id).Scan(
		&r.ID,
		&r.CharacterID,
		&r.BlockID,
		&r.PosX,
		&r.PosY,
		&r.PosZ,
		&r.AngX,
		&r.AngY,
		&r.AngZ,
		&r.MsgID,
		&r.MainMsgID,
		&r.AddMsgCateID,
		&r.Data,
		&r.Legacy,
	); err != nil {
		return nil, fmt.Errorf("query row: %w", err)
	}

	return r, nil
}

// Add adds a new replay.
func (s *SQLiteService) Add(r *Replay) error {
	return s.saveReplay(s.db, r)
}

// init initialises the database tables required by this service.
func (s *SQLiteService) init() error {
	if err := s.initTable(); err != nil {
		return err
	}

	if s.seed {
		return s.doSeed()
	}
	return nil
}

// initTable creates the database tables required by this service.
func (s *SQLiteService) initTable() error {
	ddl, err := ioutil.ReadFile("internal/service/replay/ddl.sql")
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

// doSeed seeds the database tables required by this service.
func (s *SQLiteService) doSeed() error {
	s.l.Debug().Msg("seeding legacy replays")

	f, err := os.Open("internal/service/replay/legacyreplays.bin")
	if err != nil {
		return err
	}
	defer f.Close()

	r := bufio.NewReader(f)

	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("db tx: %w", err)
	}

	buf := make([]byte, 2048)
	for {
		n, err := io.ReadFull(r, buf[:4])
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		// Replay length.
		ml := int(binary.LittleEndian.Uint32(buf[:n]))

		// Replay body.
		n, err = io.ReadFull(r, buf[:ml])
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		rp, err := NewReplayFromBytes(buf[:n])
		if err != nil {
			return fmt.Errorf("parse replay: %w", err)
		}

		if err = s.saveReplay(tx, rp); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *SQLiteService) saveReplay(tx sqlPreparer, r *Replay) error {
	cols := []string{
		"character_id",
		"block_id",
		"posx",
		"posy",
		"posz",
		"angx",
		"angy",
		"angz",
		"msg_id",
		"main_msg_id",
		"add_msg_cate_id",
		"data",
		"legacy",
	}
	vals := []interface{}{
		r.CharacterID,
		r.BlockID,
		r.PosX,
		r.PosY,
		r.PosZ,
		r.AngX,
		r.AngY,
		r.AngZ,
		r.MsgID,
		r.MainMsgID,
		r.AddMsgCateID,
		r.Data,
		r.Legacy,
	}

	// Leverage the auto increment for replays with no ID.
	if r.ID > 0 {
		cols = append([]string{"id"}, cols...)
		vals = append([]interface{}{r.ID}, vals...)
	}

	var qb strings.Builder
	qb.WriteString("INSERT OR IGNORE INTO replay (")
	for i, c := range cols {
		qb.WriteString(c)
		if i < len(cols)-1 {
			qb.WriteString(",")
		}
	}

	qb.WriteString(") VALUES (")
	for i := 0; i < len(vals); i++ {
		qb.WriteString("?")
		if i < len(vals)-1 {
			qb.WriteString(",")
		}
	}
	qb.WriteString(")")

	stmt, err := tx.Prepare(qb.String())
	if err != nil {
		return fmt.Errorf("prepare insert: %w", err)
	}

	if _, err = stmt.Exec(vals...); err != nil {
		return fmt.Errorf("save replay: %w", err)
	}

	return nil
}
