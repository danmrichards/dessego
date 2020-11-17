package gamestate

import "fmt"

// PlayerNotFoundError is returned when a player cannot be found for an IP.
type PlayerNotFoundError string

func (p PlayerNotFoundError) Error() string {
	return fmt.Sprintf("no character found with IP %q", string(p))
}
