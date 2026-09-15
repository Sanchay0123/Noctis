# Security Analysis — Implemented System

## Status

**M10.2 — PASS / CLOSED**

This document analyzes security properties demonstrated by the implemented system and the evidence recorded through M9. It distinguishes demonstrated properties from design intent and from threats that remain outside the evaluated scope.

## 1. Security Objectives

The implemented system targets:

- confidentiality of application content between communicating endpoints
- integrity and authenticity of protected application messages
- authentication of cryptographic peer identities during session establishment
- replay resistance at the application session layer
- containment of malformed and oversized network input
- bounded mesh forwarding and duplicate suppression
- separation of routing metadata from endpoint message decryption

These objectives are evaluated within the academic prototype threat model.

## 2. Security Boundaries

The principal boundary is the endpoint/relay separation:

```text
Endpoint                    Relay                     Endpoint
Alice                       Bob                      Carol
  |                           |                         |
  |-- encrypt APP_DATA ----->|-- forward ciphertext ->|
  |                           |                         |
  |<----- endpoint session ---------------------------->|
```

A relay processes routing information and ciphertext. Application decryption occurs at the endpoint application/session boundary, not in `internal/routing`.

## 3. Threat-to-Control Analysis

| Threat | Implemented control | Evidence | Assessment |
|---|---|---|---|
| Passive observation of application traffic | ChaCha20-Poly1305 E2EE | M9.1 relay evidence; M3/M4 crypto tests | Demonstrated for the evaluated relay path |
| Ciphertext modification | AEAD authentication | M9.2 tampering test | Demonstrated |
| Message replay | Sequence numbers + 64-message receive window | M9.2 replay test; M3 tests | Demonstrated |
| Identity substitution | Ed25519-authenticated handshake transcripts | M3/M9 identity/signature tests | Demonstrated against evaluated substitutions |
| Forged handshake signature | Ed25519 signature verification | M9.2 invalid-signature test | Demonstrated |
| Malformed packets | Structural validation before normal processing | M9.2 malformed-packet test; M7 boundary evidence | Demonstrated |
| Oversized input | 64 KiB frame bound | M7/M9 evidence | Demonstrated |
| Routing loops | TTL decrement/termination | M9.2 TTL test; M6 evidence | Demonstrated |
| Duplicate propagation | `(source_node, packet_id)` cache key | M9.2 duplicate-suppression test | Demonstrated |
| Incomplete-handshake exhaustion | Production `MaxPendingHandshakes = 10` | M9.2 resource-bound test | Demonstrated for the evaluated bound |
| Relay access to endpoint plaintext | Endpoint-only decryption boundary | M9.1 relay-blindness evidence | Demonstrated for the evaluated topology |

## 4. Confidentiality

Application messages are encrypted with ChaCha20-Poly1305 using directional session keys derived from the authenticated X25519 exchange and HKDF-SHA-256.

M9.1 provides the strongest system-level evidence: an Alice → Bob → Carol path was exercised and the intermediate node observed routing metadata and ciphertext while Carol received the decrypted application message.

This demonstrates endpoint E2EE across the evaluated relay path. It does not prove confidentiality against a compromised endpoint, leaked session keys, or an attacker that already possesses endpoint secrets.

## 5. Integrity and Authenticity

ChaCha20-Poly1305 provides authenticated encryption for application data. The AEAD context binds the session identifier, sequence number and internally derived direction marker through AAD.

The M9.2 ciphertext-tampering demonstration verifies that modification does not result in successful authenticated decryption.

Long-term Ed25519 identities authenticate the X25519 handshake transcripts. Invalid-signature and identity-substitution tests provide negative evidence for the evaluated authentication boundary.

## 6. Replay Resistance

Replay protection exists at two distinct layers:

1. **Application message layer:** directional sequence numbers and a 64-message receive window reject duplicate or sufficiently old authenticated messages.
2. **Mesh packet layer:** `(source_node, packet_id)` duplicate suppression limits repeated forwarding of the same packet within the managed flooding cache.

