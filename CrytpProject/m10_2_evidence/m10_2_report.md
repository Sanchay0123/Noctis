# M10.2 Experimental Campaign Report (Corrected)

## 1. Environment
**OS/Kernel**: Linux sanchayjain-B550M-DS3H-AC-R2 7.0.0-31-generic Ubuntu 24.04.1
**Hardware**: AMD Ryzen 7 5700G (16 logical cores), 14Gi RAM
**Go Version**: go1.21.13 linux/amd64
**Git Commit**: 07fa48fb63662345b9381fbd22b7791b8c6097ea
**Git Status**: Clean production source. Unstaged `m10_2_evidence/` directory and `tests/m10_benchmark_test.go` present.

## 2. Methodology
A custom, test-isolated measurement harness (`tests/m10_benchmark_test.go`) was written to exercise the unmodified `meshchat` APIs. Real `ApplicationService` nodes were launched using `app.NewNode` on the loopback interface for latency and throughput measurements. Warm-up messages (excluded from measurement) and synchronized `time.Since` mechanisms were implemented for all delivery timings to measure one-way application-level delivery latency accurately.

## 3. Experiment A — Cryptography
**Status**: PARTIAL
**Evidence**: `m10_2_evidence/m10_2_crypto_benchmarks.txt`
**Method**: Go `testing.B` for 64B, 256B, 1KiB, and 4KiB AEAD Encrypt operations, as well as `crypto.GenerateEphemeralKey`, `ComputeSharedSecret`, and `ed25519` sign/verify.
**Observations**: X25519 shared-secret computation averaged ~39.6µs/op. Ed25519 signatures averaged ~20.3µs/op, while verification was ~46.0µs/op. ChaCha20-Poly1305 encryption scaled linearly up to ~1.9µs/op for 4KiB payloads. 
**Limitations/NOT MEASURED**: AEAD Decryption performance is explicitly classified as NOT MEASURED. The previous ~7ns/op result was discarded because repeated sequential inputs triggered the production cryptographic replay-protection window, prematurely halting decryption. Rerunning isolated decryption inside `testing.B` with fresh valid nonces without modifying the production source is impractical.

## 4. Experiment B — Direct Latency
**Status**: PASS
**Evidence**: `m10_2_evidence/m10_2_direct_latency.csv`
**Method**: Alice and Bob instantiated via `ApplicationService`. Alice called `SendMessage`; one-way application-level delivery latency was captured by reading Bob's `SubscribeEvents` channel.
**Trials**: 100 non-warmup trials per size.
**Observations**:
- **64B**: Min 0.008ms, Mean 0.063ms, Median 0.049ms, P95 0.119ms, P99 0.438ms, Max 0.479ms
- **256B**: Min 0.013ms, Mean 0.068ms, Median 0.066ms, P95 0.129ms, P99 0.268ms, Max 0.351ms
- **1KiB**: Min 0.015ms, Mean 0.063ms, Median 0.052ms, P95 0.145ms, P99 0.220ms, Max 0.224ms
- **4KiB**: Min 0.019ms, Mean 0.073ms, Median 0.069ms, P95 0.143ms, P99 0.222ms, Max 0.335ms

## 5. Experiment C — Multi-Hop Latency
**Status**: PASS
**Evidence**: `m10_2_evidence/m10_2_multihop_latency.csv`
**Method**: Alice -> Bob -> Carol topology. Alice connected to Bob, Bob connected to Carol. Alice established a session with Carol, routing encrypted payloads via Bob. The measurement represents ONE-WAY APPLICATION-LEVEL DELIVERY LATENCY.
**Trials**: 100 non-warmup trials per size.
**Observations**:
- **64B**: Min 0.008ms, Mean 0.073ms, Median 0.067ms, P95 0.161ms, P99 0.223ms, Max 0.327ms
- **256B**: Min 0.011ms, Mean 0.092ms, Median 0.070ms, P95 0.257ms, P99 0.286ms, Max 0.292ms
- **1KiB**: Min 0.015ms, Mean 0.122ms, Median 0.099ms, P95 0.335ms, P99 0.408ms, Max 0.426ms
- **4KiB**: Min 0.018ms, Mean 0.132ms, Median 0.104ms, P95 0.331ms, P99 0.382ms, Max 0.393ms

## 6. Experiment D — Application Delivery Throughput
**Status**: PASS
**Evidence**: `m10_2_evidence/m10_2_throughput.csv`
**Method**: Alice attempted 5,000 continuous `SendMessage` operations of 1KiB payloads to Bob. A `SubscribeEvents` listener at Bob recorded all inbound deliveries.
**Trials**: 5,000 payloads attempted.
**Observations**: 5,000 messages delivered (0 lost, 0 duplicates). Total payload bytes delivered: 5,120,000. Total elapsed time: 2.696 seconds. 
**Calculations**:
- 5,000 / 2.696 ≈ 1,854 messages/sec
- 5,120,000 / 2.696 ≈ 1,899,109 bytes/sec

