package crypto

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"math"
	"sync"
	"testing"

	"golang.org/x/crypto/hkdf"
)

func generateTestIdentities(t *testing.T) (*NodeIdentity, *NodeIdentity) {
	_, privA, _ := ed25519.GenerateKey(nil)
	idA, _ := LoadIdentity(privA)

	_, privB, _ := ed25519.GenerateKey(nil)
	idB, _ := LoadIdentity(privB)

	return idA, idB
}

func getTestSessions(t *testing.T) (*Session, *Session) {
	alice, bob := generateTestIdentities(t)
	hA, _ := NewInitiatorHandshake(alice, bob.PublicKey())
	initMsg, _ := hA.GenerateInit()
	hB, _ := NewResponderHandshake(bob, alice.PublicKey())
	hB.ProcessInit(initMsg)
	respMsg, _ := hB.GenerateResp()
	hA.ProcessResp(respMsg)
	sessA, _ := hA.Session()
	sessB, _ := hB.Session()
	return sessA, sessB
}

// Previous Handshake Tests

func TestCompleteSuccessfulHandshake(t *testing.T) {
	alice, bob := generateTestIdentities(t)
	hA, _ := NewInitiatorHandshake(alice, bob.PublicKey())
	initMsg, err := hA.GenerateInit()
	if err != nil {
		t.Fatalf("Alice failed GenerateInit: %v", err)
	}
	hB, _ := NewResponderHandshake(bob, alice.PublicKey())
	if err := hB.ProcessInit(initMsg); err != nil {
		t.Fatalf("Bob failed ProcessInit: %v", err)
	}
	respMsg, err := hB.GenerateResp()
	if err != nil {
		t.Fatalf("Bob failed GenerateResp: %v", err)
	}
	if err := hA.ProcessResp(respMsg); err != nil {
		t.Fatalf("Alice failed ProcessResp: %v", err)
	}
	sessA, err := hA.Session()
	if err != nil {
		t.Fatalf("Alice missing session")
	}
	sessB, err := hB.Session()
	if err != nil {
		t.Fatalf("Bob missing session")
	}
	if !bytes.Equal(sessA.ID(), sessB.ID()) {
		t.Fatalf("Session IDs do not match")
	}
	if !bytes.Equal(sessA.SendKey(), sessB.ReceiveKey()) {
		t.Fatalf("Alice send != Bob receive")
	}
	if !bytes.Equal(sessA.ReceiveKey(), sessB.SendKey()) {
		t.Fatalf("Alice receive != Bob send")
	}
	if bytes.Equal(sessA.SendKey(), sessA.ReceiveKey()) {
		t.Fatalf("Alice send == Alice receive")
	}
}

func TestTranscriptTampering(t *testing.T) {
	alice, bob := generateTestIdentities(t)
	hA, _ := NewInitiatorHandshake(alice, bob.PublicKey())
	initMsg, _ := hA.GenerateInit()
	_, privC, _ := ed25519.GenerateKey(nil)
	idC, _ := LoadIdentity(privC)
	hB, _ := NewResponderHandshake(bob, idC.PublicKey())
	if err := hB.ProcessInit(initMsg); err == nil {
		t.Fatalf("ProcessInit should have failed due to tampered ID_A")
	}
	hA2, _ := NewInitiatorHandshake(alice, bob.PublicKey())
	initMsg2, _ := hA2.GenerateInit()
	initMsg2.Signature[0] ^= 0xFF
	hB2, _ := NewResponderHandshake(bob, alice.PublicKey())
	if err := hB2.ProcessInit(initMsg2); err == nil {
		t.Fatalf("ProcessInit should have failed due to tampered signature")
	}
	hA3, _ := NewInitiatorHandshake(alice, bob.PublicKey())
	initMsg3, _ := hA3.GenerateInit()
	hB3, _ := NewResponderHandshake(bob, alice.PublicKey())
	hB3.ProcessInit(initMsg3)
	respMsg, _ := hB3.GenerateResp()
	respMsg.Signature[1] ^= 0xAA
	if err := hA3.ProcessResp(respMsg); err == nil {
		t.Fatalf("ProcessResp should have failed due to tampered responder signature")
	}
}

func TestIdentitySwap(t *testing.T) {
	alice, bob := generateTestIdentities(t)
	hA, _ := NewInitiatorHandshake(alice, bob.PublicKey())
	initMsg, _ := hA.GenerateInit()
	hB, _ := NewResponderHandshake(bob, bob.PublicKey())
	if err := hB.ProcessInit(initMsg); err == nil {
		t.Fatalf("ProcessInit should fail if initiator identity doesn't match")
	}
}

