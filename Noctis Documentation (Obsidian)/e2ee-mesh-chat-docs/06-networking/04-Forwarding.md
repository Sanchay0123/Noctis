# Forwarding

## Forwarding pipeline

```mermaid
flowchart TD
    IN[Incoming packet] --> SIZE[Size validation]
    SIZE --> STRUCT[Structure validation]
    STRUCT --> DUP{Duplicate?}
    DUP -->|Yes| DROP[Drop]
    DUP -->|No| TTL{TTL valid?}
    TTL -->|No| DROP
    TTL -->|Yes| ROUTE[Route lookup]
    ROUTE --> NEXT[Select next hop]
    NEXT --> DEC[Decrement TTL / update routing field]
    DEC --> OUT[Forward unchanged E2EE payload]
```

## Critical property

The forwarding engine should not modify cryptographic ciphertext or
authenticated E2EE fields except where the protocol explicitly defines a
routing-only field that is safe to update.

## Duplicate detection

The duplicate cache must be bounded to prevent an attacker from causing
unbounded memory growth.

## Forwarding failures

Possible outcomes:

-   no route
-   TTL expired
-   connection unavailable
-   duplicate
-   malformed packet
