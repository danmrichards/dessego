package game

// Players is the interface that wraps methods that types must implement to be
// used as a service for managing players.
type Players interface {
	// EnsureCreate creates a player with the given ID and index.
	//
	// If a player with the given ID and index already exists, no error will
	// be returned.
	EnsureCreate(id string, index int) error
}

// State is the interface that wraps methods that types must implement to be
// used as a service for managing gamestate state.
type State interface {
	// Motd returns the messages of the day.
	Motd() []string

	// AddPlayer adds a new player to the current gamestate state.
	AddPlayer(ip, id string)
}
