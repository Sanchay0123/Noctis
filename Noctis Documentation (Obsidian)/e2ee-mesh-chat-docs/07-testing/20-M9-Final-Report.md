# M9 — Security Validation, Demonstration & Evidence

## 1. Executive Summary

M9 validates the implemented secure-messaging prototype through relay-blind
end-to-end encryption evidence, adversarial security demonstrations, and a live
two-node GUI demonstration.

M9.1, M9.2 and M9.3 are closed based on their approved evidence and the final
manual GUI verification performed by the Project Overseer.

M9.4 assembles and synchronizes the documentation and evidence record. This
report is submitted for Project Overseer / Project Director review and does not
self-approve M9.4.

## 2. M9 Architecture

M9 uses two distinct demonstrations:

1. **Three-node security topology:** Alice → Bob → Carol, demonstrating that a
   forwarding relay handles routing metadata and ciphertext while the endpoint
   application session remains at Alice and Carol.
2. **Two-node GUI topology:** Alice ↔ Bob, demonstrating the functional desktop
   chat workflow.

The routing layer remains separate from endpoint cryptographic processing.
M6 remains bounded managed flooding rather than route discovery or route
selection.

## 3. M9.1 — E2EE Across Relay

**PASS / CLOSED**

Alice successfully sends encrypted application data through Bob to Carol.
Bob's relay observation contains APP_DATA routing metadata and ciphertext.
Bob's implemented session manager does not contain the Alice–Carol endpoint
application session required to decrypt the packet. Carol successfully decrypts
and receives the plaintext at the endpoint.

The M9.1 observation mechanism is test-local and does not add a production
observer API.

## 4. M9.2 — Security Attack Demonstrations

**PASS / CLOSED**

Eight mandatory demonstrations passed:

1. Wrong PeerID
2. Ciphertext tampering
3. AEAD replay
4. Invalid signature
5. Malformed packet
6. TTL enforcement
7. PacketID duplicate suppression
8. Pending-handshake resource bound

The production pending-handshake bound is 10. The authoritative test attempted
15 incomplete connections and observed simultaneous admission no greater than
10.

## 5. M9.3 — GUI Demonstration

**PASS / CLOSED**

The implementation agent's environment was headless and did not provide an
interactable GUI runtime. The Project Overseer subsequently performed live
desktop verification.

Verified scenarios:

- Alice identity
- Bob identity
- secure session establishment
- Alice → Bob encrypted messaging
- Bob → Alice encrypted messaging
- duplicate Add Peer behavior
- wrong PeerID behavior

No screenshot is claimed as packaged unless the actual image artifact is present.

## 6. Test and Build Verification

Reported M9 validation passed:

```text
go test ./...
go test -race -p 1 ./...
go vet ./...
go build ./...
go build -tags gui -o meshchat-gui ./cmd/meshchat-gui
```

Docker regression remains **NOT VERIFIED** because external proxy/DNS restrictions
prevented dependency resolution in the affected environment.

## 7. Security Properties Demonstrated

- Ed25519-authenticated identities
- X25519-based authenticated session establishment
- transcript-bound HKDF-SHA-256 key derivation
- ChaCha20-Poly1305 application confidentiality/integrity
- replay protection
- authenticated packet/session binding
- relay forwarding of ciphertext without endpoint plaintext access in the
  demonstrated topology
- bounded packet forwarding and duplicate suppression
- bounded pending-handshake admission
- GUI/application separation from cryptographic internals

## 8. Security Properties NOT Demonstrated

M9 does not demonstrate:

- anonymity
- complete metadata hiding
- resistance to global traffic analysis
- endpoint compromise resistance
- guaranteed delivery
- immunity to denial of service
- production-grade network resilience
- global Sybil resistance
- deniability
- a concrete Prometheus exporter

## 9. Known Limitations

- GUI interaction was manually verified on the Project Overseer's desktop rather
  than reproduced automatically in the headless implementation environment.
- The supplied documentation archive does not contain the seven GUI screenshots.
- Docker regression remains unverified due external proxy/DNS restrictions.
- The repository telemetry subsystem remains supporting hooks/no-op behavior;
  a concrete Prometheus exporter is not implemented.
- The demonstrated relay boundary does not claim protection against a
  compromised endpoint or a relay that obtains endpoint keys through other means.

## 10. Evidence Manifest

See [[07-testing/19-M9-Evidence-Manifest]] and the individual M9 evidence records
in `07-testing`.

## 11. M9 Final Status

| Phase | Status |
|---|---|
| M9 Architecture | APPROVED |
| M9.1 | PASS / CLOSED |
| M9.2 | PASS / CLOSED |
| M9.3 | PASS / CLOSED |
| M9.4 | PASS / CLOSED — documentation/evidence audit completed |

**M9.4 was not self-approved by this document; it was subsequently audited and accepted by the Project Overseer. M9 is formally closed.**
