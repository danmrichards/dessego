package sos

import (
	"bytes"
	"encoding/binary"
	"math"
	"time"
)

// maxSOSAge is the maximum age at which we consider an SOS to be usable.
const maxSOSAge = time.Second * 30

// SOS represents a Demon's Souls SOS.
type SOS struct {
	ID            uint32
	CharacterID   string
	BlockID       int32
	PosX          float32
	PosY          float32
	PosZ          float32
	AngX          float32
	AngY          float32
	AngZ          float32
	MsgID         uint32
	MainMsgID     uint32
	AddMsgCateID  uint32
	PlayerInfo    string
	QWCWB         uint32
	QWCLR         uint32
	Black         byte
	PlayerLevel   uint32
	Ratings       []int
	TotalSessions int
	Updated       time.Time
}

func (s SOS) Bytes() []byte {
	data := new(bytes.Buffer)

	// Message ID.
	binary.Write(data, binary.LittleEndian, s.ID)

	// Character ID.
	data.WriteString(s.CharacterID)
	data.WriteByte(0x00)

	// Block ID.
	binary.Write(data, binary.LittleEndian, uint32(s.BlockID))

	// Positional data.
	binary.Write(data, binary.LittleEndian, math.Float32bits(s.PosX))
	binary.Write(data, binary.LittleEndian, math.Float32bits(s.PosY))
	binary.Write(data, binary.LittleEndian, math.Float32bits(s.PosZ))
	binary.Write(data, binary.LittleEndian, math.Float32bits(s.AngX))
	binary.Write(data, binary.LittleEndian, math.Float32bits(s.AngY))
	binary.Write(data, binary.LittleEndian, math.Float32bits(s.AngZ))

	// Metadata.
	binary.Write(data, binary.LittleEndian, s.MsgID)
	binary.Write(data, binary.LittleEndian, s.MainMsgID)
	binary.Write(data, binary.LittleEndian, s.AddMsgCateID)

	// TODO: Not sure what this is?
	binary.Write(data, binary.LittleEndian, uint32(0))

	// Ratings.
	for _, r := range s.Ratings {
		binary.Write(data, binary.LittleEndian, r)
	}

	// TODO: Not sure what this is?
	binary.Write(data, binary.LittleEndian, uint32(0))

	// Total sessions.
	binary.Write(data, binary.LittleEndian, s.TotalSessions)

	// Player info.
	data.WriteString(s.PlayerInfo)
	data.WriteByte(0x00)

	// World Tendency.
	binary.Write(data, binary.LittleEndian, s.QWCWB)
	binary.Write(data, binary.LittleEndian, s.QWCLR)

	// Black phantom.
	binary.Write(data, binary.LittleEndian, s.Black)

	return data.Bytes()
}
