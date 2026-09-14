//go:build gui
// +build gui

package main

import (
	"fmt"
	"sync"
	"testing"

	apppkg "github.com/sanchayjain/meshchat/internal/app"
	"time"
)

func TestUIState_Concurrency(t *testing.T) {
	state := NewUIState()
	var wg sync.WaitGroup

	// Simulate concurrent message arrival for Peer A and Peer B
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			peerA := "PeerA"
			peerB := "PeerB"

			state.AddMessage(peerA, Message{Sender: "Peer", Text: fmt.Sprintf("Message %d", id)})
			state.AddMessage(peerB, Message{Sender: "Peer", Text: fmt.Sprintf("Message %d", id)})

			state.GetPeerCount()
			state.GetPeerAt(0)
			state.GetCurrentPeer()
			state.SetCurrentPeer(peerA)
			state.GetConversation(peerA)
		}(i)
	}

	wg.Wait()

	if state.GetPeerCount() != 2 {
		t.Errorf("Expected 2 peers, got %d", state.GetPeerCount())
	}
}

func TestUIState_EnsureConversation(t *testing.T) {
	state := NewUIState()
	peerID := "TestPeerID123"

	// 1. New PeerID creates exactly one conversation
	isNew := state.EnsureConversation(peerID)
	if !isNew {
		t.Fatal("Expected true for new peer ID")
	}

	// 2. Peer appears in peerList
	if state.GetPeerCount() != 1 {
		t.Fatalf("Expected 1 peer, got %d", state.GetPeerCount())
	}
	if state.GetPeerAt(0) != peerID {
		t.Fatalf("Expected %s, got %s", peerID, state.GetPeerAt(0))
	}

	// 3. Conversation starts with zero messages
	conv := state.GetConversation(peerID)
	if conv != "" {
		t.Fatalf("Expected empty conversation, got %s", conv)
	}

	// 4. Second EnsureConversation for same PeerID returns false
	isNew = state.EnsureConversation(peerID)
	if isNew {
		t.Fatal("Expected false for existing peer ID")
	}

	// 5. PeerList contains no duplicate
	if state.GetPeerCount() != 1 {
		t.Fatalf("Expected 1 peer, got %d", state.GetPeerCount())
	}
}

func TestRunEventConsumer_SessionEstablished(t *testing.T) {
	events := make(chan apppkg.AppEvent, 10)
	shutdown := make(chan struct{})
	state := NewUIState()
	var wg sync.WaitGroup

	RunEventConsumer(events, state, shutdown, &wg, nil, nil, nil)

	// 2. Verify EventTypeSessionEstablished causes the conversation to become selectable
	peerID := "Peer1234567890abcdef"
	events <- apppkg.AppEvent{
		Type:   apppkg.EventTypeSessionEstablished,
		PeerID: peerID,
	}

	time.Sleep(50 * time.Millisecond) // Let consumer process

	if state.GetPeerCount() != 1 {
		t.Fatalf("Expected 1 peer, got %d", state.GetPeerCount())
	}
	if state.GetPeerAt(0) != peerID {
		t.Fatalf("Expected %s, got %s", peerID, state.GetPeerAt(0))
	}

	// 3. Verify an outbound message can now be initiated (peer is selectable)
	state.SetCurrentPeer(state.GetPeerAt(0))
	if state.GetCurrentPeer() != peerID {
		t.Fatal("Failed to select peer")
	}

	// 4. Verify EventTypeMessageReceived does not create duplicate
	events <- apppkg.AppEvent{
		Type:    apppkg.EventTypeMessageReceived,
		PeerID:  peerID,
		Message: "Hello",
	}

	time.Sleep(50 * time.Millisecond) // Let consumer process

	if state.GetPeerCount() != 1 {
		t.Fatalf("Expected still 1 peer after message, got %d", state.GetPeerCount())
	}

	close(shutdown)
	wg.Wait()
}
