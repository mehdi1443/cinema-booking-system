# Cinema Booking System

A Go/Gin backend for a cinema booking platform, organized with Clean Architecture.

## API response contract

All handlers and middleware use one response envelope. Successful responses contain `success`,
`statusCode`, `message`, and optional `data`/`meta`. Error responses also contain a stable,
machine-readable `code` such as `INVALID_INPUT`, `FORBIDDEN`, or `SEAT_UNAVAILABLE`.
Unknown routes, unsupported methods, and recovered panics use the same JSON contract. Liveness is
available at `GET /health/live`; dependency readiness is available at `GET /health/ready` and
checks PostgreSQL, Redis, MinIO, RabbitMQ, and Elasticsearch. Prometheus metrics are exposed at
`GET /metrics`.

## Authentication

Authentication is phone-based and has two entry endpoints:

1. `POST /v1/auth/otp/send`

   ```json
   { "phone": "09123456789" }
   ```

   The development SMS provider appends the OTP to `SEND_SMS.txt`. The response contains
   `expires_in_seconds`; `otp_code` is included only when `DEBUG_MODE_ENABLE=true`. Returning OTP
   codes is intended for development only and must be disabled when a production SMS provider is
   introduced.

2. `POST /v1/auth/otp/verify`

   ```json
   { "phone": "09123456789", "code": "123456" }
   ```

   A valid OTP logs in an existing active user. If the phone has no account, the service creates a
   verified account transactionally, assigns the `user` role, and returns access/refresh tokens.
   Iranian phone variants are normalized to `09xxxxxxxxx`. OTP consumption is atomic, so a code
   cannot be replayed.

Session endpoints:

- `POST /v1/auth/refresh`
- `POST /v1/auth/logout` (Bearer access token required)

Users may add an email through `PUT /v1/users/:id/profile`, then verify it using the authenticated
`POST /v1/auth/email/otp/send` and `POST /v1/auth/email/otp/verify` endpoints. In debug mode email
is appended to `SEND_EMAIL.txt`; otherwise the configured SMTP provider is used.

## Roles

- `admin`: unrestricted access through `general:all`
- `support`: user lookup and comment moderation plus normal self-service access
- `user`: own profile, wallet read, favorites, comments, and ratings
- `cinema_owner`: normal user access plus management of owned cinemas

Catalog/search reads are public. Mutations, synchronization endpoints, user administration, wallet
adjustments, and RBAC management require JWT and their matching permission. Unknown application
routes fail closed to admin-only access. Cinema mutation checks the persisted `owner_user_id` unless
the caller is an admin.

## Halls, showtimes, and temporary reservations

The first booking phase is available through these endpoint groups:

- Halls and seats: `POST/GET /v1/cinemas/:id/halls`, `GET/PUT/DELETE /v1/halls/:id`,
  `POST/GET /v1/halls/:id/seats`, and `PUT/DELETE /v1/seats/:id`
- Showtimes and pricing: `POST /v1/halls/:id/showtimes`, `GET /v1/showtimes`,
  `GET/PUT/DELETE /v1/showtimes/:id`, `PUT /v1/showtimes/:id/prices`, and
  `GET /v1/showtimes/:id/seats`
- Temporary reservations: `POST/GET /v1/reservations` and
  `GET/DELETE /v1/reservations/:id`

Creating a showtime takes a `base_price` and optional `type_prices`. The service snapshots the
hall's active seats and their prices for that showtime, preventing later hall edits from changing a
customer's selection. Overlapping active showtimes in one hall are rejected by a PostgreSQL
exclusion constraint.

Creating a reservation takes a `showtime_id` and up to ten `showtime_seat_ids`. Prices and totals
are always calculated by the server. Selected seats are locked in one database transaction, so
concurrent attempts can produce only one winner. A pending reservation expires after
`RESERVATION_HOLD_SECONDS` (default: 600); expiration releases its seats. The cinema owner may
manage only halls, seats, showtimes, and prices belonging to their own cinema. Reservation reads and
cancellations are owner-scoped for normal users, while support/admin access is permission-based.
Public showtime endpoints return only published, upcoming showtimes from active cinemas and halls.

## Mock payments

Payments are intentionally backed by a mock gateway until a real provider is configured:

- `POST /v1/reservations/:id/payments` creates a payment and requires an `Idempotency-Key` header.
- `POST /v1/payments/:id/verify` accepts `{ "success": true }` or `{ "success": false }`.
- `POST /v1/payments/:id/refund` performs an idempotent mock refund before the configured cutoff.
- `GET /v1/payments/:id` and `GET /v1/reservations/:id/payments` return payment details/history.

