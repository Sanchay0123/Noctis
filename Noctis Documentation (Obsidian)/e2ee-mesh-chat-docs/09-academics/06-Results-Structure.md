# Results Structure

> **Placeholder until implementation exists.**

## Functional results

-   node startup
-   peer discovery
-   connection establishment
-   direct messaging
-   multi-hop messaging
-   reconnection

## Security results

  Property              Test               Result   Evidence
  --------------------- ------------------ -------- ----------
  Confidentiality       Relay inspection   TBD      TBD
  Integrity             Tampering          TBD      TBD
  Replay resistance     Replay injection   TBD      TBD
  Authentication        Invalid identity   TBD      TBD
  MITM resistance       Key substitution   TBD      TBD
  Routing containment   TTL/loop           TBD      TBD

## Performance results

Only measure if performance is a project requirement.

Possible metrics:

-   message latency
-   throughput
-   connection establishment time
-   routing overhead
-   memory usage

Avoid presenting measurements without a reproducible method.
