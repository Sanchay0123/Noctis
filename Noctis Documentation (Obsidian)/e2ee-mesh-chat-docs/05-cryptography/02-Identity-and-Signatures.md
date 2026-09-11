# Identity and Signatures

## Ed25519 role

Ed25519 provides:

-   long-term node identity
-   digital signatures
-   authentication of identity-controlled protocol data

It does **not** provide:

-   encryption
-   key agreement
-   anonymity

## Identity tests

Required:

-   generate identity
-   serialize identity
-   reload identity
-   sign message
-   verify valid signature
-   reject modified message
-   reject invalid signature
-   reject malformed public key

## Key storage

Private identity keys must be treated as high-sensitivity local secrets.

The implementation must document:

-   file location
-   permissions
-   serialization format
-   whether encrypted at rest
-   backup implications
-   rotation/replacement behavior

No private key should appear in ordinary logs.
