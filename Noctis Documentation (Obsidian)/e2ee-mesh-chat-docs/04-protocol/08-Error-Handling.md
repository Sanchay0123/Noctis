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
