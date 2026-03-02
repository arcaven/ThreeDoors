# Sprint Change Proposal - CI/CD Pipeline & Alpha Release

**Date:** 2026-03-02
**Proposed by:** PM (Course Correction Workflow)
**Change Scope:** Minor - Direct implementation by dev team

---

## Section 1: Issue Summary

**Problem Statement:** The ThreeDoors project is in active multi-agent development (3 stories merged via PRs #1, #2, #4, #5) but has **zero CI/CD infrastructure**. There are no automated quality gates, no build validation on PRs, and no release artifact generation. The Architecture document explicitly deferred CI/CD to Epic 2, but operational reality requires it now.

**Discovery Context:** With multiple agents submitting PRs simultaneously, there's no automated enforcement of:
- Code formatting (`gofumpt`)
- Linting (`golangci-lint`)
- Test execution (`go test`)
- Build validation (`go build`)
- Binary packaging for releases

**Evidence:**
- `docs/architecture/infrastructure-and-deployment.md` states: "CI/CD Platform: None for Technical Demo (deferred to Epic 2)"
- No `.github/workflows/` directory exists in the repository
- Makefile has `fmt`, `lint`, `test`, `build` targets but they're only run manually
- 5 PRs have been merged with no automated quality checks

---

## Section 2: Impact Analysis

### Epic Impact
- **Epic 1 (Three Doors Technical Demo):** Add new Story 1.7 for CI/CD pipeline. No changes to existing stories 1.1-1.6.
- **Epics 2-5:** All benefit from CI/CD foundation. No negative impact.

### Story Impact
- **No existing stories require changes**
- **New story needed:** Story 1.7: CI/CD Pipeline & Alpha Release

### Artifact Conflicts
- **Architecture - Infrastructure & Deployment:** Must be updated to document CI/CD platform (GitHub Actions)
- **Architecture - Tech Stack:** Must add GitHub Actions to the technology stack table
- **PRD - Epic Details:** Must add Story 1.7 to Epic 1 story list
- **PRD - Epic List:** Must update Epic 1 deliverables to include CI/CD

### Technical Impact
- New `.github/workflows/ci.yml` file
- Possible Makefile updates for CI-specific targets
- No code changes required to existing Go source

---

## Section 3: Recommended Approach

**Selected Path:** Direct Adjustment - Add new story within existing Epic 1

**Rationale:**
- CI/CD is purely additive infrastructure - no existing work needs modification
- Low effort (1-2 hours) with immediate quality improvement
- Low risk - GitHub Actions is well-documented, Go CI is straightforward
- High value - enables quality gates for ongoing multi-agent development
- No timeline impact - can be developed in parallel with other stories

**Effort Estimate:** Low (1-2 hours)
**Risk Level:** Low
**Timeline Impact:** None - parallel work

---

## Section 4: Detailed Change Proposals

### New Story: Story 1.7 - CI/CD Pipeline & Alpha Release

**As a** developer working in a multi-agent environment,
**I want** automated CI/CD quality gates and alpha release packaging,
**so that** every PR is validated for formatting, linting, tests, and build, and merged PRs produce downloadable alpha binaries.

**Acceptance Criteria:**

1. **GitHub Actions CI workflow** (`.github/workflows/ci.yml`) created
2. **PR Quality Gates** (runs on every PR to main):
   - `gofumpt` formatting check (fail if unformatted)
   - `golangci-lint` static analysis (fail on warnings)
   - `go vet ./...` correctness check
   - `go test ./...` unit tests
   - `go build ./cmd/threedoors` build validation
3. **All jobs run on GitHub-hosted public runners** (`ubuntu-latest`)
4. **Alpha Release packaging** (runs on PR merge to main):
   - Build binary for darwin/arm64 (primary target) and darwin/amd64
   - Build binary for linux/amd64 (CI runner compatibility)
   - Create GitHub Release with tag `alpha-<short-sha>` or use artifacts
   - Upload compiled binaries as release assets
5. **Go version** pinned to match `go.mod` (1.25.4)
6. **Caching** enabled for Go modules and build cache
7. **Branch protection** recommendation documented (require CI to pass before merge)

**Estimated Time:** 60-90 minutes

### Architecture Document Updates

**File:** `docs/architecture/infrastructure-and-deployment.md`

**OLD:**
```
## Infrastructure as Code
**Tool:** Not Applicable (local execution)
**Approach:** ThreeDoors Technical Demo runs locally with no cloud infrastructure.

**CI/CD Platform:** None for Technical Demo (deferred to Epic 2)
```

**NEW:**
```
## Infrastructure as Code
**Tool:** GitHub Actions
**Approach:** ThreeDoors uses GitHub Actions for CI/CD on public runners.

**CI/CD Platform:** GitHub Actions
- PR gates: fmt, lint, vet, test, build
- Alpha release: binary packaging on merge to main
```

**File:** `docs/architecture/tech-stack.md`

**ADD to Technology Stack Table:**
```
| **CI/CD** | GitHub Actions | N/A | Continuous integration & release | Public runners, native Go support |
```

### PRD Document Updates

**File:** `docs/prd/epic-list.md`

**Epic 1 Stories list - ADD:**
```
- Story 1.7: CI/CD Pipeline & Alpha Release
```

**File:** `docs/prd/epic-details.md`

**ADD:** New Story 1.7 section (full acceptance criteria as above)

---

## Section 5: Implementation Handoff

**Change Scope Classification:** Minor - Direct implementation by development team

**Handoff:**
- **SM Agent:** Schedule Story 1.7 immediately for current sprint
- **Dev Agent:** Implement the CI/CD pipeline
- **TEA Agent:** Validate test execution in CI environment

**Success Criteria:**
- CI workflow runs successfully on PR creation
- Failed formatting/lint/tests block PR merge
- Merged PRs produce downloadable alpha binaries
- No regressions in existing functionality

---

## Approval

**Status:** Pending user approval
