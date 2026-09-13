package mesh

import "encoding/hex"

// TelemetryRecorder allows optional, out-of-band monitoring of mesh events.
type TelemetryRecorder interface {
	RecordPeerConnected(isOutbound bool)
	RecordPeerDisconnected()
	RecordHandshakeFailed(err error)
	RecordDuplicateConnection()
	RecordResourceLimitReached(limitName string)
}

// noopTelemetry is used when observability is nil to prevent crashes.
type noopTelemetry struct{}

func (n noopTelemetry) RecordPeerConnected(isOutbound bool) {}
func (n noopTelemetry) RecordPeerDisconnected()             {}
func (n noopTelemetry) RecordHandshakeFailed(err error)     {}
func (n noopTelemetry) RecordDuplicateConnection()          {}
func (n noopTelemetry) RecordResourceLimitReached(l string) {}

func SanitizeID(id []byte) string {
	if len(id) > 32 {
		id = id[:32]
	}
	return hex.EncodeToString(id)
}
