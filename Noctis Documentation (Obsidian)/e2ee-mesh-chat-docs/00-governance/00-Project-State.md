# Project State

## Current phase

**M5 complete — post-milestone documentation synchronization / M6 architecture gate**

## Status

| Area | Status |
|---|---|
| Requirements | Documented |
| Threat model | Documented |
| Architecture | Approved baseline |
| Cryptographic architecture | M3 implemented / verified for session crypto scope |
| Protocol | Frozen baseline / M3 implementation verified |
| Networking | M5 direct networking hardening / multi-peer runtime integration complete and verified |
| Observability | Approved supporting design / M5 bounded telemetry lifecycle integrated as supporting runtime facility |
| Containerization | Required / foundation implemented |
| GitHub synchronization | Required / not yet verified in this gate |
| Testing strategy | Documented / M3 crypto evidence added |
| Implementation | M5 direct networking hardening / runtime integration implemented |
| Security verification | M5 networking lifecycle, resource-boundary and arbitration evidence verified |
| UI | Not started |
| Demo | Not started |
| Final academic material | Draft |

## Current gate

**M5 — Direct Networking Hardening / Runtime Integration: COMPLETE / APPROVED**

M0, M0.1, M1.1, M2.1/M2.1-B/M2.1-C, M3.1, M3.2 and M3.3 have been
completed within their approved scopes.

M3 covers:
- fresh ephemeral X25519 key agreement
- Ed25519-authenticated session establishment
- transcript-bound session identification
- HKDF-SHA-256 directional session-key derivation
- ChaCha20-Poly1305 application encryption
- deterministic sequence-derived nonces
- authenticated associated data and role-separated directions
- per-session replay-window enforcement
- concurrency protection for send and receive session state

The M3 implementation was reviewed against the frozen protocol construction
and its reported containerized validation evidence.

## Rule

No implementation component becomes **Approved** without implementation
evidence and relevant tests. A design being frozen or accepted does not mean
that its implementation exists.

## Next action

Complete the technical-lead documentation synchronization for M5, then open
M6 architecture separately. M6 implementation must not begin merely because M5 is complete.

## M3 Completion Gate

**M3 — Authenticated Session & Encrypted Message Layer: 🟢 COMPLETE**

### M3.1 — X25519
Implemented using Go's standard `crypto/ecdh.X25519()` API with fresh
ephemeral keypairs. The implementation exposes raw shared-secret material
only to the session KDF layer, rejects malformed/invalid peer public keys,
supports an explicit `Destroy()` lifecycle boundary, and documents that
destruction is a best-effort memory-lifecycle measure rather than guaranteed
RAM zeroization. The ephemeral-key lifecycle is sequential with respect to
destruction.

### M3.2 — Authenticated Session Establishment
Implemented using the frozen Ed25519/X25519 transcript and HKDF schedule.
The implementation verifies both signed handshake transcripts, derives the
full SHA-256 transcript as the session identifier, and produces independent
initiator-to-responder and responder-to-initiator keys. TOFU remains subject
to its documented first-contact MITM limitation.

### M3.3 — AEAD Message Layer
Implemented with ChaCha20-Poly1305. Each directional session key has an
independent sequence space beginning at zero. The nonce is
`0x00000000 || uint64_be(sequence_num)`. AAD binds the fixed protocol
context, session identifier, sequence number and internally derived
direction marker. A 64-message receive replay window is enforced.

Replay state is only committed after successful AEAD authentication, so an
unauthenticated packet cannot consume receiver replay-window state.
Independent send/receive mutexes protect concurrent session use. Sequence
number exhaustion returns an error and never wraps.

### Validation evidence
The M3 implementation reports containerized `gofmt`, `go vet`, unit tests,
race-detector tests and build validation as passing, including negative and
known-answer tests for the implemented cryptographic constructions.

M3 completion does **not** imply that direct networking, mesh routing,
observability, UI, or production-grade metadata/privacy protection is
complete.

**M4 has been authorized, implemented, independently reviewed, and approved. M5 has subsequently been authorized, implemented, independently reviewed, and approved. M6 remains gated pending architecture review.**


## M4 Completion Gate

**M4 — Direct Secure Messaging Integration: 🟢 COMPLETE / APPROVED**

