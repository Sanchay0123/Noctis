# Results Structure

> **Staged results document — final results remain incomplete.**
>
> M2–M7 provide implementation and security evidence for identity, authenticated
> session/message protection, direct networking and bounded multi-hop routing.
> M8 provides application-boundary and functional GUI integration evidence.
> M9 provides relay-blind E2EE, adversarial security-demonstration and live GUI evidence.
> Clean-environment deployment reproduction remains incomplete.

## Functional results

-   node startup
-   peer discovery
-   connection establishment
-   direct messaging
-   multi-hop messaging
-   reconnection

## Security results

  Property              Test               Result       Evidence
  --------------------- ------------------ -----------  ----------------
  Ed25519 identity       Identity test suite Verified     M2 test evidence
  Confidentiality        Transport interception  Verified   M4 E2E evidence
  Integrity              Tampering               Verified   M4 tamper evidence
  Replay resistance      Replay injection        Verified*  M3 session evidence
  Authentication         Invalid identity       Verified   M3/M4 evidence
  MITM resistance        Key substitution    TBD          M3/M9 authentication evidence; full active-MITM evaluation not performed
  Routing containment    TTL/loop             Verified     M6/M9 TTL, PacketID and relay evidence

## Performance results

Only measure if performance is a project requirement.

Possible metrics:

-   message latency
-   throughput
-   connection establishment time
-   routing overhead
-   memory usage

Avoid presenting measurements without a reproducible method.


## M3 preliminary results

The M3 implementation provides evidence for authenticated session
establishment and application-message protection.

## M4 preliminary results

M4 provides evidence that the approved M3 cryptographic/session layer can be
integrated with bounded direct TCP transport and the frozen Protobuf envelope.
The direct path successfully exchanges encrypted application messages in both
directions, rejects malformed/unknown-field packets, enforces session and
identity binding, handles short writes, and rejects ciphertext tampering before
application plaintext delivery.

These are direct secure messaging results, not evidence of complete mesh
security. Multi-hop routing, malicious-relay routing demonstrations, runtime
network hardening, performance measurements and final deployment evidence
remain for later milestones.

*Replay resistance remains a property of the approved M3 session layer; M4
transports the sequence number and ciphertext through the direct packet path.


## M5 preliminary results

M5 demonstrates a functioning direct networking runtime above the approved M4
secure channel. The implementation passed race-detector validation and repeated
100-run race testing for the mesh package. Explicit arbitration tests cover all
four identity/direction permutations.

The containerized runtime starts Alice, Bob, Carol and Dave. A deliberate
simultaneous Dave↔Bob connection collision demonstrates the deterministic rule:
Dave's identity (`4765bd805913e158`) is smaller than Bob's
(`563f793f6ad353dd`), so Dave's outbound connection survives while Bob's
competing outbound connection is rejected as a duplicate. Both nodes retain one
active peer.

These results establish direct-network runtime hardening, not multi-hop mesh
routing. M6 must provide separate routing evidence.

## M6 Results

M6 provides the project's demonstrated multi-hop capability. Report the bounded managed-flooding design, TTL controls, PacketID duplicate suppression, bounded forwarding queues, endpoint/transport session separation, and relay-blind E2EE boundary.

The acceptance record reports successful one-, two- and three-hop routing, cyclic-routing termination, multi-hop handshake correlation, resource-bound tests, relay confidentiality tests, race validation and Docker multi-hop operation.


## M6 Reference

- [[00-governance/04-Architecture-Review-M6]] — authoritative M6 architecture and acceptance record


## M8 Results

M8 established a functional presentation path over the approved secure mesh. The application service now owns conversation-oriented operations and endpoint message decryption, while the Fyne GUI remains isolated from crypto, session, mesh and routing internals.

The M8.3 Alice→Bob→Carol integration test passed using real loopback TCP, PeerManager, routing/forwarding, session establishment, ChaCha20-Poly1305 encryption/decryption and `SubscribeEvents()` delivery.

Repository validation passed:

```text
go test ./...
go test -race -p 1 ./...
go vet ./...
go build ./...
```

The GUI native build and live two-node runtime acceptance were verified. Docker Compose regression remains unverified because required dependency resolution was blocked by external DNS/proxy limitations.

M8 is therefore a verified application-integration milestone with explicit environment-limited GUI runtime evidence, not a claim of clean-environment deployment reproduction.

## M8 Reference

- [[00-governance/06-Architecture-Review-M8]] — authoritative M8 review and gate
- [[07-testing/10-M8.1-Acceptance-Evidence]]
- [[07-testing/11-M8.2-Acceptance-Evidence]]
- [[07-testing/12-M8.3-Acceptance-Evidence]]


## M9 Results

M9 adds final demonstration evidence for the implemented security boundaries.
M9.1 demonstrates Alice → Bob → Carol end-to-end encrypted messaging across a
relay while Bob observes routing metadata and ciphertext rather than the endpoint
plaintext. M9.2 passes eight adversarial demonstrations covering identity
mismatch, ciphertext tampering, replay, signature failure, malformed packets,
TTL, PacketID duplicate suppression and the production pending-handshake bound.
M9.3 records successful manual desktop verification of the two-node GUI workflow
in both messaging directions, including secure-session establishment, duplicate
connection behavior and wrong-PeerID behavior.

M9 does not constitute clean-environment deployment reproduction; Docker
regression remains unverified because of external proxy/DNS restrictions.

## M9 Reference

- [[00-governance/07-Architecture-Review-M9]]
- [[07-testing/19-M9-Evidence-Manifest]]
- [[07-testing/20-M9-Final-Report]]


## M10 Technical Documentation Baseline

M10.1 establishes a synchronized technical description of the implemented
system. It consolidates the component boundaries, cryptographic construction,
wire format, mesh forwarding model, application/GUI boundary, and
reproducibility constraints without introducing new implementation claims.

Reference: [[10-technical-documentation/08-M10.1-Technical-Documentation-Map]]


## Implemented Results Source

Use [[10-technical-documentation/11-Experimental-Methodology-and-Results]] and [[10-technical-documentation/12-Security-Evidence-Matrix]] as the evidence-grounded results sources. Do not add numerical performance results without a controlled measurement campaign.
