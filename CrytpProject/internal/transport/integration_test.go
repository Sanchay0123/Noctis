package transport

import (
	"bytes"
	"net"
	"testing"
	"time"
)

func TestEndToEndDirectChannel(t *testing.T) {
	idA, idB := genIdentities(t)

	// Create interceptable pipe
	r1, w1 := net.Pipe() // A writes to w1, B reads from r1
	r2, w2 := net.Pipe() // B writes to w2, A reads from r2

	// Wrap connections
	connA := NewConnection(&interceptorConn{w1, r2, nil})

	var interceptedB []byte
	connB := NewConnection(&interceptorConn{w2, r1, &interceptedB})

	chA := NewDirectChannel(connA, idA)
	chB := NewDirectChannel(connB, idB)

	errChan := make(chan error, 2)
	go func() {
		errChan <- chA.InitiateHandshake(idB.PublicKey())
	}()
	go func() {
		errChan <- chB.AcceptHandshake(idA.PublicKey())
	}()

	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil {
			t.Fatalf("Handshake failed: %v", err)
		}
	}

	// Send Alice -> Bob
	msgA := []byte("hello from A")
	go func() {
		chA.SendPlaintext(msgA)
	}()

	out, err := chB.ReceivePlaintext()
	if err != nil {
		t.Fatalf("Bob receive failed: %v", err)
	}
	if !bytes.Equal(out, msgA) {
		t.Fatalf("Bob got wrong plaintext")
	}

	// Verify interception
	if bytes.Contains(interceptedB, msgA) {
		t.Fatalf("PLAINTEXT LEAKED TO TCP TRANSPORT B!")
	}

	// Send Bob -> Alice
	msgB := []byte("hello from B")
	go func() {
		chB.SendPlaintext(msgB)
	}()

	out2, err := chA.ReceivePlaintext()
	if err != nil {
		t.Fatalf("Alice receive failed: %v", err)
	}
	if !bytes.Equal(out2, msgB) {
		t.Fatalf("Alice got wrong plaintext")
	}

	chA.Close()
	chB.Close()
}

type interceptorConn struct {
	w           net.Conn
	r           net.Conn
	intercepted *[]byte
}

func (i *interceptorConn) Read(b []byte) (n int, err error) {
	n, err = i.r.Read(b)
	if i.intercepted != nil && n > 0 {
		*i.intercepted = append(*i.intercepted, b[:n]...)
	}
	return n, err
}

func (i *interceptorConn) Write(b []byte) (n int, err error) {
	return i.w.Write(b)
}

func (i *interceptorConn) Close() error {
	i.w.Close()
	return i.r.Close()
}

func (i *interceptorConn) LocalAddr() net.Addr {
	return i.r.LocalAddr()
}

func (i *interceptorConn) RemoteAddr() net.Addr {
	return i.w.RemoteAddr()
}

func (i *interceptorConn) SetDeadline(t time.Time) error {
	return nil
}

func (i *interceptorConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (i *interceptorConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func TestE2ESecurityBoundaryTamper(t *testing.T) {
	idA, idB := genIdentities(t)

	// Custom interceptor that mutates the ciphertext
	r1, w1 := net.Pipe()
	r2, w2 := net.Pipe()

	connA := NewConnection(&interceptorConn{w1, r2, nil})

	tamperConn := &tamperConn{w2, r1, false}
	connB := NewConnection(tamperConn)

	chA := NewDirectChannel(connA, idA)
	chB := NewDirectChannel(connB, idB)

	errChan := make(chan error, 2)
	go func() { errChan <- chA.InitiateHandshake(idB.PublicKey()) }()
	go func() { errChan <- chB.AcceptHandshake(idA.PublicKey()) }()

	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil {
			t.Fatalf("Handshake failed: %v", err)
		}
	}

	tamperConn.tamper = true

	go chA.SendPlaintext([]byte("hello from A"))

	_, err := chB.ReceivePlaintext()
	if err != ErrDecryptionFailed && err != ErrInvalidPacket {
		t.Fatalf("Expected decryption failure or invalid packet due to tampering, got: %v", err)
	}
}

type tamperConn struct {
	w      net.Conn
	r      net.Conn
	tamper bool
}

func (t *tamperConn) Read(b []byte) (n int, err error) {
	n, err = t.r.Read(b)
	if t.tamper && n > 0 {
		b[n-1] ^= 0xFF   // Flip last byte (will corrupt AEAD tag)
		t.tamper = false // Only tamper once
	}
	return n, err
}

func (t *tamperConn) Write(b []byte) (n int, err error) {
	return t.w.Write(b)
}

func (t *tamperConn) Close() error {
	t.w.Close()
	return t.r.Close()
}

func (t *tamperConn) LocalAddr() net.Addr                 { return t.r.LocalAddr() }
func (t *tamperConn) RemoteAddr() net.Addr                { return t.w.RemoteAddr() }
func (t *tamperConn) SetDeadline(tm time.Time) error      { return nil }
func (t *tamperConn) SetReadDeadline(tm time.Time) error  { return nil }
func (t *tamperConn) SetWriteDeadline(tm time.Time) error { return nil }
