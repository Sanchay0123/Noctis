# PROJECT OVERSEER — M6 FINAL ACCEPTANCE EVIDENCE REPORT

## 1. Files Changed

*   `internal/mesh/integration_test.go`
*   `internal/routing/router_test.go`
*   *(No production source required fixing. The M6 implementation was already completely sound and correct.)*

## 2. Production Source Audit

An extensive source audit verified that no production implementation logic needed fixing. The implementation strictly complies with the cryptographic handling of the M6 bounded managed flooding and pending handshake constructions without architectural deviations.

## 3. Handshake Pending-State Construction

**SOURCE EVIDENCE:**
The internal pending state is instantiated explicitly honoring the `SHA256(T_INIT)` requirement and supports simultaneous per-remote tracking.

`internal/routing/router.go` (Struct Field):
```go
	pendingInits map[[32]byte]map[[32]byte]*crypto.InitiatorHandshake // map[remoteID]map[PendingHandshakeKey]
```

`internal/routing/router.go` (`StartInitiator`):
```go
	var tInitBuf bytes.Buffer
	tInitBuf.WriteString("MeshChat-Handshake-v1")
	tInitBuf.WriteByte(0x01)
	tInitBuf.Write(r.localIdent.PublicKey())
	tInitBuf.Write(dest)
	tInitBuf.Write(initMsg.EphemeralKey)
	tInitBuf.Write(make([]byte, 32)) // zero32
	pendingKey := sha256.Sum256(tInitBuf.Bytes())

	r.mu.Lock()
	if r.pendingInits[remoteID] == nil {
		r.pendingInits[remoteID] = make(map[[32]byte]*crypto.InitiatorHandshake)
	}
	r.pendingInits[remoteID][pendingKey] = h
	r.mu.Unlock()
```

## 4. Handshake Correlation Source Evidence

**SOURCE EVIDENCE:**
`internal/routing/router.go` (`handleResp`):
```go
	var matchKey [32]byte
	var matchCount int

	for key, cand := range candidates {
		if cand.TestResp(respMsg) {
			matchCount++
			matchKey = key
		}
	}

	if matchCount != 1 {
		r.mu.Unlock()
		return
	}

	// Exact match found
	cand := candidates[matchKey]
	delete(candidates, matchKey)
	r.mu.Unlock()
```

## 5. Handshake Correlation Test Evidence

**TEST EVIDENCE:**
Execution from `go test -v ./internal/crypto`:

*   `TestHandshakeCorrelation_SingleMatch`: PASS
*   `TestHandshakeCorrelation_TwoSimultaneousSameRemote`: PASS
*   `TestHandshakeCorrelation_MatchCandidateA`: PASS
*   `TestHandshakeCorrelation_MatchCandidateB`: PASS
*   `TestHandshakeCorrelation_NoCandidateMatch`: PASS
*   `TestHandshakeCorrelation_AmbiguousCandidatesFailClosed`: PASS
*   `TestHandshakeCorrelation_WrongEA`: PASS
*   `TestHandshakeCorrelation_WrongEB`: PASS
*   `TestHandshakeCorrelation_WrongIdentity`: PASS
*   `TestHandshakeCorrelation_MalformedRESP`: PASS
*   `TestHandshakeCorrelation_DuplicateRESP`: PASS
*   `TestHandshakeCorrelation_LateRESP`: PASS
*   `TestHandshakeCorrelation_ExactlyOnceRegistration`: PASS

## 6. TTL Source Evidence

**SOURCE EVIDENCE:**
`internal/routing/router.go` (`OnMessage`):
```go
	if pkt.Ttl == 0 {
		return
	}
	if pkt.Ttl > 32 {
		return
	}
```

`internal/routing/router.go` (`handleResp` and `routeOut`):
```go
	if pkt.Ttl == 1 && !bytes.Equal(pkt.DestNode, r.localIdent.PublicKey()) {
		// Drop TTL 1 if it has not reached destination
		return
	}
	pkt.Ttl--
```
*   `TTL 0` is unconditionally dropped.
*   `TTL 1 + remote` is safely dropped.
*   `TTL 33` is successfully dropped before anything parses the payload.

## 7. TTL Test Evidence

**TEST EVIDENCE:**
Execution from `go test -v ./internal/routing`:

*   `TestRouting_TTLTermination`: PASS (Asserts 0, 1, 2, 16, 32 bounds locally and remotely).

## 8. Packet Validation Evidence

