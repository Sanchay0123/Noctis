# Results Structure

> **Staged results document — final results remain incomplete.**
>
> M2 provides implementation and test evidence for the long-term Ed25519
> identity component. Network, session, AEAD, replay, and mesh results remain TBD.

## Functional results

-   node startup
-   peer discovery
-   connection establishment
-   direct messaging
-   multi-hop messaging
-   reconnection

## Security results

  Property              Test               Result       Evidence
  --------------------- ------------------ -----------  ----------------
  Ed25519 identity       Identity test suite Verified     M2 test evidence
  Confidentiality        Relay inspection    TBD          TBD
  Integrity              Tampering           TBD          TBD
  Replay resistance      Replay injection    TBD          TBD
  Authentication         Invalid identity    TBD          TBD
  MITM resistance        Key substitution    TBD          TBD
  Routing containment    TTL/loop             TBD          TBD

## Performance results

Only measure if performance is a project requirement.

Possible metrics:

-   message latency
-   throughput
-   connection establishment time
-   routing overhead
-   memory usage

Avoid presenting measurements without a reproducible method.
