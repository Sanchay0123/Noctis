package mesh

import "errors"

var (
	ErrPeerNotFound       = errors.New("peer not found")
	ErrQueueFull          = errors.New("forwarding queue is full")
	ErrDuplicateConnection = errors.New("duplicate connection resolved locally")
	ErrPeerLimitReached   = errors.New("peer resource limit reached")
	ErrShutdown           = errors.New("peer manager is shutting down")
	ErrHandshakeTimeout   = errors.New("handshake timed out")
	ErrInvalidIdentity    = errors.New("invalid or unauthenticated identity")
	ErrNotEstablished     = errors.New("peer is not in ESTABLISHED state")
)
