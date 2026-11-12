# Chirpy HTTP Server

A learning project built with Go while working through HTTP server and backend concepts.  
The project implements a small social API called **Chirpy**, where users can create and read short messages called chirps.

## What this project demonstrates

- Building an HTTP server with Go's standard `net/http` package
- Routing requests with `http.ServeMux`
- JSON request and response handling
- PostgreSQL database integration
- SQL query generation with `sqlc`
- Database migrations with Goose
- Password hashing with Argon2id
- JWT access-token authentication
- Refresh-token creation and revocation
- Middleware for request metrics
- Automated tests with Go's testing package

## Features

- User registration and login
- Secure password storage
- JWT access tokens
- PostgreSQL-backed refresh tokens
- Refresh-token expiration and revocation
- Create chirps for the authenticated user
- List all chirps
- Retrieve a chirp by ID
- Profanity filtering
- File-server metrics
- Development-only database reset endpoint

## API endpoints

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/healthz` | Health check |
| `POST` | `/api/users` | Create a user |
| `POST` | `/api/login` | Log in and receive access and refresh tokens |
| `POST` | `/api/refresh` | Create a new access token |
| `POST` | `/api/revoke` | Revoke a refresh token |
| `POST` | `/api/chirps` | Create an authenticated chirp |
| `GET` | `/api/chirps` | List chirps |
| `GET` | `/api/chirps/{chirpID}` | Get one chirp |
| `GET` | `/api/metrics` | View file-server metrics |
| `POST` | `/admin/reset` | Reset development data |

## Tech stack

- **Language:** Go
- **HTTP:** `net/http`
- **Database:** PostgreSQL
- **Database access:** `sqlc`
- **Migrations:** Goose
- **Authentication:** JWT and refresh tokens
- **Password hashing:** Argon2id

## Getting started

### Requirements

- Go 1.27 or newer
- PostgreSQL
- `sqlc`
- Goose

### Environment variables

Create a `.env` file in the project root:

```env
PLATFORM=dev
DB_URL=postgres://username:password@localhost:5432/users?sslmode=disable
SECRET_KEY=replace-with-a-long-random-secret
```

Do not commit real credentials or secret keys.

### Run the server

```bash
go run ./backend
```

The server starts on:

```text
http://localhost:8080
```

### Run tests

```bash
go test ./...
```

### Regenerate database code

After changing SQL queries, regenerate the Go database code:

```bash
sqlc generate
```

## Project structure

```text
.
├── backend/                # Go HTTP server and its tests
│   ├── main.go             # Application startup and configuration
│   ├── routes.go           # HTTP route registration
│   ├── users.go            # Registration, login, refresh, and revoke handlers
│   ├── chirps.go           # Chirp handlers and validation
│   └── main_test.go        # HTTP handler tests
├── internal/auth/           # Password, JWT, and refresh-token helpers
├── internal/database/       # sqlc-generated database code
├── sql/schema/              # PostgreSQL schema migrations
├── sql/queries/             # SQL queries used by sqlc
└── src/                    # React frontend
```

## Learning goal

This project is being developed incrementally to build a practical understanding of backend development in Go: HTTP routing, authentication, persistence, validation, error handling, and testing.

## Chirpy frontend

The repository includes a React + Vite + TypeScript frontend. From the project root, install dependencies and start it with:

```bash
npm install
npm run dev
```

The API base defaults to `http://localhost:8080`. To use another API host, set `VITE_API_URL` in a `.env` file (see `.env.example`):

```env
VITE_API_URL=http://localhost:8080
```

Run `go run ./backend` in another terminal while developing. Because the frontend uses the configured absolute API URL, the Go server must allow the Vite origin (`http://localhost:5173`) with CORS headers during local development. Alternatively, set a same-origin API URL through a reverse proxy. The frontend itself does not change backend CORS behavior.
