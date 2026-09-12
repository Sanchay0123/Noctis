# Protocol State Machines

## Node lifecycle

```mermaid
stateDiagram-v2
    [*] --> Starting
    Starting --> Discovering
    Discovering --> Connecting
    Connecting --> Authenticating
    Authenticating --> Established
    Established --> Reconnecting
    Reconnecting --> Connecting
    Established --> Shutdown
    Starting --> Shutdown
    Shutdown --> [*]
```

## Session lifecycle

```mermaid
stateDiagram-v2
    [*] --> None
    None --> Handshaking
    Handshaking --> Established
    Handshaking --> Failed
    Established --> Rekeying
    Rekeying --> Established
    Established --> Closed
    Failed --> [*]
    Closed --> [*]
```

## Invalid transitions

Invalid state transitions must:

-   be rejected
-   not expose secrets
-   not corrupt state
-   produce bounded diagnostic information


## M3 session implementation boundary

The session state machine is implemented for the M3 handshake/session layer.
The implementation distinguishes:
- handshake state transitions
- established-session operation
- rekey/rotation boundary
- closed/failed states

Invalid transitions fail closed and do not expose cryptographic secrets.

The handshake state machine's duplicate/ordering checks are not a substitute
for application-message replay protection. Application replay is enforced by
the 64-message per-session receive window described in
[[04-protocol/05-Replay-Protection]].
