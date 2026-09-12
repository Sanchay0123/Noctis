# E2EE Mesh Chat --- Documentation Hub

> **Status:** M4 complete / M5 gated\
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

### Testing

-   [[07-testing/01-Testing-Strategy]]
-   [[07-testing/02-Cryptographic-Test-Matrix]]
-   [[07-testing/03-Protocol-Test-Matrix]]
-   [[07-testing/04-Network-Test-Matrix]]
-   [[07-testing/05-Security-Demonstrations]]
-   [[07-testing/06-Test-Evidence-Standard]]

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

**M4 — Direct Secure Messaging Integration: 🟢 Complete / Approved**

The repository has completed and verified the long-term Ed25519 identity,
fresh X25519 session establishment, transcript-bound HKDF-SHA-256 key
derivation, ChaCha20-Poly1305 application message protection, sequence-based
replay control, and direct TCP secure messaging integration.

M4 provides bounded TCP framing, packet validation, direct handshake
orchestration, session/identity binding, bidirectional encrypted messaging and
transport-path tamper rejection. Mesh routing and broader runtime networking
remain future milestones.

**M5 is NOT STARTED / GATED.**

See [[00-governance/00-Project-State]] and
[[08-implementation/01-Milestones]] for the current gate.
