package crypto

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/hkdf"
	"io"
	"math"
	"sync"
)

var (
	ErrInvalidHandshakeState   = errors.New("invalid handshake state or transition")
	ErrInvalidHandshakeMessage = errors.New("invalid handshake message")
	ErrSessionNotEstablished   = errors.New("session not established")
	ErrSequenceExhausted       = errors.New("sequence number exhausted")
	ErrReplayDetected          = errors.New("replay detected")
	ErrDecryptFailed           = errors.New("message decryption failed")
)

const (
	protocolVersion = byte(0x01)
	domainString    = "MeshChat-Handshake-v1"
)

var domainBytes = []byte(domainString)
var zero32 = make([]byte, 32)
var appDataContext = []byte("MeshChat-AppData-v1")

type InitiatorState int

const (
	InitiatorStateNew InitiatorState = iota
	InitiatorStateInitSent
	InitiatorStateEstablished
	InitiatorStateFailed
)

type ResponderState int

const (
	ResponderStateNew ResponderState = iota
	ResponderStateInitReceived
	ResponderStateRespSent
	ResponderStateEstablished
	ResponderStateFailed
)

type Session struct {
	peerID      []byte
	id          []byte
	sendKey     []byte
	receiveKey  []byte
	isInitiator bool

	sendMu  sync.Mutex
	sendSeq uint64

	receiveMu       sync.Mutex
	highestReceived uint64
	replayBitmap    uint64
}

func (s *Session) PeerID() []byte {
	out := make([]byte, len(s.peerID))
	copy(out, s.peerID)
	return out
}

func (s *Session) ID() []byte {
	out := make([]byte, 32)
	copy(out, s.id)
	return out
}

func (s *Session) SendKey() []byte {
	out := make([]byte, 32)
	copy(out, s.sendKey)
	return out
}

func (s *Session) ReceiveKey() []byte {
	out := make([]byte, 32)
	copy(out, s.receiveKey)
	return out
}

func (s *Session) EncryptMessage(plaintext []byte) (uint64, []byte, error) {
	s.sendMu.Lock()
	defer s.sendMu.Unlock()

	if s.sendSeq == math.MaxUint64 {
		return 0, nil, ErrSequenceExhausted
	}
	seq := s.sendSeq
	s.sendSeq++

	nonce := make([]byte, 12)
	binary.BigEndian.PutUint64(nonce[4:], seq)

	dir := byte(0x00)
	if s.isInitiator {
		dir = 0x01
	}

	var aad bytes.Buffer
	aad.Write(appDataContext)
	aad.Write(s.id)
	var seqBuf [8]byte
	binary.BigEndian.PutUint64(seqBuf[:], seq)
	aad.Write(seqBuf[:])
	aad.WriteByte(dir)

	aead, err := chacha20poly1305.New(s.sendKey)
	if err != nil {
		return 0, nil, err
	}

	ciphertext := aead.Seal(nil, nonce, plaintext, aad.Bytes())
	return seq, ciphertext, nil
}

func (s *Session) DecryptMessage(seq uint64, ciphertext []byte) ([]byte, error) {
	if len(ciphertext) < chacha20poly1305.Overhead {
		return nil, ErrDecryptFailed
	}

	s.receiveMu.Lock()
	defer s.receiveMu.Unlock()

	highest := s.highestReceived
	bitmap := s.replayBitmap

	if seq < highest {
		diff := highest - seq
		if diff >= 64 {
			return nil, ErrReplayDetected
		}
		if (bitmap & (1 << diff)) != 0 {
			return nil, ErrReplayDetected
		}
	} else if seq == highest {
		if (bitmap & 1) != 0 {
			return nil, ErrReplayDetected
		}
	}

	nonce := make([]byte, 12)
	binary.BigEndian.PutUint64(nonce[4:], seq)

	dir := byte(0x01)
	if s.isInitiator {
		dir = 0x00
	}

	var aad bytes.Buffer
	aad.Write(appDataContext)
	aad.Write(s.id)
	var seqBuf [8]byte
	binary.BigEndian.PutUint64(seqBuf[:], seq)
	aad.Write(seqBuf[:])
	aad.WriteByte(dir)

	aead, err := chacha20poly1305.New(s.receiveKey)
	if err != nil {
		return nil, err
	}

	plaintext, err := aead.Open(nil, nonce, ciphertext, aad.Bytes())
	if err != nil {
		return nil, ErrDecryptFailed
	}

	if seq < s.highestReceived {
		diff := s.highestReceived - seq
		s.replayBitmap |= (1 << diff)
	} else {
		if seq == s.highestReceived {
			s.replayBitmap |= 1
		} else {
			diff := seq - s.highestReceived
			if diff >= 64 {
				s.replayBitmap = 1
			} else {
				s.replayBitmap <<= diff
				s.replayBitmap |= 1
			}
			s.highestReceived = seq
		}
	}

	return plaintext, nil
}

