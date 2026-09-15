package mesh_test

import (
	"context"
	"encoding/hex"
	"fmt"
	"google.golang.org/protobuf/proto"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sanchayjain/meshchat/internal/crypto"
	"github.com/sanchayjain/meshchat/internal/mesh"
	"github.com/sanchayjain/meshchat/internal/protocol"
	"github.com/sanchayjain/meshchat/internal/routing"
	"github.com/sanchayjain/meshchat/internal/session"
)

func TestRelayCannotDecryptEndpointSession(t *testing.T) {
	// 1. Establish Alice <-> Bob <-> Carol topology locally
	aliceIdent, _ := crypto.GenerateIdentity()
	bobIdent, _ := crypto.GenerateIdentity()
	carolIdent, _ := crypto.GenerateIdentity()

	aliceSM := session.NewManager()
	bobSM := session.NewManager()
	carolSM := session.NewManager()

	aliceOnMsg := func(sessionID [32]byte, seq uint64, ciphertext []byte) {}
	bobOnMsg := func(sessionID [32]byte, seq uint64, ciphertext []byte) {}
	carolOnMsg := func(sessionID [32]byte, seq uint64, ciphertext []byte) {
		if string(ciphertext) == "Hello Carol" {
			t.Errorf("Router decrypted the payload! Architectural violation.")
		}

		session, ok := carolSM.Lookup(sessionID)
		if !ok {
			t.Errorf("Carol session not found")
			return
		}
		plaintext, err := session.DecryptMessage(seq, ciphertext)
		if err != nil {
			t.Errorf("carol decryption failed: %v", err)
			return
		}
		if string(plaintext) != "Hello Carol" {
			t.Errorf("carol received wrong message")
		}
	}

	aliceRouter := routing.NewRouter(aliceIdent, nil, aliceSM, nil)
	aliceRouter.OnAppData = aliceOnMsg
	bobRouter := routing.NewRouter(bobIdent, nil, bobSM, nil)
	bobRouter.OnAppData = bobOnMsg

	// M9.1: Test-only instrumentation for Relay Blindness Evidence
	var relayObservationCount atomic.Int32
	bobOnMsgWrapper := func(peerID []byte, rawMsg []byte) {
		var pkt protocol.MeshPacket
		if err := proto.Unmarshal(rawMsg, &pkt); err == nil {
			fmt.Printf("M9.1 EVIDENCE - RELAY PACKET OBSERVATION at Bob\n")
			fmt.Printf("  Packet Type: %v\n", pkt.Type)
			fmt.Printf("  SourceNode : %x...\n", pkt.SourceNode[:4])
			fmt.Printf("  DestNode   : %x...\n", pkt.DestNode[:4])
			fmt.Printf("  PacketID   : %x...\n", pkt.PacketId[:4])
			fmt.Printf("  TTL        : %d\n", pkt.Ttl)

			if appData, ok := pkt.Payload.(*protocol.MeshPacket_AppData); ok {
				fmt.Printf("  Payload    : APP_DATA\n")
				fmt.Printf("  SessionID  : %x...\n", appData.AppData.SessionId[:4])
				fmt.Printf("  SequenceNum: %d\n", appData.AppData.SequenceNum)
				fmt.Printf("  Ciphertext : %d bytes observed\n", len(appData.AppData.Ciphertext))

				// Verify Bob doesn't have the session
				if _, ok := bobSM.Lookup(*(*[32]byte)(appData.AppData.SessionId)); ok {
					t.Errorf("CRITICAL VIOLATION: Bob possesses the endpoint session!")
				} else {
					fmt.Printf("  Validation : Bob's implemented session manager does not contain the Alice-Carol endpoint application session required to decrypt this ciphertext.\n")
				}
				relayObservationCount.Add(1)
			}
		}
		bobRouter.OnMessage(peerID, rawMsg)
	}
	carolRouter := routing.NewRouter(carolIdent, nil, carolSM, nil)
	carolRouter.OnAppData = carolOnMsg

	alicePM := mesh.NewPeerManager(aliceIdent, nil, aliceRouter.OnMessage)
	bobPM := mesh.NewPeerManager(bobIdent, nil, bobOnMsgWrapper)
	carolPM := mesh.NewPeerManager(carolIdent, nil, carolRouter.OnMessage)

	aliceRouter.SetPeerManager(alicePM)
	bobRouter.SetPeerManager(bobPM)
	carolRouter.SetPeerManager(carolPM)

	bobPM.Listen("127.0.0.1:8200")
	carolPM.Listen("127.0.0.1:8201")

	time.Sleep(100 * time.Millisecond)

	_ = alicePM // alicePM.ExpectInbound(bobIdent.PublicKey())
	_ = bobPM   // bobPM.ExpectInbound(aliceIdent.PublicKey())
	_ = bobPM   // bobPM.ExpectInbound(carolIdent.PublicKey())
	_ = carolPM // carolPM.ExpectInbound(bobIdent.PublicKey())

	alicePM.Connect(context.Background(), "127.0.0.1:8200", bobIdent.PublicKey())
	bobPM.Connect(context.Background(), "127.0.0.1:8201", carolIdent.PublicKey())

	time.Sleep(200 * time.Millisecond)

	// Alice initiates handshake to Carol
	aliceRouter.StartInitiator(carolIdent.PublicKey())
	time.Sleep(500 * time.Millisecond)

	// Verify Alice and Carol have endpoint sessions, Bob does NOT have Alice<->Carol session
	if _, _, ok := aliceSM.LookupByPeer(hex.EncodeToString(carolIdent.PublicKey())); !ok {
		t.Errorf("Alice SM has no session")
	}
	if _, _, ok := carolSM.LookupByPeer(hex.EncodeToString(aliceIdent.PublicKey())); !ok {
		t.Errorf("Carol SM has no session")
	}
	if _, _, ok := bobSM.LookupByPeer(hex.EncodeToString(carolIdent.PublicKey())); ok {
		t.Errorf("Bob should have NO endpoint sessions")
	}

	// Send APP_DATA
	t.Logf("M9.1 EVIDENCE - SESSION OWNERSHIP")
	t.Logf("  Alice possesses Alice-Carol endpoint application session.")
	t.Logf("  Carol possesses Alice-Carol endpoint application session.")
	t.Logf("  Alice encrypting application plaintext 'Hello Carol' using existing ChaCha20-Poly1305 implementation...")
	sendAppDataHelper(aliceRouter, aliceSM, carolIdent.PublicKey(), []byte("Hello Carol"))
	time.Sleep(200 * time.Millisecond)

	if relayObservationCount.Load() == 0 {
		t.Errorf("Relay blindness evidence failed: Bob never observed the forwarded APP_DATA packet")
	}

	alicePM.Shutdown()
	bobPM.Shutdown()
	carolPM.Shutdown()
}

