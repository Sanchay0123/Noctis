# Project State

## Current phase

**M4 complete — post-milestone documentation synchronization / M5 gate**

## Status

| Area | Status |
|---|---|
| Requirements | Documented |
| Threat model | Documented |
| Architecture | Approved baseline |
| Cryptographic architecture | M3 implemented / verified for session crypto scope |
| Protocol | Frozen baseline / M3 implementation verified |
| Networking | M4 direct TCP transport integrated / M5 broader direct-networking scope pending |
| Observability | Approved supporting design / not implemented |
| Containerization | Required / foundation implemented |
| GitHub synchronization | Required / not yet verified in this gate |
| Testing strategy | Documented / M3 crypto evidence added |
| Implementation | M4 direct secure messaging integration implemented |
| Security verification | M4 direct secure messaging integration verified |
| UI | Not started |
| Demo | Not started |
| Final academic material | Draft |

## Current gate

**M3 — Cryptographic Session Layer: COMPLETE / APPROVED**

M0, M0.1, M1.1, M2.1/M2.1-B/M2.1-C, M3.1, M3.2 and M3.3 have been
completed within their approved scopes.

M3 covers:
- fresh ephemeral X25519 key agreement
- Ed25519-authenticated session establishment
- transcript-bound session identification
- HKDF-SHA-256 directional session-key derivation
- ChaCha20-Poly1305 application encryption
- deterministic sequence-derived nonces
- authenticated associated data and role-separated directions
- per-session replay-window enforcement
- concurrency protection for send and receive session state

The M3 implementation was reviewed against the frozen protocol construction
and its reported containerized validation evidence.

## Rule

No implementation component becomes **Approved** without implementation
evidence and relevant tests. A design being frozen or accepted does not mean
that its implementation exists.

## Next action

Complete the technical-lead documentation synchronization for M3, then open
M4 separately. M4 must not be started merely because M3 is complete.

## M3 Completion Gate

**M3 — Authenticated Session & Encrypted Message Layer: 🟢 COMPLETE**

### M3.1 — X25519
Implemented using Go's standard `crypto/ecdh.X25519()` API with fresh
ephemeral keypairs. The implementation exposes raw shared-secret material
only to the session KDF layer, rejects malformed/invalid peer public keys,
supports an explicit `Destroy()` lifecycle boundary, and documents that
destruction is a best-effort memory-lifecycle measure rather than guaranteed
RAM zeroization. The ephemeral-key lifecycle is sequential with respect to
destruction.

### M3.2 — Authenticated Session Establishment
Implemented using the frozen Ed25519/X25519 transcript and HKDF schedule.
The implementation verifies both signed handshake transcripts, derives the
full SHA-256 transcript as the session identifier, and produces independent
initiator-to-responder and responder-to-initiator keys. TOFU remains subject
to its documented first-contact MITM limitation.

### M3.3 — AEAD Message Layer
Implemented with ChaCha20-Poly1305. Each directional session key has an
independent sequence space beginning at zero. The nonce is
`0x00000000 || uint64_be(sequence_num)`. AAD binds the fixed protocol
context, session identifier, sequence number and internally derived
direction marker. A 64-message receive replay window is enforced.

Replay state is only committed after successful AEAD authentication, so an
unauthenticated packet cannot consume receiver replay-window state.
Independent send/receive mutexes protect concurrent session use. Sequence
number exhaustion returns an error and never wraps.

### Validation evidence
The M3 implementation reports containerized `gofmt`, `go vet`, unit tests,
race-detector tests and build validation as passing, including negative and
known-answer tests for the implemented cryptographic constructions.

M3 completion does **not** imply that direct networking, mesh routing,
observability, UI, or production-grade metadata/privacy protection is
complete.

**M4 has been authorized, implemented, independently reviewed, and approved. M5 remains gated until separately authorized by the Project Overseer.**


## M4 Completion Gate

**M4 — Direct Secure Messaging Integration: 🟢 COMPLETE / APPROVED**

M4 integrates the approved M3 cryptographic/session layer with a bounded TCP
transport and the frozen Protobuf packet format. The implementation provides a
direct two-node secure messaging path without changing M3 cryptographic
semantics.

### M4 Transport

The direct transport uses a 4-byte big-endian length prefix followed by a
Protobuf `MeshPacket`. The maximum frame size is 64 KiB and the length is
validated before body allocation. Reads use `io.ReadFull`; writes serialize
complete frames under a connection write mutex and correctly handle short
writes.

### M4 Packet boundary

Incoming packets are structurally validated before cryptographic processing.
Protocol version, packet type/oneof consistency, fixed-width identifiers,
payload-specific fields and unknown Protobuf fields are rejected according to
the version-1 policy.

### M4 Direct secure channel

`DirectChannel` binds one direct TCP connection to one established M3 session
for the M4 scope. INIT and RESP are exchanged through the existing M3 handshake
state machines. APP_DATA is accepted only after session establishment and is
bound to the expected local identity, remote identity and session identifier.

### M4 security evidence

The M4 test suite includes framing, partial/short writes, malformed input,
unknown-field rejection, session binding, end-to-end encrypted message
exchange and transport-path tampering tests. Containerized formatting, vet,
tests, race-detector tests and build validation passed.

M4 establishes direct secure messaging integration but does not establish
mesh routing, multi-hop forwarding, anonymity, metadata hiding, or endpoint
compromise resistance.

**M5 remains NOT STARTED / GATED.**
