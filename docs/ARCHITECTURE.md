# Architecture

This document describes the API Gateway architecture, data flow, and design decisions.

---

## Overview

The gateway is an HTTP service in Go that:

1. Handles `GET /dashboard`.
2. Calls three upstream services (Users, Orders, Billings) **in parallel**.
3. Aggregates results (users with orders, each order with billing).
4. Returns a single JSON response; on partial failure it uses **graceful degradation** (200 with an `errors` map).

---

## Flow diagram

```
                    ┌─────────────────────────────────────────────────────────┐
                    │                     API Gateway (:8080)                 │
                    │                                                          │
  GET /dashboard    │   Router  →  RequestContext middleware  →  Dashboard     │
  ───────────────►  │                    (trace_id, timeout)        Handler   │
                    │                                    │                     │
                    │                                    ▼                     │
                    │                         GetDashboardUseCase              │
                    │                                    │                     │
                    │              (goroutines + channels to call providers)   │
                    │              ┌─────────────────────┼─────────────────────┤
                    │              │                     │                     │
                    │              ▼                     ▼                     ▼
                    │     UsersProvider         OrdersProvider       BillingsProvider
                    │     (circuit breaker)     (circuit breaker)    (circuit breaker)
                    │              │                     │                     │
                    └──────────────┼─────────────────────┼─────────────────────┘
                                   │                     │                     │
                                   ▼                     ▼                     ▼
                            users-mock:8081       orders-mock:8082      billings-mock:8083
```

---

## Layers

| Layer | Package | Responsibility |
|-------|---------|-----------------|
| **Entrypoint** | `cmd/main.go` | Loads config, initializes telemetry, builds container, starts HTTP server, graceful shutdown. |
| **Router** | `internal/router` | Registers routes (`/health`, `/dashboard`, `/metrics`) and applies middleware. |
| **Handler** | `internal/handler` | Handles the request, extracts context (trace_id, timeout), calls use case, writes response or error. |
| **Use case** | `internal/usecase` | Launches goroutines to call the three providers, collects results via channels, builds the dashboard DTO, and handles partial failures. |
| **Providers** | `internal/contract` + `internal/clients` + `internal/resilience` | Interfaces, HTTP implementations, and circuit-breaker wrappers. |
| **Config** | `internal/config` | Reads environment variables (service URLs). |
| **Container** | `internal/container` | Instantiates config, logger, metrics, clients, circuit breakers, use case, handler, and mux. |

---

## Concurrency and context

- Calls to Users, Orders, and Billings run in parallel using **goroutines and channels**: the use case starts one goroutine per provider, each sends its result (or error) on a dedicated channel, and the use case receives from all three channels to build the response.
- The request context is propagated into each call (timeout and cancellation).
- The middleware can set a request deadline; the use case and clients respect the context.

---

## Resilience

- **Circuit breaker** (gobreaker) per provider: after N consecutive failures, the circuit opens and calls fail fast until the recovery timeout.
- **Graceful degradation:** if one or more providers fail, the gateway still returns 200 with available data and fills the `errors` map (e.g. `"orders": "upstream unavailable"`).

---

## Observability

- **Logger:** interface in `internal/observability/logger`, zerolog implementation; injected into handler, use case, and circuit breaker state-change callback.
- **Metrics:** Prometheus in `internal/observability/metrics`; counter and histogram per route; exposed at `/metrics`.
- **Tracing:** OpenTelemetry in `internal/observability/telemetry`; spans in handler and clients; OTLP export (e.g. Tempo); W3C Trace Context propagation.

---

## Dependency injection

The container wires the full dependency tree:

- Config → Clients (Users, Orders, Billings) → Circuit breakers → Use case → Handler.
- Logger, metrics, and tracer are created once and injected where needed.

This makes it easy to test with mocks and swap implementations (e.g. another logger).

---

## Graceful shutdown

In `main`:

1. Listen for SIGINT/SIGTERM.
2. On signal, call `server.Shutdown(ctx)` with a timeout (e.g. 10s).
3. Then call `telemetry.Shutdown(ctx)` to flush the tracer.
4. Exit.

In-flight requests can finish and traces are exported before the process exits.
