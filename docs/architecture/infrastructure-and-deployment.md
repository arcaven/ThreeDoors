# Infrastructure and Deployment

## Infrastructure as Code

**Tool:** Not Applicable (local execution)

**Approach:** ThreeDoors Technical Demo runs locally with no cloud infrastructure.

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

**CI/CD Platform:** None for Technical Demo (deferred to Epic 2)

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
