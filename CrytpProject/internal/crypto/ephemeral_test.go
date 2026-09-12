package crypto

import (
	"bytes"
	"crypto/ecdh"
	"encoding/hex"
	"testing"
)

func TestEphemeralKeyGeneration(t *testing.T) {
	key1, err := GenerateEphemeralKey()
	if err != nil {
		t.Fatalf("Failed to generate ephemeral key: %v", err)
	}

	pub1, err := key1.PublicKey()
	if err != nil {
		t.Fatalf("Failed to get public key: %v", err)
	}
	if len(pub1) != 32 {
		t.Errorf("Expected 32-byte public key, got %d", len(pub1))
	}

	key2, _ := GenerateEphemeralKey()
	pub2, _ := key2.PublicKey()

	if bytes.Equal(pub1, pub2) {
		t.Fatalf("Repeated generation produced identical keypairs")
	}
}

func TestSharedSecretAgreement(t *testing.T) {
	alice, _ := GenerateEphemeralKey()
	bob, _ := GenerateEphemeralKey()

	alicePub, _ := alice.PublicKey()
	bobPub, _ := bob.PublicKey()

	aliceSecret, err := alice.ComputeSharedSecret(bobPub)
	if err != nil {
		t.Fatalf("Alice failed to compute shared secret: %v", err)
	}

	bobSecret, err := bob.ComputeSharedSecret(alicePub)
	if err != nil {
		t.Fatalf("Bob failed to compute shared secret: %v", err)
	}

	if len(aliceSecret) != 32 {
		t.Errorf("Expected 32-byte shared secret, got %d", len(aliceSecret))
	}

	if !bytes.Equal(aliceSecret, bobSecret) {
		t.Errorf("Alice and Bob computed different shared secrets")
	}
}

func TestRFC7748TestVector(t *testing.T) {
	// RFC 7748 section 5.2 - Test Vector 1
	// Proves that the underlying standard library ecdh.X25519() mathematically matches the specification.
	alicePrivHex := "77076d0a7318a57d3c16c17251b26645df4c2f87ebc0992ab177fba51db92c2a"
	bobPubHex := "de9edb7d7b7dc1b4d35b61c2ece435373f8343c85b78674dadfc7e146f882b4f"
	expectedSharedHex := "4a5d9d5ba4ce2de1728e3bf480350f25e07e21c947d19e3376f09b3c1e161742"

	alicePrivBytes, _ := hex.DecodeString(alicePrivHex)
	bobPubBytes, _ := hex.DecodeString(bobPubHex)
	expectedShared, _ := hex.DecodeString(expectedSharedHex)

	aliceKey, err := ecdh.X25519().NewPrivateKey(alicePrivBytes)
	if err != nil {
		t.Fatalf("Failed to load test vector private key: %v", err)
	}

	bobKey, err := ecdh.X25519().NewPublicKey(bobPubBytes)
	if err != nil {
		t.Fatalf("Failed to load test vector public key: %v", err)
	}

	shared, err := aliceKey.ECDH(bobKey)
	if err != nil {
		t.Fatalf("ECDH failed: %v", err)
	}

	if !bytes.Equal(shared, expectedShared) {
		t.Errorf("Shared secret did not match RFC 7748 test vector")
	}
}

func TestInvalidPeerPublicKey(t *testing.T) {
	alice, _ := GenerateEphemeralKey()

	// 1. Wrong length
	badLenPub := make([]byte, 31)
	if _, err := alice.ComputeSharedSecret(badLenPub); err != ErrMalformedX25519PublicKey {
		t.Errorf("Expected ErrMalformedX25519PublicKey, got: %v", err)
	}

	// 2. Low-order point (e.g., all zeros)
	// crypto/ecdh strictly rejects the all-zero array and other known low-order points for X25519.
	zeroPub := make([]byte, 32)
	if _, err := alice.ComputeSharedSecret(zeroPub); err != ErrInvalidX25519PublicKey {
		t.Errorf("Expected low-order all-zero point to be rejected with ErrInvalidX25519PublicKey, got: %v", err)
	}

	// 3. Nil input
	if _, err := alice.ComputeSharedSecret(nil); err != ErrMalformedX25519PublicKey {
		t.Errorf("Expected nil input to be handled gracefully, got: %v", err)
	}
}

func TestEphemeralKeyOwnership(t *testing.T) {
	alice, _ := GenerateEphemeralKey()
	pub, _ := alice.PublicKey()

	// Mutate returned public key
	pub[0] ^= 0xFF

	pub2, _ := alice.PublicKey()
	if bytes.Equal(pub, pub2) {
		t.Errorf("Public key was corrupted by caller slice mutation (defensive copy failed)")
	}

	bob, _ := GenerateEphemeralKey()
	bobPub, _ := bob.PublicKey()

	secret, _ := alice.ComputeSharedSecret(bobPub)

	// Mutate returned shared secret
	secret[0] ^= 0xFF

	secret2, _ := alice.ComputeSharedSecret(bobPub)
	if bytes.Equal(secret, secret2) {
		t.Errorf("Shared secret defensive copy failed: mutating returned slice affected subsequent computations")
	}
}

func TestEphemeralKeyCleanup(t *testing.T) {
	alice, _ := GenerateEphemeralKey()
	alice.Destroy()

	if _, err := alice.PublicKey(); err != ErrEphemeralKeyDestroyed {
		t.Errorf("Expected ErrEphemeralKeyDestroyed after cleanup")
	}

	validPub := make([]byte, 32)
	validPub[0] = 9 // mathematically non-zero, but length check happens first anyway
	if _, err := alice.ComputeSharedSecret(validPub); err != ErrEphemeralKeyDestroyed {
		t.Errorf("Expected ErrEphemeralKeyDestroyed after cleanup")
	}
}
