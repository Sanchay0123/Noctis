package protocol

import (
	"errors"
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var ErrUnknownFieldsRejected = errors.New("packet contains unknown fields")

// UnmarshalStrict unmarshals the packet and explicitly enforces the Unknown-Field Policy.
func UnmarshalStrict(b []byte, m protoreflect.ProtoMessage) error {
	if err := proto.Unmarshal(b, m); err != nil {
		return err
	}
	if len(m.ProtoReflect().GetUnknown()) > 0 {
		return ErrUnknownFieldsRejected
	}
	return nil
}

func TestProtobufSerialization(t *testing.T) {
	// 1. Verify that valid schema-generated structures can marshal and unmarshal.
	packet := &MeshPacket{
		Version:    1,
		Type:       PacketType_PACKET_TYPE_INIT,
		PacketId:   []byte("1234567890123456"),
		Ttl:        16,
		SourceNode: make([]byte, 32),
		DestNode:   make([]byte, 32),
		Payload: &MeshPacket_Init{
			Init: &InitPayload{
				EphemeralKey: make([]byte, 32),
				Signature:    make([]byte, 64),
			},
		},
	}

	data, err := proto.Marshal(packet)
	if err != nil {
		t.Fatalf("Failed to marshal packet: %v", err)
	}

	unmarshaled := &MeshPacket{}
	if err := UnmarshalStrict(data, unmarshaled); err != nil {
		t.Fatalf("Failed to strict-unmarshal packet: %v", err)
	}

	if unmarshaled.Version != 1 {
		t.Errorf("Expected version 1, got %d", unmarshaled.Version)
	}
}

func TestUnknownFieldRejection(t *testing.T) {
	// Construct a packet and append a mock unknown field (tag 99, wire type 2, length 4, "junk")
	// Wire encoding: tag = (99 << 3) | 2 = 794. Varint encoded: 10011010 00000110 (0x9a 0x06)
	// Length: 4. Data: "junk"
	validPacket := &MeshPacket{
		Version: 1,
	}
	data, _ := proto.Marshal(validPacket)

	// Append unknown field to binary payload manually
	dataWithUnknown := append(data, []byte{0x9a, 0x06, 0x04, 'j', 'u', 'n', 'k'}...)

	// Normal proto.Unmarshal ignores unknown fields and succeeds
	lenientPacket := &MeshPacket{}
	if err := proto.Unmarshal(dataWithUnknown, lenientPacket); err != nil {
		t.Fatalf("Lenient unmarshal failed unexpectedly: %v", err)
	}
	if len(lenientPacket.ProtoReflect().GetUnknown()) == 0 {
		t.Fatalf("Expected unknown fields to be preserved in the lenient packet reflection")
	}

	// Strict unmarshal MUST reject the packet
	strictPacket := &MeshPacket{}
	if err := UnmarshalStrict(dataWithUnknown, strictPacket); err != ErrUnknownFieldsRejected {
		t.Errorf("Expected strict unmarshal to reject unknown fields with ErrUnknownFieldsRejected, got: %v", err)
	}
}

func TestMalformedDataRejection(t *testing.T) {
	// 3. Verify malformed/truncated serialized data is rejected
	badData := []byte{0xff, 0xff, 0xff} // Invalid protobuf payload

	unmarshaled := &MeshPacket{}
	if err := UnmarshalStrict(badData, unmarshaled); err == nil {
		t.Errorf("Expected malformed data to be rejected, but it succeeded")
	}
}
