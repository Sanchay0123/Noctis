// Package crypto defines cryptographic identities, session establishment, and authenticated encryption.
// It sits below the application layer and strictly encapsulates private keys and symmetric session keys.
// Neither raw private keys nor session keys are exposed to the transport, routing, or observability layers.
//
// M2.1 Updates:
// - Ed25519 is used as the long-term identity primitive.
// - 32-byte public keys are the public identity representation.
// - Private key handling boundary is strictly enforced.
// - Signing and verification are implemented with clear distinction between identity and future ephemeral X25519 session keys.
package crypto
