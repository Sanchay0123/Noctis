package app_test

import (
	"context"
	"encoding/hex"
	"testing"
	"time"

	"github.com/sanchayjain/meshchat/internal/app"
	"github.com/sanchayjain/meshchat/internal/crypto"
	"github.com/sanchayjain/meshchat/internal/mesh"
	"github.com/sanchayjain/meshchat/internal/routing"
	"github.com/sanchayjain/meshchat/internal/session"
)

func setupNode(t *testing.T, address string) (app.ApplicationService, mesh.PeerManager, string) {
	ident, err := crypto.GenerateIdentity()
	if err != nil {
		t.Fatal(err)
	}
	sm := session.NewManager()
	router := routing.NewRouter(ident, nil, sm, nil)
	mgr := mesh.NewPeerManager(ident, nil, router.OnMessage)
	router.SetPeerManager(mgr)

	if address != "" {
		if err := mgr.Listen(address); err != nil {
			t.Fatal(err)
		}
	}

	appService := app.NewApplicationService(ident, mgr, sm, router)
	return appService, mgr, hex.EncodeToString(ident.PublicKey())
}

func decodeHex(t *testing.T, s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestApplicationService_AliceBobCarol_Integration(t *testing.T) {
	aliceApp, _, aliceID := setupNode(t, "127.0.0.1:8100")
	bobApp, bobMgr, bobID := setupNode(t, "127.0.0.1:8101")
	carolApp, carolMgr, carolID := setupNode(t, "127.0.0.1:8102")

	defer func() {
		aliceApp.Stop()
		bobApp.Stop()
		carolApp.Stop()
	}()

	_ = bobMgr // bobMgr.ExpectInbound(decodeHex(t, aliceID))
	_ = carolMgr // carolMgr.ExpectInbound(decodeHex(t, bobID))

	carolEvents := carolApp.SubscribeEvents()
	aliceEvents := aliceApp.SubscribeEvents()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := aliceApp.DialNode(ctx, "127.0.0.1:8101", bobID)
	if err != nil {
		t.Fatalf("Alice failed to dial Bob: %v", err)
	}

	err = bobApp.DialNode(ctx, "127.0.0.1:8102", carolID)
	if err != nil {
		t.Fatalf("Bob failed to dial Carol: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	err = aliceApp.StartConversation(carolID)
	if err != nil {
		t.Fatalf("Alice failed to start conversation with Carol: %v", err)
	}

	sessionEstablished := false
	select {
	case e := <-aliceEvents:
		if e.Type == app.EventTypeSessionEstablished && e.PeerID == carolID {
			sessionEstablished = true
		}
	case <-time.After(2 * time.Second):
	}

	if !sessionEstablished {
		t.Log("Warning: Did not see Alice SessionEstablished event")
	}

	time.Sleep(500 * time.Millisecond)

	msgText := "Hello Carol, it is Alice!"
	err = aliceApp.SendMessage(carolID, msgText)
	if err != nil {
		t.Fatalf("Alice failed to send message to Carol: %v", err)
	}

	var received bool
	timeout := time.After(3 * time.Second)
	for {
		select {
		case e := <-carolEvents:
			if e.Type == app.EventTypeMessageReceived {
				if e.PeerID != aliceID {
					t.Fatalf("Carol expected message from %s, got from %s", aliceID, e.PeerID)
				}
				if e.Message != msgText {
					t.Fatalf("Carol expected message %s, got %s", msgText, e.Message)
				}
				received = true
				break
			}
		case <-timeout:
			t.Fatalf("Carol did not receive the message in time")
		}
		if received {
			break
		}
	}

	if !received {
		t.Fatalf("Carol did not receive the correct message")
	}
}
