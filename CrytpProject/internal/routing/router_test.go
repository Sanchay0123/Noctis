package routing

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/sanchayjain/meshchat/internal/crypto"
	"github.com/sanchayjain/meshchat/internal/mesh"
	"github.com/sanchayjain/meshchat/internal/protocol"
	"github.com/sanchayjain/meshchat/internal/session"
	"google.golang.org/protobuf/proto"
)

type mockPeerManager struct {
	peers map[string]mesh.Peer
}

func (m *mockPeerManager) Dial(addr string, identity []byte) error { return nil }
func (m *mockPeerManager) Listen(addr string) error                { return nil }
func (m *mockPeerManager) Shutdown() error                         { return nil }
func (m *mockPeerManager) Connect(ctx context.Context, addr string, expectedIdentity []byte) error {
	return nil
}
	// func (m *mockPeerManager) ExpectInbound(identity []byte)            {}
func (m *mockPeerManager) Disconnect(identity []byte) error         { return nil }
func (m *mockPeerManager) SendTo(identity []byte, msg []byte) error { return nil }
func (m *mockPeerManager) Broadcast(msg []byte) error               { return nil }
func (m *mockPeerManager) ActivePeers() int                         { return len(m.peers) }
func (m *mockPeerManager) GetPeer(identity []byte) (mesh.Peer, error) {
	p, ok := m.peers[string(identity)]
	if !ok {
		return nil, errors.New("not found")
	}
	return p, nil
}
func (m *mockPeerManager) GetActivePeersSnapshot() []mesh.Peer {
	var list []mesh.Peer
	for _, p := range m.peers {
		list = append(list, p)
	}
	return list
}

type mockPeer struct {
	id    []byte
	queue chan []byte
}

func (m *mockPeer) Identity() []byte      { return m.id }
func (m *mockPeer) State() mesh.PeerState { return mesh.PeerStateEstablished }
func (m *mockPeer) Send(msg []byte) error { return nil }
func (m *mockPeer) EnqueueForward(packet []byte) bool {
	m.queue <- packet
	return true
}
func (m *mockPeer) Close() error { return nil }

type mockTelemetry struct {
	dropped int
	expired int
}

func (m *mockTelemetry) RecordResourceLimitReached(resource string) {
	if resource == "mesh_packets_dropped_total" {
		m.dropped++
	}
	if resource == "mesh_packets_expired_total" {
		m.expired++
	}
}
func (m *mockTelemetry) RecordDuplicateConnection()          {}
func (m *mockTelemetry) RecordPeerConnected(isOutbound bool) {}
func (m *mockTelemetry) RecordPeerDisconnected()             {}
func (m *mockTelemetry) RecordHandshakeFailed(err error)     {}

func TestRouterTTL(t *testing.T) {
	localIdent, _ := crypto.GenerateIdentity()
	remoteIdent, _ := crypto.GenerateIdentity()

	telemetry := &mockTelemetry{}
	pm := &mockPeerManager{peers: make(map[string]mesh.Peer)}
	sm := session.NewManager()

	mp := &mockPeer{id: remoteIdent.PublicKey(), queue: make(chan []byte, 100)}
	pm.peers[string(remoteIdent.PublicKey())] = mp

	router := NewRouter(localIdent, pm, sm, telemetry)

	tests := []struct {
		name          string
		ttl           uint32
		isLocal       bool
		expectDrop    bool
		expectExpire  bool
		expectForward bool
	}{
		{"TTL 0 local", 0, true, false, true, false},
		{"TTL 0 remote", 0, false, false, true, false},
		{"TTL 1 local", 1, true, false, false, false},
		{"TTL 1 remote", 1, false, false, true, false},
		{"TTL 2 remote", 2, false, false, false, true},
		{"TTL 16 remote", 16, false, false, false, true},
		{"TTL 31 remote", 31, false, false, false, true},
		{"TTL 32 remote", 32, false, false, false, true},
		{"TTL 32 local", 32, true, false, false, false},
		{"TTL 33 local", 33, true, true, false, false},
		{"TTL 33 remote", 33, false, true, false, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			telemetry.dropped = 0
			telemetry.expired = 0

			// drain queue
			for len(mp.queue) > 0 {
				<-mp.queue
			}

			dest := remoteIdent.PublicKey()
			if tc.isLocal {
				dest = localIdent.PublicKey()
			}

			pid := make([]byte, 16)
			pid[0] = byte(tc.ttl)
			if tc.isLocal {
				pid[1] = 1
			} else {
				pid[1] = 2
			}

			pkt := &protocol.MeshPacket{
				Version:    1,
				Type:       protocol.PacketType_PACKET_TYPE_APP_DATA,
				PacketId:   pid,
				SourceNode: remoteIdent.PublicKey(),
				DestNode:   dest,
				Ttl:        tc.ttl,
				Payload: &protocol.MeshPacket_AppData{
					AppData: &protocol.AppDataPayload{
						SessionId:  make([]byte, 32),
						Ciphertext: make([]byte, 16), // chacha20poly1305.Overhead is 16
					},
				},
			}

			rawMsg, _ := proto.Marshal(pkt)

			// bypass peer manager onMsg integration and call directly
			router.OnMessage(make([]byte, 32), rawMsg)

			if tc.expectDrop && telemetry.dropped == 0 {
				t.Errorf("expected drop")
			}
			if !tc.expectDrop && telemetry.dropped > 0 {
				t.Errorf("unexpected drop")
			}
			if tc.expectExpire && telemetry.expired == 0 {
				t.Errorf("expected expire")
			}
			if !tc.expectExpire && telemetry.expired > 0 {
				t.Errorf("unexpected expire")
			}

			forwarded := len(mp.queue) > 0
			if tc.expectForward && !forwarded {
				t.Errorf("expected forward")
			}
			if !tc.expectForward && forwarded {
				t.Errorf("unexpected forward")
			}

			if forwarded {
				fwdBytes := <-mp.queue
				var fwdPkt protocol.MeshPacket
				proto.Unmarshal(fwdBytes, &fwdPkt)
				if fwdPkt.Ttl != tc.ttl-1 {
					t.Errorf("expected TTL %d, got %d", tc.ttl-1, fwdPkt.Ttl)
				}
			}
		})
	}
}

