
```
You are working as the implementation agent for the E2EE Mesh Chat project.

You have NOT been authorized to begin M1, M2, M3, or any implementation of cryptographic functionality.

Your immediate task is to perform **M0.1 — Repository & Architecture Initialization/Audit** only.

## PROJECT AUTHORITY

The architecture and security requirements are governed by the project's architecture documentation and the latest architecture review.

Security and architectural correctness take priority over implementation convenience.

Do NOT invent protocol semantics that are not specified.

Do NOT silently modify cryptographic design decisions.

Do NOT begin implementing:

- Ed25519 identity operations
    
- X25519 handshake
    
- HKDF session derivation
    
- ChaCha20-Poly1305 messaging
    
- replay protection implementation
    
- mesh routing implementation
    
- application messaging
    
- UI functionality
    

This task is strictly preparation and architecture validation.

---

# OBJECTIVE

Inspect the current repository and establish the project foundation required to begin controlled implementation.

The repository must be organized so that subsequent milestones can be implemented independently and reviewed.

---

# TASK 1 — REPOSITORY AUDIT

Inspect the entire repository.

Report:

- current directory structure
    
- existing source files
    
- existing tests
    
- existing configuration
    
- existing Docker/Compose files
    
- existing CI configuration
    
- existing documentation
    
- existing generated files
    
- existing dependencies
    
- existing secrets or potentially sensitive files
    
- existing implementation that may conflict with the approved architecture
    

Do NOT delete or rewrite existing code automatically.

If something conflicts with the architecture, report it first.

---

# TASK 2 — ESTABLISH TARGET PROJECT STRUCTURE

Propose a clean Go project structure consistent with the architecture.

The structure should clearly separate at minimum:

- application/UI
    
- cryptographic/session layer
    
- protocol/serialization
    
- networking/transport
    
- mesh/routing
    
- observability
    
- configuration
    
- tests
    

The architecture requires strict separation between:

1. routing/transport
    
2. cryptographic/session processing
    
3. application message handling
    

Relays must never require application plaintext access.

Do not implement these systems yet.

Create only the minimum skeleton required to support later milestones.

---

# TASK 3 — DEPENDENCY POLICY

Identify the exact Go dependencies that will eventually be required.

The architecture currently specifies:

- Go
    
- crypto/ed25519
    
- golang.org/x/crypto/curve25519
    
- golang.org/x/crypto/chacha20poly1305
    
- golang.org/x/crypto/hkdf
    
- crypto/rand
    
- google.golang.org/protobuf
    

Do not add unnecessary cryptographic libraries.

Do not implement cryptographic primitives.

Do not introduce custom cryptography.

If dependencies are added now, explain why each one is required.

---

# TASK 4 — SECURITY BOUNDARIES

Create or validate clear package/module boundaries preventing accidental coupling.

In particular:

- routing must not decrypt application payloads
    
- observability must not receive plaintext
    
- private keys must not be exposed to routing
    
- session keys must not be exposed to telemetry
    
- protocol serialization must not silently reinterpret security-sensitive fields
    
- configuration must not contain committed secrets
    

Document these boundaries.

---

# TASK 5 — SECRET & GENERATED-FILE POLICY

Establish:

- .gitignore
    
- secret handling policy
    
- generated-code policy
    
- local development state policy
    
- private-key storage policy
    
- environment-variable policy
    

Ensure private keys, passwords, `.env` files, generated runtime state, and similar sensitive material cannot accidentally enter Git.

Do NOT create real private keys.

---

# TASK 6 — BUILD/TEST BASELINE

Establish a deterministic baseline command for:

- formatting
    
- compilation
    
- unit tests
    

At this stage, tests may only be infrastructure/skeleton tests.

Do not write cryptographic tests yet.

Do not create fake cryptographic implementations merely to satisfy compilation.

---

# TASK 7 — CONTAINERIZATION BASELINE

Inspect whether Docker/Compose infrastructure exists.

If it does not exist, create only the minimal reproducible development/test skeleton.

Requirements:

- reproducible Go environment
    
- no secrets baked into images
    
- no host-only undocumented dependency
    
- deterministic build/test entry point
    
- architecture must remain usable if observability services are unavailable
    

Do NOT create the final multi-node demonstration topology yet.

---

# TASK 8 — DOCUMENTATION CONSISTENCY

Compare the repository against the architecture documentation.

Identify any discrepancy involving:

- protocol terminology
    
- package responsibilities
    
- cryptographic primitives
    
- packet structure
    
- networking
    
- routing
    
- observability
    
- containerization
    
- milestone boundaries
    

Do not silently resolve architectural contradictions.

List them for review.

---

# HARD SECURITY CONSTRAINTS

The following are NON-NEGOTIABLE:

1. No custom cryptographic primitives.
    
2. No plaintext application data in relay/routing code.
    
3. No private keys in Git.
    
4. No session keys in logs/metrics.
    
5. No secrets in Docker images.
    
6. No observability dependency in the E2EE critical path.
    
7. No implementation of M3 cryptographic handshake during this task.
    
8. No protocol changes without explicit architecture approval.
    
9. No claim of anonymity, unhackability, or perfect security.
    
10. Do not mark M0 complete yourself.
    

---

# REQUIRED OUTPUT

After completing the task, report exactly:

TASK:  
M0.1 — Repository & Architecture Initialization/Audit

FILES CHANGED:  
List every changed/created file.

ARCHITECTURE CHANGES:  
List any architectural decisions made.  
If none, explicitly state "None".

SECURITY CHANGES:  
List security-relevant changes.

TESTS:  
List commands executed and their results.

CONTAINERIZATION:  
List Docker/Compose changes and validation performed.

DEPENDENCIES:  
List added/removed dependencies and justification.

DOCUMENTATION:  
List documentation discrepancies discovered.

KNOWN ISSUES:  
List anything unresolved.

BLOCKERS:  
List anything that requires architecture approval.

COMMANDS:  
List the exact commands used for validation.

NEXT RECOMMENDATION:  
Recommend the next M0 task, but do not begin it automatically.

IMPORTANT:  
Do not report "M0 complete".  
Do not proceed to M1/M2/M3.  
Stop after M0.1 and wait for architectural review.
```


```
You are continuing the E2EE Mesh Chat project.

The technical lead has reviewed your M0.1 report.

## VERDICT

M0.1 is **CONDITIONALLY APPROVED**.

You are NOT authorized to begin M1 implementation yet.

Your immediate task is:

# M0.1-B — VALIDATION & CLOSURE

Do not add new application functionality.

Do not implement:

- Ed25519 identity operations
    
- X25519
    
- HKDF
    
- ChaCha20-Poly1305
    
- handshake logic
    
- replay protection
    
- encrypted messaging
    
- mesh routing
    
- peer discovery
    
- application messaging
    

This task exists only to validate the foundation you created in M0.1 and correct foundation-level problems.

---

# 1. INSPECT THE ACTUAL FILES

Review every file created during M0.1.

Report the exact contents/behavior of:

- go.mod
    
- .gitignore
    
- Dockerfile
    
- docker-compose.yml
    
- cmd/meshchat/main.go
    
- tests/skeleton_test.go
    

Also inspect the package structure under:

- internal/app
    
- internal/crypto
    
- internal/protocol
    
- internal/mesh
    
- internal/transport
    
- internal/observability
    
- internal/config
    

Do not assume that an empty directory establishes a real package boundary.

---

# 2. DOCKER-BASED GO VALIDATION

Because the host does not have Go available, validation MUST be performed through the containerized toolchain.

Run the equivalent of:

go fmt ./...  
go vet ./...  
go test ./...  
go build ./...

If the current Dockerfile cannot perform these operations reproducibly, fix the Docker development/test setup.

Do not install Go directly onto the host merely to bypass the container requirement.

Record the exact commands and results.

---

# 3. DOCKERFILE VALIDATION

Build the Docker image from a clean state.

Verify:

- build succeeds
    
- Go version is deterministic
    
- dependencies resolve successfully
    
- no host-only dependency is required
    
- no secrets are copied into the image
    
- no private keys are generated during image build
    
- repository source is sufficient to reproduce the build
    

If the Dockerfile uses a mutable base such as an unpinned `golang:1.21-alpine`, report that explicitly.

Do NOT introduce unnecessary complexity solely for image pinning unless justified.

---

# 4. COMPOSE VALIDATION

Validate the docker-compose configuration.

Verify:

- YAML/configuration is valid
    
- services can start
    
- node-alice and node-bob have the intended isolated network
    
- no host secret files are mounted
    
- no `.env` or private-key material is copied
    
- services do not depend on Prometheus/Grafana/Loki
    
- the topology remains usable without observability
    

If Compose validation requires actually starting the services, do so and record the result.

---

# 5. PACKAGE BOUNDARY CHECK

The architecture requires separation between:

Application  
↓  
Cryptographic/session layer  
↓  
Protocol  
↓  
Transport  
↓  
Mesh/routing

and observability must remain an out-of-band observer.

Check whether the current Go package structure actually supports this.

Do not implement the layers yet.

If the empty directories are not sufficient to establish package-level boundaries, create only the minimal package declarations/interfaces needed to establish the intended dependency direction.

Do not introduce fake implementations.

Do not create placeholder crypto functions that return dummy keys/ciphertext.

---

# 6. SECRET POLICY TEST

Inspect `.gitignore`.

Verify that it protects at minimum:

- `.env`
    
- `.env.*` where appropriate
    
- private key files
    
- certificate/key material where appropriate
    
- secrets/
    
- session_keys/
    
- runtime/generated state
    
- compiled binaries
    
- local development state
    

Do not over-broaden `.gitignore` in a way that hides legitimate source/documentation.

Run an appropriate repository check to verify that no obvious secret files are currently tracked.

Do not create real secrets for testing.

---

# 7. DEPENDENCY CHECK

Inspect go.mod/go.sum.

Confirm that dependencies correspond only to the approved technology plan.

Expected cryptographic dependencies include:

- golang.org/x/crypto
    
- google.golang.org/protobuf
    

Standard-library cryptography should remain preferred where already specified.

Do not add:

- alternative crypto libraries
    
- unnecessary networking frameworks
    
- unnecessary web frameworks
    
- unrelated dependencies
    

If a dependency is present but unused, report whether it should remain at this stage.

---

# 8. REPRODUCIBILITY CHECK

Starting from the repository itself, verify that another developer could reproduce the baseline using Docker.

The validation sequence should be documented.

There must be a deterministic way to:

1. build
    
2. format
    
3. vet
    
4. test
    
5. compile
    

If any step depends on an undocumented host dependency, fix or document it.

---

# 9. GITHUB READINESS

Inspect the repository for accidental generated/runtime artifacts.

Confirm:

- no secrets
    
- no private keys
    
- no credentials
    
- no local database/state
    
- no compiled binaries
    
- no unnecessary IDE state
    

Do not push anything to GitHub unless the repository's existing Git workflow explicitly permits it.

If Git metadata exists, report:

- current branch
    
- working tree status
    
- whether files are tracked
    
- whether a meaningful initial commit exists
    

Do not create a commit merely to claim completion unless instructed.

---

# 10. DO NOT START M1

This is critical.

Although your previous report recommended:

"M1 Execution: Finalizing the Protobuf schema definitions and establishing the generated-code workflow"

DO NOT perform that task yet.

M1 will be authorized separately after M0.1-B passes review.

---

# REQUIRED REPORT

Return exactly:

TASK:  
M0.1-B — Validation & Closure

FILES CHANGED:  
List every file changed.

VALIDATION:  
List exact commands and results.

DOCKER:  
Build result and relevant output.

COMPOSE:  
Validation/startup result.

PACKAGE BOUNDARIES:  
Explain the actual dependency direction and any changes made.

SECRET AUDIT:  
Result.

DEPENDENCY AUDIT:  
Result.

REPRODUCIBILITY:  
Result.

GIT STATUS:  
Current branch/status/tracking information.

SECURITY ISSUES:  
Any issues found.

ARCHITECTURE ISSUES:  
Any issues found.

KNOWN ISSUES:  
Anything unresolved.

BLOCKERS:  
Anything requiring technical-lead approval.

NEXT RECOMMENDATION:  
State what you recommend next, but DO NOT execute it.

Do not declare M0.1 complete.  
Do not begin M1.  
Stop and wait for architectural review.
```


