# Tech Stack

## Cloud Infrastructure

**Provider:** None (Local-only application)

**Key Services:** N/A - No cloud services in Technical Demo phase

**Deployment Regions:** N/A - Runs locally on developer's macOS machine

**Future Consideration (Epic 2+):** If persistence/sync needs arise, consider iCloud for Apple Notes integration or local SQLite for enrichment data.

## Technology Stack Table

| Category | Technology | Version | Purpose | Rationale |
|----------|------------|---------|---------|-----------|
| **Language** | Go | 1.25.4 | Primary development language | Excellent CLI tooling, fast compilation, simple deployment |
| **Runtime** | Go Runtime | 1.25.4 | Execution environment | Bundled with language; single binary output |
| **TUI Framework** | Bubbletea | 1.2.4 | Terminal UI framework | Elm architecture for reactive UIs, perfect for multi-view navigation |
| **TUI Styling** | Lipgloss | 1.0.0 | Terminal styling and layout | ANSI color support, box rendering for doors and detail views |
| **TUI Components** | Bubbles | 0.20.0 | Pre-built TUI components | Text input for notes, list selection for status menu |
| **Terminal Utilities** | golang.org/x/term | 0.26.0 | Terminal size detection & control | Responsive layout across terminal sizes |
| **YAML Parser** | gopkg.in/yaml.v3 | 3.0.1 | Structured data parsing | Human-readable task storage with metadata (status, notes, timestamps) |
| **UUID Generator** | github.com/google/uuid | 1.6.0 | Unique task IDs | UUID v4 generation for task identity |
| **Formatting** | gofumpt | 0.7.0 | Code formatting | Enforces consistent style |
| **Linting** | golangci-lint | 1.61.0 | Static analysis | Catches bugs early |
| **Build System** | Make | System default (macOS) | Build automation | `build`, `run`, `clean`, `lint`, `fmt` targets |
| **Dependency Mgmt** | Go Modules | 1.25.4 | Package management | `go.mod` versioning |
| **Testing** | Go testing package | 1.25.4 (stdlib) | Unit testing | Built-in; no external framework |
| **Storage** | YAML Files | N/A | Task persistence with metadata | `~/.threedoors/tasks.yaml` with status, notes, timestamps |
| **Platform** | macOS | 14+ (Sonoma+) | Target OS | Developer's primary platform |
| **Terminal** | iTerm2 / Terminal.app | Latest | Terminal emulator | 256-color support |
| **Version Control** | Git | 2.40+ | Source control | github.com/arcaven/ThreeDoors.git |

## Makefile Targets

```makefile
.PHONY: build run clean fmt lint test

build:
	go build -o bin/threedoors cmd/threedoors/main.go

run: build
	./bin/threedoors

clean:
	rm -rf bin/

fmt:
	gofumpt -l -w .

lint:
	golangci-lint run ./...

test:
	go test -v ./...
```

---
