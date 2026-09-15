package tests

import (
	"context"
	"crypto/rand"
	"encoding/csv"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sanchayjain/meshchat/internal/app"
	"github.com/sanchayjain/meshchat/internal/crypto"
)

func getFreePort() int {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

func waitForNode(port int) {
	for i := 0; i < 50; i++ {
		conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
		if err == nil {
			conn.Close()
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func genId(b *testing.B) *crypto.NodeIdentity {
	id, _ := crypto.GenerateIdentity()
	return id
}

// Exp A: Cryptography
func BenchmarkM10_Crypto_X25519(b *testing.B) {
	k1, _ := crypto.GenerateEphemeralKey()
	k2, _ := crypto.GenerateEphemeralKey()
	pub2, _ := k2.PublicKey()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k1.ComputeSharedSecret(pub2)
	}
}

func BenchmarkM10_Crypto_Ed25519_Sign(b *testing.B) {
	id, _ := crypto.GenerateIdentity()
	msg := make([]byte, 256)
	rand.Read(msg)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id.Sign(msg)
	}
}

func BenchmarkM10_Crypto_Ed25519_Verify(b *testing.B) {
	id, _ := crypto.GenerateIdentity()
	msg := make([]byte, 256)
	rand.Read(msg)
	sig, _ := id.Sign(msg)
	pub := id.PublicKey()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		crypto.VerifySignature(pub, msg, sig)
	}
}

func benchmarkAEADEncrypt(b *testing.B, size int) {
	idA := genId(b)
	idB := genId(b)
	initA, _ := crypto.NewInitiatorHandshake(idA, idB.PublicKey())
	respB, _ := crypto.NewResponderHandshake(idB, idA.PublicKey())
	initMsg, _ := initA.GenerateInit()
	respB.ProcessInit(initMsg)
	respMsg, _ := respB.GenerateResp()
	initA.ProcessResp(respMsg)
	sess, _ := initA.Session()

	plaintext := make([]byte, size)
	rand.Read(plaintext)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sess.EncryptMessage(plaintext)
	}
}

func benchmarkAEADDecrypt(b *testing.B, size int) {
	idA := genId(b)
	idB := genId(b)
	initA, _ := crypto.NewInitiatorHandshake(idA, idB.PublicKey())
	respB, _ := crypto.NewResponderHandshake(idB, idA.PublicKey())
	initMsg, _ := initA.GenerateInit()
	respB.ProcessInit(initMsg)
	respMsg, _ := respB.GenerateResp()
	initA.ProcessResp(respMsg)
	sessA, _ := initA.Session()
	sessB, _ := respB.Session()

	plaintext := make([]byte, size)
	rand.Read(plaintext)
	seq, ciphertext, _ := sessA.EncryptMessage(plaintext)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sessB.DecryptMessage(seq, ciphertext)
	}
}

func BenchmarkM10_Crypto_ChaCha20Poly1305_Encrypt_64(b *testing.B)   { benchmarkAEADEncrypt(b, 64) }
func BenchmarkM10_Crypto_ChaCha20Poly1305_Encrypt_256(b *testing.B)  { benchmarkAEADEncrypt(b, 256) }
func BenchmarkM10_Crypto_ChaCha20Poly1305_Encrypt_1K(b *testing.B)   { benchmarkAEADEncrypt(b, 1024) }
func BenchmarkM10_Crypto_ChaCha20Poly1305_Encrypt_4K(b *testing.B)   { benchmarkAEADEncrypt(b, 4096) }
func BenchmarkM10_Crypto_ChaCha20Poly1305_Decrypt_64(b *testing.B)   { benchmarkAEADDecrypt(b, 64) }
func BenchmarkM10_Crypto_ChaCha20Poly1305_Decrypt_256(b *testing.B)  { benchmarkAEADDecrypt(b, 256) }
func BenchmarkM10_Crypto_ChaCha20Poly1305_Decrypt_1K(b *testing.B)   { benchmarkAEADDecrypt(b, 1024) }
func BenchmarkM10_Crypto_ChaCha20Poly1305_Decrypt_4K(b *testing.B)   { benchmarkAEADDecrypt(b, 4096) }

