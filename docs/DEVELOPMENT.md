# Development

Guide to setting up the development environment, running tests, and building the project.

---

## Prerequisites

- **Go 1.25** or newer
- **Git**
- Optional: **Docker** and **Docker Compose** to run the full stack

---

## Environment setup

1. Clone the repository (or use the project directory).

2. Install dependencies:

   ```bash
   go mod download
   ```

3. Create a `.env` file at the project root (or export the variables):

   ```env
   USERS_URL=http://localhost:8081
   ORDERS_URL=http://localhost:8082
   BILLINGS_URL=http://localhost:8083
   ```

   For local tracing with Tempo:

   ```env
   OTLP_ENDPOINT=localhost:4318
   ```

---

## Running the gateway locally

1. Start the mocks (in separate terminals):

   ```bash
   go run cmd/users-mock/main.go
   go run cmd/orders-mock/main.go
   go run cmd/billings-mock/main.go
   ```

2. Start the gateway:

   ```bash
   go run cmd/main.go
   ```

3. Test:

   ```bash
   curl http://localhost:8080/dashboard
   curl http://localhost:8080/health
   curl http://localhost:8080/metrics
   ```

---

## Tests

Run all tests:

```bash
go test ./...
```

With coverage:

```bash
go test -cover ./...
```

Use case package only:

```bash
go test ./tests/...
```

Use case tests specifically:

```bash
go test ./tests/usecase/...
```

The use case tests cover aggregation logic and graceful degradation (one or more provider failures still return 200 with `errors`).

---

## Build

Gateway binary:

```bash
go build -o bin/gateway ./cmd/main.go
./bin/gateway
```

Mock binaries:

```bash
go build -o bin/users-mock ./cmd/users-mock/main.go
go build -o bin/orders-mock ./cmd/orders-mock/main.go
go build -o bin/billings-mock ./cmd/billings-mock/main.go
```

---

## Docker

### Gateway only

```bash
docker build -t api-gateway .
docker run -p 8080:8080 \
  -e USERS_URL=http://host.docker.internal:8081 \
  -e ORDERS_URL=http://host.docker.internal:8082 \
  -e BILLINGS_URL=http://host.docker.internal:8083 \
  api-gateway
```

(With mocks running on the host.)

### Full stack (gateway + mocks + Prometheus + Tempo + Grafana)

```bash
docker compose up --build
```

See [CONFIGURATION.md](CONFIGURATION.md) for variables and Prometheus target (`api-gateway:8080` when everything runs in Compose).

---

## Test structure

- **`tests/usecase/`** – dashboard use case tests (with provider mocks).
- **`tests/utils/`** – helper tests (e.g. `error_helper`).

The use case tests use a mock logger to keep output clean and to assert circuit breaker / graceful degradation behavior.

---

## Linting and formatting

Standard Go formatting:

```bash
go fmt ./...
```

For extra tools (e.g. `golangci-lint`), add them to CI or a Makefile.
