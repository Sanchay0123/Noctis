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

**Gate:** security test review.

## M8 --- Observability and UI

### Observability

-   Prometheus-compatible metrics
-   structured security-safe events
-   log aggregation
-   Grafana dashboard
-   telemetry failure isolation

### UI

Only after protocol stability.

**Gate:** observability safety review + UI smoke tests.

## M9 --- Demonstration environment

-   fully containerized node topology
-   observability stack
-   reproducible setup
-   scripted security demonstrations
-   documented GitHub workflow

**Gate:** clean-environment reproduction.

## M10 --- Documentation

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
