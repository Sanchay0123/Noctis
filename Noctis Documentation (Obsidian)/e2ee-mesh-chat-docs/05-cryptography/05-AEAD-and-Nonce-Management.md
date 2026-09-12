# AEAD and Nonce Management

## Selected AEAD

**ChaCha20-Poly1305** is the accepted single AEAD construction for the
prototype. Its implementation is gated to the later encrypted-messaging
milestone.

## Security properties

AEAD provides:

-   confidentiality
-   ciphertext integrity
-   authentication of associated data

## Nonce invariant

> **Never reuse a nonce with the same key.**

This must be enforced by architecture, not by developer memory.

## Planned nonce strategy

Application messages use a monotonically increasing per-direction
sequence number as protocol state. The exact mapping from sequence number
to the 12-byte ChaCha20-Poly1305 nonce must be frozen before AEAD
implementation.

A nonce must never repeat under the same directional key, including across
reconnects or session rotation.

## Never do

``` text
random nonce with accidental collision assumptions
fixed nonce
counter reset under same key
reuse nonce after reconnect with same key
```

unless the final protocol provides a rigorous construction that makes
the specific behavior safe.

## AAD

AAD should authenticate metadata whose modification must invalidate the
ciphertext.

Candidate fields:

-   version
-   message type
-   sender identity
-   recipient identity
-   session identifier
-   sequence number

Final selection requires protocol review.
