package ghost

import "time"

// Ghost represents a replay of a Demon's Souls characters actions.
type Ghost struct {
	BlockID     int32
	CharacterID string
	// ReplayData (not sure of type)
	timestamp time.Time
}

// NewGhost returns an instantiated ghost.
func NewGhost(blockID int32, characterID string) *Ghost {
	return &Ghost{
		BlockID:     blockID,
		CharacterID: characterID,
		timestamp:   time.Now(),
	}
}
