# Configuration

Reference for environment variables and stack configuration (Docker, Prometheus, Tempo).

---

## Gateway environment variables

| Variable | Required | Description | Example |
|----------|----------|-------------|---------|
| `USERS_URL` | Yes | Base URL of the users service (no trailing slash) | `http://localhost:8081` or `http://users-mock:8081` |
| `ORDERS_URL` | Yes | Base URL of the orders service | `http://localhost:8082` or `http://orders-mock:8082` |
| `BILLINGS_URL` | Yes | Base URL of the billings service | `http://localhost:8083` or `http://billings-mock:8083` |
| `OTLP_ENDPOINT` | No | OTLP endpoint for trace export (host:port, no scheme). If empty, defaults to `localhost:4318`. | `tempo:4318` or `localhost:4318` |

### Loading

- In local development, `main` loads the root `.env` file via `godotenv`.
- In Docker, variables are set in the `environment` section of the `api-gateway` service in `compose.yml`.

---

## Docker Compose

### Services

| Service | Image/Build | Ports (host:container) | Description |
|---------|-------------|------------------------|-------------|
| `api-gateway` | Build from `Dockerfile` | 8080:8080 | Main gateway |
| `users-mock` | Build `Dockerfile.mock` (arg `MOCK=users-mock`) | 8081:8081 | Users service mock |
| `orders-mock` | Build `Dockerfile.mock` (arg `MOCK=orders-mock`) | 8082:8082 | Orders service mock |
| `billings-mock` | Build `Dockerfile.mock` (arg `MOCK=billings-mock`) | 8083:8083 | Billings service mock |
| `prometheus` | `prom/prometheus:latest` | 9090:9090 | Metrics scraping |
| `tempo` | `grafana/tempo:latest` | 4318:4318, 3200:3200 | OTLP trace backend |
| `grafana` | `grafana/grafana:latest` | 3000:3000 | UI (anonymous auth enabled) |

### Network

All services use the default Compose network. The gateway reaches mocks by service name (`users-mock`, `orders-mock`, `billings-mock`).

### Prometheus (scraping the gateway)

When the gateway runs in Docker, Prometheus should scrape by **service name**, not `localhost`:

```yaml
# prometheus.yml (for Docker Compose)
scrape_configs:
  - job_name: 'api-gateway'
    static_configs:
      - targets: ['api-gateway:8080']
    metrics_path: /metrics
    scrape_interval: 30s
```

For local development (Prometheus on host), use `targets: ['localhost:8080']`.

---

## Tempo

The root `tempo.yml` is mounted at `/etc/tempo.yml` in the container. The gateway sends traces via OTLP HTTP to `tempo:4318` in Compose. `OTLP_ENDPOINT` is already set to `tempo:4318` in `compose.yml`.

---

## Grafana

Compose sets:

- `GF_AUTH_ANONYMOUS_ENABLED=true`
- `GF_AUTH_ANONYMOUS_ORG_ROLE=Admin`

So you can open http://localhost:3000 without login. For production, disable anonymous access and configure users/passwords.

To view gateway traces in Grafana, add Tempo as a data source (Tempo URL on Docker network: `http://tempo:3200` or per Tempo docs).

---

## Environment summary

| Environment | USERS_URL | ORDERS_URL | BILLINGS_URL | OTLP_ENDPOINT |
|-------------|------------|------------|--------------|---------------|
| Local (mocks on host) | http://localhost:8081 | http://localhost:8082 | http://localhost:8083 | localhost:4318 (optional) |
| Docker Compose | http://users-mock:8081 | http://orders-mock:8082 | http://billings-mock:8083 | tempo:4318 |
