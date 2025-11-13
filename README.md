# bootdev-chirpy

A small microservice built as part of the **"Learn HTTP Servers in Go"** course from Boot.dev. This project focuses on understanding how to build HTTP servers in Go using only the standard library: routing, handlers, middlewares, authentication, state management, and database operations.

## Overview

Chirpy is a lightweight Twitter-like API where users can:

- Create an account
- Log in using email and password
- Create, list, view, and delete chirps (messages up to 140 characters)
- Update their password
- Authenticate using JWT
- Refresh and revoke tokens
- Receive simulated webhook events
- Access a simple metrics dashboard
- Reset the entire system in development mode

The goal is educational: to understand how an HTTP server works from scratch.

## Tech Stack

### Language

- Go 1.22+

### Database

- PostgreSQL
- `sqlc` for type-safe query generation
- `goose` for migration management

### Dependencies

```
go get github.com/google/uuid
go get github.com/lib/pq
go get github.com/joho/godotenv
```

### Architecture Highlights

- Manual routing with `http.ServeMux`
- `apiConfig` struct for dependency injection
- Custom middlewares
- `sqlc`-generated strongly typed queries
- Environment variables loaded with `godotenv`
- Dev-only reset logic via environment flags

## Main Endpoints

### Users

- `POST /api/users` — Create a new user
- `PUT /api/users` — Update email and password
- `POST /api/login` — Authenticate
- `POST /api/refresh` — Refresh access token
- `POST /api/revoke` — Revoke refresh token

### Chirps

- `POST /api/chirps` — Create a chirp
- `GET /api/chirps` — List chirps (supports `author_id` and `sort=asc|desc`)
- `GET /api/chirps/{chirpID}` — Retrieve a chirp
- `DELETE /api/chirps/{chirpID}` — Delete a chirp (only author allowed)

### Admin

- `GET /admin/metrics` — Simple HTML metrics page
- `POST /admin/reset` — Reset all data (dev only)

### System

- `GET /api/healthz` — Health check
- `POST /api/polka/webhooks` — Fake webhook integration

## Project Structure

```
bootdev-chirpy/
 ├── cmd/
 ├── internal/
 │   └── database/      → sqlc-generated code
 ├── migrations/        → goose migration files
 ├── helpers/           → JSON helpers and utilities
 ├── main.go            → routes, server setup, controllers
 └── .env.example
```

## Environment Variables

```
DB_URL=postgres://user:pass@localhost:5432/chirpy?sslmode=disable
PLATFORM=dev
SECRET_TOKEN=your_jwt_secret
POLKA_KEY=fake_polka_webhook_key
```

## How to Run

### 1. Apply migrations

```
goose up
```

### 2. Generate database code

```
sqlc generate
```

### 3. Start the server

```
go run main.go
```

Server will run at:

```
http://localhost:8080
```

## Testing

This project relies on Boot.dev’s automated tests. Local test files are not included, as the focus is learning HTTP fundamentals.

## Notes

- Manual JWT authentication to reinforce fundamentals.
- Chirp filtering implemented with simple slice/string operations.
- Sorting implemented through query parameters.
- Metrics dashboard demonstrates atomic counters.
- Code prioritizes clarity over optimization, reflecting the educational nature of the project.

## License

Free to use for learning and experimentation.
