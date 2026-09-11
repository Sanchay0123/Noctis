# Session Establishment

## Status

**M3 implementation blocked pending final cryptographic specification.**

The revised architecture has the correct high-level model:

- Ed25519 = long-term identity/signature
- X25519 = fresh ephemeral key agreement
- HKDF-SHA-256 = session key derivation
- ChaCha20-Poly1305 = application encryption

## Identity trust

A node must already know and authenticate the expected peer's Ed25519
public key through the project's configured out-of-band trust mechanism.

Trust-on-First-Use is an application policy, not cryptographic
authentication by itself.

## Canonical transcript

The handshake must use one exact byte representation.

Define:

```text
DOMAIN = "MeshChat-Handshake-v1"

T_INIT =
    DOMAIN ||
    u8(PROTOCOL_VERSION) ||
    ID_A ||
    ID_B ||
    E_A ||
    ZERO32

T_RESP =
    DOMAIN ||
    u8(PROTOCOL_VERSION) ||
    ID_A ||
    ID_B ||
    E_A ||
    E_B
```

Where:

- `ID_A` = exactly 32 bytes Ed25519 public key
- `ID_B` = exactly 32 bytes Ed25519 public key
- `E_A` = exactly 32 bytes X25519 public key
- `E_B` = exactly 32 bytes X25519 public key
- `ZERO32` = exactly 32 zero bytes
- `PROTOCOL_VERSION` = `1`

`DOMAIN` is a fixed protocol constant and must have exactly one defined
byte representation in the implementation.

The implementation must not use ambiguous string concatenation.

## INIT

Alice:

1. generates a fresh X25519 ephemeral keypair
2. constructs `T_INIT`
3. signs `T_INIT` with Alice's Ed25519 identity key
4. sends INIT containing Alice's identity, Bob's identity, ephemeral
   public key and signature

Bob:

1. validates packet structure and sizes
2. confirms the requested responder identity is Bob
3. obtains Alice's trusted Ed25519 public key
4. verifies the INIT signature
5. checks the handshake replay/duplicate state
6. generates a fresh X25519 ephemeral keypair

## RESP

Bob constructs:

```text
T_RESP =
    DOMAIN ||
    u8(PROTOCOL_VERSION) ||
    ID_A ||
    ID_B ||
    E_A ||
    E_B
```

Bob signs `T_RESP` using Bob's Ed25519 identity key.

Bob then computes:

```text
SS = X25519(e_B, E_A)
```

Bob derives session keys and sends RESP.

## Alice finalization

Alice:

1. validates the RESP structure
2. confirms the responder identity is Bob
3. confirms the returned initiator identity is Alice
4. confirms `E_A` matches the outstanding handshake
5. verifies Bob's signature over the exact `T_RESP`
6. computes `SS = X25519(e_A, E_B)`
7. derives the same session keys
8. securely erases `e_A`
9. marks the session established

## Session identifier

**Use the full transcript hash.**

```text
session_id = SHA256(T_RESP)
```

Therefore `session_id` is exactly 32 bytes.

Do not truncate it to 32 bits. A 4-byte identifier is insufficient for
robust session-state separation.

## Key derivation

Use HKDF-SHA-256:

```text
salt = SHA256(T_RESP)

PRK = HKDF-Extract(salt, SS)

K_A_to_B =
    HKDF-Expand(
        PRK,
        "MeshChat-v1|session|initiator->responder",
        32
    )

K_B_to_A =
    HKDF-Expand(
        PRK,
        "MeshChat-v1|session|responder->initiator",
        32
    )
```

The exact byte strings are protocol constants.

Both parties therefore derive two distinct directional ChaCha20-Poly1305
keys.

## Handshake replay state

A bounded replay/duplicate cache is required.

A cache entry must identify the handshake sufficiently to distinguish:

- peer identity
- requested local identity
- initiator ephemeral public key

A duplicate INIT for a completed/active handshake must not silently
create an unrelated replacement session.

The implementation must define bounded capacity and expiration.

A timeout clears the outstanding handshake state and requires a fresh
ephemeral key for a new attempt.

## Concurrent handshakes

If multiple handshakes for the same peer are permitted, each must have
independent state and an unambiguous identifier.

If only one outstanding handshake per peer is permitted, additional INITs
must be handled deterministically and documented.

## Key erasure

Ephemeral private keys must be erased as soon as protocol processing no
longer requires them.

Session keys remain in memory only for the active session.

Application restart destroys session state and requires a fresh
handshake.

## Forward secrecy limitation

The design does not provide Signal-style post-compromise security.

Compromise of an active session key compromises messages protected by that
key until the session is replaced.

Do not claim PCS or protection against compromised endpoints.

See [[04-protocol/06-State-Machines]] and
[[05-cryptography/06-Key-Lifecycle]].
