# E-Wallet Backend API

<p align="center">
  <img src="https://img.shields.io/github/go-mod/go-version/MGT06/e-wallet_backend/development?style=for-the-badge&logo=go&logoColor=white&color=00ADD8" alt="Go version" />
  <img src="https://img.shields.io/github/last-commit/MGT06/e-wallet_backend/development?style=for-the-badge&logo=git&logoColor=white" alt="Last commit" />
  <img src="https://img.shields.io/badge/license-MIT-green?style=for-the-badge" alt="License: MIT" />
</p>

<p align="center">
  <img src="https://img.shields.io/github/issues-pr/MGT06/e-wallet_backend?style=flat-square&logo=github" alt="Pull requests" />
  <img src="https://img.shields.io/github/repo-size/MGT06/e-wallet_backend?style=flat-square&logo=github" alt="Repo size" />
  <img src="https://img.shields.io/github/commit-activity/m/MGT06/e-wallet_backend/development?style=flat-square&logo=git&logoColor=white" alt="Commit activity" />
</p>

<p align="center">
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go" /></a>
  <a href="https://github.com/gin-gonic/gin"><img src="https://img.shields.io/badge/Gin-008B8B?style=for-the-badge&logo=gin&logoColor=white" alt="Gin" /></a>
  <a href="https://www.postgresql.org/"><img src="https://img.shields.io/badge/PostgreSQL-4169E1?style=for-the-badge&logo=postgresql&logoColor=white" alt="PostgreSQL" /></a>
  <a href="https://redis.io/"><img src="https://img.shields.io/badge/Redis-DC382D?style=for-the-badge&logo=redis&logoColor=white" alt="Redis" /></a>
  <a href="https://jwt.io/"><img src="https://img.shields.io/badge/JWT-000000?style=for-the-badge&logo=jsonwebtokens&logoColor=white" alt="JWT" /></a>
  <a href="https://www.docker.com/"><img src="https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white" alt="Docker" /></a>
  <a href="https://swagger.io/"><img src="https://img.shields.io/badge/Swagger-85EA2D?style=for-the-badge&logo=swagger&logoColor=black" alt="Swagger" /></a>
  <a href="https://www.xendit.co/"><img src="https://img.shields.io/badge/Xendit-4573FF?style=for-the-badge&logoColor=white" alt="Xendit" /></a>
  <a href="https://github.com/features/actions"><img src="https://img.shields.io/badge/GitHub_Actions-2088FF?style=for-the-badge&logo=githubactions&logoColor=white" alt="GitHub Actions" /></a>
</p>

A RESTful backend for an e-wallet application, built with Go and the Gin framework. It handles authentication, user and PIN management, wallet top-ups through Xendit, peer-to-peer transfers, and transaction history. Built as a team project (Regu Mawar) during the Koda Academy bootcamp, *Gin Exercise Koda Batch 9*.

