# Security Limitations

## Metadata exposure

Even with E2EE, a relay may observe:

-   who connects to whom at the routing layer
-   timing
-   packet sizes
-   frequency
-   route information
-   node identifiers

Therefore the system is **not anonymous** by default.

## Malicious routing

A relay may still:

-   drop packets
-   delay packets
-   refuse forwarding
-   provide incorrect routing information
-   attempt resource exhaustion

Cryptographic message integrity does not guarantee delivery.

## Endpoint compromise

If Alice's endpoint is compromised, an attacker may access plaintext
before encryption or after decryption.

If Bob's endpoint is compromised, the attacker may access plaintext
after decryption.

E2EE does not protect against a compromised endpoint.

## Local key storage

The prototype may use local persistent key storage. The exact protection
mechanism must be documented once implemented.

If private keys are stored unencrypted on disk, that must be disclosed
as a limitation.

## Forward secrecy

Do not claim production-grade forward secrecy unless the final session
protocol and key lifecycle actually provide and demonstrate it.

## Deniability

The prototype does not automatically provide deniability merely because
it uses modern cryptography.

## Availability

Encryption does not prevent a malicious relay from dropping messages.

## Security-language policy

Forbidden claims unless independently demonstrated:

-   "unhackable"
-   "perfectly secure"
-   "anonymous"
-   "military-grade"
-   "completely secure"

See [[03-security/06-Security-Review-Checklist]].
