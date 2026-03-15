# Epic PENDING: Completion History & Progress View

## Status: DRAFT — Awaiting epic number allocation from project-watchdog

## Priority: P1

## Goal

Let users see what they've accomplished. New `:history` TUI view and `threedoors history` CLI command showing completed tasks with daily/weekly grouping. Aligns with SOUL.md "Progress Over Perfection" — positive reinforcement without gamification or guilt.

## Background

ThreeDoors already tracks completions in two places:
- `completed.txt` — append-only log of completed task titles with timestamps
- JSONL session logs — `task_completed` events with full metadata (task ID, title, time, session context)

But there's no way to VIEW this data from inside the app. Users who want to see "what did I do today?" or "what did I accomplish this week?" have no in-app answer. The data exists — it just needs a view.

This directly supports the SOUL.md philosophy:
> "Opening ThreeDoors should feel like a friend saying: 'Hey, here are three things you could do right now. Pick one. Any one. Let's go.'"

The history view adds: "And look — here's what you already did. Nice work."

**What this is NOT:**
- Not a habit tracker (no streaks, no scores, no "you missed a day")
- Not a productivity report (no charts, no metrics, no comparisons)
- Not gamification (no badges, no achievements, no leaderboards)

It's a simple, warm list of things you did. That's it.

## Proposed Stories

### Story Y.1: Completion Data Reader & Aggregator

Create a `CompletionReader` that reads from `completed.txt` and/or JSONL session logs, returning a unified `[]CompletionRecord` sorted by time. Support filtering by date range (today, this week, this month, all time). This is the data layer — no UI.

**Acceptance Criteria:**
- New type `CompletionRecord` with fields: Title, CompletedAt, Source (which adapter), TaskID (optional)
- `CompletionReader` reads from `completed.txt` with fallback to JSONL session logs
- Filter methods: `Today()`, `ThisWeek()`, `ThisMonth()`, `All()` returning `[]CompletionRecord`
- Records sorted newest-first by default
- Handles empty/missing files gracefully (returns empty slice, no error)
- Table-driven tests with test fixtures in `testdata/`
- 80%+ test coverage

### Story Y.2: History TUI View (`:history`)

New `ViewHistory` mode accessible via `:history` command. Shows completed tasks grouped by day, scrollable, with the same keybinding patterns as other list views (j/k scroll, q/Esc back, ? help).

**Acceptance Criteria:**
- New `ViewHistory` view mode added to `main_model.go` constants
- `:history` command registered in command dispatch
- View shows completed tasks grouped by date headers (e.g., "Today — March 15", "Yesterday — March 14")
- Each entry shows: task title, completion time (HH:MM), source badge if from an adapter
- Scrollable with j/k or arrow keys
- `q` or `Esc` returns to previous view
- Empty state: friendly message ("No completed tasks yet. Pick a door and get started!")
- Lipgloss styling consistent with existing views (help, sources, insights)
- Keybinding bar shows relevant shortcuts
- `go test -race ./internal/tui/...` passes
- Golden snapshot test for the view

### Story Y.3: History CLI Command (`threedoors history`)

New `threedoors history` CLI command with `--today`, `--week`, `--month`, `--all` flags and `--json` output support.

**Acceptance Criteria:**
- `threedoors history` shows today's completions by default
- `--today`, `--week`, `--month`, `--all` flags for date range filtering
- `--json` flag outputs structured JSON (array of CompletionRecord objects)
- Human-readable output: grouped by date, task title + time
- Exit code 0 on success, 1 on error
- Help text via `threedoors history --help`
- Tests covering all flag combinations
- Consistent with existing CLI patterns (uses Cobra, follows `internal/cli/` conventions)

## Dependency Graph

Y.1 → Y.2 (TUI depends on data layer)
Y.1 → Y.3 (CLI depends on data layer)
Y.2 and Y.3 can parallelize after Y.1

## Design Decisions

- **Data source:** Read from `completed.txt` first (simpler, always exists). Fall back to JSONL only if completed.txt is empty or missing.
- **No pagination:** Scroll-based, not page-based. Completed tasks are bounded by time filters.
- **No deletion:** History is read-only. Users cannot delete history entries (append-only audit trail).
- **No search:** Keep it simple. If search is needed, the existing `:search` command could be extended later.
- **Grouping:** By calendar day in user's local timezone, not UTC. Display dates use the user's locale.

## SOUL.md Alignment

- "Progress Over Perfection" — seeing what you've done reinforces that imperfect action > perfect planning
- "Work With Human Nature" — progress visibility is a proven psychological motivator
- "Not a habit tracker" — no streaks, no guilt, just a warm list
- "Every Interaction Should Feel Deliberate" — the view should feel like opening a journal, not a report
