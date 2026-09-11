# Replay Protection

## Two independent replay mechanisms

The protocol uses two different identifiers for two different layers.

### Mesh PacketID

Prevents repeated forwarding of the same network packet.

### Session sequence number

Prevents replay of an authenticated application message.

Neither mechanism replaces the other.

## Sequence numbers

Each directional session key has an independent sequence space.

The sender begins at sequence `0` for a newly established session and
increments exactly once for each encrypted application message.

## Nonce

The AEAD nonce is a deterministic 96-bit encoding of the sequence number:

```text
nonce = 0x00000000 || uint64_be(sequence_num)
```

This construction is safe only under the `(key, nonce)` uniqueness
invariant.

## Nonce uniqueness invariant

The implementation must guarantee:

> A given ChaCha20-Poly1305 key is never used with the same nonce twice.

This is achieved by:

- separate directional keys
- fresh X25519 ephemeral keys for every new session
- in-memory-only session state
- mandatory fresh handshake after restart
- monotonic sequence numbers
- session replacement before sequence exhaustion

## Receiver sliding window

The receiver maintains a bounded replay window.

For each incoming sequence number:

1. reject if already accepted
2. reject if too old for the window
3. authenticate/decrypt according to the protocol state machine
4. commit the sequence number as accepted only according to the defined
   successful-authentication path

The window width is initially `64`.

## Session rotation

Rotate/re-establish session keys before sequence-number exhaustion.

The initial policy is:

- maximum `2^16` application messages, or
- maximum 24 hours

whichever occurs first.

The implementation must not reuse the old key/nonce space.

## Restart

Session keys and sequence state are not persisted.

After restart:

```text
old session state = invalid
        ↓
fresh X25519 handshake
        ↓
fresh session keys
        ↓
sequence starts at 0 under new keys
```

This is a security requirement, not merely a convenience.

## Handshake replay

Handshake messages require their own duplicate/replay state.

A fresh ephemeral key does not by itself prove that an incoming handshake
is new.

See [[04-protocol/03-Session-Establishment]].
