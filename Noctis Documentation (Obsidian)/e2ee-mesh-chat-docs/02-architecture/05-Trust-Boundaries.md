# Trust Boundaries

## Boundary 1 --- Endpoint vs relay

Alice and Bob are endpoints for the application message.

Relays are not trusted with application plaintext.

## Boundary 2 --- Application vs protocol

Application content enters the cryptographic layer before being placed
into a routable packet.

## Boundary 3 --- Protocol vs transport

Transport can deliver bytes but must not interpret application
cryptographic meaning.

## Boundary 4 --- Local storage vs runtime

Persistent private identity material is more sensitive than ordinary
configuration.

## Boundary diagram

```mermaid
flowchart TB
    subgraph A[Endpoint A Trust Domain]
        A1[Identity]
        A2[Session Keys]
        A3[Plaintext]
    end

    subgraph M[Untrusted Relay Domain]
        M1[Routing Metadata]
        M2[Ciphertext]
    end

    subgraph B[Endpoint B Trust Domain]
        B1[Identity]
        B2[Session Keys]
        B3[Plaintext]
    end

    A3 --> A2
    A2 --> M2
    M2 --> B2
    B2 --> B3
```

## Review rule

Every new component must state:

1.  what it trusts
2.  what trusts it
3.  what secrets it can access
4.  what attacker-controlled input it accepts


## M4 direct transport boundary

The TCP connection is an attacker-controlled byte boundary. M4 therefore
performs bounded framing and structural packet validation before session-layer
processing. The transport does not receive plaintext or session secrets.

`DirectChannel` is the integration boundary between the application-facing
message API and the protocol/crypto layers. It must not expose plaintext from
a failed authentication or decryption attempt.

## M6 Relay Trust Boundary

The relay trust boundary is intentionally weak: a relay is trusted to forward packets according to protocol rules, but it is **not** trusted with application confidentiality.

The relay may observe routing metadata, packet type, PacketID, TTL and ciphertext. It must not obtain endpoint session keys or plaintext. A malicious relay can drop, delay, reorder, duplicate or selectively forward packets; M6 does not claim to prevent those availability attacks.

## M8 application/UI trust boundary

The UI is an unprivileged presentation client of `internal/app`.

```mermaid
flowchart LR
    UI[GUI] -->|display-safe API + app events| APP[Application Service]
    APP -->|session lookup / plaintext only at endpoint| S[Session Manager]
    APP -->|opaque routed APP_DATA| R[Routing]
    R -->|ciphertext only| M[Mesh / Relays]
```

The boundary forbids GUI access to:

- Ed25519 private keys
- X25519 ephemeral private keys
- session keys
- AEAD nonces/replay internals
- raw sockets
- raw protobuf packets
- routing internals

Routing likewise has no application-plaintext decryption path. Only the
endpoint application/session boundary may convert authenticated ciphertext
into plaintext.
