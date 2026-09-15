# M7 Acceptance Evidence

## Status

**🟢 COMPLETE / VERIFIED / APPROVED**

This document records the final M7 security-hardening evidence.

## 1. Concurrent replay

### Objective

Ensure concurrent delivery of the same authenticated application ciphertext
cannot result in multiple successful decryptions or replay-window races.

### Control

`receiveMu` covers:

```text
checkReplayWindow(sequence)
        ↓
AEAD Open
        ↓
commitReplayWindow(sequence)
```

Replay state is committed only after successful authentication.

### Test

`TestAEADConcurrentReplay`

The test performs ten concurrent decryptions of the same ciphertext and
sequence number and asserts:

- exactly 1 success
- exactly 9 `ErrReplayDetected`

**Result: PASS**

## 2. Router validation boundary

### Objective

Ensure the routing layer validates hostile input before caching or forwarding.

### Production boundary

```text
raw bytes
   ↓
proto.Unmarshal
   ↓
transport.ValidatePacket
   ↓
packet cache
   ↓
TTL / destination / forwarding
```

### Adversarial cases

`TestRouterValidationBoundary` covers:

| Case | Result |
|---|---|
| Malformed Protobuf | PASS |
| Unknown protocol version | PASS |
| Unknown Protobuf fields | PASS |
| Packet type / oneof mismatch | PASS |
| Invalid source length | PASS |
| Invalid APP_DATA session ID length | PASS |
| Invalid APP_DATA ciphertext length | PASS |
| Oversized packet | PASS |

The oversized case uses a 70,000-byte ciphertext and is rejected because the
raw router message exceeds the 64 KiB boundary.

## 3. Telemetry security and cardinality

### Objective

Prevent attacker-controlled peer identities from becoming an unbounded
telemetry resource.

### Final implementation state

The current repository has telemetry hooks and a no-op implementation but no
active Prometheus metric exporter.

The final `TelemetryRecorder` interface does not accept remote peer identity
values:

```text
RecordPeerConnected(isOutbound bool)
RecordPeerDisconnected()
RecordHandshakeFailed(err error)
RecordDuplicateConnection()
RecordResourceLimitReached(limitName string)
```

Repository searches for common Prometheus metric-vector construction and label
APIs found no active implementation.

Therefore arbitrary peer identities cannot currently become Prometheus metric
labels through the runtime telemetry interface.

### Security requirement for future M8 telemetry

Any future concrete metrics exporter must use fixed, bounded label vocabularies.
Peer identities must not be introduced as arbitrary metric labels.

Telemetry remains out-of-band and must not block, authenticate, decrypt, route
or deliver application messages.

## 4. Authentication regression

The final regression set passed:

- `TestDirectChannel_Handshake`
- `TestDirectChannel_InvalidIdentity`
- `TestDirectChannel_ReplayAttack`
- `TestHandshakeFailed`

## 5. Handshake timeout and resource release

`TestHandshakeTimeout` uses a shortened test deadline and demonstrates that:

1. a stalled handshake occupies the pending slot;
2. the handshake times out;
3. the pending slot is released;
4. a subsequent connection can acquire the slot.

`TestMaxPendingHandshakesExhaustion` confirms the configured pending-handshake
limit remains enforced.

Both tests pass.

## 6. Panic / failure audit

The final audit covered:

- `internal/protocol`
- `internal/transport`
- `internal/crypto`
- `internal/routing`
- `internal/mesh`

The review covered malformed parsing, bounded allocation, connection closure,
AEAD/replay failures, cache/queue bounds, goroutine cleanup, semaphore release,
WaitGroup ownership, and telemetry failure isolation.

No production panic path requiring remediation was identified.

Test-only panic use in failure-injection mocks is not part of the production
runtime.

## 7. Full validation

The final repository validation passed:

```text
go vet ./...
go test ./...
go test -race ./...
go build ./...
```

The full test suite includes the primary runtime packages and the `_noctis`
packages present in the repository.

## 8. Crypto terminology

Repository search:

```text
grep -r -i -E "AES-GCM|AES GCM|AESGCM|AES" <repository>
```

returned no matches.

The implemented AEAD remains ChaCha20-Poly1305.

## 9. Evidence standard

The M7 gate is supported by:

- focused security tests
- adversarial boundary tests
- concurrency/race testing
- source inspection
- full repository validation
- previously approved M3/M4/M6 integration evidence

## Gate decision

**🟢 M7 APPROVED**

See [05-Architecture-Review-M7](../00-governance/05-Architecture-Review-M7.md).
