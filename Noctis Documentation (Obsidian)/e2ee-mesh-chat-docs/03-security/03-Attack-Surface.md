# Attack Surface

## External inputs

-   transport connections
-   packet bytes
-   peer identifiers
-   discovery messages
-   routing fields
-   serialized cryptographic fields
-   user-provided message content

## Sensitive operations

-   private-key loading
-   signature verification
-   X25519 key agreement
-   HKDF derivation
-   AEAD encryption/decryption
-   replay-state updates
-   route-table updates
-   connection creation

## High-risk attack classes

### Cryptographic

-   key substitution
-   signature bypass
-   invalid key handling
-   nonce reuse
-   wrong-key acceptance
-   transcript confusion

### Protocol

-   malformed serialization
-   oversized packets
-   field confusion
-   version downgrade
-   replay
-   duplicate identifiers

### Networking

-   route loops
-   flooding
-   connection exhaustion
-   malicious discovery
-   stale routes

### Application

-   plaintext logging
-   secret leakage in UI
-   injection into logs
-   uncontrolled message size

## Attack-surface rule

Every field received from another node must be considered
attacker-controlled until validated.

## Observability attack surface

Telemetry introduces additional security-sensitive surfaces:

- log sinks
- metrics endpoints
- telemetry serialization
- metric labels
- structured log fields
- dashboard access
- log aggregation/storage
- malformed event data
- log injection
- excessive telemetry volume

### Primary risks

- plaintext leakage through logs
- private/session key leakage
- metadata leakage
- unbounded log growth
- attacker-controlled log injection
- telemetry-based resource exhaustion
- observability service compromise affecting the mesh

### Boundary rule

Observability must not be able to decrypt application messages or
access cryptographic secret material merely because it has access to
node events.

See [[02-architecture/07-Observability-and-Security-Telemetry]].

## M2 Identity Attack Surface

M2 introduces local identity-key handling as an explicit attack surface.

Relevant controls include:

- cryptographically secure identity generation
- strict key-length validation
- private-key encapsulation
- defensive copying on identity loading and public-key retrieval
- safe handling of malformed public keys and signatures
- no secret material in logs, telemetry, protocol packets, or committed fixtures

M2 does not yet address the active network attacks associated with session
establishment; those remain in the M3 threat-control scope.


## M5 Runtime Networking Attack Surface

M5 makes the direct networking layer an explicit runtime attack surface.
Relevant threats and controls include:

| Threat | M5 control |
|---|---|
| Connection exhaustion | Pending-handshake, active-peer and dial limits |
| Oversized frame allocation | 64 KiB bound checked before allocation |
| Slow/incomplete handshake | Handshake deadline |
| Dial resource retention | Dial timeout and bounded concurrency |
| Simultaneous duplicate connections | Deterministic identity-based arbitration |
| Stale peer cleanup race | Atomic registry replacement and identity-safe cleanup |
| Telemetry blocking | Bounded queue and worker outside manager WaitGroup |
| Telemetry recorder panic/failure | Recorder boundary recovery / failure isolation |

M5 does not claim protection against routing attacks, global traffic analysis,
perfect anonymity or compromised endpoints. Those properties remain outside the
approved scope.
