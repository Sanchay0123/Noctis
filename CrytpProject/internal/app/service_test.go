package app

import (
	"context"
	"encoding/hex"
	"os"
	"testing"
	"time"

	"github.com/sanchayjain/meshchat/internal/crypto"
	"github.com/sanchayjain/meshchat/internal/mesh"
	"github.com/sanchayjain/meshchat/internal/routing"
	"github.com/sanchayjain/meshchat/internal/session"
)

func createTestSessions(t *testing.T) (*crypto.Session, *crypto.Session, *crypto.NodeIdentity, *crypto.NodeIdentity) {
	identA, _ := crypto.GenerateIdentity()
	identB, _ := crypto.GenerateIdentity()

	hA, err := crypto.NewInitiatorHandshake(identA, identB.PublicKey())
	if err != nil {
		t.Fatal(err)
	}

	msgA, err := hA.GenerateInit()
	if err != nil {
		t.Fatal(err)
	}

	hB, err := crypto.NewResponderHandshake(identB, identA.PublicKey())
	if err != nil {
		t.Fatal(err)
	}

	err = hB.ProcessInit(msgA)
	if err != nil {
		t.Fatal(err)
	}

	msgB, err := hB.GenerateResp()
	if err != nil {
		t.Fatal(err)
	}

	err = hA.ProcessResp(msgB)
	if err != nil {
		t.Fatal(err)
	}

	sA, err := hA.Session()
	if err != nil {
		t.Fatal(err)
	}

	sB, err := hB.Session()
	if err != nil {
		t.Fatal(err)
	}

	return sA, sB, identA, identB
}

func TestApplicationService_SendMessage_NoSession(t *testing.T) {
	ident, _ := crypto.GenerateIdentity()
	sm := session.NewManager()
	router := routing.NewRouter(ident, nil, sm, nil)
	mgr := mesh.NewPeerManager(ident, nil, router.OnMessage)
	router.SetPeerManager(mgr)

	appService := NewApplicationService(ident, mgr, sm, router)

	peerID := "0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20"
	err := appService.SendMessage(peerID, "test message")
	if err != ErrSessionUnavailable {
		t.Fatalf("expected ErrSessionUnavailable, got %v", err)
	}
}

func TestApplicationService_SendMessage_InvalidIdentity(t *testing.T) {
	ident, _ := crypto.GenerateIdentity()
	sm := session.NewManager()
	router := routing.NewRouter(ident, nil, sm, nil)
	mgr := mesh.NewPeerManager(ident, nil, router.OnMessage)
	router.SetPeerManager(mgr)

	appService := NewApplicationService(ident, mgr, sm, router)

	err := appService.SendMessage("invalidhex", "test")
	if err != ErrInvalidIdentity {
		t.Fatalf("expected ErrInvalidIdentity, got %v", err)
	}
}

func TestApplicationService_DialNode_InvalidIdentity(t *testing.T) {
	ident, _ := crypto.GenerateIdentity()
	sm := session.NewManager()
	router := routing.NewRouter(ident, nil, sm, nil)
	mgr := mesh.NewPeerManager(ident, nil, router.OnMessage)
	router.SetPeerManager(mgr)

	appService := NewApplicationService(ident, mgr, sm, router)

	err := appService.DialNode(context.Background(), "127.0.0.1:8100", "invalidhex")
	if err != ErrInvalidIdentity {
		t.Fatalf("expected ErrInvalidIdentity, got %v", err)
	}
}

func TestApplicationService_DialNode_Success(t *testing.T) {
	ident, _ := crypto.GenerateIdentity()
	sm := session.NewManager()
	router := routing.NewRouter(ident, nil, sm, nil)
	mgr := mesh.NewPeerManager(ident, nil, router.OnMessage)
	router.SetPeerManager(mgr)

	appService := NewApplicationService(ident, mgr, sm, router)

	peerID := "0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20"

	err := appService.DialNode(context.Background(), "127.0.0.1:1", peerID)
	if err == nil || err == ErrInvalidIdentity {
		t.Fatalf("expected a connection error, got %v", err)
	}
}

func TestApplicationService_StartConversation(t *testing.T) {
	ident, _ := crypto.GenerateIdentity()
	sm := session.NewManager()
	router := routing.NewRouter(ident, nil, sm, nil)
	mgr := mesh.NewPeerManager(ident, nil, router.OnMessage)
	router.SetPeerManager(mgr)

	appService := NewApplicationService(ident, mgr, sm, router)

	peerID := "0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20"
	err := appService.StartConversation(peerID)
	if err != nil {
		t.Fatalf("StartConversation failed: %v", err)
	}

	err = appService.StartConversation("invalidhex")
	if err != ErrInvalidIdentity {
		t.Fatalf("expected ErrInvalidIdentity, got %v", err)
	}
}

