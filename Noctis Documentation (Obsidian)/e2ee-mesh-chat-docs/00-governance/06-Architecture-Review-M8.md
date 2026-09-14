# Architecture & Security Review — M8

## Milestone

**M8 — Observability and UI**

## Review authority

Project Overseer / Security Architect / Technical Lead

## Current status

**Phase 1 architecture: 🟢 APPROVED**  
**M8.1 application-service foundation: 🟢 COMPLETE / VERIFIED / APPROVED**  
**M8.2 Fyne GUI foundation: 🟢 COMPLETE / VERIFIED / APPROVED**

M8 remains open because later M8 work, including concrete observability
metrics/dashboard work, has not yet been implemented or approved.

## Approved M8 architecture boundary

M8 introduces an application-facing service layer between the UI and the
existing cryptographic, session, mesh and routing subsystems. The GUI must not
access raw crypto/session/routing internals.

The approved application boundary exposes, at minimum:

- local identity as a display-safe identifier
- explicit peer dialing
- conversation start/listing by authenticated PeerID
- message sending by PeerID and plaintext
- lifecycle start/stop
- application event subscription

A **PeerID** is the hex-encoded Ed25519 public identity. A network address is a
transport locator and is not interchangeable with a PeerID. A conversation is
an application-level view keyed by the authenticated remote PeerID; it is not
a routing path or a transport connection.

Manual peer entry is the approved initial discovery model. DHT, automatic
peer discovery and decentralized discovery protocols are not introduced by
M8.1/M8.2.

## M8.1 implementation boundary

M8.1 implemented `internal/app` and the narrow routing/application boundary
needed to deliver opaque application data upward without violating relay
blindness.

The routing layer does **not** decrypt APP_DATA. It forwards opaque session
metadata and ciphertext to the application service. The application service
looks up the endpoint session and performs `DecryptMessage`; only after
successful AEAD authentication is plaintext emitted as an application event.

The session manager maintains PeerID-to-session lookup for application
message delivery. Removal is conditional so a stale session cannot delete a
newer mapping for the same PeerID.

## M8.2 implementation boundary

M8.2 adds the primary Fyne GUI under a `gui` build tag, keeping the headless
backend build separate from GUI/CGO dependencies.

The GUI provides the approved foundation for:

- conversation selection
- message history/input
- connection/status presentation
- manual peer addition/dialing
- display of the local identity

The GUI receives application-level events and calls the application service;
it does not directly access cryptographic keys, session keys, routing state,
raw protobuf packets or raw sockets.

Mutable GUI state is encapsulated behind synchronized `UIState`. GUI shutdown
signals the event consumer and then stops the application service.

## Security constraints

- Routing remains blind to application plaintext.
- No private keys, ephemeral private keys, session keys, AEAD nonces or other
  secret cryptographic material may cross into the GUI.
- Observability remains out-of-band and must not become a dependency of the
  E2EE path.
- Future concrete metrics must use fixed, bounded label vocabularies; never
  use attacker-controlled PeerIDs as arbitrary metric labels.
- The project AEAD is **ChaCha20-Poly1305**; AES-GCM is not part of the design.
- No anonymity, complete metadata hiding or endpoint compromise resistance is
  claimed.

## Evidence qualification

M8.1 passed the reported formatting, vetting, tests, race-detector and build
validation, including the routing/application plaintext-boundary regression.

M8.2 added application-service and GUI-state concurrency tests and passed the
reported formatting, vetting, tests, race-detector and build validation.
The GUI runtime was not available for interactive smoke testing in the review
environment. The backend Docker E2E latest rerun was externally blocked by
proxy/DNS timeout; previously approved M8.1 Docker E2E evidence remains the
authoritative runtime evidence and no code regression was established.

## Gate decision

M8.1 and M8.2 are approved within their reviewed scopes. **M8.3 must not
begin until separately authorized.**
