# Cinema Booking System

[![CI](https://github.com/mehdi1443/cinema-booking-system/actions/workflows/ci.yml/badge.svg)](https://github.com/mehdi1443/cinema-booking-system/actions/workflows/ci.yml)
[![Go](https://img.shields.io/badge/Go-1.25.8-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A production-oriented cinema booking backend built with Go, Gin, PostgreSQL, Redis, RabbitMQ,
Elasticsearch, and MinIO. The application covers cinema catalog management, seat inventory,
temporary reservations, mock payments and refunds, QR/PDF tickets, email notifications, sales
reporting, settlements, audit logs, and role-based access control.

The codebase follows Clean Architecture: business rules remain independent from HTTP, database,
and third-party implementations, while infrastructure is connected at the application entry point.

> **Project status:** feature-complete backend MVP. Payment and login SMS providers are intentionally
> mocked. Production deployments must connect real providers and use production secrets.

## Table of contents

- [Highlights](#highlights)
- [Architecture](#architecture)
- [Technology stack](#technology-stack)
- [Project structure](#project-structure)
- [Quick start with Docker](#quick-start-with-docker)
- [Local development](#local-development)
- [Configuration](#configuration)
- [Authentication and authorization](#authentication-and-authorization)
- [Booking lifecycle](#booking-lifecycle)
- [Payments, tickets, and outbox](#payments-tickets-and-outbox)
- [API overview](#api-overview)
- [Response contract](#response-contract)
- [Database migrations](#database-migrations)
- [Observability](#observability)
- [Testing](#testing)
- [Security notes](#security-notes)
- [License](#license)

## Highlights

- Phone OTP login with automatic registration and atomic one-time OTP consumption
- Authenticated email verification using OTP
- RSA-signed access and refresh JWTs with Redis-backed revocation
- Centralized RBAC for `admin`, `support`, `user`, and `cinema_owner`
- Ownership checks for cinema, hall, seat, showtime, report, and settlement resources
- Movie, person, genre, cinema, facility, city, and province management
- Favorites, ratings, nested comments, votes, and moderation
- Elasticsearch search and suggestions with RabbitMQ-based synchronization
- Hall, seat, showtime, per-seat pricing, and schedule conflict management
- Concurrent-safe temporary seat reservation with automatic expiry
- Idempotent mock payment, verification, and refund workflows
- QR admission tickets and private PDF tickets stored in MinIO
- Transactional outbox with retries for PDF generation and email notifications
- One-time ticket validation scoped to the correct cinema
- Sales reports, idempotent cinema settlements, and mutation audit logs
- Unified JSON responses, structured logging, recovery, CORS, rate limiting, and metrics
- Embedded, automatically applied PostgreSQL migrations
- Docker Compose development stack and GitHub Actions CI

## Architecture

```text
HTTP client
    │
    ▼
Gin router
    │
    ├── logging / recovery / audit / metrics / CORS / rate limiting
    └── authentication / permission / ownership checks
    │
    ▼
Handler ──► DTO and validation
    │
    ▼
Application service
    │
    ▼
Domain repository or provider interface
    │
    ├── PostgreSQL / GORM
    ├── Redis
    ├── MinIO
    ├── RabbitMQ
    ├── Elasticsearch
    └── SMTP or development providers
```

Dependencies point inward: domain types and contracts do not depend on Gin, GORM, Redis, or other
infrastructure packages. Concrete implementations are assembled in `cmd/app/main.go`.

## Technology stack

| Area | Technology |
|---|---|
| Language | Go 1.25.8 |
| HTTP | Gin |
| Primary database | PostgreSQL 15, GORM, pgx |
| Cache and ephemeral state | Redis |
| Object storage | MinIO through the S3 API |
| Messaging | RabbitMQ |
| Search | Elasticsearch |
| Authentication | RSA JWT, phone/email OTP, captcha |
| Documents | GoPDF and QR Code |
| Metrics | Prometheus and Grafana |
| Logs | Structured `slog`, Filebeat, Elasticsearch, Kibana |
| Deployment | Docker, Docker Compose, GitHub Actions |

## Project structure

```text
.
├── bootstrap/                         # Environment loading and application constants
├── cmd/
│   ├── app/                           # API composition root
│   └── migrate/                       # Migration-only executable
├── internal/
│   ├── application/
│   │   ├── adapter/                   # Application service contracts
│   │   ├── dto/                       # HTTP/application input and output models
│   │   ├── service/                   # Use cases and background workers
│   │   └── validation/                # Custom validation rules
│   ├── domain/
│   │   ├── broker/                    # Payment and messaging contracts
│   │   ├── entity/                    # Core business entities
│   │   ├── enum/                      # Domain states, roles, and permissions
│   │   ├── exception/                 # Stable application errors
│   │   └── repository/                # Persistence and storage contracts
│   ├── infrastructure/
│   │   ├── broker/                    # RabbitMQ and mock payment implementations
│   │   ├── database/                  # PostgreSQL and Redis clients
│   │   ├── migration/                 # Embedded migration runner
│   │   ├── repository/                # PostgreSQL, Redis, MinIO, and search adapters
│   │   └── seed/                      # Roles, permissions, admin, and location seeds
│   └── presentation/
│       ├── middleware/                # Cross-cutting HTTP policies
│       └── v1/                        # Versioned handlers and routers
├── migrations/                        # Ordered PostgreSQL schema migrations
├── monitoring/                        # Prometheus, Alertmanager, and Grafana config
├── pkg/                               # Reusable response, RSA, image, PDF, and QR helpers
├── tests/e2e/                         # Real HTTP booking lifecycle test
├── Dockerfile
└── docker-compose.yaml
```

## Quick start with Docker

### Prerequisites

- Docker Engine with Docker Compose v2
- Git
- `curl` for the verification commands

### 1. Configure the environment

Create `.env` from the example and fill in all required credentials:

```bash
cp .env.example .env
```

PowerShell equivalent:

```powershell
Copy-Item .env.example .env
```

For a disposable local environment, `.github/ci.env` contains non-production test values and can
be used as a reference. Never reuse those credentials in a deployed environment.

At minimum, configure:

- `SUPER_ADMIN_PHONE` and `SUPER_ADMIN_PASSWORD`
- PostgreSQL `PDB_*` values
- Redis `RDB_*` values
- `RABBITMQ_URL`
- MinIO `S3_*` values and all bucket names
- JWT expiry values and all four `*_PEM_PATH` paths
- OTP, captcha, limiter, and permission-cache durations
- `ELASTICSEARCH_URL`

### 2. Start the API and required dependencies

```bash
docker compose up -d --build api
```

Compose waits for PostgreSQL, Redis, RabbitMQ, MinIO, and Elasticsearch health checks before
starting the API. On first boot the application creates RSA keys in the persistent `api_keys`
volume, applies migrations, synchronizes permissions, and seeds the administrator.

### 3. Verify the deployment

```bash
curl http://localhost:8080/health/live
curl http://localhost:8080/health/ready
```

A ready instance returns HTTP `200` with every dependency marked `UP`.

```json
{
  "success": true,
  "statusCode": 200,
  "message": "Service is ready",
  "data": {
    "status": "UP",
    "checks": {
      "postgres": "UP",
      "redis": "UP",
      "rabbitmq": "UP",
      "object_storage": "UP",
      "elasticsearch": "UP"
    }
  }
}
```

Useful lifecycle commands:

```bash
docker compose logs -f api
docker compose restart api
docker compose down
```

Use `docker compose down -v` only when you intentionally want to remove persistent development
data, including the database, objects, and generated RSA keys.

### Optional operational stack

Running the whole Compose project also starts PgAdmin, Kibana, Filebeat, Prometheus, Alertmanager,
Grafana, and exporters:

```bash
docker compose up -d --build
```

Default host ports are controlled by `.env`. Compose fallbacks include API `8080`, PostgreSQL
`15015`, Redis `15016`, PgAdmin `15017`, MinIO `15018`, MinIO Console `15019`, Kibana `15020`,
Prometheus `15021`, Alertmanager `15022`, and Grafana `15023`.

## Local development

### Prerequisites

- Go 1.25.8 or compatible version
- Running PostgreSQL, Redis, RabbitMQ, MinIO, and Elasticsearch services

Start only the infrastructure in Docker:

```bash
docker compose up -d db redis rabbitmq minio elasticsearch
```

When the API runs on the host, configure service endpoints using their host ports. For example,
PostgreSQL should use `127.0.0.1:15015`, Redis `127.0.0.1:15016`, and MinIO
`http://127.0.0.1:15018` when the default Compose ports are used.

Create the configured key and log directories before the first local run:

```bash
mkdir -p keys logs
go mod download
go run ./cmd/app
```

The RSA key pairs are loaded from the configured paths or generated when they do not exist.

Build the executables:

```bash
go build ./cmd/app
go build ./cmd/migrate
```

## Configuration

Configuration is loaded from `.env` and can be overridden by process environment variables.
Secrets are excluded from structured configuration output.

| Group | Important variables | Purpose |
|---|---|---|
| Runtime | `DEBUG_MODE_ENABLE`, `AUTO_MIGRATE`, `HTTP_DOMAIN`, `HTTP_PORT` | Application behavior and binding |
| CORS | `CORS_ALLOWED_ORIGINS` | Comma-separated trusted browser origins |
| Bootstrap admin | `SUPER_ADMIN_*` | Creates or restores the initial administrator role |
| PostgreSQL | `PDB_HOST`, `PDB_PORT`, `PDB_USER`, `PDB_PASSWORD`, `PDB_DB_NAME`, `PDB_SSL_MODE` | Primary persistence |
| Redis | `RDB_HOST`, `RDB_PORT`, `RDB_PASSWORD`, `RDB_DB_NUMBER` | OTP, captcha, rate limiting, tokens, permission cache |
| RabbitMQ | `RABBITMQ_URL` | Search synchronization events |
| Elasticsearch | `ELASTICSEARCH_URL` | Search index and suggestions |
| MinIO | `S3_ENDPOINT`, `S3_ACCESS_KEY`, `S3_SECRET_KEY`, `S3_REGION`, `S3_*_BUCKET` | Media and private ticket PDFs |
| JWT | `JWT_ACCESS_EXPIRE_MINUTES`, `JWT_REFRESH_EXPIRE_DAYS`, `*_PEM_PATH` | Token signing and expiration |
| Booking | `RESERVATION_HOLD_SECONDS`, `REFUND_CUTOFF_MINUTES` | Seat hold and refund policies |
| OTP and captcha | `OTP_*`, `CAPTCHA_*` | Authentication challenge settings |
| Email | `CINEMA_TICKET_EMAIL`, `CINEMA_TICKET_EMAIL_APP_PASSWORD`, `SMTP_HOST`, `SMTP_HOST_PORT` | Verification and payment notifications |
| Limits | `LIMITER_*`, `USER_PERMISSION_CACHE_EXPIRE_MINUTES` | API throttling and authorization caching |
| Logging | `LOGGER_LEVEL`, `LOGGER_ENCODING`, `LOGGER_NAME` | Structured application logs |

`DEBUG_MODE_ENABLE=true` is intended only for development. It exposes OTP values in responses,
bypasses the captcha requirement on phone OTP requests, writes SMS messages to `SEND_SMS.txt`, and
writes emails to `SEND_EMAIL.txt`.

## Authentication and authorization

### Phone OTP login

1. When debug mode is disabled, obtain a captcha from `POST /v1/captcha/send`.
2. Call `POST /v1/auth/otp/send` with the phone number. Supply `X-Captcha-ID` and
   `X-Captcha-Code` headers outside debug mode.
3. Call `POST /v1/auth/otp/verify` with the phone and OTP.
4. An existing active user is logged in. A new phone number creates a verified user and assigns the
   default `user` role transactionally.

Iranian phone variants are normalized to `09xxxxxxxxx`. OTP consumption is atomic, preventing code
replay.

```bash
curl -X POST http://localhost:8080/v1/auth/otp/send \
  -H "Content-Type: application/json" \
  -d '{"phone":"09123456789"}'
```

```bash
curl -X POST http://localhost:8080/v1/auth/otp/verify \
  -H "Content-Type: application/json" \
  -d '{"phone":"09123456789","code":"12345"}'
```

Authenticated requests use the access token:

```bash
curl http://localhost:8080/v1/tickets \
  -H "Authorization: Bearer <access-token>"
```

Use `POST /v1/auth/refresh` to rotate a session and `POST /v1/auth/logout` to revoke it. Refresh
tokens are tracked in Redis; revoked access tokens are rejected immediately.

### Email verification

Users first add an email through `PUT /v1/users/:id/profile`, then call:

- `POST /v1/auth/email/otp/send`
- `POST /v1/auth/email/otp/verify`

Only verified email addresses receive payment and refund notifications.

### Default roles

| Role | Default scope |
|---|---|
| `admin` | Unrestricted access through `general:all` |
| `support` | User lookup, moderation, reservation/payment support, and ticket validation |
| `user` | Self-service profile, booking, payments, tickets, favorites, comments, and ratings |
| `cinema_owner` | User capabilities plus management of resources belonging to owned cinemas |

Authorization is centralized in `middleware_route_authorization.go`. Every `/v1` route is explicitly
classified as public or protected; startup fails if a registered route lacks a policy. Unknown
protected application routes fail closed. Admins bypass ownership restrictions, while cinema owners
must own the targeted cinema, hall, seat, or showtime.

## Booking lifecycle

```text
Published showtime
    │
    ▼
Select available showtime seats
    │
    ▼
PENDING_PAYMENT reservation + HELD seats
    │
    ├── hold expires ──► EXPIRED reservation + AVAILABLE seats
    ├── user cancels ──► CANCELLED reservation + AVAILABLE seats
    └── payment succeeds
            │
            ▼
        PAID reservation + SOLD seats + ACTIVE ticket
```

Important consistency guarantees:

- Showtime creation snapshots active hall seats and prices.
- Prices and reservation totals are calculated by the server.
- A reservation can include at most ten seats.
- PostgreSQL row locks allow only one concurrent reservation to win each seat.
- A PostgreSQL exclusion constraint rejects overlapping active showtimes in the same hall.
- Pending reservations expire after `RESERVATION_HOLD_SECONDS` and release their seats.
- Public queries expose only published future showtimes from active cinemas and halls.

## Payments, tickets, and outbox

Payments currently use a mock gateway behind the `PaymentGateway` domain interface, allowing a real
provider to be introduced without changing the booking domain.

### Payment flow

- `POST /v1/reservations/:id/payments` requires an `Idempotency-Key` header.
- The payment amount is copied from the reservation; clients cannot submit an amount.
- `POST /v1/payments/:id/verify` accepts `{"success":true}` or a failed result.
- A successful verification atomically marks the payment `SUCCEEDED`, reservation `PAID`, seats
  `SOLD`, creates exactly one ticket, and enqueues outbox events.
- A failed attempt becomes `FAILED` without cancelling the still-valid reservation.
- Reusing an idempotency key returns the original payment instead of charging twice.

### Refund flow

A valid refund atomically changes the payment to `REFUNDED`, cancels the reservation and active
ticket, releases seats, and queues a notification. Refunds are idempotent, unavailable after the
configured cutoff, and forbidden after a ticket has been used. Cancelling a showtime performs the
same state changes for all affected successful payments.

### Ticket documents and validation

- Every successful payment owns exactly one ticket.
- QR payload format: `cinema-ticket:<ticket-code>`.
- PDF tickets are generated asynchronously and stored in the private MinIO ticket bucket.
- `GET /v1/tickets/:id/pdf` returns a short-lived presigned download URL.
- Missing PDFs can be repaired without issuing duplicate tickets.
- Ticket validation opens one hour before the showtime and closes when the showtime ends.
- Successful validation atomically marks the ticket `USED`; subsequent scans are rejected.

### Transactional outbox

PDF generation and verified-email notifications are driven by PostgreSQL outbox events committed in
the same transaction as payment state. Workers claim events using `FOR UPDATE SKIP LOCKED`, recover
stale claims, and retry failures with bounded backoff. Temporary SMTP or MinIO failures therefore do
not lose committed work or block the payment response.

## API overview

The API prefix is `/v1`. Catalog reads, search, showtimes, seats, public comments, and public ratings
are available without authentication. All mutations require JWT authentication and the matching
permission unless explicitly listed as public.

| Area | Main endpoints |
|---|---|
| Health and metrics | `GET /health/live`, `GET /health/ready`, `GET /metrics` |
| Captcha | `POST /v1/captcha/send`, `POST /v1/captcha/verify` |
| Authentication | `/v1/auth/otp/*`, `/v1/auth/email/otp/*`, `/v1/auth/refresh`, `/v1/auth/logout` |
| Users and wallet | `/v1/users`, `/v1/users/:id/profile`, `/v1/users/:id/roles`, `/v1/users/:id/wallet/balance` |
| RBAC | `/v1/roles/*`, `/v1/permissions/*` |
| Catalog | `/v1/genres`, `/v1/persons`, `/v1/movies`, `/v1/provinces`, `/v1/cities`, `/v1/facilities`, `/v1/cinemas` |
| Social | `/v1/favorites/*`, `/v1/comments/*`, `/v1/ratings/*` |
| Search | `GET /v1/search/`, `GET /v1/search/suggest`, `POST /v1/search/sync` |
| Cinema inventory | `/v1/cinemas/:id/halls`, `/v1/halls/:id/seats`, `/v1/showtimes/*` |
| Reservations | `POST/GET /v1/reservations`, `GET/DELETE /v1/reservations/:id` |
| Payments | `/v1/reservations/:id/payments`, `/v1/payments/:id`, `/v1/payments/:id/verify`, `/v1/payments/:id/refund` |
| Tickets | `/v1/tickets`, `/v1/tickets/:id`, `/v1/tickets/:id/qr`, `/v1/tickets/:id/pdf`, `/v1/tickets/validate` |
| Operations | `/v1/cinemas/:id/reports/sales`, `/v1/cinemas/:id/settlements`, `/v1/settlements/:id/paid`, `/v1/audit-logs` |

Collection endpoints use `page` and `limit` query parameters. Invalid values fall back to safe
defaults, and `limit` is capped at `100`. Domain-specific filters, including cinema, movie, hall,
status, and RFC3339 date ranges, are available where applicable.

## Response contract

Every handler and middleware returns the same JSON envelope.

Successful response:

```json
{
  "success": true,
  "statusCode": 200,
  "message": "Tickets retrieved successfully",
  "data": [],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 0
  }
}
```

Error response:

```json
{
  "success": false,
  "statusCode": 409,
  "code": "SEAT_UNAVAILABLE",
  "message": "One or more selected seats are unavailable"
}
```

Error codes are stable and machine-readable. Validation failures, missing routes, unsupported
methods, permission failures, and recovered panics use the same response contract. Internal errors
are logged but are not exposed to clients.

## Database migrations

SQL migrations are embedded in the compiled binary and run before repositories and workers start
when `AUTO_MIGRATE=true`, which is the default.

The migration runner provides:

- Strict version ordering and an application migration ledger
- Import of an existing Goose migration history
- PostgreSQL advisory locking for safe multi-replica startup
- Transactional migrations by default
- Statement-by-statement execution for explicit `NO TRANSACTION` migrations
- Safe retry of idempotent non-transactional schema changes

For a dedicated deployment migration step, set `DATABASE_DSN` and run:

```bash
go run ./cmd/migrate
```

Disable startup migrations with `AUTO_MIGRATE=false` only when migrations are executed separately
as part of the deployment pipeline.

## Observability

| Endpoint | Purpose |
|---|---|
| `GET /health/live` | Confirms the process is alive |
| `GET /health/ready` | Checks PostgreSQL, Redis, RabbitMQ, MinIO, and Elasticsearch |
| `GET /metrics` | Exposes Prometheus metrics |

Application logs use structured JSON and redact sensitive field names such as passwords, secrets,
tokens, authorization headers, OTP codes, and private keys. Mutation audit records include actor,
route, method, status, request ID, source IP, user agent, and duration without blocking the response.

The optional Compose observability stack provides Prometheus, Grafana, Alertmanager, Filebeat,
Elasticsearch, Kibana, Node Exporter, and PostgreSQL Exporter.

## Testing

Run the unit test suite:

```bash
go test ./...
```

Run static checks and compile both executables:

```bash
go vet ./...
go build ./cmd/app ./cmd/migrate
```

The integration suites activate when their environment variables are available:

- `PAYMENT_INTEGRATION_DSN` for transactional booking, payment, outbox, reports, and settlement tests
- `MIGRATION_INTEGRATION_DSN` for isolated fresh-database migration tests
- `TICKET_INTEGRATION_MINIO_ENDPOINT` and corresponding access/secret keys for PDF storage tests

Run the real HTTP lifecycle against an already running API:

```bash
go test -tags=e2e ./tests/e2e -v
```

The E2E flow verifies readiness, administrator OTP login, email verification, catalog and cinema
creation, hall and seat setup, showtime creation, reservation, mock payment, ticket PDF generation,
pagination, sales reporting, refund, reservation cancellation, and audit persistence.

GitHub Actions performs race-enabled tests, coverage generation, E2E verification, `go vet`, Go
builds, and a final Docker image build. CI-only credentials live in `.github/ci.env` and must never
be treated as deployable secrets.

## Security notes

- Never commit `.env`, generated private PEM files, SMTP credentials, or production database secrets.
- Keep `DEBUG_MODE_ENABLE=false` outside local development; debug mode exposes OTP values.
- Use TLS at the edge and encrypted connections to managed dependencies where available.
- Restrict `/metrics`, infrastructure consoles, and database ports at the network layer.
- Replace the mock SMS and payment providers before accepting real customers or money.
- Store deployment secrets in a secret manager instead of image layers or Compose files.
- Back up PostgreSQL and MinIO data and test the restore procedure before production launch.
- Rotate JWT keys deliberately; replacing persisted keys immediately invalidates existing sessions.

## License

This project is released under the [MIT License](LICENSE).
