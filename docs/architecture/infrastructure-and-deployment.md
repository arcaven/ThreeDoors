# Infrastructure and Deployment

## Infrastructure as Code

**Tool:** GitHub Actions

**Approach:** ThreeDoors uses GitHub Actions for CI/CD on public runners. The application runs locally with no cloud infrastructure for runtime.

## Deployment Strategy

**Strategy:** Direct Binary Distribution

**Build Process:**
```bash
make build    # Compiles to bin/threedoors
```

**Installation:**
```bash
# Option 1: Manual install
cp bin/threedoors /usr/local/bin/

# Option 2: Run from project directory
make run

# Option 3 (Future): Homebrew tap
brew install arcaven/tap/threedoors
```

**CI/CD Platform:** GitHub Actions (`.github/workflows/ci.yml`)

**PR Quality Gates** (runs on every `pull_request` to `main`):
- `gofumpt` formatting check
- `go vet` correctness check
- `golangci-lint` static analysis
- `go test` unit tests
- `go build` build validation

**Alpha Release** (runs on `push` to `main`, i.e., PR merge):
- Cross-compiles binaries for darwin/arm64, darwin/amd64, linux/amd64
- Uploads as GitHub Actions workflow artifacts (14-day retention)

**Recommended:** Enable branch protection on `main` requiring the quality-gate job to pass before merge.

## Environments

**Development:**
- Purpose: Local development and testing
- Location: Developer's macOS machine
- Data: `~/.threedoors/` (can be deleted/reset)

**Production (User Environment):**
- Purpose: End-user execution
- Location: User's macOS machine
- Data: `~/.threedoors/` (user's actual task data)

## Rollback Strategy

**Primary Method:** User keeps previous binary

**Rollback Process:**
```bash
# User manually switches to previous version
cp threedoors.old /usr/local/bin/threedoors
```

**Data Compatibility:**
- YAML schema must remain backward compatible
- Forward migrations add fields with defaults
- Never break existing tasks.yaml format

---
