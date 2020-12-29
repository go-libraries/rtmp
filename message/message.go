package message

import (
	"bytes"
)

type Message struct {
	ChunkStreamID     uint32
	Timestamp         uint32
	Size              uint32
	TypeId            uint8
	StreamID          uint32
	Buf               *bytes.Buffer
	IsInbound         bool
	AbsoluteTimestamp uint32
}

func NewMessage(csi uint32, t uint8, sid uint32, ts uint32, data []byte) *Message {
	message := &Message{
		ChunkStreamID:     csi,
		TypeId:            t,
		StreamID:          sid,
		Timestamp:         ts,
		AbsoluteTimestamp: ts,
		Buf:               new(bytes.Buffer),
	}
	if data != nil {
		message.Buf.Write(data)
		message.Size = uint32(len(data))
	}
	return message
}
