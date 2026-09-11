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
