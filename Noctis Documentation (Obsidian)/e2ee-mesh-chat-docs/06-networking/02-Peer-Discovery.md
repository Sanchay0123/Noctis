# Peer Discovery

## Objective

Allow a node to learn about reachable peers.

## Discovery must define

-   discovery transport
-   peer identifier
-   advertised address/endpoint
-   expiration
-   authentication requirements
-   duplicate handling
-   malformed advertisement handling

## Trust warning

Discovery information is not automatically trustworthy merely because it
came from a network peer.

The system must distinguish:

-   "I learned that this endpoint exists"
-   "I authenticated this endpoint as identity X"

## Initial prototype

Keep discovery simple and local unless a stronger requirement exists.

Avoid adding distributed consensus or blockchain mechanisms merely to
make discovery appear decentralized.
