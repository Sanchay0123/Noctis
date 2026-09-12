# Mesh Packet Format

## Purpose

The mesh envelope provides enough information for intermediate nodes to
route a packet without providing them the ability to decrypt application
content.

## Logical structure

``` text
MeshEnvelope
├── protocol_version
├── packet_type
├── packet_id
├── source_id
├── destination_id
├── TTL
├── routing metadata
└── E2EE payload
```

## Routing-visible information

Likely candidates:

-   protocol version
-   packet type
-   packet identifier
-   source/destination routing identifiers
-   TTL
-   route metadata

## E2EE-protected information

-   application plaintext
-   application message content
-   cryptographic session material

## Important distinction

The destination identifier may need to be visible for basic routing. If
so, this is metadata leakage and must be documented.

## Serialization decision

The M1.1 implementation froze Protobuf 3 as the wire serialization format.
The canonical handshake transcript is a separate byte-level construction
defined in [[04-protocol/03-Session-Establishment]] and is not replaced by
ordinary Protobuf serialization.

Runtime validation remains mandatory for fixed-width fields, packet-type
and `oneof` consistency, protocol version, unknown fields, and frame limits.
