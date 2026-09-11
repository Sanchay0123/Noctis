// Package mesh provides multi-hop routing, duplicate packet suppression, and TTL enforcement.
// Relay nodes route packets solely by examining unencrypted routing metadata (source, dest, packet_id, ttl).
// Intermediate nodes cannot decrypt ciphertext payloads or access cryptographic session keys.
package mesh