M4 integrates the approved M3 cryptographic/session layer with a bounded TCP
transport and the frozen Protobuf packet format. The implementation provides a
direct two-node secure messaging path without changing M3 cryptographic
semantics.

### M4 Transport

The direct transport uses a 4-byte big-endian length prefix followed by a
Protobuf `MeshPacket`. The maximum frame size is 64 KiB and the length is
validated before body allocation. Reads use `io.ReadFull`; writes serialize
complete frames under a connection write mutex and correctly handle short
writes.

### M4 Packet boundary

Incoming packets are structurally validated before cryptographic processing.
Protocol version, packet type/oneof consistency, fixed-width identifiers,
payload-specific fields and unknown Protobuf fields are rejected according to
the version-1 policy.

### M4 Direct secure channel

`DirectChannel` binds one direct TCP connection to one established M3 session
for the M4 scope. INIT and RESP are exchanged through the existing M3 handshake
state machines. APP_DATA is accepted only after session establishment and is
bound to the expected local identity, remote identity and session identifier.

### M4 security evidence

The M4 test suite includes framing, partial/short writes, malformed input,
unknown-field rejection, session binding, end-to-end encrypted message
exchange and transport-path tampering tests. Containerized formatting, vet,
tests, race-detector tests and build validation passed.

M4 establishes direct secure messaging integration but does not establish
mesh routing, multi-hop forwarding, anonymity, metadata hiding, or endpoint
compromise resistance.

**M5 has been authorized, implemented, independently reviewed, and approved. M6 remains gated pending architecture review.**


## M5 Completion Gate

**M5 — Direct Networking Hardening / Runtime Integration: 🟢 COMPLETE / APPROVED**

M5 extends the approved M4 direct secure channel into a multi-peer runtime networking layer without changing the M3 cryptographic/session construction or M4 transport semantics.

### M5 runtime networking

The implementation introduces a `PeerManager` and `Peer` lifecycle above M4. It manages inbound and outbound authenticated direct connections, peer registry state, connection limits, handshake limits, dialing limits, timeouts, shutdown and application-triggered disconnect/reconnect behavior. Automatic PeerManager reconnect remains out of scope.

### M5 duplicate arbitration

Simultaneous direct connections between the same authenticated identities are resolved deterministically. The connection initiated by the lexicographically smaller 32-byte Ed25519 identity is retained. Registry replacement is performed under the manager mutex and the incumbent connection is closed after the registry transition, preventing stale cleanup from deleting the replacement.

### M5 resource and lifecycle controls

- Maximum active peers are enforced before registry insertion.
- Maximum pending handshakes are enforced with a non-blocking semaphore path.
- Concurrent outbound dials are bounded and slots are released on all outcomes.
- TCP handshake/dial deadlines bound incomplete connection attempts.
- Oversized initial frames are rejected before body allocation.
- Peer terminal states cannot be resurrected by stale goroutines.
- Manager shutdown is synchronized with networking goroutines.

### M5 telemetry boundary

Telemetry is an out-of-band supporting subsystem. Events enter a bounded queue using a non-blocking enqueue path. The telemetry worker is deliberately outside the PeerManager `WaitGroup`, so an indefinitely blocking external recorder cannot prevent manager shutdown. Telemetry failures are recovered at the recorder boundary and queue overflow drops events rather than blocking networking. No plaintext, private keys, session keys or passwords are exported.

### M5 validation evidence

The reported validation includes `go test -race ./internal/mesh`, repeated `go test -race ./internal/mesh -count=100`, explicit four-case duplicate-arbitration tests, and a deliberate Docker Dave↔Bob simultaneous-dial collision. In the runtime collision, Dave's identity was lexicographically smaller than Bob's, Dave's outbound connection survived, Bob's competing outbound connection was rejected as a duplicate, and both nodes retained exactly one active peer. Four Docker nodes (Alice, Bob, Carol, Dave) started successfully.

M5 does not establish multi-hop routing, route discovery, forwarding, TTL processing, routing tables, DHT, flooding algorithms, anonymity, metadata hiding, or production-grade network resilience. Those remain future scope.

**M5 is formally closed. M6 requires a separate architecture/design review before implementation.**
