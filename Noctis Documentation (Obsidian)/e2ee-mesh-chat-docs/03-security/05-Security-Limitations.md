# Security Limitations

## Metadata exposure

Even with E2EE, a relay may observe:

-   who connects to whom at the routing layer
-   timing
-   packet sizes
-   frequency
-   route information
-   node identifiers

Therefore the system is **not anonymous** by default.

## Malicious routing

A relay may still:

-   drop packets
-   delay packets
-   refuse forwarding
-   provide incorrect routing information
-   attempt resource exhaustion

Cryptographic message integrity does not guarantee delivery.

## Endpoint compromise

If Alice's endpoint is compromised, an attacker may access plaintext
before encryption or after decryption.

If Bob's endpoint is compromised, the attacker may access plaintext
after decryption.

E2EE does not protect against a compromised endpoint.

## Local key storage

The prototype may use local persistent key storage. The exact protection
mechanism must be documented once implemented.

If private keys are stored unencrypted on disk, that must be disclosed
as a limitation.

## Forward secrecy

Do not claim production-grade forward secrecy unless the final session
protocol and key lifecycle actually provide and demonstrate it.

## Deniability

The prototype does not automatically provide deniability merely because
it uses modern cryptography.

## Availability

Encryption does not prevent a malicious relay from dropping messages.

## Security-language policy

Forbidden claims unless independently demonstrated:

-   "unhackable"
-   "perfectly secure"
-   "anonymous"
-   "military-grade"
-   "completely secure"

See [06-Security-Review-Checklist](06-Security-Review-Checklist.md).

## M6 Limitations

M6 provides bounded flooding, not a full routing-security protocol. Unique-packet floods can still cause CPU, bandwidth and cache churn. Shared per-peer queues provide no QoS priority, so routed flooding can compete with application traffic. Delivery is best-effort without end-to-end acknowledgements at this layer.

Routing metadata remains visible to relays. The design therefore does not provide anonymity, strong traffic-analysis resistance or complete metadata hiding. Malicious relays can also drop or selectively forward packets.


## M7 Security-Hardening Limitations

M7 hardens local runtime security but does not provide a complete network-wide
availability solution.

Telemetry at the M7 gate consists of runtime hooks/no-op behavior rather than
a concrete Prometheus exporter. M8 is responsible for implementing dashboards
and concrete metrics under the bounded-cardinality rule.

The system still exposes routing identities and other traffic metadata to
relays. No anonymity or complete metadata hiding is claimed.

Global Sybil attacks, large-scale flood traffic, malicious packet dropping and
endpoint compromise remain outside the guarantees of the prototype.


## M9 Demonstration Limitations

M9 provides evidence for the demonstrated E2EE relay boundary and selected
adversarial controls, but does not expand the system's security guarantees.
The relay can still observe routing metadata, traffic timing and packet sizes.
M9 therefore does not establish anonymity or complete metadata hiding. The GUI
demonstration was manually verified on the Project Overseer's desktop; automated
GUI interaction was unavailable in the headless implementation environment.
Docker regression remains unverified because of external proxy/DNS restrictions.
A concrete Prometheus exporter remains not implemented.
