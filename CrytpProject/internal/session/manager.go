package session

import (
	"sync"

	"github.com/sanchayjain/meshchat/internal/crypto"
)

type Manager struct {
	mu       sync.RWMutex
	sessions map[[32]byte]*crypto.Session
}

func NewManager() *Manager {
	return &Manager{
		sessions: make(map[[32]byte]*crypto.Session),
	}
}

func (m *Manager) Register(s *crypto.Session) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var id [32]byte
	copy(id[:], s.ID())
	m.sessions[id] = s
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
	delete(m.sessions, id)
}

func (m *Manager) GetFirstSession() (*crypto.Session, [32]byte, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for id, s := range m.sessions {
		return s, id, true
	}
	return nil, [32]byte{}, false
}
