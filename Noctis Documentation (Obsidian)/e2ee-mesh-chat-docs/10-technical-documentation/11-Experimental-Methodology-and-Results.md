# Experimental Methodology and Results

## Status

**M10.2 — PASS / CLOSED**

This document records the approved M10.2 experimental methodology and the controlled local measurements actually collected. Results are scoped to the stated loopback environment and are not presented as Internet deployment benchmarks.

## 1. Evaluation Strategy

The evaluation uses a layered approach:

1. cryptographic unit and negative tests
2. protocol and transport validation
3. mesh/routing integration tests
4. adversarial security demonstrations
5. multi-hop relay-blindness demonstration
6. manual desktop GUI verification
7. controlled local performance experiments

The evaluation is behavioral, security-oriented and experimentally scoped. Numerical results are included only where the campaign actually collected them.

## 2. Reproducibility Baseline

The M9 evidence records the following environment baseline:

- Project working directory: `~/Documents/Noctis/CrytpProject`
- Git repository top level: `/home/sanchayjain/Documents/Noctis`
- Recorded commit: `a67e634`
- Go: `go1.21.13 linux/amd64`
- Kernel: `Linux 7.0.0-31-generic`

The repository is located inside the larger Noctis Git repository; this does not change the project working directory.

## 3. Automated Validation

The M9 evidence records successful execution of:

```text
go vet ./...
go test ./...
go test -race -p 1 ./...
go build ./...
go build -tags gui -o meshchat-gui ./cmd/meshchat-gui
```

These commands establish compilation, static analysis, functional tests and race-detector coverage for the recorded environment. They do not establish deployment reproducibility on every operating system or container runtime.

## 4. Security Experiment Matrix

| Experiment | Setup | Expected result | Recorded result |
|---|---|---|---|
| Wrong PeerID | Dial with incorrect expected identity | Connection/authentication rejected | PASS |
| Ciphertext tampering | Modify protected application ciphertext | AEAD authentication fails | PASS |
| AEAD replay | Re-submit an accepted protected message | Replay rejected | PASS |
| Invalid signature | Alter handshake signature | Authentication fails | PASS |
| Malformed packet | Supply structurally invalid packet | Packet rejected | PASS |
| TTL enforcement | Exercise expired/near-expiry forwarding | Packet does not forward indefinitely | PASS |
| PacketID duplicate suppression | Repeat same `(source_node, packet_id)` | Duplicate forwarding suppressed | PASS |
| Pending-handshake bound | Attempt 15 incomplete handshakes | Simultaneous admission remains ≤ 10 | PASS |

## 5. Relay-Blindness Experiment

### Topology

```text
Alice  →  Bob  →  Carol
          relay
```

The intermediate Bob node observes routing metadata and ciphertext. The endpoint application message is decrypted at Carol using the endpoint session established for the communicating endpoints.

The experiment demonstrates the architectural property that routing/forwarding does not require access to the endpoint application plaintext.

## 6. GUI Experiment

The GUI experiment is a two-node desktop workflow:

```text
Alice  ↔  Bob
```

The Project Overseer manually verified:

1. Alice identity display
2. Bob identity display
3. secure session establishment
4. Alice → Bob messaging
5. Bob → Alice messaging
6. duplicate Add Peer behavior
7. wrong-PeerID behavior

Automated GUI interaction was unavailable in the headless implementation environment. The manual desktop verification is therefore recorded as manual evidence rather than automated CI evidence.

## 7. Performance Results Policy

No latency, throughput, CPU, memory, routing-overhead or connection-establishment benchmark is claimed in M10.2 because a reproducible performance measurement campaign was not part of the approved M9 evidence.

If performance evaluation becomes a project requirement, it should be conducted as a separate controlled experiment with:

- fixed hardware/software environment
- defined topology
- repeated trials
- warm-up policy
- sample counts
- statistical summary
- measurement instrumentation
- clearly stated payload sizes

Until such measurements exist, qualitative statements such as "functional" or "bounded" must not be converted into numerical performance claims.

## 8. Interpretation

The experimental evidence supports the following conclusions:

- the cryptographic/session layer behaves as specified by its tests
- application ciphertext survives multi-hop forwarding without relay-side plaintext decryption
- specified adversarial inputs are rejected by the tested controls
- pending handshake admission is bounded by the configured production limit
- the GUI workflow operates correctly in the manually verified desktop environment

The evidence does not support claims of universal attack resistance, anonymity, guaranteed availability or production deployment readiness.

## Related Records

- [[07-testing/15-M9-Test-Suite.log]]
- [[07-testing/16-M9-Attack-Evidence]]
- [[07-testing/17-M9-Relay-Blindness.log]]
- [[07-testing/18-M9-GUI-Evidence]]
- [[07-testing/19-M9-Evidence-Manifest]]
- [[10-technical-documentation/10-Security-Analysis]]
