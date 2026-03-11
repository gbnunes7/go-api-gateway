# API Gateway

A Go API gateway that aggregates data from multiple services (Users, Orders, Billings) into a single endpoint, with concurrency, resilience, and observability.

---

## Table of contents

- [Overview](#overview)
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [API](#api)
- [Architecture](#architecture)
- [Observability](#observability)
- [Development](#development)
- [Project structure](#project-structure)
- [Further documentation](#further-documentation)

---

## Overview

The gateway exposes a single **`/dashboard`** endpoint that calls three upstream services (Users, Orders, Billings) in parallel, merges the data, and returns a unified JSON response. On partial failure it uses **graceful degradation**: HTTP 200 with available data and a per-provider `errors` map.

```
Client  →  GET /dashboard  →  Gateway (:8080)
                                    │
                                    ├→ Users   (:8081)
                                    ├→ Orders  (:8082)
                                    └→ Billings (:8083)
                                    │
Client  ←  Aggregated JSON  ←  Gateway
```

---

## Features

| Area | Description |
|------|-------------|
| **Concurrency** | Calls to Users, Orders, and Billings in parallel using goroutines and channels |
| **Context** | Request context propagation (timeout, cancellation, trace ID) |
| **Resilience** | Circuit breaker (gobreaker) per provider; graceful degradation |
| **Logging** | Structured logs (zerolog) with dependency injection |
| **Metrics** | Prometheus: `http_requests_total`, `http_request_duration_seconds` |
| **Tracing** | OpenTelemetry → OTLP (e.g. Grafana Tempo), W3C propagation |
| **Graceful shutdown** | SIGINT/SIGTERM with HTTP server and tracer shutdown |

---

## Prerequisites

- **Go 1.25+** (for local development)
- **Docker and Docker Compose** (to run the full stack)

---

## Quick Start

### With Docker Compose (recommended)

Runs the gateway, mocks (users, orders, billings), Prometheus, Tempo, and Grafana:

```bash
docker compose up --build
```

- **Gateway:** http://localhost:8080  
- **Dashboard:** http://localhost:8080/dashboard  
- **Metrics:** http://localhost:8080/metrics  
- **Prometheus:** http://localhost:9090  
- **Grafana:** http://localhost:3000 (anonymous auth enabled)  
- **Tempo:** OTLP at `localhost:4318`

### Local (mocks on host)

1. Set up `.env` (or export the variables):

   ```env
   USERS_URL=http://localhost:8081
   ORDERS_URL=http://localhost:8082
   BILLINGS_URL=http://localhost:8083
   ```

2. Start the mocks (in separate terminals):

   ```bash
   go run cmd/users-mock/main.go    # :8081
   go run cmd/orders-mock/main.go   # :8082
   go run cmd/billings-mock/main.go # :8083
   ```

3. Start the gateway:

   ```bash
   go run cmd/main.go
   ```

4. Open: http://localhost:8080/dashboard

---

## Configuration

| Variable | Description | Example (Docker) | Example (local) |
|----------|-------------|------------------|-----------------|
| `USERS_URL` | Base URL of the users service | `http://users-mock:8081` | `http://localhost:8081` |
| `ORDERS_URL` | Base URL of the orders service | `http://orders-mock:8082` | `http://localhost:8082` |
| `BILLINGS_URL` | Base URL of the billings service | `http://billings-mock:8083` | `http://localhost:8083` |
| `OTLP_ENDPOINT` | OTLP endpoint for tracing (host:port) | `tempo:4318` | `localhost:4318` |

The root `.env` file is loaded in development (via `godotenv`). In Docker, set variables in `environment` in `compose.yml`.

See [docs/CONFIGURATION.md](docs/CONFIGURATION.md) for details.

---

## API

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/dashboard` | Aggregated data (users + orders + billings). Graceful degradation on partial failures. |
| `GET` | `/health` | Health check. Returns `{"status":"Server is healthy"}`. |
| `GET` | `/metrics` | Prometheus metrics. |

### Example `/dashboard` response (success)

```json
{
  "users": [
    {
      "id": "user-1",
      "name": "John Doe",
      "email": "john@example.com",
      "orders": [
        {
          "id": "order-1",
          "totalPrice": 10000,
          "createdAt": "2025-03-08T10:00:00Z",
          "billing": {
            "id": "billing-1",
            "paymentType": "credit",
            "paidValue": 10000,
            "orderId": "order-1"
          }
        }
      ]
    }
  ]
}
```

### Example with graceful degradation (one provider failed)

```json
{
  "users": [...],
  "errors": {
    "orders": "upstream unavailable"
  }
}
```

Full reference: [docs/API.md](docs/API.md).

---

## Architecture

- **Handler** → **Use case** → **Providers** (Users, Orders, Billings) wrapped with a circuit breaker.
- Configuration via `config`; dependency injection in `container`.
- Request middleware (trace ID, timeout) and shared response helpers (`WriteError`).

Diagrams and design notes: [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md).

---

## Observability

- **Logs:** zerolog; logger injected into handler, use case, and circuit breaker.
- **Metrics:** request count and duration histogram; exposed at `/metrics`.
- **Tracing:** spans in handler and HTTP clients; OTLP export to Tempo (or another backend).

When using Docker, Prometheus must scrape the gateway; on a shared network use `api-gateway:8080` in `prometheus.yml`. See [docs/CONFIGURATION.md](docs/CONFIGURATION.md).

---

## Development

```bash
# Tests
go test ./...

# Build
go build -o bin/gateway ./cmd/main.go
```

Environment and testing guide: [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md).

---

## Project structure

```
.
├── cmd/
│   ├── main.go              # Gateway entrypoint
│   ├── users-mock/          # Users service mock
│   ├── orders-mock/        # Orders service mock
│   └── billings-mock/      # Billings service mock
├── internal/
│   ├── config/              # Config (env)
│   ├── container/          # DI / wiring
│   ├── contract/           # Provider interfaces
│   ├── dto/                # Request/response DTOs
│   ├── handler/            # HTTP handlers and middleware
│   ├── router/             # Routes and bindings
│   ├── usecase/            # Dashboard aggregation logic
│   ├── clients/            # HTTP clients (users, orders, billings)
│   ├── resilience/         # Circuit breaker
│   ├── observability/      # logger, metrics, telemetry
│   └── utils/              # Helpers (trace_id, error_helper)
├── tests/                  # Use case and helper tests
├── docs/                   # Additional documentation
├── compose.yml             # Stack: gateway + mocks + Prometheus + Tempo + Grafana
├── Dockerfile              # Gateway image
├── Dockerfile.mock         # Mocks image (build arg MOCK)
├── prometheus.yml
└── tempo.yml
```

---

## Further documentation

| Document | Content |
|----------|---------|
| [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) | Architecture, data flow, and design decisions |
| [docs/CONFIGURATION.md](docs/CONFIGURATION.md) | Environment variables and stack configuration |
| [docs/API.md](docs/API.md) | API contracts and examples |
| [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md) | Local setup, tests, and Docker |

---

