# Limitations

The final limitations section should include only limitations supported
by the final implementation.

Expected areas:

-   metadata visibility
-   endpoint compromise
-   local key storage
-   malicious relay dropping
-   route manipulation
-   scalability
-   availability
-   forward secrecy scope
-   deniability
-   deployment assumptions

Do not hide weaknesses. A clearly stated limitation strengthens the
academic security analysis.

## M6 Limitations

The routing layer is intentionally simple. Managed flooding does not discover optimal paths and can generate substantial redundant traffic as topology size grows. Unique-packet floods can still consume CPU and bandwidth, shared queues have no QoS, delivery is best-effort, and routing metadata remains observable. These are deliberate prototype boundaries rather than claims of production-grade anonymity or availability.


## M6 Reference

- [04-Architecture-Review-M6](../00-governance/04-Architecture-Review-M6.md) — authoritative M6 architecture and acceptance record


## M8 Limitations

The GUI is currently an in-memory presentation layer; persistent message storage is not part of M8. Automatic peer discovery and DHT remain outside scope. Live GUI runtime acceptance passed for the two-node manual chat workflow. Docker Compose regression remains unverified because external proxy/DNS restrictions prevented dependency resolution. Timestamp and message-order behavior are implemented but were not directly covered by dedicated M8.3 tests.

M8 does not change the existing limitations around metadata visibility, malicious relay dropping, endpoint compromise, route manipulation, availability or anonymity.