func buildTInit(idA, idB, eA []byte) []byte {
	var buf bytes.Buffer
	buf.Write(domainBytes)
	buf.WriteByte(protocolVersion)
	buf.Write(idA)
	buf.Write(idB)
	buf.Write(eA)
	buf.Write(zero32)
	return buf.Bytes()
}

func buildTResp(idA, idB, eA, eB []byte) []byte {
	var buf bytes.Buffer
	buf.Write(domainBytes)
	buf.WriteByte(protocolVersion)
	buf.Write(idA)
	buf.Write(idB)
	buf.Write(eA)
	buf.Write(eB)
	return buf.Bytes()
}

type InitiatorHandshake struct {
	state     InitiatorState
	identity  *NodeIdentity
	peerID    []byte
	ephemeral *EphemeralKey
	eAPub     []byte
	session   *Session
}

func NewInitiatorHandshake(identity *NodeIdentity, peerID []byte) (*InitiatorHandshake, error) {
	if len(peerID) != 32 {
		return nil, ErrMalformedPublicKey
	}
	pid := make([]byte, 32)
	copy(pid, peerID)

	return &InitiatorHandshake{
		state:    InitiatorStateNew,
		identity: identity,
		peerID:   pid,
	}, nil
}

type InitMessage struct {
	EphemeralKey []byte
	Signature    []byte
}

func (h *InitiatorHandshake) GenerateInit() (*InitMessage, error) {
	if h.state != InitiatorStateNew {
		h.state = InitiatorStateFailed
		return nil, ErrInvalidHandshakeState
	}

	eKey, err := GenerateEphemeralKey()
	if err != nil {
		h.state = InitiatorStateFailed
		return nil, err
	}
	h.ephemeral = eKey

	pub, err := eKey.PublicKey()
	if err != nil {
		h.state = InitiatorStateFailed
		return nil, err
	}
	h.eAPub = pub

	tInit := buildTInit(h.identity.PublicKey(), h.peerID, pub)
	sig, err := h.identity.Sign(tInit)
	if err != nil {
		h.state = InitiatorStateFailed
		return nil, err
	}

	h.state = InitiatorStateInitSent
	return &InitMessage{
		EphemeralKey: pub,
		Signature:    sig,
	}, nil
}

type RespMessage struct {
	EphemeralKey []byte
	Signature    []byte
}

func (h *InitiatorHandshake) ProcessResp(msg *RespMessage) error {
	if msg == nil {
		h.state = InitiatorStateFailed
		return ErrInvalidHandshakeMessage
	}
	if h.state != InitiatorStateInitSent {
		h.state = InitiatorStateFailed
		return ErrInvalidHandshakeState
	}

	if len(msg.EphemeralKey) != 32 || len(msg.Signature) != 64 {
		h.state = InitiatorStateFailed
		return ErrInvalidHandshakeMessage
	}

	tResp := buildTResp(h.identity.PublicKey(), h.peerID, h.eAPub, msg.EphemeralKey)
	if err := VerifySignature(h.peerID, tResp, msg.Signature); err != nil {
		h.state = InitiatorStateFailed
		return err
	}

	ss, err := h.ephemeral.ComputeSharedSecret(msg.EphemeralKey)
	if err != nil {
		h.state = InitiatorStateFailed
		return err
	}

	h.session = deriveSession(tResp, ss, true, h.peerID)
	h.ephemeral.Destroy()
	h.ephemeral = nil
	h.state = InitiatorStateEstablished

	return nil
}

func (h *InitiatorHandshake) Session() (*Session, error) {
	if h.state != InitiatorStateEstablished || h.session == nil {
		return nil, ErrSessionNotEstablished
	}
	return h.session, nil
}

type ResponderHandshake struct {
	state     ResponderState
	identity  *NodeIdentity
	peerID    []byte
	ephemeral *EphemeralKey
	eAPub     []byte
	session   *Session
}