```
The technical lead has completed the M0.1-B review.

## M0.1 STATUS

**APPROVED.**

The repository foundation, containerized validation, Docker build, Compose topology, package skeleton, secret policy, and baseline tooling have passed the M0.1 gate.

You are now authorized to begin:

# M1.1 — Protobuf Schema & Generated-Code Workflow

This is the FIRST M1 implementation task.

Do NOT implement M2, M3, M4, networking, routing, or application messaging.

---

# OBJECTIVE

Establish the canonical Protocol Buffers schema and a reproducible containerized generated-code workflow.

The purpose is to make the wire-level data structures concrete while keeping all cryptographic/session behavior unimplemented.

---

# REQUIRED SCHEMA

Use the currently approved protocol structure as the starting point.

The canonical structure is:

```protobuf
syntax = "proto3";

enum PacketType {
  PACKET_TYPE_UNKNOWN = 0;
  PACKET_TYPE_INIT = 1;
  PACKET_TYPE_RESP = 2;
  PACKET_TYPE_APP_DATA = 3;
}

message MeshPacket {
  uint32 version = 1;
  PacketType type = 2;

  bytes packet_id = 3;
  uint32 ttl = 4;
  bytes source_node = 5;
  bytes dest_node = 6;

  oneof payload {
    InitPayload init = 7;
    RespPayload resp = 8;
    AppDataPayload app_data = 9;
  }
}

message InitPayload {
  bytes ephemeral_key = 1;
  bytes signature = 2;
}

message RespPayload {
  bytes ephemeral_key = 1;
  bytes signature = 2;
}

message AppDataPayload {
  bytes session_id = 1;
  uint64 sequence_num = 2;
  bytes ciphertext = 3;
}

Do NOT silently change the wire structure.

If you believe a schema change is necessary, STOP and report it as an architecture issue instead of making the change silently.

---

# 1. PROTO FILE ORGANIZATION

Create the canonical `.proto` source in an appropriate repository location.

The schema source must be clearly distinguished from generated Go code.

Do not manually write generated `.pb.go` code.

---

# 2. FIELD VALIDATION POLICY

The Protobuf schema itself cannot enforce all semantic byte-length constraints.

Document the following semantic validation requirements for later protocol implementation:

- `packet_id` = exactly 16 bytes
    
- `source_node` = exactly 32 bytes
    
- `dest_node` = exactly 32 bytes
    
- `InitPayload.ephemeral_key` = exactly 32 bytes
    
- `InitPayload.signature` = exactly 64 bytes
    
- `RespPayload.ephemeral_key` = exactly 32 bytes
    
- `RespPayload.signature` = exactly 64 bytes
    
- `AppDataPayload.session_id` = exactly 32 bytes
    
- `ciphertext` must satisfy the later AEAD minimum-size requirement
    

Do NOT implement the semantic packet validator yet unless it is strictly required for the schema/code-generation workflow.

---

# 3. PACKET TYPE / PAYLOAD CONSISTENCY

Record this as a protocol validation invariant for later implementation:

- `PACKET_TYPE_INIT` → `init`
    
- `PACKET_TYPE_RESP` → `resp`
    
- `PACKET_TYPE_APP_DATA` → `app_data`
    
- `PACKET_TYPE_UNKNOWN` → reject during protocol validation
    

Do not implement handshake or encryption to enforce this.

---

# 4. UNKNOWN-FIELD POLICY

The architecture requires security-sensitive Protobuf messages to reject unknown fields.

Establish/document the intended Go unmarshaling policy.

Important:

Do not merely write a comment claiming unknown fields are rejected.

Determine the actual `google.golang.org/protobuf` API/configuration that will be used later.

If the generated-code workflow itself cannot enforce this at schema level, document that the runtime unmarshal operation must use explicit unknown-field rejection.

Create a small test only if it is appropriate to validate the serialization/unmarshaling policy without implementing application or cryptographic behavior.

---

# 5. CODE GENERATION

Establish reproducible generation of Go code from the `.proto` schema.

Requirements:

- generator version must be controlled/documented
    
- generated code must be reproducible
    
- generation must work inside Docker
    
- host-installed protoc/Go tooling must NOT be required
    
- generated output must be placed in the intended protocol package
    
- developers must have a documented generation command
    

Do not rely on an undocumented global binary installed on the host.

---

# 6. DOCKER INTEGRATION

Extend the existing containerized development workflow as necessary.

The following must work from a clean repository:

1. install/obtain the required generation tooling inside the controlled environment
    
2. generate `.pb.go`
    
3. format generated/source Go code
    
4. run tests
    
5. build the project
    

Do not add unnecessary services to Docker Compose.

Do not add Prometheus, Grafana, Loki, databases, or external infrastructure at this stage.

---

# 7. DEPENDENCY MANAGEMENT

Resolve the dependency situation from M0.1 carefully.

The previous validation noted that `go mod tidy` on a skeleton with no imports can remove unused dependencies.

Do NOT manually preserve unused dependencies merely because the architecture eventually needs them.

Once the generated Protobuf code actually imports:

```text
google.golang.org/protobuf

the dependency should be represented naturally in the module graph.

Do not add cryptographic dependencies unless this task genuinely requires them.

M1.1 is NOT an authorization to implement crypto.

---

# 8. TESTS

Add tests appropriate to this task.

At minimum verify:

- generated types compile
    
- valid schema-generated structures can marshal
    
- valid structures can unmarshal
    
- malformed/truncated serialized data is rejected
    
- unknown-field behavior is explicitly tested/documented if runtime unmarshaling is introduced
    
- generated code does not break `go vet`
    
- complete repository still passes `go test ./...`
    
- complete repository still builds
    

Do not create fake cryptographic tests.

---

# 9. SECURITY REQUIREMENTS

Do NOT:

- put plaintext messages into telemetry
    
- add private keys
    
- add session keys
    
- implement signatures
    
- implement X25519
    
- implement HKDF
    
- implement AEAD
    
- implement handshake state machines
    
- implement replay protection
    
- implement routing
    

The schema is only the wire representation.

---

# 10. DOCUMENTATION

Update only the relevant project documentation to reflect:

- canonical `.proto` location
    
- generated-code workflow
    
- generator/tool versions
    
- semantic field-size requirements
    
- packet-type/payload invariant
    
- unknown-field policy
    
- how a clean Docker environment generates and validates the code
    

Do not rewrite unrelated documentation.

---

# ACCEPTANCE CRITERIA

M1.1 passes only if:

[ ] Canonical `.proto` exists.

[ ] Schema matches the approved architecture.

[ ] Generated Go code is produced automatically.

[ ] Generated code is not manually authored.

[ ] Code generation works in Docker.

[ ] Generator versions are controlled/documented.

[ ] Host-installed Protobuf tooling is not required.

[ ] `go fmt ./...` passes.

[ ] `go vet ./...` passes.

[ ] `go test ./...` passes.

[ ] `go build ./...` passes.

[ ] Unknown-field policy is explicitly implemented/documented.

[ ] Semantic field-size constraints are documented for runtime validation.

[ ] No cryptographic implementation has been introduced.

[ ] No networking/routing implementation has been introduced.

[ ] No secrets/private keys/session keys are present.

[ ] Docker/Compose still work.

---

# REQUIRED REPORT

Return exactly:

TASK:  
M1.1 — Protobuf Schema & Generated-Code Workflow

FILES CHANGED:  
List every created/modified file.

SCHEMA:  
Describe the `.proto` location and confirm whether it matches the approved structure.

GENERATED CODE:  
Describe generated files and generation command.

TOOLCHAIN:  
List protoc/protoc-gen-go or equivalent versions and how they are obtained.

DOCKER:  
Explain how generation and validation work inside Docker.

DEPENDENCIES:  
List changes to go.mod/go.sum.

VALIDATION:  
Provide exact commands and results for:

- formatting
    
- vet
    
- tests
    
- build
    
- protobuf generation
    

SECURITY:  
Confirm explicitly that no crypto/session/networking/routing implementation was introduced.

DOCUMENTATION:  
List documentation files changed.

KNOWN ISSUES:  
List unresolved issues.

ARCHITECTURE ISSUES:  
List anything that requires technical-lead review.

BLOCKERS:  
List blockers.

NEXT RECOMMENDATION:  
Recommend the next M1 task, but DO NOT execute it.

Do not declare M1 complete.  
Do not begin M1.2 automatically.  
Stop and wait for technical-lead review.

```


