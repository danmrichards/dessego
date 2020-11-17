package ghost

import "time"

// Ghost represents a replay of a Demon's Souls characters actions.
type Ghost struct {
	BlockID     int32
	CharacterID string
	ReplayData  []byte
	timestamp   time.Time
}

// NewGhost returns an instantiated ghost.
func NewGhost(blockID int32, characterID string, replayData []byte) *Ghost {
	return &Ghost{
		BlockID:     blockID,
		CharacterID: characterID,
		ReplayData:  replayData,
		timestamp:   time.Now(),
	}
}
