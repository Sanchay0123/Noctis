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
