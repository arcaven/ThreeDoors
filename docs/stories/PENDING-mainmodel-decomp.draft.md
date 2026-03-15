# Epic PENDING: TUI MainModel Decomposition

## Status: DRAFT — Awaiting epic number allocation from project-watchdog

## Priority: P1

## Goal

Break `internal/tui/main_model.go` (2991 lines, 32 ViewModes) into focused view controller files. Pure refactoring — zero user-visible changes. Reduces merge conflict surface, improves code navigation, and lowers bug risk for all future TUI work.

## Background

`main_model.go` is the largest file in the codebase at 2991 lines. It contains:
- 32 `ViewMode` constants and their string representations
- The `MainModel` struct with ~50+ fields
- `Init()`, `Update()`, and `View()` methods that handle all 32 views
- Command palette / search dispatch
- View transition logic
- Message routing for all view types

This is a classic "god object" anti-pattern. Every TUI PR risks merge conflicts because all changes touch this one file. The file is difficult to navigate and reason about.

## Approach

Extract logical groups of view handling into focused files while keeping `MainModel` as the central coordinator. The struct itself stays in `main_model.go`, but the `Update()` and `View()` switch arms move to dedicated files that operate on `*MainModel` via methods.

**Key constraint:** No user-visible behavior changes. Every story's AC is "before and after are identical from the user's perspective."

## Proposed Stories

### Story X.1: Extract View Transition & Navigation Logic

Extract `setViewMode()`, `goBack()`, view stack management, and the `previousView` tracking into `view_navigation.go`. This is the foundation — other stories depend on clean navigation extraction.

**Acceptance Criteria:**
- New file `internal/tui/view_navigation.go` contains all view transition logic
- `main_model.go` reduced by ~200-300 lines
- All existing tests pass unchanged
- `go test -race ./internal/tui/...` passes
- No user-visible behavior changes (golden tests unchanged)

### Story X.2: Extract Source/Sync View Controllers

Move `Update()` and `View()` handling for `ViewSources`, `ViewSourceDetail`, `ViewSyncLog`, `ViewSyncLogDetail`, `ViewConnectWizard`, `ViewDisconnect`, `ViewReauth` into `view_sources_controller.go`.

**Acceptance Criteria:**
- New file `internal/tui/view_sources_controller.go` handles 7 source-related views
- `main_model.go` reduced by ~400-500 lines
- All existing tests pass unchanged
- `go test -race ./internal/tui/...` passes

### Story X.3: Extract Planning & Task Management View Controllers

Move handling for `ViewPlanning`, `ViewAddTask`, `ViewBreakdown`, `ViewExtract`, `ViewImport`, `ViewSnooze`, `ViewDeferred` into `view_task_controller.go`.

**Acceptance Criteria:**
- New file `internal/tui/view_task_controller.go` handles 7 task-related views
- `main_model.go` reduced by ~300-400 lines
- All existing tests pass unchanged
- `go test -race ./internal/tui/...` passes

### Story X.4: Extract Auxiliary View Controllers & Command Dispatch

Move handling for `ViewHelp`, `ViewBugReport`, `ViewThemePicker`, `ViewHealth`, `ViewInsights`, `ViewMood`, `ViewFeedback`, `ViewValuesGoals`, `ViewOrphaned`, `ViewConflict`, `ViewProposals`, `ViewDevQueue` into `view_auxiliary_controller.go`. Also extract command palette dispatch (`:help`, `:bug`, `:sources`, etc.) into `command_dispatch.go`.

**Acceptance Criteria:**
- New files `internal/tui/view_auxiliary_controller.go` and `internal/tui/command_dispatch.go`
- `main_model.go` reduced to ~800-1000 lines (struct definition, Init, top-level Update/View dispatch)
- All existing tests pass unchanged
- `go test -race ./internal/tui/...` passes
- Golden snapshot tests produce identical output

## Dependency Graph

X.1 → X.2 (can parallelize with X.3)
X.1 → X.3 (can parallelize with X.2)
X.2 + X.3 → X.4

## Risk Assessment

- **Low risk:** Go methods on the same type can live in different files — this is idiomatic
- **Medium risk:** Some Update() arms may have cross-view dependencies (e.g., a source view triggers a navigation change). These need careful extraction.
- **Mitigation:** Golden snapshot tests + race detector provide strong regression safety net

## Notes

- This epic is pure tech debt reduction — no new features
- Aligns with CLAUDE.md "One primary type per file" guideline
- After this epic, future TUI stories will have smaller, focused diffs
