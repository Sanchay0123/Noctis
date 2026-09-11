// Package protocol defines the wire packet structures and serialization formats.
//
// Unknown-Field Policy:
// The architecture requires security-sensitive Protobuf messages to reject unknown fields.
// The `google.golang.org/protobuf` v2 API does not natively reject unknown fields during Unmarshal.
// Therefore, the runtime unmarshal operation MUST use explicit unknown-field rejection by checking
// for residual data after unmarshaling:
//
//	err := proto.Unmarshal(data, msg)
//	if err == nil && len(msg.ProtoReflect().GetUnknown()) > 0 {
//	    return ErrUnknownFieldsRejected
//	}
//
// Field Validation Policy:
// The Protobuf schema itself cannot enforce semantic byte-length constraints.
// The following semantic validation requirements MUST be enforced during later protocol implementation:
// - `packet_id` = exactly 16 bytes
// - `source_node` = exactly 32 bytes
// - `dest_node` = exactly 32 bytes
// - `InitPayload.ephemeral_key` = exactly 32 bytes
// - `InitPayload.signature` = exactly 64 bytes
// - `RespPayload.ephemeral_key` = exactly 32 bytes
// - `RespPayload.signature` = exactly 64 bytes
// - `AppDataPayload.session_id` = exactly 32 bytes
// - `ciphertext` must satisfy the later AEAD minimum-size requirement
//
// Packet Type / Payload Consistency Invariants:
// - `PACKET_TYPE_INIT` MUST map to `init` payload
// - `PACKET_TYPE_RESP` MUST map to `resp` payload
// - `PACKET_TYPE_APP_DATA` MUST map to `app_data` payload
// - `PACKET_TYPE_UNKNOWN` MUST be rejected during protocol validation
package protocol
