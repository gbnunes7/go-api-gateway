# API Reference

HTTP endpoint reference for the API Gateway.

---

## Base URL

- **Local:** `http://localhost:8080`
- **Docker:** `http://localhost:8080` (mapped port of `api-gateway` service)

---

## Endpoints

### GET /dashboard

Returns aggregated users, orders, and billings. Upstream calls (Users, Orders, Billings) are made in parallel. On partial failure the response is still **200 OK**, with available data and an `errors` map indicating which provider(s) failed.

#### Request

- **Method:** `GET`
- **Path:** `/dashboard`
- **Headers:** Optional; tracing headers (e.g. `traceparent`) are propagated to upstreams.

#### Response (full success)

- **Status:** `200 OK`
- **Content-Type:** `application/json`

**Body:** `DashboardResponse`

```json
{
  "users": [
    {
      "id": "string",
      "name": "string",
      "email": "string",
      "orders": [
        {
          "id": "string",
          "totalPrice": number,
          "createdAt": "string (ISO 8601)",
          "billing": {
            "id": "string",
            "paymentType": "string",
            "paidValue": number,
            "orderId": "string"
          }
        }
      ]
    }
  ]
}
```

The `errors` field is omitted when there are no failures.

#### Response (graceful degradation)

- **Status:** `200 OK`
- **Content-Type:** `application/json`

**Body:** Same shape; `users` contains available data and `errors` indicates per-provider failures:

```json
{
  "users": [...],
  "errors": {
    "users": "upstream unavailable",
    "orders": "upstream unavailable",
    "billings": "upstream unavailable"
  }
}
```

Possible keys in `errors`: `"users"`, `"orders"`, `"billings"`.

#### Response (internal / total failure)

- **Status:** `500 Internal Server Error` (or other in exceptional cases)
- **Body:** Handler-dependent (e.g. JSON error message).

#### Example (curl)

```bash
curl -s http://localhost:8080/dashboard | jq
```

---

### GET /health

Simple health check for load balancers and orchestrators.

#### Request

- **Method:** `GET`
- **Path:** `/health`

#### Response

- **Status:** `200 OK`
- **Content-Type:** `application/json`

```json
{
  "status": "Server is healthy"
}
```

---

### GET /metrics

Prometheus-format metrics. Should not be exposed publicly in production.

#### Request

- **Method:** `GET`
- **Path:** `/metrics`

#### Response

- **Status:** `200 OK`
- **Content-Type:** `text/plain; charset=utf-8` (Prometheus format)

Example metrics:

- `http_requests_total` – counter by method, route, and status.
- `http_request_duration_seconds` – duration histogram by route.

---

## Upstream service contracts (mocks)

For reference, the mocks expose the following endpoints and shapes.

### Users (GET /users)

- **Base URL:** from `USERS_URL` (e.g. `http://users-mock:8081`)
- **Response:** array of `User` with `id`, `name`, `email`, `order_id`.

### Orders (GET /orders)

- **Base URL:** `ORDERS_URL` (e.g. `http://orders-mock:8082`)
- **Response:** array of `Order` with `id`, `total_price`, `created_at`, `billing_id`, `user_id`.

### Billings (GET /billings)

- **Base URL:** `BILLINGS_URL` (e.g. `http://billings-mock:8083`)
- **Response:** array of `Billing` with `id`, `payment_type`, `paid_value`, `order_id`.

The gateway joins user → orders → billing and returns the aggregated shape described under **GET /dashboard**.
