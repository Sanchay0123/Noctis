# Architecture & Security Review — M8

## Milestone

**M8 — Application Integration, Observability/UI Foundation and Functional GUI**

## Review authority

Project Overseer / Security Architect / Technical Lead

## Status

**🟢 M8.3 COMPLETE / VERIFIED / APPROVED**

M8.1, M8.2 and M8.3 were implemented incrementally and independently reviewed.
M8 extends the approved M6/M7 backend with an application-facing service boundary
and a Fyne GUI while preserving the endpoint E2EE and relay-blind routing model.

## Scope

M8 covers:

- application/service boundary above the secure mesh backend
- conversation-oriented application state
- endpoint decryption outside the routing layer
- Fyne GUI foundation and functional chat workflow
- explicit connection and secure-session status events
- manual peer connection
- local in-memory message presentation
- GUI concurrency and shutdown lifecycle
- application-layer multi-hop integration evidence

Concrete Prometheus/Grafana exporter implementation remains deferred; M7's
telemetry hooks and security/cardinality constraints remain authoritative.

## M8.1 — Application Service Boundary

M8.1 introduced `internal/app` as the boundary between presentation and the
existing crypto/session/mesh/routing implementation.

The application API includes identity access, peer dialing, conversation
creation, message sending, service lifecycle and event subscription.

A critical M8.1 audit finding was that routing initially decrypted APP_DATA and
passed plaintext upward. This was remediated before approval. The final design
is:

```text
remote ciphertext
      ↓
internal/routing
      ↓ opaque APP_DATA
internal/app
      ↓ session lookup
session.DecryptMessage()
      ↓
Application event
      ↓
GUI
```

`internal/routing` therefore remains blind to application plaintext.

M8.1 also corrected stale peer-to-session mappings during removal.

**M8.1: 🟢 APPROVED**

## M8.2 — Fyne GUI Foundation

M8.2 introduced a Fyne GUI isolated behind the application boundary. The GUI
uses manual peer entry rather than automatic discovery or DHT mechanisms.

The GUI does not receive private keys, session keys, nonces, sequence state,
raw sockets, raw packets or crypto/session/routing internals.

Mutable GUI state is protected by `sync.RWMutex`. GUI lifecycle and event
consumption were tested at the application/state level.

The GUI build requires native graphics dependencies and CGO on the target
environment. The approved evidence did not establish visual runtime in the
headless review environment.

**M8.2: 🟢 APPROVED**

## M8.3 — Functional GUI Integration

M8.3 turned the GUI foundation into a functional chat client.

Implemented behavior includes:

- manual peer connection through `ApplicationService.DialNode`
- explicit `StartConversation` after successful dialing
- application-level session-established events
- distinction between TCP connectivity and authenticated secure-session state
- public PeerID display and copy operation
- structured in-memory messages with sender and timestamp fields
- conversation association by authenticated PeerID
- message sending through `ApplicationService.SendMessage`
- received plaintext delivery through `SubscribeEvents`
- synchronized GUI state
- deterministic GUI event-consumer shutdown using `sync.WaitGroup`

The Alice→Bob→Carol integration test uses real loopback TCP, real PeerManager,
real routing/forwarding, real session establishment and real
ChaCha20-Poly1305 encryption/decryption. The test is conservatively classified
as **APPLICATION-LAYER INTEGRATION** because all three services run inside one
test process. Remediation 4 separately verified that normal inbound transport
admission no longer depends on `ExpectInbound` pre-authorization.

## Security boundary

The GUI remains a thin presentation client:

```text
Fyne GUI
   ↓
ApplicationService
   ↓
Session / Mesh / Routing / Crypto
```

The GUI does not implement encryption, decryption, nonce generation, sequence
management, key derivation, handshake signing or X25519 operations.

Routing remains opaque to application plaintext. Source-level verification
reported no `DecryptMessage` call under `internal/routing` and no direct GUI
imports of `internal/crypto`, `internal/session`, `internal/mesh` or
`internal/routing`.

## Identity boundary

The GUI renders and copies the string returned by `GetLocalIdentity()`.
`GetLocalIdentityBytes()` used by the application for public-identity handling
returns the local Ed25519 public key only; no private key material is exposed.

## Lifecycle boundary

M8.3 remediated the GUI event-consumer shutdown lifecycle. The event-consumer
worker is tracked by a `sync.WaitGroup`; GUI shutdown signals the worker and
waits for its termination before `ApplicationService.Stop()` executes.

This replaces reliance on process termination as the lifecycle guarantee.

## Validation evidence

M8.3 implementation verification passed:

```text
go test ./...
go test -race -p 1 ./...
go vet ./...
go build ./...
```

The M8.3 Alice→Bob→Carol application integration test passed. Remediation 4
added coverage for accepting an authenticated but previously unknown inbound
peer and for preserving expected-PeerID verification after transport admission
modernization.

GUI build verification passed:

```text
go build -tags gui ./cmd/meshchat-gui
```

Docker Compose regression remained unverified because dependency resolution in
the constrained environment was blocked by external proxy/DNS limitations.

### Final live GUI acceptance

The final two-node GUI acceptance passed after Remediation 4. Clean Alice and
Bob instances were launched on `127.0.0.1:8000` and `127.0.0.1:8001`.

Acceptance results:

| Scenario | Result |
|---|---|
| Alice initiates first connection to Bob | PASS |
| Bob initiates first connection to Alice | PASS |
| Authenticated secure session establishment | PASS |
| Alice → Bob encrypted message | PASS |
| Bob → Alice encrypted message | PASS |
| Duplicate Add Peer on established connection | PASS |
| Existing conversation remains intact during duplicate attempt | PASS |

The deterministic first-initiator failure identified during live acceptance was
resolved by removing the transport-layer `ExpectInbound` pre-authorization
dependency.

## Remediation 4 — Transport Admission Modernization

The M8.3 live failure exposed an architectural mismatch between GUI manual
dialing and M5 inbound pre-authorization. The obsolete `ExpectInbound`,
`expectedInbound` and `isExpectedInbound` mechanism was removed.

The resulting admission model is:

```text
TCP accept
    ↓
bounded pending-handshake admission
    ↓
frame/protocol validation
    ↓
M3 authenticated handshake
    ↓
authenticated identity
    ↓
duplicate arbitration
    ↓
established peer/session
```

Transport admission is no longer dependent on out-of-band pre-authorization.
Inbound resource exposure remains bounded by the existing pending-handshake
limit, handshake timeout, frame-size validation and peer limits. Cryptographic
authentication and expected-PeerID verification remain enforced.

This remediation changes transport admission behavior but does not modify the
frozen Protobuf schema, M3 handshake transcript, cryptographic construction,
AEAD, replay protection or M6 routing semantics.

The security review explicitly records that removal of pre-authorization
increases exposure to unauthenticated handshake attempts; it does not provide
a claim of immunity from flooding or DoS. Existing controls bound concurrent
pending handshake work and incomplete-connection lifetime.

**Remediation 4: 🟢 APPROVED**

## Remaining limitations

M8 does not provide:

- automatic peer discovery or DHT
- persistent message storage
- anonymity or complete metadata hiding
- endpoint compromise resistance
- guaranteed delivery
- a concrete Prometheus/Grafana exporter

## Gate decision

**🟢 M8.3 APPROVED**

M8 is formally closed at the functional GUI integration gate. M9 may proceed
only after a separate milestone authorization.
