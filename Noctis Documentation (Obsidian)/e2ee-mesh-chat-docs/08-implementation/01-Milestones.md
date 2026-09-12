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
and approved. M5 is **NOT STARTED / GATED** and requires separate Project
Overseer authorization.

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

M5 extends the M4 direct secure channel into the project's intended runtime
networking environment and validates real two-node operation beyond the M4
transport/integration harness. Scope includes connection management, peer
relationships, runtime configuration and clean two-node demonstration.

**Status:** 🔒 NOT STARTED / GATED

**Gate:** Project Overseer authorization + reproducible two-node runtime demo.

## M6 --- Mesh routing

-   discovery
-   routing
-   multi-hop
-   TTL
-   duplicate detection
-   failure handling

**Gate:** four-node demonstration.

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
