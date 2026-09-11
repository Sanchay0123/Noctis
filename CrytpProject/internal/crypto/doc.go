// Package crypto defines cryptographic identities, session establishment, and authenticated encryption.
// It sits below the application layer and strictly encapsulates private keys and symmetric session keys.
// Neither raw private keys nor session keys are exposed to the transport, routing, or observability layers.
package crypto