func TestEphemeralKeySubstitution(t *testing.T) {
	alice, bob := generateTestIdentities(t)
	hA, _ := NewInitiatorHandshake(alice, bob.PublicKey())
	initMsg, _ := hA.GenerateInit()
	initMsg.EphemeralKey[0] ^= 0xFF
	hB, _ := NewResponderHandshake(bob, alice.PublicKey())
	if err := hB.ProcessInit(initMsg); err == nil {
		t.Fatalf("ProcessInit should fail if E_A is tampered")
	}
}

func TestStateMachineFailures(t *testing.T) {
	alice, bob := generateTestIdentities(t)
	hA, _ := NewInitiatorHandshake(alice, bob.PublicKey())
	err := hA.ProcessResp(&RespMessage{EphemeralKey: make([]byte, 32), Signature: make([]byte, 64)})
	if err != ErrInvalidHandshakeState {
		t.Fatalf("Expected ErrInvalidHandshakeState")
	}
	hB, _ := NewResponderHandshake(bob, alice.PublicKey())
	_, err = hB.GenerateResp()
	if err != ErrInvalidHandshakeState {
		t.Fatalf("Expected ErrInvalidHandshakeState")
	}
}

func TestRoleConfusion(t *testing.T) {
	alice, bob := generateTestIdentities(t)
	hA, _ := NewInitiatorHandshake(alice, bob.PublicKey())
	initMsg, _ := hA.GenerateInit()
	err := hA.ProcessResp(&RespMessage{EphemeralKey: initMsg.EphemeralKey, Signature: initMsg.Signature})
	if err == nil {
		t.Fatalf("Alice should reject initiator signature formatted as responder transcript")
	}
}

func TestSessionIDKnownConstruction(t *testing.T) {
	alice, bob := generateTestIdentities(t)
	hA, _ := NewInitiatorHandshake(alice, bob.PublicKey())
	initMsg, _ := hA.GenerateInit()
	hB, _ := NewResponderHandshake(bob, alice.PublicKey())
	hB.ProcessInit(initMsg)
	respMsg, _ := hB.GenerateResp()
	hA.ProcessResp(respMsg)
	s, _ := hA.Session()
	expectedTResp := buildTResp(alice.PublicKey(), bob.PublicKey(), initMsg.EphemeralKey, respMsg.EphemeralKey)
	expectedSessionID := sha256.Sum256(expectedTResp)
	if !bytes.Equal(s.ID(), expectedSessionID[:]) {
		t.Fatalf("Session ID does not strictly match SHA256(T_RESP)")
	}
}

func TestSessionIDTranscriptSensitivity(t *testing.T) {
	alice, bob := generateTestIdentities(t)
	hA, _ := NewInitiatorHandshake(alice, bob.PublicKey())
	initMsg, _ := hA.GenerateInit()
	hB, _ := NewResponderHandshake(bob, alice.PublicKey())
	hB.ProcessInit(initMsg)
	respMsg, _ := hB.GenerateResp()
	originalTResp := buildTResp(alice.PublicKey(), bob.PublicKey(), initMsg.EphemeralKey, respMsg.EphemeralKey)
	originalSessionID := sha256.Sum256(originalTResp)
	modifiedTResp := make([]byte, len(originalTResp))
	copy(modifiedTResp, originalTResp)
	modifiedTResp[0] ^= 0xFF
	modifiedSessionID := sha256.Sum256(modifiedTResp)
	if bytes.Equal(originalSessionID[:], modifiedSessionID[:]) {
		t.Fatalf("Session ID did not change when transcript was modified")
	}
}

