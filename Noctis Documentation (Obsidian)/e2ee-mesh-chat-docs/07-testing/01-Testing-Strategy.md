# Testing Strategy

## Testing layers

```mermaid
flowchart TB
    U[Unit Tests]
    P[Protocol Tests]
    I[Integration Tests]
    A[Adversarial Tests]
    D[Demonstration Tests]

    U --> P --> I --> A --> D
```

## Unit

Test individual:

-   crypto wrappers
-   serializers
-   replay state
-   routing functions
-   validation functions

## Protocol

Test valid and invalid messages.

## Integration

Test multiple node processes.

## Adversarial

Assume packets and peers can be malicious.

## Demonstration

Reproduce the actual claims made in the final presentation.

## Security test philosophy

Every major security property must have at least one negative test.

Example:

> If we claim tampering is detected, a test must deliberately tamper
> with authenticated data and expect rejection.

## Observability tests

Observability is part of the test surface.

Test at least:

- expected metrics for packet receive/forward/drop/expiry
- expected security events for replay and authentication failures
- absence of plaintext from relay logs
- absence of private/session keys from logs
- resistance to log injection through attacker-controlled fields
- bounded telemetry behavior under event floods
- continued message delivery when monitoring components are unavailable

The test suite should distinguish between:

1. **Telemetry correctness** — did the expected event/metric appear?
2. **Telemetry safety** — did the event avoid leaking protected data?
3. **Telemetry independence** — does the mesh still work without the
   observability backend?

## M2 Security Test Coverage

M2.1 testing covers:

- Ed25519 key generation and key-size invariants
- public/private correspondence
- deterministic signing
- modified-message rejection
- wrong-identity rejection
- malformed key/signature handling
- defensive copying of returned public keys
- defensive copying during private-key loading
- nil and empty-message behavior
- race-detector validation

The tests use freshly generated runtime keys and do not commit or print private
key material.


## M3 Security Test Coverage

M3.1 covers X25519 generation/agreement, RFC7748 known-answer behavior,
invalid/low-order public keys, ownership and lifecycle boundaries.

M3.2 covers authenticated handshake success/failure, transcript tampering,
identity substitution, ephemeral substitution, role/state confusion,
session-ID derivation, transcript sensitivity, HKDF known-answer values and
malformed inputs.

M3.3 covers deterministic nonce/AAD vectors, ciphertext modification,
authentication failures, truncation, duplicate messages, outside-window
messages, valid out-of-order messages, replay-state poisoning attempts,
cross-session and cross-direction rejection, sequence zero, exhaustion and
concurrent send/receive behavior.

The reported containerized validation sequence passed `gofmt`, `go vet`,
tests, race detection and build checks.