The frontend lives in a separate repository: [e-wallet_frontend](https://github.com/MGT06/e-wallet_frontend).

## Features

**Authentication**
- Register and login with JWT (HS256). The token is also tracked in Redis, so a logout invalidates it immediately.
- Passwords are hashed before they are stored.
- Forgot password flow: a 6-digit OTP is sent by email, verified, and exchanged for a short-lived reset token.
- Logout removes the session from Redis.

**User**
- View header info, full profile, and dashboard stats (balance, income, expense).
- Update profile (name, phone, email, photo) and delete the profile photo.
- Change password.
- Set, change, and verify a 6-digit transaction PIN.
- Add or remove favorite contacts for quick transfers.

**Transactions**
- Search for people, with an option to show favorites only.
- Top up through Xendit invoices. The payment methods and their tax percentage come from the database.
- Transfer balance to another user, confirmed with the PIN.
- Transaction history with search, pagination, and soft delete.
- Xendit webhook that settles top-ups automatically.

## Tech Stack

| Technology | Purpose |
| --- | --- |
| [Go](https://go.dev/) | Main programming language |
| [Gin](https://gin-gonic.com/) | HTTP router and web framework |
| [PostgreSQL](https://www.postgresql.org/) | Primary relational database |
| [pgx](https://github.com/jackc/pgx) | PostgreSQL driver and connection pool |
| [Redis](https://redis.io/) | Login sessions, OTP, rate limiting, PIN attempt counters |
| [golang-jwt](https://github.com/golang-jwt/jwt) | JSON Web Token creation and validation |
| [validator](https://github.com/go-playground/validator) | Request validation |
| [godotenv](https://github.com/joho/godotenv) | Loads `.env` configuration |
| [Resty](https://github.com/go-resty/resty) | HTTP client for the Xendit API |
| [gomail](https://github.com/go-gomail/gomail) | Sends OTP emails over SMTP |
| [Xendit](https://www.xendit.co/) | Payment gateway for top-ups |
| [swaggo](https://github.com/swaggo/swag) | Swagger/OpenAPI documentation |
| [golang-migrate](https://github.com/golang-migrate/migrate) | Database migrations |
| Docker and GitHub Actions | Containerization and image publishing |

## Architecture

The code follows a layered structure, so each layer has one job:

```
Request -> Router -> Middleware -> Controller -> Service -> Repository -> PostgreSQL
                                                    |
                                                    +-> Redis / Xendit / SMTP
```

- **Router** registers routes and attaches middleware.
- **Middleware** handles CORS, JWT and Redis session checks, and rate limiting.
- **Controller** parses and validates requests and builds responses.
- **Service** contains the business logic.
- **Repository** runs the SQL queries.

### Top-up flow

1. The client calls `POST /transaction/topup`. The API stores a `pending` top-up and creates a Xendit invoice.
2. The API returns a `payment_url` and an expiry time. The user pays on the Xendit page.
3. Xendit calls `POST /webhook/xendit` with the `x-callback-token` header.
4. The API verifies the token and ignores duplicate events (unique `event_id` and status).
5. It checks the invoice ID and paid amount, then marks the top-up `success`, adds the balance, and records a transaction. An `EXPIRED` event marks the top-up `expired`.

## Project Structure

```
.
├── .github/workflows/   # CI: build and push the Docker image
├── cmd/                 # Application entry point (main.go)
├── db/
│   ├── migrations/      # SQL migrations (up/down)
│   └── seeders/         # Dummy data
├── docs/                # Generated Swagger docs
├── internal/
│   ├── config/          # Env, PostgreSQL, and Redis setup
│   ├── controller/      # HTTP handlers
│   ├── dto/             # Request and response structs
│   ├── mailer/          # SMTP mailer and email template
│   ├── message/         # Shared error definitions
│   ├── middleware/      # Auth, CORS, and rate limit
│   ├── models/          # Database models
│   ├── repository/      # Data access layer
│   ├── router/          # Route definitions
│   ├── service/         # Business logic
│   └── utils/           # OTP, random string, Redis and validator helpers
├── pkg/                 # Reusable packages (JWT and hash)
├── Dockerfile
├── makefile
├── .env.example
└── redis.conf.example
```

## Getting Started

### Prerequisites

- [Go](https://go.dev/dl/) (the version in `go.mod`)
- [Docker](https://www.docker.com/)
- [golang-migrate](https://github.com/golang-migrate/migrate) CLI (`migrate`)
- `psql` and `make` (for seeding)
- [swag](https://github.com/swaggo/swag) CLI (only if you regenerate the docs)
- A [Xendit](https://www.xendit.co/) account (test mode is fine) and an SMTP account

### 1. Clone the repository

```bash
git clone -b development https://github.com/MGT06/e-wallet_backend.git
cd e-wallet_backend
```

### 2. Configure the environment

```bash
cp .env.example .env
```

Fill in `.env` using the [Environment Variables](#environment-variables) table below.

### 3. Configure Redis

```bash
cp redis.conf.example redis.conf
```

Adjust `redis.conf` (password, ACL users, persistence, memory limit). Every option in the example file is commented out, so enable what you need.

### 4. Start PostgreSQL and Redis

```bash
docker pull postgres:16
docker pull redis:8.10.2
```

PostgreSQL (it reads the `POSTGRES_*` values from `.env`):

```bash
docker run --name {CONTAINER_NAME} -d --env-file .env \
  -p {HOST_PORT}:5432 \
  -v {DOCKER_VOLUME}:/var/lib/postgresql \
  postgres:16
```

Redis:

```bash
docker run -d --name {CONTAINER_NAME} \
  -p {HOST_PORT}:6379 \
  -v {PATH_TO_redis.conf}:/usr/local/etc/redis/redis.conf \
  -v {DOCKER_VOLUME}:/data \
  --restart unless-stopped \
  redis:8.10.2 redis-server /usr/local/etc/redis/redis.conf
```


### 5. Migrate and seed the database

```bash
make migrate-up
make db-seed
```

### 6. Run the application

```bash
go mod download
go run ./cmd/main.go
```


## Environment Variables

| Variable | Description |
| --- | --- |
| `GIN_MODE` | Optional. Set to `release` in production |
| `HOST_SERVER`, `PORT_SERVER` | Address and port the server binds to |
| `HOST_SWAGG`, `PORT_SWAGG` | Host and port shown in the Swagger UI |
| `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB` | PostgreSQL credentials and database name |
| `POSTGRES_HOST`, `POSTGRES_PORT` | PostgreSQL address |
| `JWT_SECREET_KEY` | Secret used to sign JWTs (the spelling matches the code) |
| `JWT_ISSUER` | JWT issuer claim |
| `REDIS_USER`, `REDIS_PASSWORD` | Redis credentials |
| `REDIS_HOST`, `REDIS_PORT` | Redis address |
| `SMTP_HOST`, `SMTP_PORT`, `SMTP_USER`, `SMTP_PASS` | SMTP server used to send OTP emails |
| `PROFILE_URL` | Public URL path used to serve profile photos |
| `XENDIT_API_URL` | Xendit API base URL |
| `XENDIT_SECRET_KEY` | Xendit secret API key |
| `XENDIT_WEBHOOK_TOKEN` | Token that Xendit sends in the `x-callback-token` header |
| `XENDIT_INVOICE_DURATION` | How long an invoice stays valid |
| `XENDIT_SUCCESS_REDIRECT_URL` | Redirect URL after a successful payment |
| `XENDIT_FAILURE_REDIRECT_URL` | Redirect URL after a failed payment |

## Database

Migrations live in `db/migrations` and create these tables:

| Table | Description |
| --- | --- |
| `accounts` | Users: profile, hashed password, hashed PIN, and balance (never negative) |
| `payment_methods` | Top-up methods (`bank` or `ewallet`) with Xendit code and tax percent |
| `topups` | Top-up requests, Xendit invoice, payment URL, expiry, and status |
| `transfers` | Transfers between two different accounts |
| `transaction` | Per-account ledger with `balance_before` and `balance_after`, linked to a top-up or a transfer |
| `favorites` | Favorite contacts for each account |
| `webhook_events` | Raw webhook payloads, unique per event and status, for idempotency |

Status values: transfers are `pending`, `success`, or `failed`. Top-ups are `pending`, `success`, `failed`, or `expired`.

## API Reference

Interactive documentation is served by Swagger UI at **`/docs/index.html`** 

Protected routes require the header `Authorization: Bearer <token>`.

### Auth

| Method | Endpoint | Auth | Description |
| --- | --- | --- | --- |
| `POST` | `/auth/register` | No | Create a new account |
| `POST` | `/auth/login` | No | Log in and receive a JWT |
| `DELETE` | `/auth/logout` | Yes | Log out and invalidate the session |
| `POST` | `/auth/forgot` | No | Send an OTP to the email (rate limited) |
| `POST` | `/auth/check-otp` | No | Verify the OTP and receive a reset token |
| `POST` | `/auth/reset-password` | No | Set a new password with the reset token |

### User

| Method | Endpoint | Auth | Description |
| --- | --- | --- | --- |
| `GET` | `/user/info` | Yes | Basic user info for the app header |
| `GET` | `/user/profile` | Yes | Full profile |
| `PATCH` | `/user/profile` | Yes | Update profile (multipart form, optional photo) |
| `DELETE` | `/user/image` | Yes | Delete the profile photo |
| `PATCH` | `/user/password` | Yes | Change password |
| `POST` | `/user/set-pin` | Yes | Set the 6-digit PIN |
| `PATCH` | `/user/change-pin` | Yes | Change the PIN |
| `POST` | `/user/check-pin` | Yes | Verify the PIN |
| `POST` | `/user/favorite` | Yes | Add or remove a favorite contact |
| `GET` | `/user/history` | Yes | Transaction history (`search`, `limit`, `page`) |
| `DELETE` | `/user/history` | Yes | Delete history entries |
| `GET` | `/user/stats-dashboard` | Yes | Balance, income, and expense summary |

### Transaction

| Method | Endpoint | Auth | Description |
| --- | --- | --- | --- |
| `GET` | `/transaction/people` | Yes | Search people (`search`, `limit`, `is_favorit`) |
| `GET` | `/transaction` | Yes | Get a transfer recipient by `user_id` |
| `GET` | `/transaction/topup-methode` | Yes | List available top-up payment methods |
| `POST` | `/transaction/topup` | Yes | Create a top-up and get the Xendit payment URL |
| `POST` | `/transaction/transfer` | Yes | Transfer balance (`receiver_id`, `amount`, `notes`, `pin`) |

### Webhook

| Method | Endpoint | Auth | Description |
| --- | --- | --- | --- |
| `POST` | `/webhook/xendit` | `x-callback-token` header | Receives Xendit invoice status updates |

## Security Notes

- **Sessions:** a JWT is valid only while its session exists in Redis (60 minutes). Logging out removes it right away.
- **Passwords and PINs:** both are hashed. Passwords need at least 8 characters, and the PIN is exactly 6 digits.
- **Forgot password:** the endpoint is limited to 10 requests per hour per IP and one OTP per email every 60 seconds. It returns the same response whether or not the email exists. The OTP is stored hashed and expires after 3 minutes, with at most 3 attempts. The reset token lasts 10 minutes.
- **Webhook:** requests are checked with a constant-time token comparison, and duplicate events are ignored.
- **Profile photos:** only `.jpg`, `.jpeg`, and `.png` files up to 2 MB are accepted.
- **CORS:** allowed origins are defined in `internal/middleware/cors.middleware.go`. Update them for your own frontend.

## Docker and CI/CD

The `Dockerfile` is a multi-stage build: it compiles a static Go binary and runs it in a small Alpine image.

```bash
docker build -t e-wallet-backend .
docker run -d --env-file .env -p 5000:5000 e-wallet-backend
```

The GitHub Actions workflow (`.github/workflows/action.yml`) runs on every push to the `development` branch. It builds the image and pushes it to GitHub Container Registry with two tags: `dev` and `sha-<short commit>`.

```bash
docker pull ghcr.io/mgt06/e-wallet_backend:dev
```

## Makefile Commands

| Command | Description |
| --- | --- |
| `make migrate-create NAME=<name>` | Create a new migration pair |
| `make migrate-up` | Apply all migrations |
| `make migrate-down N=<count>` | Roll back N migrations |
| `make db-seed` | Insert dummy data |
| `make db-reset` | Roll back, migrate, and seed again |
| `make swag` | Regenerate the Swagger docs in `docs/` |

## Changelog

| Version | Desctiption |
| ------- | ----------- |
| 1.0.0   | Initial app |

- 1.1.1
  - Feat history .... by [almarufh](https://github.com/almarufh/)
  - Feat transfer .... by [given](https://github.com/MGT06/)
  - Feat topup .... by [given](https://github.com/MGT06/)
  - Feat forgot password .... by [given](https://github.com/almarufh/)
  - Feat cors .... by [almarufh](https://github.com/almarufh/)
  - Feat hash password .... by [almarufh](https://github.com/almarufh)
  - Feat Autentication .... by [almarufh](https://github.com/almarufh)


## How to Contribute

1. Fork this repository.
2. Create a feature branch: `git checkout -b feat/your-feature`.
3. Commit your changes.
4. Push the branch and open a Pull Request.

## License

This project is licensed under the [MIT License](https://opensource.org/license/mit).

## Related Project

[![frontend](https://img.shields.io/badge/FRONTEND-REACT-blue?logo=react&logoColor=WHITE)](https://github.com/MGT06/e-wallet_frontend)
