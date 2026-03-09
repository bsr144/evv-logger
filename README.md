# EVV Logger

Electronic Visit Verification (EVV) Logger API built with Go. Manages caregiver schedules and task tracking for home healthcare visits.

## Prerequisites

- **Go 1.24+** - [Install Go](https://go.dev/doc/install)
- **Docker & Docker Compose** (optional, for containerized deployment)

## Quick Start

### Run Locally

```bash
# Clone the repository
git clone https://github.com/bsr144/evv-logger.git
cd evv-logger

# Download dependencies
go mod download

# Run the server (starts on port 8080 with seed data)
go run ./cmd/http

# Or build and run the binary
go build -o evv-logger ./cmd/http
./evv-logger
```

The server starts on port `8080` by default. Set the `PORT` environment variable to change it:

```bash
PORT=3000 go run ./cmd/http
```

### Run with Docker

```bash
cd deployment/docker
docker compose up --build
```

The container exposes port `8080`. To customize, create a `.env` file in the project root (see `deployment/docker/.env.example`).

## Tech Stack & Key Decisions

| Technology | Why |
|------------|-----|
| **Go + Fiber v2** | Lightweight, high-performance HTTP framework well-suited for microservices. Fiber's Express-like API makes route handling straightforward, and its built-in `app.Test()` method simplifies integration testing without starting a real server. |
| **In-memory store with `sync.RWMutex`** | Keeps the project dependency-free -- no database driver, no migrations, no external process to manage. `RWMutex` provides thread-safe concurrent access with readers never blocking each other, which is sufficient for the assignment scope. |
| **Clean Architecture** | Enforces separation of concerns across three layers: `controller -> usecase -> repository`. Each layer depends only on the interface of the layer below it, making every component independently testable and replaceable. |
| **Manual DI in `main.go`** | All wiring is explicit in a single file. There is no framework magic, no reflection, and no hidden initialization order. This makes the dependency graph easy to read and debug. |
| **validator v10** | Provides declarative struct-level validation via field tags (`validate:"required"`). This keeps validation logic co-located with the DTO definitions and avoids hand-written checks. |
| **testify** | Idiomatic Go test assertions (`assert`, `require`) that produce clear failure messages. Widely adopted in the Go ecosystem with minimal learning curve. |

## Dependencies

| Package | Purpose |
|---------|---------|
| [Fiber v2](https://github.com/gofiber/fiber) | HTTP framework |
| [validator v10](https://github.com/go-playground/validator) | Request validation |
| [testify](https://github.com/stretchr/testify) | Test assertions |

No external services (databases, caches, etc.) are required. All data is stored in-memory using Go maps with `sync.RWMutex` for thread safety. Data resets on restart.

## Running Tests

```bash
# Run all tests
go test ./...

# Verbose output
go test -v ./...

# With coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out    # open HTML report in browser

# Run with race detector
go test -race ./...

# Run a single test
go test -run TestCreateSchedule ./...

# Run a specific package
go test -v ./tests/integration/...
go test -v ./tests/unit/...
```

### Test Structure

```
tests/
  integration/     # API endpoint tests using Fiber's app.Test()
    schedule_api_test.go
    task_api_test.go
  unit/            # Unit tests for store and usecases
    memory_store_test.go
    schedule_usecase_test.go
    task_usecase_test.go
```

## Linting

```bash
go vet ./...
golangci-lint run ./...   # requires golangci-lint installed
```

## API Reference

Base URL: `http://localhost:8080/api/v1`

### Health Check

```
GET /api/v1/health
```

Response: `{"status": "ok"}`

### Schedules

#### Create Schedule

```
POST /api/v1/schedules
Content-Type: application/json

{
  "caregiver_name": "Alice Johnson",
  "patient_name": "Robert Smith",
  "date": "2026-03-08",
  "start_time": "08:00",
  "end_time": "12:00",
  "location": "123 Main St"
}
```

#### List Schedules

```
GET /api/v1/schedules
GET /api/v1/schedules?date=2026-03-08
```

#### Get Schedule by ID

```
GET /api/v1/schedules/:id
```

#### Update Schedule (Clock In/Out, Status)

```
PATCH /api/v1/schedules/:id
Content-Type: application/json

{
  "status": "in_progress",
  "clock_in_time": "2026-03-08T08:05:00Z",
  "clock_in_lat": 40.7128,
  "clock_in_lng": -74.0060
}
```

Schedule statuses: `upcoming`, `in_progress`, `completed`, `missed`

### Tasks

#### Create Task

```
POST /api/v1/schedules/:id/tasks
Content-Type: application/json

{
  "description": "Check vital signs"
}
```

#### List Tasks for Schedule

```
GET /api/v1/schedules/:id/tasks
```

#### Update Task

```
PATCH /api/v1/schedules/:id/tasks/:taskId
Content-Type: application/json

{
  "status": "completed",
  "notes": "BP 120/80, all normal"
}
```

Task statuses: `pending`, `completed`, `not_completed`

> Note: `notes` is required when setting status to `not_completed`.

### Response Format

All endpoints return a consistent JSON envelope:

```json
// Success
{
  "success": true,
  "data": { ... }
}

// Error
{
  "success": false,
  "error": "error message"
}
```

## Project Structure

```
cmd/http/main.go                    # Entry point, DI wiring, seed data
entity/                             # Domain models (Schedule, Task)
evv/schedule/                       # Schedule usecase (interface, logic, request/response)
evv/task/                           # Task usecase
delivery/http/controller/           # HTTP controllers (Fiber handlers)
delivery/http/dto/                  # Shared response types
delivery/http/router/               # Route registration
delivery/http/middleware/           # CORS, request logging
impl/memory/                        # In-memory store (sync.RWMutex)
pkg/errors/                         # Shared error types
tests/                              # Integration + unit tests
deployment/docker/                  # Dockerfile and docker-compose
```

## Assumptions

- **Single caregiver user.** The assignment specifies "no auth needed," so the API assumes a single caregiver and does not implement authentication or authorization.
- **Data resets on server restart.** All state lives in-memory by design. This is intentional for the assignment scope; a production system would use persistent storage.
- **Schedule dates and times are strings.** Dates use ISO format (`2006-01-02`) and times use `HH:MM` format. This simplifies parsing and avoids timezone ambiguity for a demo application.
- **Geolocation is captured by the frontend.** The backend accepts `clock_in_lat`, `clock_in_lng`, `clock_out_lat`, and `clock_out_lng` as optional `float64` fields in the PATCH request. It is the frontend's responsibility to obtain these values via the browser Geolocation API.
- **Task status `not_completed` requires notes.** When a caregiver marks a task as not completed, a reason must be provided. This is enforced server-side with a validation check.
- **Auto-increment integer IDs.** All entities use simple auto-increment integers as primary keys. UUIDs are unnecessary for a single-instance in-memory store.
- **Seed data is loaded on startup.** The server pre-populates sample schedules and tasks so the application is immediately usable for demo and review purposes.

## Geolocation Handling

Per the assignment requirement, geolocation is mandatory for EVV compliance. The approach is split between frontend and backend responsibilities:

**Backend (this repository):**
- The `PATCH /api/v1/schedules/:id` endpoint accepts `clock_in_lat`, `clock_in_lng`, `clock_out_lat`, and `clock_out_lng` as optional `float64` fields.
- These fields are stored alongside the clock-in/clock-out timestamps on the schedule record.
- The backend does not enforce that coordinates must be present -- it treats them as optional to support graceful degradation.

**Frontend (expected behavior):**
- On "Start Visit" and "End Visit" actions, the frontend calls `navigator.geolocation.getCurrentPosition()` to obtain the device's latitude and longitude.
- The obtained coordinates are included in the PATCH request body sent to the backend.

**Fallback handling:**
- If the user denies geolocation permission or the browser does not support the Geolocation API, the frontend should still allow the clock-in/clock-out to proceed.
- In this case, `clock_in_lat` and `clock_in_lng` (or their clock-out equivalents) are sent as `null` or omitted entirely. The backend accepts the request without coordinates.
- This graceful degradation ensures the core EVV workflow -- capturing timestamps and updating visit status -- is never blocked by geolocation availability.

## Architecture

Clean Architecture with dependency flow: `controller -> usecase -> repository (impl)`.

All wiring is done manually in `cmd/http/main.go` - no DI framework. Controllers depend on usecase interfaces, usecases depend on repository interfaces, and only `main.go` imports concrete implementations.

## What I Would Improve

Given more time, the following enhancements would strengthen the application for production use:

- **Persistent storage (PostgreSQL).** Replace the in-memory store with a relational database to survive restarts and support multi-instance deployments.
- **Authentication and authorization (JWT).** Add token-based auth so multiple caregivers can securely access only their own schedules and tasks.
- **WebSocket for real-time updates.** Push schedule and task status changes to connected clients without polling, improving the caregiver experience during active visits.
- **Pagination for list endpoints.** The current `GET /schedules` and `GET /tasks` endpoints return all records. Cursor-based or offset pagination would be necessary at scale.
- **Request rate limiting.** Protect the API from abuse with per-IP or per-token rate limits using middleware.
- **Structured audit logging for EVV compliance.** Log every clock-in, clock-out, and task status change with timestamps, coordinates, and user identity to an append-only audit trail.
- **CI/CD pipeline.** Automate linting, testing, race detection, and Docker image builds on every push using GitHub Actions or a similar platform.
