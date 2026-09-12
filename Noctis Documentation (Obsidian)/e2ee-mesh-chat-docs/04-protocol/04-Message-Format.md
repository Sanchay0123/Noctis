# Message Format

## Status

**M1.1 schema frozen and verified; M2 identity complete; M3 cryptographic
implementation gated.**

The Protobuf wire schema and runtime validation policy are frozen for the
current protocol version. This does not mean that X25519, HKDF, or AEAD
implementation exists yet.

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

## M1.1 Implementation Freeze — Protobuf Schema & Code Generation

The M1.1 implementation establishes the canonical Protobuf source at
`internal/protocol/meshchat.proto` and generated Go output at
`internal/protocol/meshchat.pb.go`.

The generated-code workflow is containerized through
`scripts/generate_proto.sh`; host-installed Protobuf tooling is not required.
The implementation report records `protoc-gen-go` v1.33.0 and the Go 1.21
Alpine container as the controlled generation environment.

### Runtime validation invariants

The Protobuf schema does not itself enforce all security-sensitive byte
lengths. Runtime validation must enforce the documented exact sizes for
packet IDs, node identities, ephemeral keys, signatures, and session IDs.

The runtime packet validator must also enforce consistency between
`PacketType` and the selected `oneof` payload:

- `PACKET_TYPE_INIT` -> `init`
- `PACKET_TYPE_RESP` -> `resp`
- `PACKET_TYPE_APP_DATA` -> `app_data`
- `PACKET_TYPE_UNKNOWN` -> reject

Unknown Protobuf fields are not accepted by the protocol. The implementation
currently validates this policy at runtime after unmarshaling; this behavior
must remain covered by an explicit regression test.

### Generation validation

Changes to the `.proto` source must be followed by deterministic regeneration
of the generated Go code and the normal formatting, vet, test, and build
validation sequence.
