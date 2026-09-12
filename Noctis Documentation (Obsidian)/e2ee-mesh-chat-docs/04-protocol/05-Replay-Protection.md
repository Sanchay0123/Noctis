# Replay Protection

## Two independent replay mechanisms

The protocol uses two identifiers for two different layers:

### Mesh PacketID

`packet_id` prevents repeated forwarding of the same network packet.

### Session sequence number

The application sequence number prevents replay of an authenticated
application message within a session.

Neither mechanism replaces the other.

## Sequence numbers

Each directional session key has an independent sequence space. The sender
begins at sequence `0` for a newly established session and increments once
per successfully encrypted application message.

The sender never wraps the sequence number. At `MaxUint64`, encryption
returns `ErrSequenceExhausted`.

## Nonce

The ChaCha20-Poly1305 nonce is:

```text
nonce = 0x00000000 || uint64_be(sequence_num)
```

Nonce uniqueness is guaranteed by the directional key/sequence lifecycle:
fresh session keys, fresh ephemeral handshakes, monotonic sequence state,
restart invalidation and session rotation before the project threshold.

## Receiver sliding window

The receive window is **64 messages**.

Let `highest` be the highest accepted sequence number:
- bit 0 represents `highest`
- bit N represents `highest-N`
- a duplicate set bit is rejected
- a sequence number exactly 64 below `highest` is outside the window and is
  rejected

Newer sequence numbers advance the window. Valid older sequence numbers
inside the window are accepted once; duplicates are rejected.

## Authentication ordering

A packet is not allowed to consume replay state merely because its sequence
number appears acceptable.

The implementation:
1. validates the ciphertext structure
2. checks the current replay window without mutating it
3. authenticates/decrypts with ChaCha20-Poly1305
4. commits replay-window state only after successful authentication

This prevents forged or modified packets from poisoning receiver replay
state. Concurrent duplicate arrivals are rechecked before replay-state
commit.

## Session scope

Application replay state belongs to one established session. It is not a
network-wide replay mechanism.

Session replacement creates fresh keys and fresh sequence/replay state.
Session keys and sequence state are not persisted across restart.

## Handshake replay

Handshake state handling is separate from application replay protection.
Handshake state rejects invalid or duplicate state transitions, but it does
not replace the application ciphertext replay window.

See [[04-protocol/03-Session-Establishment]].
