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
