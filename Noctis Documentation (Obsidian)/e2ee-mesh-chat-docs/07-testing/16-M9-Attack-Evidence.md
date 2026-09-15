# M9.2 — Security Attack Demonstrations

**Status: 🟢 PASS / CLOSED**

M9.2 validates eight mandatory adversarial/security demonstrations against the
approved implementation.

| # | Demonstration | Expected behavior | Result |
|---|---|---|---|
| 1 | Wrong PeerID | Authenticated identity mismatch is rejected | PASS |
| 2 | Ciphertext tampering | AEAD authentication rejects modified ciphertext | PASS |
| 3 | AEAD replay | Previously accepted sequence is rejected by replay protection | PASS |
| 4 | Invalid signature | Handshake authentication fails | PASS |
| 5 | Malformed packet | Structural validation rejects malformed input | PASS |
| 6 | TTL enforcement | Expired/non-forwardable packets are dropped according to TTL rules | PASS |
| 7 | PacketID duplicate suppression | Duplicate packet is suppressed by the bounded seen-cache | PASS |
| 8 | Pending-handshake resource bound | Concurrent incomplete handshakes are bounded by production admission limits | PASS |

## Pending-handshake resource bound

The authoritative production configuration is:

- `MaxPendingHandshakes = 10`
- maximum active peers = 50
- handshake timeout = 5 seconds
- pending-handshake semaphore capacity = 10

The M9.2 remediation test attempted **15 concurrent incomplete TCP handshakes**.
The observed number of simultaneously admitted pending handshakes was **no more
than 10**.

This is evidence of a production-configured resource bound; it is not an
artificial test configuration of one pending handshake.

## Evidence boundary

These demonstrations are test-only security validation. They do not create a
production attack/debug API and do not establish immunity from all denial of
service, routing, traffic-analysis, or endpoint-compromise attacks.
