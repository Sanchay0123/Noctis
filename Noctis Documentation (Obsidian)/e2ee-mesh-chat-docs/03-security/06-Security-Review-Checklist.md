# Security Review Checklist

## M2 and M3 reviewed scope

-   [x] Established library primitives only for Ed25519
-   [x] Ed25519 identity handling reviewed
-   [x] Ed25519 malformed-input behavior reviewed
-   [x] Ed25519 key ownership/defensive-copy behavior reviewed
-   [x] Ed25519 concurrent-signing behavior tested

The remaining unchecked items are primarily networking, integration,
deployment and final security-demonstration controls and must not be
interpreted as completed by the M3 gate.


## Cryptography

-   [ ] Established library primitives only
-   [ ] Ed25519 identity handling reviewed
-   [ ] X25519 key exchange reviewed
-   [ ] Identity binding reviewed
-   [ ] Transcript binding reviewed
-   [ ] HKDF context separation reviewed
-   [ ] AEAD usage reviewed
-   [ ] Nonce generation reviewed
-   [ ] Nonce reuse impossible by design
-   [ ] Key lifecycle reviewed
-   [ ] Invalid cryptographic inputs fail safely

## Protocol

-   [ ] Version checked
-   [ ] Message types validated
-   [ ] Field lengths bounded
-   [ ] Unknown fields handled intentionally
-   [x] Authentication occurs before trust-sensitive session-state commitment
-   [x] Replay protection enforced for authenticated application messages
-   [x] Duplicate detection enforced for session messages
-   [ ] Error messages do not leak secrets

## Networking

-   [ ] TTL enforced
-   [ ] Routing loops bounded
-   [ ] Duplicate packets bounded
-   [ ] Connection limits exist
-   [ ] Timeouts exist
-   [ ] Malformed packets cannot crash the node
-   [ ] Resource use is bounded

## Privacy

-   [ ] Plaintext absent from relay logs
-   [ ] Private keys absent from logs
-   [ ] Session keys absent from logs
-   [ ] Sensitive metadata documented
-   [ ] Debug mode reviewed

## Evidence

-   [ ] Security tests exist
-   [ ] Negative tests exist
-   [ ] Integration tests exist
-   [ ] Demonstration evidence captured
-   [ ] Documentation claims match actual evidence


## M3 review evidence

M3.1/M3.2/M3.3 were reviewed against the frozen protocol construction.
Negative tests cover transcript tampering, identity and ephemeral
substitution, malformed inputs, invalid ciphertext/AAD, duplicate and
out-of-window sequence numbers, unauthenticated replay-state poisoning,
cross-session and cross-direction use, and sequence exhaustion.

The reported containerized validation sequence passed formatting, vetting,
tests, race detection and build checks.

The following remain future review scope: direct-network integration,
multi-node behavior, mesh forwarding/routing controls, basic DoS resistance,
telemetry leakage under the integrated runtime, and end-to-end security
demonstrations.


## M7 Security-Hardening Review

- [x] Established cryptographic primitives remain unchanged and library-backed
- [x] X25519 session construction remains frozen and reviewed
- [x] Identity binding regression tested
- [x] Transcript binding remains frozen
- [x] HKDF context separation remains frozen
- [x] ChaCha20-Poly1305 usage remains frozen and reviewed
- [x] Nonce/sequence exhaustion behavior remains bounded
- [x] Concurrent replay state is serialized
- [x] Invalid cryptographic inputs fail safely
- [x] Protocol version checked at runtime boundary
- [x] Message types and oneof relationships validated
- [x] Field lengths bounded
- [x] Unknown Protobuf fields rejected
- [x] Router validates before cache/routing
- [x] Oversized router input rejected
- [x] TTL enforced
- [x] Routing loops bounded
- [x] Duplicate packet state bounded
- [x] Connection and handshake limits exist
- [x] Handshake timeouts verified
- [x] Malformed packets cannot crash the node in tested paths
- [x] Local resource use bounded
- [x] Peer identities removed from the telemetry recorder interface
- [x] No active dynamic Prometheus label implementation exists at M7
- [x] Telemetry remains out-of-band
- [x] Full repository race validation passed
- [x] M7 acceptance evidence recorded

See [05-Architecture-Review-M7](../00-governance/05-Architecture-Review-M7.md) and
[09-M7-Acceptance-Evidence](../07-testing/09-M7-Acceptance-Evidence.md).
