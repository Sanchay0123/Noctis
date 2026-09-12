package mesh

// TelemetryRecorder allows optional, out-of-band monitoring of mesh events.
type TelemetryRecorder interface {
	RecordPeerConnected(remoteID []byte, isOutbound bool)
	RecordPeerDisconnected(remoteID []byte)
	RecordHandshakeFailed(remoteID []byte, err error)
	RecordDuplicateConnection(remoteID []byte)
	RecordResourceLimitReached(limitName string)
}

// noopTelemetry is used when observability is nil to prevent crashes.
type noopTelemetry struct{}

func (n noopTelemetry) RecordPeerConnected(remoteID []byte, isOutbound bool) {}
func (n noopTelemetry) RecordPeerDisconnected(remoteID []byte)               {}
func (n noopTelemetry) RecordHandshakeFailed(remoteID []byte, err error)     {}
func (n noopTelemetry) RecordDuplicateConnection(remoteID []byte)            {}
func (n noopTelemetry) RecordResourceLimitReached(limitName string)          {}