```
The technical lead has reviewed M1.1.

# M1.1 STATUS

**APPROVED.**

The canonical Protobuf schema and containerized generated-code workflow have passed review.

You are now authorized to begin:

# M2.1 — Ed25519 Identity Core

This is the first cryptographic implementation task.

## STRICT SCOPE

Implement ONLY the long-term Ed25519 node identity foundation.

Do NOT implement:

- X25519
    
- Diffie-Hellman
    
- HKDF
    
- ChaCha20-Poly1305
    
- handshake/session establishment
    
- session keys
    
- encrypted messaging
    
- replay protection
    
- routing
    
- peer discovery
    
- network communication
    
- key exchange
    

M2.1 is isolated identity functionality only.

---

# OBJECTIVE

Create a small, well-tested identity abstraction around the Go standard library:

`crypto/ed25519`

and secure randomness:

`crypto/rand`

The implementation must not implement cryptographic primitives itself.

---

# 1. KEY GENERATION

Implement secure Ed25519 key generation using the approved standard-library API and cryptographically secure randomness.

Requirements:

- use `crypto/ed25519`
    
- use `crypto/rand`
    
- never use math/rand
    
- never derive identity keys from predictable values
    
- never hard-code keys
    
- never generate identity keys during Docker image construction
    
- never generate identity keys at package initialization
    

Each generated identity must have:

- Ed25519 private key
    
- Ed25519 public key
    

The public key must correspond exactly to the private key.

---

# 2. IDENTITY ABSTRACTION

Create a minimal identity abstraction in:

`internal/crypto`

It should make the intended boundary clear.

The abstraction should support, at minimum:

- generation
    
- public-key retrieval
    
- signing
    
- signature verification
    

Do not expose unnecessary implementation details.

Do not expose session-key concepts.

Do not create interfaces merely for abstraction's sake if a concrete type is clearer.

Keep the implementation small.

---

# 3. IDENTITY REPRESENTATION

The project's protocol defines node identity as the 32-byte Ed25519 public key.

Therefore:

- public identity representation = exactly 32 bytes
    
- private Ed25519 key remains local secret material
    
- identity serialization must be explicit
    
- malformed public/private key lengths must be rejected
    

Do not introduce usernames, passwords, account IDs, UUID-based identities, or alternate identity schemes.

---

# 4. SIGNING

Implement Ed25519 signing using the standard library.

Requirements:

- signing accepts arbitrary message bytes
    
- signature is exactly the Ed25519 signature size
    
- signing requires possession of the private key
    
- verification uses the public key
    
- invalid signatures must fail verification
    
- modified messages must fail verification
    
- signatures from another identity must fail verification
    

Important:

Do not confuse "deterministic signing" with deterministic key generation.

Ed25519 signatures are deterministic for a given private key and message; identity key generation must use cryptographically secure randomness.

---

# 5. INPUT VALIDATION

Define behavior for:

- nil/empty message
    
- malformed public key
    
- malformed private key
    
- malformed signature
    
- mismatched public/private key
    

Do not panic on attacker-controlled byte slices.

Prefer explicit errors where appropriate.

Do not silently truncate or pad keys.

Do not accept alternate key sizes.

---

# 6. PRIVATE KEY HANDLING

Private keys are sensitive.

Requirements:

- never log private keys
    
- never include private keys in errors
    
- never include private keys in metrics
    
- never serialize private keys into protocol packets
    
- never write private keys to disk in this task
    
- never place private keys in test fixtures committed to Git
    
- never place private keys in Docker images
    

Do not claim that ordinary Go memory handling guarantees physical memory erasure.

No secure-storage implementation is required yet.

---

# 7. TESTS

Create comprehensive unit tests for the identity abstraction.

At minimum:

### Generation

- generated keypair has correct sizes
    
- generated public key matches private key
    
- two independently generated identities are overwhelmingly unlikely to have the same public key
    

### Signing

- valid signature verifies
    
- modified message fails verification
    
- signature from identity A fails under identity B's public key
    
- repeated signing of the same message with the same identity produces the expected Ed25519 deterministic signature behavior
    

### Validation

- malformed public key rejected
    
- malformed signature rejected
    
- malformed private key rejected
    
- mismatched key material handled safely
    
- nil/empty inputs behave according to documented policy
    

### Security

Tests must not print private keys.

Tests must not commit fixed private keys.

Use freshly generated test keys.

---

# 8. PROTOCOL SEPARATION

Do NOT modify the handshake protocol.

Do NOT connect Ed25519 signing to X25519.

Do NOT construct:

- T_INIT
    
- T_RESP
    
- session IDs
    
- HKDF inputs
    
- session keys
    

Those belong to M3.

M2.1 only provides the identity primitive that M3 will later consume.

---

# 9. DEPENDENCY POLICY

Prefer the Go standard library.

No new external cryptographic dependency should be introduced for Ed25519.

Do not add another crypto library.

---

# 10. CONTAINERIZED VALIDATION

All validation must continue to work without host-installed Go.

Run the project validation through Docker.

At minimum:

- gofmt
    
- go vet
    
- go test
    
- go build
    

The existing Protobuf generation workflow must remain functional.

Docker/Compose must continue to work.

---

# 11. DOCUMENTATION

Update only the relevant documentation.

Document:

- Ed25519 as the long-term identity primitive
    
- 32-byte public identity representation
    
- private-key handling boundary
    
- signing/verification responsibility
    
- distinction between identity and future ephemeral X25519 session keys
    
- M2.1 test coverage
    

Do not rewrite the entire security architecture.

---

# ACCEPTANCE CRITERIA

M2.1 is acceptable only if:

[ ] Ed25519 uses Go's standard library.

[ ] Key generation uses cryptographically secure randomness.

[ ] Identity public key is exactly 32 bytes.

[ ] Private key is never exposed through logs/metrics/errors.

[ ] Signing and verification work correctly.

[ ] Modified messages fail verification.

[ ] Wrong identities fail verification.

[ ] Malformed inputs are handled safely.

[ ] No X25519 implementation exists.

[ ] No HKDF implementation exists.

[ ] No AEAD implementation exists.

[ ] No handshake implementation exists.

[ ] No session-key implementation exists.

[ ] No networking/routing implementation exists.

[ ] No private keys are persisted.

[ ] No secrets are committed.

[ ] Existing Protobuf functionality remains intact.

[ ] Docker-based fmt/vet/test/build pass.

[ ] Relevant documentation is synchronized.

---

# REQUIRED REPORT

Return exactly:

TASK:  
M2.1 — Ed25519 Identity Core

FILES CHANGED:  
List every created/modified file.

CRYPTOGRAPHY:  
Describe exactly which standard-library APIs were used.

IDENTITY:  
Describe the identity representation and key lifecycle.

SIGNING:  
Describe signing and verification behavior.

PRIVATE KEY HANDLING:  
Describe all protections applied.

TESTS:  
List every identity/security test and its result.

VALIDATION:  
Provide exact Docker-based commands and results for:

- fmt
    
- vet
    
- test
    
- build
    
- Protobuf generation
    

DEPENDENCIES:  
List go.mod/go.sum changes.

DOCUMENTATION:  
List documentation files changed.

SECURITY ISSUES:  
List any issues.

ARCHITECTURE ISSUES:  
List anything requiring technical-lead review.

KNOWN ISSUES:  
List unresolved issues.

BLOCKERS:  
List blockers.

NEXT RECOMMENDATION:  
Recommend the next M2 task, but DO NOT execute it.

Do not declare M2 complete.  
Do not begin M2.2 automatically.  
Do not begin M3.  
Stop and wait for technical-lead review.
```


