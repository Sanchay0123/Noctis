# Results Structure

> **Staged results document — final results remain incomplete.**
>
> M2 and M3 provide implementation and test evidence for the long-term
> identity and authenticated session/message-protection layers. Network and
> mesh results remain TBD.

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


## M3 preliminary results

The M3 implementation provides evidence for authenticated session
establishment and application-message protection. These results are
implementation-layer results, not a claim of complete end-to-end network
security.

Mesh routing, direct network delivery, malicious-relay demonstrations,
performance measurements and final deployment evidence remain to be
collected in later milestones.
