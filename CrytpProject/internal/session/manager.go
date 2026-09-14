package session

import (
	"encoding/hex"
	"sync"

	"github.com/sanchayjain/meshchat/internal/crypto"
)

type Manager struct {
	mu       sync.RWMutex
	sessions map[[32]byte]*crypto.Session
	peerMap  map[string][32]byte
}

func NewManager() *Manager {
	return &Manager{
		sessions: make(map[[32]byte]*crypto.Session),
		peerMap:  make(map[string][32]byte),
	}
}

func (m *Manager) Register(s *crypto.Session) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var id [32]byte
	copy(id[:], s.ID())
	m.sessions[id] = s

	peerHex := hex.EncodeToString(s.PeerID())
	m.peerMap[peerHex] = id
}

func (m *Manager) Lookup(id [32]byte) (*crypto.Session, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.sessions[id]
	return s, ok
}

func (m *Manager) Remove(id [32]byte) {
	m.mu.Lock()
	defer m.mu.Unlock()
	s, ok := m.sessions[id]
	if ok {
		peerHex := hex.EncodeToString(s.PeerID())
		if m.peerMap[peerHex] == id {
			delete(m.peerMap, peerHex)
		}
		delete(m.sessions, id)
	}
}

func (m *Manager) LookupByPeer(peerID string) (*crypto.Session, [32]byte, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	id, ok := m.peerMap[peerID]
	if !ok {
		return nil, [32]byte{}, false
	}
	s, ok := m.sessions[id]
	return s, id, ok
}
