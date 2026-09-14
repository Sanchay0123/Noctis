//go:build gui
// +build gui

package main

import (
	"sync"
	"testing"
	"time"

	apppkg "github.com/sanchayjain/meshchat/internal/app"
)

func TestGUI_EventConsumerLifecycle(t *testing.T) {
	events := make(chan apppkg.AppEvent)
	shutdown := make(chan struct{})
	state := NewUIState()
	var wg sync.WaitGroup

	RunEventConsumer(events, state, shutdown, &wg, nil, nil, nil)

	finished := make(chan struct{})
	go func() {
		wg.Wait()
		close(finished)
	}()

	select {
	case <-finished:
		t.Fatal("WaitGroup exited before shutdown signal")
	case <-time.After(100 * time.Millisecond):
	}

	close(shutdown)

	select {
	case <-finished:
	case <-time.After(1 * time.Second):
		t.Fatal("Event consumer failed to exit after shutdown signal")
	}
}
