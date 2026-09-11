# Implementation Milestones

## M0 — Architecture

**Status: Revised review complete**

Architecture is **approved with changes**.

M1 and M2 are authorized. M3 remains gated on final cryptographic
protocol specification.

## M1 — Project Skeleton, Tooling & Containerization

**Status: AUTHORIZED**

- Go module
- repository structure
- Protobuf schema skeleton
- generated-code workflow
- Dockerfile
- Docker Compose topology skeleton
- CI workflow
- dependency management
- `.gitignore`
- configuration/secrets policy
- deterministic test command

### M1 gate

Clean environment must be able to:

1. build containers
2. generate/compile Protobuf
3. run tests
4. run the minimal node skeleton

No custom cryptographic handshake implementation in M1.

## M2 — Cryptographic Primitives & Identity

**Status: AUTHORIZED WITH SCOPE**

Implement and test:

- Ed25519 identity generation
- public/private key representation
- signing
- verification
- secure randomness

Do not invent session protocol semantics during M2.

## M3 — Secure Session Establishment

**Status: BLOCKED**

Before implementation, finalize:

- exact canonical transcript
- INIT/RESP wire messages
- exact signature inputs
- full 32-byte session ID
- exact HKDF construction
- handshake replay cache/state
- retransmission behavior
- timeout behavior
- key erasure lifecycle

### M3 gate

Cryptographic protocol review must pass.

## M4 — Encrypted Messaging

After M3 approval:

- ChaCha20-Poly1305
- canonical AAD
- nonce construction
- sequence numbers
- replay window
- key rotation

### Gate

Crypto + replay review.

## M5 — Direct Networking

- TCP
- 4-byte length-prefixed framing
- strict maximum frame size
- Protobuf parsing
- connection lifecycle
- backpressure

### Gate

Two-node integration tests.

## M6 — Mesh Routing

- bootstrap peers
- PacketID seen-cache
- TTL
- managed flooding
- resource controls
- four-node integration

### Gate

Mesh/routing review.

## M7 — UI & Observability

- CLI/UI
- Prometheus-compatible metrics
- structured security-safe events
- log aggregation
- Grafana dashboard
- telemetry failure isolation

### Gate

Observability security review.

## M8 — Security Hardening

- malformed input
- replay
- impersonation
- tampering
- malicious relay
- resource exhaustion
- log/telemetry injection
- secret/plaintext leakage
- dependency review

### Gate

Security review.

## M9 — Demonstration Environment

- complete Docker Compose topology
- four-node demonstration
- security demonstrations
- observability dashboard
- clean-environment reproduction

### Gate

Reproducibility test.

## M10 — Documentation & Release

- final report
- synchronized protocol documentation
- limitations
- test evidence
- GitHub repository state
- meaningful milestone commits

### Gate

Final technical-lead approval.
