package message

// Chunk Message Header - "fmt" field values
const (
	HeaderFmtFull                = 0x00
	HeaderFmtSameStream          = 0x01
	HeaderFmtSameLengthAbdStream = 0x02
	HeaderFmtContinuation        = 0x03
)

// Result codes
const (
	ResultConnectOk           = "NetConnection.Connect.Success"
	ResultConnectRejected     = "NetConnection.Connect.Rejected"
	ResultConnectOkDesc       = "Connection success."
	ResultConnectRejectedDesc = "[ AccessManager.Reject ] : [ code=400 ] : "
	NetStreamPlayStart        = "NetStream.Play.Start"
	NetStreamPlayReset        = "NetStream.Play.Reset"
	NetStreamPublishStart     = "NetStream.Publish.Start"
)

// Chunk stream ID
const (
	CsIdProtocolControl = uint32(2)
	CsIdCommand         = uint32(3)
	CsIdUserControl     = uint32(4)
)

// Message type
const (
	// Set Chunk Size
	//
	// Protocol control message 1, Set Chunk Size, is used to notify the
	// peer a new maximum chunk size to use.

	// The value of the chunk size is carried as 4-byte message payload. A
	// default value exists for chunk size, but if the sender wants to
	// change this value it notifies the peer about it through this
	// protocol message. For example, a client wants to send 131 bytes of
	// data and the chunk size is at its default value of 128. So every
	// message from the client gets split into two chunks. The client can
	// choose to change the chunk size to 131 so that every message get
	// split into two chunks. The client MUST send this protocol message to
	// the server to notify that the chunk size is set to 131 bytes.
	// The maximum chunk size can be 65536 bytes. Chunk size is maintained
	// independently for server to client communication and client to server
	// communication.
	//
	//  0                   1                   2                   3
	//  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |                          chunk size (4 bytes)                 |
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// Figure 2 Pay load for the protocol message ‘Set Chunk Size’
	//
	// chunk size: 32 bits
	//   This field holds the new chunk size, which will be used for all
	//   future chunks sent by this chunk stream.
	SetChunkSize = uint8(1)

	// Abort Message
	//
	// Protocol control message 2, Abort Message, is used to notify the peer
	// if it is waiting for chunks to complete a message, then to discard
	// the partially received message over a chunk stream and abort
	// processing of that message. The peer receives the chunk stream ID of
	// the message to be discarded as payload of this protocol message. This
	// message is sent when the sender has sent part of a message, but wants
	// to tell the receiver that the rest of the message will not be sent.
	//
	//  0                   1                   2                   3
	//  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |                        chunk stream id (4 bytes)              |
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// Figure 3 Pay load for the protocol message ‘Abort Message’.
	//
	//
	// chunk stream ID: 32 bits
	//   This field holds the chunk stream ID, whose message is to be
	//   discarded.
	AbortMessage = uint8(2)

	// Acknowledgement
	//
	// The client or the server sends the acknowledgment to the peer after
	// receiving bytes equal to the window size. The window size is the
	// maximum number of bytes that the sender sends without receiving
	// acknowledgment from the receiver. The server sends the window size to
	// the client after application connects. This message specifies the
	// sequence number, which is the number of the bytes received so far.
	//  0                   1                   2                   3
	//  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |                        sequence number (4 bytes)              |
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// Figure 4 Pay load for the protocol message ‘Acknowledgement’.
	//
	// sequence number: 32 bits
	//   This field holds the number of bytes received so far.
	Acknowledgement = uint8(3)

	// User Control Message
	//
	// The client or the server sends this message to notify the peer about
	// the user control events. This message carries Event type and Event
	// data.
	// +------------------------------+-------------------------
	// |     Event Type ( 2- bytes ) | Event Data
	// +------------------------------+-------------------------
	// Figure 5 Pay load for the ‘User Control Message’.
	//
	//
	// The first 2 bytes of the message data are used to identify the Event
	// type. Event type is followed by Event data. Size of Event data field
	// is variable.
	UserControlMessage = uint8(4)

	// Window Acknowledgement Size
	//
	// The client or the server sends this message to inform the peer which
	// window size to use when sending acknowledgment. For example, a server
	// expects acknowledgment from the client every time the server sends
	// bytes equivalent to the window size. The server updates the client
	// about its window size after successful processing of a connect
	// request from the client.
	//
	//  0                   1                   2                   3
	//  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |                   Acknowledgement Window size (4 bytes)       |
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// Figure 6 Pay load for ‘Window Acknowledgement Size’.
	WindowAcknowledgementSize = uint8(5)

	// Set Peer Bandwidth
	//
	// The client or the server sends this message to update the output
	// bandwidth of the peer. The output bandwidth value is the same as the
	// window size for the peer. The peer sends ‘Window Acknowledgement
	// Size’ back if its present window size is different from the one
	// received in the message.
	//  0                   1                   2                   3
	//  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |                   Acknowledgement Window size                 |
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// | Limit type    |
	// +-+-+-+-+-+-+-+-+
	// Figure 7 Pay load for ‘Set Peer Bandwidth’
	//
	// The sender can mark this message hard (0), soft (1), or dynamic (2)
	// using the Limit type field. In a hard (0) request, the peer must send
	// the data in the provided bandwidth. In a soft (1) request, the
	// bandwidth is at the discretion of the peer and the sender can limit
	// the bandwidth. In a dynamic (2) request, the bandwidth can be hard or
	// soft.
	SetPeerBandwidth = uint8(6)

	// Audio message
	//
	// The client or the server sends this message to send audio data to the
	// peer. The message type value of 8 is reserved for audio messages.
	AudioType = uint8(8)

	// Video message
	//
	// The client or the server sends this message to send video data to the
	// peer. The message type value of 9 is reserved for video messages.
	// These messages are large and can delay the sending of other type of
	// messages. To avoid such a situation, the video message is assigned
	// the lowest priority.
	VideoType = uint8(9)

	// Aggregate message
	//
	// An aggregate message is a single message that contains a list of sub-
	// messages. The message type value of 22 is reserved for aggregate
	// messages.
	AggregateType = uint8(22)

	// Shared object message
	//
	// A shared object is a Flash object (a collection of name value pairs)
	// that are in synchronization across multiple clients, instances, and
	// so on. The message types kMsgContainer=19 for AMF0 and
	// kMsgContainerEx=16 for AMF3 are reserved for shared object events.
	// Each message can contain multiple events.
	SharedObjectAmf0 = uint8(19)
	SharedObjectAmf3 = uint8(16)
	// Data message
	//
	// The client or the server sends this message to send Metadata or any
	// user data to the peer. Metadata includes details about the
	// data(audio, video etc.) like creation time, duration, theme and so
	// on. These messages have been assigned message type value of 18 for
	// AMF0 and message type value of 15 for AMF3.
	DataAmf0 = uint8(18)
	DataAmf3 = uint8(15)

	// Command message
	//
	// Command messages carry the AMF-encoded commands between the client
	// and the server. These messages have been assigned message type value
	// of 20 for AMF0 encoding and message type value of 17 for AMF3
	// encoding. These messages are sent to perform some operations like
	// connect, createStream, publish, play, pause on the peer. Command
	// messages like onstatus, result etc. are used to inform the sender
	// about the status of the requested commands. A command message
	// consists of command name, transaction ID, and command object that
	// contains related parameters. A client or a server can request Remote
	// Procedure Calls (RPC) over streams that are communicated using the
	// command messages to the peer.
	CommandAmf0 = uint8(20)
	CommandAmf3 = uint8(17) // Keng-die!!! Just ignore one byte before AMF0.
)

