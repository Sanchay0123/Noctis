# Project State

## Current Authoritative Status

> **M8.3 — Functional GUI Integration: COMPLETE / VERIFIED / APPROVED**
>
> M0–M7 remain the approved architectural and implementation baseline. M8 adds the ApplicationService boundary and a Fyne presentation layer while preserving endpoint E2EE and relay-blind routing.
>
> M8.1, M8.2 and M8.3 acceptance records are the authoritative sources for the M8 implementation status. Historical review documents may describe earlier states; those descriptions are not the current project status.


## Current phase

**M8.3 complete — live acceptance passed / post-milestone documentation synchronization**

## Status

| Area | Status |
|---|---|
| Requirements | Documented |
| Threat model | Documented |
| Architecture | Approved baseline |
| Cryptographic architecture | M3 implemented / verified for session crypto scope |
| Protocol | Frozen baseline / M3 implementation verified |
| Networking | M6 bounded mesh routing / multi-hop runtime complete and verified |
| Observability | Approved supporting design / runtime telemetry hooks; concrete Prometheus exporter remains deferred |
| Containerization | Required / foundation implemented |
| GitHub synchronization | Required / not yet verified in this gate |
| Testing strategy | Documented / M3 crypto evidence added |
| Implementation | M8.3 functional GUI integration implemented and approved above M7 security baseline |
| Security verification | M8.3 GUI/application boundary and routing plaintext-blindness verified; M6/M7 security evidence retained |
| UI | M8.1/M8.2/M8.3 implemented / functional integration approved; live two-node GUI acceptance passed |
| Demo | Not started / M9 pending separate authorization |
| Final academic material | Draft |

## Current gate

**M8.3 — Functional GUI Integration: COMPLETE / APPROVED**

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

M8 documentation synchronization is complete. M9 — Demonstration Environment — is the next milestone and requires a separate architecture/design and implementation authorization.

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


## M6 Completion Gate

**M6 — Mesh Routing / Multi-Hop E2EE: 🟢 COMPLETE / APPROVED**

M6 implemented bounded managed flooding above the M5 peer runtime. It added TTL propagation control, 16-byte PacketID duplicate suppression, bounded per-peer forwarding queues, multi-hop INIT/RESP forwarding with exact cryptographic correlation, and endpoint E2EE session management independent of transport peers.

Frozen bounds are: 10,000 cache entries with 2-minute lifetime; maximum accepted TTL 32; default initial TTL 16; 1,000 queued packets or 2 MiB per peer; and 50 maximum active peers.

The final acceptance record reports passing handshake-correlation, TTL, packet-validation, cache, queue, one-/two-/three-hop, cyclic-routing, relay-confidentiality, peer-loss/session-separation, sustained-flood, race, build and Docker multi-hop evidence.

M6 remains deliberately limited to bounded managed flooding. DHT, route discovery, route selection, anonymity, complete metadata hiding and guaranteed delivery remain out of scope.

See [[00-governance/04-Architecture-Review-M6]] and [[07-testing/08-M6-Acceptance-Evidence]].

## Governance Reviews

- [[00-governance/03-Architecture-Review-M5]] — M5 architecture review
- [[00-governance/04-Architecture-Review-M6]] — M6 architecture review


## M7 --- Security hardening

**Status:** 🟢 COMPLETE / APPROVED / VERIFIED

M7 hardened the approved M6 runtime without changing the frozen M3, M4 or M6
architecture.

### M7 Completion Record

The concurrent replay race was closed by serializing replay-window validation,
AEAD authentication and replay-state commitment under the receive mutex.
`TestAEADConcurrentReplay` verifies exactly one successful decryption and nine
replay rejections for ten concurrent attempts.

`Router.OnMessage` now rejects malformed and oversized input at the routing
boundary before normal cache/routing processing. The final boundary evidence
covers malformed Protobuf, unknown protocol version, unknown fields,
type/oneof mismatch, invalid fixed-width fields and oversized input.

Telemetry was hardened by removing arbitrary remote peer identities from the
`TelemetryRecorder` interface. The current repository contains telemetry
hooks/no-op behavior but no active Prometheus exporter or metric-vector
implementation. Future M8 metrics must use fixed, bounded label vocabularies
and must not use attacker-controlled peer identities as arbitrary labels.

Handshake deadlines and deferred pending-slot cleanup were verified, including
timeout and subsequent slot reuse. Authentication regressions and failure
paths were reviewed.

Final validation passed `go vet ./...`, `go test ./...`,
`go test -race ./...` and `go build ./...`.

See [[00-governance/05-Architecture-Review-M7]] and
[[07-testing/09-M7-Acceptance-Evidence]].

M7 does not provide anonymity, complete metadata hiding, global Sybil/flood
resistance, endpoint compromise resistance or guaranteed delivery.



## M8 Completion Gate

**M8.3 — Functional GUI Integration: 🟢 COMPLETE / VERIFIED / APPROVED**

M8.1 established the ApplicationService boundary and moved endpoint APP_DATA decryption out of routing. M8.2 introduced the Fyne GUI foundation with synchronized UI state. M8.3 completed the functional chat workflow, including manual peer dialing, explicit conversation startup, secure-session status events, public identity copy, structured in-memory message history, message send/receive integration and deterministic GUI event-consumer shutdown.

The M8.3 application integration test uses real loopback TCP, PeerManager, mesh routing/forwarding, session establishment and ChaCha20-Poly1305 encryption/decryption, delivering the final plaintext through `SubscribeEvents()`. M8.3 Remediation 4 removed the obsolete `ExpectInbound` transport pre-authorization dependency; inbound peers now enter the bounded authenticated handshake directly.

The GUI remains a thin presentation client and has no direct imports of the crypto, session, mesh or routing packages. `internal/routing` remains plaintext-blind.

Repository validation passed `go test ./...`, `go test -race -p 1 ./...`, `go vet ./...` and `go build ./...`. GUI visual runtime and Docker regression remain unverified because the review environment could not obtain required native graphics/dependency packages due external DNS/proxy restrictions.

### M8.1

**COMPLETE / APPROVED** — ApplicationService boundary and endpoint decryption placement.

### M8.2

**COMPLETE / APPROVED** — Fyne GUI foundation and synchronized UI state.

### M8.3

**COMPLETE / APPROVED** — Functional GUI integration and application-layer multi-hop messaging evidence.

M8 is formally closed. M9 requires separate authorization.


## M8.3 Live Acceptance Record

**Status: 🟢 PASSED / M8.3 CLOSED**

Final live GUI acceptance was performed after Remediation 4. Two clean GUI
instances were run with distinct loopback listeners (`127.0.0.1:8000` and
`127.0.0.1:8001`). The first connection succeeded when Alice initiated and
when Bob initiated in a clean run. Secure-session establishment completed in
both directions. Bidirectional encrypted messaging succeeded. Repeated Add
Peer on an already established connection remained idempotent and did not
disrupt the active conversation state.

Remediation 4 removed `ExpectInbound`, `expectedInbound` and
`isExpectedInbound` from the transport admission path. The authenticated
handshake, expected-PeerID verification, pending-handshake bound, timeout,
frame-size validation and duplicate arbitration remain enforced.

M8.3 is therefore formally closed. M9 remains subject to separate
authorization.
