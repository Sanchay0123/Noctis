# Data Flow

## Outbound message

```mermaid
flowchart LR
    P[Plaintext] --> APP[Chat Application]
    APP --> S[Session Manager]
    S --> E[AEAD Encryption]
    E --> C[Ciphertext]
    C --> W[Protocol Envelope]
    W --> R[Routing]
    R --> T[Transport]
```

## Relay processing

```mermaid
flowchart LR
    T[Transport] --> V[Packet Validation]
    V --> R[Routing Decision]
    R --> F[Forward]
    F --> T2[Next Hop]
```

A relay must not have an application-level path from the received packet
to a plaintext decryption operation.

## Destination processing

```mermaid
flowchart LR
    T[Transport] --> V[Packet Validation]
    V --> S[Session Lookup]
    S --> A[AEAD Authentication]
    A --> D[Decrypt]
    D --> APP[Chat Application]
```

## Sensitive-data rule

Sensitive material must have a documented lifetime.

Examples:

-   private identity key: persistent secret
-   session key: bounded session lifetime
-   plaintext message: application-controlled lifetime
-   nonce: protocol field; never treated as secret
-   ciphertext: transport-visible but confidential

See [[05-cryptography/06-Key-Lifecycle]].
