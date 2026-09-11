# Abstract --- Draft

> **Status: Draft. Do not submit until implementation and results are
> available.**

This project presents the design and implementation of an educational
end-to-end encrypted messaging system operating over a decentralized
mesh network. The system combines cryptographic identity, authenticated
key establishment, authenticated encryption, replay protection and
multi-hop peer-to-peer forwarding. Unlike conventional demonstrations
that focus only on direct encrypted communication, the proposed system
separates the mesh routing layer from the application cryptographic
layer, allowing intermediate nodes to forward encrypted packets without
accessing the underlying message content.

The prototype uses established cryptographic primitives rather than
custom cryptography and evaluates the system through functional,
integration and adversarial tests. Security demonstrations focus on
message confidentiality, message integrity, replay resistance, basic
impersonation resistance and the behavior of malicious or failed relay
nodes.

The project is intentionally scoped as an academic prototype. It does
not claim perfect anonymity, complete metadata protection, protection
against compromised endpoints, or universal resistance to routing
attacks. The resulting system is intended to demonstrate the practical
interaction between applied cryptography, computer networking and
distributed systems while maintaining explicit and testable security
boundaries.

**Note:** Replace generic claims with measured implementation results
before final submission.
