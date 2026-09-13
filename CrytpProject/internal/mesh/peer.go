package mesh

import (
	"encoding/hex"
	"sync"

	"github.com/sanchayjain/meshchat/internal/transport"
)

type PeerState int

const (
	PeerStateConnecting PeerState = iota
	PeerStateHandshaking
	PeerStateEstablished
	PeerStateClosing
	PeerStateFailed
)

type Peer interface {
	Identity() []byte
	State() PeerState
	Send(msg []byte) error
	EnqueueForward(packet []byte) bool
	Close() error
}

type peer struct {
	mu        sync.RWMutex
	state     PeerState
	initiator bool
	id        []byte
	channel   *transport.DirectChannel
	manager   *peerManager
	onMsg     func(sender []byte, msg []byte)

	queueMu     sync.Mutex
	queue       chan []byte
	queuedCount int
	queuedBytes int

	closeOnce sync.Once
	closed    chan struct{}
}

func newPeer(ch *transport.DirectChannel, id []byte, initiator bool, mgr *peerManager, onMsg func([]byte, []byte)) *peer {
	return &peer{
		state:     PeerStateEstablished,
		initiator: initiator,
		id:        id,
		channel:   ch,
		manager:   mgr,
		onMsg:     onMsg,
		queue:     make(chan []byte, 1000), // Max capacity in channel
		closed:    make(chan struct{}),
	}
}

func (p *peer) Identity() []byte {
	return p.id
}

func (p *peer) String() string {
	return hex.EncodeToString(p.id)
}

func (p *peer) State() PeerState {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.state
}

func (p *peer) setState(newState PeerState) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Terminal states cannot be resurrected
	if p.state == PeerStateClosing || p.state == PeerStateFailed {
		return false
	}
	p.state = newState
	return true
}

func (p *peer) Send(msg []byte) error {
	p.mu.RLock()
	state := p.state
	p.mu.RUnlock()

	if state != PeerStateEstablished {
		return ErrNotEstablished
	}

	if !p.EnqueueForward(msg) {
		return ErrQueueFull
	}
	return nil
}

func (p *peer) EnqueueForward(packet []byte) bool {
	p.mu.RLock()
	state := p.state
	p.mu.RUnlock()

	if state != PeerStateEstablished {
		return false
	}

	p.queueMu.Lock()
	defer p.queueMu.Unlock()

	if p.queuedCount >= 1000 || p.queuedBytes+len(packet) > 2*1024*1024 {
		return false
	}

	p.queuedCount++
	p.queuedBytes += len(packet)

	select {
	case p.queue <- packet:
		return true
	default:
		p.queuedCount--
		p.queuedBytes -= len(packet)
		return false
	}
}

func (p *peer) writeLoop() {
	for {
		select {
		case <-p.closed:
			return
		case msg, ok := <-p.queue:
			if !ok {
				return
			}
			p.queueMu.Lock()
			p.queuedCount--
			p.queuedBytes -= len(msg)
			p.queueMu.Unlock()

			err := p.channel.SendPlaintext(msg)
			if err != nil {
				p.Close()
				return
			}
		}
	}
}

func (p *peer) Close() error {
	var err error
	p.closeOnce.Do(func() {
		if !p.setState(PeerStateClosing) {
			return
		}
		close(p.closed)
		if p.channel != nil {
			err = p.channel.Close()
		}
		p.manager.removeIfCurrent(p.id, p)
	})
	return err
}

func (p *peer) readLoop() {
	defer p.Close()

	for {
		select {
		case <-p.closed:
			return
		default:
		}

		msg, err := p.channel.ReceivePlaintext()
		if err != nil {
			// p.setState(PeerStateFailed) - let Close handle it
			return
		}

		if p.onMsg != nil {
			p.onMsg(p.id, msg)
		}
	}
}
