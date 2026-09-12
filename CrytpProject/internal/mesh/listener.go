package mesh

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"sync"
	"time"

	"github.com/sanchayjain/meshchat/internal/protocol"
	"github.com/sanchayjain/meshchat/internal/transport"
	"google.golang.org/protobuf/proto"
)

const MaxFrameSize = 65536

type replayConn struct {
	net.Conn
	r io.Reader
}

func (c *replayConn) Read(b []byte) (int, error) {
	return c.r.Read(b)
}

type Listener struct {
	pm                   *peerManager
	listener             net.Listener
	
	maxActivePeers       int
	maxPendingHandshakes int
	handshakeTimeout     time.Duration
	
	pendingSem           chan struct{}
	
	mu                   sync.Mutex
	closed               bool
}

func NewListener(pm *peerManager) *Listener {
	return &Listener{
		pm:                   pm,
		maxActivePeers:       50,
		maxPendingHandshakes: 10,
		handshakeTimeout:     5 * time.Second,
		pendingSem:           make(chan struct{}, 10),
	}
}

func (l *Listener) Listen(addr string) error {
	l.mu.Lock()
	if l.closed {
		l.mu.Unlock()
		return errors.New("listener closed")
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		l.mu.Unlock()
		return err
	}
	l.listener = ln
	l.mu.Unlock()

	go l.acceptLoop()
	return nil
}

func (l *Listener) acceptLoop() {
	for {
		conn, err := l.listener.Accept()
		if err != nil {
			return
		}

		if l.pm.shuttingDown.Load() {
			conn.Close()
			return
		}

		select {
		case l.pendingSem <- struct{}{}:
			l.pm.wg.Add(1)
			go l.handleInbound(conn)
		default:
			conn.Close()
			l.pm.enqueueTelemetry(func() {
				l.pm.telemetry.RecordResourceLimitReached("MaxPendingHandshakes")
			})
		}
	}
}

func (l *Listener) handleInbound(conn net.Conn) {
	defer func() {
		<-l.pendingSem
		l.pm.wg.Done()
	}()

	conn.SetDeadline(time.Now().Add(l.handshakeTimeout))

	var lenBuf [4]byte
	if _, err := io.ReadFull(conn, lenBuf[:]); err != nil {
		conn.Close()
		return
	}
	
	// Validate length BEFORE allocation
	length := binary.BigEndian.Uint32(lenBuf[:])
	if length == 0 || length > MaxFrameSize {
		conn.Close()
		return
	}

	data := make([]byte, length)
	if _, err := io.ReadFull(conn, data); err != nil {
		conn.Close()
		return
	}

	var pkt protocol.MeshPacket
	if err := proto.Unmarshal(data, &pkt); err != nil {
		conn.Close()
		return
	}

	expectedID := pkt.SourceNode

	if !l.pm.isExpectedInbound(expectedID) {
		conn.Close()
		return
	}

	fullPacket := make([]byte, 4+length)
	copy(fullPacket[0:4], lenBuf[:])
	copy(fullPacket[4:], data)

	rConn := &replayConn{
		Conn: conn,
		r:    io.MultiReader(bytes.NewReader(fullPacket), conn),
	}

	tconn := transport.NewConnection(rConn)
	ch := transport.NewDirectChannel(tconn, l.pm.localIdent)

	err := ch.AcceptHandshake(expectedID)
	rConn.SetDeadline(time.Time{}) // Clear deadline

	if err != nil {
		ch.Close()
		l.pm.enqueueTelemetry(func() {
			l.pm.telemetry.RecordHandshakeFailed(expectedID, err)
		})
		return
	}

	err = l.pm.register(ch, expectedID, false)
	if err != nil {
		// ch.Close() handled by register
	}
}

func (l *Listener) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.closed {
		return nil
	}
	l.closed = true
	if l.listener != nil {
		return l.listener.Close()
	}
	return nil
}
