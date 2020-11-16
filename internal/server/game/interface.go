package game

import (
	"github.com/danmrichards/dessego/internal/service/msg"
	"github.com/danmrichards/dessego/internal/service/player"
)

// Players is the interface that wraps methods that types must implement to be
// used as a service for managing players.
type Players interface {
	// EnsureCreate creates a player with the given ID and index.
	//
	// If a player with the given ID and index already exists, no error will
	// be returned.
	EnsureCreate(id string) error

	// DesiredTendency returns the desired tendency for the player with the
	// given ID.
	DesiredTendency(id string) (int, error)

	// Stats returns a map of statistics for the given player.
	Stats(id string) (*player.Stats, error)

	// MsgRating returns the message rating for the player with the given ID.
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
	// Get returns n messages for the given player and block ID.
	Get(playerID string, blockID, n int) ([]msg.BloodMsg, error)
}
