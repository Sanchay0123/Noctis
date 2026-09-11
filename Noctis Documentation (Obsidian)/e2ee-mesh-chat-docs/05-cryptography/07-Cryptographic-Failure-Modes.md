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