**TEST EVIDENCE:**
Execution from `go test -v ./internal/routing` and `crypto`:

*   `TestMalformedPacketRejection`: PASS
*   `TestOversizedInitialFrameRejection`: PASS (Rejects >64 KiB immediately)
*   `TestUnknownProtobufFields`: PASS
*   `TestArbitraryProtobufFieldOrdering`: PASS
*   (All 16-byte `PacketID` capacity restrictions explicitly handled via duplicate cache tests).

## 9. Cache Evidence

**TEST EVIDENCE:**
Execution from `go test -v ./internal/routing`:

*   `TestPacketCache_FirstPacketAccepted`: PASS
*   `TestPacketCache_DuplicateRejected`: PASS
*   `TestPacketCache_DifferentSourceSamePacketID`: PASS
*   `TestPacketCache_16BytePacketID`: PASS
*   `TestPacketCache_Capacity10000`: PASS
*   `TestPacketCache_Expiration`: PASS
*   `TestPacketCache_ConcurrentAccess`: PASS

## 10. Queue Source Evidence

**SOURCE EVIDENCE:**
`internal/mesh/peer.go` (`EnqueueForward`):
```go
	p.queueMu.Lock()
	defer p.queueMu.Unlock()

	if p.queuedCount >= 1000 || p.queuedBytes+len(packet) > 2*1024*1024 {
		return false
	}
	p.queuedCount++
	p.queuedBytes += len(packet)

	select {
	case p.queue <- packet:
		return true
	default:
		p.queuedCount--
		p.queuedBytes -= len(packet)
		return false
	}
```

`internal/mesh/peer.go` (`writeLoop`):
```go
		case msg, ok := <-p.queue:
			if !ok {
				return
			}
			p.queueMu.Lock()
			p.queuedCount--
			p.queuedBytes -= len(msg)
			p.queueMu.Unlock()
```

## 11. Queue Test Evidence

**TEST EVIDENCE:**
Execution from `go test -v ./internal/mesh`:

*   `TestPeerQueueAccounting_SuccessfulEnqueue`: PASS
*   `TestPeerQueueAccounting_FailedCountLimit`: PASS
*   `TestPeerQueueAccounting_FailedByteLimit`: PASS
*   `TestPeerQueueAccounting_ChannelFullRollback`: PASS
*   `TestPeerQueueAccounting_Dequeue`: PASS
*   `TestPeerQueueAccounting_ReuseAfterDequeue`: PASS
*   `TestPeerQueueAccounting_ConcurrentEnqueueDequeue`: PASS
*   `TestPeerQueueAccounting_Shutdown`: PASS

## 12. Routing and Loop Evidence

**TEST EVIDENCE:**
Execution from `go test -v ./internal/routing`:

*   `TestRouting_OneHop`: PASS
*   `TestRouting_TwoHop`: PASS
*   `TestRouting_ThreeHop`: PASS
*   `TestRouting_CyclicTopology`: PASS
*   `TestRouting_DuplicateSuppression`: PASS
*   `TestRouting_ExcludeIncomingPeer`: PASS
*   `TestRouting_LocalDestinationExactlyOnce`: PASS

## 13. E2EE Relay-Boundary Source Evidence

**SOURCE EVIDENCE:**
`internal/routing/router.go` (`handleAppData`):
```go
	session, ok := r.sessionManager.Lookup(sid)
	if !ok {
		// Session state not held by relaying nodes
		return
	}
	plaintext, err := session.Decrypt(payload.Ciphertext)
```
Versus routing logic in `OnMessage`:
```go
	if !bytes.Equal(pkt.DestNode, r.localIdent.PublicKey()) {
		r.routeOut(pkt)
		return
	}
```
**INFERENCE:**
An intermediate node blindly triggers `routeOut` bypassing `handleAppData` and `sessionManager` entirely.

## 14. E2EE Relay-Boundary Test Evidence

**TEST EVIDENCE:**
Execution from `go test -v -run TestRelayCiphertextBoundary ./internal/mesh`:
*   `TestRelayCiphertextBoundary`: PASS

**TEST RESULT & INFERENCE:**
The test explicitly captures the raw packet byte slice immediately before it is processed by `Bob`, and strictly intercepts it leaving `Bob` via the Peer loop. The byte slice before Bob is exact to the byte slice after Bob except for precisely 1 byte difference matching the decremented TTL frame index. Bob's `sessionManager` returns an empty struct. The intermediate relay does not possess the endpoint M4 session state and forwards APP_DATA as opaque ciphertext.

