# Security Controls

  Threat                      Control
  --------------------------- -----------------------------------------------
  Packet modification         AEAD authentication
  Identity substitution       Ed25519-authenticated key exchange
  Replay                      Message identifiers + sequence/replay state
  Route loops                 TTL + duplicate detection
  Malformed packet            Strict validation
  Key exposure through logs   Secret-aware logging policy
  Nonce reuse                 Explicit nonce/key lifecycle
  Unknown protocol versions   Version validation
  Resource exhaustion         Size limits, connection limits, bounded state
  Relay plaintext exposure    End-to-end encryption

## Defense-in-depth

No single control should be treated as solving every threat.

For example:

-   TLS-like transport protection does not replace E2EE.
-   signatures do not provide confidentiality.
-   AEAD does not solve routing loops.
-   TTL does not authenticate a sender.
-   replay detection does not prevent packet dropping.

## Observability controls

| Threat | Control |
|---|---|
| Plaintext log leakage | Explicit telemetry sanitization |
| Key leakage | Secret-aware logging; keys excluded from event payloads |
| Metadata leakage | Minimize sensitive labels and routing metadata |
| Log injection | Structured logging with validated fields |
| Log/telemetry DoS | Bounded queues, sizes and rates |
| Monitoring outage | Observability decoupled from message path |
| False security evidence | Treat protocol/test evidence as authoritative |
| Dashboard exposure | Restrict observability interfaces in deployment |

The dashboard is a visualization of evidence; it is not itself a
security control.

See [07-Observability-and-Security-Telemetry](../02-architecture/07-Observability-and-Security-Telemetry.md).

## M2 Implemented Security Controls

The following identity controls are implemented and tested:

1. Ed25519 identity generation uses `crypto/rand`.
2. Public identity is constrained to 32 bytes.
3. Loaded private-key input is defensively copied before derivation.
4. Stored public-key material is independently owned.
5. Returned public keys are defensive copies.
6. Malformed key/signature inputs are rejected without panics.
7. Private key material is not exposed through the identity API.
8. Ed25519 signing and verification use the standard library.

These controls do not imply session confidentiality or forward secrecy; those
depend on the later X25519/session protocol.


## M3 Implemented Security Controls

The following session-layer controls are implemented and verified:

1. Fresh ephemeral X25519 keypairs are generated for new sessions.
2. X25519 is authenticated by Ed25519 signatures over the canonical handshake
   transcripts.
3. The full responder transcript hash binds the session identifier and HKDF
   salt.
4. HKDF role-specific labels derive independent directional session keys.
5. ChaCha20-Poly1305 authenticates application ciphertext and canonical AAD.
6. The deterministic 96-bit nonce is derived from the per-direction
   sequence number and is never intentionally reused under a session key.
7. Direction is derived from the established session role rather than
   caller-controlled input.
8. Sequence numbers begin at zero and never wrap; exhaustion returns an
   explicit error.
9. A 64-message receive replay window rejects duplicates and packets that are
   too old.
10. Replay-window state is committed only after successful AEAD
    authentication.
11. Send and receive session state use independent mutexes for concurrent use.
12. Session keys and sequence/replay state are memory-only and fresh
    sessions are required after restart.

These controls establish the M3 cryptographic/session boundary. They do not
establish secure direct networking, mesh routing security, anonymity,
metadata hiding or endpoint compromise resistance.


## M4 Implemented Security Controls

The following direct-messaging controls are implemented and verified:

1. TCP frames use a 4-byte big-endian length prefix.
2. Frame length is bounded to 64 KiB before allocation.
3. `io.ReadFull` handles complete frame reads, including partial TCP reads.
4. Complete writes are serialized and short writes are handled explicitly.
5. Unknown Protobuf fields are rejected at the top-level and populated payload
   level.
