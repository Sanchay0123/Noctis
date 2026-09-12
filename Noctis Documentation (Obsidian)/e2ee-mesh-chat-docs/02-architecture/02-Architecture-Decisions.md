# Architecture Decisions

## ADR-001 --- Layer separation

**Status:** Proposed

### Decision

Separate transport, mesh, protocol, cryptographic/session and
application layers.

### Rationale

Prevents routing code from becoming responsible for cryptographic state
and prevents cryptographic code from depending on routing internals.

### Consequence

More interfaces must be defined, but security analysis becomes clearer.

------------------------------------------------------------------------

## ADR-002 --- Ed25519 for long-term identity

**Status:** Implemented / verified in M2

### Decision

Use Ed25519 for persistent node identity and signatures.

### Rationale

The project requires a well-established asymmetric signature primitive
and explicitly prefers Ed25519.

### Consequence

The implementation must distinguish identity keys from X25519
key-agreement keys.

------------------------------------------------------------------------

## ADR-003 --- X25519 for key agreement

**Status:** Accepted / implementation gated to M3

### Decision

Use X25519 for establishing shared key material.

### Rationale

Well-established elliptic-curve Diffie-Hellman construction and
explicitly requested by the project.

### Consequence

The X25519 exchange must be authenticated by the Ed25519 identity.

------------------------------------------------------------------------

## ADR-004 --- ChaCha20-Poly1305

**Status:** Proposed

### Decision

Prefer ChaCha20-Poly1305 as the single application AEAD.

### Rationale

Avoid unnecessary implementation complexity from supporting multiple
AEAD algorithms.

### Consequence

Nonce lifecycle becomes a central protocol invariant.

------------------------------------------------------------------------

## ADR-005 --- HKDF-SHA-256

**Status:** Proposed

### Decision

Use HKDF-SHA-256 for deriving session/application keys.

### Rationale

Provides a standard KDF with explicit context separation.

------------------------------------------------------------------------

## ADR-006 --- No custom cryptography

**Status:** Accepted

### Decision

Cryptographic primitives must come from mature libraries.

### Rationale

Custom cryptographic algorithms and constructions create
disproportionate security risk.

------------------------------------------------------------------------

## ADR-007 --- Architecture before implementation

**Status:** Accepted

### Decision

No large-scale implementation begins until M0 architecture is reviewed.

### Rationale

The project is security-sensitive and protocol mistakes are expensive to
retrofit.

------------------------------------------------------------------------

## ADR process

Every significant design decision should record:

-   context
-   alternatives
-   decision
-   rationale
-   security consequences
-   operational consequences
-   status

See [[08-implementation/03-Code-Review-Checklist]].

------------------------------------------------------------------------

## ADR-008 --- Out-of-band observability

**Status:** Proposed

### Decision

Add an Observability & Security Telemetry subsystem using
Prometheus-compatible metrics and structured logs. Grafana is the
preferred visualization layer and Loki is the preferred initial log
aggregation candidate.

### Rationale

The project benefits from reproducible operational and security
demonstrations, especially for replay rejection, authentication
failures, tampering, routing behavior and node health.

### Security constraint

Observability must remain outside the E2EE data path and must never
receive plaintext messages, private keys, session keys or other secret
cryptographic material.

### Availability constraint

The mesh must continue functioning when Grafana, Prometheus or log
collection is unavailable.

### Consequence

The project gains measurable and demonstrable runtime behavior, but must
also enforce telemetry sanitization, bounded logging and metadata
minimization.

### Scope constraint

Do not add a large observability stack solely for visual complexity.
Additional systems require an explicit architectural justification.

See [[02-architecture/07-Observability-and-Security-Telemetry]].

------------------------------------------------------------------------

## ADR-009 --- Go as implementation language

**Status:** Accepted

Go is approved for the prototype. The project does not require Rust.
The implementation must use established cryptographic libraries and
must not implement cryptographic primitives from scratch.

------------------------------------------------------------------------

## ADR-010 --- Protocol Buffers serialization

**Status:** Accepted

Protocol Buffers is approved as the wire serialization format.

The Protobuf schema is part of the protocol specification and must be
versioned and tested for malformed input, compatibility behavior and
field-size limits.

------------------------------------------------------------------------

## ADR-011 --- Managed flooding for prototype routing

**Status:** Accepted with safeguards

Managed flooding is approved because the project prioritizes educational
clarity, resilience and correctness over routing efficiency.

The implementation must use bounded duplicate detection, TTL, resource
limits and malformed-input handling.

------------------------------------------------------------------------

## ADR-012 --- Signed ephemeral X25519 handshake

**Status:** Accepted / implementation gated to M3
specification

The project will use distinct Ed25519 identity keys and X25519 ephemeral
keys. Ed25519 signatures authenticate the X25519 exchange.

The exact canonical transcript, wire messages, KDF context, session ID,
handshake replay behavior and failure state machine must be specified
before implementation.

See [[00-governance/01-Architecture-Review-M0]].

## M2 Decision Record

### ADR — Long-Term Identity Implementation

**Decision:** Implement node identity with Go's standard `crypto/ed25519` and
entropy from `crypto/rand`.

**Rationale:** Ed25519 provides the approved long-term signing identity
primitive without custom cryptographic implementation. The public identity is
represented by its exact 32-byte public key.

**Security boundary:** Private keys are encapsulated inside `NodeIdentity`;
callers receive defensive copies of public identity material, and loaded
private keys are copied before cryptographic derivation.

**Status:** Implemented and audited in M2.1/M2.1-B/M2.1-C.
