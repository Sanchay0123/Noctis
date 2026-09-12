# Cryptographic Failure Modes

## Critical failures

### Nonce reuse

Potential catastrophic AEAD compromise.

**Required response:** reject design or implementation that can reuse a
nonce under the same key.

### Unauthenticated key exchange

Enables MITM.

**Required response:** bind X25519 material to authenticated identities.

### Custom cryptography

High probability of subtle design/implementation errors.

**Required response:** use established library primitives.

### Plaintext logging

Direct confidentiality failure.

**Required response:** remove sensitive logging and add regression
checks.

### Wrong-key acceptance

Indicates catastrophic protocol failure.

**Required response:** mandatory negative test.

## Other failures

-   invalid signature accepted
-   malformed public key accepted
-   invalid ciphertext accepted
-   changed AAD accepted
-   replayed message accepted
-   key context confusion
-   session-role confusion
-   identity mismatch silently ignored

## Review principle

Cryptographic failure handling should generally be fail-closed.


## M3 verified failure handling

The M3 implementation includes negative tests for:
- invalid/low-order X25519 peer public keys
- transcript and identity substitution
- malformed handshake state
- modified ciphertext
- modified AAD/context
- truncated ciphertext
- duplicate and outside-window sequence numbers
- cross-session and cross-direction ciphertext
- sequence-number exhaustion

The AEAD receive path authenticates ciphertext before committing replay-window
state. Cryptographic failures therefore do not advance the receiver's replay
state.

The implementation uses typed errors for the principal X25519/session failure
conditions rather than exposing secret material in diagnostics.
