# M2 Synchronization Note

This document records the earlier M0 revised architecture review. Its
historical M3 blockers have since been resolved at the specification level:
the canonical handshake transcript, session identifier, and HKDF schedule
are now frozen in [03-Session-Establishment](../04-protocol/03-Session-Establishment.md).

M3 implementation is still **blocked**, but the remaining gate is now
implementation-level review and authorization, including the exact X25519
API choice, low-order/all-zero handling, ephemeral-key lifecycle,
zeroization expectations, and corresponding tests. Do not interpret the
historical blocker list below as the current implementation task.

---

# M0 Revised Architecture Review — Antigravity Plan

## Review status

**🟡 APPROVED WITH CHANGES**

The revised plan addresses the major M0 feedback: explicit bootstrap
peers, bounded flooding, TCP framing, canonical transcript intent,
fresh ephemeral X25519 keys, in-memory session state, Protobuf schema,
containerization, CI, and observability.

The plan is sufficient to authorize **M1 project skeleton/tooling work**.

It is **not yet sufficient to authorize M3 cryptographic protocol
implementation**. The handshake and session identifier require the
mandatory refinements below.

## Evidence reviewed

Antigravity states that this revision incorporates the M0 feedback and
defines the revised architecture before implementation. The proposal
explicitly separates Ed25519 identity keys from ephemeral X25519 keys,
defines a canonical transcript, adds strict TCP framing, bounded
managed flooding, and adds containerization and non-blocking
observability. fileciteturn1file0L4-L16

The revised proposal also specifies the transcript structure and
session-key lifecycle, including in-memory-only session state across
restart. fileciteturn1file0L34-L60 fileciteturn1file0L68-L87

## Decisions

### 1. Go

**✅ APPROVED**

Go remains appropriate. No requirement forces Rust.

### 2. Protobuf

**✅ APPROVED**

Protobuf is approved as the wire serialization format.

The `.proto` schema must be version-controlled and tested as a security
boundary.

### 3. TCP with length-prefixed framing

**✅ APPROVED**

The explicit statement that TCP is a byte stream and that one `Read()`
does not equal one packet is correct. The proposed 4-byte length prefix
is suitable for the prototype. fileciteturn1file0L143-L149

M1/M5 must also define maximum frame length and reject oversized frames
before allocating unbounded memory.

### 4. Explicit bootstrap peers

**✅ APPROVED**

Deferring UDP broadcast is a good reduction in initial attack surface and
implementation complexity. fileciteturn1file0L145-L146

### 5. Managed flooding

**✅ APPROVED WITH CHANGES**

Managed flooding remains suitable for the educational prototype.

The bounded seen-cache, TTL and fan-out controls are good additions.
However, the proposed values are **engineering parameters, not
cryptographic guarantees**.

The implementation must define:

- maximum cache entries
- cache expiration
- eviction behavior
- packet/frame size
- maximum peers
- forwarding/backpressure limits
- behavior when the cache is full

The statement that a 60-second cache is safe merely because TTL is 16 is
not sufficient as a standalone proof. Timing, queueing and delayed
delivery must also be considered.

### 6. Observability

**✅ APPROVED**

The non-blocking observer architecture is correct.

Observability must remain outside the critical cryptographic and routing
paths, must not receive plaintext or secret keys, and must not be
required for message delivery. The proposal explicitly adopts this
separation. fileciteturn1file0L201-L209

### 7. Containerization and GitHub

**✅ APPROVED**

M1 must include reproducible Docker/Compose development and test
execution, dependency management, CI, and secret exclusion. The proposal
already places containerization in M1 and reproducibility in M9.
fileciteturn1file0L216-L218 fileciteturn1file0L241-L245

GitHub synchronization remains a project-wide Definition of Done.

## Cryptographic review

### A. Canonical transcript

**🟡 APPROVED IN PRINCIPLE — SPECIFICATION MUST BE FROZEN BEFORE M3**

The proposal correctly binds:

- protocol version
- domain label
- initiator identity
- responder identity
- both ephemeral X25519 public keys

and distinguishes INIT from RESP by the responder ephemeral field.
fileciteturn1file0L34-L49

However, the implementation must not independently invent byte encoding.
The canonical encoding must be specified as an exact function, including
the domain-label representation and field boundaries.

Recommended construction:

```text
T =
  "MeshChat-Handshake-v1" ||
  u8(protocol_version) ||
  ID_initiator ||
  ID_responder ||
  E_initiator ||
  E_responder
```

with every field fixed-width except the domain label, which should be
defined as a fixed protocol constant rather than an ambiguously encoded
string.

The final implementation specification must also explicitly state that
the INIT signature authenticates a transcript containing a fixed
all-zero responder ephemeral placeholder, while the RESP signature
authenticates the complete transcript.

### B. Session ID

**🛑 BLOCKER**

The proposed:

```text
SHA256(T_RESP)[0:4]
```

