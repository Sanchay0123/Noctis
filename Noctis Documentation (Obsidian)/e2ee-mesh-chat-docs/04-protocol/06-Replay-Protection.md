# Replay Protection

## Threat

An attacker captures an accepted ciphertext and later sends it again.

Without replay protection:

```text
Alice -> Bob: TRANSFER/COMMAND/MESSAGE
Attacker captures packet
Attacker -> Bob: same packet again
```

The application may incorrectly process it twice.

## Required controls

Use:
- unique packet identifiers at the mesh layer
- per-direction sequence numbers at the session layer
- authenticated session identifiers
- receiver replay state
- a bounded replay window

Timestamps alone are insufficient.

## Receiver flow

```mermaid
flowchart TD
    P[Packet received] --> V[Validate structure]
    V --> R[Check replay window without mutation]
    R --> A[Authenticate AEAD]
    A -->|Failure| X[Reject without state change]
    A -->|Success| C[Commit replay state]
    C --> D[Deliver plaintext]
```

The M3.3 implementation uses a 64-message receive window and commits state
only after successful AEAD authentication.

## Requirements

- duplicate application ciphertext rejected
- sequence numbers outside the receive window rejected
- valid out-of-order ciphertext inside the window accepted once
- invalid ciphertext does not advance replay state
- cross-session ciphertext rejected
- cross-direction ciphertext rejected
- restart invalidates old session replay state
- mesh PacketID duplicate handling remains independent of session replay

## M6 Network Duplicate Suppression vs M4 Replay Protection

M6 duplicate suppression and M4 cryptographic replay protection solve different problems.

- **M6:** short-lived network duplicate cache, keyed by `(source_node, packet_id)`, used to terminate flooding loops and repeated propagation.
- **M4/M3:** authenticated sequence numbers and a 64-message replay window, bound to the endpoint E2EE session, used to reject cryptographic replays.

Neither mechanism replaces the other.
