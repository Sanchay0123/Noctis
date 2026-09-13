package mesh

import (
	"sync"
	"testing"
	"time"
)

func TestPeerQueueAccounting_SuccessfulEnqueue(t *testing.T) {
	p := &peer{
		state: PeerStateEstablished,
		queue: make(chan []byte, 1000),
	}

	msg := make([]byte, 100)
	ok := p.EnqueueForward(msg)
	if !ok {
		t.Fatal("expected enqueue to succeed")
	}

	p.queueMu.Lock()
	defer p.queueMu.Unlock()
	if p.queuedCount != 1 {
		t.Errorf("expected count 1, got %d", p.queuedCount)
	}
	if p.queuedBytes != 100 {
		t.Errorf("expected bytes 100, got %d", p.queuedBytes)
	}
}

func TestPeerQueueAccounting_FailedCountLimit(t *testing.T) {
	p := &peer{
		state:       PeerStateEstablished,
		queue:       make(chan []byte, 1000),
		queuedCount: 1000,
		queuedBytes: 0,
	}

	msg := make([]byte, 10)
	ok := p.EnqueueForward(msg)
	if ok {
		t.Fatal("expected enqueue to fail")
	}

	p.queueMu.Lock()
	defer p.queueMu.Unlock()
	if p.queuedCount != 1000 {
		t.Errorf("expected count 1000, got %d", p.queuedCount)
	}
	if p.queuedBytes != 0 {
		t.Errorf("expected bytes 0, got %d", p.queuedBytes)
	}
}

func TestPeerQueueAccounting_FailedByteLimit(t *testing.T) {
	p := &peer{
		state:       PeerStateEstablished,
		queue:       make(chan []byte, 1000),
		queuedCount: 1,
		queuedBytes: 2 * 1024 * 1024,
	}

	msg := make([]byte, 10)
	ok := p.EnqueueForward(msg)
	if ok {
		t.Fatal("expected enqueue to fail")
	}

	p.queueMu.Lock()
	defer p.queueMu.Unlock()
	if p.queuedCount != 1 {
		t.Errorf("expected count 1, got %d", p.queuedCount)
	}
	if p.queuedBytes != 2*1024*1024 {
		t.Errorf("expected bytes 2MiB, got %d", p.queuedBytes)
	}
}

func TestPeerQueueAccounting_ChannelFullRollback(t *testing.T) {
	p := &peer{
		state:       PeerStateEstablished,
		queue:       make(chan []byte, 1),
		queuedCount: 0,
		queuedBytes: 0,
	}

	// Fill channel manually
	p.queue <- make([]byte, 10)
	p.queuedCount = 1
	p.queuedBytes = 10

	msg := make([]byte, 20)
	ok := p.EnqueueForward(msg)
	if ok {
		t.Fatal("expected enqueue to fail")
	}

	p.queueMu.Lock()
	defer p.queueMu.Unlock()
	if p.queuedCount != 1 {
		t.Errorf("expected count to rollback to 1, got %d", p.queuedCount)
	}
	if p.queuedBytes != 10 {
		t.Errorf("expected bytes to rollback to 10, got %d", p.queuedBytes)
	}
}

func TestPeerQueueAccounting_Dequeue(t *testing.T) {
	p := &peer{
		state:  PeerStateEstablished,
		queue:  make(chan []byte, 1000),
		closed: make(chan struct{}),
	}

	msg := make([]byte, 100)
	p.EnqueueForward(msg)

	// Simulate dequeue
	<-p.queue
	p.queueMu.Lock()
	p.queuedCount--
	p.queuedBytes -= 100
	p.queueMu.Unlock()

	p.queueMu.Lock()
	defer p.queueMu.Unlock()
	if p.queuedCount != 0 {
		t.Errorf("expected count 0, got %d", p.queuedCount)
	}
	if p.queuedBytes != 0 {
		t.Errorf("expected bytes 0, got %d", p.queuedBytes)
	}
}

func TestPeerQueueAccounting_ReuseAfterDequeue(t *testing.T) {
	p := &peer{
		state:  PeerStateEstablished,
		queue:  make(chan []byte, 1000),
		closed: make(chan struct{}),
	}

	msg := make([]byte, 100)
	for i := 0; i < 1000; i++ {
		p.EnqueueForward(msg)
	}

	if p.EnqueueForward(msg) {
		t.Fatal("expected 1001st to fail")
	}

	// Dequeue
	<-p.queue
	p.queueMu.Lock()
	p.queuedCount--
	p.queuedBytes -= 100
	p.queueMu.Unlock()

	if !p.EnqueueForward(msg) {
		t.Fatal("expected enqueue to succeed after manual dequeue")
	}
}

func TestPeerQueueAccounting_ConcurrentEnqueueDequeue(t *testing.T) {
	p := &peer{
		state:  PeerStateEstablished,
		queue:  make(chan []byte, 1000),
		closed: make(chan struct{}),
	}

	// Consumer loop
	go func() {
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
			}
		}
	}()

	var wg sync.WaitGroup
	msg := make([]byte, 50)
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				p.EnqueueForward(msg)
				time.Sleep(1 * time.Millisecond)
			}
		}()
	}

	wg.Wait()
	time.Sleep(100 * time.Millisecond)
	close(p.closed)

	p.queueMu.Lock()
	defer p.queueMu.Unlock()
	if p.queuedCount != 0 {
		t.Errorf("concurrent count leak: %d", p.queuedCount)
	}
	if p.queuedBytes != 0 {
		t.Errorf("concurrent byte leak: %d", p.queuedBytes)
	}
}

func TestPeerQueueAccounting_Shutdown(t *testing.T) {
	p := &peer{
		state:   PeerStateEstablished,
		queue:   make(chan []byte, 100),
		closed:  make(chan struct{}),
		channel: nil,
		manager: NewPeerManager(nil, nil, nil),
	}

	p.EnqueueForward(make([]byte, 10))
	p.EnqueueForward(make([]byte, 10))

	p.Close() // Shuts down gracefully

	p.queueMu.Lock()
	defer p.queueMu.Unlock()
	if p.queuedCount != 2 {
		t.Errorf("accounting corrupted during shutdown: %d", p.queuedCount)
	}
}