func TestHKDFKnownAnswer(t *testing.T) {
	idA := bytes.Repeat([]byte("A"), 32)
	idB := bytes.Repeat([]byte("B"), 32)
	eA := bytes.Repeat([]byte("X"), 32)
	eB := bytes.Repeat([]byte("Y"), 32)
	ss := bytes.Repeat([]byte("S"), 32)
	tResp := buildTResp(idA, idB, eA, eB)
	expectedSessionID, _ := hex.DecodeString("45422eb08ee361b5cbf1c81650f7d97c371f23a95698b89dc775349779a3d5ee")
	tRespSum := sha256.Sum256(tResp)
	if !bytes.Equal(tRespSum[:], expectedSessionID) {
		t.Fatalf("T_RESP construction mismatch")
	}
	session := deriveSession(tResp, ss, true, make([]byte, 32))
	expectedSend, _ := hex.DecodeString("750b44a95d57c29b536ca57990a2e2e19e5265bf64262e8aebd970785c76e6f1")
	expectedReceive, _ := hex.DecodeString("d016ddec6d26a85c38a3fdedb084d018ded4b75ad6d04b579035c75848e17e91")
	if !bytes.Equal(session.ID(), expectedSessionID) {
		t.Fatalf("Session ID mismatch against known answer")
	}
	if !bytes.Equal(session.SendKey(), expectedSend) {
		t.Fatalf("Initiator SendKey mismatch against known answer")
	}
	if !bytes.Equal(session.ReceiveKey(), expectedReceive) {
		t.Fatalf("Initiator ReceiveKey mismatch against known answer")
	}
	salt := make([]byte, 32)
	copy(salt, expectedSessionID)
	prk := hkdf.Extract(sha256.New, ss, salt)
	expectedPRK, _ := hex.DecodeString("6f14e26edb1b552f2ee4864e8f8d91354fdee2e7f3a610b620be0ac0fd4df79b")
	if !bytes.Equal(prk, expectedPRK) {
		t.Fatalf("PRK mismatch against known answer")
	}
}

func TestTwoIndependentHandshakes(t *testing.T) {
	alice, bob := generateTestIdentities(t)
	hA1, _ := NewInitiatorHandshake(alice, bob.PublicKey())
	iM1, _ := hA1.GenerateInit()
	hB1, _ := NewResponderHandshake(bob, alice.PublicKey())
	hB1.ProcessInit(iM1)
	rM1, _ := hB1.GenerateResp()
	hA1.ProcessResp(rM1)
	s1, _ := hA1.Session()
	hA2, _ := NewInitiatorHandshake(alice, bob.PublicKey())
	iM2, _ := hA2.GenerateInit()
	hB2, _ := NewResponderHandshake(bob, alice.PublicKey())
	hB2.ProcessInit(iM2)
	rM2, _ := hB2.GenerateResp()
	hA2.ProcessResp(rM2)
	s2, _ := hA2.Session()
	if bytes.Equal(iM1.EphemeralKey, iM2.EphemeralKey) {
		t.Fatalf("Ephemeral keys were reused across handshakes")
	}
	if bytes.Equal(s1.ID(), s2.ID()) {
		t.Fatalf("Session IDs match across independent handshakes")
	}
	if bytes.Equal(s1.SendKey(), s2.SendKey()) {
		t.Fatalf("SendKeys match across independent handshakes")
	}
	if bytes.Equal(s1.ReceiveKey(), s2.ReceiveKey()) {
		t.Fatalf("ReceiveKeys match across independent handshakes")
	}
}

func TestMalformedInputLengths(t *testing.T) {
	alice, bob := generateTestIdentities(t)
	_, err := NewInitiatorHandshake(alice, make([]byte, 31))
	if err == nil {
		t.Fatalf("Initiator did not reject wrong ID_B length")
	}
	_, err = NewResponderHandshake(bob, make([]byte, 31))
	if err == nil {
		t.Fatalf("Responder did not reject wrong ID_A length")
	}
	hA, _ := NewInitiatorHandshake(alice, bob.PublicKey())
	hA.GenerateInit()
	hB, _ := NewResponderHandshake(bob, alice.PublicKey())
	err = hB.ProcessInit(&InitMessage{EphemeralKey: make([]byte, 31), Signature: make([]byte, 64)})
	if err == nil {
		t.Fatalf("Responder did not reject wrong E_A length")
	}
	err = hB.ProcessInit(&InitMessage{EphemeralKey: make([]byte, 32), Signature: make([]byte, 63)})
	if err == nil {
		t.Fatalf("Responder did not reject wrong signature length")
	}
}

func TestDifferentIdentitySignature(t *testing.T) {
	alice, bob := generateTestIdentities(t)
	_, charlie := generateTestIdentities(t)
	hA, _ := NewInitiatorHandshake(alice, bob.PublicKey())
	initMsg, _ := hA.GenerateInit()
	tInit := buildTInit(alice.PublicKey(), bob.PublicKey(), initMsg.EphemeralKey)
	charlieSig, _ := charlie.Sign(tInit)
	initMsg.Signature = charlieSig
	hB, _ := NewResponderHandshake(bob, alice.PublicKey())
	if err := hB.ProcessInit(initMsg); err == nil {
		t.Fatalf("ProcessInit succeeded with signature from wrong identity")
	}
}

