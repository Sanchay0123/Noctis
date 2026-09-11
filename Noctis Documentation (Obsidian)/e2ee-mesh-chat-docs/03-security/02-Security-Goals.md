# Security Goals

  ----------------------------------------------------------------------------
  ID                      Goal                    Required evidence
  ----------------------- ----------------------- ----------------------------
  SG-01                   Message confidentiality Relay inspection test

  SG-02                   Message integrity       Ciphertext tampering test

  SG-03                   Identity authentication Handshake/authentication
                                                  tests

  SG-04                   Replay resistance       Replay regression tests

  SG-05                   Basic impersonation     Identity-substitution/MITM
                          resistance              test

  SG-06                   Nonce uniqueness        Design invariant + tests
                                                  where practical

  SG-07                   Secret non-disclosure   Logging/storage review

  SG-08                   Malformed input         Fuzz/property/negative tests
                          resilience              

  SG-09                   Routing containment     TTL/loop/duplicate tests
  ----------------------------------------------------------------------------

## Claim discipline

A security claim is only **Verified** when there is evidence.

A passing happy-path test does not prove a security property by itself.

For example:

> "Messages are encrypted."

is weaker than:

> "An integration test captures the application packet at an
> intermediate node and demonstrates that the relay does not possess the
> destination session key and cannot successfully authenticate/decrypt
> the ciphertext."

See [[07-testing/06-Test-Evidence-Standard]].
