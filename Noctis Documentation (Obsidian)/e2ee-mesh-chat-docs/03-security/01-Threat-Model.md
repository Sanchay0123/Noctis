# Threat Model

## Method

This threat model defines:

-   assets
-   trust assumptions
-   attacker capabilities
-   security goals
-   security boundaries
-   limitations

It is intentionally scoped to an academic prototype.

## Assets

### A1 --- Private identity keys

Long-term Ed25519 private keys.

### A2 --- Session key material

Keys used to protect E2EE application messages.

### A3 --- Message plaintext

Application content.

### A4 --- Message authenticity

The ability to distinguish authentic messages from modified/injected
ones.

### A5 --- Replay state

Receiver state used to prevent duplicate acceptance.

### A6 --- Routing state

Peer and route information.

## Attacker capabilities

An attacker may:

-   observe traffic
-   capture packets
-   modify packets
-   replay packets
-   drop packets
-   inject malformed packets
-   operate a malicious relay
-   attempt identity impersonation
-   attempt key-substitution attacks
-   attempt basic resource exhaustion
-   attempt to infer communication metadata

## Security goals

The system should provide:

-   confidentiality of application content
-   integrity/authenticity of messages
-   authentication of communicating identities
-   replay resistance
-   basic impersonation resistance
-   relay confidentiality

## Attacker examples

### Passive relay

Can observe:

-   packet timing
-   packet size
-   routing metadata
-   ciphertext

Should not obtain:

-   application plaintext
-   endpoint session keys

### Malicious relay

Can:

-   drop packets
-   delay packets
-   duplicate packets
-   modify packets
-   inject malformed traffic

Should not be able to:

-   silently modify authenticated application content
-   decrypt correctly protected E2EE messages

### Identity attacker

Attempts to replace a legitimate peer's cryptographic identity.

Expected defense:

-   authenticated key establishment
-   identity binding
-   transcript authentication

## Out of scope

Not guaranteed:

-   endpoint compromise resistance
-   global traffic analysis resistance
-   perfect anonymity
-   perfect metadata hiding
-   universal routing-attack resistance
-   production-grade availability

See [[03-security/02-Security-Goals]] and
[[03-security/05-Security-Limitations]].
