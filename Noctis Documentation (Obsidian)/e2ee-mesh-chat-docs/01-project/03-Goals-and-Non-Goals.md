# Goals and Non-Goals

## Goals

The project aims to demonstrate:

- authenticated cryptographic identities
    
- authenticated X25519-based key establishment
    
- HKDF-based session key derivation
    
- AEAD-protected application messages
    
- replay protection
    
- multi-hop forwarding
    
- separation between routing and application cryptography
    
- basic resilience to malicious packet manipulation
    
- reproducible security demonstrations
    
- security-safe observability and runtime telemetry
    
- reproducible containerized execution
    
- traceable implementation and documentation through GitHub
    

## Non-goals

The prototype does **not** promise:

- perfect anonymity
    
- complete metadata hiding
    
- resistance to global traffic analysis
    
- production-grade deniability
    
- protection of a compromised endpoint
    
- protection against every routing attack
    
- production-grade availability
    
- complete forward secrecy unless specifically demonstrated
    
- military-grade or "unhackable" security
    

## Why non-goals matter

A security project is stronger academically when it precisely states  
what it does **not** solve.

For example, a relay may be unable to decrypt a message while still  
learning:

- that two nodes are communicating
    
- packet timing
    
- packet size
    
- source/destination routing identifiers
    
- path-related information
    

Therefore:

> **Confidentiality of message content is not equivalent to anonymity.**

See [[03-security/05-Security-Limitations]].

## Observability boundary

Observability improves visibility into the system but does not expand  
the security claims.

The project does not aim to provide complete metadata privacy through  
its telemetry system. Telemetry must instead minimize unnecessary  
metadata and never expose application plaintext or cryptographic  
secrets.