is too small for a protocol session identifier.

A 32-bit identifier has an unacceptable collision probability as the
number of sessions grows. It can also create ambiguous session-state
lookup.

**Required change:**

Use a substantially larger identifier, preferably:

```text
session_id = SHA-256(T_RESP)
```

represented as a fixed 32-byte value.

If implementation simplicity requires a shorter identifier, the design
must instead provide an explicit collision-handling mechanism. The
preferred decision is the full 32-byte transcript hash.

### C. KDF

**🟡 APPROVED WITH REFINEMENT**

Using the transcript hash as HKDF salt and directional `info` labels is
reasonable. fileciteturn1file0L54-L62

The final specification must make the KDF deterministic and exact:

```text
PRK = HKDF-Extract(
    salt = SHA256(T_RESP),
    IKM  = SS
)

K_A_to_B = HKDF-Expand(
    PRK,
    info = "MeshChat-v1|session|initiator->responder",
    L = 32
)

K_B_to_A = HKDF-Expand(
    PRK,
    info = "MeshChat-v1|session|responder->initiator",
    L = 32
)
```

The exact strings/bytes become protocol constants.

### D. Forward secrecy wording

**🟡 APPROVED WITH WORDING CONSTRAINT**

The design provides forward secrecy for completed sessions under the
stated endpoint/key-compromise assumptions, assuming ephemeral private
keys are securely erased and long-term keys are not used for DH.

Do not claim protection against endpoint compromise while a session key
is active.

The proposal correctly states that it does not provide PCS. fileciteturn1file0L81-L87

### E. Handshake replay

**🛑 MUST BE SPECIFIED BEFORE M3**

The statement:

> "Bob will reject an INIT with an E_A he has already processed"

is directionally correct but insufficient as a protocol definition.
fileciteturn1file0L64-L66

Required:

- bounded handshake replay cache
- cache key definition
- expiration
- duplicate INIT behavior
- duplicate RESP behavior
- concurrent INIT handling
- timeout behavior
- expected-peer/session binding
- behavior after a completed session receives an old handshake

Do not rely solely on "fresh ephemeral keys" for replay protection.

## Packet schema review

**🟡 APPROVED AS DRAFT**

The oneof-based packet structure is a significant improvement.

However:

### Unknown fields

The proposal says unknown fields will cause rejection. Standard Protobuf
unknown-field behavior does not automatically mean "reject unknown
fields." The application must explicitly implement any such policy.

For version 1, rejecting an unsupported **protocol version** is mandatory.
Unknown fields should either be explicitly rejected by a documented
validation layer or tolerated according to a deliberate compatibility
policy. Do not accidentally depend on library behavior.

### Routing versus payload

The relay must inspect only the minimum routing envelope required for
forwarding. It must not need to parse application content.

### TTL

The forwarding rule must distinguish destination processing from
forwarding:

```text
if packet is for local destination:
    process according to protocol state
else if ttl <= 1:
    drop
else:
    decrement ttl
    forward
```

Do not accidentally drop a valid packet arriving at its destination with
TTL == 1.

### Size limits

The 64 KiB proposal is acceptable as an initial prototype limit, but the
length prefix must be checked before allocating the frame buffer.

## M1 authorization

**M1 MAY PROCEED.**

M1 scope:

- Go module
- repository layout
- Protobuf tooling/schema skeleton
- Dockerfile
- Compose topology skeleton
- CI
- dependency management
- `.gitignore`
- configuration/secrets policy
- deterministic test command

M1 must not silently implement the cryptographic handshake.

## M2 authorization

**M2 MAY PROCEED** for isolated Ed25519 identity primitives and tests,
provided the code does not invent the final session protocol.

## M3 authorization

**🛑 BLOCKED**

Before M3 begins, Antigravity must produce the exact final:

1. handshake wire messages
2. canonical transcript function
3. signature inputs
4. session-ID construction
5. HKDF extract/expand construction
6. handshake replay state machine
7. duplicate/retransmission behavior
8. timeout behavior
9. key-erasure lifecycle

## Required next response from Antigravity

Revise the implementation plan again only for the blockers above.

Do not redesign the whole project.

The next plan should specifically resolve:

- 32-byte session ID
- exact canonical transcript encoding
- exact KDF construction
- bounded handshake replay state
- exact handshake state transitions
- Protobuf unknown-field policy
- TTL destination/forwarding semantics
- pre-allocation frame-size enforcement

Once those are resolved, M3 can be re-reviewed for authorization.

## Overall conclusion

The revised architecture is materially better and demonstrates that the
previous M0 feedback was incorporated.

**M1: AUTHORIZED**

**M2: AUTHORIZED WITH SCOPE**

**M3: BLOCKED PENDING CRYPTOGRAPHIC SPECIFICATION**

This is a deliberate gate: the project should build the scaffolding now,
but should not freeze a custom cryptographic protocol from an
underspecified document.
