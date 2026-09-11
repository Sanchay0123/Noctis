# Security Demonstrations

## Demo 1 --- Normal E2EE

``` text
Alice → Relay 1 → Relay 2 → Bob
```

Show:

1.  Alice enters plaintext.
2.  Alice creates ciphertext.
3.  Relay 1 receives ciphertext.
4.  Relay 1 forwards it.
5.  Relay 2 forwards it.
6.  Bob decrypts it.
7.  Bob sees the original plaintext.

## Demo 2 --- Relay cannot decrypt

At Relay 1:

-   inspect packet
-   show ciphertext
-   show absence of destination session key
-   attempt decryption using relay's own cryptographic state
-   demonstrate failure

## Demo 3 --- Tampering

Modify one byte of ciphertext or authenticated metadata.

Expected:

``` text
AEAD authentication failure
        ↓
message rejected
        ↓
no plaintext delivered
```

## Demo 4 --- Replay

Capture a valid packet and inject it again.

Expected:

``` text
replay detected
    ↓
packet/message rejected
```

## Demo 5 --- Impersonation

Replace the legitimate peer's key material during session establishment.

Expected:

``` text
identity/transcript verification failure
        ↓
session rejected
```

## Demo 6 --- Route failure

Disable a relay.

Expected behavior depends on the final routing implementation. Do not
claim alternate routing unless it is actually implemented and tested.

## Demo 7 --- Observability

Use a four-node topology with Prometheus/Grafana available.

Show:

1. nodes becoming connected
2. a message traversing multiple hops
3. packet receive/forward counters changing
4. no relay decryption occurring
5. a replay attempt
6. replay rejection counter/event increasing
7. a tampering attempt
8. corresponding security failure telemetry

Example:

```text
Alice → Relay 1 → Relay 2 → Bob
             │
             └── sanitized events → Prometheus/Grafana
```

The dashboard must not display the application plaintext.

## Demo 8 --- Observability failure isolation

Stop or disconnect the monitoring backend.

Expected:

```text
Grafana/Prometheus unavailable
        ↓
mesh continues operating
        ↓
encrypted message still delivered
```

This demonstrates that observability is out-of-band and not a
cryptographic or routing dependency.

## Demo 9 --- Telemetry leakage test

Inspect relay logs during a normal encrypted conversation.

Expected:

```text
routing/security telemetry present
        ↓
application plaintext absent
        ↓
private/session keys absent
```

A dashboard screenshot is supporting evidence only. Automated tests and
direct protocol behavior remain authoritative.
