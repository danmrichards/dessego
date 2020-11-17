package ghost

import (
	"reflect"
	"testing"
	"time"
)

func TestMemory_Get(t *testing.T) {
	gt := time.Date(2020, 11, 17, 13, 0, 0, 0, time.UTC)

	gm := Memory{
		ghosts: map[string]*Ghost{
			"test234": {
				BlockID:     123,
				CharacterID: "test234",
				timestamp:   gt,
			},
		},
		ghostAge: []string{"test234"},
	}

	tcs := []struct {
		name        string
		characterID string
		blockID     int32
		expGhosts   []*Ghost
	}{
		{
			name:        "non-matching character",
			characterID: "test123",
			blockID:     123,
			expGhosts: []*Ghost{{
				BlockID:     123,
				CharacterID: "test234",
				timestamp:   gt,
			}},
		},
		{
			name:        "matching character",
			characterID: "test234",
			blockID:     123,
			expGhosts:   []*Ghost{},
		},
		{
			name:        "incorrect block",
			characterID: "test123",
			blockID:     456,
			expGhosts:   []*Ghost{},
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ghosts := gm.Get(tc.characterID, tc.blockID, 10)
			if !reflect.DeepEqual(ghosts, tc.expGhosts) {
				t.Fatalf("expected: %d ghosts got: %d", len(tc.expGhosts), len(ghosts))
			}
		})
	}
}

func TestMemory_ClearBefore(t *testing.T) {
	gt := time.Date(2020, 11, 17, 13, 0, 0, 0, time.UTC)

	gm := Memory{
		ghosts: map[string]*Ghost{
			"test234": {
				BlockID:     123,
				CharacterID: "test234",
				timestamp:   gt,
			},
			"test456": {
				BlockID:     123,
				CharacterID: "test456",
				timestamp:   gt.Add(-1 * time.Minute),
			},
			"test678": {
				BlockID:     123,
				CharacterID: "test678",
				timestamp:   gt.Add(-30 * time.Second),
			},
		},
		ghostAge: []string{"test456", "test678", "test234"},
	}

	gm.ClearBefore(gt.Add(-35 * time.Second))

	expGhosts := map[string]*Ghost{
		"test234": {
			BlockID:     123,
			CharacterID: "test234",
			timestamp:   gt,
		},
		"test678": {
			BlockID:     123,
			CharacterID: "test678",
			timestamp:   gt.Add(-30 * time.Second),
		},
	}
	expGhostAge := []string{"test678", "test234"}

	if !reflect.DeepEqual(expGhosts, gm.ghosts) {
		t.Fatalf("expected %d ghosts got %d", len(expGhosts), len(gm.ghosts))
	}
	if !reflect.DeepEqual(expGhostAge, gm.ghostAge) {
		t.Fatalf("expected %d ghost ages got %d", len(expGhosts), len(gm.ghosts))
	}
}
