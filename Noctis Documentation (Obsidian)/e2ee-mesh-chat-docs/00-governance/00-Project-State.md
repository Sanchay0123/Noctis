# Project State

## Current phase

**M2 complete — post-milestone documentation synchronization / M3 gate**

## Status

  Area                         Status
  --------------------------- -------------
  Requirements                 Documented
  Threat model                 Documented
  Architecture                 Approved baseline
  Cryptographic architecture   Approved baseline
  Protocol                     Frozen baseline / M3 implementation gated
  Networking                   Proposed / not implemented
  Observability                Approved supporting design / not implemented
  Containerization             Required / foundation implemented
  GitHub synchronization       Required / not yet verified in this gate
  Testing strategy             Documented
  Implementation               M2 identity scope implemented
  Security verification        M2 identity scope verified
  UI                           Not started
  Demo                         Not started
  Final academic material      Draft

## Current gate

**M2 --- Documentation Synchronization / Technical-Lead Review**

M0, M0.1, M1.1, and M2.1/M2.1-B/M2.1-C have been completed within their
approved scopes. This synchronization aligns the vault with the completed
M2 identity implementation and the frozen protocol decisions that govern
future M3 work.

## Rule

No implementation component becomes **Approved** without implementation
evidence and relevant tests. A design being frozen or accepted does not
mean its implementation exists.

## Next action

Complete the technical-lead consistency review of this synchronized
documentation set. M3.1 remains unauthorized until that review explicitly
opens the next implementation gate.

## M2 Completion Gate

**M2 — Cryptographic Primitives & Identity: 🟢 COMPLETE**

M2.1 implemented the long-term Ed25519 node identity foundation using Go's
standard `crypto/ed25519` and `crypto/rand`. M2.1-B audited the implementation,
and M2.1-C completed the final key-ownership correction.

Validated properties include secure key generation, 32-byte public identity
representation, private-key encapsulation, defensive copies, safe
malformed-input handling, deterministic Ed25519 signing, and successful
race-detector validation.

**M3 remains unauthorized until the post-M2 documentation synchronization and
technical-lead gate are complete.**
