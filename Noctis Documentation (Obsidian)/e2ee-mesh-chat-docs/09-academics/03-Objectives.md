# Objectives

## Primary objective

Design and implement a working prototype of an end-to-end encrypted
messaging system over a decentralized mesh network.

## Specific objectives

1.  Design a cryptographic identity system using Ed25519.
2.  Establish authenticated shared session material using X25519.
3.  Derive session keys using HKDF-SHA-256.
4.  Protect messages using ChaCha20-Poly1305.
5.  Prevent AEAD nonce reuse.
6.  Implement replay protection.
7.  Separate mesh routing from application cryptography.
8.  Implement multi-hop forwarding.
9.  Implement TTL and duplicate detection.
10. Test packet tampering and impersonation resistance.
11. Handle basic peer disconnection/reconnection.
12. Build a usable interface.
13. Create a reproducible security demonstration.
14. Document security limitations honestly.
