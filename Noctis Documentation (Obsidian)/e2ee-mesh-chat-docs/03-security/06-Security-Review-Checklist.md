# Security Review Checklist

## Cryptography

-   [ ] Established library primitives only
-   [ ] Ed25519 identity handling reviewed
-   [ ] X25519 key exchange reviewed
-   [ ] Identity binding reviewed
-   [ ] Transcript binding reviewed
-   [ ] HKDF context separation reviewed
-   [ ] AEAD usage reviewed
-   [ ] Nonce generation reviewed
-   [ ] Nonce reuse impossible by design
-   [ ] Key lifecycle reviewed
-   [ ] Invalid cryptographic inputs fail safely

## Protocol

-   [ ] Version checked
-   [ ] Message types validated
-   [ ] Field lengths bounded
-   [ ] Unknown fields handled intentionally
-   [ ] Authentication occurs before trust-sensitive processing
-   [ ] Replay protection enforced
-   [ ] Duplicate detection enforced
-   [ ] Error messages do not leak secrets

## Networking

-   [ ] TTL enforced
-   [ ] Routing loops bounded
-   [ ] Duplicate packets bounded
-   [ ] Connection limits exist
-   [ ] Timeouts exist
-   [ ] Malformed packets cannot crash the node
-   [ ] Resource use is bounded

## Privacy

-   [ ] Plaintext absent from relay logs
-   [ ] Private keys absent from logs
-   [ ] Session keys absent from logs
-   [ ] Sensitive metadata documented
-   [ ] Debug mode reviewed

## Evidence

-   [ ] Security tests exist
-   [ ] Negative tests exist
-   [ ] Integration tests exist
-   [ ] Demonstration evidence captured
-   [ ] Documentation claims match actual evidence
