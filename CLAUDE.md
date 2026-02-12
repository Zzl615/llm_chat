# CLAUDE.md

This file provides guidance to AI assistants working with the `llm_chat` codebase.

## Project Overview

A standalone, high-concurrency LLM chat demo built in Go. It uses Gin + gorilla/websocket to accept WebSocket connections, manages per-connection sessions with goroutines/channels, and simulates an LLM streaming response pipeline through an in-memory mock message queue.

There is no database, no persistent storage, and no real LLM integration — all inference is simulated.

## Repository Structure

```
llm_chat/
├── main.go                  # Entry point: wires Manager, MockQueue, routes; starts server on :8080
├── go.mod / go.sum          # Go module (module name: llm-chat, Go 1.23.4)
├── internal/                # Private internal packages
│   ├── manager.go           # Manager: thread-safe session registry (sync.RWMutex)
│   ├── session.go           # Session: per-connection struct with ReadPump/WritePump goroutines
│   ├── mock_queue.go        # MockQueue: channel-based request/result queue simulating Kafka/Redis
│   └── ws_handler.go        # RegisterRoutes: Gin route handler that upgrades HTTP to WebSocket
├── test/                    # Unit tests (separate package)
│   └── mock_queue_test.go   # Tests for MockQueue (5 test functions)
├── README.md                # Project docs (Chinese)
└── .gitignore               # Standard Go ignores + IDE files + compiled binary
```

## Key Architecture

### Core Components

1. **Manager** (`internal/manager.go`) — Thread-safe map of session ID to `*Session`. Methods: `Register`, `Unregister`, `Get`. Protected by `sync.RWMutex`.

2. **Session** (`internal/session.go`) — Represents one WebSocket connection. Contains a `Send` channel (buffered at 128) and a `CloseCh`. Runs two goroutines:
   - `ReadPump`: reads client messages, forwards them via a callback to the queue
   - `WritePump`: drains the `Send` channel and writes to the WebSocket; sends periodic pings

3. **MockQueue** (`internal/mock_queue.go`) — Simulates a Kafka/Redis message queue using Go channels (both buffered at 1024). `StartMockModelWorker` consumes requests and produces 5 streaming chunks per request (400ms apart). `SubscribeResults` delivers results to a callback.

4. **WebSocket Handler** (`internal/ws_handler.go`) — Registers `GET /ws` via Gin. Upgrades the connection, creates a Session, registers it, and spawns read/write goroutines.

### Data Flow

```
Client → WebSocket /ws → Session.ReadPump → MockQueue.PublishRequest
                                                    ↓
                                          MockQueue worker (goroutine)
                                                    ↓
                                          MockQueue.SubscribeResults → callback
                                                    ↓
                                          Session.Send channel → Session.WritePump → Client
```

### Concurrency Model

- Each WebSocket connection gets 2 goroutines (read + write pump)
- All inter-component communication uses Go channels (no shared mutable state except Manager's map behind RWMutex)
- Non-blocking sends to `Session.Send` with overflow drop + warning log

## Build & Run Commands

```bash
# Install dependencies
go mod download

# Run the server (starts on :8080)
go run .

# Build binary
go build -o llm_chat .

# Run all tests
go test ./test -v

# Run a specific test
go test ./test -v -run TestStartMockModelWorker

# Check test coverage
go test ./test -cover
```

## Dependencies

Only 2 direct dependencies:
- `github.com/gin-gonic/gin v1.10.0` — HTTP web framework
- `github.com/gorilla/websocket v1.5.1` — WebSocket protocol support

## API

Single endpoint:
- **`GET /ws`** — WebSocket upgrade. Send plain text messages, receive streaming chunks in format `chunk N: {content}` (5 chunks per message, 400ms apart).

## Code Conventions

- **Go version**: 1.23.4
- **Module name**: `llm-chat`
- **Package layout**: All non-main code lives in `internal/` (Go's enforced private package convention). Tests live in `test/` as a separate package.
- **File headers**: JSDoc-style comment blocks with `@Author`, `@Date`, `@Last Modified by/time`
- **Naming**: Standard Go conventions — exported names are PascalCase, unexported are camelCase
- **Comments**: Inline comments in Chinese describing component purpose
- **Error handling**: Log errors via `log.Printf` with bracketed context tags (e.g., `[session %s]`, `[MockModel]`, `[WARN]`)
- **Constants**: WebSocket timing constants (`writeWait`, `pongWait`, `pingPeriod`) are package-level in `session.go`
- **Channel buffer sizes**: Session.Send=128, MockQueue requests/results=1024

## Testing Conventions

- Tests are in `test/` directory (package `test`), not alongside source files
- Uses Go's standard `testing` package only (no third-party test frameworks)
- Async assertions use channels + `select` with `time.After` timeouts
- Concurrency tests use `sync.Mutex` for shared result collection
- Test function naming: `TestDescriptiveName` (e.g., `TestConcurrentPublish`)
- Run tests with `go test ./test -v`

## Known Limitations / Notes

- `ws_handler.go:28` generates session IDs by reading `mgr.sessions` map length without holding the lock — this is a known race condition in the demo code
- `CheckOrigin` returns `true` for all origins (no CORS restriction)
- No graceful shutdown handling
- MockQueue processes requests sequentially (single worker goroutine), so concurrent requests queue up
- No CI/CD, no Dockerfile, no Makefile configured

## Planned Future Work (from README)

- RESTful API endpoints
- Real Kafka queue integration
- Real LLM model integration
- Rate limiting and circuit breaking
- Message persistence
- Prometheus metrics
- Distributed deployment support