func TestPeerLossDoesNotInvalidateEndpointSession(t *testing.T) {
	aliceIdent, _ := crypto.GenerateIdentity()
	bobIdent, _ := crypto.GenerateIdentity()
	carolIdent, _ := crypto.GenerateIdentity()

	aliceSM := session.NewManager()
	bobSM := session.NewManager()
	carolSM := session.NewManager()

	aliceOnMsg := func(sessionID [32]byte, seq uint64, ciphertext []byte) {}
	bobOnMsg := func(sessionID [32]byte, seq uint64, ciphertext []byte) {}
	carolOnMsg := func(sessionID [32]byte, seq uint64, ciphertext []byte) {}

	aliceRouter := routing.NewRouter(aliceIdent, nil, aliceSM, nil)
	aliceRouter.OnAppData = aliceOnMsg
	bobRouter := routing.NewRouter(bobIdent, nil, bobSM, nil)
	bobRouter.OnAppData = bobOnMsg

	// M9.1: Test-only instrumentation for Relay Blindness Evidence
	var relayObservationCount atomic.Int32
	bobOnMsgWrapper := func(peerID []byte, rawMsg []byte) {
		var pkt protocol.MeshPacket
		if err := proto.Unmarshal(rawMsg, &pkt); err == nil {
			fmt.Printf("M9.1 EVIDENCE - RELAY PACKET OBSERVATION at Bob\n")
			fmt.Printf("  Packet Type: %v\n", pkt.Type)
			fmt.Printf("  SourceNode : %x...\n", pkt.SourceNode[:4])
			fmt.Printf("  DestNode   : %x...\n", pkt.DestNode[:4])
			fmt.Printf("  PacketID   : %x...\n", pkt.PacketId[:4])
			fmt.Printf("  TTL        : %d\n", pkt.Ttl)

			if appData, ok := pkt.Payload.(*protocol.MeshPacket_AppData); ok {
				fmt.Printf("  Payload    : APP_DATA\n")
				fmt.Printf("  SessionID  : %x...\n", appData.AppData.SessionId[:4])
				fmt.Printf("  SequenceNum: %d\n", appData.AppData.SequenceNum)
				fmt.Printf("  Ciphertext : %d bytes observed\n", len(appData.AppData.Ciphertext))

				// Verify Bob doesn't have the session
				if _, ok := bobSM.Lookup(*(*[32]byte)(appData.AppData.SessionId)); ok {
					t.Errorf("CRITICAL VIOLATION: Bob possesses the endpoint session!")
				} else {
					fmt.Printf("  Validation : Bob's implemented session manager does not contain the Alice-Carol endpoint application session required to decrypt this ciphertext.\n")
				}
				relayObservationCount.Add(1)
			}
		}
		bobRouter.OnMessage(peerID, rawMsg)
	}
	carolRouter := routing.NewRouter(carolIdent, nil, carolSM, nil)
	carolRouter.OnAppData = carolOnMsg

	alicePM := mesh.NewPeerManager(aliceIdent, nil, aliceRouter.OnMessage)
	bobPM := mesh.NewPeerManager(bobIdent, nil, bobOnMsgWrapper)
	carolPM := mesh.NewPeerManager(carolIdent, nil, carolRouter.OnMessage)

	aliceRouter.SetPeerManager(alicePM)
	bobRouter.SetPeerManager(bobPM)
	carolRouter.SetPeerManager(carolPM)

	bobPM.Listen("127.0.0.1:8200")
	carolPM.Listen("127.0.0.1:8201")

	time.Sleep(100 * time.Millisecond)

	_ = alicePM // alicePM.ExpectInbound(bobIdent.PublicKey())
	_ = bobPM   // bobPM.ExpectInbound(aliceIdent.PublicKey())
	_ = bobPM   // bobPM.ExpectInbound(carolIdent.PublicKey())
	_ = carolPM // carolPM.ExpectInbound(bobIdent.PublicKey())

	alicePM.Connect(context.Background(), "127.0.0.1:8200", bobIdent.PublicKey())
	bobPM.Connect(context.Background(), "127.0.0.1:8201", carolIdent.PublicKey())

	time.Sleep(200 * time.Millisecond)

	aliceRouter.StartInitiator(carolIdent.PublicKey())
	time.Sleep(500 * time.Millisecond)

	if _, _, ok := aliceSM.LookupByPeer(hex.EncodeToString(carolIdent.PublicKey())); !ok {
		t.Fatalf("Session failed to establish on Alice")
	}
	if _, _, ok := carolSM.LookupByPeer(hex.EncodeToString(aliceIdent.PublicKey())); !ok {
		t.Fatalf("Session failed to establish on Carol")
	}

	// Close Bob
	bobPM.Shutdown()
	time.Sleep(200 * time.Millisecond)

	// Confirm SessionManager still has Alice <-> Carol
	if _, _, ok := aliceSM.LookupByPeer(hex.EncodeToString(carolIdent.PublicKey())); !ok {
		t.Errorf("Alice SM lost session after peer loss")
	}
	if _, _, ok := carolSM.LookupByPeer(hex.EncodeToString(aliceIdent.PublicKey())); !ok {
		t.Errorf("Carol SM lost session after peer loss")
	}

	alicePM.Shutdown()
	carolPM.Shutdown()
}