```
The technical lead has reviewed the M2.1 report.

## VERDICT

M2.1 is **CONDITIONALLY APPROVED PENDING CRYPTOGRAPHIC CODE AUDIT**.

Do NOT begin M3.

Do NOT implement X25519, `crypto/ecdh`, HKDF, ChaCha20-Poly1305, handshake logic, session establishment, or any other M3 functionality.

Your task is to perform a focused security audit of the actual M2.1 Ed25519 implementation you just created.

This is an audit/correction task, not a feature-expansion task.

---

# 1. AUDIT THE ACTUAL IMPLEMENTATION

Inspect:

- `internal/crypto/identity.go`
    
- `internal/crypto/identity_test.go`
    
- `internal/crypto/doc.go`
    

Do not rely on the previous report.

Review the actual source line by line for:

- key generation
    
- private-key storage
    
- public-key derivation
    
- defensive copying
    
- signing
    
- verification
    
- error handling
    
- malformed input handling
    
- nil handling
    
- exported APIs
    
- accidental private-key exposure
    

---

# 2. PRIVATE-KEY EXPOSURE AUDIT

Prove that there is no public API that can accidentally expose or mutate the internal private key.

Check specifically for:

- exported fields
    
- methods returning private-key slices
    
- methods returning references to internal byte arrays
    
- unsafe conversions
    
- shared backing arrays
    
- logging
    
- `%v`, `%x`, `%s`, or similar formatting of private keys
    
- error messages containing key material
    
- protocol serialization of private keys
    
- telemetry containing key material
    

If defensive copies are used, verify that they actually copy the underlying bytes rather than merely returning another slice pointing at the same backing array.

Add regression tests where appropriate.

---

# 3. PUBLIC-KEY DEFENSIVE COPYING

Verify that:

`PublicKey()`

returns a defensive copy.

Create a test demonstrating:

1. obtain public key
    
2. mutate returned byte slice
    
3. obtain public key again
    
4. internal identity public key remains unchanged
    

Do not weaken the identity invariant.

---

# 4. PRIVATE-KEY LOADING

Audit `LoadIdentity(privateKey)`.

Verify:

- exact Ed25519 private-key length is required
    
- malformed private keys are rejected safely
    
- the private key is copied into internal storage
    
- the public key is derived from the private key
    
- independently supplied public-key data cannot create a mismatch
    
- mutation of the caller's input after loading cannot mutate internal key material
    

Add a regression test for caller-side mutation of the supplied private-key slice.

---

# 5. SIGNING SEMANTICS

Verify that signing uses the standard:

`crypto/ed25519`

implementation directly or through a thin, auditable wrapper.

Do not reimplement Ed25519.

Verify:

- signature length is correct
    
- same key + same message gives deterministic Ed25519 signatures
    
- modified message fails verification
    
- wrong public key fails verification
    

Do not alter Ed25519's cryptographic semantics merely to make tests pass.

---

# 6. NIL / EMPTY MESSAGE POLICY

Review the current behavior for nil and empty messages.

Important distinction:

- nil byte slice
    
- zero-length byte slice
    

Ed25519 itself can operate on arbitrary byte messages, including zero-length messages.

If the project intentionally rejects nil/empty messages as an application-level policy, document the exact policy and explain why.

If the implementation is rejecting nil merely because it was assumed to be cryptographically unsafe, correct that assumption.

Do not modify the standard Ed25519 primitive.

The policy must be intentional and documented.

---

# 7. VERIFICATION API

Audit:

`VerifySignature(pubKey, message, signature)`

It must:

- validate public-key length
    
- validate signature length
    
- safely reject malformed input
    
- never panic on attacker-controlled slices
    
- return a clear failure result/error
    
- never log supplied key/signature/message contents
    

Do not add unnecessary custom cryptographic logic around `ed25519.Verify`.

---

# 8. ERROR TYPES

Review the custom errors such as:

- `ErrMalformedPublicKey`
    
- `ErrMalformedSignature`
    
- other identity errors
    

Ensure errors do not contain:

- private keys
    
- public-key material unnecessarily
    
- signatures
    
- message contents
    

Do not create a complicated error hierarchy unless it is justified.

---

# 9. TEST QUALITY AUDIT

Inspect the tests themselves.

Ensure they actually prove the claimed properties instead of merely checking expected return values.

At minimum, verify coverage for:

### Generation

- correct key sizes
    
- public/private correspondence
    
- independent generation produces distinct identities
    

### Signing

- successful verification
    
- deterministic Ed25519 behavior
    
- modified message rejection
    
- wrong identity rejection
    

### Input validation

- malformed public key
    
- malformed signature
    
- malformed private key
    
- caller mutation of loaded private key
    
- mutation of returned public key
    
- nil/empty-message policy
    

### Safety

- no panic on malformed attacker-controlled byte slices
    

Tests must not print secrets.

---

# 10. RACE / CONCURRENCY CONSIDERATION

Determine whether the current `NodeIdentity` API is safe for concurrent signing.

Do not add synchronization unless the implementation actually requires mutable shared state.

If the identity is immutable after construction, document that property.

If necessary, run:

`go test -race ./...`

inside the Docker toolchain.

Do not introduce unnecessary locks.

---

# 11. CRYPTOGRAPHIC LIBRARY POLICY

Confirm that M2.1 uses only:

- `crypto/ed25519`
    
- `crypto/rand`
    

for the Ed25519 identity implementation.

Do not add external cryptographic libraries.

Do not implement cryptographic primitives manually.

---

# 12. M3 MUST REMAIN BLOCKED

Do NOT implement any of the following:

- X25519
    
- `crypto/ecdh`
    
- `curve25519`
    
- ephemeral keys
    
- shared secrets
    
- HKDF
    
- session IDs
    
- T_INIT
    
- T_RESP
    
- signatures over handshake transcripts
    
- session state machine
    

The technical lead will separately authorize M3 after M2.1 is formally approved.

---

# 13. VALIDATION

Run the complete containerized validation:

- go fmt ./...
    
- go vet ./...
    
- go test ./...
    
- go test -race ./...
    
- go build ./...
    

Also verify the existing Protobuf generation workflow remains functional.

Do not use host-installed Go.

---

# 14. DOCUMENTATION

Only update documentation if the audit discovers that the current documentation does not accurately describe the implementation.

Do not rewrite unrelated documentation.

If the nil/empty-message behavior changes, document that policy.

---

# REQUIRED REPORT

Return exactly:

TASK:  
M2.1-B — Ed25519 Cryptographic Implementation Audit

FILES CHANGED:  
List every file changed.

IMPLEMENTATION AUDIT:  
Summarize the actual identity implementation and any issues found.

PRIVATE-KEY AUDIT:  
Explain exactly how private-key exposure and mutation are prevented.

PUBLIC-KEY AUDIT:  
Explain defensive-copy behavior and test evidence.

LOADING AUDIT:  
Explain private-key loading, validation, copying, and public-key derivation.

SIGNING AUDIT:  
Explain the exact Ed25519 API usage.

NIL/EMPTY MESSAGE POLICY:  
State the exact behavior and why it is intentional.

ERROR AUDIT:  
List custom errors and confirm they do not leak sensitive material.

TEST AUDIT:  
List all relevant tests and what security property each proves.

RACE TEST:  
Provide the exact command and result.

VALIDATION:  
Provide exact results for:

- fmt
    
- vet
    
- tests
    
- race tests
    
- build
    
- Protobuf generation
    

DEPENDENCIES:  
List any changes.

DOCUMENTATION:  
List any documentation changes.

SECURITY ISSUES:  
List every issue found, including issues that were corrected.

ARCHITECTURE ISSUES:  
List anything requiring technical-lead review.

KNOWN ISSUES:  
List unresolved issues.

BLOCKERS:  
List blockers.

NEXT RECOMMENDATION:  
Recommend the next task, but DO NOT implement it.

IMPORTANT:  
Do not begin M3.  
Do not implement X25519.  
Do not choose between `crypto/ecdh` and `golang.org/x/crypto/curve25519`.  
Stop after the audit and wait for technical-lead approval.
```


```
The technical lead has reviewed M2.1-B.

## STATUS

M2.1 is **conditionally approved**, pending one final key-ownership correction.

Do NOT begin M3.

Do NOT implement:

- X25519
    
- crypto/ecdh
    
- curve25519
    
- HKDF
    
- session IDs
    
- handshake transcripts
    
- handshake state machines
    
- AEAD
    
- session keys
    
- networking
    
- routing
    

---

# REQUIRED CORRECTION

Review `internal/crypto/identity.go`, specifically `LoadIdentity()`.

The current ownership model must ensure that caller-owned private-key memory is copied BEFORE it is used as the internal source for identity construction.

The intended sequence is:

1. Validate exact private-key length.
    
2. Make a defensive copy of the caller's private-key bytes.
    
3. Construct the internal `ed25519.PrivateKey` from that defensive copy.
    
4. Derive the public key from the internal copied private key.
    
5. Store an appropriately owned/independent public-key representation.
    

Do not derive identity material from caller-owned memory before the defensive copy.

Do not introduce unnecessary cryptographic transformations.

---

# REQUIRED REGRESSION TEST

Ensure the tests demonstrate:

1. Caller creates a private-key byte slice.
    
2. `LoadIdentity()` is called.
    
3. Caller mutates the original private-key slice.
    
4. Identity remains fully functional.
    
5. `PublicKey()` remains unchanged.
    
6. Signing still produces signatures that verify against the original public identity.
    

The test must not print or persist private key material.

---

# PUBLIC-KEY OWNERSHIP

Also verify that the stored public key cannot accidentally share mutable backing memory with:

- the caller's private-key input
    
- a temporary public-key slice
    
- a returned `PublicKey()` result
    

Maintain the existing defensive-copy guarantee.

---

# THREAD-SAFETY WORDING

Do not change implementation merely for this point.

Correct the M2.1-B report/documentation wording if necessary:

`go test -race` demonstrates that the exercised concurrent test cases did not detect data races.

Do NOT claim that a race-detector pass mathematically proves complete thread safety under every possible usage.

If the implementation is immutable after construction and concurrent signing has been tested successfully, state that precisely.

---

# VALIDATION

Run:

- go fmt ./...
    
- go vet ./...
    
- go test ./...
    
- go test -race ./...
    
- go build ./...
    

All must pass inside Docker.

Verify the existing Protobuf generation workflow remains intact.

---

# DOCUMENTATION

Only modify documentation if the implementation or wording currently requires correction.

Do not rewrite unrelated documentation.

---

# REQUIRED REPORT

Return exactly:

TASK:  
M2.1-C — Final Key Ownership Verification

FILES CHANGED:  
List every changed file.

KEY OWNERSHIP:  
Explain the final LoadIdentity ownership sequence.

PRIVATE KEY:  
Explain why caller mutation cannot affect the loaded identity.

PUBLIC KEY:  
Explain the ownership/defensive-copy guarantees.

TESTS:  
List the relevant regression tests and their results.

RACE:  
State precisely what the race-detector test demonstrates.

VALIDATION:  
Report fmt, vet, test, race, build, and Protobuf-generation results.

SECURITY ISSUES:  
List anything corrected.

DOCUMENTATION:  
List any documentation changes.

KNOWN ISSUES:  
List unresolved issues.

BLOCKERS:  
List blockers.

NEXT RECOMMENDATION:  
Recommend the next task but DO NOT execute it.

IMPORTANT:  
Do not begin M3.  
Stop after this correction and wait for technical-lead approval.
```


