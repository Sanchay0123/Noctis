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
	Send(plaintext []byte) error
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

func (p *peer) Send(plaintext []byte) error {
	p.mu.RLock()
	state := p.state
	p.mu.RUnlock()

	if state != PeerStateEstablished {
		return ErrNotEstablished
	}

	return p.channel.SendPlaintext(plaintext)
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
