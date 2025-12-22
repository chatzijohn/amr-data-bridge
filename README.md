# AMR Data Bridge

A minimal, performant Go service that ingests external data, transforms it, and persists it to an internal datastore. Designed for clarity, testability, and a strict separation of concerns.

---

## Guiding Principle: Separation of Concerns

The architecture of this project is explicitly designed to separate responsibilities into distinct layers. This makes the codebase easier to understand, maintain, test, and extend.

-   **Handler (`transport`):** Manages HTTP only. It decodes requests, validates input, calls the service layer, and encodes the response. It contains no business logic.
-   **Service (`service`):** Contains the core business logic. It orchestrates operations but is completely unaware of the database. It calls methods on the repository layer to persist data.
-   **Repository (`repository`):** Manages all communication with the database. It is responsible for executing SQL queries, managing transactions, and handling database-specific types.

---

## Architecture Overview

This project follows a classic layered architecture. A request flows from the HTTP transport layer down to the database and back.

```
+--------------------------------+
|      HTTP Request/Response     |
+--------------------------------+
                |
                v
+---------------------------------------+
|    Transport (internal/transport/http)|
|   - Server startup                    |
|   - Routing (router)                  |
|   - Request/Response (handler)        |
+---------------------------------------+
                |
                v
+--------------------------------------+
|      Service (internal/service)      |
|   - Core business logic              |
|   - Orchestrates operations          |
|   - Database-agnostic                |
+--------------------------------------+
                |
                v
+--------------------------------------+
|    Repository (internal/repository)  |
|   - Database communication           |
|   - Transaction management           |
|   - Implements Querier               |
+--------------------------------------+
                |
                v
+--------------------------------------+
|       Database (internal/db)         |
|   - DB connection (db.go)            |
|   - Querier interface                |
|   - sqlc-generated code              |
+--------------------------------------+
|       SQL (db/query, db/schema)      |
+--------------------------------------+

```

### Request Lifecycle

1.  **Entrypoint (`cmd/server/main.go`):** Initializes the database connection pool, configuration, and all dependencies. It starts the HTTP server.
2.  **Router (`internal/transport/http/router`):** Receives the HTTP request and directs it to the appropriate handler based on the URL path.
3.  **Handler (`internal/transport/http/handler`):** Decodes the request body into a DTO (`internal/dto`), validates the input, and calls the appropriate method on the `service` layer.
4.  **Service (`internal/service`):** Executes the core business logic. For data persistence, it calls methods on its `repository` dependency.
5.  **Repository (`internal/repository`):** Receives the call from the service and executes the necessary database operations using the `sqlc`-generated queries. This is the layer where all SQL queries and transactions are managed.
6.  **Response:** The data flows back up the chain. The repository returns database models to the service. The service returns DTOs to the handler. The handler encodes the DTOs into a JSON response and sends it back to the client.

---

## Getting Started

### 1. Clone and configure

```bash
git clone https://github.com/chatzijohn/amr-data-bridge.git
cd amr-data-bridge
cp .env.example .env
# Edit .env with your database connection details
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
- sqlc for SQL generation
- PostgreSQL or a compatible datastore