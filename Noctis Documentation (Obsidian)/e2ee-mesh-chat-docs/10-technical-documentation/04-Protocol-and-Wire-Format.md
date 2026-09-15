# Protocol and Wire Format — Implemented Version 1

## Packet schema

The protocol uses Protocol Buffers.

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
  bytes packet_id = 3;
  uint32 ttl = 4;
  bytes source_node = 5;
  bytes dest_node = 6;
  oneof payload {
    InitPayload init = 7;
    RespPayload resp = 8;
    AppDataPayload app_data = 9;
  }
}

message InitPayload {
  bytes ephemeral_key = 1;
  bytes signature = 2;
}

message RespPayload {
  bytes ephemeral_key = 1;
  bytes signature = 2;
}

message AppDataPayload {
  bytes session_id = 1;
  uint64 sequence_num = 2;
  bytes ciphertext = 3;
}
```

## Fixed-size invariants

| Field | Required size |
|---|---:|
| `packet_id` | 16 bytes |
| `source_node` | 32 bytes |
| `dest_node` | 32 bytes |
| `InitPayload.ephemeral_key` | 32 bytes |
| `RespPayload.ephemeral_key` | 32 bytes |
| signatures | 64 bytes |
| `session_id` | 32 bytes |
| APP_DATA ciphertext | AEAD ciphertext plus 16-byte tag |

## Packet-type/oneof consistency

The packet type and selected Protobuf `oneof` must agree:

```text
INIT      -> init
RESP      -> resp
APP_DATA  -> app_data
UNKNOWN   -> reject
```

Unsupported protocol versions are rejected.

## TCP framing

The transport uses a four-byte big-endian length prefix followed by the
serialized Protobuf frame.

The declared length is checked before allocation and must not exceed **64 KiB**.

This prevents an attacker-controlled length prefix from directly causing an
unbounded frame allocation.

## Routing metadata

The outer packet contains routing metadata such as source identity,
destination identity, packet ID and TTL. These fields are intentionally
available to intermediate routing nodes.

The application ciphertext is not decrypted by the router.

## TTL semantics

```text
if destination == local node:
    process according to packet type/state
else if ttl <= 1:
    drop
else:
    decrement ttl
    forward
```

Thus a destination may process a packet arriving with TTL 1, while a
non-destination cannot forward a packet whose TTL is 1.

## Duplicate suppression

The routing duplicate cache uses:

```text
(source_node, packet_id)
```

as its cache key.

`packet_id` is a 16-byte random identifier. The implementation does not claim
that a PacketID is mathematically or globally unique.

The bounded cache prevents repeated forwarding of the same source/packet pair
within its retention window.

## Validation boundary

Before expensive routing or application processing, the implementation
validates the frame and packet structure, including:

1. maximum frame size
2. Protobuf decoding policy
3. protocol version
4. packet type
5. fixed-size identifiers
6. payload/oneof consistency
7. packet-specific constraints
8. destination and session identifiers where applicable
9. ciphertext structural constraints

Unknown Protobuf fields are rejected under the version-1 runtime validation
policy.

## Security interpretation

The wire format intentionally separates:

```text
Routing envelope  = visible to relays
Application data  = encrypted endpoint payload
```

This provides relay-compatible forwarding without requiring relays to possess
endpoint application session keys.

## Related records

- [04-Message-Format](../04-protocol/04-Message-Format.md)
- [05-Mesh-Packet-Format](../04-protocol/05-Mesh-Packet-Format.md)
- [07-Protocol-State-Machines](../04-protocol/07-Protocol-State-Machines.md)
- [03-Routing](../06-networking/03-Routing.md)
- [16-M9-Attack-Evidence](../07-testing/16-M9-Attack-Evidence.md)