func NewResponderHandshake(identity *NodeIdentity, peerID []byte) (*ResponderHandshake, error) {
	if len(peerID) != 32 {
		return nil, ErrMalformedPublicKey
	}
	pid := make([]byte, 32)
	copy(pid, peerID)

	return &ResponderHandshake{
		state:    ResponderStateNew,
		identity: identity,
		peerID:   pid,
	}, nil
}

func (h *ResponderHandshake) ProcessInit(msg *InitMessage) error {
	if msg == nil {
		h.state = ResponderStateFailed
		return ErrInvalidHandshakeMessage
	}
	if h.state != ResponderStateNew {
		h.state = ResponderStateFailed
		return ErrInvalidHandshakeState
	}

	if len(msg.EphemeralKey) != 32 || len(msg.Signature) != 64 {
		h.state = ResponderStateFailed
		return ErrInvalidHandshakeMessage
	}

	tInit := buildTInit(h.peerID, h.identity.PublicKey(), msg.EphemeralKey)
	if err := VerifySignature(h.peerID, tInit, msg.Signature); err != nil {
		h.state = ResponderStateFailed
		return err
	}

	h.eAPub = make([]byte, 32)
	copy(h.eAPub, msg.EphemeralKey)
	h.state = ResponderStateInitReceived
	return nil
}

func (h *ResponderHandshake) GenerateResp() (*RespMessage, error) {
	if h.state != ResponderStateInitReceived {
		h.state = ResponderStateFailed
		return nil, ErrInvalidHandshakeState
	}

	eKey, err := GenerateEphemeralKey()
	if err != nil {
		h.state = ResponderStateFailed
		return nil, err
	}
	h.ephemeral = eKey

	pub, err := eKey.PublicKey()
	if err != nil {
		h.state = ResponderStateFailed
		return nil, err
	}

	tResp := buildTResp(h.peerID, h.identity.PublicKey(), h.eAPub, pub)
	sig, err := h.identity.Sign(tResp)
	if err != nil {
		h.state = ResponderStateFailed
		return nil, err
	}

	ss, err := h.ephemeral.ComputeSharedSecret(h.eAPub)
	if err != nil {
		h.state = ResponderStateFailed
		return nil, err
	}

	h.session = deriveSession(tResp, ss, false, h.peerID)
	h.ephemeral.Destroy()
	h.ephemeral = nil
	h.state = ResponderStateEstablished

	return &RespMessage{
		EphemeralKey: pub,
		Signature:    sig,
	}, nil
}

func (h *ResponderHandshake) Session() (*Session, error) {
	if h.state != ResponderStateEstablished || h.session == nil {
		return nil, ErrSessionNotEstablished
	}
	return h.session, nil
}

func deriveSession(tResp []byte, ss []byte, isInitiator bool, peerID []byte) *Session {
	tRespHd := sha256.Sum256(tResp)
	salt := make([]byte, 32)
	copy(salt, tRespHd[:])

	prk := hkdf.Extract(sha256.New, ss, salt)

	ka2bReader := hkdf.Expand(sha256.New, prk, []byte("MeshChat-v1|session|initiator->responder"))
	kb2aReader := hkdf.Expand(sha256.New, prk, []byte("MeshChat-v1|session|responder->initiator"))

	ka2b := make([]byte, 32)
	kb2a := make([]byte, 32)
	io.ReadFull(ka2bReader, ka2b)
	io.ReadFull(kb2aReader, kb2a)

	pid := make([]byte, len(peerID))
	copy(pid, peerID)

	s := &Session{
		peerID:      pid,
		id:          salt,
		isInitiator: isInitiator,
	}

	if isInitiator {
		s.sendKey = ka2b
		s.receiveKey = kb2a
	} else {
		s.sendKey = kb2a
		s.receiveKey = ka2b
	}

	return s
}

func (h *InitiatorHandshake) TestResp(msg *RespMessage) bool {
	if msg == nil || h.state != InitiatorStateInitSent {
		return false
	}
	if len(msg.EphemeralKey) != 32 || len(msg.Signature) != 64 {
		return false
	}
	tResp := buildTResp(h.identity.PublicKey(), h.peerID, h.eAPub, msg.EphemeralKey)
	err := VerifySignature(h.peerID, tResp, msg.Signature)
	return err == nil
}
