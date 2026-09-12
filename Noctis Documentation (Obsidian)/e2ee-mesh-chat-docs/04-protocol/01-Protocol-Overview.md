# Protocol Overview

## Status

**Architecture approved with implementation gates.**

The protocol separates:

1. routing metadata
2. cryptographic/session metadata
3. opaque E2EE application content
4. handshake messages

Relays may inspect and modify routing metadata as required by forwarding,
but they must never receive application plaintext or session secrets.

## Packet classes

The protocol shall define an explicit packet type, including at least:

- handshake initiation
- handshake response
- encrypted application message

The exact numeric assignments belong in the versioned Protobuf schema.

## Canonicalization requirement

Every cryptographically signed or authenticated structure must have one
unambiguous byte representation.

Do not sign ad-hoc concatenated strings whose encoding can vary.

The implementation must define:

- field order
- length encoding
- protocol version
- domain-separation label
- identity encoding
- ephemeral-key encoding
- transcript encoding

## Handshake transcript

The authenticated exchange must bind:

```text
protocol_version
handshake_domain
initiator_identity
responder_identity
initiator_ephemeral_X25519_public_key
responder_ephemeral_X25519_public_key
```

Both identity keys remain Ed25519 keys. They are not reused as X25519
keys.

The final transcript construction must be documented in
[[04-protocol/03-Session-Establishment]] before M3 implementation.

## Encryption context

The encrypted message authentication context must bind the message to:

```text
protocol_version
sender_identity
receiver_identity
session_id
sequence_number
```

The exact canonical AAD encoding is part of the protocol specification.

## Versioning

A packet with an unsupported protocol version must be rejected safely.

Unknown fields must not change the cryptographic meaning of a message.

## Size limits

All packet and field sizes must have explicit upper bounds before
deserialization or expensive cryptographic processing.

## Relay invariant

A relay forwards opaque ciphertext.

It must never need:

- the recipient's session key
- the sender's session key
- application plaintext

See [[02-architecture/07-Observability-and-Security-Telemetry]] for the
telemetry boundary.

## M2 Protocol Status

M2 does not alter the frozen wire protocol. It supplies the long-term Ed25519
identity primitive that later authenticated handshakes will use.

No X25519 ephemeral keys, handshake transcripts, session IDs, or session keys
are implemented at this stage.
