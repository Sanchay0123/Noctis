# Deployment Architecture

## Reference demonstration

```
flowchart LR
    A[Alice] --- N1[Node 1]
    N1 --- N2[Node 2]
    N2 --- B[Bob]

    N1 -. telemetry .-> OBS[Observability]
    N2 -. telemetry .-> OBS
    A -. telemetry .-> OBS
    B -. telemetry .-> OBS

    OBS --> PROM[Prometheus]
    PROM --> GRAF[Grafana]
```

## Recommended academic test topology

Four independently identifiable node processes:

- Alice endpoint
    
- Node 1 relay
    
- Node 2 relay
    
- Bob endpoint
    

They should be reproducible as containers for the standard development  
and demonstration environment.

The project should not require manual host installation of application  
dependencies.

## Containerized observability

The demonstration environment may additionally contain:

- Prometheus
    
- Grafana
    
- Loki or another documented structured-log backend
    

The exact log backend is an implementation decision.

Observability services must remain outside the E2EE message path.

## Observation points

At Node 1 and Node 2, the demonstration may show:

- packet received
    
- source/destination routing identifiers where permitted
    
- TTL
    
- packet identifier
    
- ciphertext presence
    
- forwarding decision
    
- security-safe counters/events
    

It must **not** reveal:

- plaintext
    
- private keys
    
- session keys
    
- passwords
    
- decrypted application content
    

## Observability failure

If Grafana, Prometheus or log aggregation is stopped, the mesh must  
continue to provide its core communication behavior.

## Failure scenario

At minimum, one relay should be disabled during a demonstration to show  
the behavior of the routing layer.

Whether automatic alternate routing is supported is an implementation  
decision and must not be claimed until tested.

See [[02-architecture/07-Observability-and-Security-Telemetry]] and  
[[07-testing/05-Security-Demonstrations]].