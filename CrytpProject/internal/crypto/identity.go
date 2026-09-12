package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"errors"
)

var (
	// ErrMalformedPublicKey is returned when a public key is not exactly 32 bytes.
	ErrMalformedPublicKey = errors.New("malformed public key: must be exactly 32 bytes")
	// ErrMalformedPrivateKey is returned when a private key is not exactly 64 bytes.
	ErrMalformedPrivateKey = errors.New("malformed private key: must be exactly 64 bytes")
	// ErrMalformedSignature is returned when a signature is not exactly 64 bytes.
	ErrMalformedSignature = errors.New("malformed signature: must be exactly 64 bytes")
	// ErrInvalidSignature is returned when signature verification mathematically fails.
	ErrInvalidSignature = errors.New("invalid signature")
)

// NodeIdentity represents the long-term cryptographic identity of a node.
// It strictly encapsulates the private Ed25519 key, ensuring it cannot be
// accessed directly by other layers.
type NodeIdentity struct {
	privateKey ed25519.PrivateKey
	publicKey  ed25519.PublicKey
}

// GenerateIdentity securely creates a new NodeIdentity using crypto/rand.
// It returns a fresh identity guaranteeing the public key correctly corresponds to the private key.
func GenerateIdentity() (*NodeIdentity, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	return &NodeIdentity{
		privateKey: priv,
		publicKey:  pub,
	}, nil
}

// LoadIdentity instantiates a NodeIdentity from an existing 64-byte private key.
// The public key is deterministically derived from the private key to prevent mismatch vulnerabilities.
func LoadIdentity(privateKey []byte) (*NodeIdentity, error) {
	if len(privateKey) != ed25519.PrivateKeySize {
		return nil, ErrMalformedPrivateKey
	}

	// 1. Make a defensive copy BEFORE any cryptographic derivation to guarantee ownership
	privCopy := make([]byte, len(privateKey))
	copy(privCopy, privateKey)

	// 2. Construct the internal ed25519.PrivateKey from the defensive copy
	internalPriv := ed25519.PrivateKey(privCopy)

	// 3. Derive the public key from the safely copied internal memory
	derivedPub := internalPriv.Public().(ed25519.PublicKey)

	// 4. Make a defensive copy of the derived public key to ensure it shares no backing memory
	pubCopy := make([]byte, len(derivedPub))
	copy(pubCopy, derivedPub)

	return &NodeIdentity{
		privateKey: internalPriv,
		publicKey:  ed25519.PublicKey(pubCopy),
	}, nil
}

// PublicKey returns the exact 32-byte public key representation of this identity.
func (id *NodeIdentity) PublicKey() []byte {
	// Return a copy to prevent mutation of the internal state
	pub := make([]byte, len(id.publicKey))
	copy(pub, id.publicKey)
	return pub
}

// Sign deterministically signs the given message using the identity's private key.
// It returns exactly a 64-byte signature. Nil or empty messages are safely handled by the primitive.
func (id *NodeIdentity) Sign(message []byte) ([]byte, error) {
	// ed25519.Sign is strictly deterministic for a given private key and message.
	return ed25519.Sign(id.privateKey, message), nil
}

// VerifySignature validates that a signature was produced by the identity
// corresponding to the given public key over the exact message.
func VerifySignature(pubKey []byte, message []byte, signature []byte) error {
	if len(pubKey) != ed25519.PublicKeySize {
		return ErrMalformedPublicKey
	}
	if len(signature) != ed25519.SignatureSize {
		return ErrMalformedSignature
	}

	if !ed25519.Verify(pubKey, message, signature) {
		return ErrInvalidSignature
	}

	return nil
}
