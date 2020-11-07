package gamestate

import (
	"strconv"
	"sync"
)

// Memory is an in-memory gamestate state manager.
type Memory struct {
	sync.Mutex

	// players stores connected players by IP address.
	players map[string]string
}

// NewMemory returns a new in-memory gamestate state manager.
func NewMemory() *Memory {
	return &Memory{
		players: make(map[string]string),
	}
}

// Motd returns the messages of the day.
func (m *Memory) Motd() []string {
	motd := "Welcome to DeSSE Go\r\n"
	motd += "A server emulator for Demon's Souls implemented in Go\r\n"
	motd += "Source code:\r\n"
	motd += "https://github.com/danmrichards/dessego\r\n"

	motd2 := "Current players online: " + strconv.Itoa(m.playerCount())

	return []string{motd, motd2}
}

// AddPlayer adds a new player to the current gamestate state.
func (m *Memory) AddPlayer(ip, id string) {
	m.Lock()
	defer m.Unlock()

	m.players[ip] = id
}

func (m *Memory) playerCount() int {
	m.Lock()
	defer m.Unlock()

	return len(m.players)
}