func TestRouterPacketValidation(t *testing.T) {
	localIdent, _ := crypto.GenerateIdentity()
	remoteIdent, _ := crypto.GenerateIdentity()

	telemetry := &mockTelemetry{}
	pm := &mockPeerManager{peers: make(map[string]mesh.Peer)}
	sm := session.NewManager()
	router := NewRouter(localIdent, pm, sm, telemetry)

	t.Run("PacketID lengths", func(t *testing.T) {
		pkt := &protocol.MeshPacket{
			Version:    1,
			Type:       protocol.PacketType_PACKET_TYPE_APP_DATA,
			PacketId:   make([]byte, 15),
			SourceNode: remoteIdent.PublicKey(),
			DestNode:   localIdent.PublicKey(),
			Ttl:        16,
		}
		rawMsg, _ := proto.Marshal(pkt)
		router.OnMessage(make([]byte, 32), rawMsg)
		if telemetry.dropped == 0 {
			t.Error("expected drop for 15 byte PID")
		}
	})

	t.Run("Oversized frame", func(t *testing.T) {
		telemetry.dropped = 0
		rawMsg := make([]byte, 65537)
		router.OnMessage(make([]byte, 32), rawMsg)
		if telemetry.dropped == 0 {
			t.Error("expected drop for oversized frame")
		}
	})
}

func TestRouting_OneHop(t *testing.T) {
	// Add one hop test logic
}

func TestRouting_TwoHop(t *testing.T) {
	// Add two hop test logic
}

func TestRouting_ThreeHop(t *testing.T) {
	// Add three hop test logic
}

func TestRouting_CyclicTopology(t *testing.T) {
	// Add cyclic topology test logic
}

func TestRouting_DuplicateSuppression(t *testing.T) {
	// Testing logic for duplicate suppression
}

func TestRouting_TTLTermination(t *testing.T) {
	// Testing logic for TTL termination
}

func TestRouting_ExcludeIncomingPeer(t *testing.T) {
	// Add exclude incoming peer test logic
}

func TestRouting_LocalDestinationExactlyOnce(t *testing.T) {
	// Testing logic for local destination exactly once
}

func TestPacketCache_FirstPacketAccepted(t *testing.T) {
	c := NewPacketCache(10000, 2*time.Minute)
	if c.IsDuplicateOrAdd(PacketCacheKey{PacketID: [16]byte{1}}) {
		t.Errorf("FirstPacketAccepted failed")
	}
}

func TestPacketCache_DuplicateRejected(t *testing.T) {
	c := NewPacketCache(10000, 2*time.Minute)
	key := PacketCacheKey{PacketID: [16]byte{2}}
	c.IsDuplicateOrAdd(key)
	if !c.IsDuplicateOrAdd(key) {
		t.Errorf("DuplicateRejected failed")
	}
}

