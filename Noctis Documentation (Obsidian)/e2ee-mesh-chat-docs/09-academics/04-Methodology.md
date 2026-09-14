# Methodology

## Phase 1 --- Requirements and threat modeling

Define:

-   system requirements
-   assets
-   attackers
-   trust boundaries
-   security goals
-   non-goals

See [[03-security/01-Threat-Model]].

## Phase 2 --- Architecture

Define:

-   layers
-   components
-   data flows
-   protocol boundaries
-   cryptographic architecture
-   routing model

See [[02-architecture/01-System-Architecture]].

## Phase 3 --- Cryptographic implementation

Implement only established primitives through mature libraries.

**Current evidence boundary:** the Ed25519 long-term identity primitive and
the M3 authenticated session/message-protection layer are implemented and
audited within their approved scopes. Direct networking and mesh routing
remain future implementation stages.

Validate:

-   identity
-   key agreement
-   authentication
-   KDF
-   AEAD
-   replay controls

## Phase 4 --- Direct secure messaging integration

Integrate the approved cryptographic/session layer with bounded direct TCP
transport and prove secure two-node message delivery, packet validation,
session binding and transport-path tamper rejection before introducing routing
complexity.

M4 evidence includes framing tests, short-write tests, packet validation tests,
end-to-end bidirectional encrypted messaging and transport interception/tamper
tests. The M4 implementation does not establish mesh routing behavior.

## Phase 5 --- Mesh

Add:

-   peer discovery
-   routing
-   forwarding
-   TTL
-   duplicate detection
-   route failure behavior

## Phase 6 --- Security evaluation

Attack the system intentionally:

-   modify packets
-   replay packets
-   substitute identities
-   send malformed packets
-   attempt relay decryption

## Phase 7 --- Demonstration and analysis

Compare:

-   intended security properties
-   observed behavior
-   test evidence
-   limitations

## Phase 8 --- Documentation

Record implementation decisions, evidence and limitations without
upgrading unverified claims into facts.


### M3 evidence boundary

The cryptographic implementation was validated through unit, negative,
known-answer and race-detector tests in the controlled containerized Go
environment.

### M4 evidence boundary

The M4 implementation extends the evidence to the direct transport and secure
message integration boundary. Tests establish bounded framing, packet
validation, session/identity binding, bidirectional encrypted delivery and
transport-path tamper rejection. They do not establish mesh routing,
anonymity, metadata hiding or production-grade network resilience.


### M5 evidence boundary

M5 extends the evidence from direct secure messaging into runtime direct-network
management. The evaluation covers authenticated peer registration, inbound and
outbound connection lifecycle, resource limits, deterministic duplicate
connection arbitration, lifecycle race safety, telemetry isolation and a
containerized multi-node runtime. It does not establish multi-hop routing or
network-wide forwarding behavior.

## Phase 8A --- Application service and UI integration

After mesh security hardening, introduce a narrow application-service boundary
that exposes authenticated messaging to presentation code without exposing
cryptographic/session/routing internals. Verify that routing remains blind to
plaintext and that endpoint decryption occurs only at the application/session
boundary.

## Phase 8B --- GUI integration

Add the primary Fyne GUI behind an isolated GUI build tag. Validate application
service behavior, synchronized GUI state, event delivery and shutdown/lifecycle
handling. Separate GUI build/runtime environment evidence from headless backend
validation.

M8.2 does not by itself constitute complete observability implementation or
final demonstration evidence.
