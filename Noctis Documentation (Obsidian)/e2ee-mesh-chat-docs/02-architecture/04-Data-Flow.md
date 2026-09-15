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


## M4 direct secure messaging flow

```mermaid
flowchart LR
    P[Plaintext] --> E[M3 Session EncryptMessage]
    E --> A[AppDataPayload]
    A --> M[MeshPacket]
    M --> F[4-byte length prefix + Protobuf]
    F --> TCP[TCP Connection]
    TCP --> UF[Frame + packet validation]
    UF --> S[M3 Session binding]
    S --> D[DecryptMessage / AEAD authentication]
    D --> Q[Authenticated plaintext]
```

For M4 the path is direct and no forwarding occurs. The transport can observe
outer packet metadata and ciphertext, but not application plaintext or session
keys.

### M4 handshake flow

```text
Initiator: GenerateInit -> send INIT -> receive RESP -> ProcessResp -> Session
Responder: receive INIT -> ProcessInit -> GenerateResp -> send RESP -> Session
```

APP_DATA is rejected until the session is established.

## M6 Multi-Hop Data Flow

```text
A encrypts APP_DATA for C
        ↓
A creates MeshPacket
        ↓
B validates envelope + PacketID + TTL
        ↓
B decrements TTL and queues packet
        ↓
C receives the same endpoint ciphertext
        ↓
C authenticates/decrypts with the A-C E2EE session
```

For handshake traffic, INIT and RESP are similarly forwarded through relays, while the endpoint transcript and X25519-derived session remain between the actual initiator and responder.


## M8 Message Flow

### Outbound

```text
Fyne GUI
  ↓ SendMessage(peerID, plaintext)
ApplicationService
  ↓ session lookup
Session / AEAD
  ↓ ciphertext
Mesh routing
  ↓
Relay peers
  ↓
Destination session
```

### Inbound

```text
Relay / transport
  ↓ opaque APP_DATA
Routing
  ↓ session_id + sequence + ciphertext
ApplicationService
  ↓ Lookup(session_id)
Session.DecryptMessage()
  ↓ plaintext
SubscribeEvents()
  ↓
Fyne GUI
```

Routing does not decrypt application payloads.