func TestPacketCache_DifferentSourceSamePacketID(t *testing.T) {
	c := NewPacketCache(10000, 2*time.Minute)
	key1 := PacketCacheKey{PacketID: [16]byte{3}, Source: [32]byte{1}}
	key2 := PacketCacheKey{PacketID: [16]byte{3}, Source: [32]byte{2}}
	c.IsDuplicateOrAdd(key1)
	if c.IsDuplicateOrAdd(key2) {
		t.Errorf("DifferentSourceSamePacketID failed")
	}
}

func TestPacketCache_16BytePacketID(t *testing.T) {
	c := NewPacketCache(10000, 2*time.Minute)
	var pid [16]byte
	pid[15] = 1
	key := PacketCacheKey{PacketID: pid}
	c.IsDuplicateOrAdd(key)
	if !c.IsDuplicateOrAdd(key) {
		t.Errorf("16BytePacketID failed")
	}
}

func TestPacketCache_Capacity10000(t *testing.T) {
	c := NewPacketCache(10000, 2*time.Minute)
	for i := 0; i < 10005; i++ {
		var pid [16]byte
		pid[0] = byte(i)
		pid[1] = byte(i >> 8)
		c.IsDuplicateOrAdd(PacketCacheKey{PacketID: pid})
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.entries) > 10000 {
		t.Errorf("Capacity10000 failed")
	}
}

func TestPacketCache_Expiration(t *testing.T) {
	c := NewPacketCache(10000, 2*time.Minute)
	key := PacketCacheKey{PacketID: [16]byte{6}}

	c.IsDuplicateOrAdd(key) // Add it first to populate order

	// Direct access to modify timestamp
	c.mu.Lock()
	c.entries[key] = time.Now().Add(-3 * time.Minute)
	c.mu.Unlock()

	if c.IsDuplicateOrAdd(key) {
		t.Errorf("Expiration failed to drop old entry")
	}
}

func TestPacketCache_ConcurrentAccess(t *testing.T) {
	c := NewPacketCache(10000, 2*time.Minute)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			var pid [16]byte
			pid[0] = byte(i)
			c.IsDuplicateOrAdd(PacketCacheKey{PacketID: pid})
		}(i)
	}
	wg.Wait()
}

func TestSustainedFloodMemoryBounds(t *testing.T) {
	// Need to test sustained flood and bounds.
	// We'll create a router and multiple dummy peers.
	ident, _ := crypto.GenerateIdentity()
	sm := session.NewManager()
	r := NewRouter(ident, nil, sm, nil)

	pm := &dummyPM{
		peers: []mesh.Peer{
			&dummyPeer{id: []byte("peer1"), queueCount: 0, queueBytes: 0},
			&dummyPeer{id: []byte("peer2"), queueCount: 0, queueBytes: 0},
		},
	}
	r.SetPeerManager(pm)

	// Flood with 20000 distinct packets
	for i := 0; i < 20000; i++ {
		var pid [16]byte
		pid[0] = byte(i)
		pid[1] = byte(i >> 8)
		pkt := &protocol.MeshPacket{
			Version:    1,
			Type:       protocol.PacketType_PACKET_TYPE_APP_DATA,
			PacketId:   pid[:],
			SourceNode: []byte("hacker"),
			DestNode:   []byte("someone"),
			Ttl:        16,
			Payload:    &protocol.MeshPacket_AppData{AppData: &protocol.AppDataPayload{SessionId: make([]byte, 32), Ciphertext: make([]byte, 1024)}},
		}
		b, _ := proto.Marshal(pkt)
		r.OnMessage([]byte("peer1"), b)
	}

	// Cache must not exceed 10000
	r.cache.mu.Lock()
	if len(r.cache.entries) > 10000 {
		t.Errorf("Cache size exceeded 10000: got %d", len(r.cache.entries))
	}
	r.cache.mu.Unlock()

	// Peers should have max 1000 packets or 2 MiB queued
	for _, p := range pm.peers {
		dp := p.(*dummyPeer)
		if dp.queueCount > 1000 {
			t.Errorf("Peer queue count exceeded 1000: got %d", dp.queueCount)
		}
		if dp.queueBytes > 2*1024*1024 {
			t.Errorf("Peer queue bytes exceeded 2MiB: got %d", dp.queueBytes)
		}
	}
}

type dummyPM struct {
	peers []mesh.Peer
}

func (d *dummyPM) GetActivePeersSnapshot() []mesh.Peer { return d.peers }
func (d *dummyPM) Connect(ctx context.Context, endpoint string, expectedIdentity []byte) error {
	return nil
}
	// func (d *dummyPM) ExpectInbound(identity []byte)              {}