func TestApplicationService_SendMessage_Success(t *testing.T) {
	sA, _, identA, identB := createTestSessions(t)

	sm := session.NewManager()
	sm.Register(sA)

	router := routing.NewRouter(identA, nil, sm, nil)
	mgr := mesh.NewPeerManager(identA, nil, router.OnMessage)
	router.SetPeerManager(mgr)

	appService := NewApplicationService(identA, mgr, sm, router)

	peerIDHex := hex.EncodeToString(identB.PublicKey())
	err := appService.SendMessage(peerIDHex, "Hello World")
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
}

func TestApplicationService_EventMessageReceived(t *testing.T) {
	sA, sB, identA, identB := createTestSessions(t)

	sm := session.NewManager()
	sm.Register(sB) // Receiver manager has sB

	router := routing.NewRouter(identB, nil, sm, nil)
	mgr := mesh.NewPeerManager(identB, nil, router.OnMessage)
	router.SetPeerManager(mgr)

	appService := NewApplicationService(identB, mgr, sm, router).(*appService)

	msg := []byte("Hello Event")
	seq, ct, err := sA.EncryptMessage(msg) // Encrypt with sender session
	if err != nil {
		t.Fatal(err)
	}

	// Simulate Router calling OnAppData
	go func() {
		var id [32]byte
		copy(id[:], sB.ID())
		appService.handleAppDataReceived(id, seq, ct)
	}()

	select {
	case event := <-appService.SubscribeEvents():
		if event.Type != EventTypeMessageReceived {
			t.Fatalf("expected EventTypeMessageReceived, got %v", event.Type)
		}
		if event.Message != "Hello Event" {
			t.Fatalf("expected message 'Hello Event', got '%s'", event.Message)
		}
		if event.PeerID != hex.EncodeToString(identA.PublicKey()) {
			t.Fatalf("expected peer ID %x, got %s", identA.PublicKey(), event.PeerID)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for event")
	}
}

func TestApplicationService_Shutdown(t *testing.T) {
	ident, _ := crypto.GenerateIdentity()
	sm := session.NewManager()
	router := routing.NewRouter(ident, nil, sm, nil)
	mgr := mesh.NewPeerManager(ident, nil, router.OnMessage)
	router.SetPeerManager(mgr)

	appService := NewApplicationService(ident, mgr, sm, router)
	err := appService.Stop()
	if err != nil {
		t.Fatal(err)
	}
}

func TestApplicationService_NewNode_ConfigurableListener(t *testing.T) {
	// Test configured listener
	os.Setenv("MESH_LISTEN_ADDR", "127.0.0.1:8088")
	defer os.Unsetenv("MESH_LISTEN_ADDR")

	node1, err := NewNode("peer")
	if err != nil {
		t.Fatalf("Failed to create node: %v", err)
	}
	defer node1.Stop()

	// Verify we can connect to 127.0.0.1:8088
	node2, err := NewNode("standalone")
	if err != nil {
		t.Fatalf("Failed to create standalone node: %v", err)
	}
	defer node2.Stop()

	// node1 must expect node2's inbound connection
	pub2, _ := hex.DecodeString(node2.GetLocalIdentity())
	node1.(*appService).meshManager.ExpectInbound(pub2)
	time.Sleep(50 * time.Millisecond) // Wait for listener to bind

	peerIDHex := node1.GetLocalIdentity()
	err = node2.DialNode(context.Background(), "127.0.0.1:8088", peerIDHex)
	if err != nil {
		t.Fatalf("Expected connection to succeed on configured port, got %v", err)
	}
}

func TestApplicationService_NewNode_DefaultListener(t *testing.T) {
	// Ensure unset
	os.Unsetenv("MESH_LISTEN_ADDR")

	node1, err := NewNode("peer")
	if err != nil {
		t.Fatalf("Failed to create node: %v", err)
	}
	defer node1.Stop()

	// Verify we can connect to 0.0.0.0:8000
	node2, err := NewNode("standalone")
	if err != nil {
		t.Fatalf("Failed to create standalone node: %v", err)
	}
	defer node2.Stop()

	pub2, _ := hex.DecodeString(node2.GetLocalIdentity())
	node1.(*appService).meshManager.ExpectInbound(pub2)
	time.Sleep(50 * time.Millisecond) // Wait for listener to bind

	peerIDHex := node1.GetLocalIdentity()
	err = node2.DialNode(context.Background(), "127.0.0.1:8000", peerIDHex)
	if err != nil {
		t.Fatalf("Expected connection to succeed on default port, got %v", err)
	}
}
