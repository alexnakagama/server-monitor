# Server Monitor

Server Monitor is a system for collecting and monitoring resource usage from servers. It is composed of two components:

- **Server** (`cmd/server`): an HTTP API that manages users, servers, agents and the metrics they report. Data is persisted in PostgreSQL.
- **Agent** (`cmd/agent`): a lightweight program deployed on each monitored host. It periodically collects CPU, memory, disk and network statistics and reports them to the server.

The collected metrics can be queried through the API and are intended to be consumed by dashboards or other monitoring tooling.

## Contents

- [Architecture](#architecture)
- [Features](#features)
- [Tech stack](#tech-stack)
- [Repository layout](#repository-layout)
- [Prerequisites](#prerequisites)
- [Setup](#setup)
- [Configuration](#configuration)
- [Database migrations](#database-migrations)
- [API reference](#api-reference)
- [Contributing](#contributing)

## Architecture

```
+------------------+        +------------------+        +---------------+
|  Monitored host   |        |  Monitored host   |        |  Monitored host
|    Agent          |        |    Agent          |        |    Agent
+------------------+        +------------------+        +---------------+
        |                          |                           |
        +--------------------------+---------------------------+
                                   |
                          (HTTPS/HTTP, Bearer token)
                                   |
                          POST /agents/metrics
                                   |
                        +---------------------+
                        |      Server API      |
                        |      (Go, net/http)  |
                        +---------------------+
                                   |
                                   v
                        +---------------------+
                        |      PostgreSQL      |
                        +---------------------+
```

The agent authenticates with a token issued at provisioning time. The server stores only the SHA-256 hash of the token. Metrics are associated with the server that the agent belongs to.

## Features

- User registration, login and profile management.
- User authentication with PASETO v4 symmetric tokens.
- Passwords hashed with Argon2id.
- Role-based access control with `admin` and `user` roles.
- Server inventory management (register, query and update servers).
- Agent provisioning and revocation per server.
- Automatic metric collection by the agent (CPU, memory, disk and network).
- Liveness and readiness health endpoints.
- Graceful shutdown on `SIGINT` / `SIGTERM`.

## Tech stack

- [Go](https://go.dev) 1.22+ (standard library `net/http` with method and path routing)
- [PostgreSQL](https://www.postgresql.org)
- [pgx](https://github.com/jackc/pgx) v5 (driver and connection pool)
- [go-paseto](https://github.com/aidanwoods/go-paseto) v4 (token handling)
- [argon2id](https://github.com/alexedwards/argon2id) (password hashing)
- [godotenv](https://github.com/joho/godotenv) (environment configuration)
- [golang-migrate](https://github.com/golang-migrate/migrate) (database migrations)

## Repository layout

```
cmd/
  server/                 Server entrypoint
  agent/                  Agent entrypoint
internal/
  agent/                  Agent client and metric collectors
    collector/            CPU, memory, disk and network collectors
  auth/                   PASETO helpers, auth and authorization middleware
  model/                  Domain models and validation
  server/
    db/                   PostgreSQL connection
    handler/              HTTP handlers
    repository/           Data access layer
    service/              Business logic
    errors_custom/        Domain errors
db/
  migrations/             SQL schema migrations
web/                      Web assets (reserved)
docker-compose.yml        PostgreSQL for local development
Makefile                  Migration commands
```

## Prerequisites

- Go 1.22 or later
- PostgreSQL 14 or later (or Docker with Docker Compose)
- [golang-migrate](https://github.com/golang-migrate/migrate) CLI (for running migrations)

## Setup

### 1. Create the configuration file

Copy `.env` and adjust the values to match your environment:

```
DATABASE_URL=postgres://postgres:postgres@localhost:5432/server_monitor?sslmode=disable
PASETO_KEY=<32-byte hex key>
API_URL=http://localhost:8080
AGENT_TOKEN=<agent token>
```

`PASETO_KEY` must be a 32-byte key encoded as 64 hexadecimal characters.

### 2. Start PostgreSQL

```bash
docker compose up -d
```

### 3. Run the migrations

```bash
make migrate-up
```

### 4. Run the server

```bash
go run ./cmd/server
```

The server listens on `:8080` by default.

### 5. Provision an agent

Agents are created through the API. They belong to a server, so a server must exist first.

```bash
# create a server (admin role required)
curl -X POST http://localhost:8080/servers \
  -H "Authorization: Bearer <user-token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"web-01","hostname":"web-01.example.com","os":"linux"}'

# create an agent for the server (admin role required)
curl -X POST http://localhost:8080/servers/1/agents \
  -H "Authorization: Bearer <user-token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"web-01-agent"}'
```

The response contains a token that must be stored by the agent operator. It is only shown once and cannot be retrieved later.

### 6. Run the agent

Set `API_URL` and `AGENT_TOKEN` in `.env`, then start the agent on the monitored host:

```bash
go run ./cmd/agent
```

The agent collects and reports metrics every 100 seconds by default.

## Configuration

The application is configured through environment variables loaded from `.env`.

| Variable       | Component | Description                                             |
| -------------- | --------- | ------------------------------------------------------- |
| `DATABASE_URL` | server    | PostgreSQL connection string.                           |
| `PASETO_KEY`   | server    | 64-character hex key used to encrypt and verify tokens. |
| `API_URL`      | agent     | Base URL of the server API.                             |
| `AGENT_TOKEN`  | agent     | Token used by the agent to authenticate with the server. |

## Database migrations

Migrations live in `db/migrations` and are managed with golang-migrate.

```bash
make migrate-up       # apply all pending migrations
make migrate-down     # roll back the last applied migration
```

## API reference

All endpoints respond with JSON. Authenticated endpoints require an `Authorization: Bearer <token>` header.

### Public

| Method | Path                 | Description                |
| ------ | -------------------- | -------------------------- |
| POST   | `/users/register`    | Register a new user.       |
| POST   | `/users/login`       | Login and obtain a token.  |
| GET    | `/health/live`       | Liveness probe.            |
| GET    | `/health/ready`      | Readiness probe (DB check).|

### User profile (authenticated)

| Method | Path                         | Description            |
| ------ | ---------------------------- | ---------------------- |
| GET    | `/users/profile/me`          | Get the current user.  |
| PUT    | `/users/profile/me`          | Update the profile.    |
| DELETE | `/users/profile/me`          | Delete the profile.    |
| PUT    | `/users/profile/me/password` | Change the password.   |

### Servers (authenticated)

| Method | Path                                      | Access          | Description                  |
| ------ | ----------------------------------------- | --------------- | ---------------------------- |
| POST   | `/servers`                                | admin           | Create a server.             |
| GET    | `/servers`                                | authenticated   | List all servers.            |
| GET    | `/servers/{name}`                         | authenticated   | Get a server by name.        |
| GET    | `/servers/os/{os}`                        | authenticated   | Get servers by OS.           |
| GET    | `/servers/hostname/{hostname}`            | authenticated   | Get a server by hostname.    |
| PUT    | `/servers/hostname/{hostname}`            | admin           | Update name and OS.          |
| PATCH  | `/servers/hostname/{hostname}/name`       | admin           | Update the server name.      |
| PATCH  | `/servers/hostname/{hostname}/os`         | admin           | Update the server OS.        |
| DELETE | `/servers/hostname/{hostname}`            | admin           | Delete a server by hostname. |

### Agents (authenticated)

| Method | Path                         | Access        | Description                        |
| ------ | ---------------------------- | ------------- | ---------------------------------- |
| POST   | `/servers/{serverID}/agents` | admin         | Provision an agent for a server.   |
| GET    | `/agents/{agentID}`          | authenticated | Get an agent by ID.                |
| DELETE | `/agents/{agentID}`          | admin         | Delete (revoke) an agent.          |

### Metrics

| Method | Path                          | Auth          | Description                                |
| ------ | ----------------------------- | ------------- | ------------------------------------------ |
| POST   | `/agents/metrics`             | agent token   | Report a metric sample (used by the agent).|
| GET    | `/metrics/{metricID}`         | user token    | Get a metric by ID.                        |
| GET    | `/metrics/server/{serverID}`  | user token    | Get metrics for a server.                  |

`GET /metrics/server/{serverID}` accepts an optional `limit` query parameter between 1 and 1000. It defaults to 100.

## Contributing

1. Fork the repository and create a feature branch.
2. Keep changes scoped to a single concern.
3. Run `go vet ./...` and `go build ./...` before submitting a pull request.
4. Describe the change and any migration requirements in the pull request description.
