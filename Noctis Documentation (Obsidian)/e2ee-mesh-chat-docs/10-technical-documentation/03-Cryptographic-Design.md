# Cryptographic Design — Implemented Protocol

## Status

**Implemented and validated through M3 and integrated through M9.**

## Primitive roles

| Primitive | Role |
|---|---|
| Ed25519 | Long-term node identity and signatures |
| X25519 | Fresh ephemeral Diffie–Hellman key agreement |
| HKDF-SHA-256 | Session key derivation |
| SHA-256 | Session identifier and handshake KDF salt |
| ChaCha20-Poly1305 | Authenticated application encryption |

No cryptographic primitive is implemented from scratch.

## Identity

Each node owns a long-term Ed25519 keypair. The public key is the node's
cryptographic identity. The private key is secret and is not an application or
routing payload.

Ed25519 is used for authentication/signatures; it is not reused as an X25519
key.

## Handshake

The frozen handshake domain is:

```text
DOMAIN = "MeshChat-Handshake-v1"
PROTOCOL_VERSION = 0x01
ZERO32 = 32 zero bytes
```

For initiator identity `ID_A`, responder identity `ID_B`, and ephemeral public
keys `E_A`, `E_B`:

```text
T_INIT = DOMAIN || u8(PROTOCOL_VERSION) || ID_A || ID_B || E_A || ZERO32
T_RESP = DOMAIN || u8(PROTOCOL_VERSION) || ID_A || ID_B || E_A || E_B
```

The initiator signs `T_INIT`. The responder verifies it and signs `T_RESP`.
The initiator verifies the responder signature.

## Shared secret and KDF

After authentication:

```text
SS = X25519(ephemeral_private, peer_ephemeral_public)
session_id = SHA256(T_RESP)
salt = SHA256(T_RESP)
PRK = HKDF-Extract(salt, SS)
```

Directional keys are derived with distinct labels:

```text
K_A_to_B = HKDF-Expand(
    PRK,
    "MeshChat-v1|session|initiator->responder",
    32
)

K_B_to_A = HKDF-Expand(
    PRK,
    "MeshChat-v1|session|responder->initiator",
    32
)
```

The full 32-byte transcript digest is used as the session identifier.

## Application AEAD

Application messages use ChaCha20-Poly1305.

Nonce construction:

```text
nonce = 4 zero bytes || uint64_be(sequence_number)
```

The application AAD is:

```text
"MeshChat-AppData-v1" || session_id || uint64_be(sequence_number) || direction
```

Direction is derived from the established role:

```text
0x01 = initiator -> responder
0x00 = responder -> initiator
```

The caller cannot independently choose the direction marker.

## Replay protection

The receive side maintains a 64-message replay window. Replay state is
updated only after successful AEAD authentication. Duplicate or out-of-window
messages are rejected according to the session receive policy.

The first outbound sequence is zero. Sequence exhaustion is treated as an
error rather than wrapping the nonce counter.

## Key lifecycle

Fresh ephemeral keys are generated for session establishment. Session keys
remain in memory for the active session. Restart invalidates previous session
state and a fresh handshake establishes new keys.

The implementation does not claim post-compromise security. Forward secrecy is
conditional on fresh ephemeral keys, appropriate private-key erasure, and
uncompromised endpoints.

## Trust model

Identity public keys may be supplied through a managed/pre-shared directory or
TOFU. TOFU alone does not authenticate first contact against an active MITM.

## Evidence

M3 known-answer and negative tests validate transcript construction, session
identifier derivation, HKDF outputs, identity/ephemeral substitution handling,
role separation, malformed input handling, AEAD integrity, nonce/sequence
behavior, replay windows, and concurrency safety. M9 extends the evidence to
multi-hop relay behavior and adversarial demonstrations.

## Related records

- [[04-protocol/03-Session-Establishment]]
- [[04-protocol/04-Message-Format]]
- [[05-cryptography/01-Cryptographic-Architecture]]
- [[05-cryptography/04-Key-Derivation]]
- [[05-cryptography/05-AEAD-and-Nonce-Management]]
- [[03-security/05-Security-Limitations]]