These mechanisms solve different problems and must not be conflated. A PacketID duplicate is not the same security mechanism as an application sequence-number replay.

## 7. Resource Containment

The runtime applies several bounds:

- maximum accepted frame size: 64 KiB
- maximum active peers: 50
- maximum pending handshakes: 10
- handshake timeout: 5 seconds
- bounded routing duplicate cache
- bounded forwarding queues
- TTL-based forwarding termination

M9.2 directly exercised the production pending-handshake limit with 15 incomplete connection attempts and observed simultaneous admission bounded at 10 or below.

These controls provide bounded resource consumption for the evaluated conditions; they do not constitute complete denial-of-service resistance.

## 8. Relay Blindness

The relay-blindness property is an architectural consequence of separating endpoint application sessions from mesh forwarding.

Bob forwards APP_DATA without accessing the Alice–Carol endpoint application session. The demonstrated M9.1 topology confirms that Carol can decrypt and deliver the message after multi-hop forwarding.

The precise claim is:

> Bob's implemented session manager does not contain the Alice–Carol endpoint application session required to decrypt the packet.

This is not a claim of mathematical impossibility under arbitrary endpoint compromise or key leakage.

## 9. Authentication Caveat

The implementation supports managed/pre-shared public-key knowledge and TOFU. TOFU does not authenticate first contact against an active man-in-the-middle attacker when the initial identity is substituted before it is trusted.

Accordingly, the project should claim authentication resistance only within the configured trust model and the evaluated identity/signature substitution tests. It should not claim universal MITM immunity.

## 10. Forward Secrecy and Key Lifecycle

Fresh ephemeral X25519 keys are used for session establishment and session keys are maintained in memory for the active session. Restart requires a new handshake.

The project does not claim post-compromise security. Forward-secrecy claims remain conditional on fresh ephemeral key use, appropriate private-key lifecycle handling and uncompromised endpoints.

## 11. Security Properties Not Demonstrated

The evidence does **not** establish:

- anonymity
- complete metadata hiding
- global traffic-analysis resistance
- guaranteed message delivery
- resistance to a globally coordinated denial-of-service attack
- complete Sybil resistance
- routing optimality
- endpoint-compromise resistance
- post-compromise security
- production deployment readiness

These are explicit limitations, not failed security controls.

## 12. Observability Boundary

Telemetry is supporting and out-of-band. The current implementation contains telemetry hooks/no-op behavior but no active Prometheus exporter or metric-vector implementation.

The security-critical message path does not depend on telemetry availability. Sensitive material such as plaintext, private keys, session keys and passwords is not intended to be exported as telemetry.

## 13. Evidence Strength

The evidence has different levels of scope:

- **Unit/negative tests:** establish behavior of individual cryptographic, protocol and runtime controls under specified inputs.
- **Integration tests:** establish composition of crypto, transport, routing and application boundaries.
- **M9.1 relay test:** establishes the evaluated multi-hop endpoint E2EE boundary.
- **M9.2 adversarial tests:** establish the listed rejection and resource-bound behaviors.
- **M9.3 manual GUI verification:** establishes the observed two-node desktop workflow.

A passing test is evidence for its tested condition, not proof of all possible attacks.

## 14. Overall Security Assessment

Within the declared academic threat model, the implementation demonstrates a coherent endpoint E2EE construction integrated with bounded mesh forwarding. The strongest evidence concerns authenticated session establishment, authenticated encryption, replay handling, relay-blind multi-hop message delivery, malformed-input rejection and bounded handshake admission.

The implementation should be presented as a security-focused academic prototype with explicit boundaries rather than as an anonymous, production-ready or universally attack-resistant messaging system.

## Related Records

- [[03-security/01-Threat-Model]]
- [[03-security/02-Security-Goals]]
- [[03-security/04-Security-Controls]]
- [[03-security/05-Security-Limitations]]
- [[07-testing/16-M9-Attack-Evidence]]
- [[07-testing/17-M9-Relay-Blindness.log]]
- [[07-testing/18-M9-GUI-Evidence]]
- [[07-testing/19-M9-Evidence-Manifest]]
