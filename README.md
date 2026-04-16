# go-project

## Swagger

- UI: `GET /swagger`
- OpenAPI spec: `GET /swagger/openapi.yaml`

## Run

1. Apply migrations:
   - `go run ./cmd/migrator -command up`
2. Start app:
   - `go run ./cmd/app`

## Docker Compose

1. Start infra (`postgres`, `kafka`, `zookeeper`, `kafka-ui`):
   - `docker compose up -d postgres zookeeper kafka kafka-ui`
2. Run migrations in docker network:
   - `docker compose --profile tools run --rm migrator`
3. Start app in docker:
   - `docker compose --profile app up -d app`

Useful:
- API: `http://localhost:8080`
- Swagger: `http://localhost:8080/swagger`
- Kafka UI: `http://localhost:8081`
