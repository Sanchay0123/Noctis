# Mesh Architecture

## Responsibility

The mesh layer handles delivery, not application decryption.

## Core components

-   discovery
-   peer manager
-   connection manager
-   routing table
-   forwarding engine
-   duplicate cache
-   TTL enforcement

## Conceptual architecture

```mermaid
flowchart LR
    D[Discovery] --> P[Peer Manager]
    P --> R[Routing Table]
    R --> F[Forwarding Engine]
    F --> C[Connections]
```

## Relay invariant

``` text
Receive packet
   ↓
Validate mesh envelope
   ↓
Check duplicate / TTL
   ↓
Choose next hop
   ↓
Forward packet
```

No E2EE plaintext is required for this process.

## Routing attacks

The mesh layer must be tested against:

-   loops
-   duplicate packets
-   stale routes
-   malicious dropping
-   malformed routing fields
-   excessive route state
