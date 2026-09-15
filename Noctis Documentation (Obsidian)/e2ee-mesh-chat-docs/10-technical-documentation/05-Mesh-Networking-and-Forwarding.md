# Mesh Networking and Forwarding — Implemented System

## Overview

The mesh runtime combines authenticated direct peer connections with bounded
managed flooding. It does **not** implement DHT-based routing, route discovery,
or route selection.

```mermaid
flowchart LR
    A[Alice] --> B[Bob]
    B --> C[Carol]
    C --> D[Dave]
```

A packet can be forwarded through active peers until it reaches its
(destination) node or its TTL prevents further forwarding.

## PeerManager

`internal/mesh` owns direct peer lifecycle.

Important production bounds include:

```text
MaxActivePeers       = 50
MaxPendingHandshakes = 10
Handshake timeout    = 5 seconds
```

The pending-handshake bound is enforced with a semaphore. The M9.2 production-
bound demonstration attempted 15 incomplete TCP handshakes and observed that
simultaneous admission remained at or below 10.

## Peer states

The direct peer runtime uses states including:

```text
CONNECTING
HANDSHAKING
ESTABLISHED
CLOSING
FAILED
```

Peer identity is established from the authenticated Ed25519 identity rather
than trusted solely from a network address.

## Duplicate connection arbitration

When two nodes create simultaneous connections to the same authenticated
identity, deterministic identity ordering resolves the collision. The
lexicographically smaller identity wins the connection it initiated.

The rule is evaluated after authenticated identity is available; it is not an
address-based trust decision.

## Routing lifecycle

A received packet follows this conceptual sequence:

```text
receive
  ↓
validate envelope
  ↓
check duplicate cache
  ↓
record (source_node, packet_id)
  ↓
local destination?
  ├── yes → deliver/process
  └── no
       ↓
     TTL <= 1?
       ├── yes → drop
       └── no → decrement + bounded forward
```

The exact cache insertion/TTL semantics are important: a structurally valid
packet is eligible for duplicate suppression before forwarding decisions, so
repeated packets cannot continuously amplify through the mesh.

## Forwarding queues

Forwarding uses bounded per-peer queues. The implementation bounds both packet
count and queue memory and avoids performing network I/O while holding routing
state locks.

When a peer closes, its queued forwarding work is discarded rather than being
allowed to grow without bound.

## Relay confidentiality

The router handles the outer packet required for forwarding. For `APP_DATA`,
it passes the session identifier, sequence number and ciphertext toward the
application boundary without calling `DecryptMessage` itself.

This separation was demonstrated in M9.1 using Alice → Bob → Carol.

## Failure handling

The implementation contains explicit limits and failure paths for malformed
input, oversized frames, failed authentication, handshake timeout, duplicate
packets, TTL exhaustion, and peer lifecycle failure.

The system does not claim guaranteed delivery or immunity to malicious packet
dropping.

## Related records

- [[06-networking/01-Mesh-Architecture]]
- [[06-networking/03-Routing]]
- [[06-networking/04-Forwarding]]
- [[06-networking/05-TTL-and-Loop-Prevention]]
- [[06-networking/06-Connection-Management]]
- [[06-networking/08-M6-Resource-Bounds]]
- [[07-testing/16-M9-Attack-Evidence]]
