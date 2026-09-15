# Security Evidence Matrix

## Status

**M10.2 — PASS / CLOSED**

This matrix maps implemented security claims to the strongest available project evidence.

| Security claim | Primary evidence | Scope | Status |
|---|---|---|---|
| Endpoint application confidentiality | M9.1 relay-blindness integration | Alice → Bob → Carol evaluated path | Verified |
| AEAD integrity | M9.2 ciphertext-tampering test | Tested ciphertext modification | Verified |
| Session identity authentication | M3/M9 signature and identity tests | Evaluated handshake substitutions | Verified |
| Replay resistance | M3/M9 replay tests | Session sequence/replay behavior | Verified |
| Malformed packet rejection | M7/M9 validation tests | Tested malformed/oversized cases | Verified |
| TTL containment | M6/M9 routing tests | Tested TTL behavior | Verified |
| Duplicate forwarding suppression | M6/M9 routing tests | `(source_node, packet_id)` cache behavior | Verified |
| Pending handshake bound | M9.2 resource test | 15 attempts; configured maximum 10 | Verified |
| Relay/application-layer separation | M8.1/M9.1 evidence | Routing does not decrypt APP_DATA | Verified |
| GUI secure workflow | M9.3 manual verification | Two-node desktop workflow | Verified manually |
| Anonymity | No corresponding implementation/evidence | Explicit non-goal | Not claimed |
| Metadata hiding | No corresponding implementation/evidence | Explicit non-goal | Not claimed |
| Global traffic-analysis resistance | No corresponding implementation/evidence | Explicit non-goal | Not claimed |
| Guaranteed delivery | No corresponding implementation/evidence | Explicit non-goal | Not claimed |
| Post-compromise security | No corresponding implementation/evidence | Explicit limitation | Not claimed |
| Production-grade DoS resistance | No corresponding implementation/evidence | Explicit limitation | Not claimed |
| Active Prometheus exporter | No implementation | Deferred | Not implemented |
| Clean Docker regression | External environment prevented verification | M9 limitation | Not verified |

## Evidence Interpretation Rule

"Verified" means that the corresponding behavior was exercised by the recorded test or demonstration within its stated scope. It does not mean that every possible attack or deployment condition has been exhaustively evaluated.
