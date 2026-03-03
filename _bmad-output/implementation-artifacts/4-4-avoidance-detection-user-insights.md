# Story 4.4: Avoidance Detection & User Insights

Status: ready-for-dev

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a user,
I want to be gently informed about my avoidance patterns,
So that I can make conscious decisions about deferred tasks.

## Acceptance Criteria

1. **Given** a task has been shown in doors 5+ times without selection, **When** that task appears in doors again, **Then** a subtle indicator appears (e.g., "You've seen this task 7 times")
2. **And** the system does NOT nag or guilt — framing is informational
3. **And** a `:insights` command shows a summary of patterns ("When stressed, you avoid technical tasks")
4. **And** persistent avoidance (10+ bypasses) triggers a gentle prompt: "This task keeps appearing. Would you like to: [R]econsider, [B]reak down, [D]efer, [A]rchive?"

## Tasks / Subtasks

- [ ] Task 1: Enhanced Avoidance Indicator on Door Cards (AC: #1, #2)
  - [ ] 1.1: Read avoidance data from PatternReport.AvoidanceList in DoorsView
  - [ ] 1.2: Cross-reference current door tasks against avoidance list by task text
  - [ ] 1.3: Render subtle "Seen X times" indicator on door cards for tasks with 5+ bypasses. **IMPORTANT:** The indicator should display `TimesShown` (total appearances) for the "Seen X times" text, NOT `TimesBypassed`. Update `SetAvoidanceData()` to store the full `AvoidanceEntry` or `TimesShown` value. The 5+ threshold check should still use `TimesBypassed` (bypassed 5+ times to show the indicator), but the displayed count should be `TimesShown`.
  - [ ] 1.4: Ensure indicator language is informational, not guilt-inducing (NOT "Avoided")
  - [ ] 1.5: Write tests for avoidance indicator rendering logic

- [ ] Task 2: Persistent Avoidance Prompt (AC: #4)
  - [ ] 2.1: Create AvoidancePromptView in `internal/tui/avoidance_prompt_view.go` for tasks with 10+ bypasses. **Implementation details:**
    - Add `ViewAvoidancePrompt` to the `ViewMode` enum in `main_model.go`
    - Add `ShowAvoidancePromptMsg{Task *tasks.Task}` and `AvoidanceActionMsg{Task *tasks.Task, Action string}` to `messages.go`
    - Add `avoidancePromptView *AvoidancePromptView` field to `MainModel`
    - Add `updateAvoidancePrompt` delegate method
    - Wire in `View()` switch case
    - Follow the exact pattern of `feedbackView`/`ShowFeedbackMsg`/`DoorFeedbackMsg`
    - **Prompt text:** Use this exact text: `"This task has appeared {N} times. What would you like to do?\n\n [R] Reconsider - take it on now\n [B] Break down - look at it closer\n [D] Defer - set aside for later\n [A] Archive - remove from your list\n [Esc] Dismiss"`. Avoid words like 'avoid', 'ignore', 'neglect'.
  - [ ] 2.2: Detect when a 10+ bypassed task is shown in doors and trigger prompt. **Integration point:** After `RefreshDoors()` populates `currentDoors`, iterate over them and check `avoidanceMap`. If any task has `TimesBypassed >= 10`, emit a `ShowAvoidancePromptMsg{Task: task}` that `MainModel` handles by transitioning to `ViewAvoidancePrompt`. Only prompt for the FIRST qualifying task in the set.
  - [ ] 2.3: Implement prompt options: [R]econsider, [B]reak down, [D]efer, [A]rchive
  - [ ] 2.4: Handle each action:
    - Reconsider: select the task (move to in-progress via existing `UpdateStatus`)
    - Break down: navigate to the existing DetailView for the task. The user can then manually add notes or context. There is no automatic decomposition in this story — simply transitioning to the detail view is sufficient.
    - Defer: set status to `deferred`, task excluded from door selection pool
    - Archive: set status to `archived`, set `CompletedAt`, append to `completed.txt` via existing persistence
  - [ ] 2.5: Ensure prompt appears at most once per task per session (prevent nagging). **Mechanism:** Add a `promptedTasks map[string]bool` field to `MainModel` (initialized in `NewMainModel`). Before emitting `ShowAvoidancePromptMsg`, check if `m.promptedTasks[task.Text]` is true. After the prompt is shown (regardless of action), set `m.promptedTasks[task.Text] = true`.
  - [ ] 2.6: Write tests for prompt trigger logic and action handlers. **Required test scenarios:**
    - Task with exactly 9 bypasses does NOT trigger prompt
    - Task with exactly 10 bypasses DOES trigger prompt
    - Prompt shown once for a task, then that task appears again in same session — prompt does NOT appear again
    - User selects [R]econsider — task status changes to `in-progress`
    - User selects [D]efer — task status changes to `deferred`
    - User selects [A]rchive — task is removed from pool
    - User presses Escape — prompt dismissed, no action taken
    - Multiple 10+ tasks in same door set — only prompt for the first one

- [ ] Task 3: Enhanced `:insights` Command (AC: #3)
  - [ ] 3.1: Add a `FormatAvoidanceInsights(report *PatternReport) string` function in `insights_formatter.go`. Wire it via `:insights avoidance` subcommand. **IMPORTANT:** The `:insights mood` subcommand routing was never wired in Story 4.3 — the `case "insights":` handler in `search_view.go` ignores args. You need to add args-based routing (e.g., `args == "mood"` -> `FormatMoodInsights`, `args == "avoidance"` -> `FormatAvoidanceInsights`, no args -> `FormatInsights`) as part of this task.
  - [ ] 3.2: Show avoidance summary: "When [mood], you tend to avoid [type] tasks". **NEW ANALYSIS REQUIRED:** The current `MoodCorrelation` tracks what users SELECT by mood, not what they AVOID. In `pattern_analyzer.go`, extend `analyzeMoodCorrelations` to also track bypassed tasks per mood. For each mood, find the most-bypassed task type (using `taskCategories` map). Store as a new field `AvoidedType string` on `MoodCorrelation`. Note: this changes the `MoodCorrelation` struct which may require updating Story 4.3 tests.
  - [ ] 3.3: Show most avoided task categories (by type and effort)
  - [ ] 3.4: Show avoidance trends (improving/worsening over time). **Data approach:** Calculate trends by comparing avoidance rates in recent sessions (last 5) vs older sessions. Use the `SessionMetrics` timestamps to split sessions into 'recent' and 'historical' buckets. If the bypass rate for a task type is higher in recent sessions, it's 'worsening'; if lower, 'improving'. Display as simple text like 'Technical task avoidance: improving' or 'worsening'.
  - [ ] 3.5: Ensure graceful display when insufficient data (< 5 sessions)
  - [ ] 3.6: Write tests for insights formatting with avoidance data. **NOTE:** Avoidance insights may be multi-line. Verify that `FlashMsg` display works for longer text. If the text is too long for a flash, consider using a dedicated view (like `HealthView`) instead.

- [ ] Task 4: Integration & Edge Cases
  - [ ] 4.1: Ensure avoidance count persists across sessions (stored in patterns.json)
  - [ ] 4.2: Handle edge case: task text changes slightly between sessions. **Decision:** Accept exact-match limitation for this story. Document as a known limitation. Session data (`SessionMetrics.TaskBypasses`) only stores task text, not task ID, so fuzzy matching is out of scope.
  - [ ] 4.3: Handle edge case: task is completed then recreated. **Decision:** A recreated task has a new UUID but potentially the same text. Since avoidance tracking is text-based, the count carries over. Accept this as a known limitation — the count is still technically accurate (the user has seen that text N times).
  - [ ] 4.4: Ensure defer/archive actions properly update task store
  - [ ] 4.5: Verify backward compatibility with existing sessions.jsonl format
  - [ ] 4.6: Run full test suite, ensure no regressions in Stories 4.1-4.3
  - [ ] 4.7: Add `StatusColor` mapping for `deferred` and `archived` statuses in TUI styles. Add `deferred`/`archived` to any status-based filtering in door selection. Update `:help` text if it lists statuses.

## Dev Notes

### Architecture Patterns & Constraints

- **Language:** Go 1.25.4, single binary
- **TUI Framework:** Bubbletea 1.2.4 with Lipgloss 1.0.0 (Elm architecture / MVU pattern)
- **Storage:** YAML tasks file + sessions.jsonl + patterns.json (atomic write pattern)
- **Code Quality:** gofumpt formatting, golangci-lint
- **Error Handling:** Idiomatic Go `(result, error)` tuples, wrap with `%w`
- **Testing:** Table-driven tests, deterministic RNG (seed=42), t.TempDir() for file I/O

### Key Existing Code to Reuse (DO NOT REINVENT)

**Avoidance data already exists from Story 4.2:**
- `internal/tasks/pattern_analyzer.go` — `PatternReport.AvoidanceList` contains tasks with bypass counts
- Avoidance threshold for detection: 3+ times (already in analyzer)
- `AvoidanceEntry` struct has fields: `TaskText string`, `TimesBypassed int`, `TimesShown int`, `NeverSelected bool` — **NOTE: There is NO `LastSeen` or `BypassCount` field. Use the correct field names.**

**Avoidance indicator rendering already partially exists from Story 4.2:**
- `internal/tui/doors_view.go` — Already has avoidance indicator rendering ("Seen X times")
- `DoorsView.SetAvoidanceData()` method exists for passing data
- This story needs to ENHANCE the existing indicator, not recreate it

**Insights formatting already exists from Story 4.2/4.3:**
- `internal/tasks/insights_formatter.go` — `FormatInsights()` and `FormatMoodInsights()`
- Extend with avoidance-specific insights (avoidance by type, mood-avoidance correlation)
- **NOTE:** `:insights mood` subcommand routing was NEVER wired in `search_view.go`. The `parseCommand()` function already splits args, but the `case "insights":` handler ignores args. You must fix this routing as part of Task 3.

**Mood correlation data from Story 4.3:**
- `internal/tasks/mood_selector.go` — MoodAlignmentScore(), SelectDoorsWithMood()
- `internal/tasks/pattern_analyzer.go` — MoodCorrelation with PreferredType/PreferredEffort
- Cross-reference mood data with avoidance patterns for insights

**Session tracking from Stories 1.5/4.2:**
- `internal/tasks/session_tracker.go` — RecordDoorSelection(), RecordRefresh(), LatestMood()
- `TaskBypasses` field tracks which tasks were shown but not selected

### What's NEW in This Story

1. **Persistent Avoidance Prompt** (10+ threshold) — new TUI view/overlay (new file: `avoidance_prompt_view.go`)
2. **Enhanced `:insights` command** — avoidance-specific section with mood-avoidance correlation + args routing fix
3. **Action handlers** for Reconsider/Break Down/Defer/Archive from avoidance prompt
4. **Defer/Archive actions** — new task status transitions (ripples through multiple files, see below)
5. **Mood-avoidance cross-reference** — new analysis logic in PatternAnalyzer

### Critical: What NOT to Do

- Do NOT recreate avoidance detection logic — it's in PatternAnalyzer already
- Do NOT change the avoidance threshold (3+ for detection, 5+ for indicator, 10+ for prompt)
- Do NOT use nagging/guilt language — informational only ("You've seen this" not "You're avoiding this")
- Do NOT make the prompt modal/blocking — it should be dismissible with single keypress
- Do NOT add external dependencies — use stdlib only
- Do NOT modify sessions.jsonl format — backward compatibility required

### Project Structure Notes

```
internal/
  tasks/
    task_status.go           # ADD StatusDeferred, StatusArchived constants + validation + transitions
    task.go                  # Update Validate() for new statuses, archived sets CompletedAt
    pattern_analyzer.go      # Extend analyzeMoodCorrelations with AvoidedType field
    insights_formatter.go    # ADD FormatAvoidanceInsights() function
    session_tracker.go       # Read-only for this story
    mood_selector.go         # Read-only for this story
    door_selector.go         # Verify deferred tasks are excluded from selection
  tui/
    messages.go              # ADD ShowAvoidancePromptMsg, AvoidanceActionMsg
    main_model.go            # ADD ViewAvoidancePrompt to ViewMode, wire view, add promptedTasks map
    doors_view.go            # Enhance avoidance indicator, trigger prompt for 10+ tasks
    avoidance_prompt_view.go # NEW: Avoidance prompt overlay view
    search_view.go           # Fix :insights args routing, add :insights avoidance subcommand
    styles.go                # ADD StatusColor mapping for deferred/archived
```

### Task Status Transitions for New Actions

Current valid statuses in `task_status.go`: `todo`, `blocked`, `in-progress`, `in-review`, `complete`

**New statuses to add:**
- `deferred` — temporarily removed from active pool, can be un-deferred later. Does NOT set `CompletedAt`.
- `archived` — permanently removed from active pool. Sets `CompletedAt`. Appended to `completed.txt`.

**Implementation in `task_status.go`:**
1. Add `StatusDeferred TaskStatus = "deferred"` and `StatusArchived TaskStatus = "archived"` constants
2. Add them to `ValidateStatus()` switch cases
3. Add transition entries to `validTransitions` map:
   - `StatusTodo: {...existing..., StatusDeferred, StatusArchived}`
   - `StatusDeferred: {StatusTodo}` (un-defer returns to todo)
   - `StatusArchived: {}` (terminal state, no transitions out)
4. Verify `SelectDoors()` in `door_selector.go` filters by status — if it only picks `todo` tasks, `deferred` status is sufficient for exclusion. If it picks all tasks, add explicit filtering.

### Avoidance Thresholds (from PRD + Stories 4.2)

| Threshold | Behavior | Source |
|-----------|----------|--------|
| 3+ bypasses | Detected in PatternAnalyzer.AvoidanceList | Story 4.2 |
| 5+ bypasses | Subtle "Seen X times" indicator on door card (display TimesShown, threshold on TimesBypassed) | Story 4.4 AC#1 |
| 10+ bypasses | Gentle prompt with action options (threshold on TimesBypassed) | Story 4.4 AC#4 |

### Testing Standards

- **Table-driven tests** (Go convention)
- **Deterministic RNG** with fixed seed (seed=42) for reproducible assertions
- **t.TempDir()** for file I/O tests
- **No mocking frameworks** — use interfaces and simple stubs
- **Test helpers**: Use existing `makeTestSession`, `newCategorizedTestTask` helpers
- **Coverage target**: 70%+ for internal/tasks/, 20%+ for internal/tui/

### Known Limitations (accepted for this story)

1. **Exact text matching for avoidance:** Avoidance tracking uses exact task text. If a user edits task text, bypass count resets. Session data only stores text, not task ID.
2. **Completed-then-recreated tasks:** A recreated task with the same text will inherit the old bypass count. Acceptable since the user has genuinely seen that text N times.
3. **`:insights mood` routing prerequisite:** Must fix the args routing in search_view.go that was missed in Story 4.3.

### References

- [Source: docs/prd/epics-and-stories.md#Story 4.4]
- [Source: _bmad-output/planning-artifacts/epics.md#Story 4.4: Avoidance Detection & User Insights]
- [Source: docs/architecture/index.md - Five-Layer Architecture]
- [Source: internal/tasks/pattern_analyzer.go - AvoidanceList, AvoidanceEntry struct]
- [Source: internal/tasks/task_status.go - ValidateStatus, validTransitions, IsValidTransition]
- [Source: internal/tui/doors_view.go - existing avoidance indicator, SetAvoidanceData]
- [Source: internal/tui/main_model.go - ViewMode enum, feedbackView pattern]
- [Source: internal/tui/messages.go - existing message types pattern]
- [Source: internal/tasks/insights_formatter.go - FormatInsights(), FormatMoodInsights()]
- [Source: internal/tasks/session_tracker.go - TaskBypasses, RecordRefresh()]

### Previous Story Intelligence

**From Story 4.1:**
- Diversity-preferring door selection with 10-candidate scoring
- Category badges rendered on door cards (type icon + effort label)
- Tag parser for inline categories (#type @effort +location)
- State machine pattern for tag editing view

**From Story 4.2:**
- Pattern analyzer with 6 dimensions including avoidance detection
- AvoidanceList populated with tasks shown 3+ times without selection
- Avoidance indicator ("Seen X times") already rendered on door cards
- `:insights` command shows pattern summary
- Background goroutine analysis on startup, atomic.Pointer sharing

**From Story 4.3:**
- Mood-aware door selection with combined diversity+mood scoring
- MoodCorrelation with PreferredType/PreferredEffort
- `:insights mood` subcommand for mood-focused analysis (NOTE: routing not wired)
- LatestMood() accessor on SessionTracker
- Diversity floor: at least 1 non-matching door

### Git Intelligence

Recent commits show pattern:
- Stories 4.1-4.3 follow consistent patterns: core logic in `internal/tasks/`, TUI in `internal/tui/`
- Test files mirror source files with `_test.go` suffix
- Implementation artifacts saved to `_bmad-output/implementation-artifacts/`
- gofumpt formatting required before commit

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

### File List
