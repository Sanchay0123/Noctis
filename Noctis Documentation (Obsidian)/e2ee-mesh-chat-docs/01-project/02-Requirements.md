# Requirements

## Functional requirements

### R-F01 --- Node identity

Each node shall have a cryptographic identity based on a long-term
Ed25519 keypair.

### R-F02 --- Peer discovery

Nodes shall be able to discover peers using a defined discovery
mechanism.

### R-F03 --- Peer connectivity

Nodes shall establish and maintain connections to peers.

### R-F04 --- Authenticated session establishment

Communicating peers shall authenticate the identities involved in
session establishment.

### R-F05 --- End-to-end encryption

Application message content shall be encrypted so relay nodes cannot
decrypt it.

### R-F06 --- Multi-hop forwarding

A message shall be able to traverse one or more intermediate nodes.

### R-F07 --- Message integrity

Modified ciphertext or authenticated metadata shall be rejected.

### R-F08 --- Replay protection

Previously accepted messages shall not be accepted again as new
messages.

### R-F09 --- Routing controls

Packets shall support TTL/hop limits and duplicate detection.

### R-F10 --- Reconnection

The system shall provide basic handling for peer disconnection and
reconnection.

### R-F11 --- User interface

The system shall provide a usable interface for conversations, peers and
connection state.

### R-F12 --- Testing

Automated tests shall cover cryptographic, protocol, networking and
adversarial cases.

### R-F13 --- Demonstration

A reproducible multi-node demonstration shall exist.

### R-F14 --- Documentation

Architecture, protocol, threat model, cryptography, networking, testing
and limitations shall be documented.

### R-F15 --- Observability

The system shall provide security-safe operational and security telemetry
sufficient to support debugging, performance evaluation and reproducible
demonstrations.

### R-F16 --- Containerized reproducibility

The application and demonstration environment shall be reproducible
through containerized configuration.

### R-F17 --- Repository synchronization

Implementation, tests, documentation and deployment configuration shall
be maintained in the project Git repository, with meaningful milestone
changes committed and pushed to the canonical GitHub repository.

## Security requirements

### R-S01

No custom cryptographic primitive shall be implemented.

### R-S02

Only established cryptographic library implementations shall be used.

### R-S03

AEAD nonce reuse under the same key shall be prevented by design.

### R-S04

The key exchange shall be authenticated against peer identities.

### R-S05

Private keys shall not appear in normal logs.

### R-S06

Application plaintext shall not be logged by relay nodes.

### R-S07

Malformed packets shall not cause uncontrolled process failure.

### R-S08

Security failures shall fail closed where practical.

### R-S09

Replay state shall be maintained by the receiver rather than relying
only on timestamps.

### R-S10

Security claims in documentation shall be traceable to implementation
and test evidence.

### R-S11

Observability shall not expose application plaintext, private keys,
session keys, passwords or other secret cryptographic material.

### R-S12

Observability shall remain outside the E2EE message path. Failure of
monitoring or log aggregation shall not prevent core message delivery.

### R-S13

Telemetry shall use bounded, structured and sanitized data to reduce
metadata leakage and log-injection risk.

## Quality requirements

-   deterministic and reproducible tests
-   minimal dependency footprint
-   clear separation of concerns
-   explicit protocol versioning
-   documented error handling
-   observable but non-secret debug information
-   reproducible containerized execution
-   synchronized implementation and documentation in GitHub
-   cross-platform development where practical

See [03-Goals-and-Non-Goals](03-Goals-and-Non-Goals.md) and
[06-Test-Evidence-Standard](../07-testing/06-Test-Evidence-Standard.md).

## M2 Implementation Status

The identity requirements are now implemented and tested:

- Long-term node identity uses Ed25519.
- Public identity representation is exactly 32 bytes.
- Identity generation uses cryptographically secure randomness.
- Private identity material remains local to the cryptographic identity abstraction.
- Signing and verification use the standard Ed25519 implementation.
- Identity functionality is isolated from session establishment, routing, transport, and application messaging.

X25519 session establishment and encrypted messaging are now implemented
within M3. Direct networking, mesh routing and end-to-end delivery remain
future milestones.
