package app

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/sanchayjain/meshchat/internal/crypto"
	"github.com/sanchayjain/meshchat/internal/mesh"
	"github.com/sanchayjain/meshchat/internal/routing"
	"github.com/sanchayjain/meshchat/internal/session"
)

var (
	ErrNetworkOffline     = errors.New("network offline")
	ErrSessionUnavailable = errors.New("session unavailable")
	ErrMessageRejected    = errors.New("message rejected")
	ErrInvalidIdentity    = errors.New("invalid identity")
)

type ApplicationService interface {
	GetLocalIdentity() string
	DialNode(ctx context.Context, address string, expectedPeerID string) error
	StartConversation(peerID string) error
	SendMessage(peerID string, plaintext string) error
	Start() error
	Stop() error
	SubscribeEvents() <-chan AppEvent
}

type appService struct {
	localIdent  *crypto.NodeIdentity
	meshManager mesh.PeerManager
	sessions    *session.Manager
	router      *routing.Router
	events      chan AppEvent
	done        chan struct{}
}

func NewApplicationService(
	localIdent *crypto.NodeIdentity,
	mManager mesh.PeerManager,
	sManager *session.Manager,
	r *routing.Router,
) ApplicationService {
	s := &appService{
		localIdent:  localIdent,
		meshManager: mManager,
		sessions:    sManager,
		router:      r,
		events:      make(chan AppEvent, 1000), // Bounded buffer
		done:        make(chan struct{}),
	}
	s.router.OnAppData = s.handleAppDataReceived
	s.router.OnSessionEstablished = s.handleSessionEstablished
	return s
}

func (s *appService) GetLocalIdentity() string {
	return hex.EncodeToString(s.localIdent.PublicKey())
}

func (s *appService) StartConversation(peerID string) error {
	pub, err := hex.DecodeString(peerID)
	if err != nil || len(pub) != 32 {
		return ErrInvalidIdentity
	}
	return s.router.StartInitiator(pub)
}

func (s *appService) SendMessage(peerID string, plaintext string) error {
	pub, err := hex.DecodeString(peerID)
	if err != nil || len(pub) != 32 {
		return ErrInvalidIdentity
	}

	session, sessionID, ok := s.sessions.LookupByPeer(peerID)
	if !ok {
		return ErrSessionUnavailable
	}

	seq, ciphertext, err := session.EncryptMessage([]byte(plaintext))
	if err != nil {
		return err
	}

	return s.router.RouteAppData(pub, sessionID, seq, ciphertext)
}

func (s *appService) handleAppDataReceived(sessionID [32]byte, sequenceNum uint64, ciphertext []byte) {
	session, ok := s.sessions.Lookup(sessionID)
	if !ok {
		return
	}

	plaintext, err := session.DecryptMessage(sequenceNum, ciphertext)
	if err != nil {
		return
	}

	peerID := hex.EncodeToString(session.PeerID())
	event := EventMessageReceived(peerID, string(plaintext))

	select {
	case s.events <- event:
	default:
		// Drop event if channel is full to avoid blocking network goroutine
		fmt.Println("[AppService] Warning: Event channel full, dropping message")
	}
}

func (s *appService) Start() error {
	return nil
}

func (s *appService) Stop() error {
	close(s.done)
	// Safe shutdown of event channel
	return nil
}

func (s *appService) SubscribeEvents() <-chan AppEvent {
	return s.events
}

func (s *appService) DialNode(ctx context.Context, address string, expectedPeerID string) error {
	peerPubKey, err := hex.DecodeString(expectedPeerID)
	if err != nil || len(peerPubKey) != 32 {
		return ErrInvalidIdentity
	}
	s.meshManager.ExpectInbound(peerPubKey)
	return s.meshManager.Connect(ctx, address, peerPubKey)
}

func (s *appService) handleSessionEstablished(peerID []byte) {
	peerHex := hex.EncodeToString(peerID)
	event := AppEvent{
		Type:   EventTypeSessionEstablished,
		PeerID: peerHex,
	}
	select {
	case s.events <- event:
	default:
		fmt.Println("[AppService] Warning: Event channel full, dropping session event")
	}
}
