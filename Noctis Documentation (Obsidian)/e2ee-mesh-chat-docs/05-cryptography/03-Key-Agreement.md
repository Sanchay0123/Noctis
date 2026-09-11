# Key Agreement

## X25519

X25519 provides shared secret agreement between endpoints.

It does not, by itself, authenticate who owns a public key.

Therefore:

> X25519 must be authenticated by the Ed25519 identity mechanism.

## Conceptual exchange

``` text
Alice:
  Ed25519 identity A
  X25519 ephemeral a

Bob:
  Ed25519 identity B
  X25519 ephemeral b

Alice computes X25519(a, B_ephemeral)
Bob computes X25519(b, A_ephemeral)

Both derive the same shared secret.

Ed25519 signatures bind the ephemeral exchange to identities.
```

## Required properties

-   both parties derive the same shared secret
-   unauthorized key substitution fails
-   session context is bound
-   roles/directions are separated
-   invalid authentication terminates the exchange

## Required security test

MITM identity substitution:

1.  attacker intercepts handshake
2.  attacker substitutes X25519 public key
3.  legitimate Ed25519 authentication no longer matches
4.  handshake is rejected
