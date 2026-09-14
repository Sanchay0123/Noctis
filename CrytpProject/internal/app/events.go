package app

// AppEventType defines the safe/sanitized categories of events the GUI can handle.
type AppEventType int

const (
	EventTypeMessageReceived AppEventType = iota
	EventTypePeerConnected
	EventTypePeerDisconnected
	EventTypeSecurityAlert
	EventTypeSessionEstablished
	EventTypeSessionFailed
	EventTypeConnectionFailed
)

// AppEvent is the bounded structure pushed to the GUI.
type AppEvent struct {
	Type        AppEventType
	PeerID      string
	Message     string
	AlertReason string
}

func EventMessageReceived(peerID string, plaintext string) AppEvent {
	return AppEvent{
		Type:    EventTypeMessageReceived,
		PeerID:  peerID,
		Message: plaintext,
	}
}

func EventPeerConnected(peerID string) AppEvent {
	return AppEvent{
		Type:   EventTypePeerConnected,
		PeerID: peerID,
	}
}

func EventPeerDisconnected(peerID string) AppEvent {
	return AppEvent{
		Type:   EventTypePeerDisconnected,
		PeerID: peerID,
	}
}

func EventSecurityAlert(peerID string, reason string) AppEvent {
	return AppEvent{
		Type:        EventTypeSecurityAlert,
		PeerID:      peerID,
		AlertReason: reason,
	}
}
