# Network Test Matrix

  Scenario                       Expected behavior
  ------------------------------ ----------------------------------
  Two nodes connect              Connection established
  Peer authentication succeeds   Peer becomes trusted for session
  Peer authentication fails      Connection/session rejected
  Relay forwards packet          Destination receives
  Relay inspects packet          No E2EE plaintext available
  TTL expires                    Packet dropped
  Duplicate packet               Not repeatedly forwarded
  Route loop                     Eventually terminated
  Relay disconnects              Failure detected
  Route unavailable              Controlled delivery failure
  Malformed packet               Node remains alive
  Connection flooding            Bounded resource behavior
  Reconnect                      Explicit session policy followed


## M5 Runtime Networking Matrix

| Scenario | Expected behavior | Evidence |
|---|---|---|
| Outbound direct connection | Authenticated peer becomes established | PASS |
| Inbound direct connection | Expected remote identity is authenticated before registration | PASS |
| Pending handshake limit | Excess sockets rejected without listener blocking | PASS |
| Active peer limit | New non-duplicate peer rejected at capacity | PASS |
| Concurrent dial limit | Excess dials rejected and slots released | PASS |
| Oversized initial frame | Rejected before body allocation | PASS |
| Duplicate direct connection | Smaller identity's initiated connection retained | PASS |
| Stale incumbent cleanup | Replacement cannot be removed by stale goroutine | PASS |
| Blocking telemetry recorder | Manager shutdown remains independent | PASS |
| Failing telemetry recorder | Network lifecycle remains operational | PASS |
| Simultaneous Dave↔Bob Docker crossover | Exactly one connection survives | PASS |

M6-only behaviors such as multi-hop forwarding, route discovery, TTL forwarding
and network-wide PacketID suppression remain intentionally untested here.

## M6 Network Test Matrix

| Area | Required evidence | Result |
|---|---|---|
| Handshake correlation | simultaneous same-remote candidates; wrong/late/duplicate RESP | PASS |
| TTL | 0, 1, 2, 16, 32 and >32 boundaries | PASS |
| Routing | one-, two- and three-hop delivery | PASS |
| Cycles | cyclic flooding terminates | PASS |
| Packet validation | malformed, oversized, unknown fields | PASS |
| Duplicate cache | first, duplicate, source separation, capacity, expiration, concurrency | PASS |
| Forward queue | count/byte bounds and rollback accounting | PASS |
| Relay confidentiality | relay cannot decrypt endpoint ciphertext | PASS |
| Peer loss | transport loss does not invalidate endpoint session | PASS |
| Flood bound | sustained unique-packet load remains bounded | PASS |
| Race safety | repository race-detector validation | PASS |
| Container demo | Docker multi-hop E2EE delivery | PASS |

Detailed evidence: [[07-testing/08-M6-Acceptance-Evidence]].