const (
	EventStreamBEGIN    = uint16(0)
	EventStreamEOF      = uint16(1)
	EventStreamDRY      = uint16(2)
	EventSetBuffLength  = uint16(3)
	EventStreamIsRecord = uint16(4)
	EventPingREQUEST    = uint16(6)
	EventPingRESPONSE   = uint16(7)
	EventRequestVerify  = uint16(0x1a)
	EventResponseVerify = uint16(0x1b)
	EventBufferEmpty    = uint16(0x1f)
	EventBufferReady    = uint16(0x20)
)

const (
	BindWidthLimitHard    = uint8(0)
	BindWidthLimitSoft    = uint8(1)
	BindWidthLimitDynamic = uint8(2)
)

var (
	MinBufferLength = uint32(256)
	FmsVersion      = []byte{0x04, 0x05, 0x00, 0x01}
	FmsVersionStr   = "4,5,0,297"
)

const (
	MaxTimestamp                    = uint32(2000000000)
	AutoTimestamp                   = uint32(0XFFFFFFFF)
	DefaultHighPriorityBufferSize   = 2048
	DefaultMiddlePriorityBufferSize = 128
	DefaultLowPriorityBufferSize    = 64
	DefaultChunkSize                = uint32(128)
	DefaultWindowSize               = 2500000
	DefaultCapabilities             = float64(15)
	DefaultAudioCodecs              = float64(4071)
	DefaultVideoCodecs              = float64(252)
	FmsCapabilities                 = uint32(255)
	FmsMode                         = uint32(2)
	SetPeerBandWidthHard            = byte(0)
	SetPeerBandWidthSoft            = byte(1)
	SetPeerBandWidthDynamic         = byte(2)
)