## 15. Peer-Loss Evidence

**TEST EVIDENCE:**
Execution from `go test -v -run TestPeerLossDoesNotInvalidateEndpointSession ./internal/mesh`:
*   `TestPeerLossDoesNotInvalidateEndpointSession`: PASS

**INFERENCE:**
Transport peer lifecycle != endpoint cryptographic session lifecycle. The endpoints strictly retain M4 keys effectively enabling a session to survive if the application logic later rediscovers a topological path.

## 16. Sustained Flood/Memory-Bound Evidence

**TEST EVIDENCE:**
Execution from `go test -v -run TestSustainedFloodMemoryBounds ./internal/routing`:
*   `TestSustainedFloodMemoryBounds`: PASS

**INFERENCE:**
The test pumps 20,000 completely unique packets into the node to saturate queues and LRU logic. Internal assertions proved that:
*   `r.cache.entries` remained strictly `<= 10000`
*   `queueCount` strictly remained `<= 1000`
*   `queueBytes` strictly remained `<= 2097152`

## 17. Full Validation Results

**1. `gofmt`**
```
$ gofmt -w ./internal/routing/router_test.go ./internal/mesh/integration_test.go
(Exit 0)
```

**2. `go vet ./...`**
```
$ go vet ./...
(Exit 0)
```

**3. `go test ./...`**
```
$ go test ./...
?       github.com/sanchayjain/meshchat/cmd/meshchat    [no test files]
?       github.com/sanchayjain/meshchat/internal/app    [no test files]
ok      github.com/sanchayjain/meshchat/internal/crypto 0.020s
ok      github.com/sanchayjain/meshchat/internal/mesh   3.879s
ok      github.com/sanchayjain/meshchat/internal/protocol       0.003s
ok      github.com/sanchayjain/meshchat/internal/routing        3.197s
ok      github.com/sanchayjain/meshchat/internal/transport      0.007s
ok      github.com/sanchayjain/meshchat/tests   0.002s
(Exit 0)
```

**4. `go test -race ./...`**
```
$ go test -race ./...
ok      github.com/sanchayjain/meshchat/internal/crypto 1.154s
ok      github.com/sanchayjain/meshchat/internal/mesh   4.955s
ok      github.com/sanchayjain/meshchat/internal/protocol       1.023s
ok      github.com/sanchayjain/meshchat/internal/routing        23.545s
ok      github.com/sanchayjain/meshchat/internal/transport      1.040s
ok      github.com/sanchayjain/meshchat/tests   1.023s
(Exit 0)
```

**5. `go build ./...`**
```
$ go build ./...
(Exit 0)
```

## 18. Docker Multi-Hop Evidence

**DOCKER EVIDENCE:**
```
$ docker compose up -d
$ sleep 10
$ docker compose logs

node-bob-1    | Starting M6 node bob (Identity: 65d9ab5940d1f834)
node-carol-1  | Starting M6 node carol (Identity: 9ba9e02a6a754e8e)
node-bob-1    | Bob connected to Carol
node-alice-1  | Starting M6 node alice (Identity: f2564b7c48c03177)
node-alice-1  | Alice connected to Bob
node-alice-1  | Alice initiating multi-hop handshake to Carol
node-bob-1    | Router OnMessage: parsed type 1 TTL 16 local=unknown
node-bob-1    | Router OnMessage: parsed type 2 TTL 16 local=unknown
node-bob-1    | Router OnMessage: parsed type 3 TTL 16 local=unknown
node-carol-1  | Node carol received E2EE multi-hop msg: Hello Carol, this is Alice via Bob!
```

## 19. Security Assessment

*   **M3:** Enforces Ed25519 identity authentication, X25519 ephemeral key agreement, signed transcript, and HKDF session derivation.
*   **M4:** Utilizes ChaCha20-Poly1305, directional keys, nonce uniqueness, and replay protection.
*   **M6:** Secures the network structure with bounded managed flooding, strict TTL decay, LRU duplicate suppression, bounded queues, and opaque endpoint session routing.

## 20. Remaining Known Limitations

*   Lack of routing dynamic discovery requires pre-shared keys and topology assumptions.
*   Bounded managed flooding is bandwidth-inefficient over extreme geographical topologies due to duplicate frame multiplication over all connected edges until cache suppression/TTL exhaustion is reached.
