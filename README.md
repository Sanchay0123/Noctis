# E2EE Mesh Chat --- Documentation Hub

> **Status:** M8.3 complete / verified / approved\
> **Project type:** Academic cybersecurity and networking prototype

## Purpose

This vault documents the design of an **End-to-End Encrypted Messaging
System over a Decentralized Mesh Network**.

The project combines:

-   cybersecurity
-   applied cryptography
-   computer networking
-   distributed systems
-   peer-to-peer communication
-   secure software engineering

The documentation is deliberately separated into small, linked Markdown
notes so it can be used directly in **Obsidian**.

## Security priority

> **Security \> Correctness \> Simplicity \> Performance \> Features**

No component is considered secure merely because an implementation agent
reports that it is complete. Security claims require implementation
evidence and tests.

## Start here

-   [[01-project/01-Project-Overview]]
-   [[01-project/02-Requirements]]
-   [[03-security/01-Threat-Model]]
-   [[02-architecture/01-System-Architecture]]
-   [[02-architecture/02-Architecture-Decisions]]
-   [[04-protocol/01-Protocol-Overview]]
-   [[05-cryptography/01-Cryptographic-Architecture]]
-   [[07-testing/01-Testing-Strategy]]
-   [[08-implementation/01-Milestones]]

## Documentation map

### Governance

-   [[00-governance/00-Project-State]]
-   [[00-governance/01-Vault-Index]]
-   [[00-governance/03-Architecture-Review-M5]]
-   [[00-governance/06-Architecture-Review-M8]]
-   [[00-governance/05-Architecture-Review-M7]]
-   [[00-governance/04-Architecture-Review-M6]]

### Project

-   [[01-project/01-Project-Overview]]
-   [[01-project/02-Requirements]]
-   [[01-project/03-Goals-and-Non-Goals]]
-   [[01-project/04-Assumptions-and-Constraints]]

### Architecture

-   [[02-architecture/01-System-Architecture]]
-   [[02-architecture/02-Architecture-Decisions]]
-   [[02-architecture/03-Component-Architecture]]
-   [[02-architecture/04-Data-Flow]]
-   [[02-architecture/05-Trust-Boundaries]]
-   [[02-architecture/06-Deployment-Architecture]]

### Security

-   [[03-security/01-Threat-Model]]
-   [[03-security/02-Security-Goals]]
-   [[03-security/03-Attack-Surface]]
-   [[03-security/04-Security-Controls]]
-   [[03-security/05-Security-Limitations]]
-   [[03-security/06-Security-Review-Checklist]]

### Protocol

-   [[04-protocol/01-Protocol-Overview]]
-   [[04-protocol/02-Node-Identity]]
-   [[04-protocol/03-Session-Establishment]]
-   [[04-protocol/04-Message-Format]]
-   [[04-protocol/05-Mesh-Packet-Format]]
-   [[04-protocol/06-Replay-Protection]]
-   [[04-protocol/07-Protocol-State-Machines]]
-   [[04-protocol/08-Error-Handling]]

### Cryptography

-   [[05-cryptography/01-Cryptographic-Architecture]]
-   [[05-cryptography/02-Identity-and-Signatures]]
-   [[05-cryptography/03-Key-Agreement]]
-   [[05-cryptography/04-Key-Derivation]]
-   [[05-cryptography/05-AEAD-and-Nonce-Management]]
-   [[05-cryptography/06-Key-Lifecycle]]
-   [[05-cryptography/07-Cryptographic-Failure-Modes]]

### Networking

-   [[06-networking/01-Mesh-Architecture]]
-   [[06-networking/02-Peer-Discovery]]
-   [[06-networking/03-Routing]]
-   [[06-networking/04-Forwarding]]
-   [[06-networking/05-TTL-and-Loop-Prevention]]
-   [[06-networking/06-Connection-Management]]
-   [[06-networking/07-M6-Handshake-Correlation]]
-   [[06-networking/08-M6-Resource-Bounds]]

### Testing

-   [[07-testing/01-Testing-Strategy]]
-   [[07-testing/02-Cryptographic-Test-Matrix]]
-   [[07-testing/03-Protocol-Test-Matrix]]
-   [[07-testing/04-Network-Test-Matrix]]
-   [[07-testing/05-Security-Demonstrations]]
-   [[07-testing/06-Test-Evidence-Standard]]
-   [[07-testing/08-M6-Acceptance-Evidence]]
-   [[07-testing/09-M7-Acceptance-Evidence]]
-   [[07-testing/10-M8.1-Acceptance-Evidence]]
-   [[07-testing/11-M8.2-Acceptance-Evidence]]
-   [[07-testing/12-M8.3-Acceptance-Evidence]]

### Implementation

-   [[08-implementation/01-Milestones]]
-   [[08-implementation/02-Task-Template]]
-   [[08-implementation/03-Code-Review-Checklist]]
-   [[08-implementation/04-Antigravity-Workflow]]
-   [[08-implementation/05-Definition-of-Done]]

### Academic

-   [[09-academics/01-Abstract-Draft]]
-   [[09-academics/02-Problem-Statement]]
-   [[09-academics/03-Objectives]]
-   [[09-academics/04-Methodology]]
-   [[09-academics/05-Security-Analysis-Structure]]
-   [[09-academics/06-Results-Structure]]
-   [[09-academics/07-Limitations]]
-   [[09-academics/08-Future-Work]]
-   [[09-academics/09-Viva-Question-Bank]]

## Decision status legend

-   **Proposed** --- reasonable design candidate; not approved.
-   **Accepted** --- selected for implementation.
-   **Implemented** --- code exists.
-   **Verified** --- tests/evidence support the claim.
-   **Approved** --- technical lead has reviewed implementation and
    evidence.
-   **Rejected** --- explicitly not permitted.
-   **Deferred** --- intentionally postponed.

## Important rule

If a document contains an unresolved design decision, it must say so.
Never silently convert a proposal into an implementation requirement.

## Supporting subsystem

- [[02-architecture/07-Observability-and-Security-Telemetry]] —
  Observability, security telemetry, Prometheus/Grafana/Loki boundary
  and telemetry security policy.

Observability is deliberately out-of-band. The mesh must continue to
operate if monitoring is unavailable.

## Current Implementation Status

**M8.3 — Functional GUI Integration: 🟢 Complete / Verified / Approved**

The repository has completed the approved Ed25519 identity, X25519 authenticated session establishment, transcript-bound HKDF-SHA-256 key derivation, ChaCha20-Poly1305 application protection, replay control, bounded TCP transport, authenticated multi-peer networking, bounded multi-hop forwarding, security hardening, the ApplicationService boundary and functional Fyne GUI integration.

M8 preserves endpoint E2EE: relays can route ciphertext but cannot decrypt endpoint application data. The GUI is isolated behind ApplicationService and does not directly access crypto, session, mesh or routing internals.

M8.3 application integration passed repository tests, race validation, vetting and build validation. Final live two-node GUI acceptance passed in both first-initiator directions, including secure-session establishment, bidirectional encrypted messaging and idempotent duplicate Add Peer behavior. Docker regression remains unverified because external dependency resolution is blocked by proxy/DNS restrictions.

**M8.3 is complete, live accepted and approved. M9 requires separate authorization.**

See [[00-governance/00-Project-State]], [[00-governance/06-Architecture-Review-M8]] and [[08-implementation/01-Milestones]] for the current gate.
