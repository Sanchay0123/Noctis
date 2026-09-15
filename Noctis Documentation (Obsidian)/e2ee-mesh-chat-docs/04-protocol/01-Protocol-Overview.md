# Protocol Overview

## Status

**Architecture approved; M4 direct secure messaging integration implemented, verified and approved.**

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

The canonical transcript construction is frozen and implemented in
[03-Session-Establishment](03-Session-Establishment.md).

## Encryption context

The encrypted message authentication context uses the canonical M3 AAD:

```text
"MeshChat-AppData-v1" ||
session_id ||
uint64_be(sequence_number) ||
direction
```

The direction marker is derived from the established session role. Additional
outer routing fields are not implicitly part of the M3 AAD construction.

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

See [07-Observability-and-Security-Telemetry](../02-architecture/07-Observability-and-Security-Telemetry.md) for the
telemetry boundary.

## M3 Protocol Status

M3 implements the authenticated session and application message-protection
layer over the frozen protocol baseline.

Implemented:
- fresh ephemeral X25519 key agreement
- Ed25519-authenticated canonical handshake transcripts
- SHA-256 session identifiers
- HKDF-SHA-256 directional session keys
- ChaCha20-Poly1305 application ciphertext
- sequence-derived nonces
- role-derived direction markers
- 64-message authenticated replay window

M5 direct networking runtime integration is implemented and approved. Mesh forwarding remains a future implementation stage under M6.


## M4 Direct Transport Integration

M4 carries the frozen protocol packets over a direct TCP connection using a
4-byte big-endian length prefix. The frame length is limited to 64 KiB and is
validated before allocation.

The direct channel uses the existing INIT/RESP handshake state machines. After
the M3 session is established, application messages are represented as
`PACKET_TYPE_APP_DATA` packets. APP_DATA contains the session identifier,
application sequence number and AEAD ciphertext.

For M4, a direct TCP channel owns one established cryptographic session. This
is an implementation scope decision; it does not change the protocol's future
ability to support session multiplexing.

Transport and protocol validation occur before application plaintext delivery.
Unknown fields, unsupported versions, malformed payloads and invalid
fixed-width fields are rejected.