// Exp B & E: Direct Latency & Message Size Scaling
func TestM10_DirectLatency(t *testing.T) {
	portA := getFreePort()
	portB := getFreePort()
	os.Setenv("MESH_LISTEN_ADDR", fmt.Sprintf("127.0.0.1:%d", portA))
	nodeA, _ := app.NewNode("alice")
	go nodeA.Start()
	os.Setenv("MESH_LISTEN_ADDR", fmt.Sprintf("127.0.0.1:%d", portB))
	nodeB, _ := app.NewNode("bob")
	go nodeB.Start()
	defer nodeA.Stop()
	defer nodeB.Stop()

	waitForNode(portA)
	waitForNode(portB)

	nodeA.DialNode(context.Background(), fmt.Sprintf("127.0.0.1:%d", portB), nodeB.GetLocalIdentity())
	time.Sleep(200 * time.Millisecond)

	nodeA.StartConversation(nodeB.GetLocalIdentity())
	time.Sleep(200 * time.Millisecond)

	os.MkdirAll("../m10_2_evidence", 0755)
	f, _ := os.Create("../m10_2_evidence/m10_2_direct_latency.csv")
	defer f.Close()
	writer := csv.NewWriter(f)
	writer.Write([]string{"trial", "payload_bytes", "latency_ns"})

	fSize, _ := os.Create("../m10_2_evidence/m10_2_message_size.csv")
	defer fSize.Close()
	writerSize := csv.NewWriter(fSize)
	writerSize.Write([]string{"trial", "payload_bytes", "latency_ns", "delivered"})

	sizes := []int{64, 256, 1024, 4096, 16384, 32768, 49152, 61440}
	trials := 100

	eventsB := nodeB.SubscribeEvents()

	for _, size := range sizes {
		payload := make([]byte, size)
		rand.Read(payload)
		payloadStr := string(payload)

		// Warmup
		for i := 0; i < 5; i++ {
			nodeA.SendMessage(nodeB.GetLocalIdentity(), payloadStr)
			select {
			case <-eventsB:
			case <-time.After(50 * time.Millisecond):
			}
		}

		for i := 1; i <= trials; i++ {
			start := time.Now()
			err := nodeA.SendMessage(nodeB.GetLocalIdentity(), payloadStr)
			if err != nil {
				writerSize.Write([]string{strconv.Itoa(i), strconv.Itoa(size), "0", "false"})
				continue
			}

			select {
			case <-eventsB:
				latency := time.Since(start).Nanoseconds()
				writerSize.Write([]string{strconv.Itoa(i), strconv.Itoa(size), strconv.FormatInt(latency, 10), "true"})
				if size <= 4096 {
					writer.Write([]string{strconv.Itoa(i), strconv.Itoa(size), strconv.FormatInt(latency, 10)})
				}
			case <-time.After(100 * time.Millisecond):
				writerSize.Write([]string{strconv.Itoa(i), strconv.Itoa(size), "0", "false"})
			}
		}
	}
	writer.Flush()
	writerSize.Flush()
}

// Exp C: Multi-Hop Latency
func TestM10_MultiHopLatency(t *testing.T) {
	portA := getFreePort()
	portB := getFreePort()
	portC := getFreePort()

	os.Setenv("MESH_LISTEN_ADDR", fmt.Sprintf("127.0.0.1:%d", portA))
	nodeA, _ := app.NewNode("alice")
	go nodeA.Start()
	os.Setenv("MESH_LISTEN_ADDR", fmt.Sprintf("127.0.0.1:%d", portB))
	nodeB, _ := app.NewNode("bob")
	go nodeB.Start()
	os.Setenv("MESH_LISTEN_ADDR", fmt.Sprintf("127.0.0.1:%d", portC))
	nodeC, _ := app.NewNode("carol")
	go nodeC.Start()

	defer nodeA.Stop()
	defer nodeB.Stop()
	defer nodeC.Stop()

	waitForNode(portA)
	waitForNode(portB)
	waitForNode(portC)

	nodeA.DialNode(context.Background(), fmt.Sprintf("127.0.0.1:%d", portB), nodeB.GetLocalIdentity())
	nodeB.DialNode(context.Background(), fmt.Sprintf("127.0.0.1:%d", portC), nodeC.GetLocalIdentity())

	time.Sleep(500 * time.Millisecond)

	nodeA.StartConversation(nodeC.GetLocalIdentity())
	time.Sleep(500 * time.Millisecond)

	f, _ := os.Create("../m10_2_evidence/m10_2_multihop_latency.csv")
	defer f.Close()
	writer := csv.NewWriter(f)
	writer.Write([]string{"trial", "payload_bytes", "latency_ns"})

	sizes := []int{64, 256, 1024, 4096}
	trials := 100
	eventsC := nodeC.SubscribeEvents()

	for _, size := range sizes {
		payload := make([]byte, size)
		rand.Read(payload)
		payloadStr := string(payload)

		// Warmup
		for i := 0; i < 5; i++ {
			nodeA.SendMessage(nodeC.GetLocalIdentity(), payloadStr)
			select {
			case <-eventsC:
			case <-time.After(50 * time.Millisecond):
			}
		}

		for i := 1; i <= trials; i++ {
			start := time.Now()
			nodeA.SendMessage(nodeC.GetLocalIdentity(), payloadStr)
			select {
			case <-eventsC:
				latency := time.Since(start).Nanoseconds()
				writer.Write([]string{strconv.Itoa(i), strconv.Itoa(size), strconv.FormatInt(latency, 10)})
			case <-time.After(100 * time.Millisecond):
				writer.Write([]string{strconv.Itoa(i), strconv.Itoa(size), "0"})
			}
		}
	}
	writer.Flush()
}

