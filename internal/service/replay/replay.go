package replay

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"

	"github.com/danmrichards/dessego/internal/service/gamestate"
)

// LegacyType indicates if a replay is legacy or not.
type LegacyType int

const (
	// NonLegacy indicates a reply created by Demon's Souls interacting with
	// this server.
	NonLegacy LegacyType = iota

	// Legacy indicates a replay imported from a legacy data dump.
	Legacy
)

// Replay represents a replay of a players actions in Demon's Souls.
//
// Demon's Souls exchanges replays in the following binary format:
//
// MessageID 		= 4 bytes
// CharacterID 		= n bytes (terminated by a zero byte)
// BlockID 			= 4 bytes
// Positional data 	= 24 bytes (each element at 4 bytes each)
//   - X position
//   - Y position
//   - Z position
//   - X angle
//   - Y angle
//   - Z angle
// Metadata 		= 16 bytes (each element at 4 bytes each)
//	 - Message ID
//   - Main message ID
//   - Add Message Cate ID (?)
//   - Rating
// Data 			= n bytes (terminated by a zero byte)
type Replay struct {
	ID           uint32
	CharacterID  string
	BlockID      int32
	PosX         float32
	PosY         float32
	PosZ         float32
	AngX         float32
	AngY         float32
	AngZ         float32
	MsgID        uint32
	MainMsgID    uint32
	AddMsgCateID uint32
	Data         []byte
	Legacy       uint32
}

// NewReplayFromBytes returns a replay parsed from the given byte slice.
func NewReplayFromBytes(b []byte) (r *Replay, err error) {
	r = &Replay{
		Legacy: 1,
	}

	if len(b) < 4 {
		return nil, nil
	}

	// Message ID.
	cursor := 4
	r.ID = binary.LittleEndian.Uint32(b[:cursor])

	// Character ID.
	for i := cursor; ; i++ {
		c := b[i : i+1][0]
		if c == 0x00 {
			cursor = i + 1
			break
		}
		r.CharacterID += string(c)
	}

	// Block ID.
	r.BlockID = int32(binary.LittleEndian.Uint32(b[cursor : cursor+4]))
	cursor += 4

	// Positional data.
	r.PosX = math.Float32frombits(binary.LittleEndian.Uint32(b[cursor : cursor+4]))
	cursor += 4
	r.PosY = math.Float32frombits(binary.LittleEndian.Uint32(b[cursor : cursor+4]))
	cursor += 4
	r.PosZ = math.Float32frombits(binary.LittleEndian.Uint32(b[cursor : cursor+4]))
	cursor += 4
	r.AngX = math.Float32frombits(binary.LittleEndian.Uint32(b[cursor : cursor+4]))
	cursor += 4
	r.AngY = math.Float32frombits(binary.LittleEndian.Uint32(b[cursor : cursor+4]))
	cursor += 4
	r.AngZ = math.Float32frombits(binary.LittleEndian.Uint32(b[cursor : cursor+4]))
	cursor += 4

	// Metadata.
	r.MsgID = binary.LittleEndian.Uint32(b[cursor : cursor+4])
	cursor += 4
	r.MainMsgID = binary.LittleEndian.Uint32(b[cursor : cursor+4])
	cursor += 4
	r.AddMsgCateID = binary.LittleEndian.Uint32(b[cursor : cursor+4])
	cursor += 4

	// Replay data.
	for i := cursor; ; i++ {
		c := b[i : i+1][0]
		if c == 0x00 {
			break
		}
		r.Data = append(r.Data, c)
	}

	return r, nil
}

// Header returns the header information for the replay, binary encoded into a
// byte slice.
func (r Replay) Header() []byte {
	data := new(bytes.Buffer)

	// Message ID.
	binary.Write(data, binary.LittleEndian, r.ID)

	// Character ID.
	data.WriteString(r.CharacterID)
	data.WriteByte(0x00)

	// Block ID.
	binary.Write(data, binary.LittleEndian, uint32(r.BlockID))

	// Positional data.
	binary.Write(data, binary.LittleEndian, math.Float32bits(r.PosX))
	binary.Write(data, binary.LittleEndian, math.Float32bits(r.PosY))
	binary.Write(data, binary.LittleEndian, math.Float32bits(r.PosZ))
	binary.Write(data, binary.LittleEndian, math.Float32bits(r.AngX))
	binary.Write(data, binary.LittleEndian, math.Float32bits(r.AngY))
	binary.Write(data, binary.LittleEndian, math.Float32bits(r.AngZ))

	// Metadata.
	binary.Write(data, binary.LittleEndian, r.MsgID)
	binary.Write(data, binary.LittleEndian, r.MainMsgID)
	binary.Write(data, binary.LittleEndian, r.AddMsgCateID)

	return data.Bytes()
}

// String implements fmt.Stringer
func (r Replay) String() string {
	return fmt.Sprintf(
		"id: %d character: %q block: %q",
		r.ID, r.CharacterID, gamestate.Block(r.BlockID),
	)
}