func TestOldResponseNewInit(t *testing.T) {
	alice, bob := generateTestIdentities(t)
	hA1, _ := NewInitiatorHandshake(alice, bob.PublicKey())
	iM1, _ := hA1.GenerateInit()
	hB1, _ := NewResponderHandshake(bob, alice.PublicKey())
	hB1.ProcessInit(iM1)
	rM1, _ := hB1.GenerateResp()
	hA2, _ := NewInitiatorHandshake(alice, bob.PublicKey())
	hA2.GenerateInit()
	if err := hA2.ProcessResp(rM1); err == nil {
		t.Fatalf("Alice improperly processed an old response with a new state")
	}
}

func TestDuplicateInvalidState(t *testing.T) {
	alice, bob := generateTestIdentities(t)
	hA, _ := NewInitiatorHandshake(alice, bob.PublicKey())
	if err := hA.ProcessResp(&RespMessage{make([]byte, 32), make([]byte, 64)}); err == nil {
		t.Fatalf("Allowed ProcessResp before GenerateInit")
	}
	iM, _ := hA.GenerateInit()
	hB, _ := NewResponderHandshake(bob, alice.PublicKey())
	hB.ProcessInit(iM)
	if err := hB.ProcessInit(iM); err == nil {
		t.Fatalf("Allowed duplicate ProcessInit")
	}
	rM, _ := hB.GenerateResp()
	if _, err := hB.GenerateResp(); err == nil {
		t.Fatalf("Allowed GenerateResp after state reached Established")
	}
	hA.ProcessResp(rM)
	if err := hA.ProcessResp(rM); err == nil {
		t.Fatalf("Allowed duplicate ProcessResp")
	}
	hB2, _ := NewResponderHandshake(bob, alice.PublicKey())
	hB2.ProcessInit(&InitMessage{EphemeralKey: make([]byte, 32), Signature: make([]byte, 64)}) // fails
	if _, err := hB2.GenerateResp(); err == nil {
		t.Fatalf("Allowed GenerateResp after state reached Failed")
	}
}

// M3.3 AEAD Tests

func TestAEADKnownAnswer(t *testing.T) {
	// Values from python KAT script
	key, _ := hex.DecodeString("4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b")
	sessID, _ := hex.DecodeString("5353535353535353535353535353535353535353535353535353535353535353")
	seq := uint64(1)
	plaintext, _ := hex.DecodeString("48656c6c6f2c204d6573684368617421")
	expectedCiphertextWithTag, _ := hex.DecodeString("24a202cb5dea7e950eadf8810133655210a18dd96f79487b87957c8a084c1ab5")

	// Inject test session
	s := &Session{
		id:          sessID,
		sendKey:     key,
		receiveKey:  key,
		isInitiator: true, // Direction 0x01
		sendSeq:     1,
	}

	outSeq, ct, err := s.EncryptMessage(plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}
	if outSeq != seq {
		t.Fatalf("Seq mismatch")
	}
	if !bytes.Equal(ct, expectedCiphertextWithTag) {
		t.Fatalf("Ciphertext mismatch against KAT")
	}

	// Also verify decrypt
	s.highestReceived = 0 // To allow 1
	s.replayBitmap = 0
	// For decryption, the receiver must be the opposite role.
	sReceiver := &Session{
		id:          sessID,
		receiveKey:  key,
		isInitiator: false, // Expects 0x01
	}

	pt, err := sReceiver.DecryptMessage(outSeq, ct)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}
	if !bytes.Equal(pt, plaintext) {
		t.Fatalf("Plaintext mismatch")
	}
}

func TestAEADCiphertextModification(t *testing.T) {
	sA, sB := getTestSessions(t)
	seq, ct, _ := sA.EncryptMessage([]byte("hello"))

	// Modify tag (last 16 bytes)
	ct[len(ct)-1] ^= 0xFF
	_, err := sB.DecryptMessage(seq, ct)
	if err == nil {
		t.Fatalf("Accepted modified tag")
	}

	// Restore tag, modify ciphertext payload
	ct[len(ct)-1] ^= 0xFF
	ct[0] ^= 0xFF
	_, err = sB.DecryptMessage(seq, ct)
	if err == nil {
		t.Fatalf("Accepted modified ciphertext")
	}
}

