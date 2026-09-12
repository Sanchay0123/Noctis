package transport

import (
	"crypto/rand"
	"github.com/sanchayjain/meshchat/internal/protocol"
	"golang.org/x/crypto/chacha20poly1305"
)

func GeneratePacketID() ([]byte, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return nil, ErrPacketIDGeneration
	}
	return b, nil
}

func ValidatePacket(pkt *protocol.MeshPacket) error {
	if pkt == nil {
		return ErrInvalidPacket
	}

	// 1. Check for unknown protobuf fields at all levels
	if len(pkt.ProtoReflect().GetUnknown()) > 0 {
		return ErrUnknownProtobufFields
	}

	switch p := pkt.Payload.(type) {
	case *protocol.MeshPacket_Init:
		if p.Init != nil && len(p.Init.ProtoReflect().GetUnknown()) > 0 {
			return ErrUnknownProtobufFields
		}
	case *protocol.MeshPacket_Resp:
		if p.Resp != nil && len(p.Resp.ProtoReflect().GetUnknown()) > 0 {
			return ErrUnknownProtobufFields
		}
	case *protocol.MeshPacket_AppData:
		if p.AppData != nil && len(p.AppData.ProtoReflect().GetUnknown()) > 0 {
			return ErrUnknownProtobufFields
		}
	}

	// 2. Protocol Version
	if pkt.Version != 1 {
		return ErrUnsupportedVersion
	}

	// 3. Lengths
	if len(pkt.PacketId) != 16 {
		return ErrInvalidPacket
	}
	if len(pkt.SourceNode) != 32 {
		return ErrInvalidPacket
	}
	if len(pkt.DestNode) != 32 {
		return ErrInvalidPacket
	}

	// 4. Packet Type Validations
	switch pkt.Type {
	case protocol.PacketType_PACKET_TYPE_INIT:
		initPayload, ok := pkt.Payload.(*protocol.MeshPacket_Init)
		if !ok || initPayload.Init == nil {
			return ErrInvalidPacketType
		}
		if len(initPayload.Init.EphemeralKey) != 32 || len(initPayload.Init.Signature) != 64 {
			return ErrInvalidPacket
		}

	case protocol.PacketType_PACKET_TYPE_RESP:
		respPayload, ok := pkt.Payload.(*protocol.MeshPacket_Resp)
		if !ok || respPayload.Resp == nil {
			return ErrInvalidPacketType
		}
		if len(respPayload.Resp.EphemeralKey) != 32 || len(respPayload.Resp.Signature) != 64 {
			return ErrInvalidPacket
		}

	case protocol.PacketType_PACKET_TYPE_APP_DATA:
		appPayload, ok := pkt.Payload.(*protocol.MeshPacket_AppData)
		if !ok || appPayload.AppData == nil {
			return ErrInvalidPacketType
		}
		if len(appPayload.AppData.SessionId) != 32 {
			return ErrInvalidPacket
		}
		// Ciphertext must be at least the Poly1305 overhead size
		if len(appPayload.AppData.Ciphertext) < chacha20poly1305.Overhead {
			return ErrInvalidPacket
		}

	default:
		return ErrInvalidPacketType
	}

	return nil
}
