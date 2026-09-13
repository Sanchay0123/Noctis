# Protocol Error Handling

## Principles

1.  Treat remote input as hostile.
2.  Validate before use.
3.  Fail closed for authentication failures.
4.  Do not reveal cryptographic secrets in errors.
5.  Bound error-triggered resource usage.
6.  Do not crash the node because a peer sent malformed data.

## Error classes

### Structural

Malformed serialization, missing fields, invalid lengths.

### Protocol

Unsupported version, invalid message type, illegal state.

### Cryptographic

Invalid signature, invalid key, failed AEAD authentication.

### Routing

TTL expired, no route, duplicate packet, loop prevention.

### Transport

Timeout, connection reset, unreachable peer.

## Logging policy

Safe:

``` text
packet authentication failed
peer connection closed
route unavailable
packet rejected: invalid version
```

Unsafe:

``` text
session key = ...
private key = ...
plaintext = ...
password = ...
```

## Error oracle caution

Avoid detailed errors that allow an attacker to distinguish sensitive
internal states unless that distinction is intentionally required.


## M7 Failure-Handling Controls

M7 verifies that hostile input and failure paths remain fail-closed and
resource-bounded:

- router parsing errors terminate processing without routing;
- structural validation failures terminate processing before cache insertion;
- oversized router input is rejected before normal parsing/routing;
- AEAD authentication failures do not commit replay state;
- concurrent replay attempts are serialized;
- stalled handshakes terminate on deadline;
- pending handshake resources are released through deferred cleanup;
- telemetry is out-of-band and must not be required for core message delivery.

No production panic vulnerability requiring remediation was identified in the
final M7 audit.
