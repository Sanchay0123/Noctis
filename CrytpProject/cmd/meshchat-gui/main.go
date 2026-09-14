//go:build gui
// +build gui

package main

import (
	"context"
	"fmt"
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
)

type Message struct {
	Sender    string
	Text      string
	Timestamp time.Time
}

type UIState struct {
	mu            sync.RWMutex
	conversations map[string][]Message
	peerList      []string
	currentPeer   string
}

func NewUIState() *UIState {
	return &UIState{
		conversations: make(map[string][]Message),
	}
}

func (s *UIState) AddMessage(peerID string, msg Message) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	isNew := false
	if _, exists := s.conversations[peerID]; !exists {
		s.peerList = append(s.peerList, peerID)
		isNew = true
	}
	s.conversations[peerID] = append(s.conversations[peerID], msg)
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
	msgs := s.conversations[peerID]
	var sb strings.Builder
	for _, m := range msgs {
		sb.WriteString(fmt.Sprintf("[%s] %s: %s\n", m.Timestamp.Format("15:04:05"), m.Sender, m.Text))
	}
	return sb.String()
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

func main() {
	fyneApp := app.New()
	window := fyneApp.NewWindow("MeshChat")
	window.Resize(fyne.NewSize(800, 600))

	var appService apppkg.ApplicationService
	state := NewUIState()

	statusLabel := widget.NewLabel("Status: Starting...")
	identityLabel := widget.NewLabel("Identity: Unknown")

	copyIdBtn := widget.NewButton("Copy ID", func() {
		if appService != nil {
			window.Clipboard().SetContent(appService.GetLocalIdentity())
		}
	})

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

		state.AddMessage(curr, Message{
			Sender:    "You",
			Text:      text,
			Timestamp: time.Now(),
		})
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
					statusLabel.SetText("Status: Connected to " + addr + ". Handshaking...")
					errStart := appService.StartConversation(pid)
					if errStart != nil {
						statusLabel.SetText("Status: Handshake init failed: " + errStart.Error())
					}
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

	identBox := container.NewHBox(identityLabel, copyIdBtn)
	rightTop := container.NewVBox(identBox, statusLabel)
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

		appSvc, err := apppkg.NewNode(role)
		if err != nil {
			statusLabel.SetText("Status: Initialization failed")
			return
		}
		appService = appSvc

		if role != "standalone" {
			os.WriteFile(filepath.Join("/shared", role+".pub"), appService.GetLocalIdentityBytes(), 0644)
		}

		myPeerID := appService.GetLocalIdentity()
		identityLabel.SetText("Identity: " + myPeerID[:8] + "...")
		statusLabel.SetText("Status: Disconnected")
		appService.Start()

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
						isNew := state.AddMessage(e.PeerID, Message{
							Sender:    "Peer",
							Text:      e.Message,
							Timestamp: time.Now(),
						})
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
					case apppkg.EventTypeSessionEstablished:
						statusLabel.SetText("Status: Secure Session Established with " + e.PeerID[:8])
					case apppkg.EventTypeConnectionFailed:
						statusLabel.SetText("Status: Connection Failed for " + e.PeerID[:8])
					case apppkg.EventTypeSessionFailed:
						statusLabel.SetText("Status: Authentication Failed for " + e.PeerID[:8])
					case apppkg.EventTypeSecurityAlert:
						statusLabel.SetText("Status: Security Alert - " + e.Message)
					}
				}
			}
		}()
	}()

	window.ShowAndRun()

	close(shutdown)
	if appService != nil {
		appService.Stop()
	}
}
