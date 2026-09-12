package crypto

import (
	"bytes"
	"crypto/ed25519"
	"testing"
)

func TestIdentityGeneration(t *testing.T) {
	id, err := GenerateIdentity()
	if err != nil {
		t.Fatalf("Failed to generate identity: %v", err)
	}

	pub := id.PublicKey()
	if len(pub) != ed25519.PublicKeySize {
		t.Errorf("Generated public key has incorrect size: got %d, expected %d", len(pub), ed25519.PublicKeySize)
	}
	if len(id.privateKey) != ed25519.PrivateKeySize {
		t.Errorf("Generated private key has incorrect size")
	}

	// Verify that the explicitly exposed public key matches the one derived from the private key
	derivedPub := id.privateKey.Public().(ed25519.PublicKey)
	if !bytes.Equal(pub, derivedPub) {
		t.Errorf("Public key does not correctly map to the private key")
	}

	id2, err := GenerateIdentity()
	if err != nil {
		t.Fatalf("Failed to generate second identity: %v", err)
	}
	if bytes.Equal(pub, id2.PublicKey()) {
		t.Fatalf("Two independently generated identities produced the same public key")
	}
}

func TestIdentitySigningAndVerification(t *testing.T) {
	idA, _ := GenerateIdentity()
	idB, _ := GenerateIdentity()

	msg := []byte("hello mesh")

	sig, err := idA.Sign(msg)
	if err != nil {
		t.Fatalf("Signing failed: %v", err)
	}

	if len(sig) != ed25519.SignatureSize {
		t.Errorf("Signature has incorrect size: got %d, expected %d", len(sig), ed25519.SignatureSize)
	}

	// 1. Valid signature verifies
	if err := VerifySignature(idA.PublicKey(), msg, sig); err != nil {
		t.Errorf("Valid signature failed verification: %v", err)
	}

	// 2. Modified message fails verification
	modifiedMsg := []byte("hello mesh!")
	if err := VerifySignature(idA.PublicKey(), modifiedMsg, sig); err != ErrInvalidSignature {
		t.Errorf("Modified message should fail verification")
	}

	// 3. Signature from identity A fails under identity B's public key
	if err := VerifySignature(idB.PublicKey(), msg, sig); err != ErrInvalidSignature {
		t.Errorf("Signature from A verified successfully against B's public key")
	}

	// 4. Repeated deterministic signing
	sig2, _ := idA.Sign(msg)
	if !bytes.Equal(sig, sig2) {
		t.Errorf("Ed25519 signatures must be deterministic for the same message and private key")
	}
}

func TestMalformedInputs(t *testing.T) {
	id, _ := GenerateIdentity()
	msg := []byte("test")
	sig, _ := id.Sign(msg)

	// Malformed public key rejected
	badPubKey := make([]byte, 31)
	if err := VerifySignature(badPubKey, msg, sig); err != ErrMalformedPublicKey {
		t.Errorf("Expected ErrMalformedPublicKey, got: %v", err)
	}

	// Malformed signature rejected
	badSig := make([]byte, 63)
	if err := VerifySignature(id.PublicKey(), msg, badSig); err != ErrMalformedSignature {
		t.Errorf("Expected ErrMalformedSignature, got: %v", err)
	}

	// Nil and Empty messages are natively supported and mathematically deterministic
	nilSig, err := id.Sign(nil)
	if err != nil {
		t.Errorf("Failed to sign nil message: %v", err)
	}
	emptySig, err := id.Sign([]byte{})
	if err != nil {
		t.Errorf("Failed to sign empty message: %v", err)
	}
	if !bytes.Equal(nilSig, emptySig) {
		t.Errorf("Nil and empty message signatures must mathematically match")
	}
	if err := VerifySignature(id.PublicKey(), nil, nilSig); err != nil {
		t.Errorf("Failed to verify nil message: %v", err)
	}
	if err := VerifySignature(id.PublicKey(), []byte{}, emptySig); err != nil {
		t.Errorf("Failed to verify empty message: %v", err)
	}

	// Malformed private key loading rejected
	badPrivKey := make([]byte, 63)
	if _, err := LoadIdentity(badPrivKey); err != ErrMalformedPrivateKey {
		t.Errorf("Expected ErrMalformedPrivateKey when loading identity, got: %v", err)
	}
}

func TestMismatchedKeyMaterialSafelyHandled(t *testing.T) {
	idA, _ := GenerateIdentity()

	// LoadIdentity automatically derives the public key from the private key,
	// guaranteeing they always match. Passing just the private key prevents mismatch.
	loadedId, err := LoadIdentity(idA.privateKey)
	if err != nil {
		t.Fatalf("Failed to load valid identity: %v", err)
	}

	if !bytes.Equal(loadedId.PublicKey(), idA.PublicKey()) {
		t.Errorf("Loaded identity public key does not match original")
	}
}

func TestPublicKeyDefensiveCopy(t *testing.T) {
	id, _ := GenerateIdentity()
	pub1 := id.PublicKey()

	// Mutate returned slice
	pub1[0] ^= 0xFF

	pub2 := id.PublicKey()
	if bytes.Equal(pub1, pub2) {
		t.Errorf("Public key defensive copy failed: mutation affected internal state")
	}
	if len(pub2) != ed25519.PublicKeySize {
		t.Errorf("Public key size incorrect")
	}
}

func TestPrivateKeyLoadingMutation(t *testing.T) {
	// crypto/ed25519.GenerateKey(nil) uses crypto/rand by default if nil is passed
	pub, originalPriv, _ := ed25519.GenerateKey(nil)
	privCopy := make([]byte, len(originalPriv))
	copy(privCopy, originalPriv)

	id, err := LoadIdentity(privCopy)
	if err != nil {
		t.Fatalf("Failed to load identity: %v", err)
	}

	// Mutate caller's slice
	privCopy[0] ^= 0xFF

	// Ensure identity still works with the original correct private key and public key
	msg := []byte("test mutation")
	sig, _ := id.Sign(msg)

	if err := VerifySignature(pub, msg, sig); err != nil {
		t.Errorf("Loaded identity was corrupted by caller slice mutation")
	}
}
