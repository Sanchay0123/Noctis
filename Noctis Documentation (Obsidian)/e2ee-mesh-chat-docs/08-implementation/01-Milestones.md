# Implementation Milestones

## M0 --- Architecture

Deliver:

-   requirements
-   threat model
-   architecture
-   component model
-   data flows
-   protocol proposal
-   cryptographic design
-   technology decision
-   repository plan
-   testing strategy

**Gate:** Technical lead approval.

## M1 --- Skeleton

-   repository
-   build system
-   configuration
-   logging
-   test framework
-   startup

**Gate:** build and tests pass.

## M2 --- Identity

-   Ed25519 generation
-   storage
-   loading
-   signatures
-   verification

**Gate:** crypto review + tests.

## M3 --- Authenticated session & encrypted message layer

-   X25519
-   Ed25519 authentication
-   canonical transcript binding
-   HKDF-SHA-256
-   session lifecycle
-   ChaCha20-Poly1305
-   nonce management
-   AAD
-   directional sequence spaces
-   bounded replay protection

**Gate:** cryptographic/session implementation review + negative tests +
known-answer vectors + race-detector validation.

**Status:** 🟢 COMPLETE / APPROVED

### M3 Completion Record

M3.1, M3.2 and M3.3 were implemented and independently reviewed against
the frozen protocol construction.

M3.1 established fresh X25519 ephemeral key agreement.

M3.2 established mutually authenticated sessions using Ed25519 signatures,
the canonical handshake transcript, SHA-256 session IDs, and the frozen
HKDF-SHA-256 directional key schedule.

M3.3 established the application AEAD layer using ChaCha20-Poly1305,
sequence-derived nonces, authenticated associated data, role-separated
directions, a 64-message replay window, and concurrency-safe session state.

The reported containerized validation suite passed formatting, vetting, tests,
race detection, and build checks.

M3 is now closed. M4 has subsequently been implemented, independently reviewed,
and approved. M5 has now been implemented, independently reviewed, and approved
by the Project Overseer.

## M4 --- Direct secure messaging integration

M4 integrates the completed M3 session/AEAD layer with the frozen Protobuf
protocol and bounded TCP transport. The scope includes direct-channel
handshake orchestration, packet validation, session/identity binding,
application plaintext delivery after successful AEAD authentication, and
transport-path tampering evidence.

For M4, one established cryptographic session is associated with one direct TCP
channel. This is a milestone scope constraint, not a protocol limitation.

**Status:** 🟢 COMPLETE / APPROVED

### M4 Completion Record

M4 implemented `Connection`, packet validation and `DirectChannel` integration.
The transport uses 4-byte big-endian framing with a 64 KiB bound, bounded reads,
serialized complete writes including short-write handling, and explicit
unknown-field rejection. Direct APP_DATA is bound to the active session and
endpoint identities. End-to-end tests demonstrate bidirectional encrypted
messaging and rejection of transport-path ciphertext tampering.

Containerized formatting, vetting, tests, race detection and build validation
passed. M4 was independently reviewed and approved by the Project Overseer.

## M5 --- Direct networking hardening / runtime integration

M5 extends the M4 direct secure channel into the project's runtime networking
environment. Scope includes connection management, peer relationships,
inbound/outbound dialing, authenticated peer registration, deterministic
duplicate-connection arbitration, resource limits, lifecycle/shutdown safety,
bounded telemetry isolation, and reproducible containerized multi-node runtime
evidence.

**Status:** 🟢 COMPLETE / APPROVED

**Gate:** Project Overseer final evidence review — PASSED.

### M5 Completion Record

M5 implemented the `PeerManager` / `Peer` runtime layer above the approved M4
transport boundary. It added inbound listening, outbound dialing, peer lifecycle
states, identity-bound peer registration, handshake and dial limits, active-peer
limits, timeout enforcement, stale-peer-safe registry replacement, and clean
shutdown.

Duplicate direct connections are resolved by the frozen rule that the connection
initiated by the lexicographically smaller Ed25519 identity wins. The final
Docker evidence deliberately created simultaneous Dave↔Bob outbound dials; Dave
(`4765bd805913e158`) was smaller than Bob (`563f793f6ad353dd`), Dave's outbound
connection survived, Bob's competing connection was rejected as a duplicate, and
both nodes retained one active peer.

Telemetry uses a bounded non-blocking queue and a worker deliberately outside
the PeerManager `WaitGroup`, isolating blocking/failing external recorders from
network shutdown and delivery.

Validation passed `go test -race ./internal/mesh` and the repeated
`go test -race ./internal/mesh -count=100` run, plus the explicit four-case
`TestDuplicateArbitrationCases`. Docker Compose successfully built and started
the four-node Alice/Bob/Carol/Dave environment.

M5 does not implement M6 routing or multi-hop forwarding.

