package ghost

import (
	"sync"
	"time"

	"github.com/rs/zerolog"
)

// Memory is an in-memory ghost manager.
type Memory struct {
	ghosts   map[string]*Ghost
	ghostAge []string
	l        zerolog.Logger

	sync.Mutex
}

// NewMemory returns a new ghost manager.
//
// It stores both a map of ghosts by character ID and also an ordered set of
// ghosts by timestamp (oldest first).
func NewMemory(l zerolog.Logger) *Memory {
	return &Memory{
		ghosts: make(map[string]*Ghost),
		l:      l,
	}
}

// Get returns n ghosts for anyone other than the given character and
// within the given block ID.
func (m *Memory) Get(characterID string, blockID int32, n int) []*Ghost {
	m.Lock()
	defer m.Unlock()

	g := make([]*Ghost, 0, n)
	var i int
	for _, mg := range m.ghosts {
		if i == n {
			break
		} else if mg.CharacterID == characterID || mg.BlockID != blockID {
			continue
		}

		g = append(g, mg)

		i++
	}

	return g
}

// ClearBefore clears any ghosts before the given time.
func (m *Memory) ClearBefore(t time.Time) {
	m.Lock()
	defer m.Unlock()

	for i, g := range m.ghostAge {
		if !m.ghosts[g].timestamp.Before(t) {
			continue
		}

		m.l.Debug().Msgf(
			"deleting stale ghost for character: %q", m.ghosts[g].CharacterID,
		)

		// Clear the ghost from the map and ordered set.
		delete(m.ghosts, g)
		m.ghostAge = append(m.ghostAge[:i], m.ghostAge[i+1:]...)
	}
}

// Character returns a ghost, if it exists, for the given character.
func (m *Memory) Character(characterID string) (*Ghost, error) {
	m.Lock()
	defer m.Unlock()

	g, ok := m.ghosts[characterID]
	if !ok {
		return nil, CharacterGhostNotFoundError(characterID)
	}

	return g, nil
}

// Set sets the ghost for the given character.
func (m *Memory) Set(characterID string, g *Ghost) {
	m.Lock()
	defer m.Unlock()

	m.ghosts[characterID] = g
}