func TestAEADWrongAADFields(t *testing.T) {
	sA, sB := getTestSessions(t)
	seq, ct, _ := sA.EncryptMessage([]byte("hello"))

	// Wrong sequence
	_, err := sB.DecryptMessage(seq+1, ct)
	if err == nil {
		t.Fatalf("Accepted with wrong sequence in AAD")
	}

	// Wrong Session ID
	originalID := sB.id[0]
	sB.id[0] ^= 0xFF
	_, err = sB.DecryptMessage(seq, ct)
	if err == nil {
		t.Fatalf("Accepted with wrong Session ID")
	}
	sB.id[0] = originalID

	// Wrong Key
	originalKey := sB.receiveKey[0]
	sB.receiveKey[0] ^= 0xFF
	_, err = sB.DecryptMessage(seq, ct)
	if err == nil {
		t.Fatalf("Accepted with wrong key")
	}
	sB.receiveKey[0] = originalKey

	// Wrong Direction
	sB.isInitiator = true // Should be false for Bob. Forces expectation of 0x00 instead of 0x01
	_, err = sB.DecryptMessage(seq, ct)
	if err == nil {
		t.Fatalf("Accepted with wrong direction flag")
	}
	sB.isInitiator = false
}

func TestAEADWrongProtocolContext(t *testing.T) {
	sA, sB := getTestSessions(t)

	// Alter the global context temporarily
	originalCtx := appDataContext[0]
	appDataContext[0] ^= 0xFF
	seq, ct, _ := sA.EncryptMessage([]byte("hello"))
	appDataContext[0] = originalCtx // Restore for Bob

	_, err := sB.DecryptMessage(seq, ct)
	if err == nil {
		t.Fatalf("Accepted wrong protocol context")
	}
}

func TestAEADTruncatedCiphertext(t *testing.T) {
	sA, sB := getTestSessions(t)
	seq, ct, _ := sA.EncryptMessage([]byte("hello"))

	ct = ct[:15] // Less than chacha20poly1305.Overhead
	_, err := sB.DecryptMessage(seq, ct)
	if err != ErrDecryptFailed {
		t.Fatalf("Did not fail truncated ciphertext with generic error: %v", err)
	}
}

func TestAEADReplayDuplicate(t *testing.T) {
	sA, sB := getTestSessions(t)
	seq, ct, _ := sA.EncryptMessage([]byte("hello"))

	_, err := sB.DecryptMessage(seq, ct)
	if err != nil {
		t.Fatalf("Initial decrypt failed: %v", err)
	}

	// Replay exact packet
	_, err = sB.DecryptMessage(seq, ct)
	if err != ErrReplayDetected {
		t.Fatalf("Failed to detect duplicate replay")
	}
}

func TestAEADReplayOutsideWindow(t *testing.T) {
	sA, sB := getTestSessions(t)
	seq0, ct0, _ := sA.EncryptMessage([]byte("msg0"))

	// Advance window to 64
	sA.sendSeq = 64
	seq64, ct64, _ := sA.EncryptMessage([]byte("msg64"))

	sB.DecryptMessage(seq64, ct64) // Advances window. highestReceived = 64

	// Attempt to process seq0
	_, err := sB.DecryptMessage(seq0, ct0)
	if err != ErrReplayDetected {
		t.Fatalf("Failed to detect replay outside 64-message window: %v", err)
	}
}

func TestAEADValidOutOfOrder(t *testing.T) {
	sA, sB := getTestSessions(t)
	seq0, ct0, _ := sA.EncryptMessage([]byte("msg0"))
	seq1, ct1, _ := sA.EncryptMessage([]byte("msg1"))
	seq2, ct2, _ := sA.EncryptMessage([]byte("msg2"))

	sB.DecryptMessage(seq2, ct2) // Window advanced to 2

	_, err := sB.DecryptMessage(seq0, ct0)
	if err != nil {
		t.Fatalf("Rejected valid out-of-order seq 0 inside window: %v", err)
	}

	_, err = sB.DecryptMessage(seq1, ct1)
	if err != nil {
		t.Fatalf("Rejected valid out-of-order seq 1 inside window: %v", err)
	}

	// Check they are now marked as received
	_, err = sB.DecryptMessage(seq0, ct0)
	if err != ErrReplayDetected {
		t.Fatalf("Failed to detect duplicate after out of order receive")
	}
}

