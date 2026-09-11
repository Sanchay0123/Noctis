// Package observability provides metrics, structured logging, and health telemetry.
// Observability is strictly out-of-band and non-blocking; its availability never affects routing or encryption.
// Telemetry MUST NEVER receive, log, or export plaintext, private keys, or session keys.
package observability
