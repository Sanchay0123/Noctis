//go:build gui
// +build gui

package main

import (
	"fmt"
	"sync"
	"testing"
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
