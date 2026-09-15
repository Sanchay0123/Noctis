package mesh

import (
	"bytes"
	"context"
	"sync"
	"sync/atomic"

	"github.com/sanchayjain/meshchat/internal/crypto"
	"github.com/sanchayjain/meshchat/internal/transport"
)

type PeerManager interface {
	Connect(ctx context.Context, endpoint string, expectedIdentity []byte) error
	Listen(addr string) error
	Disconnect(identity []byte) error
	GetPeer(identity []byte) (Peer, error)
	ActivePeers() int
	GetActivePeersSnapshot() []Peer
	Shutdown() error
}

type peerManager struct {
	mu    sync.RWMutex
	peers map[string]*peer

	localIdent *crypto.NodeIdentity
	telemetry  TelemetryRecorder

	dialer   *Dialer
	listener *Listener

	onMsg func(sender []byte, msg []byte)

	shuttingDown atomic.Bool
	wg           sync.WaitGroup

	telemetryQ    chan func()
	quitTelemetry chan struct{}
}

func NewPeerManager(localIdent *crypto.NodeIdentity, telemetry TelemetryRecorder, onMsg func([]byte, []byte)) *peerManager {
	if telemetry == nil {
		telemetry = noopTelemetry{}
	}

	pm := &peerManager{
		peers:         make(map[string]*peer),
		localIdent:    localIdent,
		telemetry:     telemetry,
		onMsg:         onMsg,
		telemetryQ:    make(chan func(), 100), // Bounded non-blocking queue
		quitTelemetry: make(chan struct{}),
	}

	pm.dialer = NewDialer(pm)
	pm.listener = NewListener(pm)

	go pm.telemetryWorker()

	return pm
}

func (pm *peerManager) enqueueTelemetry(fn func()) {
	select {
	case pm.telemetryQ <- fn:
	default:
		// Queue full, drop event to prevent blocking networking
	}
}

func (pm *peerManager) telemetryWorker() {

	for {
		select {
		case fn := <-pm.telemetryQ:
			func() {
				defer func() {
					if r := recover(); r != nil {
						// Recover failing telemetry to prevent crash
					}
				}()
				fn()
			}()
		case <-pm.quitTelemetry:
			return
		}
	}
}

func (pm *peerManager) Connect(ctx context.Context, endpoint string, expectedIdentity []byte) error {
	if pm.shuttingDown.Load() {
		return ErrShutdown
	}
	return pm.dialer.Dial(ctx, endpoint, expectedIdentity)
}

func (pm *peerManager) removeIfCurrent(identity []byte, p *peer) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	idStr := string(identity)
	if existing, ok := pm.peers[idStr]; ok && existing == p {
		delete(pm.peers, idStr)
		if pm.telemetry != nil {
			pm.enqueueTelemetry(func() {
				pm.telemetry.RecordPeerDisconnected()
			})
		}
	}
}

func (pm *peerManager) Disconnect(identity []byte) error {
	pm.mu.RLock()
	p, ok := pm.peers[string(identity)]
	pm.mu.RUnlock()

	if !ok {
		return ErrPeerNotFound
	}

	return p.Close()
}

func (pm *peerManager) GetPeer(identity []byte) (Peer, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	if p, ok := pm.peers[string(identity)]; ok {
		return p, nil
	}
	return nil, ErrPeerNotFound
}

func (pm *peerManager) ActivePeers() int {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return len(pm.peers)
}

func (pm *peerManager) Shutdown() error {
	if !pm.shuttingDown.CompareAndSwap(false, true) {
		return nil // already shutting down
	}

	pm.listener.Close()

	pm.mu.RLock()
	peersToClose := make([]*peer, 0, len(pm.peers))
	for _, p := range pm.peers {
		peersToClose = append(peersToClose, p)
	}
	pm.mu.RUnlock()

	for _, p := range peersToClose {
		p.Close()
	}

	close(pm.quitTelemetry)
	pm.wg.Wait()
	return nil
}

func (pm *peerManager) register(ch *transport.DirectChannel, remoteID []byte, initiator bool) error {
	if pm.shuttingDown.Load() {
		ch.Close()
		return ErrShutdown
	}

	pm.mu.Lock()
	idStr := string(remoteID)
	existing, ok := pm.peers[idStr]

	if !ok && len(pm.peers) >= pm.listener.maxActivePeers {
		pm.mu.Unlock()
		ch.Close()
		pm.enqueueTelemetry(func() {
			pm.telemetry.RecordResourceLimitReached("MaxActivePeers")
		})
		return ErrPeerLimitReached
	}
	if ok {
		cmp := bytes.Compare(pm.localIdent.PublicKey(), remoteID)
		weAreSmaller := cmp < 0
		var keepNew bool
		if weAreSmaller {
			keepNew = initiator
		} else {
			keepNew = !initiator
		}

		if !keepNew {
			pm.mu.Unlock()
			ch.Close()
			pm.enqueueTelemetry(func() {
				pm.telemetry.RecordDuplicateConnection()
			})
			return ErrDuplicateConnection
		} else {
			// ATOMIC REPLACEMENT
			newP := newPeer(ch, remoteID, initiator, pm, pm.onMsg)
			pm.peers[idStr] = newP
			pm.wg.Add(1)
			pm.mu.Unlock()

			// Close incumbent OUTSIDE registry lock
			existing.Close()

			pm.wg.Add(1)
			go func() {
				defer pm.wg.Done()
				newP.writeLoop()
			}()
			go func() {
				defer pm.wg.Done()
				newP.readLoop()
			}()

			pm.enqueueTelemetry(func() {
				pm.telemetry.RecordPeerConnected(initiator)
			})
			return nil
		}
	}

	newP := newPeer(ch, remoteID, initiator, pm, pm.onMsg)
	pm.peers[idStr] = newP
	pm.wg.Add(1)
	pm.mu.Unlock()

	pm.wg.Add(1)
	go func() {
		defer pm.wg.Done()
		newP.writeLoop()
	}()
	go func() {
		defer pm.wg.Done()
		newP.readLoop()
	}()

	pm.enqueueTelemetry(func() {
		pm.telemetry.RecordPeerConnected(initiator)
	})
	return nil
}

func (pm *peerManager) Listen(addr string) error {
	return pm.listener.Listen(addr)
}
func (pm *peerManager) GetActivePeersSnapshot() []Peer {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	peers := make([]Peer, 0, len(pm.peers))
	for _, p := range pm.peers {
		peers = append(peers, p)
	}
	return peers
}
