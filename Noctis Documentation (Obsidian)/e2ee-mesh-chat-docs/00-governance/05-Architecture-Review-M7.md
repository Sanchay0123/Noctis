# Architecture & Security Review — M7

## Milestone

**M7 — Security Hardening**

## Review authority

Project Overseer / Security Architect / Technical Lead

## Status

**🟢 COMPLETE / VERIFIED / APPROVED**

M7 was independently reviewed after implementation and final evidence remediation.

M7 does not redesign the frozen M3, M4 or M6 architecture. It hardens the
existing implementation against malformed input, concurrent replay,
authentication failures, packet tampering, malicious relay behavior, local
resource exhaustion, telemetry leakage/cardinality problems, and failure-path
issues.

## Scope

M7 reviewed:

- malformed and oversized input
- replay correctness under concurrency
- authentication and identity binding regression
- packet tampering
- malicious relay boundary
- basic resource exhaustion
- telemetry security and metric-cardinality safety
- panic/failure handling
- connection/session lifecycle
- concurrency and resource ownership

## Security decisions

### Concurrent replay

`DecryptMessage` serializes replay-window verification, AEAD authentication and
replay-state commitment under the receive mutex.

The state transition is therefore:

```text
replay check
    ↓
AEAD authentication
    ↓
replay-state commit
```

Replay state is not consumed by an unauthenticated ciphertext.

`TestAEADConcurrentReplay` asserts exactly one successful decryption and nine
`ErrReplayDetected` results when ten concurrent attempts process the same
ciphertext and sequence number.

### Router validation boundary

`Router.OnMessage` rejects malformed input before cache insertion or routing.
The boundary invokes the existing transport packet validator after Protobuf
unmarshalling.

The final boundary test covers:

- malformed Protobuf
- unknown protocol version
- unknown Protobuf fields
- packet type / oneof mismatch
- invalid source identity length
- invalid APP_DATA session identifier length
- invalid APP_DATA ciphertext length
- oversized packet input

The router rejects raw messages above the 64 KiB limit before normal
unmarshalling/routing processing.

### Telemetry boundary and cardinality

The runtime currently provides telemetry hooks and a no-op implementation;
there is no active Prometheus exporter or metric-vector implementation in the
repository at M7.

Peer identity parameters were removed from the `TelemetryRecorder` interface:

```text
RecordPeerConnected(isOutbound bool)
RecordPeerDisconnected()
RecordHandshakeFailed(err error)
RecordDuplicateConnection()
RecordResourceLimitReached(limitName string)
```

Consequently, arbitrary peer identities cannot enter telemetry through the
runtime recorder interface or become Prometheus label values in the current
implementation.

No dynamic metric names or Prometheus label-vector construction were found.

This is a current implementation property. A future concrete metrics exporter
must preserve bounded, fixed-vocabulary labels and must not reintroduce
attacker-controlled peer identities as metric labels.

### Handshake and resource lifecycle

Inbound handshakes acquire a bounded pending slot, apply a handshake deadline,
validate the initial packet before processing, and release the pending slot in
a deferred cleanup path.

`TestHandshakeTimeout` demonstrates that a stalled handshake times out and
that a subsequent connection can acquire the released slot.

### Failure and panic safety

The M7 source audit reviewed protocol parsing, transport framing/allocation,
cryptographic failure paths, routing cache/queue lifecycle, and mesh
connection/semaphore/WaitGroup cleanup.

The repository contains no production panic paths identified by the final
audit. Panic usage found in test-only failure-injection mocks is intentionally
test infrastructure.

### Frozen architecture preserved

M7 does not change:

- Ed25519 identity construction
- X25519 key agreement
- canonical M3 handshake transcript
- SHA-256 session identifier
- HKDF-SHA-256 derivation
- ChaCha20-Poly1305
- sequence-derived nonce construction
- 64-message replay window
- M4 TCP framing or packet schema
- M6 managed flooding
- M6 PacketID cache semantics
- M6 TTL semantics
- M6 bounded forwarding queues
- endpoint-to-endpoint E2EE relay boundary

## Validation evidence

Final repository validation passed:

```text
go vet ./...
go test ./...
go test -race ./...
go build ./...
```

Focused M7 evidence passed:

```text
TestAEADConcurrentReplay
TestRouterValidationBoundary
TestRouterValidationBoundary/oversized_packet
TestMaxPendingHandshakesExhaustion
TestHandshakeTimeout
TestDirectChannel_Handshake
TestDirectChannel_InvalidIdentity
TestDirectChannel_ReplayAttack
TestHandshakeFailed
```

The complete test and race suites also pass for the primary runtime packages
and the `_noctis` packages present in the repository.

## Crypto terminology audit

Repository search for:

```text
AES-GCM
AES GCM
AESGCM
AES
```

returned no matches. The frozen AEAD construction remains
ChaCha20-Poly1305.

## Docker

No Docker architecture or runtime behavior changed in M7. Existing M6
containerized multi-hop evidence remains authoritative.

## Remaining limitations

M7 does not provide:

- anonymity
- complete metadata hiding
- global traffic-analysis resistance
- global Sybil/flood resistance
- protection against compromised endpoints
- guaranteed message delivery

M7 provides local resource bounds and security hardening, not a complete
availability solution against an arbitrarily large adversarial network.

## Gate decision

**🟢 M7 APPROVED**

M7 is formally closed. M8 may proceed through its own architecture and
implementation gates.
