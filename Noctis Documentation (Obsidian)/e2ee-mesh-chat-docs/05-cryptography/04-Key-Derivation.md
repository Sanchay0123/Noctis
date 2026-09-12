# Key Derivation

## HKDF-SHA-256

**Implemented / verified in M3.2.**

HKDF-SHA-256 converts authenticated X25519 shared secret material into
purpose-specific directional session keys. The construction is frozen;
the implementation is not yet present.

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

## Frozen M3 schedule

The approved schedule is:

```text
salt = SHA256(T_RESP)
PRK  = HKDF-Extract(salt, SS)

K_A_to_B = HKDF-Expand(
    PRK,
    "MeshChat-v1|session|initiator->responder",
    32
)

K_B_to_A = HKDF-Expand(
    PRK,
    "MeshChat-v1|session|responder->initiator",
    32
)
```

The transcript hash therefore binds the session context, while the
directional labels provide key separation. The `info` labels are encoded
as ASCII bytes.

This schedule must be implemented exactly as specified unless a future
technical-lead review explicitly changes the protocol.

## Requirements

-   no ad-hoc hash concatenation
-   no custom KDF
-   explicit context
-   deterministic test vectors for the implemented schedule
-   different sessions should not accidentally reuse identical key
    material


## M3 implementation evidence

The session implementation derives both directional keys from the same
HKDF-Extract output but uses distinct role-specific `info` labels:

```text
K_A_to_B = HKDF-Expand(PRK, "MeshChat-v1|session|initiator->responder", 32)
K_B_to_A = HKDF-Expand(PRK, "MeshChat-v1|session|responder->initiator", 32)
```

The full SHA-256 digest of `T_RESP` is both the session identifier and HKDF
salt. The implementation was checked against a deterministic known-answer
vector. Session derivation tests also cover transcript sensitivity and
independent handshakes.
