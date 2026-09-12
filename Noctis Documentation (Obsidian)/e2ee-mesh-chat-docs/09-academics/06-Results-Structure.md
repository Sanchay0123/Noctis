# Results Structure

> **Staged results document — final results remain incomplete.**
>
> M2, M3 and M4 provide implementation and test evidence for identity,
> authenticated session/message protection and direct secure messaging. Mesh
> routing results remain TBD.

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
  Confidentiality        Transport interception  Verified   M4 E2E evidence
  Integrity              Tampering               Verified   M4 tamper evidence
  Replay resistance      Replay injection        Verified*  M3 session evidence
  Authentication         Invalid identity       Verified   M3/M4 evidence
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
establishment and application-message protection.

## M4 preliminary results

M4 provides evidence that the approved M3 cryptographic/session layer can be
integrated with bounded direct TCP transport and the frozen Protobuf envelope.
The direct path successfully exchanges encrypted application messages in both
directions, rejects malformed/unknown-field packets, enforces session and
identity binding, handles short writes, and rejects ciphertext tampering before
application plaintext delivery.

These are direct secure messaging results, not evidence of complete mesh
security. Multi-hop routing, malicious-relay routing demonstrations, runtime
network hardening, performance measurements and final deployment evidence
remain for later milestones.

*Replay resistance remains a property of the approved M3 session layer; M4
transports the sequence number and ciphertext through the direct packet path.


## M5 preliminary results

M5 demonstrates a functioning direct networking runtime above the approved M4
secure channel. The implementation passed race-detector validation and repeated
100-run race testing for the mesh package. Explicit arbitration tests cover all
four identity/direction permutations.

The containerized runtime starts Alice, Bob, Carol and Dave. A deliberate
simultaneous Dave↔Bob connection collision demonstrates the deterministic rule:
Dave's identity (`4765bd805913e158`) is smaller than Bob's
(`563f793f6ad353dd`), so Dave's outbound connection survives while Bob's
competing outbound connection is rejected as a duplicate. Both nodes retain one
active peer.

These results establish direct-network runtime hardening, not multi-hop mesh
routing. M6 must provide separate routing evidence.
