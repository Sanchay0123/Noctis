# AEAD and Nonce Management

## Selected AEAD

**ChaCha20-Poly1305 — implemented / verified in M3.3**

The prototype uses a single AEAD construction to avoid unnecessary
algorithm-selection complexity.

## Nonce construction

Each directional session key has its own monotonically increasing sequence
space. The first outbound sequence number is `0`.

The 12-byte nonce is:

```text
nonce = 0x00000000 || uint64_be(sequence_num)
```

Nonce uniqueness therefore depends on the invariant that a directional key
is never reused with the same sequence number.

The sender never wraps the sequence number. When the sequence number reaches
`MaxUint64`, encryption returns `ErrSequenceExhausted` and does not continue
with a wrapped nonce.

## Associated data

The canonical application AAD is:

```text
"MeshChat-AppData-v1" ||
session_id ||
uint64_be(sequence_num) ||
direction
```

where:
- `"MeshChat-AppData-v1"` is the fixed ASCII protocol context
- `session_id` is the full 32-byte session identifier
- `sequence_num` is the 8-byte big-endian sequence number
- `direction` is one byte derived internally from the session role:
  - `0x01` initiator -> responder
  - `0x00` responder -> initiator

The caller does not supply the direction marker.

Changing any authenticated field causes AEAD authentication to fail.

## Replay protection

The receiver maintains a 64-message sliding replay window. Bit 0 represents
the current highest accepted sequence number; bit N represents
`highest-N`.

A sequence number exactly 64 behind the current highest is outside the
window and is rejected. Duplicates are rejected.

Critically, replay-window state is committed **only after successful AEAD
authentication**. A forged or modified ciphertext therefore cannot consume
receiver replay state.

## Concurrency

The session uses independent send and receive mutexes. Concurrent outbound
messages serialize sequence allocation under the send lock; concurrent
inbound messages serialize replay-state commitment under the receive lock.

## Session lifecycle

Session keys and sequence state are memory-only. Fresh X25519 handshakes are
required after restart and during session rotation.

The initial rotation policy is the earlier of:
- `2^16` application messages
- 24 hours

## Evidence

M3.3 includes a deterministic known-answer vector for nonce, AAD and
ChaCha20-Poly1305 output, plus negative tests for ciphertext modification,
AAD/context changes, truncation, duplicates, outside-window packets,
unauthenticated replay-state poisoning, cross-session use, cross-direction
use, sequence zero and sequence exhaustion.

Containerized formatting, vetting, tests, race detection and build validation
passed.