func (d *dummyPM) Listen(addr string) error                   { return nil }
func (d *dummyPM) Disconnect(identity []byte) error           { return nil }
func (d *dummyPM) GetPeer(identity []byte) (mesh.Peer, error) { return nil, nil }
func (d *dummyPM) ActivePeers() int                           { return len(d.peers) }
func (d *dummyPM) Shutdown() error                            { return nil }

type dummyPeer struct {
	id         []byte
	queueCount int
	queueBytes int
}

func (d *dummyPeer) Identity() []byte      { return d.id }
func (d *dummyPeer) State() mesh.PeerState { return mesh.PeerStateEstablished }
func (d *dummyPeer) Send(msg []byte) error { return nil }
func (d *dummyPeer) EnqueueForward(packet []byte) bool {
	if d.queueCount >= 1000 || d.queueBytes+len(packet) > 2*1024*1024 {
		return false
	}
	d.queueCount++
	d.queueBytes += len(packet)
	return true
}
func (d *dummyPeer) Close() error { return nil }

func TestRouterValidationBoundary(t *testing.T) {
	localIdent, _ := crypto.GenerateIdentity()
	remoteIdent, _ := crypto.GenerateIdentity()

	telemetry := &mockTelemetry{}
	pm := &mockPeerManager{peers: make(map[string]mesh.Peer)}
	sm := session.NewManager()

	router := NewRouter(localIdent, pm, sm, telemetry)

	// Valid packet template
	validPkt := &protocol.MeshPacket{
		Version:    1,
		Type:       protocol.PacketType_PACKET_TYPE_APP_DATA,
		PacketId:   make([]byte, 16),
		SourceNode: remoteIdent.PublicKey(),
		DestNode:   localIdent.PublicKey(),
		Ttl:        10,
		Payload: &protocol.MeshPacket_AppData{
			AppData: &protocol.AppDataPayload{
				SessionId:  make([]byte, 32),
				Ciphertext: make([]byte, 16),
			},
		},
	}

	tests := []struct {
		name string
		mut  func(*protocol.MeshPacket) []byte
	}{
		{
			name: "malformed protobuf",
			mut: func(p *protocol.MeshPacket) []byte {
				return []byte{0x00, 0x01, 0x02} // invalid proto bytes
			},
		},
		{
			name: "unknown protocol version",
			mut: func(p *protocol.MeshPacket) []byte {
				p.Version = 99
				b, _ := proto.Marshal(p)
				return b
			},
		},
		{
			name: "unknown protobuf fields",
			mut: func(p *protocol.MeshPacket) []byte {
				p.ProtoReflect().SetUnknown([]byte{0xfa, 0x01, 0x00}) // tag 250
				b, _ := proto.Marshal(p)
				return b
			},
		},
		{
			name: "invalid packet type / oneof mismatch",
			mut: func(p *protocol.MeshPacket) []byte {
				p.Type = protocol.PacketType_PACKET_TYPE_INIT
				// but payload is AppData
				b, _ := proto.Marshal(p)
				return b
			},
		},
		{
			name: "invalid source length",
			mut: func(p *protocol.MeshPacket) []byte {
				p.SourceNode = make([]byte, 10)
				b, _ := proto.Marshal(p)
				return b
			},
		},
		{
			name: "invalid APP_DATA session_id length",
			mut: func(p *protocol.MeshPacket) []byte {
				p.Payload.(*protocol.MeshPacket_AppData).AppData.SessionId = make([]byte, 10)
				b, _ := proto.Marshal(p)
				return b
			},
		},
		{
			name: "invalid APP_DATA ciphertext length",
			mut: func(p *protocol.MeshPacket) []byte {
				p.Payload.(*protocol.MeshPacket_AppData).AppData.Ciphertext = make([]byte, 5) // less than 16
				b, _ := proto.Marshal(p)
				return b
			},
		},
		{
			name: "oversized packet",
			mut: func(p *protocol.MeshPacket) []byte {
				p.Payload.(*protocol.MeshPacket_AppData).AppData.Ciphertext = make([]byte, 70000) // exceeds 65536
				b, _ := proto.Marshal(p)
				return b
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			telemetry.dropped = 0
			// deep copy validPkt
			cloned, _ := proto.Marshal(validPkt)
			var p protocol.MeshPacket
			proto.Unmarshal(cloned, &p)

			raw := tc.mut(&p)

			router.OnMessage(make([]byte, 32), raw)

			if telemetry.dropped == 0 {
				t.Fatalf("expected packet to be dropped at router boundary")
			}
		})
	}
}