6. Packet type and Protobuf `oneof` payload are required to agree.
7. INIT, RESP and APP_DATA fixed-width fields are validated before crypto.
8. APP_DATA is rejected before an M3 session is established.
9. APP_DATA is bound to the active session ID and expected source/destination
   identities.
10. Application plaintext is delivered only after successful AEAD
    authentication.
11. Direct M4 channels use one established session per TCP connection.
12. The transport layer remains independent of cryptographic implementation
    details.

M4 evidence demonstrates direct secure messaging and transport-path tamper
rejection. It does not establish mesh routing security or anonymity.


## M5 Implemented Security and Resilience Controls

M5 adds runtime network controls without changing the M3/M4 cryptographic
construction:

1. Peer registration occurs only after authenticated direct-channel setup.
2. Peer identity is canonicalized from the authenticated Ed25519 identity.
3. Pending handshakes are bounded and rejected without blocking the listener.
4. Active peer count is bounded.
5. Concurrent outbound dials are bounded and their slots are released on all
   terminal outcomes.
6. Handshake and dial deadlines limit incomplete connection resource use.
7. Initial oversized frames are rejected before large allocation.
8. Duplicate direct connections are deterministically resolved using the
   lexicographic identity rule.
9. Registry replacement is synchronized and stale cleanup cannot delete the
   replacement.
10. Terminal peer states cannot be resurrected.
11. Telemetry uses a bounded non-blocking queue and is outside the E2EE path.
12. The telemetry worker is outside the manager WaitGroup, so an indefinitely
    blocking external recorder cannot prevent manager shutdown.
13. Telemetry records no plaintext, private keys, session keys or passwords.

M5 evidence establishes direct-network runtime controls; it does not establish
resistance to all routing attacks or network-wide metadata analysis.

## M6 Mesh Security Controls

- Strict Protobuf and frame validation before routing.
- Exactly 16-byte PacketIDs.
- Duplicate suppression keyed by `(source_node, packet_id)`.
- 10,000-entry / 2-minute bounded duplicate cache.
- Maximum accepted TTL of 32; default initial TTL 16.
- Incoming-peer exclusion during flooding.
- Per-peer queues capped at 1,000 packets / 2 MiB.
- Non-blocking router enqueue path.
- Endpoint session state independent from relay transport state.
- Exact M3 transcript correlation for multi-hop RESP messages.
- No relay decryption of endpoint APP_DATA.


## M7 Security Controls

The following hardening controls were implemented and verified:

1. Router input is structurally validated before cache insertion or routing.
2. Raw router input above the 64 KiB boundary is rejected at the router entry
   boundary.
3. Concurrent replay verification, AEAD authentication and replay-state
   commitment are serialized under the receive mutex.
4. Replay state is committed only after successful AEAD authentication.
5. Authentication and identity-binding regressions remain passing after
   hardening.
6. Handshake deadlines bound stalled authentication attempts.
7. Pending-handshake slots are released on timeout/failure through deferred
   cleanup.
8. Telemetry no longer propagates arbitrary remote peer identities through the
   runtime recorder interface.
9. The current repository has no active Prometheus metric exporter or
   metric-vector implementation.
10. Future concrete metrics must use fixed, bounded label vocabularies and
    must not use attacker-controlled peer identities as labels.
11. Telemetry remains out-of-band and must not become a dependency of routing,
    encryption, authentication or delivery.
12. Full repository tests, race detection, vetting and build validation pass.

M7 does not add anonymity, metadata hiding or global network-wide DoS
protection.


## M8.3 inbound admission control

Transport admission no longer requires out-of-band pre-authorization. Unknown
inbound peers may consume a bounded pending-handshake slot, subject to the
existing handshake timeout, frame-size validation and connection/peer limits.
Successful peer registration still requires cryptographic authentication.

The change increases exposure to unauthenticated handshake attempts compared
with the former whitelist model. It is not a claim of complete DoS prevention;
the documented controls bound concurrent pending work and incomplete connection
lifetime.
