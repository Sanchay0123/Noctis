> Predecessor review: [[00-governance/03-Architecture-Review-M5]]

# Architecture Review — M6 Mesh Routing

## Decision

**M6 — APPROVED / IMPLEMENTED / VERIFIED / CLOSED**

M6 extends the approved M5 direct multi-peer runtime with bounded managed flooding, TTL propagation control, PacketID duplicate suppression, bounded per-peer forwarding queues, endpoint E2EE session management independent of transport peers, and multi-hop M3 handshake forwarding/correlation.

## Scope

- bounded managed flooding
- TTL-based propagation control
- 16-byte PacketID duplicate suppression
- per-peer packet-count and byte-bounded forwarding queues
- multi-hop handshake forwarding and correlation
- relay-blind endpoint E2EE

## Explicit non-goals

- DHT
- route discovery or shortest-path routing
- dynamic routing protocols
- anonymity or traffic-analysis resistance
- guaranteed delivery
- post-compromise security
- complete metadata hiding

## Frozen parameters

| Parameter | Value |
|---|---:|
| PacketID | 16 bytes |
| Duplicate cache | 10,000 entries |
| Cache lifetime | 2 minutes |
| Maximum accepted TTL | 32 |
| Default initial TTL | 16 |
| Per-peer forward queue | 1,000 packets |
| Per-peer queue bytes | 2 MiB |
| Maximum active peers | 50 |

## Routing model

M6 uses bounded managed flooding rather than route discovery. A structurally valid packet is duplicate-checked, TTL-checked, delivered locally when appropriate, and otherwise forwarded to eligible active peers other than the incoming peer. The router does not perform network I/O itself; forwarding is handed to per-peer non-blocking queues.

## Multi-hop handshake

Endpoint session establishment remains the M3 cryptographic operation. A relay forwards INIT/RESP packets but does not terminate the endpoint session. Initiator pending state is scoped by remote identity and `SHA256(T_INIT)` so simultaneous handshakes with the same remote can be correlated safely.

A response is accepted only when exactly one pending candidate reconstructs the exact frozen M3 response transcript and verifies the responder signature. Zero matches and ambiguous matches fail closed.

## Security boundary

`source_node` and `dest_node` are routing metadata, not cryptographic authentication. Endpoint identity is authenticated only by the M3 Ed25519 transcript signatures and the resulting session establishment. Relays must not decrypt APP_DATA, derive endpoint session keys, or terminate endpoint E2EE sessions.

## Resource bounds

The duplicate cache and forwarding queues are explicitly bounded. Queue insertion is non-blocking and fails when either packet-count or byte limits are reached. Cache and queue accounting are concurrency-safe. These controls bound memory growth but do not make the prototype immune to malicious traffic or bandwidth exhaustion.

## Acceptance evidence

The M6 final acceptance record reports passing handshake-correlation, TTL, malformed/oversized packet, unknown-field, cache, queue-accounting, multi-hop routing, cyclic-routing, relay-confidentiality, peer-loss/session-separation, sustained-flood, race-detector, build, and Docker multi-hop validation.

See [[07-testing/08-M6-Acceptance-Evidence]] for the detailed evidence record.
