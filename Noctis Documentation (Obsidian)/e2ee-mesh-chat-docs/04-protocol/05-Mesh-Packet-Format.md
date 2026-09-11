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

Candidates should be evaluated for:

-   deterministic encoding
-   schema evolution
-   compactness
-   validation
-   interoperability
-   ease of debugging

No serialization format is approved by this document alone.
