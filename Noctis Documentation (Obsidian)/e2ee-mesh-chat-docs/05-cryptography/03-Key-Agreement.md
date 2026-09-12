# Key Agreement

## X25519

X25519 provides shared secret agreement between authenticated endpoints.

It does not, by itself, authenticate who owns a public key. The MeshChat
session protocol therefore authenticates the ephemeral exchange with
long-term Ed25519 signatures.

## M3 implementation

**Status: Implemented / verified**

The implementation uses Go's standard `crypto/ecdh.X25519()` API.

Properties:
- a fresh ephemeral X25519 keypair is generated for each new session
- the public key is exactly 32 bytes
- the shared secret is exactly 32 bytes
- long-term Ed25519 keys are never reused as X25519 private keys
- malformed and invalid/low-order peer public keys are rejected
- shared-secret material is returned only as the input to the session KDF
- the ephemeral private key has an explicit `Destroy()` lifecycle boundary

`Destroy()` is a best-effort memory-lifecycle measure. It must not be described
as guaranteed RAM zeroization. The key object is not concurrency-safe with
destruction; callers must use it sequentially with respect to `Destroy()`.

## Conceptual exchange

```text
Alice:
  Ed25519 identity A
  X25519 ephemeral a

Bob:
  Ed25519 identity B
  X25519 ephemeral b

Alice computes X25519(a, B_ephemeral)
Bob computes X25519(b, A_ephemeral)

Both derive the same shared secret.

Ed25519 signatures bind the ephemeral exchange to identities.
```

## Security boundary

X25519 output is raw shared-secret material. It is not used directly as an
application encryption key. M3.2 passes it into the frozen HKDF-SHA-256
schedule with the transcript hash as salt.

## Evidence

M3.1 testing covers key generation, agreement, an RFC7748 known-answer vector,
invalid/low-order public keys, ownership/defensive-copy behavior and cleanup
semantics. Containerized formatting, vetting, tests, race detection and build
validation passed.
