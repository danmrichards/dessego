package msg

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"os"

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

// NewSQLiteService returns an initialised SQLite character service.
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

// Character returns n messages for the given character and within the given
// block ID.
func (s *SQLiteService) Character(playerID string, blockID int32, n int) (bms []BloodMsg, err error) {
	bms = make([]BloodMsg, 0, n)

	var stmt *sql.Stmt
	stmt, err = s.db.Prepare(
		`SELECT *
		FROM message 
		WHERE character_id = ?
		AND block_id = ?
		AND legacy = ?
		ORDER BY random()
		LIMIT ?`,
	)
	if err != nil {
		return nil, fmt.Errorf("prepare select: %w", err)
	}

	var rows *sql.Rows
	rows, err = stmt.Query(playerID, blockID, 0, n)
	if err != nil {
		return nil, fmt.Errorf("query rows: %w", err)
	}

	for rows.Next() {
		var bm BloodMsg
		if err = rows.Scan(
			&bm.ID,
			&bm.CharacterID,
			&bm.BlockID,
			&bm.PosX,
			&bm.PosY,
			&bm.PosZ,
			&bm.AngX,
			&bm.AngY,
			&bm.AngZ,
			&bm.MsgID,
			&bm.MainMsgID,
			&bm.AddMsgCateID,
			&bm.Rating,
			&bm.Legacy,
		); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}

		bms = append(bms, bm)
	}

	return bms, nil
}

// NonCharacter returns n messages for anyone other than the given character and
// within the given block ID.
func (s *SQLiteService) NonCharacter(playerID string, blockID int32, n int) (bms []BloodMsg, err error) {
	bms = make([]BloodMsg, 0, n)

	var stmt *sql.Stmt
	stmt, err = s.db.Prepare(
		`SELECT *
		FROM message 
		WHERE character_id != ?
		AND block_id = ?
		AND legacy = ?
		ORDER BY random()
		LIMIT ?`,
	)
	if err != nil {
		return nil, fmt.Errorf("prepare select: %w", err)
	}

	var rows *sql.Rows
	rows, err = stmt.Query(playerID, blockID, 0, n)
	if err != nil {
		return nil, fmt.Errorf("query rows: %w", err)
	}

	for rows.Next() {
		var bm BloodMsg
		if err = rows.Scan(
			&bm.ID,
			&bm.CharacterID,
			&bm.BlockID,
			&bm.PosX,
			&bm.PosY,
			&bm.PosZ,
			&bm.AngX,
			&bm.AngY,
			&bm.AngZ,
			&bm.MsgID,
			&bm.MainMsgID,
			&bm.AddMsgCateID,
			&bm.Rating,
			&bm.Legacy,
		); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}

		bms = append(bms, bm)
	}

	return bms, nil
}

// Legacy returns n legacy messages within the given block ID.
func (s *SQLiteService) Legacy(blockID int32, n int) (bms []BloodMsg, err error) {
	bms = make([]BloodMsg, 0, n)

	var stmt *sql.Stmt
	stmt, err = s.db.Prepare(
		`SELECT *
		FROM message 
		WHERE block_id = ?
		AND legacy = ?
		ORDER BY random()
		LIMIT ?`,
	)
	if err != nil {
		return nil, fmt.Errorf("prepare select: %w", err)
	}

	var rows *sql.Rows
	rows, err = stmt.Query(blockID, 1, n)
	if err != nil {
		return nil, fmt.Errorf("query rows: %w", err)
	}

	for rows.Next() {
		var bm BloodMsg
		if err = rows.Scan(
			&bm.ID,
			&bm.CharacterID,
			&bm.BlockID,
			&bm.PosX,
			&bm.PosY,
			&bm.PosZ,
			&bm.AngX,
			&bm.AngY,
			&bm.AngZ,
			&bm.MsgID,
			&bm.MainMsgID,
			&bm.AddMsgCateID,
			&bm.Rating,
			&bm.Legacy,
		); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}

		bms = append(bms, bm)
	}

	return bms, nil
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
	ddl, err := ioutil.ReadFile("internal/service/msg/ddl.sql")
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
	s.l.Debug().Msg("seeding legacy messages")

	f, err := os.Open("internal/service/msg/legacymessages.bin")
	if err != nil {
		return err
	}
	defer f.Close()

	r := bufio.NewReader(f)

	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("db tx: %w", err)
	}

	buf := make([]byte, 512)
	for {
		n, err := io.ReadFull(r, buf[:4])
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		// Message length.
		ml := int(binary.LittleEndian.Uint32(buf[:n]))

		// Message body.
		n, err = io.ReadFull(r, buf[:ml])
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		msg, err := NewBloodMsgFromBytes(buf[:n])
		if err != nil {
			return fmt.Errorf("parse blood message: %w", err)
		}

		if err = s.saveMsg(tx, msg); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *SQLiteService) saveMsg(tx sqlPreparer, msg *BloodMsg) error {
	stmt, err := tx.Prepare(
		`INSERT OR IGNORE INTO message (
			id,
			character_id,
			block_id,
			posx,
			posy,
			posz,
			angx,
			angy,
			angz,
			msg_id,
			main_msg_id,
			add_msg_cate_id,
			rating,
			legacy
		) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
	)
	if err != nil {
		return fmt.Errorf("prepare insert: %w", err)
	}

	if _, err = stmt.Exec(
		msg.ID,
		msg.CharacterID,
		msg.BlockID,
		msg.PosX,
		msg.PosY,
		msg.PosZ,
		msg.AngX,
		msg.AngY,
		msg.AngZ,
		msg.MsgID,
		msg.MainMsgID,
		msg.AddMsgCateID,
		msg.Rating,
		msg.Legacy,
	); err != nil {
		return fmt.Errorf("create character: %w", err)
	}

	return nil
}
