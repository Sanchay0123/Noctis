# Message Format

## Status

**M1.1 schema frozen and verified; M2 identity complete; M3 cryptographic
session/message protection implemented and verified.**

The Protobuf wire schema and runtime validation policy are frozen for the
current protocol version. The M3 X25519, HKDF and AEAD implementations are
complete, and M4 integrates them with direct TCP messaging.

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

## Application AEAD AAD

For M3 application encryption, the canonical AAD is:

```text
"MeshChat-AppData-v1" ||
SessionID ||
uint64_be(SequenceNumber) ||
Direction
```

Where:
- `SessionID` is the full 32-byte session identifier
- `SequenceNumber` is the 8-byte big-endian application sequence number
- `Direction` is one byte derived internally from the session role:
  `0x01` initiator -> responder, `0x00` responder -> initiator

The direction marker is not caller-controlled.

Routing metadata remains part of the outer mesh packet and is not implicitly
authenticated by this application-layer AAD construction. Any future
cryptographic binding of additional envelope fields requires an explicit
protocol decision.

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


## M3 application payload status

`AppDataPayload.sequence_num` carries the application sequence number used by
the M3 AEAD layer. `AppDataPayload.ciphertext` contains the encrypted
application bytes together with the 16-byte Poly1305 authentication tag.

The M3 implementation rejects truncated ciphertext and authenticates the
session ID, sequence number and direction context before delivering
plaintext.

Relay nodes forward the encrypted payload without access to application
plaintext or session keys.


## M4 Transport Framing

Direct TCP transport encodes each serialized `MeshPacket` as:

```text
4-byte big-endian unsigned length
+
Protobuf MeshPacket bytes
```

The maximum body length is 64 KiB. The receiver validates the declared length
before allocating the body buffer and uses complete-read semantics. The sender
serializes concurrent writes and handles short writes until the complete frame
has been transmitted or an error occurs.

## M4 Direct APP_DATA Binding

An established direct channel constructs APP_DATA with:

- `source_node` = local Ed25519 public identity
- `dest_node` = expected remote Ed25519 public identity
- `session_id` = active M3 session identifier
- `sequence_num` = M3 directional sequence number
- `ciphertext` = M3 ChaCha20-Poly1305 output

The receiver requires all three identity/session bindings to match the active
channel before delivering authenticated plaintext.