type spyPM struct {
	mesh.PeerManager
	capturedAfter chan []byte
}

func (s *spyPM) GetActivePeersSnapshot() []mesh.Peer {
	peers := s.PeerManager.GetActivePeersSnapshot()
	var spied []mesh.Peer
	for _, p := range peers {
		spied = append(spied, &spyPeer{Peer: p, pm: s})
	}
	return spied
}

type spyPeer struct {
	mesh.Peer
	pm *spyPM
}

func (s *spyPeer) EnqueueForward(packet []byte) bool {
	select {
	case s.pm.capturedAfter <- append([]byte(nil), packet...):
	default:
	}
	return s.Peer.EnqueueForward(packet)
}

func TestRelayCiphertextBoundary(t *testing.T) {
	aliceIdent, _ := crypto.GenerateIdentity()
	bobIdent, _ := crypto.GenerateIdentity()
	carolIdent, _ := crypto.GenerateIdentity()

	aliceSM := session.NewManager()
	bobSM := session.NewManager()
	carolSM := session.NewManager()

	aliceOnMsg := func(sessionID [32]byte, seq uint64, ciphertext []byte) {}

	capturedBeforeBob := make(chan []byte, 100)

	carolOnMsg := func(sessionID [32]byte, seq uint64, ciphertext []byte) {}

	aliceRouter := routing.NewRouter(aliceIdent, nil, aliceSM, nil)
	aliceRouter.OnAppData = aliceOnMsg
	bobRouter := routing.NewRouter(bobIdent, nil, bobSM, nil)
	bobRouter.OnAppData = func(sessionID [32]byte, seq uint64, ciphertext []byte) {}
	carolRouter := routing.NewRouter(carolIdent, nil, carolSM, nil)
	carolRouter.OnAppData = carolOnMsg

	alicePM := mesh.NewPeerManager(aliceIdent, nil, aliceRouter.OnMessage)

	bobInterceptor := func(sender []byte, msg []byte) {
		select {
		case capturedBeforeBob <- append([]byte(nil), msg...):
		default:
		}
		bobRouter.OnMessage(sender, msg)
	}

	baseBobPM := mesh.NewPeerManager(bobIdent, nil, bobInterceptor)
	bobSpyPM := &spyPM{PeerManager: baseBobPM, capturedAfter: make(chan []byte, 100)}

	carolPM := mesh.NewPeerManager(carolIdent, nil, carolRouter.OnMessage)

	aliceRouter.SetPeerManager(alicePM)
	bobRouter.SetPeerManager(bobSpyPM)
	carolRouter.SetPeerManager(carolPM)

	baseBobPM.Listen("127.0.0.1:8300")
	carolPM.Listen("127.0.0.1:8301")
	time.Sleep(100 * time.Millisecond)

	_ = alicePM   // alicePM.ExpectInbound(bobIdent.PublicKey())
	_ = baseBobPM // baseBobPM.ExpectInbound(aliceIdent.PublicKey())
	_ = baseBobPM // baseBobPM.ExpectInbound(carolIdent.PublicKey())
	_ = carolPM   // carolPM.ExpectInbound(bobIdent.PublicKey())

	alicePM.Connect(context.Background(), "127.0.0.1:8300", bobIdent.PublicKey())
	baseBobPM.Connect(context.Background(), "127.0.0.1:8301", carolIdent.PublicKey())
	time.Sleep(200 * time.Millisecond)

	aliceRouter.StartInitiator(carolIdent.PublicKey())
	time.Sleep(500 * time.Millisecond)

	// Drain channels before sending APP_DATA
	for len(capturedBeforeBob) > 0 {
		<-capturedBeforeBob
	}
	for len(bobSpyPM.capturedAfter) > 0 {
		<-bobSpyPM.capturedAfter
	}

	// Send APP_DATA
	sendAppDataHelper(aliceRouter, aliceSM, carolIdent.PublicKey(), []byte("SECRET_RELAY_TEST"))
	time.Sleep(200 * time.Millisecond)

	if len(capturedBeforeBob) == 0 {
		t.Fatalf("Bob never received anything")
	}
	if len(bobSpyPM.capturedAfter) == 0 {
		t.Fatalf("Bob never forwarded anything")
	}

	// The last message is the APP_DATA
	var lastBeforeBob []byte
	for len(capturedBeforeBob) > 0 {
		lastBeforeBob = <-capturedBeforeBob
	}
	var capturedAfterBob []byte
	for len(bobSpyPM.capturedAfter) > 0 {
		capturedAfterBob = <-bobSpyPM.capturedAfter
	}

	// They must match exactly except for TTL byte which is at index 21 (assuming standard proto serialization)
	// We can just verify lengths are equal, and if we flip the TTL back, bytes are equal.
	if len(lastBeforeBob) != len(capturedAfterBob) {
		t.Fatalf("Ciphertext length changed during relay")
	}

	diffCount := 0
	for i := 0; i < len(lastBeforeBob); i++ {
		if lastBeforeBob[i] != capturedAfterBob[i] {
			diffCount++
		}
	}
	if diffCount > 1 {
		t.Fatalf("Ciphertext altered in more than 1 byte (TTL). Diffs: %d", diffCount)
	}

	if _, _, ok := bobSM.LookupByPeer(hex.EncodeToString(carolIdent.PublicKey())); ok {
		t.Fatalf("Bob possesses endpoint M4 session state")
	}

	alicePM.Shutdown()
	baseBobPM.Shutdown()
	carolPM.Shutdown()
}

func sendAppDataHelper(r *routing.Router, sm *session.Manager, destPubKey []byte, plaintext []byte) error {
	session, id, ok := sm.LookupByPeer(hex.EncodeToString(destPubKey))
	if !ok {
		return fmt.Errorf("no session")
	}
	seq, ciphertext, err := session.EncryptMessage(plaintext)
	if err != nil {
		return err
	}
	return r.RouteAppData(destPubKey, id, seq, ciphertext)
}