## M6 --- Mesh routing

**Status:** 🟢 COMPLETE / APPROVED / VERIFIED

M6 extends M5 into a bounded multi-hop mesh. The implementation uses managed flooding with TTL termination, PacketID duplicate suppression, bounded per-peer forwarding queues, and endpoint E2EE sessions independent from transport peers.

Frozen parameters:
- PacketID: 16 bytes
- duplicate cache: 10,000 entries / 2 minutes
- default TTL: 16
- maximum accepted TTL: 32
- per-peer forwarding queue: 1,000 packets / 2 MiB
- maximum active peers: 50

Multi-hop INIT/RESP messages are correlated using remote identity plus `SHA256(T_INIT)` and exact M3 response reconstruction. Relays forward ciphertext but do not terminate endpoint E2EE sessions.

Validation evidence includes one-, two- and three-hop routing, cyclic-loop termination, handshake correlation, malformed/oversized input rejection, cache and queue bounds, relay-blind E2EE, peer-loss/session separation, sustained flood bounds, race validation and Docker multi-hop operation.

See [[00-governance/04-Architecture-Review-M6]] and [[07-testing/08-M6-Acceptance-Evidence]].

## M7 --- Security hardening

-   malformed input
-   replay
-   impersonation
-   packet tampering
-   malicious relay
-   basic DoS
-   telemetry leakage
-   telemetry/log injection
-   telemetry resource bounds
-   panic/failure safety
-   concurrency/resource lifecycle

**Status:** 🟢 COMPLETE / APPROVED / VERIFIED

**Gate:** Project Overseer final security evidence review — PASSED.

### M7 Completion Record

M7 hardened the approved M6 implementation without changing the frozen M3,
M4 or M6 architecture.

The concurrent replay race was closed by serializing replay validation, AEAD
authentication and replay-state commitment.

Router input is validated at the routing boundary before cache insertion or
forwarding. Final adversarial coverage includes malformed Protobuf, unknown
protocol version, unknown fields, type/oneof mismatch, invalid fixed-width
fields and oversized input.

Telemetry was reviewed for leakage, injection and metric cardinality. Peer
identity was removed from the `TelemetryRecorder` interface. The current
repository contains telemetry hooks/no-op behavior but no active Prometheus
metric exporter or metric-vector implementation. Future M8 metrics must not
introduce attacker-controlled peer identities as labels.

Handshake deadlines and pending-slot release were verified under timeout.
Authentication regressions, failure paths, concurrency behavior and resource
bounds were reviewed.

Final validation passed:

```text
go vet ./...
go test ./...
go test -race ./...
go build ./...
```

See [[00-governance/05-Architecture-Review-M7]] and
[[07-testing/09-M7-Acceptance-Evidence]].

M7 does not establish anonymity, complete metadata hiding, global Sybil/flood
resistance, endpoint compromise resistance or guaranteed delivery.

## M8 --- Observability and UI

**Status:** 🟢 COMPLETE / APPROVED / VERIFIED for the implemented M8.1–M8.3 scope

### M8.1 — ApplicationService boundary

Implemented the application-facing service and event boundary. A critical initial routing-decryption violation was remediated before approval: routing now forwards opaque APP_DATA and `internal/app` performs endpoint decryption after session lookup.

### M8.2 — Fyne GUI foundation

Implemented the Fyne GUI foundation with manual peer entry, synchronized UI state and lifecycle/event handling while keeping GUI code isolated from crypto/session/mesh/routing internals.

### M8.3 — Functional GUI integration

Implemented manual dialing followed by explicit conversation/session establishment, secure-session status events, public identity copy, structured in-memory messages, conversation handling, application-layer multi-hop integration and deterministic GUI event-consumer shutdown.

The Alice→Bob→Carol integration test uses real loopback TCP, PeerManager, routing/forwarding, session establishment and ChaCha20-Poly1305 encryption/decryption. It is classified as **APPLICATION-LAYER INTEGRATION** because the three ApplicationServices execute within one test process.

Validation passed `go test ./...`, `go test -race -p 1 ./...`, `go vet ./...` and `go build ./...`; GUI build verification also passed. Docker regression remains unverified because of external proxy/DNS limitations.

Concrete Prometheus/Grafana exporter implementation remains deferred; M7 telemetry hooks and security/cardinality constraints remain authoritative.

**Gate:** Project Overseer final evidence review — PASSED.

### M8.3 Remediation 4 — Transport Admission Modernization

A deterministic live GUI failure was traced to the M5 `ExpectInbound` transport
pre-authorization mechanism. The first initiator could not connect because the
remote listener had not pre-authorized that identity; the failed attempt left
local authorization state that allowed the reverse direction.

