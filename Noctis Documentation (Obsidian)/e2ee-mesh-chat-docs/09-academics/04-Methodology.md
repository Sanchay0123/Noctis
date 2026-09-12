# Methodology

## Phase 1 --- Requirements and threat modeling

Define:

-   system requirements
-   assets
-   attackers
-   trust boundaries
-   security goals
-   non-goals

See [[03-security/01-Threat-Model]].

## Phase 2 --- Architecture

Define:

-   layers
-   components
-   data flows
-   protocol boundaries
-   cryptographic architecture
-   routing model

See [[02-architecture/01-System-Architecture]].

## Phase 3 --- Cryptographic implementation

Implement only established primitives through mature libraries.

**Current evidence boundary:** the Ed25519 long-term identity primitive is
implemented and audited. X25519, HKDF, AEAD, and the authenticated session
protocol remain future implementation stages.

Validate:

-   identity
-   key agreement
-   authentication
-   KDF
-   AEAD
-   replay controls

## Phase 4 --- Direct networking

Prove secure two-node communication before introducing routing
complexity.

## Phase 5 --- Mesh

Add:

-   peer discovery
-   routing
-   forwarding
-   TTL
-   duplicate detection
-   route failure behavior

## Phase 6 --- Security evaluation

Attack the system intentionally:

-   modify packets
-   replay packets
-   substitute identities
-   send malformed packets
-   attempt relay decryption

## Phase 7 --- Demonstration and analysis

Compare:

-   intended security properties
-   observed behavior
-   test evidence
-   limitations

## Phase 8 --- Documentation

Record implementation decisions, evidence and limitations without
upgrading unverified claims into facts.
