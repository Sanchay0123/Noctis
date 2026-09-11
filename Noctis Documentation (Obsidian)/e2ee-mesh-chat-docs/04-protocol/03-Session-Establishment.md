# Session Establishment

## Status

**Design under review — M3 implementation blocked until the canonical  
construction is approved.**

## Key roles

Each node has:

- long-term Ed25519 identity keypair
    
- fresh ephemeral X25519 keypair per session
    

The Ed25519 key is used for authentication/signatures only.

The X25519 key is used for key agreement only.

## Initiation

Alice generates:

```
eA = ephemeral X25519 private key
EA = corresponding public key
```

Alice creates an authenticated INIT structure containing at minimum:

```
protocol_version
handshake_domain
initiator_identity = ID_A
responder_identity = ID_B
initiator_ephemeral = EA
```

Alice signs the canonical encoding of that structure with her Ed25519  
private key.

## Response

Bob verifies Alice's identity signature and confirms that the requested  
responder identity is Bob.

Bob generates:

```
eB = ephemeral X25519 private key
EB = corresponding public key
```

Bob signs a canonical response transcript containing:

```
protocol_version
handshake_domain
initiator_identity = ID_A
responder_identity = ID_B
initiator_ephemeral = EA
responder_ephemeral = EB
```

Alice verifies Bob's signature.

## Shared secret

Both parties compute:

```
SS = X25519(ephemeral_private, peer_ephemeral_public)
```

The implementation must reject an invalid/low-order result according to  
the selected library's documented behavior.

## Key derivation

The final KDF must be specified before M3.

It must bind:

- the shared secret
    
- protocol/domain context
    
- authenticated handshake transcript
    
- both identities
    
- both ephemeral public keys
    
- directional role labels
    

Two independent directional AEAD keys are required.

Do not use a generic string such as `"MeshChat_v1_A_to_B"` without  
binding it to the actual session transcript.

## Session identifier

`session_id` is a protocol identifier for selecting the correct session  
state. It is not a cryptographic secret.

The implementation must define how it is generated and ensure that it  
cannot cause ambiguous session-state lookup.

A cryptographic transcript hash may be used as the basis for deriving a  
session identifier, subject to the final schema/type decision.

## Key lifecycle

After a successful handshake:

- retain only the session state required for operation
    
- erase ephemeral private keys as soon as they are no longer required
    
- never persist session keys by default
    
- establish fresh session keys after restart
    
- establish fresh session keys during rotation
    

## Handshake replay

Handshake messages must have their own duplicate/replay handling.

A captured INIT or RESP must not cause an endpoint to silently replace a  
live session with attacker-controlled state.

The implementation must define session-state transitions and acceptable  
retransmission behavior.

## Failure states

At minimum:

```
stateDiagram-v2
    [*] --> Idle
    Idle --> InitSent: send INIT
    InitSent --> Established: valid RESP
    InitSent --> Failed: invalid/timeout
    Idle --> Failed: invalid INIT
    Established --> Rekeying: rotation threshold
    Rekeying --> Established: valid new handshake
    Established --> Closed: disconnect
    Failed --> Idle
    Closed --> Idle
```

Invalid signatures, malformed messages, unexpected state transitions and  
timeouts must fail closed.

See [[04-protocol/06-State-Machines]].