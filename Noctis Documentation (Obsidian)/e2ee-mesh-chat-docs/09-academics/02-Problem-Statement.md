# Problem Statement

Traditional introductory messaging demonstrations often show encryption
between two endpoints but do not address what happens when communication
must traverse intermediate network participants.

A decentralized mesh introduces additional challenges:

-   messages may traverse untrusted relay nodes
-   routing metadata may be visible to intermediate nodes
-   packets may be modified, duplicated or replayed
-   malicious nodes may drop or misroute traffic
-   cryptographic identity must be separated from routing identity
-   secure sessions must remain end-to-end despite multi-hop transport

The problem addressed by this project is therefore:

> **How can an educational multi-node messaging system provide
> authenticated end-to-end confidentiality and integrity while allowing
> application ciphertext to traverse untrusted intermediate mesh
> nodes?**

The solution must remain simple enough to implement, test and explain
while using established cryptographic constructions rather than custom
cryptography.
