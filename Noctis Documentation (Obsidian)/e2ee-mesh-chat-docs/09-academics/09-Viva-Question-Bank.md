# Viva Question Bank

## Architecture

### Why separate routing from cryptography?

Because relay nodes need routing information but should not require
application plaintext or E2EE keys.

### Why not encrypt the whole routing packet?

Intermediate nodes need some routing information to deliver it.
Encrypting all metadata would require a different routing architecture.

### Can a relay read the message?

It should not be able to decrypt the application payload if the E2EE
design is implemented correctly.

## Cryptography

### Why Ed25519?

It provides established public-key signatures suitable for persistent
cryptographic identity.

### Why X25519?

It provides established elliptic-curve Diffie-Hellman key agreement.

### Why do we need both?

They solve different problems: Ed25519 authenticates identity; X25519
establishes shared secret material.

### Why HKDF?

To derive purpose-specific keys from shared secret material with
explicit context separation.

### Why ChaCha20-Poly1305?

It provides authenticated encryption and avoids unnecessary support for
multiple AEAD schemes in the prototype.

### Why not implement AES/ChaCha ourselves?

Cryptographic primitives are security-critical and should come from
reviewed libraries.

## Security

### What happens if a packet is modified?

AEAD authentication should fail and the message should not be delivered.

### What happens if a packet is replayed?

Receiver replay state should detect and reject it.

### Can a malicious relay drop packets?

Yes. Confidentiality/integrity do not guarantee availability.

### Is the system anonymous?

No. Routing and transport metadata may remain visible.

### What if Alice's computer is compromised?

E2EE cannot protect plaintext already exposed at a compromised endpoint.

### What if a private identity key is stolen?

The attacker may impersonate that identity depending on the protocol
state and key-revocation model. The final implementation must document
this precisely.

## Networking

### Why TTL?

To bound packet lifetime and prevent indefinite routing loops.

### Why duplicate detection?

To prevent repeated forwarding of the same packet.

### What does a relay actually do?

Validate the mesh envelope, enforce routing controls, choose a next hop
and forward the packet.

## Testing

### Why negative tests?

A security property is meaningful only if attempts to violate it are
rejected.

### Is one passing encryption test enough?

No. We need tampering, wrong-key, metadata modification, nonce and
replay tests as appropriate.

## Academic integrity

### What is the most important limitation?

The answer must reflect the actual implementation and evidence, not an
aspirational design.
