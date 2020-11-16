package msg

import (
	"bytes"
	"encoding/binary"
	"math"
)

// BloodMsg represents a blood message.
//
// Demon's Souls exchanges messages in the following binary format:
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
type BloodMsg struct {
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
	Rating       uint32
	Legacy       uint32
}

// NewBloodMsgFromBytes returns a blood parsed from the given byte slice.
func NewBloodMsgFromBytes(b []byte) (bm *BloodMsg, err error) {
	bm = &BloodMsg{
		Legacy: 1,
	}

	if len(b) < 4 {
		return nil, nil
	}

	// Message ID.
	cursor := 4
	bm.ID = binary.LittleEndian.Uint32(b[:cursor])

	// Character ID.
	for i := cursor; ; i++ {
		c := b[i : i+1][0]
		if c == 0x00 {
			cursor = i + 1
			break
		}
		bm.CharacterID += string(c)
	}

	// Block ID.
	bm.BlockID = int32(binary.LittleEndian.Uint32(b[cursor : cursor+4]))
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
	bm.MsgID = binary.LittleEndian.Uint32(b[cursor : cursor+4])
	cursor += 4
	bm.MainMsgID = binary.LittleEndian.Uint32(b[cursor : cursor+4])
	cursor += 4
	bm.AddMsgCateID = binary.LittleEndian.Uint32(b[cursor : cursor+4])
	cursor += 4
	bm.Rating = binary.LittleEndian.Uint32(b[cursor : cursor+4])

	return bm, nil
}

// Bytes returns a serialised version of the message as a byte slice.
func (bm BloodMsg) Bytes() []byte {
	data := new(bytes.Buffer)

	// Message ID.
	binary.Write(data, binary.LittleEndian, bm.ID)

	// Character ID.
	data.WriteString(bm.CharacterID)
	data.WriteByte(0x00)

	// Block ID.
	binary.Write(data, binary.LittleEndian, uint32(bm.BlockID))

	// Positional data.
	binary.Write(data, binary.LittleEndian, math.Float32bits(bm.PosX))
	binary.Write(data, binary.LittleEndian, math.Float32bits(bm.PosY))
	binary.Write(data, binary.LittleEndian, math.Float32bits(bm.PosZ))
	binary.Write(data, binary.LittleEndian, math.Float32bits(bm.AngX))
	binary.Write(data, binary.LittleEndian, math.Float32bits(bm.AngY))
	binary.Write(data, binary.LittleEndian, math.Float32bits(bm.AngZ))

	// Metadata.
	binary.Write(data, binary.LittleEndian, bm.MsgID)
	binary.Write(data, binary.LittleEndian, bm.MainMsgID)
	binary.Write(data, binary.LittleEndian, bm.AddMsgCateID)
	binary.Write(data, binary.LittleEndian, bm.Rating)

	return data.Bytes()
}
