package sos

import (
	"sync"
	"time"

	"github.com/rs/zerolog"
)

type Manager struct {
	active map[string]SOS
	l      zerolog.Logger

	sync.Mutex
}

// NewManager returns an instantiated SOS manager.
func NewManager(l zerolog.Logger) *Manager {
	return &Manager{
		l: l,
	}
}

// Get returns n SOS entries in the given block.
func (m *Manager) Get(blockID int32, n int) []SOS {
	m.Lock()
	defer m.Unlock()

	var (
		found int
		sos   = make([]SOS, 0, n)
	)
	for cid, a := range m.active {
		if a.Updated.Add(maxSOSAge).Before(time.Now()) {
			m.l.Info().Msgf("deleted SOS %q due to inactivity", a.ID)
			delete(m.active, cid)
			continue
		} else if a.BlockID != blockID {
			// Only include SOS for the given block.
			continue
		}

		if found == n {
			break
		}
		sos = append(sos, a)
		found++
	}

	return sos
}

// Delete deletes the SOS for a given character.
func (m *Manager) Delete(characterID string) {
	m.Lock()
	defer m.Unlock()

	delete(m.active, characterID)
}
