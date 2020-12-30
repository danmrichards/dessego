package sos

import (
	"time"

	"github.com/rs/zerolog"
)

type Manager struct {
	active []SOS
	l      zerolog.Logger
}

// NewManager returns an instantiated SOS manager.
func NewManager(l zerolog.Logger) *Manager {
	return &Manager{
		l: l,
	}
}

// Get returns n SOS entries in the given block.
func (m *Manager) Get(blockID int32, n int) []SOS {
	var (
		found int
		sos   = make([]SOS, 0, n)
	)
	for i, a := range m.active {
		if a.Updated.Add(maxSOSAge).Before(time.Now()) {
			m.l.Info().Msgf("deleted SOS %q due to inactivity", a.ID)
			m.active = append(m.active[:i], m.active[i+1:]...)
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
