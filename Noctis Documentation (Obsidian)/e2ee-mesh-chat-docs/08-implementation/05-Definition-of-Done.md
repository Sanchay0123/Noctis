# Definition of Done

A task is **not done** merely because code exists.

## Required

### Implementation

-   [ ] Requirements implemented
-   [ ] Scope unchanged
-   [ ] Architecture boundaries preserved

### Tests

-   [ ] Relevant tests added
-   [ ] Existing tests pass
-   [ ] Negative/security tests pass

### Security

-   [ ] Threat model considered
-   [ ] No known critical issue
-   [ ] Secrets not exposed
-   [ ] Cryptographic choices reviewed

### Documentation

-   [ ] Protocol changes documented
-   [ ] Architecture changes documented
-   [ ] Security implications documented

### Review

-   [ ] Code reviewed
-   [ ] Tests reviewed
-   [ ] Logs reviewed where relevant
-   [ ] Telemetry reviewed for secret/plaintext leakage
-   [ ] Observability failure isolation verified where relevant
-   [ ] Acceptance criteria verified

## Status progression

``` text
Proposed
  ↓
Implemented
  ↓
Tested
  ↓
Verified
  ↓
Approved
```

Only the technical lead can mark a milestone **Approved**.

## Containerization

For tasks affecting runtime/deployment:

-   [ ] Container configuration updated
-   [ ] Clean containerized startup tested
-   [ ] No host-only dependency introduced without documentation
-   [ ] Secrets excluded from images/repository

## GitHub synchronization

For milestone-level work:

-   [ ] Relevant source, tests and documentation committed
-   [ ] Repository status is clean or intentional changes are documented
-   [ ] Commit represents the reviewed implementation state
-   [ ] Canonical GitHub repository updated

## M2 Completion Evidence

M2 is considered complete for the current scope because the Ed25519 identity
foundation has:

- an approved implementation boundary
- standard-library cryptography
- secure randomness
- explicit key ownership
- defensive-copy guarantees
- malformed-input handling
- security-focused regression tests
- containerized validation
- race-detector validation
- synchronized implementation documentation

M2 completion does not constitute approval of M3 session establishment.


## M4 Completion Evidence

M4 is complete and approved because the direct secure messaging integration
met the milestone acceptance boundary:

- bounded 4-byte length-prefixed TCP framing
- correct short-write handling under a serialized writer lock
- explicit unknown-Protobuf-field rejection
- structural validation for INIT, RESP and APP_DATA
- established-session requirement for APP_DATA
- session ID and endpoint identity binding
- authenticated plaintext delivery boundary
- bidirectional end-to-end encrypted messaging test
- transport-path ciphertext tampering test
- short-write and failing-write regression tests
- `go vet`, `go test`, `go test -race` and build validation passed
- no M3 cryptographic semantics were changed

M4 approval does not constitute approval of mesh routing, multi-hop forwarding,
anonymity, metadata hiding or production-grade endpoint security.

**Next milestone:** M5 — Direct networking hardening / runtime integration.
