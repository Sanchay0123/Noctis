# Key Derivation

## HKDF-SHA-256

HKDF should convert authenticated shared secret material into
purpose-specific session keys.

## Conceptual schedule

``` text
X25519 shared secret
        |
        v
HKDF Extract
        |
        v
PRK
        |
        v
HKDF Expand + context
        |
        +--> Alice→Bob encryption key
        |
        +--> Bob→Alice encryption key
```

## Domain separation

The derivation context should distinguish:

-   protocol
-   protocol version
-   session
-   direction
-   purpose

This prevents accidental reuse of the same derived bytes for unrelated
functions.

## Requirements

-   no ad-hoc hash concatenation
-   no custom KDF
-   explicit context
-   deterministic test vectors for the implemented schedule
-   different sessions should not accidentally reuse identical key
    material