```
You are implementing Milestone M3.1 of the E2EE Mesh Chat project.

IMPORTANT GOVERNANCE:
- M2 is complete and documentation has been synchronized.
- You are NOT authorized to implement M3.2 or any later milestone.
- Do not self-advance the project.
- Implement ONLY the scope explicitly defined below.
- Do not modify protocol wire format, networking, routing, UI, observability architecture, or session establishment beyond what is necessary to add the isolated X25519 primitive.
- Do not implement HKDF, Ed25519 handshake signatures, canonical handshake transcripts, session IDs, AEAD, message encryption, or network handshake logic in this milestone.

MILESTONE:
M3.1 — X25519 Ephemeral Key Agreement Foundation

OBJECTIVE:
Implement and test a small, auditable X25519 abstraction that can later be used by the M3 handshake implementation.

ARCHITECTURE REQUIREMENTS:

1. CRYPTOGRAPHIC LIBRARY/API
- Inspect the existing go.mod and architecture documentation before implementation.
- Choose the X25519 implementation deliberately and document the exact choice.
- Prefer a well-maintained standard/library implementation rather than implementing Curve25519/X25519 arithmetic manually.
- Do NOT implement cryptographic primitives from scratch.
- If the repository already specifies a library/API, follow that decision unless there is a concrete technical reason it is unsuitable; if you believe the existing decision must change, STOP and report the conflict instead of silently changing the architecture.

2. EPHEMERAL KEY GENERATION
Provide an API for generating a fresh X25519 ephemeral keypair.

Requirements:
- Private key must be generated using a cryptographically secure randomness source.
- Public key must correspond exactly to the private key.
- Generate a fresh keypair for each requested session/key exchange.
- Never reuse the long-term Ed25519 identity key as an X25519 private key.
- Do not derive X25519 private keys from Ed25519 keys.
- Do not persist ephemeral private keys as project state.
- Do not log private keys or public/private key material unnecessarily.

3. KEY SIZES
Enforce/document the X25519 sizes used by the implementation:
- private key: 32 bytes
- public key: 32 bytes
- shared secret: 32 bytes

Use the selected library's canonical representation/API rather than inventing a custom encoding.

4. SHARED SECRET
Provide an API that computes:

    shared_secret = X25519(our_ephemeral_private_key, peer_ephemeral_public_key)

Requirements:
- Return exactly the raw 32-byte X25519 shared secret.
- Do NOT run HKDF inside this API.
- Do NOT hash, encode, truncate, expand, or otherwise transform the shared secret.
- HKDF will be implemented in a later milestone.
- Treat the shared secret as sensitive material.

5. INVALID / LOW-ORDER PUBLIC KEYS
Handle invalid peer public keys according to the guarantees and error semantics of the selected X25519 implementation.

In particular:
- Do not silently accept an invalid/low-order input if the selected API provides a way to detect/reject it.
- Do not substitute an arbitrary shared secret on failure.
- Do not convert cryptographic failure into a successful-looking zero/constant secret.
- Explicitly test the failure behavior.
- Document exactly what the chosen Go API guarantees regarding low-order/all-zero shared secrets.

If the selected API requires explicit all-zero shared-secret rejection, implement that rejection.
If the API already guarantees rejection, document and test that behavior rather than duplicating incompatible logic.

6. KEY OWNERSHIP / MUTABILITY
Follow the same defensive ownership principles established during M2.1-C.

Requirements:
- Do not retain mutable caller-owned byte slices without deliberate ownership semantics.
- Do not expose internal private-key buffers directly.
- Public/private key accessors must return defensive copies where appropriate.
- Clearly document whether key objects are immutable after construction.
- Avoid accidental aliasing between internal key state and caller-provided buffers.

7. PRIVATE-KEY LIFETIME / ZEROIZATION
Ephemeral private keys and shared secrets are sensitive.

Implement reasonable best-effort cleanup semantics where practical in Go.

IMPORTANT:
- Do NOT claim that Go provides guaranteed memory zeroization.
- Do NOT claim that garbage collection makes secrets securely erased.
- Do NOT introduce unsafe memory tricks merely to claim "secure deletion."
- If explicit zeroization is implemented, document its limitations clearly.
- Ensure cleanup does not corrupt memory or cause races.
- The API must make it clear which object owns sensitive key material and when it should no longer be used.

8. ERROR HANDLING
Cryptographic failures must be explicit.

Requirements:
- Invalid key lengths must fail.
- Invalid/nil inputs must not cause panics in normal API usage.
- X25519 computation failures must return an error.
- Do not leak secret material through error messages.
- Do not include private keys, shared secrets, or peer secret material in errors/logs.

9. LOGGING / OBSERVABILITY
No cryptographic secret material may enter:
- application logs
- structured logs
- Prometheus metrics
- tracing
- debug output
- test output committed as artifacts

Do not add crypto secrets to telemetry labels or metric values.

10. TESTS
Add comprehensive tests for the X25519 abstraction.

At minimum test:

A. Key generation
- key generation succeeds
- private key is 32 bytes
- public key is 32 bytes
- two independently generated keypairs are not deterministically identical

B. Public/private consistency
- generated public key corresponds to its private key

C. Shared-secret agreement
Generate Alice and Bob keypairs and verify:

    Alice(priv) × Bob(pub) == Bob(priv) × Alice(pub)

The resulting shared secret must be exactly 32 bytes.

D. Freshness
- repeated key generation produces independent ephemeral keys
- demonstrate that the API does not reuse one static ephemeral keypair

E. Invalid input
- wrong private-key length
- wrong public-key length
- nil inputs where applicable
- malformed inputs must fail cleanly

F. Low-order / invalid peer key behavior
- include known test cases appropriate to the selected implementation
- verify that invalid/low-order inputs are handled according to the API's documented security semantics
- specifically test all-zero shared-secret handling if applicable

G. Ownership
- mutate caller buffers after passing them into the API and verify that internal state is not unexpectedly affected
- mutate returned key material and verify internal state remains protected where the API promises defensive copies

H. Cleanup
- test cleanup behavior if an explicit cleanup API is implemented
- verify cleanup does not panic
- do not write tests that falsely claim cryptographic erasure is guaranteed by Go

I. Concurrency
- run relevant crypto object operations under the race detector
- do not claim race-detector success proves mathematical thread safety; report precisely what was tested

11. TEST VECTORS / INTEROPERABILITY
Where practical, include at least one known-good X25519 test vector from a reputable public specification/reference or use a well-established library test vector.

Do not copy large copyrighted material into the repository.

The test should verify that the selected implementation agrees with an external/reference value, not merely that two calls to the same abstraction agree with each other.

12. API DESIGN
Keep the API small.

Do not expose unnecessary cryptographic internals.

Before implementing, inspect existing internal/crypto conventions and make the new X25519 API consistent with the existing Ed25519 identity implementation.

Do not redesign NodeIdentity.

Do not mix Ed25519 and X25519 into one key type.

13. DOCUMENTATION
Update only documentation directly required for M3.1.

At minimum, update the relevant cryptography documentation to state:
- X25519 is now implemented for ephemeral key agreement
- exact library/API selected
- key sizes
- shared-secret semantics
- invalid/low-order input behavior
- ownership/lifetime expectations
- zeroization limitations
- explicit statement that HKDF/session establishment is NOT yet implemented

Do not mark M3 complete.
Do not mark the authenticated handshake complete.

14. VALIDATION
Run the complete applicable validation in the project's containerized environment.

At minimum:
- gofmt
- go vet
- go test
- go test -race
- go build
- Docker image build
- Docker Compose configuration validation
- any relevant existing project tests

Do not rely only on host tooling if the project is intended to be containerized.

15. SECURITY REVIEW REPORT
When finished, provide a concise implementation report containing:

A. Files added
B. Files modified
C. Exact X25519 library/API selected and why
D. Public API introduced
E. Key generation behavior
F. Shared-secret behavior
G. Low-order/all-zero handling
H. Ownership/defensive-copy behavior
I. Cleanup/zeroization approach and limitations
J. Tests added
K. Validation commands and results
L. Any unresolved security concerns
M. Confirmation that no M3.2 functionality was implemented

IMPORTANT:
Do not report "secure", "production-ready", "unhackable", or equivalent absolute claims.

STOP after M3.1 and report the result to the Project Overseer for review.
```


```
M3.1 SECURITY REVIEW — CORRECTIONS REQUIRED

The Project Overseer has reviewed the M3.1 report.

M3.1 is NOT yet approved. Do not begin M3.2.

Make the following corrections/clarifications only.

1. EPHEMERAL FRESHNESS CLAIM

The current implementation/report states that the absence of a private-key loading API "mathematically enforces freshness."

This claim is incorrect.

Do NOT claim that the API mathematically enforces one-time use merely because callers cannot load an arbitrary private key.

The EphemeralKey abstraction may remain reusable at the primitive level unless there is a strong reason to make it single-use.

Preferred design:
- GenerateEphemeralKey() creates a fresh keypair.
- The object represents one ephemeral keypair.
- The primitive may compute X25519 with multiple peer public keys if the API permits it.
- The FUTURE handshake/session layer is responsible for creating a fresh EphemeralKey for every handshake/session establishment.
- Document this explicitly.
- Do not implement session lifecycle enforcement in M3.1.

2. DESTROY CONCURRENCY / LIFECYCLE

Inspect the actual implementation of:
- PublicKey()
- ComputeSharedSecret()
- Destroy()

Determine whether concurrent calls involving Destroy() can access k.privateKey concurrently.

There must be no unsynchronized read/write race on privateKey.

Choose one clean design:

OPTION A — synchronization:
Use an appropriate synchronization primitive to make lifecycle access safe.

OR

OPTION B — explicit lifecycle contract:
If the object is deliberately NOT concurrency-safe and must not be used concurrently with Destroy(), make that contract explicit and ensure the implementation does not falsely imply otherwise.

Do not use a design that passes ordinary race tests but leaves Destroy() racing with cryptographic operations.

Prefer the simplest auditable solution.

3. DESTROY SECURITY CLAIM

Keep the existing limitation that Destroy() is best-effort only.

Do not claim:
- guaranteed RAM zeroization
- secure erasure
- protection against core dumps
- protection against swap analysis

If Destroy() only removes the application's reference to the opaque *ecdh.PrivateKey, describe it accurately as an application-level lifecycle boundary.

4. ACTUAL LOW-ORDER-POINT SEMANTICS

Inspect the exact Go crypto/ecdh X25519 behavior used by this project.

Verify:
- malformed 32-byte public keys
- invalid public keys
- low-order public keys
- all-zero public key
- all-zero shared-secret behavior, if applicable

Do not merely state that crypto/ecdh "mathematically rejects" these inputs.

Document the exact behavior guaranteed by the selected Go API/version.

If the API returns an error, preserve that error boundary and map it to the project's typed error without leaking secret material.

5. TEST THE LIFECYCLE EDGE CASE

Add a test covering the relevant Destroy/ComputeSharedSecret lifecycle behavior.

If the API is synchronization-safe:
- test concurrent lifecycle operations under `go test -race`.

If the API explicitly forbids concurrent Destroy/use:
- document the contract and test valid sequential lifecycle behavior.
- Do not claim concurrency safety.

6. RFC 7748 TEST VECTOR

Inspect TestRFC7748TestVector carefully.

Confirm that:
- the scalar/private input is interpreted exactly as required by RFC 7748,
- the peer public input is correct,
- the expected output is the RFC reference output,
- the test is genuinely an external known-answer test and not merely a round-trip generated by the same implementation.

7. API REVIEW

Re-review the API for unnecessary complexity.

Current API:
- EphemeralKey
- GenerateEphemeralKey
- PublicKey
- ComputeSharedSecret
- Destroy

Keep it small unless a concrete issue requires change.

Do not add:
- HKDF
- session state
- Ed25519
- handshake transcript
- AEAD
- packet handling
- networking

8. UPDATE DOCUMENTATION

Correct any statement that implies:
"the API mathematically enforces freshness."

Replace it with an accurate lifecycle statement explaining that:
- GenerateEphemeralKey creates a fresh keypair;
- each handshake/session establishment must create a new ephemeral object;
- enforcement of that protocol lifecycle belongs to the future session/handshake layer.

9. RE-RUN VALIDATION

Run:

go fmt ./...
go vet ./...
go test -v ./...
go test -race -v ./...
go build ./...

Also run the relevant Docker-based project validation.

10. REPORT

Return:

A. Files changed
B. Exact changes made
C. Destroy concurrency/lifecycle semantics
D. Exact crypto/ecdh invalid/low-order behavior verified
E. RFC7748 vector verification
F. Tests added/changed
G. Validation results
H. Confirmation that M3.2 remains unimplemented

STOP after these corrections and return the report.

Do not claim M3.1 is approved. The Project Overseer will make the gate decision.

```