The obsolete `ExpectInbound`, `expectedInbound` and `isExpectedInbound`
mechanism was removed. Inbound transport now proceeds through bounded
handshake admission and normal M3 authentication before peer registration.
Expected-PeerID verification remains enforced at the application boundary and
duplicate arbitration remains unchanged. No protocol/schema or cryptographic
construction changed.

Regression coverage includes `TestListener_AcceptsUnknownPeer` and
`TestExpectedPeerIDStillEnforced`. Full tests, race tests, vet and builds pass.

Final live GUI acceptance passed in both initiation directions, with secure
session establishment, bidirectional encrypted messaging and idempotent
duplicate Add Peer behavior.

**M8.3: 🟢 COMPLETE / LIVE ACCEPTED / APPROVED**

## M9 --- Demonstration & Security Validation

**Status:** 🟢 COMPLETE / CLOSED — M9.1/M9.2/M9.3 CLOSED; M9.4 AUDITED / CLOSED

M9 validates the implemented system through security demonstrations, relay-blind
E2EE evidence and live GUI verification. It does not change the approved M3–M8
cryptographic, routing or presentation architecture.

### M9.1 — E2EE / Relay Blindness

**🟢 COMPLETE / CLOSED**

Alice → Bob → Carol demonstrates application ciphertext traversing a relay while
the endpoint application session remains at the communicating endpoints in the
demonstrated implementation.

See [[07-testing/17-M9-Relay-Blindness.log]] and [[07-testing/19-M9-Evidence-Manifest]].

### M9.2 — Security Attack Demonstrations

**🟢 COMPLETE / CLOSED**

Eight adversarial demonstrations passed: Wrong PeerID, ciphertext tampering,
AEAD replay, invalid signature, malformed packet, TTL enforcement, PacketID
duplicate suppression and the production pending-handshake resource bound.
The authoritative pending-handshake limit is 10.

See [[07-testing/16-M9-Attack-Evidence]].

### M9.3 — GUI Demonstration

**🟢 COMPLETE / CLOSED**

The two-node GUI workflow was manually verified by the Project Overseer on the
desktop. Automated GUI interaction was unavailable in the headless implementation
environment.

See [[07-testing/18-M9-GUI-Evidence]].

### M9.4 — Final Evidence Assembly & Documentation Synchronization

**🟢 COMPLETE / CLOSED**

M9 evidence records and documentation synchronization were assembled and audited.
The Project Overseer audit passed the package without requiring production-code
changes. M9 is formally closed.

See [[00-governance/07-Architecture-Review-M9]],
[[07-testing/19-M9-Evidence-Manifest]], [[07-testing/20-M9-Final-Report]] and
[[07-testing/21-M9.4-Overseer-Audit]].

## M10 --- Final Technical & Academic Documentation

Final technical and academic documentation.

See [[08-implementation/05-Definition-of-Done]].

## M2 Completion Record

**Status: 🟢 COMPLETE**

M2.1 Ed25519 Identity Core was implemented, followed by M2.1-B cryptographic
audit and M2.1-C final key-ownership correction.

M2 acceptance evidence includes containerized formatting, vetting, unit tests,
race-detector tests, and builds.

**Next gate:** post-M2 documentation synchronization and review.

M3 must not begin until that documentation gate is completed and separately
authorized.


## M10 — Final Technical & Academic Documentation

**Status: 🟢 COMPLETE / CLOSED — FINAL DOCUMENTATION FREEZE**

M10 produced the final technical baseline, security analysis, controlled experimental results, academic report material, presentation structure, viva preparation and evidence index. Final system validation subsequently passed all executable checks available in the validation environment.

### M10.0 — Reconnaissance

**🟢 PASS / CLOSED.** The implementation was inventoried without modifying production behavior.

### M10.1 — Technical Documentation

**🟢 PASS / CLOSED.** The technical baseline documents the implemented architecture, components, cryptography, protocol, mesh forwarding, application/GUI boundary and reproducibility constraints.

### M10.2 — Security Analysis & Experimental Results

**🟢 PASS / CLOSED.** Approved M9 security evidence was mapped to security properties and controls, and a controlled local experimental campaign produced the documented cryptographic, latency, throughput, message-size, resource and handshake-repeatability measurements. AEAD decryption performance, CPU utilization percentage and concurrent handshake scaling remain explicitly NOT MEASURED.

### M10.3 — Final Academic Deliverables

**🟢 PASS / COMPLETE / CLOSED.** The final academic report, presentation outline, viva Q&A and evidence index are complete.

### Final System Validation

**🟢 PASS.** The final validation campaign passed the executable cryptographic, protocol/framing, direct E2EE, three-node mesh, routing, security-regression and resource-bound checks. GUI final re-execution was NOT MEASURED in the headless environment; previously accepted manual GUI evidence remains part of the baseline.

### Documentation Freeze

**🟢 FROZEN.** No further implementation or protocol changes are authorized unless a genuine defect or factual inconsistency is discovered during submission review.
