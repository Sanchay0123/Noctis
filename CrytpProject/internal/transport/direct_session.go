package transport

import (
	"bytes"
	"github.com/sanchayjain/meshchat/internal/crypto"
	"github.com/sanchayjain/meshchat/internal/protocol"
)

type DirectChannel struct {
	conn       *Connection
	localIdent *crypto.NodeIdentity
	remoteID   []byte
	session    *crypto.Session
}

func NewDirectChannel(conn *Connection, localIdent *crypto.NodeIdentity) *DirectChannel {
	return &DirectChannel{
		conn:       conn,
		localIdent: localIdent,
	}
}

// InitiateHandshake runs the initiator side of the M3 handshake over the direct connection.
func (ch *DirectChannel) InitiateHandshake(remoteID []byte) error {
	hs, err := crypto.NewInitiatorHandshake(ch.localIdent, remoteID)
	if err != nil {
		return err
	}

	initMsg, err := hs.GenerateInit()
	if err != nil {
		return err
	}

	pktID, err := GeneratePacketID()
	if err != nil {
		return err
	}

	initPkt := &protocol.MeshPacket{
		Version:    1,
		Type:       protocol.PacketType_PACKET_TYPE_INIT,
		PacketId:   pktID,
		Ttl:        1,
		SourceNode: ch.localIdent.PublicKey(),
		DestNode:   remoteID,
		Payload: &protocol.MeshPacket_Init{
			Init: &protocol.InitPayload{
				EphemeralKey: initMsg.EphemeralKey,
				Signature:    initMsg.Signature,
			},
		},
	}

	if err := ch.conn.WritePacket(initPkt); err != nil {
		return err
	}

	respPkt, err := ch.conn.ReadPacket()
	if err != nil {
		return err
	}

	// Validate it's a response intended for us
	if respPkt.Type != protocol.PacketType_PACKET_TYPE_RESP {
		return ErrInvalidPacketType
	}
	if !bytes.Equal(respPkt.DestNode, ch.localIdent.PublicKey()) {
		return ErrWrongDestination
	}
	if !bytes.Equal(respPkt.SourceNode, remoteID) {
		return ErrWrongSource
	}

	respPayload := respPkt.Payload.(*protocol.MeshPacket_Resp).Resp
	err = hs.ProcessResp(&crypto.RespMessage{
		EphemeralKey: respPayload.EphemeralKey,
		Signature:    respPayload.Signature,
	})
	if err != nil {
		return err
	}

	ch.remoteID = remoteID
	ch.session, err = hs.Session()
	return err
}

// AcceptHandshake runs the responder side of the M3 handshake over the direct connection.
func (ch *DirectChannel) AcceptHandshake(expectedRemoteID []byte) error {
	hs, err := crypto.NewResponderHandshake(ch.localIdent, expectedRemoteID)
	if err != nil {
		return err
	}

	initPkt, err := ch.conn.ReadPacket()
	if err != nil {
		return err
	}

	if initPkt.Type != protocol.PacketType_PACKET_TYPE_INIT {
		return ErrInvalidPacketType
	}
	if !bytes.Equal(initPkt.DestNode, ch.localIdent.PublicKey()) {
		return ErrWrongDestination
	}
	if !bytes.Equal(initPkt.SourceNode, expectedRemoteID) {
		return ErrWrongSource
	}

	initPayload := initPkt.Payload.(*protocol.MeshPacket_Init).Init
	err = hs.ProcessInit(&crypto.InitMessage{
		EphemeralKey: initPayload.EphemeralKey,
		Signature:    initPayload.Signature,
	})
	if err != nil {
		return err
	}

	respMsg, err := hs.GenerateResp()
	if err != nil {
		return err
	}

	pktID, err := GeneratePacketID()
	if err != nil {
		return err
	}

	respPkt := &protocol.MeshPacket{
		Version:    1,
		Type:       protocol.PacketType_PACKET_TYPE_RESP,
		PacketId:   pktID,
		Ttl:        1,
		SourceNode: ch.localIdent.PublicKey(),
		DestNode:   expectedRemoteID,
		Payload: &protocol.MeshPacket_Resp{
			Resp: &protocol.RespPayload{
				EphemeralKey: respMsg.EphemeralKey,
				Signature:    respMsg.Signature,
			},
		},
	}

	if err := ch.conn.WritePacket(respPkt); err != nil {
		return err
	}

	ch.remoteID = expectedRemoteID
	ch.session, err = hs.Session()
	return err
}

// SendPlaintext encrypts a plaintext application message and dispatches it.
func (ch *DirectChannel) SendPlaintext(plaintext []byte) error {
	if ch.session == nil {
		return ErrSessionNotEstablished
	}

	seq, ct, err := ch.session.EncryptMessage(plaintext)
	if err != nil {
		return err
	}

	pktID, err := GeneratePacketID()
	if err != nil {
		return err
	}

	appPkt := &protocol.MeshPacket{
		Version:    1,
		Type:       protocol.PacketType_PACKET_TYPE_APP_DATA,
		PacketId:   pktID,
		Ttl:        1,
		SourceNode: ch.localIdent.PublicKey(),
		DestNode:   ch.remoteID,
		Payload: &protocol.MeshPacket_AppData{
			AppData: &protocol.AppDataPayload{
				SessionId:   ch.session.ID(),
				SequenceNum: seq,
				Ciphertext:  ct,
			},
		},
	}

	return ch.conn.WritePacket(appPkt)
}

// ReceivePlaintext blocks and receives the next authenticated plaintext message.
func (ch *DirectChannel) ReceivePlaintext() ([]byte, error) {
	for {
		pkt, err := ch.conn.ReadPacket()
		if err != nil {
			return nil, err
		}

		// Reject non-AppData immediately in the app stream.
		if pkt.Type != protocol.PacketType_PACKET_TYPE_APP_DATA {
			return nil, ErrInvalidPacketType
		}

		// We MUST enforce session establishment
		if ch.session == nil {
			return nil, ErrSessionNotEstablished
		}

		// Verify packet destination matches our identity
		if !bytes.Equal(pkt.DestNode, ch.localIdent.PublicKey()) {
			return nil, ErrWrongDestination
		}

		// Verify packet source matches the established remote identity
		if !bytes.Equal(pkt.SourceNode, ch.remoteID) {
			return nil, ErrWrongSource
		}

		appPayload := pkt.Payload.(*protocol.MeshPacket_AppData).AppData

		// Verify session ID matches our active session
		if !bytes.Equal(appPayload.SessionId, ch.session.ID()) {
			return nil, ErrSessionMismatch
		}

		// Attempt AEAD Decryption / Authentication.
		// If it fails, plaintext is never delivered.
		plaintext, err := ch.session.DecryptMessage(appPayload.SequenceNum, appPayload.Ciphertext)
		if err != nil {
			return nil, ErrDecryptionFailed
		}

		return plaintext, nil
	}
}

// Close gracefully shuts down the underlying connection
func (ch *DirectChannel) Close() error {
	return ch.conn.Close()
}
