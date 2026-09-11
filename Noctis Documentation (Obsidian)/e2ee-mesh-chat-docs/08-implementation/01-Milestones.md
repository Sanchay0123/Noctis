# Implementation Milestones

## M0 --- Architecture

Deliver:

- requirements
- threat model
- architecture
- component model
- data flows
- protocol proposal
- cryptographic design
- technology decision
- repository plan
- testing strategy

### M0.1 --- Repository & Architecture Initialization/Audit

- Go module and repository structure
- package boundary skeleton
- `.gitignore` and secret policy
- Docker/Compose development skeleton
- CLI entrypoint stub
- deterministic build/test command definition
- initial dependency policy
- architecture consistency audit

**Status:** Completed as a foundation task; implementation validation remains subject to technical-lead review.

**Gate:** Technical lead review of actual repository artifacts and validation evidence.

### M0.2 --- Architecture Consistency / Specification Freeze

- verify repository structure against architecture
- freeze protocol terminology and packet schema
- verify security boundaries
- verify containerization constraints
- verify testing and security acceptance criteria

**Gate:** Technical lead approval before protocol implementation.

## M1 --- Skeleton

- repository
- build system
- configuration
- logging
- test framework
- startup
- Protobuf schema/tooling
- generated-code workflow
- CI baseline
- containerized deterministic validation

**Constraint:** No custom cryptographic handshake, identity implementation, encrypted messaging, or mesh-routing implementation.

**Gate:** build and tests pass; generated-code workflow is reproducible; technical lead approval.

## M2 --- Identity

- Ed25519 generation
- storage
- loading
- signatures
- verification

**Gate:** crypto review + tests.

## M3 --- Session establishment

- X25519
- Ed25519 authentication
- HKDF
- session lifecycle

**Gate:** MITM test passes.

## M4 --- Encrypted messaging

- ChaCha20-Poly1305
- nonce management
- AAD
- message serialization
- replay model

**Gate:** negative crypto tests pass.

## M5 --- Direct networking

- two-node connection
- peer authentication
- direct encrypted message

**Gate:** end-to-end two-node demo.

## M6 --- Mesh routing

- discovery
- routing
- multi-hop
- TTL
- duplicate detection
- failure handling

**Gate:** four-node demonstration.

## M7 --- Security hardening

- malformed input
- replay
- impersonation
- packet tampering
- malicious relay
- basic DoS
- telemetry leakage
- telemetry/log injection
- telemetry resource bounds

**Gate:** security test review.

## M8 --- Observability and UI

### Observability

- Prometheus-compatible metrics
- structured security-safe events
- log aggregation
- Grafana dashboard
- telemetry failure isolation

### UI

Only after protocol stability.

**Gate:** observability safety review + UI smoke tests.

## M9 --- Demonstration environment

- fully containerized node topology
- observability stack
- reproducible setup
- scripted security demonstrations
- documented GitHub workflow

**Gate:** clean-environment reproduction.

## M10 --- Documentation

Final technical and academic documentation.

See [[08-implementation/05-Definition-of-Done]].
