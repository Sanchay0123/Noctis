# Replay Protection

## Message sequence

Each direction of a session maintains an independent monotonically
increasing sequence number.

A new session receives fresh session keys, so sequence numbering may
restart at zero only with those new keys.

## Nonce invariant

ChaCha20-Poly1305 requires nonce uniqueness for a given key.

The protocol invariant is:

> No `(AEAD key, nonce)` pair may ever be reused.

The proposed nonce construction may encode a 64-bit sequence number into
the 96-bit nonce, provided:

- the mapping is injective
- the directional key is unique
- sequence state never repeats under that key
- the session key is never reused after restart
- overflow is handled by terminating/rotating the session before reuse

The implementation must test this invariant.

## Sliding window

The receiver maintains a bounded replay window.

A candidate sequence number is:

1. rejected if already accepted
2. rejected if outside the permitted old-message window
3. authenticated before being committed as accepted, according to the
   final protocol state machine
4. accepted and recorded if valid

The exact window width is an implementation parameter and must be
documented.

## PacketID versus sequence number

These solve different problems:

- `PacketID` prevents mesh forwarding loops and duplicate flooding.
- `sequence_number` protects an authenticated session against message
  replay.

Do not substitute one for the other.

## Rotation

Session rotation must occur before any sequence/nonce exhaustion limit.

The proposed `2^16` message threshold is acceptable as a conservative
prototype policy but must be implemented as a policy constant rather
than an assumption in the cryptographic primitive.

See [[05-cryptography/04-AEAD-and-Nonce-Management]].
