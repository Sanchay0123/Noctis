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

M3 is now closed. M4 is **NOT STARTED / GATED** and requires a separate
Project Overseer authorization.

## M4 --- Direct encrypted messaging integration

The milestone previously labelled “M4 — Encrypted messaging” is now
reframed as the next integration stage because the approved M3 implementation
already contains the cryptographic message-protection layer.

M4 must integrate the completed session/AEAD layer with the protocol and
application message path without changing the frozen cryptographic
construction unless separately reviewed.

**Gate:** Project Overseer authorization, followed by integration tests and
end-to-end evidence.

## M5 --- Direct networking

-   two-node connection
-   peer authentication
-   direct encrypted message

**Gate:** end-to-end two-node demo.

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
