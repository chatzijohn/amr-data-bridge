# AMR Data Bridge

A minimal, performant Go service that ingests external data, transforms it, and persists it to an internal datastore. Designed for clarity, testability, and ease of extension — with zero unnecessary dependencies.

---

## Project Goals

- ✅ Clean, idiomatic Go architecture (no frameworks)
- ✅ Minimal dependencies — stdlib where possible
- ✅ Type-safe DB access via `sqlc`
- ✅ Modular structure: config, transport, business logic, persistence
- ✅ Easy to run, test, and deploy

Future roadmap includes:

- Observability (tracing)
- Input validation and rate limiting
- Contract tests and schema versioning
- Graceful shutdown and retries

---

## Core Components

| Layer                    | Description                                        |
| ------------------------ | -------------------------------------------------- |
| `cmd/`                   | Entry point + service bootstrap logic              |
| `internal/`              | Domain logic, HTTP handlers, middleware, DB        |
| `config/`                | Environment-based config loading & validation      |
| `db/`                    | SQL queries generated via [sqlc](https://sqlc.dev) |
| `scripts/`               | Setup, seed, migration scripts (optional)          |
| `observability/metrics/` | Prometheus metrics handler                         |

---

## Getting Started

### 1. Clone and configure

```bash
git clone https://github.com/chatzijohn/amr-data-bridge.git
cd amr-data-bridge
cp .env.example .env
```

### 2. Build and run

```bash
go build -o amr-bridge ./cmd/server
./amr-bridge
```

Or run directly:

```bash
go run ./cmd/server
```

## Dependencies

- Go 1.25.1+

- sqlc
  for SQL generation

- PostgreSQL
  or compatible datastore
