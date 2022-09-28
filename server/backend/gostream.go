package backend

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/golang/protobuf/proto"
)

// This is forked from https://github.com/robinpowered/go-proto

type MessageCollection []proto.Message

type UnmarshalFunc func([]byte) (proto.Message, error)

const (
	MaximumDataRecordLength = 1024 * 10
)

// ReadLengthPrefixedCollection reads a collection of protocol buffer messages from the supplied reader.
// Each message is presumed prefixed by a 32 bit varint which represents the size of the ensuing message.
// The UnmarshalFunc argument is a supplied callback used to convert the raw bytes read as a message to the desired message type.
// The protocol buffer message collection is returned, along with any error arising.
// For more detailed information on this approach, see the official protocol buffer documentation https://developers.google.com/protocol-buffers/docs/techniques#streaming.
func ReadLengthPrefixedCollection(ctx context.Context, maximumMessageLength uint64, r io.Reader, f UnmarshalFunc) (pbs MessageCollection, totalBytesRead int, err error) {
	position := 0

	for {
		var prefixBuf [binary.MaxVarintLen32]byte
		var bytesRead, varIntBytes int
		var messageLength uint64
		for varIntBytes == 0 { // i.e. no varint has been decoded yet.
			if bytesRead >= len(prefixBuf) {
				return pbs, position, fmt.Errorf("invalid varint32 encountered (position = %d)", position)
			}
			// We have to read byte by byte here to avoid reading more bytes
			// than required. Each read byte is appended to what we have
			// read before.
			newBytesRead, err := r.Read(prefixBuf[bytesRead : bytesRead+1])
			if newBytesRead == 0 {
				if io.EOF == err {
					return pbs, position, nil
				} else if err != nil {
					return pbs, position, fmt.Errorf("stream:read (position = %d): %w", position, err)
				}
				// A Reader should not return (0, nil), but if it does,
				// it should be treated as no-op (according to the
				// Reader contract). So let's go on...
				continue
			}
			bytesRead += newBytesRead
			position += newBytesRead
			// Now present everything read so far to the varint decoder and
			// see if a varint can be decoded already.
			messageLength, varIntBytes = proto.DecodeVarint(prefixBuf[:bytesRead])
		}

		if messageLength > maximumMessageLength {
			return pbs, position, fmt.Errorf("stream:message-too-big (%d bytes) (position = %d)", messageLength, position)
		}

		beforeRead := position
		messageBuf := make([]byte, messageLength)
		newBytesRead, err := io.ReadFull(r, messageBuf)
		bytesRead += newBytesRead
		if err != nil {
			return pbs, position, fmt.Errorf("stream:incomplete-message (length = %d) (started = %d) (bytes-read = %d) (position = %d) (read-end = %d): %w",
				messageLength, beforeRead, newBytesRead, position, position+newBytesRead, err)
		}

		pb, err := f(messageBuf)
		if nil != err {
			return nil, position, fmt.Errorf("stream:error (position = %d): %w", position, err)
		}

		position += newBytesRead

		pbs = append(pbs, pb)
	}
}
