# Security Controls

Threat Control

---

Packet modification AEAD authentication  
Identity substitution Ed25519-authenticated key exchange  
Replay Message identifiers + sequence/replay state  
Route loops TTL + duplicate detection  
Malformed packet Strict validation  
Key exposure through logs Secret-aware logging policy  
Nonce reuse Explicit nonce/key lifecycle  
Unknown protocol versions Version validation  
Resource exhaustion Size limits, connection limits, bounded state  
Relay plaintext exposure End-to-end encryption

## Defense-in-depth

No single control should be treated as solving every threat.

For example:

- TLS-like transport protection does not replace E2EE.
    
- signatures do not provide confidentiality.
    
- AEAD does not solve routing loops.
    
- TTL does not authenticate a sender.
    
- replay detection does not prevent packet dropping.
    

## Observability controls

|Threat|Control|
|---|---|
|Plaintext log leakage|Explicit telemetry sanitization|
|Key leakage|Secret-aware logging; keys excluded from event payloads|
|Metadata leakage|Minimize sensitive labels and routing metadata|
|Log injection|Structured logging with validated fields|
|Log/telemetry DoS|Bounded queues, sizes and rates|
|Monitoring outage|Observability decoupled from message path|
|False security evidence|Treat protocol/test evidence as authoritative|
|Dashboard exposure|Restrict observability interfaces in deployment|

The dashboard is a visualization of evidence; it is not itself a  
security control.

See [[02-architecture/07-Observability-and-Security-Telemetry]].