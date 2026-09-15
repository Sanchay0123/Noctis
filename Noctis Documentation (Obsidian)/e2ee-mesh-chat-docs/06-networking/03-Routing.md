# Routing

## Strategy

**Managed flooding** is approved for the educational prototype.

## Forwarding algorithm

For each received packet:

```text
1. Validate frame size and basic packet structure.
2. Extract PacketID.
3. Check bounded seen-cache.
4. If already seen: drop.
5. Record PacketID according to cache policy.
6. If destination is local: pass to protocol/application handling.
7. Otherwise:
   a. if TTL <= 1: drop
   b. decrement TTL
   c. forward to eligible peers
```

The implementation must ensure duplicate suppression occurs before
forwarding.

## TTL

Initial TTL is `16`.

A packet arriving at its destination with TTL `1` may be processed.

A non-destination packet with TTL `1` is not forwarded.

A packet with invalid/zero TTL must not be forwarded.

## Seen-cache

The cache is bounded by:

- maximum entries
- expiration time
- eviction policy

The initial target may be 10,000 entries with a 60-second expiration,
but these are configurable engineering parameters and must be validated
against the expected demonstration topology and packet lifetime.

A cache-full condition must have deterministic behavior.

## Flooding resource controls

At minimum define:

- maximum peers/connections
- maximum frame size
- maximum forwarding queue
- backpressure behavior
- per-connection read/write limits

Optional rate limiting may be added during hardening.

Managed flooding does not eliminate denial-of-service risk.

## Relay confidentiality

Relays operate on routing metadata only.

They do not decrypt application data and must not inspect application
plaintext to make routing decisions.

## Bootstrap peers

Initial implementation uses explicit bootstrap peers.

UDP broadcast discovery is deferred.

## Testing

Required tests:

- linear topology
- ring topology
- duplicate PacketID
- TTL expiry
- TTL == 1 at destination
- oversized frame
- malformed Protobuf
- cache saturation
- forwarding queue pressure

## M6 Final Implementation

M6 freezes the routing strategy as **bounded managed flooding**. The cache key is the pair `(source_node, packet_id)` where `packet_id` is exactly 16 bytes. The cache is bounded to 10,000 entries with a 2-minute lifetime.

The default initial TTL is 16 and the maximum accepted TTL is 32. TTL 0 is dropped; TTL 1 is deliverable only when the local node is the destination; TTL greater than 1 is delivered locally or decremented exactly once before remote forwarding.

Forwarding is to eligible active peers except the incoming peer. The router performs no network I/O directly and uses bounded per-peer queues. See [07-M6-Handshake-Correlation](07-M6-Handshake-Correlation.md) and [08-M6-Resource-Bounds](08-M6-Resource-Bounds.md).

**Status: 🟢 COMPLETE / VERIFIED / APPROVED**
