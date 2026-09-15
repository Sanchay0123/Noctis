# Architecture Review — M9

**Status: 🟢 APPROVED / M9 CLOSED**

## Scope

M9 is the demonstration and security-validation milestone. It validates the
implemented system through three complementary phases:

- **M9.1 — E2EE / relay blindness**
- **M9.2 — adversarial security demonstrations**
- **M9.3 — live GUI demonstration**
- **M9.4 — evidence assembly and documentation synchronization**

## Approved constraints

M9 does not change the frozen cryptographic construction, packet schema, routing
model, or GUI architecture. M6 remains bounded managed flooding. ChaCha20-Poly1305
remains the application AEAD.

The relay demonstration uses Alice → Bob → Carol. Bob forwards ciphertext and
routing metadata; endpoint application-session ownership remains at Alice and
Carol in the demonstrated implementation.

The GUI demonstration uses a separate Alice ↔ Bob two-node topology.

## M9.2 security demonstrations

The approved adversarial set is:

1. Wrong PeerID
2. Ciphertext tampering
3. AEAD replay
4. Invalid signature
5. Malformed packet
6. TTL enforcement
7. PacketID duplicate suppression
8. Pending-handshake resource bound

The pending-handshake test must use the actual production bound of 10.

## M9.3 evidence provenance

The automated implementation environment was headless. The final GUI behavior
was manually verified by the Project Overseer on the desktop. No fabricated
GUI evidence is permitted.

## Non-goals

M9 does not establish anonymity, complete metadata hiding, global traffic-analysis
resistance, endpoint compromise resistance, guaranteed delivery, perfect
security, or immunity from denial of service.

## Gate

M9 architecture was approved before M9.1–M9.3 execution. M9.4 documentation and
evidence assembly was subsequently audited by the Project Overseer. The audit found
the evidence package consistent with the approved M9.1–M9.3 results and the M9
constraints. M9 is therefore formally closed.
