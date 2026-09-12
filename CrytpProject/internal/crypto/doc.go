// Package crypto defines cryptographic identities, session establishment, and authenticated encryption.
// It sits below the application layer and strictly encapsulates private keys and symmetric session keys.
// Neither raw private keys nor session keys are exposed to the transport, routing, or observability layers.
//
// Cryptographic Dependency Policy:
// - Go standard crypto/* is preferred.
// - Approved golang.org/x/crypto/* modules may be used when required (e.g. HKDF).
// - Arbitrary third-party cryptographic libraries are not permitted.
// - Cryptographic primitives must not be implemented manually.
//
// M2.1 Updates:
// - Ed25519 is used as the long-term identity primitive.
// - 32-byte public keys are the public identity representation.
//
// M3.1 Updates:
// - X25519 (via standard library crypto/ecdh) is implemented for ephemeral key agreement.
// - Ephemeral keys enforce exactly 32-byte private keys, public keys, and shared secrets.
//
// M3.2 Updates:
// - Authenticated session establishment implemented.
// - Exact T_INIT and T_RESP transcript bindings are signed and verified.
// - Session ID is strictly SHA256(T_RESP).
// - HKDF-SHA-256 (golang.org/x/crypto/hkdf) derives directional keys based on role.
// - State machine strictly prevents replay, ordering, and transition failures.
// - Failed authentications do not leak key material or establish sessions.
//
// M3.3 Updates (AEAD):
// - Message encryption uses golang.org/x/crypto/chacha20poly1305.
// - Nonce (12 bytes) = [4 bytes zero padding] || [8 bytes big-endian sequence number].
// - AAD = "MeshChat-AppData-v1" || session_id (32 bytes) || sequence_number (8 bytes big-endian) || direction marker (1 byte).
// - Direction marker: 0x01 (Initiator -> Responder), 0x00 (Responder -> Initiator).
// - Sequence numbers strictly start at 0, increment by 1, and panic/exhaust on math.MaxUint64 (never wrap).
// - Replay protection: 64-message explicit sliding window using a uint64 bitmap. Unauthenticated packets do not advance window state.
//
// Note: Networking handshakes and mesh routing are NOT implemented.
package crypto
