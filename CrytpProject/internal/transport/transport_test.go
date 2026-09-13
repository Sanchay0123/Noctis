package transport

import (
	"crypto/ed25519"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"sync"
	"testing"

	"github.com/sanchayjain/meshchat/internal/crypto"
	"github.com/sanchayjain/meshchat/internal/protocol"
	"google.golang.org/protobuf/proto"
)

func genIdentities(t *testing.T) (*crypto.NodeIdentity, *crypto.NodeIdentity) {
	_, privA, _ := ed25519.GenerateKey(nil)
	idA, _ := crypto.LoadIdentity(privA)
	_, privB, _ := ed25519.GenerateKey(nil)
	idB, _ := crypto.LoadIdentity(privB)
	return idA, idB
}

// Memory pipe helper
func getPipe() (*Connection, *Connection) {
	c1, c2 := net.Pipe()
	return NewConnection(c1), NewConnection(c2)
}

func TestTCPFramingValidRoundtrip(t *testing.T) {
	c1, c2 := getPipe()
	defer c1.Close()
	defer c2.Close()

	id, _ := GeneratePacketID()
	pkt := &protocol.MeshPacket{
		Version:  1,
		Type:     protocol.PacketType_PACKET_TYPE_UNKNOWN, // Force validation to fail, but transport framing shouldn't care
		PacketId: id,
	}

	go func() {
		// Bypass validation on write! We just test framing.
		data, _ := proto.Marshal(pkt)
		frame := make([]byte, 4+len(data))
		binary.BigEndian.PutUint32(frame, uint32(len(data)))
		copy(frame[4:], data)
		c1.conn.Write(frame)
	}()

	rpkt, err := c2.ReadPacket()
	// Should fail validation, but error is returned after framing parses correctly
	if err != ErrInvalidPacket {
		t.Fatalf("Expected ErrInvalidPacket, got: %v", err)
	}
	if rpkt != nil {
		t.Fatalf("Packet should be nil on error")
	}
}

func TestTCPFramingOversized(t *testing.T) {
	c1, c2 := getPipe()
	defer c1.Close()
	defer c2.Close()

	go func() {
		frame := make([]byte, 4)
		binary.BigEndian.PutUint32(frame, MaxFrameSize+1)
		c1.conn.Write(frame)
	}()

	_, err := c2.ReadPacket()
	if err != ErrOversizedFrame {
		t.Fatalf("Did not reject oversized frame")
	}
}

func TestTCPFramingEOF(t *testing.T) {
	c1, c2 := getPipe()
	c1.Close()

	_, err := c2.ReadPacket()
	if err != ErrMalformedFrame && err != ErrConnectionClosed {
		t.Fatalf("Expected ErrConnectionClosed, got %v", err)
	}
}

func TestTCPFramingPartialRead(t *testing.T) {
	c1, c2 := getPipe()
	defer c1.Close()
	defer c2.Close()

	go func() {
		// Write exactly 2 bytes (partial length) then close
		c1.conn.Write([]byte{0, 0})
		c1.Close()
	}()

	_, err := c2.ReadPacket()
	if err != ErrMalformedFrame && err != ErrConnectionClosed {
		t.Fatalf("Expected ErrMalformedFrame on partial read of length, got %v", err)
	}
}

func TestTCPFramingConcurrentWrites(t *testing.T) {
	c1, c2 := getPipe()
	defer c1.Close()
	defer c2.Close()

	var wg sync.WaitGroup
	// 50 concurrent writers to avoid memory pipe deadlocks
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			pid, _ := GeneratePacketID()
			pkt := &protocol.MeshPacket{
				Version:  1,
				Type:     protocol.PacketType_PACKET_TYPE_UNKNOWN, // We'll ignore the validation error, just check framing doesn't panic
				PacketId: pid,
			}
			// WritePacket encodes and writes safely
			_ = c1.WritePacket(pkt)
		}()
	}

	go func() {
		wg.Wait()
		c1.Close()
	}()

	count := 0
	for {
		_, err := c2.ReadPacket()
		if err == ErrConnectionClosed {
			break
		}
		if err == ErrInvalidPacket {
			count++
		} else {
			t.Fatalf("Unexpected error: %v", err)
		}
	}

	if count != 50 {
		t.Fatalf("Expected 50 packets, got %d", count)
	}
}

