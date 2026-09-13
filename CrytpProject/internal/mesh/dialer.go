package mesh

import (
	"context"
	"net"
	"time"

	"github.com/sanchayjain/meshchat/internal/transport"
)

type Dialer struct {
	pm                 *peerManager
	maxConcurrentDials int
	dialTimeout        time.Duration
	handshakeTimeout   time.Duration

	dialSem chan struct{}
}

func NewDialer(pm *peerManager) *Dialer {
	return &Dialer{
		pm:                 pm,
		maxConcurrentDials: 10,
		dialTimeout:        5 * time.Second,
		handshakeTimeout:   5 * time.Second,
		dialSem:            make(chan struct{}, 10),
	}
}

func (d *Dialer) Dial(ctx context.Context, endpoint string, expectedIdentity []byte) error {
	select {
	case d.dialSem <- struct{}{}:
	case <-ctx.Done():
		return ctx.Err()
	}

	defer func() { <-d.dialSem }()

	dialCtx, cancel := context.WithTimeout(ctx, d.dialTimeout)
	defer cancel()

	conn, err := net.DialTimeout("tcp", endpoint, d.dialTimeout)
	if err != nil {
		// Respect context cancellation
		select {
		case <-dialCtx.Done():
			return dialCtx.Err()
		default:
			return err
		}
	}

	if d.pm.shuttingDown.Load() {
		conn.Close()
		return ErrShutdown
	}

	conn.SetDeadline(time.Now().Add(d.handshakeTimeout))

	tconn := transport.NewConnection(conn)
	ch := transport.NewDirectChannel(tconn, d.pm.localIdent)

	err = ch.InitiateHandshake(expectedIdentity)
	conn.SetDeadline(time.Time{})

	if err != nil {
		ch.Close()
		d.pm.enqueueTelemetry(func() {
			d.pm.telemetry.RecordHandshakeFailed(err)
		})
		return err
	}

	return d.pm.register(ch, expectedIdentity, true)
}
