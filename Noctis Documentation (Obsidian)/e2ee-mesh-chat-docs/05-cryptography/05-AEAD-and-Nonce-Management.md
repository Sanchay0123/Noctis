# AEAD and Nonce Management

## Selected AEAD

**ChaCha20-Poly1305** is the preferred single AEAD construction.

## Security properties

AEAD provides:

-   confidentiality
-   ciphertext integrity
-   authentication of associated data

## Nonce invariant

> **Never reuse a nonce with the same key.**

This must be enforced by architecture, not by developer memory.

## Candidate strategy

A session may use:

``` text
session key
+
monotonically allocated message sequence
+
defined nonce construction
```

The exact nonce construction must be reviewed before implementation.

## Never do

``` text
random nonce with accidental collision assumptions
fixed nonce
counter reset under same key
reuse nonce after reconnect with same key
```

unless the final protocol provides a rigorous construction that makes
the specific behavior safe.

## AAD

AAD should authenticate metadata whose modification must invalidate the
ciphertext.

Candidate fields:

-   version
-   message type
-   sender identity
-   recipient identity
-   session identifier
-   sequence number

Final selection requires protocol review.
