# ThreeDoors User Guide

Three doors, one choice. Move forward.

ThreeDoors is a terminal-based task manager that reduces decision friction by showing you only three tasks at a time. Instead of staring at a long backlog, you pick a door and take action.

---

## Table of Contents

- [Getting Started](#getting-started)
- [Core Concepts](#core-concepts)
- [Basic Usage](#basic-usage)
- [Task Management](#task-management)
- [Search and Commands](#search-and-commands)
- [Task Sources](#task-sources)
- [Obsidian Integration](#obsidian-integration)
- [Apple Notes Integration](#apple-notes-integration)
- [Intelligent Features](#intelligent-features)
- [Offline and Sync](#offline-and-sync)
- [Session Metrics](#session-metrics)
- [Configuration](#configuration)
- [Troubleshooting](#troubleshooting)

---

## Getting Started

### Installation

**Homebrew (macOS — recommended):**

```bash
brew install arcaven/tap/threedoors
```

**Pre-built binaries:**

Download from [GitHub Releases](https://github.com/arcaven/ThreeDoors/releases):

```bash
# macOS Apple Silicon
chmod +x threedoors-darwin-arm64
mv threedoors-darwin-arm64 /usr/local/bin/threedoors

# macOS Intel
chmod +x threedoors-darwin-amd64
mv threedoors-darwin-amd64 /usr/local/bin/threedoors

# Linux x86_64
chmod +x threedoors-linux-amd64
mv threedoors-linux-amd64 /usr/local/bin/threedoors
```

**Go install:**

```bash
go install github.com/arcaven/ThreeDoors/cmd/threedoors@latest
```

**Build from source (Go 1.25.4+):**

```bash
git clone https://github.com/arcaven/ThreeDoors.git
cd ThreeDoors
make build
# Binary at bin/threedoors
```

### First Launch

Run `threedoors` in your terminal. On first launch, the onboarding wizard walks you through five steps:

1. **Welcome** — Introduction to the Three Doors concept
2. **Key bindings** — Interactive tutorial where you press keys to learn the controls
3. **Values** — Set 1–5 personal values or goals (displayed as a reminder while you work)
4. **Import** — Optionally import tasks from CSV or Apple Notes
5. **Done** — You're ready to go

After onboarding, ThreeDoors creates its data directory at `~/.threedoors/` with a default text file provider.

### Version

```bash
threedoors --version
```

---

## Core Concepts

### The Three Doors Philosophy

Traditional task lists create choice paralysis. When you have 50+ tasks, picking one becomes its own task. ThreeDoors solves this by presenting exactly three options — pick one, take action, move on.

This is grounded in behavioral science:

- **Choice overload** (Iyengar & Lepper) — too many options reduces satisfaction and action
- **Cognitive capacity** (Cowan) — working memory holds 3–5 chunks; three is the sweet spot
- **Decision fatigue** (Baumeister) — fewer decisions preserves energy for actual work
- **Hick's Law** — response time scales with options; three minimizes decision latency

### Progress Over Perfection

ThreeDoors doesn't try to find the "optimal" task. It presents three reasonable options and trusts you to pick. Any forward motion beats standing still.

### How Door Selection Works

When you open ThreeDoors, three tasks are randomly selected from your pool. The selection algorithm:

- Excludes blocked, deferred, and archived tasks
- Avoids showing recently-displayed tasks
- Prefers diversity across task types and effort levels
- If you've logged a mood, biases toward tasks that correlate with productivity for that mood (based on your history)

Don't like your options? Press `s` or `↓` to refresh and get three new doors.

---

## Basic Usage

### Key Bindings — Doors View

| Key | Action |
|-----|--------|
| `a` / `←` | Select left door |
| `w` / `↑` | Select center door |
| `d` / `→` | Select right door |
| `enter` | Open the selected door (view task details) |
| `s` / `↓` | Refresh — get three new doors |
| `n` | Give feedback on a door (without opening it) |
| `m` | Log your current mood |
| `/` | Search tasks |
| `:` | Open command palette |
| `q` / `ctrl+c` | Quit |

### Key Bindings — Detail View

| Key | Action |
|-----|--------|
| `c` | Complete the task |
| `i` | Mark as in-progress |
| `b` | Mark as blocked (prompts for reason) |
| `p` | Procrastinate — return task to pool |
| `r` | Rework — return task to pool |
| `e` | Expand task (planned) |
| `f` | Fork task (planned) |
| `l` | Link to another task |
| `x` | Browse cross-references |
| `m` | Log mood |
| `esc` | Return to doors |

### Typical Workflow

1. Launch `threedoors`
2. Look at your three doors
3. Pick one (`a`, `w`, or `d`, then `enter`)
4. Take action: complete it (`c`), mark it in-progress (`i`), or return it to the pool (`p`)
5. Repeat

---

## Task Management

### Task Statuses

| Status | Description |
|--------|-------------|
| **Todo** | Default state for new tasks |
| **In-Progress** | Actively being worked on |
| **Blocked** | Cannot proceed; requires a blocker reason |
| **Complete** | Done (removed from active pool) |
| **Deferred** | Intentionally postponed |
| **Archived** | Removed from pool without completing |
| **In-Review** | Awaiting review (from in-progress) |

### Status Transitions

```
Todo ──→ In-Progress ──→ Complete
  │          │    │
  │          │    └──→ In-Review ──→ Complete
  │          │              │
  │          └──→ Blocked ──┘
  │
  ├──→ Blocked ──→ Todo / In-Progress / Complete
  ├──→ Complete (terminal)
  ├──→ Deferred ──→ Todo
  └──→ Archived (terminal)
```

### Action Keys in Detail View

- **`c` — Complete:** Marks the task done. It's removed from your active pool and logged to `completed.txt`. You'll see a celebration message with your daily completion count.
- **`i` — In-Progress:** Signals you're actively working on this task.
- **`b` — Blocked:** Prompts you to type a reason (e.g., "waiting on API key"). Press `enter` to confirm or `esc` to cancel.
- **`p` — Procrastinate:** Returns the task to the pool. No judgment — it'll come back around.
- **`r` — Rework:** Same as procrastinate, but signals the task needs more thought.

### Task Categorization

Press `:tag` in command mode to categorize a task across three dimensions:

**Type:**
- Creative (🎨) — design, writing, ideation
- Administrative (📋) — email, scheduling, paperwork
- Technical (🔧) — coding, debugging, system work
- Physical (💪) — exercise, errands, hands-on tasks

**Effort:**
- Quick Win — under 15 minutes
- Medium — 15–60 minutes
- Deep Work — focused, extended effort

**Location:**
- Home, Work, Errands, Anywhere

Categorization improves door selection over time — ThreeDoors learns which task types you prefer in different moods.

### Mood Tracking

Press `m` at any time to log how you're feeling:

1. Focused
2. Tired
3. Stressed
4. Energized
5. Distracted
6. Calm
7. Other (type your own)

Mood entries are timestamped and stored in your session data. Over time, ThreeDoors correlates mood with task completion patterns and adjusts door selection accordingly.

### Door Feedback

Press `n` on a door (without opening it) to give quick feedback:

1. **Blocked** — can't do this right now
2. **Not now** — not the right time
3. **Needs breakdown** — too big, needs splitting
4. **Other** — free-text comment

Feedback is recorded in session metrics for pattern analysis.

### Cross-References and Linking

Press `l` in detail view to link the current task to another task. Press `x` to browse existing links and navigate between related tasks. Links are stored in a local SQLite database (`enrichment.db`).

---

## Search and Commands

### Search Mode (`/`)

Press `/` to search tasks by text. Results filter in real-time as you type.

- `j` / `↓` — navigate down through results
- `k` / `↑` — navigate up through results
- `enter` — open selected task in detail view
- `esc` — exit search

After viewing a task from search, you return to search results (not the doors view) so you can continue browsing.

### Command Palette (`:`)

Press `:` to enter command mode. Available commands:

| Command | Description |
|---------|-------------|
| `:add [text]` | Add a new task. Without text, opens a prompt. |
| `:add-ctx [text]` | Add a task with context — two steps: task text, then why/context. |
| `:add --why [text]` | Same as `:add-ctx` — captures task and reason. |
| `:mood [mood]` | Log mood. Without argument, opens the mood dialog. |
| `:stats` | Show session metrics (tasks completed, duration, refreshes). |
| `:health` | Run system health checks on providers and data files. |
| `:goals` | View your values and goals. |
| `:goals edit` | Edit your values and goals. |
| `:tag` | Categorize the selected task (type, effort, location). |
| `:dashboard` | Open the insights dashboard. |
| `:insights` | Same as `:dashboard`. Accepts optional filter: `:insights mood` or `:insights avoidance`. |
| `:help` | Display all available commands. |
| `:quit` / `:exit` | Exit the application. |

### Adding Tasks

Three ways to add tasks:

**Quick add:**
```
:add Buy groceries
```

**Add with context (two-step):**
```
:add-ctx Refactor auth module
```
Then you're prompted: "Why is this important?" — your answer is stored as context.

**Inline tags:** Task text is parsed for inline categorization automatically.

---

## Task Sources

ThreeDoors supports multiple task storage backends ("providers"). Configure them in `~/.threedoors/config.yaml`.

### Text File (Default)

Tasks are stored as YAML in `~/.threedoors/tasks.yaml`. This is the default provider — no configuration needed.

```yaml
# Simple config.yaml
provider: textfile
```

### Apple Notes

Pull tasks from an Apple Notes note using checkbox syntax.

```yaml
provider: applenotes
note_title: ThreeDoors Tasks
```

See [Apple Notes Integration](#apple-notes-integration) for details.

### Obsidian

Pull tasks from Markdown files in an Obsidian vault.

```yaml
provider: obsidian
```

With settings:

```yaml
providers:
  - name: obsidian
    settings:
      vault_path: /path/to/your/vault
      tasks_folder: tasks
      file_pattern: "*.md"
      daily_notes: true
      daily_notes_folder: Daily
      daily_notes_heading: "## Tasks"
```

See [Obsidian Integration](#obsidian-integration) for details.

### Multiple Providers

You can configure multiple providers with fallback behavior:

```yaml
providers:
  - name: applenotes
    settings:
      note_title: ThreeDoors Tasks
  - name: textfile
    settings:
      task_file: ~/.threedoors/tasks.yaml
```

If the primary provider fails (e.g., Apple Notes is inaccessible), ThreeDoors falls back to the next provider automatically.

---

## Obsidian Integration

ThreeDoors integrates with Obsidian vaults, reading and writing tasks from Markdown files.

### Setup

1. Set your vault path in `config.yaml`:

```yaml
providers:
  - name: obsidian
    settings:
      vault_path: /Users/you/Documents/MyVault
      tasks_folder: tasks
      file_pattern: "*.md"
```

2. ThreeDoors scans the specified folder for Markdown files containing checkbox tasks.

### Checkbox Syntax

ThreeDoors recognizes standard Markdown checkboxes:

```markdown
- [ ] Uncompleted task (imported as Todo)
- [x] Completed task (imported as Complete)
- [/] In-progress task (imported as In-Progress)
```

Both `-` and `*` list markers are supported.

### Metadata Parsing

ThreeDoors extracts metadata from Obsidian-style annotations:

- **Due dates:** `📅 2026-03-15`
- **Priority:** `⏫` (high → Deep Work), `🔼` (medium → Medium), `🔽` (low → Quick Win)
- **Tags:** `#project` `#urgent` — stored in task context

### Task ID Tracking

ThreeDoors embeds unique IDs in HTML comments to track tasks across syncs:

```markdown
- [ ] Buy groceries <!-- td:a1b2c3d4 -->
```

These are invisible in Obsidian's reading view but allow ThreeDoors to maintain identity across edits.

### Daily Notes

When enabled, ThreeDoors appends new tasks to your daily note under a configurable heading:

```yaml
daily_notes: true
daily_notes_folder: Daily
daily_notes_heading: "## Tasks"
daily_notes_format: "2006-01-02.md"
```

Tasks added via ThreeDoors appear in today's daily note (e.g., `Daily/2026-03-03.md`):

```markdown
## Tasks
- [ ] New task from ThreeDoors <!-- td:uuid -->
```

### Real-Time Sync

ThreeDoors watches your vault for changes using filesystem notifications. When you edit a task in Obsidian, ThreeDoors picks up the change within 100ms. Self-writes are tracked to avoid echo loops (2-second cooldown).

### Path Safety

The Obsidian adapter validates all paths:
- No `..` traversal allowed
- No null bytes
- No absolute paths within vault-relative settings

---

## Apple Notes Integration

### Setup

1. Set your provider to `applenotes` in `config.yaml`:

```yaml
provider: applenotes
note_title: ThreeDoors Tasks
```

2. Create a note in Apple Notes with the exact title you specified.

3. Add tasks using checkbox syntax:

```
- [ ] Task one
- [ ] Task two
- [x] Already done
```

### Requirements

- macOS only (uses AppleScript via `osascript`)
- Full Disk Access may be required in System Settings → Privacy & Security
- Read timeout: 2 seconds, write timeout: 5 seconds

### Bidirectional Sync

Changes flow both ways:
- Tasks added in Apple Notes appear in ThreeDoors
- Tasks completed in ThreeDoors are checked off in Apple Notes
- Task IDs are embedded as HTML comments (invisible in Apple Notes UI)

### Health Check

Run `:health` to verify Apple Notes connectivity. The health checker tests:
- AppleScript execution access
- Note accessibility
- Full Disk Access permissions

---

## Intelligent Features

### Pattern Recognition

After 3+ sessions, ThreeDoors analyzes your behavior:

- **Door position bias** — do you always pick the left door?
- **Task type preferences** — which categories do you complete most?
- **Time-of-day patterns** — when are you most productive?
- **Mood correlations** — which moods lead to more completions?

View these with `:dashboard` or `:insights`.

### Mood-Aware Door Selection

Once you have sufficient session data, logging a mood (`m`) influences which tasks appear. If you historically complete more technical tasks when focused, ThreeDoors will surface technical tasks when you log "Focused."

The algorithm:
1. Looks up your mood in the pattern report
2. Finds your preferred task type and effort level for that mood
3. Scores candidate door sets for diversity + mood alignment
4. Selects the highest-scoring set
5. Enforces a diversity floor (won't show three identical types)

### Avoidance Detection

If a task has been shown 10+ times and you've never selected it, ThreeDoors gently asks what's going on:

- **Reconsider** — set it to in-progress and tackle it now
- **Break down** — open it in detail view to rethink it
- **Defer** — explicitly postpone it
- **Archive** — remove it from the pool entirely

This happens at most once per task per session.

### Insights Dashboard

Access via `:dashboard` or `:insights`. Shows:

- **Completion trends** — last 7 days with sparkline visualization
- **Streaks** — current and longest consecutive completion days
- **Mood and productivity** — average completions per mood state
- **Door position preferences** — left/center/right selection percentages

Requires 3+ sessions for meaningful data.

### Completion Tracking

After completing a task, ThreeDoors shows:
- How many tasks you've completed today
- Comparison with yesterday's count
- Your current streak (consecutive days with at least one completion)

---

## Offline and Sync

### Offline-First Architecture

ThreeDoors uses a Write-Ahead Log (WAL) for crash-safe, offline-first operation. All writes go to a local queue (`sync-queue.jsonl`) before being applied to the provider.

If the provider is unavailable (e.g., Apple Notes is closed, Obsidian vault is on a disconnected drive):
- Operations queue locally
- Automatic replay with exponential backoff when the provider comes back
- Maximum 10,000 queued operations (oldest evicted first)

### Sync Status

The footer shows sync status for each provider:

- **Synced** — up to date
- **Syncing** — sync in progress
- **Pending (N)** — N operations waiting to sync
- **Error** — last sync failed (check `:health` for details)

### Conflict Resolution

For bidirectional providers (Apple Notes, Obsidian), ThreeDoors uses a three-way sync engine. When conflicts occur (same task edited in both places), the most recent change wins. Task IDs embedded as HTML comments ensure identity is preserved.

---

## Session Metrics

### What Is Tracked

Every session records:

| Metric | Description |
|--------|-------------|
| Session ID | Unique identifier per session |
| Start/end time | UTC timestamps |
| Duration | Session length in seconds |
| Tasks completed | Count of completed tasks |
| Doors viewed | How many door sets you saw |
| Refreshes used | Number of times you pressed `s`/`↓` |
| Detail views | How many tasks you opened |
| Door selections | Which position (left/center/right) and which task |
| Task bypasses | Which tasks were shown but not selected |
| Mood entries | Timestamped mood logs |
| Door feedback | Feedback given via `n` key |
| Time to first door | Seconds from launch to first door selection |

### Data Files

- `sessions.jsonl` — one JSON object per line, one line per session
- `patterns.json` — cached analysis report (regenerated when new sessions arrive)
- `completed.txt` — append-only completion log

### Analysis

Run `make analyze` to execute the analysis scripts in `scripts/` against your session data.

View quick stats in-app with `:stats`.

---

## Configuration

### Data Directory

All data lives in `~/.threedoors/`:

```
~/.threedoors/
├── config.yaml         # Provider configuration
├── tasks.yaml          # Active tasks (text file provider)
├── completed.txt       # Completed task log
├── values.yaml         # Your values and goals
├── sessions.jsonl      # Session metrics
├── patterns.json       # Cached pattern analysis
├── enrichment.db       # Cross-reference database
├── improvements.txt    # Exit survey responses
├── sync-queue.jsonl    # WAL pending operations
└── .onboarded          # First-run marker
```

### config.yaml Reference

**Simple mode (single provider):**

```yaml
provider: textfile
note_title: ThreeDoors Tasks
```

**Advanced mode (multiple providers):**

```yaml
providers:
  - name: textfile
    settings:
      task_file: ~/.threedoors/tasks.yaml

  - name: applenotes
    settings:
      note_title: ThreeDoors Tasks

  - name: obsidian
    settings:
      vault_path: /path/to/your/vault
      tasks_folder: tasks
      file_pattern: "*.md"
      daily_notes: true
      daily_notes_folder: Daily
      daily_notes_heading: "## Tasks"
      daily_notes_format: "2006-01-02.md"
```

**Provider options:**

| Provider | Key Settings |
|----------|-------------|
| `textfile` | `task_file` — path to YAML task file |
| `applenotes` | `note_title` — name of the Apple Notes note |
| `obsidian` | `vault_path`, `tasks_folder`, `file_pattern`, `daily_notes`, `daily_notes_folder`, `daily_notes_heading`, `daily_notes_format` |

### Values and Goals

Set up to 5 personal values or goals that display as a reminder while you work. Configure during onboarding or anytime with `:goals edit`.

Values are stored in `~/.threedoors/values.yaml` and displayed in the footer across multiple views.

---

## Troubleshooting

### Health Check

Run `:health` to diagnose issues. The health checker verifies:

1. **Task file** — exists, readable, writable (tests atomic write)
2. **Database** — can load and parse tasks, shows task count
3. **Sync status** — last sync time, warns if stale (>24 hours)
4. **Apple Notes** — AppleScript access, Full Disk Access permissions

Status levels:
- **OK** — everything is working
- **WARN** — something needs attention but isn't broken
- **FAIL** — something is broken and needs fixing

### Common Issues

**"No tasks found"**
- Check that your task file exists: `ls ~/.threedoors/tasks.yaml`
- If using Apple Notes, verify the note title matches your config exactly
- If using Obsidian, verify `vault_path` and `tasks_folder` are correct

**Apple Notes not syncing**
- Grant Full Disk Access: System Settings → Privacy & Security → Full Disk Access → add your terminal app
- Verify the note exists with the exact title from your config
- Run `:health` to see specific error messages

**Obsidian tasks not appearing**
- Confirm `vault_path` points to your vault root (the folder containing `.obsidian/`)
- Check `tasks_folder` is relative to vault root
- Verify your files match `file_pattern` (default: `*.md`)
- Ensure tasks use checkbox syntax: `- [ ] Task text`

**Sync stuck in "Pending"**
- Check `:health` for provider errors
- Verify the provider is accessible (Apple Notes open, vault mounted)
- Pending operations replay automatically with backoff — wait a few seconds
- If stuck, check `sync-queue.jsonl` for error details

**Session data not showing in insights**
- Insights require 3+ completed sessions
- Run `:insights` to see how many more sessions are needed
- Each app launch and quit counts as one session

**Onboarding keeps appearing**
- The `.onboarded` marker file may be missing from `~/.threedoors/`
- Run ThreeDoors and complete the onboarding wizard — it creates the marker automatically

### Resetting Data

To start fresh, remove the data directory:

```bash
rm -rf ~/.threedoors
```

ThreeDoors will recreate it and show onboarding on next launch.

To reset only specific data:
- Delete `sessions.jsonl` to clear session history
- Delete `patterns.json` to force reanalysis
- Delete `values.yaml` to reset goals
- Delete `.onboarded` to re-trigger onboarding

---

*Progress over perfection. Three doors. One choice. Move forward.*