func TestValidationUnknownField(t *testing.T) {
	// Construct a raw protobuf payload with an unknown field
	// Field 99 = 1 (varint) -> 99<<3 | 0 = 792 -> [0x18, 0x06]
	// protobuf encoding for Field 99, varint value 1:
	// 99 * 8 = 792. 792 in varint: 0x98 0x06.
	raw := []byte{0x98, 0x06, 0x01}

	// Prepend valid required fields to pass early checks so it reaches unmarshal check
	// But it actually rejects immediately after unmarshal.
	var pkt protocol.MeshPacket
	proto.Unmarshal(raw, &pkt)

	err := ValidatePacket(&pkt)
	if err != ErrUnknownProtobufFields {
		t.Fatalf("Did not reject unknown protobuf field, got: %v", err)
	}
}

func TestValidationInvalidVersion(t *testing.T) {
	pkt := &protocol.MeshPacket{Version: 2}
	if ValidatePacket(pkt) != ErrUnsupportedVersion {
		t.Fatalf("Failed to reject invalid version")
	}
}

func TestValidationInvalidLengths(t *testing.T) {
	pid, _ := GeneratePacketID()
	pkt := &protocol.MeshPacket{
		Version:  1,
		PacketId: pid[:15], // Too short
	}
	if ValidatePacket(pkt) != ErrInvalidPacket {
		t.Fatalf("Failed to reject bad packet ID length")
	}
}

func TestSessionBindingNoSession(t *testing.T) {
	idA, idB := genIdentities(t)
	c1, c2 := getPipe()
	defer c1.Close()
	defer c2.Close()

	chB := NewDirectChannel(c2, idB)

	// A sends APP_DATA before handshake
	pid, _ := GeneratePacketID()
	app := &protocol.MeshPacket{
		Version:    1,
		Type:       protocol.PacketType_PACKET_TYPE_APP_DATA,
		PacketId:   pid,
		SourceNode: idA.PublicKey(),
		DestNode:   idB.PublicKey(),
		Payload: &protocol.MeshPacket_AppData{
			AppData: &protocol.AppDataPayload{
				SessionId:   make([]byte, 32),
				SequenceNum: 0,
				Ciphertext:  make([]byte, 16), // min chacha overhead
			},
		},
	}
	go c1.WritePacket(app)

	_, err := chB.ReceivePlaintext()
	if err != ErrSessionNotEstablished {
		t.Fatalf("Accepted APP_DATA without established session")
	}
}

type shortWriterConn struct {
	net.Conn
	writeSize int
}

func (s *shortWriterConn) Write(b []byte) (n int, err error) {
	if len(b) > s.writeSize {
		return s.Conn.Write(b[:s.writeSize])
	}
	return s.Conn.Write(b)
}

func TestTCPFramingShortWrite(t *testing.T) {
	c1, c2 := net.Pipe()
	defer c1.Close()
	defer c2.Close()

	// Simulate a connection that only writes 2 bytes at a time
	swConn := &shortWriterConn{Conn: c1, writeSize: 2}
	conn1 := NewConnection(swConn)
	conn2 := NewConnection(c2)

	id, _ := GeneratePacketID()
	pkt := &protocol.MeshPacket{
		Version:  1,
		Type:     protocol.PacketType_PACKET_TYPE_UNKNOWN,
		PacketId: id,
	}

	go func() {
		_ = conn1.WritePacket(pkt)
	}()

	_, err := conn2.ReadPacket()
	if err != ErrInvalidPacket {
		t.Fatalf("Expected ErrInvalidPacket, got %v", err)
	}
}

type failingShortWriterConn struct {
	net.Conn
	writes int
}

func (f *failingShortWriterConn) Write(b []byte) (n int, err error) {
	f.writes++
	if f.writes == 1 {
		return f.Conn.Write(b[:1])
	}
	return 0, errors.New("simulated failure")
}

func TestTCPFramingFailingShortWrite(t *testing.T) {
	c1, c2 := net.Pipe()
	defer c1.Close()
	defer c2.Close()

	fwConn := &failingShortWriterConn{Conn: c1}
	conn1 := NewConnection(fwConn)

	id, _ := GeneratePacketID()
	pkt := &protocol.MeshPacket{
		Version:  1,
		Type:     protocol.PacketType_PACKET_TYPE_UNKNOWN,
		PacketId: id,
	}

	go func() { io.ReadAll(c2) }()
	err := conn1.WritePacket(pkt)
	if err == nil || err.Error() != "simulated failure" {
		t.Fatalf("Expected simulated failure, got %v", err)
	}
}
