package crypto

import (
	"crypto/ecdh"
	"crypto/rand"
	"errors"
)

var (
	// ErrMalformedX25519PublicKey is returned when an X25519 public key is not exactly 32 bytes.
	ErrMalformedX25519PublicKey = errors.New("malformed X25519 public key: must be exactly 32 bytes")
	// ErrInvalidX25519PublicKey is returned when an X25519 public key is structurally invalid or represents a low-order point.
	ErrInvalidX25519PublicKey = errors.New("invalid X25519 public key: low order or invalid point")
	// ErrEphemeralKeyDestroyed is returned if an operation is attempted after the key is destroyed.
	ErrEphemeralKeyDestroyed = errors.New("ephemeral key has been destroyed")
)

// EphemeralKey represents a short-lived X25519 keypair for key agreement.
// It is strictly separated from long-term Ed25519 identities.
//
// Lifecycle and Concurrency Contract:
// The primitive represents a single ephemeral keypair. The future handshake/session
// layer is responsible for creating a fresh EphemeralKey for every session establishment.
// EphemeralKey is NOT concurrency-safe with respect to destruction. The Destroy() method
// must NOT be called concurrently with PublicKey() or ComputeSharedSecret().
type EphemeralKey struct {
	privateKey *ecdh.PrivateKey
}

// GenerateEphemeralKey securely creates a fresh X25519 keypair using crypto/rand.
func GenerateEphemeralKey() (*EphemeralKey, error) {
	priv, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	return &EphemeralKey{privateKey: priv}, nil
}

// PublicKey returns a defensive copy of the exact 32-byte X25519 public key.
func (k *EphemeralKey) PublicKey() ([]byte, error) {
	if k.privateKey == nil {
		return nil, ErrEphemeralKeyDestroyed
	}
	// ecdh.PublicKey.Bytes() returns a fresh copy, but we explicitly copy it again
	// to strictly guarantee defensive ownership semantics mathematically against future stdlib changes.
	pubBytes := k.privateKey.PublicKey().Bytes()
	pubCopy := make([]byte, len(pubBytes))
	copy(pubCopy, pubBytes)
	return pubCopy, nil
}

// ComputeSharedSecret computes the X25519 ECDH shared secret with a peer's public key.
// It requires an exact 32-byte input. It relies on the crypto/ecdh guarantees:
// NewPublicKey natively rejects malformed or invalid points.
// ECDH natively rejects low-order points (e.g. all-zero shared secrets) by returning an error.
// These standard-library errors are caught and securely mapped to ErrInvalidX25519PublicKey.
func (k *EphemeralKey) ComputeSharedSecret(peerPubKey []byte) ([]byte, error) {
	if k.privateKey == nil {
		return nil, ErrEphemeralKeyDestroyed
	}
	if len(peerPubKey) != 32 {
		return nil, ErrMalformedX25519PublicKey
	}

	peerKey, err := ecdh.X25519().NewPublicKey(peerPubKey)
	if err != nil {
		// NewPublicKey natively rejects malformed/invalid/low-order points.
		return nil, ErrInvalidX25519PublicKey
	}

	// ECDH returns exactly a 32-byte shared secret, performing standard X25519 checks.
	secret, err := k.privateKey.ECDH(peerKey)
	if err != nil {
		return nil, ErrInvalidX25519PublicKey
	}

	// Defensive copy of the raw 32-byte shared secret to guarantee ownership boundaries
	secretCopy := make([]byte, len(secret))
	copy(secretCopy, secret)
	return secretCopy, nil
}

// Destroy performs a best-effort cleanup of the ephemeral key by dropping the reference
// to the internal private key. This is an application-level lifecycle boundary only.
// It does NOT guarantee RAM zeroization, secure erasure, or protection against
// core dumps or swap file analysis.
// This method must be called sequentially and never concurrently with cryptographic operations.
func (k *EphemeralKey) Destroy() {
	k.privateKey = nil
}
