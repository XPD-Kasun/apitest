# Hexagonal Architecture Suitability Test

Hexagonal architecture is by far one of the best software architecture to model any component at C4 level 2. This architecture uses ports and adapters to decouple the business logic from its environment. Ports and adapters used here is somewhat similar to standard adapter pattern where the port is the interface which the adapters implement or use. We have two types of ports and adapters : driving and driven. The flow into the hexagonal(app core) is handled by the driving ports and adapters.The outflow is handled by driven ports and adapters.

