//go:build gui
// +build gui

package main

import (
	"context"
	"fyne.io/fyne/v2/widget"
	"testing"

	apppkg "github.com/sanchayjain/meshchat/internal/app"
)

type mockAppService struct {
	dialCalled bool
}

func (m *mockAppService) GetLocalIdentity() string                          { return "" }
func (m *mockAppService) GetLocalIdentityBytes() []byte                     { return nil }
func (m *mockAppService) StartConversation(peerID string) error             { return nil }
func (m *mockAppService) SendMessage(peerID string, plaintext string) error { return nil }
func (m *mockAppService) Start() error                                      { return nil }
func (m *mockAppService) Stop() error                                       { return nil }
func (m *mockAppService) SubscribeEvents() <-chan apppkg.AppEvent           { return nil }
func (m *mockAppService) DialNode(ctx context.Context, address string, expectedPeerID string) error {
	m.dialCalled = true
	return nil
}

func TestValidateAndDialPeer_EmptyValidation(t *testing.T) {
	label := widget.NewLabel("")
	mockSvc := &mockAppService{}

	// Test empty address
	ValidateAndDialPeer("", "validpeerid", label, mockSvc)
	if mockSvc.dialCalled {
		t.Fatal("DialNode should not be called when address is empty")
	}
	if label.Text != "Status: Address and PeerID are required" {
		t.Fatalf("Expected validation error message, got: %s", label.Text)
	}

	// Test empty peer id
	label.SetText("")
	ValidateAndDialPeer("127.0.0.1:8100", "  ", label, mockSvc)
	if mockSvc.dialCalled {
		t.Fatal("DialNode should not be called when peer ID is empty (or whitespace)")
	}
	if label.Text != "Status: Address and PeerID are required" {
		t.Fatalf("Expected validation error message, got: %s", label.Text)
	}

	// Test valid
	label.SetText("")
	ValidateAndDialPeer("127.0.0.1:8100", "validpeerid", label, mockSvc)
	if label.Text != "Status: Connecting to 127.0.0.1:8100" {
		t.Fatalf("Expected connection message, got: %s", label.Text)
	}
}
