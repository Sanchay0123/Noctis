# Message Format

## Status

**Proposed wire schema — requires implementation review before freezing.**

Protocol Buffers is the selected serialization format.

## Design principles

- fixed protocol version
- explicit packet type
- bounded fields
- binary identifiers represented as `bytes`
- no application plaintext in relay-visible fields
- cryptographic metadata included in authenticated context
- handshake messages defined separately from encrypted application data

## Draft schema

```protobuf
syntax = "proto3";

message MeshPacket {
  uint32 protocol_version = 1;
  PacketType packet_type = 2;

  // Routing metadata.
  bytes packet_id = 3;       // fixed 16-byte identifier
  uint32 ttl = 4;
  bytes source_node = 5;     // Ed25519 public key, fixed 32 bytes
  bytes dest_node = 6;       // Ed25519 public key, fixed 32 bytes

  // Session/message metadata.
  bytes session_id = 7;      // exact size/derivation to be finalized
  uint64 sequence_num = 8;

  // Opaque E2EE application ciphertext + authentication tag.
  bytes ciphertext = 9;
}

enum PacketType {
  PACKET_TYPE_UNSPECIFIED = 0;
  HANDSHAKE_INIT = 1;
  HANDSHAKE_RESP = 2;
  ENCRYPTED_MESSAGE = 3;
}
```

This is a design draft, not the final schema.

## Important corrections

Do not encode Ed25519 public keys as base64 strings inside Protobuf
unless interoperability requirements justify it. `bytes` avoids
unnecessary encoding/decoding and makes fixed-size validation clearer.

The exact `session_id` representation must be finalized.

## AAD

For an encrypted application packet, the canonical AAD must include at
least:

```text
protocol_version
source_node
dest_node
session_id
sequence_num
```

The exact byte encoding must be specified so both endpoints authenticate
the same bytes.

## Routing mutation

A relay may decrement TTL.

Any fields covered by AAD must be treated as cryptographically
authenticated. A relay modifying them must cause destination
authentication failure.

If TTL is deliberately excluded from AAD, its mutation remains a
routing-layer operation and must be protected by routing validation
rather than application AEAD.

This distinction must be explicit in the final protocol specification.

## Handshake packets

Handshake packets must define exact Protobuf messages for:

- INIT
- RESP

Do not place an underspecified binary blob in a generic `ciphertext`
field and call it a protocol.

## Validation

Before cryptographic processing:

- validate protocol version
- validate packet type
- validate required fields
- validate fixed-size identifiers
- validate TTL
- validate field/packet size limits
- reject malformed input safely
