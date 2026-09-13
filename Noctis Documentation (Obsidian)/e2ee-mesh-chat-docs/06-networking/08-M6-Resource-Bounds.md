# M6 Resource Bounds

## Duplicate cache

- Key: `(source_node, packet_id)`.
- Source identity component: 32 bytes.
- PacketID: exactly 16 bytes.
- Capacity: 10,000 entries.
- Lifetime: 2 minutes.
- Concurrent access is synchronized.

The cache suppresses repeated propagation and routing loops. Unique malicious PacketIDs can still churn the cache and consume processing/bandwidth.

## Forwarding queue

Each active peer has a bounded forwarding queue with both limits enforced:

- 1,000 queued packets
- 2 MiB queued packet bytes

Enqueue is non-blocking. If either limit would be exceeded, the packet is dropped. Failed enqueue accounting is rolled back. When the peer writer dequeues a packet, its queue accounting is released.

## Router I/O rule

The router never performs blocking network I/O while deciding or initiating forwarding. It selects eligible peers and enqueues work; the peer writer owns transport I/O.

## Aggregate bound

With the M5 maximum of 50 active peers, the nominal queued-payload ceiling from the per-peer byte limits is 100 MiB, excluding small fixed bookkeeping and the single packet temporarily held by a writer during transport I/O.

## Security interpretation

These are bounded-availability controls, not a complete DoS defense. The system may still experience packet drops, CPU pressure, or bandwidth exhaustion under sufficiently aggressive unique-packet floods.


## M7 Resource-Hardening Verification

M7 revalidated the existing M6 resource bounds and added router-entry
oversized-input rejection. Concurrent replay and handshake timeout paths were
also reviewed for bounded lifecycle behavior.

M6 cache and forwarding-queue limits remain frozen and unchanged.

**Status: 🟢 M7 VERIFIED**
