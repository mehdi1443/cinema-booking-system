# Cinema Booking System

A Go/Gin backend for a cinema booking platform, organized with Clean Architecture.

## API response contract

All handlers and middleware use one response envelope. Successful responses contain `success`,
`statusCode`, `message`, and optional `data`/`meta`. Error responses also contain a stable,
machine-readable `code` such as `INVALID_INPUT`, `FORBIDDEN`, or `SEAT_UNAVAILABLE`.

## Authentication

Authentication is phone-based and has two entry endpoints:

1. `POST /v1/auth/otp/send`

   ```json
   { "phone": "09123456789" }
   ```

   The development SMS provider appends the OTP to `SEND_SMS.txt`. The response also contains
   `otp_code` and `expires_in_seconds`. Returning OTP codes is intended for development only and
   must be disabled when a production SMS provider is introduced.

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
- `GET /v1/payments/:id` and `GET /v1/reservations/:id/payments` return payment details/history.

The amount always comes from the reservation and cannot be supplied by the client. Successful
verification atomically changes the payment to `SUCCEEDED`, the reservation to `PAID`, and all held
seats to `SOLD`. A failed mock verification changes only the payment attempt to `FAILED`, allowing a
new attempt with a new idempotency key. Reusing the same key returns the original payment.

## Database changes

Run the Goose migrations before starting the application. The latest migrations rename the legacy
`customer` role to `user`, add `cinemas.owner_user_id`, and create the booking/payment schemas. Cinema and
movie records use soft deletion so historical showtime and reservation references remain intact.

## Verification

```bash
CGO_ENABLED=0 go test ./...
```

Image decoding and JPEG encoding use pure-Go packages, so the project no longer requires a local C
compiler for tests or builds.
