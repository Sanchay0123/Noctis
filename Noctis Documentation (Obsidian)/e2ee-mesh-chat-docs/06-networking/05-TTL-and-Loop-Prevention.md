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
