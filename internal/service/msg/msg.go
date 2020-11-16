package msg

import (
	"encoding/binary"
	"math"
)

// BloodMsg represents a blood message.
type BloodMsg struct {
	ID           int
	CharacterID  string
	BlockID      int
	PosX         float32
	PosY         float32
	PosZ         float32
	AngX         float32
	AngY         float32
	AngZ         float32
	MsgID        int
	MainMsgID    int
	AddMsgCateID int
	Rating       int
	Legacy       int
}

// NewBloodMsgFromBytes returns a blood parsed from the given byte slice.
//
// Messages expected in the format:
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
func NewBloodMsgFromBytes(b []byte) (bm *BloodMsg, err error) {
	bm = &BloodMsg{
		Legacy: 1,
	}

	if len(b) < 4 {
		return nil, nil
	}

	// Message ID.
	cursor := 4
	bm.ID = int(binary.LittleEndian.Uint32(b[:cursor]))

	// Character ID.
	for i := cursor; ; i++ {
		c := b[i : i+1][0]
		if c == 0x00 {
			cursor = i + 1
			break
		}
		bm.CharacterID += string(c)
	}

	// BlockID.
	bm.BlockID = int(binary.LittleEndian.Uint32(b[cursor : cursor+4]))
	cursor += 4

	// Positional data.
	bm.PosX = math.Float32frombits(binary.LittleEndian.Uint32(b[cursor : cursor+4]))
	cursor += 4
	bm.PosY = math.Float32frombits(binary.LittleEndian.Uint32(b[cursor : cursor+4]))
	cursor += 4
	bm.PosZ = math.Float32frombits(binary.LittleEndian.Uint32(b[cursor : cursor+4]))
	cursor += 4
	bm.AngX = math.Float32frombits(binary.LittleEndian.Uint32(b[cursor : cursor+4]))
	cursor += 4
	bm.AngY = math.Float32frombits(binary.LittleEndian.Uint32(b[cursor : cursor+4]))
	cursor += 4
	bm.AngZ = math.Float32frombits(binary.LittleEndian.Uint32(b[cursor : cursor+4]))
	cursor += 4

	// Metadata.
	bm.MsgID = int(binary.LittleEndian.Uint32(b[cursor : cursor+4]))
	cursor += 4
	bm.MainMsgID = int(binary.LittleEndian.Uint32(b[cursor : cursor+4]))
	cursor += 4
	bm.AddMsgCateID = int(binary.LittleEndian.Uint32(b[cursor : cursor+4]))
	cursor += 4
	bm.Rating = int(binary.LittleEndian.Uint32(b[cursor : cursor+4]))

	return bm, nil
}

func (bm BloodMsg) Bytes() []byte {
	// TODO

	return []byte{}
}