```
TASK: M3.2 — Authenticated Session Establishment

PROJECT GOVERNANCE:
- M0, M0.1, M1.1, M2.1, M2.1-B, M2.1-C, M2 documentation synchronization, and M3.1 are COMPLETE.
- M3.1 X25519 ephemeral key agreement is APPROVED.
- You are now authorized to implement ONLY M3.2.
- Do NOT implement M3.3 or any later milestone.
- Do NOT implement AEAD application-message encryption yet.
- Do NOT implement mesh routing changes.
- Do NOT implement application messaging.
- Do NOT modify the frozen wire-format decisions unless a blocking inconsistency is discovered. If one is discovered, STOP and report it rather than silently changing the architecture.

MILESTONE:
M3.2 — Authenticated Session Establishment

OBJECTIVE:
Build the authenticated cryptographic session-establishment layer on top of:
- existing Ed25519 NodeIdentity
- existing M3.1 X25519 EphemeralKey

The result must allow two authenticated peers to derive the same directional session keys from an exact, canonical handshake transcript.

This milestone is about HANDSHAKE AND SESSION KEY ESTABLISHMENT ONLY.

==================================================
1. FROZEN CRYPTOGRAPHIC CONSTRUCTION
==================================================

Do not reinterpret or simplify the following construction.

All concatenation below is raw byte concatenation.

DOMAIN:

    "MeshChat-Handshake-v1"

encoded as ASCII bytes.

PROTOCOL_VERSION:

    0x01

ZERO32:

    exactly 32 zero bytes.

IDENTITIES:

    ID_A = exactly 32-byte Ed25519 public key of initiator
    ID_B = exactly 32-byte Ed25519 public key of responder

EPHEMERALS:

    E_A = exactly 32-byte X25519 public key generated by initiator
    E_B = exactly 32-byte X25519 public key generated by responder

INITIATOR TRANSCRIPT:

    T_INIT =
        DOMAIN
        || u8(PROTOCOL_VERSION)
        || ID_A
        || ID_B
        || E_A
        || ZERO32

RESPONDER TRANSCRIPT:

    T_RESP =
        DOMAIN
        || u8(PROTOCOL_VERSION)
        || ID_A
        || ID_B
        || E_A
        || E_B

SIGNATURES:

    Alice/initiator signs T_INIT using her Ed25519 identity private key.

    Bob/responder signs T_RESP using his Ed25519 identity private key.

Verification:

    Initiator verifies responder signature over T_RESP using ID_B.

    Responder verifies initiator signature over T_INIT using ID_A.

IMPORTANT:
- Do not sign only the ephemeral key.
- Do not sign only the identity.
- Do not sign serialized protobuf messages as a substitute.
- Do not change field ordering.
- Do not use T_RESP for the initiator signature.
- Do not use T_INIT for the responder signature.
- Do not omit ZERO32 from T_INIT.
- Do not add timestamps, packet IDs, TTL, route metadata, or other fields to these transcripts.

SESSION ID:

    session_id = SHA256(T_RESP)

The full 32-byte SHA-256 digest is the session ID.

Do not truncate it.

==================================================
2. X25519 SHARED SECRET
==================================================

Use the approved M3.1 EphemeralKey abstraction.

For initiator:

    SS = X25519(E_A_private, E_B_public)

For responder:

    SS = X25519(E_B_private, E_A_public)

Both sides must obtain the identical 32-byte shared secret.

Do not:
- hash SS before HKDF
- truncate SS
- concatenate SS with other secrets
- use Ed25519 private keys in X25519
- reuse long-term Ed25519 keys as ephemeral DH keys

Every handshake must generate a fresh X25519 ephemeral keypair.

==================================================
3. HKDF
==================================================

Use HKDF-SHA-256.

Do not implement HKDF from scratch.

Construction:

    salt = SHA256(T_RESP)

    PRK = HKDF-Extract(
        salt,
        SS
    )

Initiator-to-responder key:

    K_A_to_B =
        HKDF-Expand(
            PRK,
            "MeshChat-v1|session|initiator->responder",
            32
        )

Responder-to-initiator key:

    K_B_to_A =
        HKDF-Expand(
            PRK,
            "MeshChat-v1|session|responder->initiator",
            32
        )

The labels above are exact ASCII byte strings.

Do not modify capitalization, punctuation, arrows, separators, or spacing.

The two directional keys MUST be distinct.

Do not derive both directions from separate independent DH computations.

==================================================
4. LIBRARY CHOICE
==================================================

Use the existing approved crypto implementations.

X25519:
- existing internal/crypto/EphemeralKey abstraction
- ultimately Go crypto/ecdh.X25519()

Ed25519:
- existing NodeIdentity abstraction
- Go crypto/ed25519

SHA-256:
- Go standard crypto/sha256

HKDF:
- use a maintained implementation.
- Prefer the Go standard-library-supported/approved implementation available to the project.
- Do not implement HKDF manually.

Before adding a dependency, inspect go.mod and existing architecture decisions.

If the selected HKDF API conflicts with the documented architecture, STOP and report the conflict instead of silently changing the architecture.

==================================================
5. SESSION API
==================================================

Design a small, auditable API.

The session-establishment layer should conceptually provide:

INITIATOR:
    create handshake state
    generate ephemeral X25519 key
    construct T_INIT
    sign T_INIT
    produce initiator handshake material

RESPONDER:
    receive initiator material
    validate identities/ephemeral key
    construct T_RESP
    verify initiator signature
    generate fresh responder ephemeral key
    sign T_RESP
    derive session keys

INITIATOR:
    receive responder material
    reconstruct T_RESP
    verify responder signature
    derive session keys

Do not unnecessarily expose:
- private Ed25519 keys
- X25519 private keys
- raw shared secrets
- HKDF PRK

The final session object should expose only what the later authenticated-encryption layer will actually need, such as:
- session ID
- local/remote identity public keys where appropriate
- directional key material through controlled internal access

Avoid returning secret byte slices without clear ownership semantics.

==================================================
6. KEY OWNERSHIP
==================================================

Follow M2.1-C defensive ownership principles and M3.1 lifecycle principles.

Requirements:
- no caller-owned mutable slice may unexpectedly alias internal cryptographic state
- public identity values returned to callers must use defensive copies where appropriate
- session secret material must not be directly exposed unnecessarily
- ephemeral private keys must have a clear lifecycle
- ephemeral keys must not be reused across separate handshake instances
- Destroy ephemeral state after successful or failed handshake when no longer required, according to the existing M3.1 lifecycle contract

Do not claim guaranteed zeroization.

==================================================
7. AUTHENTICATION / TRUST MODEL
==================================================

The existing architecture uses Ed25519 identity keys.

Verification must use the supplied peer identity public key.

Do NOT introduce a new authentication mechanism.

TOFU/pre-shared identity directory remains an architectural trust decision:
- cryptographic signatures authenticate possession of the corresponding Ed25519 private key
- TOFU does NOT magically make first contact resistant to an active MITM
- do not claim first-contact authentication unless a trusted key directory is actually present

Do not implement certificate authorities, passwords, OAuth, or other identity systems in M3.2.

==================================================
8. TRANSCRIPT BINDING SECURITY
==================================================

The implementation must ensure that signatures authenticate the exact intended transcript.

The following changes MUST cause verification failure:

- changing ID_A
- changing ID_B
- changing E_A
- changing E_B
- changing protocol version
- changing DOMAIN
- changing T_INIT/T_RESP field ordering
- replacing ZERO32
- swapping initiator/responder identities
- substituting a different responder ephemeral key
- modifying any signed transcript byte

The session ID must also change when T_RESP changes.

==================================================
9. ROLE / DIRECTION BINDING
==================================================

The directional keys must be role-bound.

Initiator derives:

    send_key = K_A_to_B
    receive_key = K_B_to_A

Responder derives:

    send_key = K_B_to_A
    receive_key = K_A_to_B

Do not simply assign the same key to both directions.

The two derived keys must not be equal.

A reflected message must not automatically become valid in the opposite direction.

Do not implement application-message encryption yet; only establish the directional key separation required by the future AEAD layer.

==================================================
10. SESSION ID
==================================================

Session ID must be:

    SHA256(T_RESP)

exactly 32 bytes.

It must be identical on both sides after successful authentication.

It must change if any byte of T_RESP changes.

Do not use:
- random session IDs
- packet IDs
- timestamps
- identity hashes alone
- ephemeral public keys alone

==================================================
11. FAILURE HANDLING
==================================================

Handshake failure must be explicit and fail closed.

Reject:
- malformed identities
- wrong identity length
- malformed ephemeral keys
- wrong ephemeral key length
- unsupported protocol version
- invalid initiator signature
- invalid responder signature
- transcript mismatch
- X25519 failure
- HKDF failure if applicable
- inconsistent role/state
- unexpected handshake message
- missing required fields

Do not continue to session establishment after authentication failure.

Do not generate session keys for an unauthenticated peer.

Do not return partially established sessions.

Do not expose secret material in errors.

==================================================
12. STATE MACHINE
==================================================

Implement explicit handshake states.

At minimum distinguish:

Initiator:
    NEW
    INIT_SENT
    ESTABLISHED
    FAILED

Responder:
    NEW
    INIT_RECEIVED
    RESP_SENT
    ESTABLISHED
    FAILED

Invalid state transitions must fail.

Examples:
- responder cannot receive RESP before INIT
- initiator cannot process RESP before INIT was sent
- established session cannot be silently reset by an unexpected handshake message
- failed session cannot continue deriving keys

Keep the state machine local to M3.2.

Do not build the networking transport state machine yet.

==================================================
13. REPLAY / DUPLICATE HANDSHAKE CONSIDERATIONS
==================================================

M3.2 is not yet responsible for complete network-level replay protection.

However:
- handshake state must reject unexpected duplicate/out-of-order messages
- a completed handshake must not blindly accept an old RESP as a new successful state transition
- do not introduce timestamps as a substitute for cryptographic transcript binding
- do not claim network-level replay resistance is complete until the later networking/protocol work implements it

Document this boundary clearly.

==================================================
14. PROTOCOL MESSAGE INTEGRATION
==================================================

The existing MeshPacket protobuf schema contains:

    PACKET_TYPE_INIT
    PACKET_TYPE_RESP

and:

    InitPayload {
        bytes ephemeral_key
        bytes signature
    }

    RespPayload {
        bytes ephemeral_key
        bytes signature
    }

M3.2 may add the necessary handshake construction/parsing helpers around these existing payloads.

Do NOT redesign the frozen protobuf schema unless absolutely necessary.

The existing outer packet fields remain conceptually separate from the cryptographic transcript.

The canonical cryptographic transcript MUST NOT be derived implicitly from protobuf serialization.

Construct T_INIT and T_RESP explicitly from their specified fields.

If packet construction is not yet appropriate because transport is not implemented, keep packet integration at a clean helper/API boundary and do not implement networking.

==================================================
15. TESTS — REQUIRED SECURITY CASES
==================================================

Add comprehensive tests.

A. COMPLETE SUCCESSFUL HANDSHAKE
- Alice initiates
- Bob responds
- both authenticate each other
- both derive the same session ID
- both derive matching directional keys

B. KEY DIRECTION
Verify:

    Alice.send == Bob.receive
    Alice.receive == Bob.send
    Alice.send != Alice.receive

C. TRANSCRIPT TAMPERING
For every field below, modify exactly one value and verify signature verification fails:

- protocol version
- ID_A
- ID_B
- E_A
- E_B
- ZERO32
- DOMAIN

D. IDENTITY SWAP
- swap initiator/responder identities
- verify authentication fails

E. EPHEMERAL KEY SUBSTITUTION
- replace E_A
- replace E_B
- verify authentication/session establishment fails

F. SIGNATURE SUBSTITUTION
- replace signature with random bytes
- verify failure
- use valid signature from a different identity
- verify failure

G. SESSION ID
- verify session_id == SHA256(T_RESP)
- modify one T_RESP byte
- verify session ID changes

H. KEY SEPARATION
- verify K_A_to_B != K_B_to_A
- verify changing transcript input causes derived keys to change

I. ROLE CONFUSION / REFLECTION
- attempt to process initiator material as responder material and vice versa
- verify state/authentication failure
- ensure directional labels prevent key-role confusion

J. WRONG PROTOCOL VERSION
- verify handshake rejects unsupported version

K. MALFORMED LENGTHS
Test malformed:
- identity
- ephemeral key
- signature
- session-related values

L. STATE MACHINE
Test invalid ordering:
- RESP before INIT
- duplicate RESP
- processing messages after FAILED
- invalid transitions after ESTABLISHED

M. FRESH EPHEMERALS
Perform two independent handshakes between the same identities and verify:
- ephemeral public keys differ with overwhelming probability
- session IDs differ
- directional session keys differ

N. REPLAY-LIKE OLD RESPONSE
Attempt to use an old responder response with a new initiator handshake state.
Verify authentication/session establishment fails.

O. KNOWN-ANSWER / CROSS-CHECK
Where practical, verify the exact transcript hash/session ID and HKDF outputs using independently computed expected values rather than relying exclusively on the implementation under test.

Do not copy large copyrighted test vectors.

==================================================
16. NEGATIVE CRYPTOGRAPHIC TESTING
==================================================

Specifically ensure that:
- no session object is returned after signature failure
- no session keys are exposed after handshake failure
- failed X25519 computation cannot continue into HKDF
- invalid peer public keys fail closed
- malformed signatures fail closed

==================================================
17. SECRET MATERIAL AND OBSERVABILITY
==================================================

Never log:
- Ed25519 private keys
- X25519 private keys
- X25519 shared secret
- HKDF PRK
- directional session keys

Do not place secret material into:
- Prometheus labels
- metrics
- structured logs
- tracing
- error strings
- debug output

Telemetry architecture must remain out-of-band.

==================================================
18. CONCURRENCY
==================================================

Do not claim the handshake state machine is concurrency-safe unless it actually is.

If handshake objects are intentionally single-threaded:
- document the lifecycle contract.

If synchronization is introduced:
- test it under `go test -race`.

Avoid unnecessary synchronization complexity.

==================================================
19. DOCUMENTATION
==================================================

Update the relevant documentation required by M3.2.

At minimum ensure the cryptography/session-establishment documentation records:

- M3.2 implemented status
- exact T_INIT
- exact T_RESP
- signature roles
- SHA256(T_RESP) session ID
- X25519 shared-secret input
- HKDF-SHA-256 construction
- exact directional labels
- role-to-key mapping
- authentication/trust boundary
- handshake state machine
- failure semantics
- replay boundary
- ephemeral lifecycle
- zeroization limitations

Do not mark AEAD/application-message encryption complete.

Do not mark the entire E2EE messaging system complete.

==================================================
20. VALIDATION
==================================================

Run all relevant containerized validation:

    gofmt / gofmt -w as appropriate
    go vet ./...
    go test -v ./...
    go test -race -v ./...
    go build ./...

Also verify:
- protobuf generation remains intact
- existing M0-M3.1 tests remain green
- Docker build remains green
- Docker Compose configuration remains valid

If a failure is caused by pre-existing unrelated infrastructure, report it explicitly rather than hiding it.

==================================================
21. CODE REVIEW REQUIREMENT
==================================================

Before reporting completion, inspect the actual generated/implemented code for:

- transcript byte construction
- signature input
- identity binding
- role binding
- HKDF labels
- session ID calculation
- key ownership
- state transitions
- failure paths

Do not rely solely on tests.

==================================================
22. REPORT FORMAT
==================================================

When complete, report:

A. Files added
B. Files modified
C. Exact public/internal APIs introduced
D. Exact handshake construction
E. Exact signature construction/verification
F. Exact session-ID construction
G. Exact HKDF construction
H. Directional key mapping
I. Authentication/trust semantics
J. State machine
K. Error/failure behavior
L. Tests added
M. Security-negative tests
N. Validation commands/results
O. Documentation changes
P. Any unresolved security concerns
Q. Explicit confirmation that:
   - AEAD is NOT implemented
   - application messaging encryption is NOT implemented
   - networking handshake is NOT implemented
   - M3.3+ is NOT implemented

STOP after M3.2.

Do not claim the milestone is approved.
Return the report to the Project Overseer for independent security review.

```

