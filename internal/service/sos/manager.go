package sos

import (
	"sync"
	"time"

	"github.com/rs/zerolog"
)

var monkBlockIDs = []int32{40070, 40071, 40072, 40073, 40074, 40170, 40171, 40172, 40270}

// Manager is an in-memory SOS management service.
type Manager struct {
	index int32

	// Active SOS requests.
	active map[string]*SOS

	// Pending players (character ID -> room ID).
	pending map[string]string

	// Pending monks (character ID -> room ID).
	monks map[string]string

	l zerolog.Logger

	sync.Mutex
}

// NewManager returns an instantiated SOS manager.
func NewManager(l zerolog.Logger) *Manager {
	return &Manager{
		active: make(map[string]*SOS),
		l:      l,
	}
}

// List returns n SOS entries in the given block.
func (m *Manager) List(blockID int32, n int) []*SOS {
	m.Lock()
	defer m.Unlock()

	var (
		found int
		sos   = make([]*SOS, 0, n)
	)
	for cid, a := range m.active {
		if a.Updated.Add(maxSOSAge).Before(time.Now()) {
			m.l.Info().Msgf("deleted SOS %d due to inactivity", a.ID)
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

// Add adds a new SOS and returns its ID.
func (m *Manager) Add(s *SOS) {
	m.Lock()
	defer m.Unlock()

	m.index++
	s.ID = m.index
	m.active[s.CharacterID] = s
}

// Delete deletes the SOS for a given character.
func (m *Manager) Delete(characterID string) {
	m.Lock()
	defer m.Unlock()

	delete(m.active, characterID)
}

// Check checks for a matching player to fulfill an SOS and returns the ID
// of the room for the match.
func (m *Manager) Check(characterID string) string {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.active[characterID]; ok {
		m.active[characterID].Updated = time.Now()
	}

	if len(m.pending) > 0 || len(m.monks) > 0 {
		m.l.Debug().Msgf("potential connect data %v %v", m.pending, m.monks)
	}

	if rid, ok := m.monks[characterID]; ok {
		m.l.Info().Msgf("summoning for monk player %q", characterID)
		delete(m.monks, characterID)
		return rid
	}

	if rid, ok := m.pending[characterID]; ok {
		m.l.Info().Msgf("connecting player %q", characterID)
		delete(m.pending, characterID)
		return rid
	}

	return ""
}

// Summon returns true if the given SOS ID was able to be summoned.
func (m *Manager) Summon(id int32, room string) bool {
	m.Lock()
	defer m.Unlock()

	for _, a := range m.active {
		if a.ID == id {
			m.pending[a.CharacterID] = room
			m.l.Info().Msgf(
				"added pending summon for character %q in room %q",
				a.CharacterID,
				room,
			)
			return true
		}
	}

	return false
}

// Monk returns true if a monk was able to be summoned to the given room.
func (m *Manager) Monk(room string) bool {
	m.Lock()
	defer m.Unlock()

	for _, a := range m.active {
		if monkBlock(a.BlockID) {
			m.monks[a.CharacterID] = room
			m.l.Info().Msgf("added pending request for monk in room %q", room)
			return true
		}
	}

	return false
}

func monkBlock(id int32) bool {
	for _, mb := range monkBlockIDs {
		if id == mb {
			return true
		}
	}

	return false
}
