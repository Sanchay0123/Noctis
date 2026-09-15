# M0 Architecture Review — Antigravity Proposal

## Review status

**🟡 APPROVED WITH CHANGES — implementation must not begin until the  
mandatory changes below are incorporated into the implementation plan.**

Antigravity's proposal establishes the correct major separation:

- Application/UI
    
- Cryptographic/session layer
    
- Protocol/serialization
    
- Mesh/routing
    
- Transport
    
- supporting observability
    

The choice of Go + Protobuf + ChaCha20-Poly1305 is compatible with the  
project direction. Managed flooding is acceptable for the educational  
prototype, provided its resource and forwarding behavior is bounded.

The handshake concept is promising but requires protocol hardening  
before implementation.

## Decisions

### 1. Language — APPROVED

**Go is acceptable.**

Rust is not a project requirement in the governing specification.  
Go's memory safety, concurrency model and mature cryptographic APIs are  
appropriate for this prototype.

The implementation must use maintained, well-established cryptographic  
libraries and must not implement primitives from scratch.

### 2. Serialization — APPROVED

**Protocol Buffers is acceptable and preferred for the initial design.**

The project does not require human-readable wire traffic. Protobuf gives  
the protocol an explicit schema and makes malformed-input testing  
straightforward.

The `.proto` schema becomes part of the protocol specification and must  
be versioned with the implementation.

### 3. Routing — APPROVED WITH CHANGES

**Managed flooding is acceptable for the educational prototype.**

It must include:

- bounded PacketID/seen-cache
    
- duplicate detection before forwarding
    
- strict TTL/hop-limit handling
    
- forwarding fan-out control
    
- malformed-packet rejection
    
- resource limits
    
- deterministic behavior suitable for tests
    

A 60-second cache must not be treated as inherently sufficient. The  
cache lifetime and size must be justified against the maximum packet  
lifetime and expected network behavior.

### 4. Handshake — BLOCKED PENDING REFINEMENT

The signed ephemeral-DH design is acceptable as a project-specific  
protocol, but the current proposal is not yet sufficiently specified to  
implement safely.

Before M3 implementation, the following must be resolved:

1. Define the exact canonical transcript that both parties sign.
    
2. Signatures must cover both identities and both ephemeral public keys.
    
3. Include protocol version and a domain-separation label.
    
4. Define exact INIT and RESP wire structures.
    
5. Define how the recipient obtains/authenticates the Ed25519 public key.
    
6. Define session ID derivation/uniqueness.
    
7. Define key erasure after session establishment/rotation.
    
8. Define handshake replay/duplicate handling.
    
9. Define failure behavior and state transitions.
    
10. Explicitly state whether a standardized Noise construction is rejected  
    or deferred and why.
    

The current text's use of `"INIT"` and `"RESP"` is not enough by itself to  
define a cryptographically canonical transcript.

## Critical cryptographic corrections

### Nonce derivation

The proposed 64-bit sequence number mapped into a 96-bit nonce is  
acceptable **only if the following invariant is enforced**:

> A `(key, nonce)` pair must never be reused.

Therefore:

- each direction has a distinct key
    
- every session receives fresh session keys
    
- sequence numbers are monotonic per direction
    
- sequence state cannot reset while reusing a session key
    
- session keys must not survive a restart unless sequence state is  
    persisted safely and the design explicitly guarantees uniqueness
    

The safer prototype choice is to generate fresh ephemeral X25519  
session keys for every new session and discard session state on restart.

### HKDF

Do not use an all-zero salt or a generic constant casually.

The final KDF should bind the protocol context and handshake transcript.  
At minimum, the design should derive keys from the X25519 shared secret  
with a protocol/domain context and role-specific labels.

The exact construction must be documented before implementation.

### Forward secrecy

The proposal may claim forward secrecy only with the required  
conditions:

- ephemeral X25519 keys are genuinely fresh per session
    
- ephemeral private keys are erased after use
    
- long-term identity keys are not used as X25519 private keys
    
- session keys are rotated/replaced as specified
    
- the limitations of no post-compromise security are stated
    

Do not describe this as a Signal-style ratchet or claim post-compromise  
security.

## Packet-format corrections

The proposed packet format is a useful starting point but is **not yet  
the final wire specification**.

Required before implementation:

- explicit protocol version
    
- explicit packet type
    
- canonical field types for node IDs
    
- exact session identifier semantics
    
- exact nonce representation
    
- exact AAD construction
    
- handshake packet schemas
    
- size limits
    
- unknown-field/version behavior
    
- malformed packet behavior
    

Avoid storing Ed25519 public keys as base64 strings inside the binary  
Protobuf packet unless there is a concrete reason. Fixed-width `bytes`  
fields are preferable for cryptographic identifiers.

## Routing corrections

Managed flooding must not become uncontrolled broadcast.

Required:

- bounded seen-cache
    
- bounded packet size
    
- bounded peer fan-out where appropriate
    
- per-peer/per-node rate limits if needed
    
- duplicate suppression before forwarding
    
- TTL checked before forwarding
    
- malformed packets rejected without expensive work
    
- cache insertion strategy documented
    
- no plaintext inspection by relays
    

A malicious node can still consume network resources by injecting unique  
PacketIDs. This remains a documented DoS limitation unless rate limiting  
or admission controls are added.

## Discovery recommendation

For the first implementation, use **explicit bootstrap peers**.

UDP broadcast discovery should be treated as a later optional feature  
because it adds platform, container-network and spoofing complexity.

If discovery is added later, its security properties must be documented  
separately.

## Transport

TCP is approved for the first implementation.

The design must account for:

- message framing
    
- connection lifecycle
    
- partial reads/writes
    
- connection timeouts
    
- peer identity association
    
- graceful disconnect
    
- malformed input
    
- backpressure
    

Do not assume one TCP `Read()` corresponds to one Protobuf packet.

## Observability

The previously approved Observability & Security Telemetry design  
remains applicable.

The proposal must add telemetry without making telemetry a dependency  
of the cryptographic or routing path.

See [07-Observability-and-Security-Telemetry](../02-architecture/07-Observability-and-Security-Telemetry.md).

## Containerization and GitHub

These are project-wide requirements and must be included in M1/M9:

- reproducible containerized development/test/demo environment
    
- pinned or deliberately managed dependencies
    
- no secrets in images or repository
    
- automated tests in CI
    
- Docker/Compose topology for multi-node demonstrations
    
- repository documentation synchronized with implementation
    
- meaningful milestone commits
    
- clean-environment reproduction test
    

## Required next action

Antigravity should revise `implementation_plan.md` to address the  
BLOCKED handshake items and the packet/routing requirements above.

**No source-code implementation should begin until the revised M0 plan  
is reviewed.**