## 7. Experiment E — Message Size Scaling
**Status**: PASS
**Evidence**: `m10_2_evidence/m10_2_message_size.csv`
**Method**: Sizes of 64B, 256B, 1KiB, 4KiB, 16KiB, 32KiB, 48KiB, and 60KiB (61440B) were tested over 100 trials each.
**Observations**: All 100 trials for 60KiB successfully traversed the current 64KiB framed-packet constraint. 60KiB was the largest tested application payload and successfully traversed the current implementation.

## 8. Experiment F — Mesh Forwarding
**Status**: PASS
**Evidence**: Multi-hop execution logged implicitly via `m10_2_multihop_latency.csv`.
**Method**: Used standard Alice -> Bob -> Carol deployment.
**Observations**: The experiment demonstrates the operational forwarding behavior of the unmodified `Router` across intermediate peers. 
**Important Note**: This experiment only demonstrates *forwarding behavior*. Relay blindness (i.e. Bob routing ciphertext without endpoint keys) is NOT claimed as independently proven by this latency experiment. Relay confidentiality and blindness were previously independently established in M9.1 using test-local observation.

## 9. Experiment G — Resource Usage
**Status**: PASS
**Evidence**: `m10_2_evidence/m10_2_resource_usage.txt`
**Method**: Used `/usr/bin/time -v` executing `go test -v ./tests -run TestM10_Throughput` to observe OS-level process resource utilization.
**Observations**:
- Elapsed (wall clock) time: 3.28 seconds
- User CPU time: 3.08 seconds
- System CPU time: 1.03 seconds
- Maximum resident set size (RSS): 99,004 kbytes (~96 MiB)
**Limitations**: CPU utilization percentage is not meaningfully established purely from this test context and is classified as NOT MEASURED.

## 10. Experiment H — Handshake Repeatability
**Status**: PASS
**Evidence**: `m10_2_evidence/m10_2_handshake.csv`
**Method**: Repeated direct TCP dialing and `StartConversation` handshake initialization 50 times sequentially.
**Trials**: 50 full sequential handshakes.
**Observations**: Demonstrates Handshake Repeatability. 
**Limitations/NOT MEASURED**: This experiment does NOT establish concurrent scaling. Scaling is NOT MEASURED here. The existing M9.2 pending-handshake bound remains authoritative: `MaxPendingHandshakes = 10`.

## 11. Experiment I — Concurrency Validation
**Status**: PASS
**Method**: `go test -race -p 1 ./...`
**Observations**: Passed successfully. No race conditions were triggered by the experimental campaign.

## 12. Raw Evidence
All generated artifacts remain inside `~/Documents/Noctis/CrytpProject/m10_2_evidence/`. Per-trial granularity is preserved within the CSVs. 

## 13. Reproducibility
- **Topology**: Standard `app.NewNode` executing over `127.0.0.1` ephemeral ports.
- **Commands**: Automated via isolated `tests/m10_benchmark_test.go` executed with standard `go test` and `/usr/bin/time -v`.

## 14. Limitations
Measurements captured represent non-networked CPU and loopback TCP operations within a unified local context. They do not simulate high-latency internet environments, realistic packet drops, or active adversarial DoS load characteristics. 

## 15. Test/Build Verification
- `go test ./...` -> `ok github.com/sanchayjain/meshchat/tests 15.814s`
- `go test -race -p 1 ./...` -> `ok github.com/sanchayjain/meshchat/tests 7.683s`
- `go vet ./...` -> Success (no output)
- `go build ./...` -> Success
- `go build -tags gui -o meshchat-gui ./cmd/meshchat-gui` -> Success

## 16. Git Integrity
Checked with `git diff --stat` and `git status --short`.
Zero production source, protocol schemas, or configuration files were modified. Modifications are restricted entirely to the unstaged `m10_2_evidence/` output directory and `tests/m10_benchmark_test.go`. The Obsidian documentation vault was completely unedited. No commits were generated.

## 17. Corrected Results Summary

| Experiment | Status | Measured Metric / Notes |
|---|---|---|
| A (Cryptography) | PARTIAL | Encryption speed measured. Decryption speed NOT MEASURED. |
| B (Direct Latency) | PASS | Application-level delivery latency measured locally. |
| C (Multi-Hop Latency)| PASS | Application-level delivery latency measured locally. |
| D (Throughput) | PASS | Application delivery throughput measured locally. |
| E (Message Size) | PASS | 60KiB traversed the transport framing constraint. |
| F (Mesh Forwarding) | PASS | Forwarding behavior validated. Blindness relies on M9.1. |
| G (Resource Usage) | PASS | RSS and User/Sys CPU time measured. |
| H (Handshake) | PASS | Repeatability measured. Concurrency scaling NOT MEASURED. |
| I (Concurrency) | PASS | `go test -race` validated. |

## 18. Explicit unresolved/NOT MEASURED items
- **AEAD Decryption Performance**: NOT MEASURED.
- **Concurrent Handshake Scaling**: NOT MEASURED.
- **CPU Utilization Percentage**: NOT MEASURED.
