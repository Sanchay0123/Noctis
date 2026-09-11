# Message Format

## Status

**Draft — M1 schema work authorized; cryptographic fields must not be
considered frozen until M3 review.**

## Protobuf schema

```protobuf
syntax = "proto3";

enum PacketType {
  PACKET_TYPE_UNKNOWN = 0;
  PACKET_TYPE_INIT = 1;
  PACKET_TYPE_RESP = 2;
  PACKET_TYPE_APP_DATA = 3;
}

message MeshPacket {
  uint32 version = 1;
  PacketType type = 2;

  // Routing metadata.
  bytes packet_id = 3;       // exactly 16 bytes
  uint32 ttl = 4;
  bytes source_node = 5;     // exactly 32 bytes
  bytes dest_node = 6;       // exactly 32 bytes

  oneof payload {
    InitPayload init = 7;
    RespPayload resp = 8;
    AppDataPayload app_data = 9;
  }
}

message InitPayload {
  bytes ephemeral_key = 1;   // exactly 32 bytes
  bytes signature = 2;       // exactly 64 bytes
}

message RespPayload {
  bytes ephemeral_key = 1;   // exactly 32 bytes
  bytes signature = 2;       // exactly 64 bytes
}

message AppDataPayload {
  bytes session_id = 1;      // exactly 32 bytes
  uint64 sequence_num = 2;
  bytes ciphertext = 3;      // ciphertext + 16-byte AEAD tag
}
```

## Session identifier

`session_id` is the SHA-256 hash of the canonical RESP transcript and is
therefore exactly 32 bytes.

Do not truncate it to four bytes.

## Validation

Before expensive processing:

1. enforce maximum TCP frame size
2. parse Protobuf
3. validate protocol version
4. validate packet type
5. validate required fields
6. validate fixed-width identifier sizes
7. validate payload structure
8. validate protocol-specific limits

## Unknown versions

An unsupported `version` must be rejected.

## Unknown fields

Do not assume the Protobuf library automatically rejects unknown fields.

The implementation must make an explicit decision:

- reject unknown fields through a validation policy, or
- tolerate them according to a documented compatibility policy

For protocol version 1, security-sensitive canonicalized structures must
not acquire alternate interpretations because of ignored fields.

## TTL

Routing semantics:

```text
if destination == local_node:
    process according to packet type/state
else:
    if ttl <= 1:
        drop
    else:
        ttl = ttl - 1
        forward
```

The destination must not be accidentally dropped merely because its
incoming TTL is `1`.

## AAD

For application data, the canonical AAD is:

```text
ProtocolVersion ||
PacketType ||
SessionID ||
SequenceNumber ||
SourceNode ||
DestNode
```

All fields must have deterministic fixed-width encodings.

## Relay behavior

A relay only needs to inspect the routing envelope required for
forwarding and duplicate/TTL processing.

It must not access:

- application plaintext
- session keys
- private keys

## Frame limit

The initial maximum frame size may be 64 KiB.

The 4-byte TCP length prefix must be checked against this limit **before
allocating the frame buffer**.

Oversized frames are rejected and must not trigger unbounded memory
allocation.
