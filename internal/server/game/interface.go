package game

import (
	"time"

	"github.com/danmrichards/dessego/internal/service/character"
	"github.com/danmrichards/dessego/internal/service/ghost"
	"github.com/danmrichards/dessego/internal/service/msg"
)

// Characters is the interface that wraps methods that types must implement to
// be used as a service for managing characters.
type Characters interface {
	// EnsureCreate creates a character with the given ID and index.
	//
	// If a character with the given ID and index already exists, no error will
	// be returned.
	EnsureCreate(id string) error

	// DesiredTendency returns the desired tendency for the character with the
	// given ID.
	DesiredTendency(id string) (int, error)

	// Stats returns a map of statistics for the given character.
	Stats(id string) (*character.Stats, error)

	// MsgRating returns the message rating for the character with the given ID.
	MsgRating(id string) (int, error)
}

// State is the interface that wraps methods that types must implement to be
// used as a service for managing gamestate state.
type State interface {
	// Motd returns the messages of the day.
	Motd() []string

	// AddPlayer adds a new player to the current gamestate state.
	AddPlayer(ip, id string)

	// Player returns the ID of a player with the given IP address
	Player(ip string) (string, error)

	// TODO: Add Ghost

	// TODO: Get Ghost
}

// Messages is the interface that wraps methods that types must implement to be
// used as a service for managing messages.
type Messages interface {
	// Character returns n messages for the given character and within the given
	// block ID.
	Character(characterID string, blockID int32, n int) ([]msg.BloodMsg, error)

	// NonCharacter returns n messages for anyone other than the given character
	// and within the given block ID.
	NonCharacter(characterID string, blockID int32, n int) ([]msg.BloodMsg, error)

	// Legacy returns n legacy messages within the given block ID.
	Legacy(blockID int32, n int) ([]msg.BloodMsg, error)
}

// Ghosts is the interface that wraps methods that types must implement to be
// used as a service for managing messages.
type Ghosts interface {
	// Get returns n ghosts for anyone other than the given character and
	// within the given block ID.
	Get(characterID string, blockID int32, n int) []*ghost.Ghost

	// ClearBefore clears any ghosts before the given time.
	ClearBefore(t time.Time)

	// Character returns a ghost, if it exists, for the given character.
	Character(characterID string) (*ghost.Ghost, error)

	// Set sets the ghost for the given character.
	Set(characterID string, g *ghost.Ghost)
}
