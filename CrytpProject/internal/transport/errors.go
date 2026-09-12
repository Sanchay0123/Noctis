package transport

import "errors"

var (
	ErrMalformedFrame        = errors.New("malformed tcp frame")
	ErrOversizedFrame        = errors.New("oversized tcp frame")
	ErrUnknownProtobufFields = errors.New("unknown protobuf fields detected")
	ErrInvalidPacket         = errors.New("invalid protocol packet")
	ErrUnsupportedVersion    = errors.New("unsupported protocol version")
	ErrInvalidPacketType     = errors.New("invalid packet type or oneof mismatch")
	ErrSessionMismatch       = errors.New("session mismatch or unknown session")
	ErrWrongDestination      = errors.New("wrong destination node identity")
	ErrWrongSource           = errors.New("wrong source node identity")
	ErrDecryptionFailed      = errors.New("decryption or authentication failed")
	ErrConnectionClosed      = errors.New("connection closed")
	ErrPacketIDGeneration    = errors.New("failed to generate secure packet ID")
	ErrSessionNotEstablished = errors.New("session not established")
)
