# M5 Architecture & Implementation Review

## Review status

**🟢 COMPLETE / APPROVED — 12 September 2026**

M5 implements direct networking hardening and runtime peer management above the
approved M4 direct secure messaging boundary.

## Scope

- inbound TCP listener
- outbound dialing
- authenticated peer registration
- peer lifecycle management
- connection and handshake resource limits
- dial/handshake timeouts
- deterministic duplicate direct-connection arbitration
- stale-peer-safe registry replacement
- clean shutdown and race safety
- bounded out-of-band telemetry
- reproducible Docker runtime

## Frozen decisions implemented

### Duplicate arbitration

For two authenticated direct connections between the same pair, compare the
canonical 32-byte Ed25519 identities. The connection initiated by the
lexicographically smaller identity is retained.

### Resource limits

Pending handshakes and concurrent dials are bounded. Active peers are capped.
Exhausted pending-handshake capacity rejects immediately rather than blocking
the listener. Oversized initial frames are rejected before body allocation.

### Lifecycle ownership

Peer terminal states cannot be resurrected. Registry replacement is performed
under the manager mutex and the incumbent is closed after unlocking. Stale
cleanup cannot delete a newer replacement.

### Telemetry isolation

Telemetry uses a bounded non-blocking queue. The worker is intentionally outside
the PeerManager WaitGroup so an indefinitely blocking external recorder cannot
prevent manager shutdown. Telemetry failures are isolated and secret material
is excluded.

## Validation evidence

- `go test -race ./internal/mesh` — PASS
- `go test -race ./internal/mesh -count=100` — PASS
- `TestDuplicateArbitrationCases` — PASS
- Docker Compose configuration/build/startup — PASS
- Four-node runtime: Alice, Bob, Carol, Dave — PASS
- Deliberate simultaneous Dave↔Bob collision — PASS

Runtime collision evidence:

- Dave: `4765bd805913e158`
- Bob: `563f793f6ad353dd`
- Dave is lexicographically smaller.
- Dave's outbound connection survives.
- Bob's competing outbound connection is rejected as a duplicate.
- Dave and Bob each retain exactly one active peer.

## Boundary confirmation

M5 calls the approved M4 transport/channel APIs and introduces no replacement
cryptographic primitives, AEAD implementation, handshake transcript changes or
nonce construction. M6 routing and multi-hop forwarding are not implemented.

## Final gate

**M5: 🟢 APPROVED AND FROZEN**

**Next gate:** M6 architecture/design review.