func TestAEADUnauthenticatedPacketDoesNotConsumeWindow(t *testing.T) {
	sA, sB := getTestSessions(t)
	seq1, ct1, _ := sA.EncryptMessage([]byte("msg1"))

	ct1[0] ^= 0xFF               // Corrupt it
	sB.DecryptMessage(seq1, ct1) // Fails authentication

	// Must still be able to receive legitimate seq1 if it arrives
	ct1[0] ^= 0xFF // Fix it
	_, err := sB.DecryptMessage(seq1, ct1)
	if err != nil {
		t.Fatalf("Unauthenticated packet permanently consumed replay state! %v", err)
	}
}

func TestAEADCrossSession(t *testing.T) {
	sA1, _ := getTestSessions(t)
	_, sB2 := getTestSessions(t) // Independent session

	seq, ct, _ := sA1.EncryptMessage([]byte("hello"))
	_, err := sB2.DecryptMessage(seq, ct)
	if err == nil {
		t.Fatalf("Accepted ciphertext from different session")
	}
}

func TestAEADCrossDirection(t *testing.T) {
	sA, _ := getTestSessions(t)
	seq, ct, _ := sA.EncryptMessage([]byte("hello"))

	// Alice attempts to decrypt her own outbound message using receive logic
	// She will expect a different directional marker, AND use a different key.
	// But even if keys matched, the marker would fail.
	_, err := sA.DecryptMessage(seq, ct)
	if err == nil {
		t.Fatalf("Accepted cross-direction ciphertext (reflected)")
	}

	// Explicitly force keys to match to prove AAD protection works independently
	sA.receiveKey = sA.sendKey
	_, err = sA.DecryptMessage(seq, ct)
	if err == nil {
		t.Fatalf("Accepted cross-direction ciphertext when keys were matched! AAD direction failed.")
	}
}

func TestAEADSequenceZero(t *testing.T) {
	sA, sB := getTestSessions(t)
	seq, ct, _ := sA.EncryptMessage([]byte("first message"))
	if seq != 0 {
		t.Fatalf("Expected first sequence to be 0")
	}
	_, err := sB.DecryptMessage(seq, ct)
	if err != nil {
		t.Fatalf("Failed to decrypt sequence 0 properly: %v", err)
	}
}

func TestAEADSequenceExhaustion(t *testing.T) {
	sA, _ := getTestSessions(t)
	sA.sendSeq = math.MaxUint64 - 1

	seq, _, err := sA.EncryptMessage([]byte("almost"))
	if err != nil || seq != math.MaxUint64-1 {
		t.Fatalf("Failed encrypt before bounds")
	}

	_, _, err = sA.EncryptMessage([]byte("exhausted"))
	if err != ErrSequenceExhausted {
		t.Fatalf("Did not return ErrSequenceExhausted on MaxUint64")
	}
}

func TestAEADConcurrentExecution(t *testing.T) {
	sA, sB := getTestSessions(t)
	var wg sync.WaitGroup

	// Concurrently send 1000 messages
	numMessages := 50

	type Msg struct {
		seq uint64
		ct  []byte
	}
	msgChan := make(chan Msg, numMessages)

	wg.Add(numMessages)
	for i := 0; i < numMessages; i++ {
		go func() {
			defer wg.Done()
			seq, ct, _ := sA.EncryptMessage([]byte("concurrent data"))
			msgChan <- Msg{seq: seq, ct: ct}
		}()
	}

	wg.Wait()
	close(msgChan)

	// Concurrently decrypt them
	var decWg sync.WaitGroup
	decWg.Add(numMessages)
	for msg := range msgChan {
		m := msg
		go func() {
			defer decWg.Done()
			_, err := sB.DecryptMessage(m.seq, m.ct)
			if err != nil {
				t.Errorf("Decrypt failed during concurrency test: %v", err)
			}
		}()
	}
	decWg.Wait()

	if sB.highestReceived != 49 {
		t.Fatalf("Highest received did not reach 49: %d", sB.highestReceived)
	}
}

