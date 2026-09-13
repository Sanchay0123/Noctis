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

- [[00-governance/04-Architecture-Review-M6]] — authoritative M6 architecture and acceptance record
