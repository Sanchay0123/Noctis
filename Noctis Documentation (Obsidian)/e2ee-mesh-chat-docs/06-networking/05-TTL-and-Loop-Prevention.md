# TTL and Loop Prevention

## TTL

Every routable packet should have a hop limit.

Each forwarding hop consumes one unit according to the final protocol
semantics.

## Loop example

``` text
Alice → B → C → B → C → ...
```

Without TTL, packets may circulate indefinitely.

## Required controls

-   TTL/hop limit
-   packet identifier
-   duplicate cache
-   bounded route state

## Tests

-   TTL=0 rejected
-   TTL expires after expected number of hops
-   repeated packet identifier is not forwarded indefinitely
-   routing loop eventually terminates

## M6 Frozen TTL and Loop Controls

- Default initial TTL: **16**.
- Maximum accepted TTL: **32**.
- TTL 0: drop.
- TTL 1 at destination: deliver.
- TTL 1 at non-destination: drop.
- TTL >1 at destination: deliver.
- TTL >1 at non-destination: decrement once and forward.

Loop prevention combines TTL exhaustion, exclusion of the incoming peer, and the 10,000-entry/2-minute duplicate cache keyed by `(source_node, packet_id)`.

A cyclic topology therefore terminates without requiring route discovery or routing tables.
