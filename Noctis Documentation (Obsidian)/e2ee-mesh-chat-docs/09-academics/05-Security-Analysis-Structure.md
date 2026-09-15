# Security Analysis Structure

## 1. Threat model

Reference [01-Threat-Model](../03-security/01-Threat-Model.md).

## 2. Security goals

Reference [02-Security-Goals](../03-security/02-Security-Goals.md).

## 3. Cryptographic analysis

Discuss:

-   Ed25519
-   X25519
-   authenticated exchange
-   HKDF
-   ChaCha20-Poly1305
-   nonce management

## 4. Relay confidentiality

Demonstrate why an intermediate node cannot decrypt the application
payload.

## 5. Integrity

Show tampering rejection.

## 6. Replay

Show duplicate/replayed message rejection.

## 7. Impersonation

Show failed identity substitution.

## 8. Routing security

Discuss:

-   TTL
-   loops
-   malicious dropping
-   duplicate forwarding
-   malformed packets

## 9. Privacy

Document visible metadata.

## 10. Limitations

State what the implementation does not solve.

## Evidence rule

Every claim must point to either:

-   architecture/design rationale
-   implementation
-   test result
-   demonstration evidence

Prefer all four for major security claims.


## Implemented Analysis Source

Use [10-Security-Analysis](../10-technical-documentation/10-Security-Analysis.md) as the implementation-grounded security analysis source for the final academic report.