// Exp D: Throughput
func TestM10_Throughput(t *testing.T) {
	portA := getFreePort()
	portB := getFreePort()
	os.Setenv("MESH_LISTEN_ADDR", fmt.Sprintf("127.0.0.1:%d", portA))
	nodeA, _ := app.NewNode("alice")
	go nodeA.Start()
	os.Setenv("MESH_LISTEN_ADDR", fmt.Sprintf("127.0.0.1:%d", portB))
	nodeB, _ := app.NewNode("bob")
	go nodeB.Start()
	defer nodeA.Stop()
	defer nodeB.Stop()
	waitForNode(portB)

	nodeA.DialNode(context.Background(), fmt.Sprintf("127.0.0.1:%d", portB), nodeB.GetLocalIdentity())
	time.Sleep(200 * time.Millisecond)

	nodeA.StartConversation(nodeB.GetLocalIdentity())
	time.Sleep(200 * time.Millisecond)

	eventsB := nodeB.SubscribeEvents()

	// Drain any session events
	drain := true
	for drain {
		select {
		case <-eventsB:
		default:
			drain = false
		}
	}

	payloadSize := 1024
	payload := make([]byte, payloadSize)
	rand.Read(payload)
	payloadStr := string(payload)
	totalMessages := 5000

	f, _ := os.Create("../m10_2_evidence/m10_2_throughput.csv")
	defer f.Close()
	writer := csv.NewWriter(f)
	writer.Write([]string{"messages_attempted", "messages_delivered", "elapsed_ns", "payload_bytes_delivered"})

	var delivered int32
	start := time.Now()

	go func() {
		for i := 0; i < totalMessages; i++ {
			nodeA.SendMessage(nodeB.GetLocalIdentity(), payloadStr)
			time.Sleep(100 * time.Microsecond)
		}
	}()

	timeout := time.After(10 * time.Second)
	done := false
	for !done {
		select {
		case <-eventsB:
			atomic.AddInt32(&delivered, 1)
			if atomic.LoadInt32(&delivered) == int32(totalMessages) {
				done = true
			}
		case <-timeout:
			done = true
		}
	}
	elapsed := time.Since(start).Nanoseconds()

	writer.Write([]string{
		strconv.Itoa(totalMessages),
		strconv.Itoa(int(delivered)),
		strconv.FormatInt(elapsed, 10),
		strconv.Itoa(int(delivered) * payloadSize),
	})
	writer.Flush()
}

// Exp H: Handshake Scaling
func TestM10_HandshakeScaling(t *testing.T) {
	portA := getFreePort()
	portB := getFreePort()
	os.Setenv("MESH_LISTEN_ADDR", fmt.Sprintf("127.0.0.1:%d", portA))
	nodeA, _ := app.NewNode("alice")
	go nodeA.Start()
	os.Setenv("MESH_LISTEN_ADDR", fmt.Sprintf("127.0.0.1:%d", portB))
	nodeB, _ := app.NewNode("bob")
	go nodeB.Start()
	defer nodeA.Stop()
	defer nodeB.Stop()
	waitForNode(portB)
	
	f, _ := os.Create("../m10_2_evidence/m10_2_handshake.csv")
	defer f.Close()
	writer := csv.NewWriter(f)
	writer.Write([]string{"trial", "latency_ns"})
	
	for i := 1; i <= 50; i++ {
		start := time.Now()
		nodeA.DialNode(context.Background(), fmt.Sprintf("127.0.0.1:%d", portB), nodeB.GetLocalIdentity())
		time.Sleep(50 * time.Millisecond) // Ensure tcp is ready
		nodeA.StartConversation(nodeB.GetLocalIdentity())
		latency := time.Since(start).Nanoseconds()
		writer.Write([]string{strconv.Itoa(i), strconv.FormatInt(latency, 10)})
		time.Sleep(10 * time.Millisecond)
	}
	writer.Flush()
}
