# Test Evidence Standard

## Evidence levels

### Level 0 --- Claim

A developer says a property exists.

**Not sufficient.**

### Level 1 --- Unit test

A focused test exists.

Useful but may not prove system-level behavior.

### Level 2 --- Integration test

Multiple components demonstrate the property.

### Level 3 --- Adversarial test

The test deliberately attempts to violate the property.

### Level 4 --- Demonstration evidence

A reproducible multi-node scenario shows the behavior.

## Approval rule

A security property should generally have Level 2--4 evidence depending
on its scope.

## Evidence record

Every important test result should record:

-   test name
-   objective
-   setup
-   exact command
-   expected result
-   actual result
-   pass/fail
-   relevant logs
-   screenshots where useful
-   commit/version
-   known limitations

## No silent claims

If a test does not exist, documentation must say:

> **Not yet verified.**

rather than implying security.


## M6 Reference

- [[00-governance/04-Architecture-Review-M6]] — authoritative M6 architecture and acceptance record
