# Assumptions and Constraints

## Assumptions

The initial prototype assumes:

1.  Endpoints are not fully compromised.
2.  The underlying operating system and cryptographic library are
    trusted.
3.  Cryptographic randomness is provided by the platform/library.
4.  Nodes can communicate over an IP-capable network.
5.  A node can maintain local persistent state where required.
6.  The project is primarily intended for controlled academic
    demonstrations.

## Constraints

-   Do not implement cryptography from scratch.
-   Do not introduce unnecessary infrastructure.
-   Do not implement blockchain or consensus unless a concrete
    requirement appears.
-   Do not add onion routing without a justified threat-model
    requirement.
-   Keep routing and E2EE logically independent.
-   Keep the implementation explainable during a university viva.

## Security boundary

The system is not expected to defend against:

-   a compromised endpoint
-   a compromised host OS
-   theft of an unprotected local private-key file
-   a global passive adversary
-   all metadata analysis

Any stronger claim must be supported by new architecture and evidence.
