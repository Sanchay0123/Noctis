package transport

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
	"sync"

	"github.com/sanchayjain/meshchat/internal/protocol"
	"google.golang.org/protobuf/proto"
)

const MaxFrameSize = 65536 // 64 KiB

type Connection struct {
	conn    net.Conn
	writeMu sync.Mutex
}

func NewConnection(conn net.Conn) *Connection {
	return &Connection{
		conn: conn,
	}
}

// WritePacket encodes a MeshPacket, prepends a 4-byte length prefix, and writes it to the socket.
// It is safe for concurrent use by multiple goroutines.
func (c *Connection) WritePacket(pkt *protocol.MeshPacket) error {
	data, err := proto.MarshalOptions{Deterministic: true}.Marshal(pkt)
	if err != nil {
		return err
	}

	if len(data) > MaxFrameSize {
		return ErrOversizedFrame
	}

	frame := make([]byte, 4+len(data))
	binary.BigEndian.PutUint32(frame[0:4], uint32(len(data)))
	copy(frame[4:], data)

	c.writeMu.Lock()
	defer c.writeMu.Unlock()

	written := 0
	for written < len(frame) {
		n, err := c.conn.Write(frame[written:])
		written += n
		if err != nil {
			return err
		}
	}
	return nil
}

// ReadPacket reads a 4-byte length prefix, followed by the frame payload, and decodes the MeshPacket.
// ONLY one goroutine should call this continuously (e.g., in a read loop).
func (c *Connection) ReadPacket() (*protocol.MeshPacket, error) {
	var lenBuf [4]byte
	if _, err := io.ReadFull(c.conn, lenBuf[:]); err != nil {
		if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
			return nil, ErrConnectionClosed
		}
		return nil, ErrMalformedFrame
	}

	length := binary.BigEndian.Uint32(lenBuf[:])
	if length == 0 || length > MaxFrameSize {
		return nil, ErrOversizedFrame
	}

	data := make([]byte, length)
	if _, err := io.ReadFull(c.conn, data); err != nil {
		if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
			return nil, ErrConnectionClosed
		}
		return nil, ErrMalformedFrame
	}

	var pkt protocol.MeshPacket
	// Use strict validation for unknown fields via our ValidatePacket
	if err := proto.Unmarshal(data, &pkt); err != nil {
		return nil, ErrInvalidPacket
	}

	// Reject immediately if validation fails (which includes unknown fields check)
	if err := ValidatePacket(&pkt); err != nil {
		return nil, err
	}

	return &pkt, nil
}

func (c *Connection) Close() error {
	return c.conn.Close()
}