```
M3.2 ARCHITECTURE CONFLICT RESOLUTION — AUTHORIZATION GRANTED

The Project Overseer has reviewed the reported HKDF architecture conflict.

DECISION:
Choose OPTION A.

You are explicitly authorized to add:

    golang.org/x/crypto/hkdf

to go.mod and use its maintained HKDF implementation.

DO NOT implement HKDF manually.

==================================================
ARCHITECTURE AMENDMENT
==================================================

The previous wording:

    "No external cryptographic libraries; strictly standard library (crypto/*)."

was too restrictive because the required HKDF primitive is not provided by the Go standard library.

The governing cryptography rule is now:

    "Use Go's standard library cryptographic primitives
    (crypto/*) and approved Go cryptographic sub-repositories
    such as golang.org/x/crypto/* where the standard library
    does not provide the required primitive. Do not implement
    cryptographic primitives manually. Avoid arbitrary or
    unnecessary
    third-party cryptographic dependencies."

For the currently approved architecture:

    Ed25519       -> crypto/ed25519
    X25519        -> crypto/ecdh
    SHA-256       -> crypto/sha256
    HKDF-SHA-256  -> golang.org/x/crypto/hkdf

Future ChaCha20-Poly1305 may use:

    golang.org/x/crypto/chacha20poly1305

but that is NOT authorization to implement M3.3 now.

==================================================
DEPENDENCY REQUIREMENTS
==================================================

Add the smallest appropriate version of:

    golang.org/x/crypto

required by the selected hkdf API.

Do not add unrelated dependencies.

Record the dependency in go.mod/go.sum.

Verify the dependency is reproducible in the containerized build.

Do not use a random third-party HKDF package.

==================================================
M3.2 RESUMPTION
==================================================

Resume M3.2 implementation from the exact task specification previously provided.

In particular, implement:

    salt = SHA256(T_RESP)

    PRK = HKDF-Extract(
        salt,
        SS
    )

    K_A_to_B =
        HKDF-Expand(
            PRK,
            "MeshChat-v1|session|initiator->responder",
            32
        )

    K_B_to_A =
        HKDF-Expand(
            PRK,
            "MeshChat-v1|session|responder->initiator",
            32
        )

Do not alter any transcript, label, ordering, or length.

==================================================
IMPORTANT SCOPE
==================================================

This authorization ONLY resolves the HKDF dependency conflict.

It does NOT authorize:
- M3.3
- AEAD
- application message encryption
- mesh networking
- routing changes
- session rotation
- application messaging
- UI changes

Continue with M3.2 only.

==================================================
DOCUMENTATION
==================================================

Update the relevant architecture/cryptography documentation to reflect
the amended cryptographic dependency policy.

Explicitly state that:
- Go standard crypto/* is preferred;
- approved golang.org/x/crypto/* modules may be used when required;
- arbitrary third-party cryptographic libraries are not permitted;
- cryptographic primitives must not be implemented manually.

Record HKDF as:

    golang.org/x/crypto/hkdf

Do not claim that x/crypto is part of the standard library.

==================================================
VALIDATION
==================================================

After implementation run:

    gofmt
    go vet ./...
    go test -v ./...
    go test -race -v ./...
    go build ./...

and the project's containerized validation.

Verify:
- go.mod is reproducible
- go.sum is present/consistent
- existing M0-M3.1 tests remain green
- protobuf generation remains intact
- Docker build remains green
- Docker Compose configuration remains valid

==================================================
REPORT
==================================================

Return the normal M3.2 implementation report.

Additionally report:

    - exact golang.org/x/crypto version selected
    - why that version/API was selected
    - exact files modified because of the architecture amendment
    - confirmation that HKDF was NOT manually implemented
    - confirmation that no M3.3 functionality was implemented

STOP after M3.2.

Do not claim M3.2 is approved.
Return the implementation for independent Project Overseer security review.

```