func TestAEADConcurrentReplay(t *testing.T) {
	sA, sB := getTestSessions(t)
	seq, ct, _ := sA.EncryptMessage([]byte("exact same message"))

	var decWg sync.WaitGroup
	numConcurrent := 10
	decWg.Add(numConcurrent)

	successCount := 0
	replayCount := 0
	var mu sync.Mutex

	for i := 0; i < numConcurrent; i++ {
		go func() {
			defer decWg.Done()
			_, err := sB.DecryptMessage(seq, ct)
			mu.Lock()
			if err == nil {
				successCount++
			} else if err == ErrReplayDetected {
				replayCount++
			}
			mu.Unlock()
		}()
	}

	decWg.Wait()

	if successCount != 1 {
		t.Fatalf("Expected exactly 1 successful decryption, got %d", successCount)
	}
	if replayCount != 9 {
		t.Fatalf("Expected exactly 9 ErrReplayDetected, got %d", replayCount)
	}
}

func TestHandshakeCorrelation_SingleMatch(t *testing.T) {
	alice, bob := generateTestIdentities(t)
	hA, _ := NewInitiatorHandshake(alice, bob.PublicKey())
	initMsg, _ := hA.GenerateInit()
	hB, _ := NewResponderHandshake(bob, alice.PublicKey())
	hB.ProcessInit(initMsg)
	respMsg, _ := hB.GenerateResp()

	if !hA.TestResp(respMsg) {
		t.Errorf("SingleMatch failed")
	}
}

func TestHandshakeCorrelation_TwoSimultaneousSameRemote(t *testing.T) {
	alice, bob := generateTestIdentities(t)
	hA1, _ := NewInitiatorHandshake(alice, bob.PublicKey())
	hA2, _ := NewInitiatorHandshake(alice, bob.PublicKey())

	init1, _ := hA1.GenerateInit()
	init2, _ := hA2.GenerateInit()

	hB1, _ := NewResponderHandshake(bob, alice.PublicKey())
	hB1.ProcessInit(init1)
	resp1, _ := hB1.GenerateResp()

	hB2, _ := NewResponderHandshake(bob, alice.PublicKey())
	hB2.ProcessInit(init2)
	resp2, _ := hB2.GenerateResp()

	if !hA1.TestResp(resp1) {
		t.Errorf("expected resp1 to match hA1")
	}
	if hA1.TestResp(resp2) {
		t.Errorf("expected resp2 NOT to match hA1")
	}
	if hA2.TestResp(resp1) {
		t.Errorf("expected resp1 NOT to match hA2")
	}
	if !hA2.TestResp(resp2) {
		t.Errorf("expected resp2 to match hA2")
	}
}

func TestHandshakeCorrelation_MatchCandidateA(t *testing.T) {
	// Covered by TwoSimultaneousSameRemote
}

func TestHandshakeCorrelation_MatchCandidateB(t *testing.T) {
	// Covered by TwoSimultaneousSameRemote
}

func TestHandshakeCorrelation_NoCandidateMatch(t *testing.T) {
	alice, bob := generateTestIdentities(t)
	hA, _ := NewInitiatorHandshake(alice, bob.PublicKey())
	hA.GenerateInit()

	hB, _ := NewResponderHandshake(bob, alice.PublicKey())
	// Let's generate a totally random RESP
	_, eve := generateTestIdentities(t)
	hE, _ := NewInitiatorHandshake(eve, bob.PublicKey())
	initEve, _ := hE.GenerateInit()
	hB.ProcessInit(initEve)
	respMsg, _ := hB.GenerateResp()

	if hA.TestResp(respMsg) {
		t.Errorf("Expected NoCandidateMatch to fail")
	}
}

func TestHandshakeCorrelation_AmbiguousCandidatesFailClosed(t *testing.T) {
	alice, bob := generateTestIdentities(t)
	hA1, _ := NewInitiatorHandshake(alice, bob.PublicKey())
	hA2, _ := NewInitiatorHandshake(alice, bob.PublicKey())

	init1, _ := hA1.GenerateInit()
	hA2.GenerateInit()

	// Force hA2 to use the exact same E_A (simulating a broken RNG or explicit ambiguous state)
	hA2.eAPub = hA1.eAPub

	hB, _ := NewResponderHandshake(bob, alice.PublicKey())
	hB.ProcessInit(init1)
	respMsg, _ := hB.GenerateResp()

	if !hA1.TestResp(respMsg) {
		t.Errorf("Expected hA1 to match")
	}
	if !hA2.TestResp(respMsg) {
		t.Errorf("Expected hA2 to match due to ambiguous E_A. state=%d eAPub=%x, expectedState=%d", hA2.state, hA2.eAPub, InitiatorStateInitSent)
	}
	// This proves that if two candidates match mathematically, the manager must fail closed.
	// Manager test in router_test.go will prove the failure behavior.
}

