# M6 Handshake Correlation

## Purpose

Multi-hop forwarding separates the transport peer carrying a packet from the endpoint E2EE session represented by that packet. Because multiple initiator handshakes to the same remote can be in flight simultaneously, a RESP must be correlated to the exact INIT candidate rather than merely to the remote identity.

## Pending-state key

Conceptually:

```text
pendingInits[remoteIdentity][SHA256(T_INIT)] -> InitiatorHandshake
```

The hash is internal state only; it is not added to the wire protocol.

## Response processing

1. Obtain the responder identity from the packet metadata.
2. Select only pending candidates for that remote identity.
3. Test every candidate against the exact frozen M3 `T_RESP` reconstruction.
4. Require exactly one match.
5. On one match, consume the candidate before session registration.
6. On zero or multiple matches, fail closed.

This prevents late, duplicate, wrong-identity, wrong-ephemeral-key, and ambiguous RESP messages from being attached to the wrong session.

## Security invariant

Routing metadata never authenticates an endpoint. Only successful M3 transcript verification and key establishment can establish an endpoint session.
