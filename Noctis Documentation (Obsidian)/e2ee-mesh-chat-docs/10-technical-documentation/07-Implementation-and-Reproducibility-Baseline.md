# Implementation and Reproducibility Baseline

## Purpose

This page records the implementation baseline used to describe the system in
M10. It separates reproducible repository checks from environment-dependent
evidence.

## Repository

Authoritative implementation directory:

```text
~/Documents/Noctis/CrytpProject
```

The project resides inside the larger Git working tree:

```text
/home/sanchayjain/Documents/Noctis
```

The larger Git root does not change the authoritative project directory.

## Standard validation commands

The project uses the following baseline checks:

```bash
go test ./...
go test -race -p 1 ./...
go vet ./...
go build ./...
go build -tags gui -o meshchat-gui ./cmd/meshchat-gui
```

The M9 evidence package records successful execution of the relevant validation
commands at the M9 gate.

## Containerization

The repository contains a multi-stage Alpine Go Dockerfile and a Docker Compose
configuration capable of defining a four-node topology.

Docker runtime regression is **not currently verified** in the final evidence
because external proxy/DNS conditions prevented a clean verification run.
This limitation must not be rewritten as a successful deployment test.

## Observability

Current implementation:

- telemetry interface/hooks
- bounded non-blocking no-op telemetry recorder

Not implemented:

- active Prometheus exporter
- metric-vector implementation
- Grafana deployment

Telemetry is not part of the E2EE critical path.

## GUI reproducibility

The GUI requires native graphics support. The headless verification environment
could compile the GUI but could not provide the same interactive desktop
runtime. M9.3 therefore relies on explicit manual desktop verification for the
visual/runtime scenarios.

## Evidence provenance

The project distinguishes:

1. current command output
2. repository tests
3. approved milestone evidence
4. manual desktop verification
5. environment-limited/unverified activities

A result without an available raw artifact is not represented as though that
artifact exists.

## Security reproducibility

The most important security demonstrations are encoded as tests and integration
flows rather than screenshots alone. M9.2 includes eight adversarial cases;
M9.1 demonstrates endpoint E2EE across a relay.

## Related records

- [06-Test-Evidence-Standard](../07-testing/06-Test-Evidence-Standard.md)
- [19-M9-Evidence-Manifest](../07-testing/19-M9-Evidence-Manifest.md)
- [20-M9-Final-Report](../07-testing/20-M9-Final-Report.md)
- [05-Security-Limitations](../03-security/05-Security-Limitations.md)