Reservation, reservation-payment, and ticket list endpoints accept `page` and `limit` (maximum
100) and return the standard pagination metadata.

The amount always comes from the reservation and cannot be supplied by the client. Successful
verification atomically changes the payment to `SUCCEEDED`, the reservation to `PAID`, and all held
seats to `SOLD`. A failed mock verification changes only the payment attempt to `FAILED`, allowing a
new attempt with a new idempotency key. Reusing the same key returns the original payment.

A successful refund atomically marks the payment `REFUNDED`, cancels its reservation and active
ticket, and releases the sold seats. Refunds are accepted only before
`REFUND_CUTOFF_MINUTES` (default: 120) ahead of showtime, and a used ticket cannot be refunded.
Cancelling an entire showtime automatically mock-refunds its successful payments, cancels active
tickets and reservations, and releases seats in the same database transaction.

## Tickets and admission

A successful payment issues exactly one ticket in the same database transaction. Users can access
their tickets through `GET /v1/tickets`, `GET /v1/tickets/:id`, and the PNG QR endpoint
`GET /v1/tickets/:id/qr`. The QR payload uses the form `cinema-ticket:<ticket-code>`.

Each successful payment queues generation of a PDF admission ticket, uploads it to the private MinIO
bucket configured by `S3_TICKET_PDF_BUCKET`, and stores its object URL/key in PostgreSQL. The
authorized owner (or a user with `tickets:read:any`) can call `GET /v1/tickets/:id/pdf` to receive a
15-minute presigned download URL. Repeating payment verification repairs a missing PDF without
issuing a duplicate ticket.

PDF generation and verified-email notifications for payment success/refund use a transactional
PostgreSQL outbox. Events are committed in the same transaction as the payment state and processed
with bounded retries, so a temporary SMTP or MinIO outage does not lose the work.

Cinema owners validate tickets for their own cinemas with `POST /v1/tickets/validate`; support and
admin users may validate across cinemas. Validation opens one hour before the showtime and closes at
its end time. A successful scan marks the ticket `USED` atomically, so subsequent scans are rejected.

## Reports, settlements, and audit

- `GET /v1/cinemas/:id/reports/sales?from=<RFC3339>&to=<RFC3339>` returns gross, refunded, and net sales.
- `GET /v1/cinemas/:id/settlements` is available to the cinema owner and administrators.
- `POST /v1/cinemas/:id/settlements` snapshots completed-showtime payments exactly once; only an
  administrator may create it.
- `PATCH /v1/settlements/:id/paid` idempotently marks the settlement paid.
- `GET /v1/audit-logs` provides an administrator-only, paginated mutation audit trail.

Settlement items have a unique payment constraint. Only successful payments for showtimes that
have ended are eligible, preventing refundable future sales and duplicate payouts.

## Database changes

Migrations are embedded in the binary and applied automatically before repositories/workers start
when `AUTO_MIGRATE=true` (the default). The runner imports an existing Goose migration ledger and
uses a PostgreSQL advisory lock so concurrent replicas cannot migrate simultaneously. For a
migration-only deployment step, set `DATABASE_DSN` and run `go run ./cmd/migrate`.

## Runtime configuration

Copy `.env.example` to `.env` and set the database, Redis, storage, RabbitMQ, JWT, and bootstrap
administrator values before starting the application. `SUPER_ADMIN_PHONE` and
`SUPER_ADMIN_PASSWORD` are required for the first startup; later starts may omit them once an active
administrator exists. Credentials are deliberately not stored in source code. The
local Docker Compose file includes the API image plus PostgreSQL, Redis, RabbitMQ, MinIO,
Elasticsearch and monitoring services. Start the API and its required dependencies with:

```bash
docker compose up -d --build api
```

## Verification

```bash
CGO_ENABLED=0 go test ./...
```

The GitHub Actions workflow runs tests, vet, build, container build, PostgreSQL/MinIO integration
tests, and the tagged HTTP flow from OTP login through reservation, mock payment, ticket PDF,
reporting, refund, and audit. To run that HTTP test against a running local API:

```bash
go test -tags=e2e ./tests/e2e -v
```

Image decoding and JPEG encoding use pure-Go packages, so the project no longer requires a local C
compiler for tests or builds.