func TestHandshakeCorrelation_WrongEA(t *testing.T) {
	alice, bob := generateTestIdentities(t)
	hA, _ := NewInitiatorHandshake(alice, bob.PublicKey())
	initMsg, _ := hA.GenerateInit()

	hB, _ := NewResponderHandshake(bob, alice.PublicKey())
	// Modify E_A in initMsg before Bob sees it
	initMsg.EphemeralKey[0] ^= 0xFF
	err := hB.ProcessInit(initMsg)
	if err == nil {
		t.Errorf("Should not process bad E_A init")
	}
}

func TestHandshakeCorrelation_WrongEB(t *testing.T) {
	alice, bob := generateTestIdentities(t)
	hA, _ := NewInitiatorHandshake(alice, bob.PublicKey())
	initMsg, _ := hA.GenerateInit()

	hB, _ := NewResponderHandshake(bob, alice.PublicKey())
	hB.ProcessInit(initMsg)
	respMsg, _ := hB.GenerateResp()

	// Modify E_B in transit
	respMsg.EphemeralKey[0] ^= 0xFF
	if hA.TestResp(respMsg) {
		t.Errorf("WrongEB matched")
	}
}

func TestHandshakeCorrelation_WrongIdentity(t *testing.T) {
	alice, bob := generateTestIdentities(t)
	_, charlie := generateTestIdentities(t)
	hA, _ := NewInitiatorHandshake(alice, bob.PublicKey())
	initMsg, _ := hA.GenerateInit()

	hC, _ := NewResponderHandshake(charlie, alice.PublicKey())
	hC.ProcessInit(initMsg)
	respMsg, _ := hC.GenerateResp()

	if hA.TestResp(respMsg) {
		t.Errorf("WrongIdentity matched")
	}
}

func TestHandshakeCorrelation_MalformedRESP(t *testing.T) {
	alice, bob := generateTestIdentities(t)
	hA, _ := NewInitiatorHandshake(alice, bob.PublicKey())
	hA.GenerateInit()

	if hA.TestResp(&RespMessage{EphemeralKey: make([]byte, 31), Signature: make([]byte, 64)}) {
		t.Errorf("MalformedRESP E_B matched")
	}
	if hA.TestResp(&RespMessage{EphemeralKey: make([]byte, 32), Signature: make([]byte, 63)}) {
		t.Errorf("MalformedRESP Sig matched")
	}
	if hA.TestResp(nil) {
		t.Errorf("MalformedRESP nil matched")
	}
}

func TestHandshakeCorrelation_DuplicateRESP(t *testing.T) {
	alice, bob := generateTestIdentities(t)
	hA, _ := NewInitiatorHandshake(alice, bob.PublicKey())
	initMsg, _ := hA.GenerateInit()

	hB, _ := NewResponderHandshake(bob, alice.PublicKey())
	hB.ProcessInit(initMsg)
	respMsg, _ := hB.GenerateResp()

	hA.TestResp(respMsg)
	err := hA.ProcessResp(respMsg)
	if err != nil {
		t.Errorf("first process failed")
	}

	err = hA.ProcessResp(respMsg)
	if err == nil {
		t.Errorf("duplicate ProcessResp succeeded")
	}
}

func TestHandshakeCorrelation_LateRESP(t *testing.T) {
	alice, bob := generateTestIdentities(t)
	hA, _ := NewInitiatorHandshake(alice, bob.PublicKey())
	initMsg, _ := hA.GenerateInit()

	hB, _ := NewResponderHandshake(bob, alice.PublicKey())
	hB.ProcessInit(initMsg)
	respMsg, _ := hB.GenerateResp()

	hA.ProcessResp(respMsg) // Completes

	if hA.TestResp(respMsg) {
		t.Errorf("LateRESP matched after completion")
	}
}

func TestHandshakeCorrelation_ExactlyOnceRegistration(t *testing.T) {
	alice, bob := generateTestIdentities(t)
	hA, _ := NewInitiatorHandshake(alice, bob.PublicKey())
	initMsg, _ := hA.GenerateInit()

	hB, _ := NewResponderHandshake(bob, alice.PublicKey())
	hB.ProcessInit(initMsg)
	respMsg, _ := hB.GenerateResp()

	hA.ProcessResp(respMsg)

	s1, _ := hA.Session()
	s2, _ := hA.Session()
	if s1 == nil || s2 == nil || s1 != s2 {
		t.Errorf("ExactlyOnceRegistration failed")
	}
}
