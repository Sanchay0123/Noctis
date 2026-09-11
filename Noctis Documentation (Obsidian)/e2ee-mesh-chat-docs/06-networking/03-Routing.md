# Routing

## Strategy

**Managed flooding** is approved for the educational prototype.

It is intentionally simpler than a full distance-vector or link-state
protocol.

## Forwarding

A node that receives a packet not destined for itself may forward it to
eligible peers except the incoming peer.

Forwarding must occur only after basic packet validation and duplicate
checking.

## Duplicate detection

Every packet has a sender-generated PacketID.

Each node maintains a bounded seen-cache.

The cache must define:

- maximum number of entries
- expiration policy
- eviction policy
- memory limit

A fixed 60-second lifetime is not automatically sufficient. It must be
consistent with the maximum expected packet lifetime.

## TTL

Packets have a finite hop limit.

On receipt:

1. reject invalid TTL values
2. if the packet is for the local node, process it as appropriate
3. otherwise decrement TTL before forwarding
4. do not forward when the resulting TTL is zero

## Flooding risks

Managed flooding does not prevent a malicious node from injecting many
unique PacketIDs.

The implementation should therefore include reasonable:

- packet size limits
- peer limits
- connection limits
- forwarding/backpressure limits
- rate limits where justified

These are DoS mitigations, not a claim of complete DoS protection.

## Loop prevention

PacketID caching prevents repeated forwarding of the same packet through
cycles.

Sequence numbers are unrelated to mesh duplicate suppression.

## Relay confidentiality

Routing logic must not require access to:

- application plaintext
- recipient session keys
- sender session keys

See [[04-protocol/04-Message-Format]].
