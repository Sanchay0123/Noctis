
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
