package ghost

import "fmt"

// CharacterGhostNotFoundError is returned when a ghost cannot be found for a character.
type CharacterGhostNotFoundError string

func (g CharacterGhostNotFoundError) Error() string {
	return fmt.Sprintf("no ghost found for character %q", string(g))
}