```

M3.2 SECURITY REVIEW — SOURCE AUDIT REQUIRED

The Project Overseer has completed the first review of the M3.2 report.

M3.2 is NOT yet approved.

Do not modify the implementation unless required by the additional test requirements below.

The Project Overseer requires the actual source for independent security audit.

Provide the COMPLETE contents of:

    internal/crypto/session.go
    internal/crypto/session_test.go
    go.mod
    internal/crypto/doc.go

Do not provide excerpts or summaries.

Additionally, audit the existing M3.2 tests and add any missing tests from the following list:

1. SESSION ID KNOWN CONSTRUCTION
Verify explicitly:

    session_id == SHA256(T_RESP)

using the exact bytes constructed by the implementation.

2. SESSION ID TRANSCRIPT SENSITIVITY
Modify exactly one byte of T_RESP and verify:

    SHA256(T_RESP_modified) != original_session_id

3. HKDF KNOWN-ANSWER
Add an independently calculated known-answer test for the M3.2 HKDF construction.

The test must verify the exact:
- salt
- PRK
- K_A_to_B
- K_B_to_A

Do not merely verify that Alice and Bob derive equal keys.

The expected values must be independently established rather than generated by calling the same session implementation being tested.

4. TWO-INDEPENDENT-HANDSHAKES TEST
Perform two complete handshakes using the same identities.

Verify:
- fresh ephemeral keys are generated
- session IDs differ
- initiator->responder keys differ
- responder->initiator keys differ

5. MALFORMED INPUT TESTS
Explicitly test:
- wrong ID_A length
- wrong ID_B length
- wrong E_A length
- wrong E_B length
- wrong signature length
- unsupported protocol version

All must fail closed.

6. DIFFERENT-IDENTITY SIGNATURE
Create a valid Ed25519 signature using an unrelated identity and present it where the expected peer signature is required.

Verify authentication fails.

7. OLD RESPONSE / NEW INIT TEST
Perform:
- handshake A -> response A
- create a new initiator handshake with fresh E_A
- attempt to process old response A

Verify the new handshake does NOT establish a session.

8. DUPLICATE / INVALID STATE TESTS
Explicitly test:
- RESP before INIT
- duplicate RESP
- processing after FAILED
- processing after ESTABLISHED
- duplicate/invalid handshake progression

9. SECRET EXPOSURE REVIEW
Inspect the source for:
- raw shared secret returned unnecessarily
- PRK exposed publicly
- session keys exposed through mutable internal slices
- secret material in error strings
- secret material in logs
- secret material in tests/output

10. DO NOT CHANGE THE CRYPTOGRAPHIC CONSTRUCTION

The following remain frozen:

    T_INIT =
        DOMAIN
        || 0x01
        || ID_A
        || ID_B
        || E_A
        || ZERO32

    T_RESP =
        DOMAIN
        || 0x01
        || ID_A
        || ID_B
        || E_A
        || E_B

    session_id = SHA256(T_RESP)

    salt = SHA256(T_RESP)

    PRK = HKDF-Extract(salt, SS)

    K_A_to_B =
        HKDF-Expand(
            PRK,
            "MeshChat-v1|session|initiator->responder",
            32
        )

    K_B_to_A =
        HKDF-Expand(
            PRK,
            "MeshChat-v1|session|responder->initiator",
            32
        )

11. VALIDATION

After any test additions run:

    gofmt
    go vet ./...
    go test -v ./...
    go test -race -v ./...
    go build ./...

Use the containerized project environment.

12. REPORT

Return:
- complete source files requested above
- tests added
- exact known-answer values
- validation results
- any discovered issues

STOP.

Do not proceed to M3.3.
Do not claim M3.2 approval.

```

```
M3.2 has been reviewed by the Project Overseer and is APPROVED.

Proceed to M3.3 ONLY.

Before writing code, inspect the current frozen protocol and cryptographic documentation and preserve all existing M0–M3.2 decisions.

M3.3 objective:  
Define and implement the authenticated-encryption message-layer construction using the already-established directional session keys.

Requirements:

1. Use ChaCha20-Poly1305 through the approved Go cryptographic dependency policy. Do NOT implement any cryptographic primitive manually.
    
2. Freeze and document the exact nonce construction before implementation. It MUST guarantee that a nonce is never reused with the same directional key.
    
3. Bind the following authenticated data exactly as specified by the protocol design:
    
    - protocol/version context
        
    - session_id
        
    - direction/context
        
    - sequence number
        
    - any other fields explicitly frozen by the protocol  
        Do not invent fields silently.
        
4. Use the existing Session object and directional keys. Do not expose raw private keys, X25519 shared secrets, PRKs, or internal mutable secret state.
    
5. Implement:
    
    - EncryptMessage(...)
        
    - DecryptMessage(...)  
        or an equivalent minimal API with the same security properties.  
        Keep the API small and fail-closed.
        
6. Sequence handling:
    
    - reject malformed sequence values
        
    - reject replayed sequence numbers according to the frozen replay model
        
    - do not rely solely on timestamps
        
    - make the receive-side replay state explicit
        
    - do not claim network-wide replay protection yet
        
7. Nonce handling:
    
    - explicitly document ownership and uniqueness guarantees
        
    - test boundary values
        
    - test that a nonce cannot accidentally repeat for the same key
        
    - test that opposite directions cannot collide through the construction
        
8. Negative/security tests MUST include at minimum:
    
    - ciphertext modification
        
    - AAD modification
        
    - wrong session ID
        
    - wrong direction
        
    - wrong key
        
    - wrong sequence number
        
    - truncated ciphertext/tag
        
    - malformed nonce
        
    - replayed message
        
    - cross-session ciphertext
        
    - cross-direction ciphertext
        
    - sequence boundary behavior
        
9. Add an independently computed known-answer test for the complete AEAD construction, including nonce, AAD, plaintext, ciphertext, and authentication tag.
    
10. Verify that plaintext, keys, nonces, and authentication failures are never logged or emitted through observability.
    
11. Run:
    

- gofmt
    
- go vet ./...
    
- go test -v ./...
    
- go test -race -v ./...
    
- go build ./...  
    inside the existing containerized Go 1.21 environment.
    

12. Update only the documentation directly affected by M3.3.
    
13. Do NOT begin M4, networking, mesh routing, UI, or observability implementation.
    
14. At the end, stop and report:
    

- exact source files changed
    
- complete relevant source files
    
- exact nonce construction
    
- exact AAD construction
    
- exact sequence/replay rules
    
- known-answer values
    
- negative tests
    
- validation results
    
- any unresolved security concerns
    

Important:  
The Project Overseer will independently audit the source. Do not report “secure” or “replay protected” in absolute terms. State precisely what the implementation guarantees and what remains outside M3.3 scope.

```


```
M3.3 is CONDITIONALLY APPROVED. Implement it, but freeze the following corrections before coding.

1. ChaCha20-Poly1305:
    
    - Use golang.org/x/crypto/chacha20poly1305.
        
    - Never implement cryptographic primitives manually.
        
    - Use the existing 32-byte directional session keys.
        
2. Nonce construction is FROZEN:  
    nonce = 4 zero bytes || uint64(sequence_number) encoded big-endian.  
    Total nonce length = 12 bytes.
    
    Sequence numbering:
    
    - First outbound sequence = 0.
        
    - Increment by exactly 1 after successful allocation of the current sequence.
        
    - NEVER wrap from math.MaxUint64 to 0.
        
    - If the sender has exhausted the uint64 sequence space, return a permanent sequence-exhausted error.
        
    - Nonce uniqueness depends on this no-wrap rule plus fresh session keys.
        
3. Direction marker is FROZEN:
    
    - 0x01 = Initiator -> Responder
        
    - 0x00 = Responder -> Initiator
        
    
    EncryptMessage must derive its direction internally from the Session role.  
    DecryptMessage must expect the opposite direction internally.  
    The caller must not supply the direction marker.
    
4. AAD is FROZEN:  
    ASCII bytes:  
    "MeshChat-AppData-v1"  
    followed by:  
    session_id (32 bytes)  
    sequence_number (8-byte big-endian)  
    direction_marker (1 byte)
    
    Do not silently add, remove, reorder, or reinterpret fields.
    
5. Replay window is FROZEN:
    
    - 64-message sliding window.
        
    - highest_received = greatest successfully authenticated sequence.
        
    - bitmap bit 0 represents highest_received.
        
    - bit N represents highest_received - N.
        
    - Therefore the currently accepted historical window is:  
        highest_received-63 through highest_received.
        
    - A sequence < highest_received-63 is too old.
        
    - A sequence already represented by a set bitmap bit is a replay.
        
    - A sequence greater than highest_received advances the window.
        
    - Handle startup explicitly because no highest_received exists before the first authenticated message.
        
    - Do not use uint64 arithmetic that can underflow at low sequence numbers.
        
6. CRITICAL AUTHENTICATION ORDER:  
    Never mutate replay state before AEAD authentication succeeds.
    
    Correct DecryptMessage order:  
    a. Validate ciphertext minimum length.  
    b. Validate sequence against replay-window policy without mutating state.  
    c. Construct nonce and AAD.  
    d. Call AEAD.Open().  
    e. If authentication fails, return error and DO NOT mutate replay state.  
    f. If authentication succeeds, update highest_received/replay bitmap.  
    g. Return plaintext.
    
    An unauthenticated packet must never consume a sequence number.
    
7. Concurrency:  
    Session methods must be safe for concurrent use.  
    Prefer:
    
    - sendMu protecting send sequence state
        
    - receiveMu protecting highest_received and replay bitmap
        
    
    Do not serialize sending and receiving through one unnecessary global mutex.
    
8. Minimal API:  
    EncryptMessage(plaintext []byte) (seq uint64, ciphertext []byte, err error)  
    DecryptMessage(seq uint64, ciphertext []byte) ([]byte, error)
    
9. Defensive security behavior:
    
    - Do not expose session keys.
        
    - Do not log plaintext, ciphertext contents, keys, nonces, or authentication secrets.
        
    - Authentication failures must use bounded/generic errors and must not leak secret material.
        
    - Do not claim network-wide replay protection. This is per-session message-layer replay protection.
        
10. Required negative tests:
    

- ciphertext modification
    
- tag modification
    
- wrong session ID in AAD
    
- wrong sequence in AAD
    
- wrong direction
    
- wrong key
    
- wrong protocol context
    
- truncated ciphertext
    
- duplicate/replayed packet
    
- packet outside the 64-message window
    
- valid out-of-order packet inside the window
    
- bitmap advancement
    
- authenticated packet followed by duplicate
    
- unauthenticated packet MUST NOT consume replay-window state
    
- cross-session ciphertext
    
- cross-direction ciphertext
    
- sequence 0
    
- low sequence-number boundary without underflow
    
- math.MaxUint64 exhaustion behavior
    
- concurrent EncryptMessage calls
    
- concurrent DecryptMessage calls
    
- concurrent send and receive operations
    

11. Known-answer test:  
    Produce an independently generated test vector containing:
    

- session key
    
- session ID
    
- sequence number
    
- direction
    
- exact nonce
    
- exact AAD
    
- plaintext
    
- ciphertext including authentication tag
    

Verify the implementation against those values.

12. Documentation:  
    Freeze the exact nonce, AAD, direction, sequence, and replay-window constructions in the appropriate protocol/cryptography documentation.
    
13. Run the complete validation inside the existing containerized Go 1.21 environment:
    

- gofmt
    
- go vet ./...
    
- go test -v ./...
    
- go test -race -v ./...
    
- go build ./...
    

14. Do NOT proceed to networking, mesh routing, UI, or M4.
    

At completion, stop and provide:

- complete changed source files
    
- exact nonce construction
    
- exact AAD bytes/layout
    
- exact sequence rules
    
- exact replay-window algorithm
    
- exact direction semantics
    
- known-answer vector
    
- all security/negative tests
    
- validation output
    
- unresolved security concerns
    

The Project Overseer will independently audit the implementation before M3.3 is approved.

```



```


```