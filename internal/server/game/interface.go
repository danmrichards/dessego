package game

import (
	"time"

	"github.com/danmrichards/dessego/internal/service/character"
	"github.com/danmrichards/dessego/internal/service/ghost"
	"github.com/danmrichards/dessego/internal/service/msg"
	"github.com/danmrichards/dessego/internal/service/replay"
	"github.com/danmrichards/dessego/internal/service/sos"
)

// Characters is the interface that wraps methods that types must implement to
// be used as a service for managing characters.
type Characters interface {
	// EnsureCreate creates a character with the given ID and index.
	//
	// If a character with the given ID and index already exists, no error will
	// be returned.
	EnsureCreate(id string) error

	// WorldTendency returns a maximum of n world tendency entries.
	WorldTendency(n int) ([]character.WorldTendency, error)

	// SetTendency sets the world tendency for the character with the given ID.
	SetTendency(id string, wt character.WorldTendency) error

	// Stats returns a map of statistics for the given character.
	Stats(id string) (*character.Stats, error)

	// MsgRating returns the message rating for the character with the given ID.
	MsgRating(id string) (int, error)

	// UpdateMsgRating updates the message rating for the character with the
	// given ID.
	UpdateMsgRating(id string) error
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

	// Add adds a new message.
	Add(bm msg.BloodMsg) error

	// Delete deletes the message with the given ID.
	Delete(id int) error

	// Get returns the message with the given ID.
	Get(id int) (*msg.BloodMsg, error)

	// UpdateRating updates the rating for the message with the given ID.
	UpdateRating(id int) error
}

// Ghosts is the interface that wraps methods that types must implement to be
// used as a service for managing ghosts.
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

// Replays is the interface that wraps methods that types must implement to be
// used as a service for managing replays.
type Replays interface {
	// List returns n replays for the given block ID and legacy type.
	List(blockID int32, n int, legacy replay.LegacyType) ([]replay.Replay, error)

	// Get returns a given replay.
	Get(id uint32) (*replay.Replay, error)

	// Add adds a new replay.
	Add(r *replay.Replay) error
}

// SOS is the interface that wraps methods that types must implement to be used
// as a service for managing SOS data.
type SOS interface {
	// Get returns n SOS entries, from the requested list in the given block.
	Get(blockID int32, n int) []*sos.SOS

	// Add adds a new SOS.
	Add(s *sos.SOS)

	// Delete deletes the SOS for a given character.
	Delete(characterID string)

	// Check checks for a matching player to fulfill an SOS and returns the ID
	// of the room for the match.
	Check(characterID string) string
}
