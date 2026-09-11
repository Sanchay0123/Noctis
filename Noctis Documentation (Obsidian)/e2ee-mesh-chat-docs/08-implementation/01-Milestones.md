# Implementation Milestones

## M0 — Architecture

**Status: Under review**

- architecture baseline
- threat model
- protocol design
- crypto design
- observability design
- containerization strategy
- GitHub/CI strategy

**Gate:** architecture approval.

## M1 — Project Skeleton & Tooling

- Go module
- repository structure
- Protobuf schema and generated code
- Dockerfile
- Compose demonstration topology
- deterministic test command
- CI workflow
- lint/static analysis where justified
- `.gitignore`
- configuration/secrets policy

**Gate:** clean containerized test run and CI pass.

## M2 — Cryptographic Primitives & Identity

- Ed25519 identity
- secure key generation
- key serialization/storage policy
- signature tests

**Gate:** crypto review.

## M3 — Secure Session Establishment

- canonical signed handshake
- X25519 ephemeral agreement
- transcript-bound KDF
- directional session keys
- session state machine
- handshake replay handling
- key erasure/rotation

**Gate:** cryptographic protocol review.

## M4 — Encrypted Messaging

- ChaCha20-Poly1305
- canonical AAD
- sequence numbers
- replay window
- nonce uniqueness enforcement

**Gate:** crypto + replay test review.

## M5 — Direct Networking

- TCP
- message framing
- Protobuf parsing
- connection lifecycle
- backpressure/error handling

**Gate:** two-node integration tests.

## M6 — Mesh Routing

- peer management
- PacketID cache
- TTL
- managed flooding
- resource limits
- four-node integration tests

**Gate:** mesh test review.

## M7 — Observability & Security Telemetry

- Prometheus-compatible metrics
- structured security-safe events
- log aggregation candidate (Loki preferred)
- Grafana dashboard
- telemetry failure isolation
- telemetry leakage tests

**Gate:** observability/security review.

## M8 — UI & Demonstration Environment

- CLI or web UI
- reproducible four-node topology
- security demonstrations
- dashboard demonstrations

**Gate:** clean-environment demo.

## M9 — Security Hardening

- malformed input
- replay
- impersonation
- tampering
- malicious relay
- resource exhaustion
- telemetry leakage/injection
- dependency review

**Gate:** security review.

## M10 — Documentation & Release

- final report
- architecture/protocol documentation
- limitations
- test evidence
- reproducibility instructions
- GitHub release state

**Gate:** final technical-lead approval.
