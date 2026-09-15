# M9 Evidence Manifest

| Evidence | Type | Source | Status |
|---|---|---|---|
| Environment capture | Environment | Recorded M9 project state | Recorded |
| Topology A | Mermaid | M9.1 approved topology | Recorded |
| Topology B | Mermaid | M9.3 GUI topology | Recorded |
| Repository tests | Automated | M9 validation evidence | PASS |
| Race tests | Automated | M9 validation evidence | PASS |
| Static analysis | Automated | `go vet ./...` | PASS |
| Backend build | Automated | `go build ./...` | PASS |
| GUI build | Automated | `go build -tags gui -o meshchat-gui ./cmd/meshchat-gui` | PASS |
| Relay blindness | Integration / demonstration | M9.1 approved evidence | PASS / CLOSED |
| Wrong PeerID | Adversarial | M9.2 test | PASS |
| Ciphertext tampering | Adversarial | M9.2 test | PASS |
| AEAD replay | Adversarial | M9.2 test | PASS |
| Invalid signature | Adversarial | M9.2 test | PASS |
| Malformed packet | Adversarial | M9.2 test | PASS |
| TTL enforcement | Adversarial | M9.2 test | PASS |
| PacketID duplicate suppression | Adversarial | M9.2 test | PASS |
| Pending-handshake bound | Adversarial/resource | M9.2 production-bound test | PASS |
| GUI identity/session/messaging scenarios | Manual demonstration | Project Overseer desktop verification | PASS / CLOSED |
| GUI screenshots | Visual artifact | Not present in supplied documentation archive | NOT PACKAGED |
| M9.4 Overseer audit | Governance / documentation audit | Project Overseer review of synchronized package | PASS / CLOSED |
| Docker regression | Environment-dependent | External proxy/DNS restriction | NOT VERIFIED |
| Prometheus exporter | Implementation feature | Project scope | NOT IMPLEMENTED |

## Evidence rule

No missing raw artifact is represented as if it existed. Where a result is known
from an approved acceptance record but raw stdout or screenshots are not present
in the supplied archive, the provenance is stated explicitly.
