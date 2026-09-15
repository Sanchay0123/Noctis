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
defined in [03-Session-Establishment](03-Session-Establishment.md) and is not replaced by
ordinary Protobuf serialization.

Runtime validation remains mandatory for fixed-width fields, packet-type
and `oneof` consistency, protocol version, unknown fields, and frame limits.

## M6 Routing Envelope Rules

M6 uses the existing version-1 Protobuf envelope as the routable outer structure. Routing metadata includes packet type, PacketID, source/destination node identifiers and TTL. The relay can inspect this envelope for forwarding but must not gain access to endpoint plaintext or E2EE session keys.

PacketID is exactly 16 bytes. Incoming frames remain subject to the M4 64 KiB frame bound and strict structural validation, including protocol version, oneof consistency, fixed-width fields and unknown-field rejection.

The routing envelope and endpoint cryptographic session are deliberately separate: routing fields identify where a packet should be propagated, while M3 authentication and AEAD protect the application content.
