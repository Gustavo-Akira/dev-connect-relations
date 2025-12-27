# Dev-connect-relations

Go project to manage relationships, stacks, and recommendations using Neo4j as a graph database.

It integrates with other DevConnect projects:

- [DevConnect (profile ,auth and projects)](https://github.com/Gustavo-Akira/devconnect)
- [DevConnect frontend(React frontend)](https://github.com/Gustavo-Akira/devconnect-frontend)

**Overview**

HTTP REST service written in Go using gin.

Graph database: Neo4j (driver neo4j-go-driver/v5).

Messaging: Kafka (consumer for profile creation events).

Architecture: inbound/outbound adapters, domain, application (use cases) — clear separation of responsibilities.

**Main Components**

- cmd/api/main.go: application entry point, initializes infrastructure, routes, and Kafka consumers.

- internal/adapters/inbound/rest: HTTP controllers by resource (profile, relation, stack, city, recommendation).

- internal/adapters/outbound/repository: Neo4j implementations of domain repositories.

- internal/domain: domain entities and services (profile, city, stack, relation, recommendation).

- internal/application: composed use cases (e.g., relation pagination).

**Technical Decisions**

- Go + Gin: performance, simplicity, and static typing.

- Neo4j: natural model for relationships between profiles, cities, and stacks; graph-based recommendation queries.

- Kafka: event-driven decoupling (e.g., profile.created) to asynchronously react to profile creation.

- Hexagonal architecture (ports/adapters): improves testability and infrastructure replaceability.

- Neo4j reconnection logic: retry strategy during startup to support Docker Compose orchestration (prevents failures when Neo4j is still starting).

**How to Run (Local / Docker)**

- Example using Docker Compose (assumes docker-compose.yml / docker-compose-infra.yml are already configured):

- Start infrastructure (Kafka, etc.)
docker compose -f docker-compose-infra.yml up -d

- Start a Neo4j container and wait for Kafka, Zookeper and Neo4j containers to be ready, then run the API locally:
(or build and run using the Dockerfile if preferred)
go run ./cmd/api


- The application waits for Neo4j for N attempts before failing. This can be controlled via environment variables (optional):

    - NEO4J_RETRY_COUNT (default: 10)

    - NEO4J_RETRY_INTERVAL_SECONDS (default: 5)

**Relevant Environment Variables**

- NEO4J_URI (e.g., neo4j://neo4j:7687)

- NEO4J_USER (e.g., neo4j)

- NEO4J_PASSWORD

- KAFKA_SERVER (e.g., localhost:9092)

- KAFKA_PROFILE_CREATED_TOPIC (e.g., dev-profile.created.v1)

- KAFKA_GROUP_ID (consumer group)

- AUTH_URL (authentication service endpoint used by the middleware)

- NEO4J_RETRY_COUNT

- NEO4J_RETRY_INTERVAL_SECONDS

**Health / Robustness**

- Startup reconnection: the application repeatedly attempts to create the Neo4j driver and verify connectivity before starting, reducing errors in Docker-based environments where Neo4j may take time to become available.

- It is recommended to add healthcheck endpoints (liveness/readiness) when orchestrating with Kubernetes.

**Tests**

- Unit tests are organized by package (domain, adapters, application).

- For Neo4j integration tests, there is support for running a Neo4j container during tests (see internal/tests/neo4j_container.go).

Contributions and issues are welcome.