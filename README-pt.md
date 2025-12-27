# Dev-connect-relations

Projeto Go para gerenciar relações, stacks e recomendações usando Neo4j como banco de grafos.

Faz conexão com outros projetos do dev-connect como:
- [DevConnect (perfis,auth e projetos)](https://github.com/Gustavo-Akira/devconnect)
- [DevConnect frontend(frontend React)](https://github.com/Gustavo-Akira/devconnect-frontend)

**Visão Geral**
- Serviço HTTP REST escrito em Go usando `gin`.
- Banco de grafos: Neo4j (driver `neo4j-go-driver/v5`).
- Mensageria: Kafka (consumidor para eventos de perfil criado).
- Arquitetura: adaptadores inbound/outbound, domínio, aplicação (use-cases) — separação clara de responsabilidades.

**Principais Componentes**
- `cmd/api/main.go`: entrada da aplicação, inicializa infra, rotas e consumidores Kafka.
- `internal/adapters/inbound/rest`: controladores HTTP por recurso (profile, relation, stack, city, recommendation).
- `internal/adapters/outbound/repository`: implementações Neo4j para repositórios de domínio.
- `internal/domain`: entidades e serviços do domínio (profile, city, stack, relation, recommendation).
- `internal/application`: casos de uso compostos (ex.: paginação de relações).

**Decisões Técnicas**
- Go + Gin: performance, simplicidade e tipagem estática.
- Neo4j: modelo natural para relações entre perfis, cidades e stacks; consultas de recomendação baseadas em grafos.
- Kafka: desacoplamento via eventos (ex.: profile.created) para ouvvir criação de perfis assíncrona.
- Estrutura hexagonal (ports/adapters) para facilitar testes e troca de infra.
- Reconexão ao Neo4j: adicionada lógica de retry na inicialização para suportar orquestração com Docker Compose (evita falha quando Neo4j ainda está iniciando).

**Como rodar (local / Docker)**
- Exemplo usando Docker Compose (assume arquivos `docker-compose.yml` / `docker-compose-infra.yml` já configurados):

```bash
# Levanta infra (Neo4j, Kafka etc.)
docker compose -f docker-compose-infra.yml up -d

# Aguarde Neo4j e Kafka criarem seus containers; em seguida, rode a API (localmente):
# (ou construa e rode via Dockerfile se preferir)
go run ./cmd/api
```

- A aplicação agora espera pelo Neo4j até N tentativas antes de falhar. Controle via variáveis de ambiente (opcional):
  - `NEO4J_RETRY_COUNT` (padrão `10`)
  - `NEO4J_RETRY_INTERVAL_SECONDS` (padrão `5`)

**Variáveis de ambiente relevantes**
- `NEO4J_URI` (ex.: `neo4j://neo4j:7687`)
- `NEO4J_USER` (ex.: `neo4j`)
- `NEO4J_PASSWORD`
- `KAFKA_SERVER` (ex.: `localhost:9092`)
- `KAFKA_PROFILE_CREATED_TOPIC` (ex.: `dev-profile.created.v1`)
- `KAFKA_GROUP_ID` (ex.: grupo do consumidor)
- `AUTH_URL` (endpoint do serviço de autenticação usado pelo middleware)
- `NEO4J_RETRY_COUNT`
- `NEO4J_RETRY_INTERVAL_SECONDS`

**Health / Robustez**
- Reconexão de startup: a aplicação tenta criar o driver Neo4j e verificar conectividade repetidas vezes antes de iniciar, reduzindo erros em arranjos Docker onde o Neo4j pode demorar para ficar pronto.
- Recomenda-se adicionar endpoints de healthcheck (liveness/readiness) caso for orquestrar em Kubernetes.

**Testes**
- Testes unitários estão organizados por pacote (`domain`, `adapters`, `application`).
- Para testes de integração com Neo4j, há suporte para container Neo4j nos testes (ver `internal/tests/neo4j_container.go`).


Contribuições e issues são bem-vindas.
