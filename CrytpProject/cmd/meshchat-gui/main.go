//go:build gui
// +build gui

package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	apppkg "github.com/sanchayjain/meshchat/internal/app"
	"github.com/sanchayjain/meshchat/internal/crypto"
	"github.com/sanchayjain/meshchat/internal/mesh"
	"github.com/sanchayjain/meshchat/internal/routing"
	"github.com/sanchayjain/meshchat/internal/session"
)

type UIState struct {
	mu            sync.RWMutex
	conversations map[string]string
	peerList      []string
	currentPeer   string
}

func NewUIState() *UIState {
	return &UIState{
		conversations: make(map[string]string),
	}
}

func (s *UIState) AddMessage(peerID, message string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	isNew := false
	if _, exists := s.conversations[peerID]; !exists {
		s.peerList = append(s.peerList, peerID)
		isNew = true
	}
	s.conversations[peerID] += message + "\n"
	return isNew
}

func (s *UIState) GetPeerCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.peerList)
}

func (s *UIState) GetPeerAt(index int) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if index >= 0 && index < len(s.peerList) {
		return s.peerList[index]
	}
	return ""
}

func (s *UIState) GetConversation(peerID string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.conversations[peerID]
}

func (s *UIState) SetCurrentPeer(peerID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.currentPeer = peerID
}

func (s *UIState) GetCurrentPeer() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.currentPeer
}

func waitForKey(role string) []byte {
	path := filepath.Join("/shared", role+".pub")
	for {
		b, err := ioutil.ReadFile(path)
		if err == nil && len(b) == 32 {
			return b
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	fyneApp := app.New()
	window := fyneApp.NewWindow("MeshChat")
	window.Resize(fyne.NewSize(800, 600))

	var appService apppkg.ApplicationService
	var myPeerID string
	state := NewUIState()

	statusLabel := widget.NewLabel("Status: Starting...")
	identityLabel := widget.NewLabel("Identity: Unknown")

	chatHistory := widget.NewLabel("Select a conversation to start chatting.")
	chatHistory.Wrapping = fyne.TextWrapWord
	scrollContainer := container.NewVScroll(chatHistory)

	messageInput := widget.NewEntry()
	messageInput.PlaceHolder = "Type a message..."

	var conversationList *widget.List

	refreshChat := func() {
		curr := state.GetCurrentPeer()
		if curr == "" {
			chatHistory.SetText("Select a conversation to start chatting.")
		} else {
			chatHistory.SetText(state.GetConversation(curr))
		}
		scrollContainer.ScrollToBottom()
	}

	conversationList = widget.NewList(
		func() int { return state.GetPeerCount() },
		func() fyne.CanvasObject { return widget.NewLabel("Template...") },
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			p := state.GetPeerAt(id)
			if len(p) > 8 {
				obj.(*widget.Label).SetText(p[:8] + "...")
			}
		},
	)

	conversationList.OnSelected = func(id widget.ListItemID) {
		state.SetCurrentPeer(state.GetPeerAt(id))
		refreshChat()
	}

	sendButton := widget.NewButton("Send", func() {
		text := strings.TrimSpace(messageInput.Text)
		curr := state.GetCurrentPeer()
		if text == "" || curr == "" {
			return
		}

		err := appService.SendMessage(curr, text)
		if err != nil {
			statusLabel.SetText(fmt.Sprintf("Status: Send failed - %v", err))
			return
		}

		state.AddMessage(curr, fmt.Sprintf("You: %s", text))
		messageInput.SetText("")
		refreshChat()
	})

	addPeerButton := widget.NewButton("Add Peer", func() {
		addrEntry := widget.NewEntry()
		addrEntry.PlaceHolder = "127.0.0.1:8100"
		peerEntry := widget.NewEntry()
		peerEntry.PlaceHolder = "Optional expected PeerID (hex)"

		items := []*widget.FormItem{
			widget.NewFormItem("Address", addrEntry),
			widget.NewFormItem("PeerID", peerEntry),
		}

		dialog.ShowForm("Add Peer", "Connect", "Cancel", items, func(b bool) {
			if !b {
				return
			}
			addr := strings.TrimSpace(addrEntry.Text)
			pid := strings.TrimSpace(peerEntry.Text)

			if addr == "" || pid == "" {
				return
			}

			statusLabel.SetText("Status: Connecting to " + addr)

			go func() {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				err := appService.DialNode(ctx, addr, pid)

				if err != nil {
					statusLabel.SetText("Status: Connection failed: " + err.Error())
				} else {
					statusLabel.SetText("Status: Connected to " + addr)
				}
			}()
		}, window)
	})

	leftPanel := container.NewBorder(
		widget.NewLabel("Conversations"),
		addPeerButton,
		nil, nil,
		conversationList,
	)

	rightTop := container.NewVBox(identityLabel, statusLabel)
	rightBottom := container.NewBorder(nil, nil, nil, sendButton, messageInput)

	rightPanel := container.NewBorder(rightTop, rightBottom, nil, nil, scrollContainer)

	split := container.NewHSplit(leftPanel, rightPanel)
	split.Offset = 0.3

	window.SetContent(split)

	shutdown := make(chan struct{})

	go func() {
		role := os.Getenv("MESH_ROLE")
		if role == "" {
			role = "standalone"
		}

		ident, _ := crypto.GenerateIdentity()
		if role != "standalone" {
			ioutil.WriteFile(filepath.Join("/shared", role+".pub"), ident.PublicKey(), 0644)
		}

		myPeerID = hex.EncodeToString(ident.PublicKey())

		identityLabel.SetText("Identity: " + myPeerID[:8] + "...")
		statusLabel.SetText("Status: Disconnected")

		sm := session.NewManager()
		router := routing.NewRouter(ident, nil, sm, nil)
		mgr := mesh.NewPeerManager(ident, nil, router.OnMessage)
		router.SetPeerManager(mgr)

		appSvc := apppkg.NewApplicationService(ident, mgr, sm, router)
		appService = appSvc

		go func() {
			for {
				select {
				case <-shutdown:
					return
				case e, ok := <-appService.SubscribeEvents():
					if !ok {
						return
					}
					switch e.Type {
					case apppkg.EventTypeMessageReceived:
						isNew := state.AddMessage(e.PeerID, fmt.Sprintf("Peer: %s", e.Message))
						if isNew {
							conversationList.Refresh()
						}

						curr := state.GetCurrentPeer()
						if curr == e.PeerID {
							refreshChat()
						} else {
							statusLabel.SetText("Status: New message from " + e.PeerID[:8])
						}
					case apppkg.EventTypePeerConnected:
						statusLabel.SetText("Status: Peer connected " + e.PeerID[:8])
					case apppkg.EventTypePeerDisconnected:
						statusLabel.SetText("Status: Peer disconnected " + e.PeerID[:8])
					case apppkg.EventTypeSecurityAlert:
						statusLabel.SetText("Status: Security Alert - " + e.Message)
					}
				}
			}
		}()

		if role != "standalone" {
			mgr.Listen("0.0.0.0:8000")
			statusLabel.SetText("Status: Listening on 0.0.0.0:8000")
		}
	}()

	window.ShowAndRun()

	close(shutdown)
	if appService != nil {
		appService.Stop()
	}
}
