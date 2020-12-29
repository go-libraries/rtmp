package conn

import (
	"bytes"
	"encoding/binary"

	"io"
	"net"
)

func read(reader io.Reader) error {
	defer func() {

		if r := recover(); r != nil {
			if reader.err == nil {
				conn.err = r.(error)
				logger.ModulePrintf(logHandler, log.LOG_LEVEL_WARNING,
					"readLoop panic:", conn.err)
			}
		}
		conn.Close()
		conn.handler.OnClosed(conn)
	}()

	var found bool
	var chunkstream *InboundChunkStream
	var remain uint32
	for !conn.closed {
		// Read base header
		n, vfmt, csi, err := ReadBaseHeader(conn.br)
		CheckError(err, "ReadBaseHeader")
		conn.inBytes += uint32(n)
		// Get chunk stream
		chunkstream, found = conn.inChunkStreams[csi]
		if !found || chunkstream == nil {
			logger.ModulePrintf(logHandler, log.LOG_LEVEL_TRACE, "New stream 1 csi: %d, fmt: %d\n", csi, vfmt)
			chunkstream = NewInboundChunkStream(csi)
			conn.inChunkStreams[csi] = chunkstream
		}
		// Read header
		header := &Header{}
		n, err = header.ReadHeader(conn.br, vfmt, csi, chunkstream.lastHeader)
		CheckError(err, "ReadHeader")
		if !found {
			logger.ModulePrintf(logHandler, log.LOG_LEVEL_TRACE, "New stream 2 csi: %d, fmt: %d, header: %+v\n", csi, vfmt, header)
		}
		conn.inBytes += uint32(n)
		var absoluteTimestamp uint32
		var message *Message
		switch vfmt {
		case HEADER_FMT_FULL:
			chunkstream.lastHeader = header
			absoluteTimestamp = header.Timestamp
		case HEADER_FMT_SAME_STREAM:
			// A new message with same stream ID
			if chunkstream.lastHeader == nil {
				logger.ModulePrintf(logHandler, log.LOG_LEVEL_WARNING,
					"A new message with fmt: %d, csi: %d\n", vfmt, csi)
				header.Dump("err")
			} else {
				header.MessageStreamID = chunkstream.lastHeader.MessageStreamID
			}
			chunkstream.lastHeader = header
			absoluteTimestamp = chunkstream.lastInAbsoluteTimestamp + header.Timestamp
		case HEADER_FMT_SAME_LENGTH_AND_STREAM:
			// A new message with same stream ID, message length and message type
			if chunkstream.lastHeader == nil {
				logger.ModulePrintf(logHandler, log.LOG_LEVEL_WARNING,
					"A new message with fmt: %d, csi: %d\n", vfmt, csi)
				header.Dump("err")
			}
			header.MessageStreamID = chunkstream.lastHeader.MessageStreamID
			header.MessageLength = chunkstream.lastHeader.MessageLength
			header.MessageTypeID = chunkstream.lastHeader.MessageTypeID
			chunkstream.lastHeader = header
			absoluteTimestamp = chunkstream.lastInAbsoluteTimestamp + header.Timestamp
		case HEADER_FMT_CONTINUATION:
			if chunkstream.receivedMessage != nil {
				// Continuation the previous unfinished message
				message = chunkstream.receivedMessage
			}
			if chunkstream.lastHeader == nil {
				logger.ModulePrintf(logHandler, log.LOG_LEVEL_WARNING,
					"A new message with fmt: %d, csi: %d\n", vfmt, csi)
				header.Dump("err")
			} else {
				header.MessageStreamID = chunkstream.lastHeader.MessageStreamID
				header.MessageLength = chunkstream.lastHeader.MessageLength
				header.MessageTypeID = chunkstream.lastHeader.MessageTypeID
				header.Timestamp = chunkstream.lastHeader.Timestamp
			}
			chunkstream.lastHeader = header
			absoluteTimestamp = chunkstream.lastInAbsoluteTimestamp
		}
		if message == nil {
			// New message
			message = &Message{
				ChunkStreamID:     csi,
				Type:              header.MessageTypeID,
				Timestamp:         header.RealTimestamp(),
				Size:              header.MessageLength,
				StreamID:          header.MessageStreamID,
				Buf:               new(bytes.Buffer),
				IsInbound:         true,
				AbsoluteTimestamp: absoluteTimestamp,
			}
		}
		chunkstream.lastInAbsoluteTimestamp = absoluteTimestamp
		// Read data
		remain = message.Remain()
		var n64 int64
		if remain <= conn.inChunkSize {
			// One chunk message
			for {
				// n64, err = CopyNFromNetwork(message.Buf, conn.br, int64(remain))
				n64, err = io.CopyN(message.Buf, conn.br, int64(remain))
				if err == nil {
					conn.inBytes += uint32(n64)
					if remain <= uint32(n64) {
						break
					} else {
						remain -= uint32(n64)
						logger.ModulePrintf(logHandler, log.LOG_LEVEL_TRACE,
							"Message continue copy remain: %d\n", remain)
						continue
					}
				}
				netErr, ok := err.(net.Error)
				if !ok || !netErr.Temporary() {
					CheckError(err, "Read data 1")
				}
				logger.ModulePrintf(logHandler, log.LOG_LEVEL_TRACE,
					"Message copy blocked!\n")
			}
			// Finished message
			conn.received(message)
			chunkstream.receivedMessage = nil
		} else {
			// Unfinish
			logger.ModulePrintf(logHandler, log.LOG_LEVEL_DEBUG,
				"Unfinish message(remain: %d, chunksize: %d)\n", remain, conn.inChunkSize)

			remain = conn.inChunkSize
			for {
				// n64, err = CopyNFromNetwork(message.Buf, conn.br, int64(remain))
				n64, err = io.CopyN(message.Buf, conn.br, int64(remain))
				if err == nil {
					conn.inBytes += uint32(n64)
					if remain <= uint32(n64) {
						break
					} else {
						remain -= uint32(n64)
						logger.ModulePrintf(logHandler, log.LOG_LEVEL_TRACE,
							"Unfinish message continue copy remain: %d\n", remain)
						continue
					}
					break
				}
				netErr, ok := err.(net.Error)
				if !ok || !netErr.Temporary() {
					CheckError(err, "Read data 2")
				}
				logger.ModulePrintf(logHandler, log.LOG_LEVEL_TRACE,
					"Unfinish message copy blocked!\n")
			}
			chunkstream.receivedMessage = message
		}

		// Check window
		if conn.inBytes > (conn.inBytesPreWindow + conn.inWindowSize) {
			// Send window acknowledgement
			ackmessage := NewMessage(CS_ID_PROTOCOL_CONTROL, ACKNOWLEDGEMENT, 0, absoluteTimestamp+1, nil)
			err = binary.Write(ackmessage.Buf, binary.BigEndian, conn.inBytes)
			CheckError(err, "ACK Message write data")
			conn.inBytesPreWindow = conn.inBytes
			conn.Send(ackmessage)
		}
	}
}
