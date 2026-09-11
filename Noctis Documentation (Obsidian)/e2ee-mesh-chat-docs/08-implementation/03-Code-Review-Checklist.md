# Code Review Checklist

## 1. Correctness

-   [ ] Does it work?
-   [ ] Are edge cases handled?
-   [ ] Are failure paths tested?

## 2. Security

-   [ ] Can an attacker decrypt messages?
-   [ ] Can an attacker impersonate a peer?
-   [ ] Can packets be modified undetected?
-   [ ] Can packets be replayed?
-   [ ] Can a nonce be reused?
-   [ ] Can keys leak?
-   [ ] Can malformed packets crash the node?

## 3. Cryptography

-   [ ] Correct primitive
-   [ ] Correct key sizes
-   [ ] Correct nonce handling
-   [ ] Secure randomness
-   [ ] Correct KDF
-   [ ] Correct authentication
-   [ ] Correct key lifecycle
-   [ ] Safe error handling

## 4. Networking

-   [ ] Routing loops
-   [ ] TTL
-   [ ] duplicate packets
-   [ ] malformed input
-   [ ] resource exhaustion
-   [ ] timeouts
-   [ ] race conditions

## 5. Privacy

-   [ ] no plaintext relay logs
-   [ ] no private keys in logs
-   [ ] no session keys in logs
-   [ ] metadata exposure documented

## 6. Maintainability

-   [ ] separation of concerns
-   [ ] minimal dependencies
-   [ ] no unnecessary abstraction
-   [ ] testability
-   [ ] no duplicated security logic

## 7. Evidence

-   [ ] tests actually exercise the security property
-   [ ] negative tests exist
-   [ ] integration behavior is demonstrated
