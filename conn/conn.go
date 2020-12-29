package conn

import "rtmp/message"

type Conn interface {
	Close()
	Send(message *message.Message) error
	CreateChunkStream(id uint32) (*OutboundChunkStream, error)
	CloseChunkStream(id uint32)
	NewTransactionId() uint32
	CreateMediaChunkStream() (*OutboundChunkStream, error)
	CloseMediaChunkStream(id uint32)
	SetStreamBufferSize(streamId uint32, size uint32)
	OutboundChunkStream(id uint32) (chunkStream *OutboundChunkStream, found bool)
	InboundChunkStream(id uint32) (chunkStream *InboundChunkStream, found bool)
	SetWindowAcknowledgementSize()
	SetPeerBandwidth(peerBandwidth uint32, limitType byte)
	SetChunkSize(chunkSize uint32)
	SendUserControlMessage(eventId uint16)
}
