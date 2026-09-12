# Identity and Signatures

## Ed25519 role

Ed25519 provides:

-   long-term node identity
-   digital signatures
-   authentication of identity-controlled protocol data

It does **not** provide:

-   encryption
-   key agreement
-   anonymity

## M2 implementation status

The long-term Ed25519 identity core is implemented and audited.

Validated scope includes:

-   cryptographically random Ed25519 key generation
-   32-byte public identity representation
-   private-key encapsulation
-   defensive copies for caller-provided private-key material and returned
    public-key material
-   malformed-input rejection
-   deterministic Ed25519 signing and signature verification
-   concurrent signing validation with the race detector

The implementation uses Go's standard `crypto/ed25519` and `crypto/rand`
packages and does not implement Ed25519 primitives from scratch.

## Identity tests

Implemented tests cover:

-   generate identity
-   load identity
-   sign message
-   verify valid signature
-   reject modified message
-   reject invalid signature
-   reject malformed public key
-   defensive ownership behavior
-   concurrent signing behavior

## Key storage

Private identity keys must be treated as high-sensitivity local secrets.

Persistent storage, encryption at rest, permissions, backup implications,
and identity rotation/replacement remain deployment/protocol work outside
the completed M2 in-memory identity primitive.

No private key should appear in ordinary logs or telemetry.
