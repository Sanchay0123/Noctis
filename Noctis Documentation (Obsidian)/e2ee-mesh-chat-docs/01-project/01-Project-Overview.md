# Project Overview

## Working title

**Design and Implementation of an End-to-End Encrypted Messaging System
over a Decentralized Mesh Network**

## One-sentence definition

An educational multi-node messaging prototype in which application
messages remain end-to-end encrypted while being transported through
intermediate mesh nodes.

## Core demonstration

``` text
Alice
  |
  | encrypted application payload
  v
Node B
  |
  | forwards ciphertext
  v
Node C
  |
  | forwards ciphertext
  v
Bob
```

The critical demonstration is that **Node B and Node C can forward the
packet but cannot decrypt Alice's application message**.

## Why this project exists

The project demonstrates how several security and distributed-systems
concepts interact:

1.  cryptographic identity
2.  authenticated key establishment
3.  authenticated encryption
4.  replay resistance
5.  peer-to-peer networking
6.  multi-hop forwarding
7.  failure handling
8.  adversarial testing

## Academic character

This is a **research/educational prototype**, not a production messaging
platform.

The project must be technically defensible without overstating its
security properties.

## Primary success criterion

A reproducible demonstration must establish:

-   Alice can communicate with Bob.
-   The route can contain intermediate nodes.
-   Intermediate nodes see routing information and ciphertext.
-   Intermediate nodes cannot decrypt the application content.
-   Tampering causes authentication failure.
-   Replaying an accepted message does not cause it to be accepted
    again.
-   Basic impersonation attempts fail.
-   Route failure is handled to the extent implemented.

See [02-Requirements](02-Requirements.md) and
[01-Threat-Model](../03-security/01-Threat-Model.md).
