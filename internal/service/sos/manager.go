package sos

import (
	"sync"
	"time"

	"github.com/rs/zerolog"
)

// Manager is an in-memory SOS management service.
type Manager struct {
	index uint32

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

// Get returns n SOS entries in the given block.
func (m *Manager) Get(blockID int32, n int) []*SOS {
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
	//if characterID in self.monkPending[serverport]:
	//logging.info("Summoning for monk player %r" % characterID)
	//data = self.monkPending[serverport][characterID]
	//del self.monkPending[serverport][characterID]
	//
	//elif characterID in self.playerPending:
	//logging.info("Connecting player %r" % characterID)
	//data = self.playerPending[characterID]
	//del self.playerPending[characterID]
	//
	//else:
	//data = "\x00"
}
