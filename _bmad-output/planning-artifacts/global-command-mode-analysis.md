# Global Command Mode (`:`) Accessibility & Autocomplete Analysis

**Date:** 2026-03-09
**Trigger:** Course correction — users must navigate back to doors view to type `:` commands
**Participants:** UX Designer (analysis), Dev (architecture), PM (story creation)

---

## Problem Statement

The `:` command entry (for `:help`, `:dashboard`, `:stats`, `:mood`, etc.) currently only works from the doors view. Users on the help screen, detail view, insights view, or any other screen must first navigate back to doors view before typing a command. This creates unnecessary friction — especially since `q` (quit) and `?` (help) are already globally accessible from any view.

Additionally, the command palette lacks discoverability — users must memorize all 16+ available commands. A completion/autocomplete system would reduce friction and aid discovery.

## Research Findings

### Current State (main_model.go)

| Key | Scope | Guard | Location |
|-----|-------|-------|----------|
| `q` | **Global** (any view) | `isTextInputActive()` | Line 911 |
| `?` | **Global** (any view) | `isTextInputActive()` | Line 916 |
| `:` | **Doors only** | None (inside `updateDoors()`) | Line 1033-1040 |
| `/` | **Doors only** | None (inside `updateDoors()`) | Line 1027-1031 |

### Views Where `:` Does NOT Work (17 of 18 non-doors views)

1. ViewDetail — no conflict (`:` not used)
2. ViewMood — potential conflict when `isCustom` (text input active)
3. ViewSearch — MUST NOT intercept (text input for search/commands)
4. ViewHealth — no conflict
5. ViewAddTask — MUST NOT intercept (text input)
6. ViewValuesGoals — potential conflict when text input focused
7. ViewFeedback — potential conflict when `isCustom`
8. ViewImprovement — MUST NOT intercept (text input)
9. ViewNextSteps — no conflict
10. ViewAvoidancePrompt — no conflict
11. ViewInsights — no conflict
12. ViewOnboarding — potential conflict (values step has text input)
13. ViewConflict — no conflict
14. ViewSyncLog — no conflict
15. ViewThemePicker — no conflict
16. ViewDevQueue — no conflict
17. ViewProposals — no conflict
18. ViewHelp — no conflict

### Guard Analysis

The existing `isTextInputActive()` function (line 1184-1210) already correctly identifies all views where text input is active. This is the same guard used for `q` and `?`. Using it for `:` prevents the exact same set of conflicts.

**Key insight:** If `q` is safe to intercept globally with this guard, then `:` is equally safe. The guard was purpose-built for this pattern.

### Command Inventory (16 commands)

| Command | Arguments | Description |
|---------|-----------|-------------|
| `:add` | `<text>` or `--why` | Create new task |
| `:add-ctx` | `<text>` | Create task with context prompt |
| `:mood` | `[mood_string]` | Record mood or open dialog |
| `:stats` | — | Show session statistics |
| `:health` | — | Run health check |
| `:dashboard` | — | Open insights dashboard |
| `:insights` | `[mood\|avoidance]` | Show pattern insights |
| `:goals` | `[edit]` | View or edit values/goals |
| `:synclog` | — | View sync operation log |
| `:tag` | — | Edit task categories/tags |
| `:theme` | — | Open theme picker |
| `:devqueue` | — | View dev dispatch queue |
| `:suggestions` | — | View AI task suggestions |
| `:dispatch` | — | Dev dispatch info |
| `:help` | — | Show help screen |
| `:quit` / `:exit` | — | Exit application |

---

## Part 1: Global `:` Command Mode

### Adopted: Global `:` via MainModel-level interception (same as D-059 pattern)

**Rationale:**
- Exact precedent exists: `q` (D-059) and `?` (D-087) use this pattern
- `isTextInputActive()` guard prevents all text input conflicts
- One code change location (move from `updateDoors()` to global section)
- Zero per-view changes needed
- Consistent user mental model: "global shortcuts always work unless I'm typing"

**Implementation:**
1. Add `:` interception between the `?` handler and the view delegation switch (lines 917-920)
2. Guard with `!m.isTextInputActive()`
3. Remove `:` case from `updateDoors()` (lines 1033-1040)
4. The `previousView` field should capture `m.viewMode` (not hardcode `ViewDoors`)

### Rejected: Per-view `:` handler

**Why rejected:** Would require modifying 17 view Update() methods or 17 `updateXXX()` methods in main_model.go. Higher maintenance burden. Same approach was rejected for `q` in X-024.

### Rejected: Keep `:` doors-only, add `:` to "a few more" views

**Why rejected:** Half-measures create inconsistent behavior. Users would need to memorize which views support `:` and which don't. The global approach is simpler and more predictable.

### Out of Scope: Making `/` (search) global

While similar, `/` opens a search view that's conceptually different from command entry. Search makes most sense from the doors view where you're browsing tasks. Command mode (`:`) is a meta-action that transcends any particular view. `/` globalization could be a separate story if demand emerges.

### Impact on Existing Behavior

- **Doors view:** Identical behavior (`:` still works, just handled at a different code level)
- **Text input views:** No change (guard prevents interception)
- **All other views:** `:` now opens command palette — new capability, no conflicts
- **previousView tracking:** Commands entered from non-doors views will return to the originating view on cancel/completion (currently hardcoded to `ViewDoors`)

---

## Part 2: Command Autocomplete/Completion

### Charm Ecosystem Research

**Available components in `bubbles v1.0.0` (already imported):**
- `bubbles/textinput` — currently used for search/command input
- `bubbles/list` — has built-in fuzzy filtering via `sahilm/fuzzy`, ranked results, pagination

**No dedicated autocomplete component exists in the Charm ecosystem.** The `textinput` component has no built-in suggestion/completion API. The `list` component has filtering but is a full-featured list widget, not a lightweight suggestion dropdown.

**How other Bubbletea apps handle command completion:**
- Most use custom implementations (simple prefix matching on a small command list)
- Some pair `textinput` with a custom filtered list rendered below the input
- The `bubbles/list` component is overkill for 16 items — it adds fuzzy matching, ranking, pagination, and delegate rendering infrastructure meant for hundreds of items

### Adopted: Custom lightweight completion in SearchView

**Rationale:**
- Only 16 commands — instant prefix filtering with a simple `strings.HasPrefix` loop
- No external dependency or complex component needed
- Consistent with existing `filterTasks()` pattern (case-insensitive substring match)
- Full control over rendering (can show descriptions alongside command names)

**UX Design:**

1. **Trigger:** When `isCommandMode` is true and user has typed at least `:` + one character
2. **Filtering:** Case-insensitive prefix match on command name (`:a` → `:add`, `:add-ctx`)
3. **Display:** Suggestion list renders between the header and the text input, showing:
   - Command name (bold/highlighted)
   - Brief description (dim, right-aligned or after dash separator)
   - Example: `  :add       — Create a new task`
4. **Navigation:** Arrow keys (↑/↓) or Tab to cycle through suggestions
5. **Selection:** Enter on a highlighted suggestion fills the command, Tab completes it
6. **Backspace:** Dynamically updates the suggestion list
7. **Layout:** Suggestions rendered inline (push content down), not overlay — consistent with how search results already render in SearchView
8. **Empty state:** When `:` alone is typed with no further characters, show all commands (serves as a command reference)

### Rejected: `bubbles/list` component for command suggestions

**Why rejected:** Heavyweight for 16 items. Adds fuzzy matching complexity (user types `:dsh` expecting `:dashboard` — fuzzy matching would show it, but users expect prefix matching from vim/command-line muscle memory). The fuzzy matching library adds a transitive dependency for minimal value.

### Rejected: Overlay-style dropdown

**Why rejected:** No overlay rendering system exists in the TUI. All existing patterns (search results, help content, proposals) use inline rendering that pushes content down. An overlay would require Z-ordering infrastructure that doesn't exist and would be the only instance of it.

### Rejected: Argument-level completion (e.g., `:insights m` → `:insights mood`)

**Why rejected for MVP:** Only 2 commands have arguments that could be completed (`:insights mood|avoidance`, `:goals edit`). The complexity of context-aware argument completion is disproportionate to the value for 2 commands. Can be added later if demand emerges.

### Story Split Decision: Two stories, not one

**Rationale:** Global `:` and autocomplete are independent improvements:
- Global `:` is a minimal, high-value change (move 8 lines of code, fix `previousView`)
- Autocomplete is a larger UX feature with its own ACs, rendering, and test surface
- A user benefits from global `:` even without autocomplete
- Autocomplete benefits from global `:` but doesn't depend on it architecturally

**Stories:**
- **39.7:** Global `:` Command Mode (S estimate) — move `:` to MainModel-level with `isTextInputActive()` guard
- **39.8:** Command Autocomplete/Completion (M estimate) — dynamic command suggestions in SearchView

---

## Decisions Summary

| ID | Decision | Rationale |
|----|----------|-----------|
| Adopted | Global `:` via MainModel-level `isTextInputActive()` guard | Follows D-059/D-087 pattern; one change location; zero per-view modifications |
| Adopted | Custom lightweight completion (not bubbles/list) | 16 commands; prefix match sufficient; no dependency needed |
| Adopted | Inline suggestion rendering (push content down) | Consistent with SearchView pattern; no overlay infrastructure needed |
| Adopted | Two separate stories (39.7 + 39.8) | Independent improvements; different scope; user benefits from each alone |
| Rejected | Per-view `:` handlers | 17 views to modify; maintenance burden; rejected for `q` in X-024 |
| Rejected | Selective view additions | Inconsistent UX; users must memorize which views support `:` |
| Rejected | `bubbles/list` for command suggestions | Heavyweight; fuzzy matching wrong for command completion |
| Rejected | Overlay-style dropdown | No overlay infrastructure; only inline patterns exist |
| Rejected | Argument-level completion (MVP) | Only 2 commands benefit; disproportionate complexity |
| Out of scope | Global `/` (search) | Different semantics; search is task-browsing, commands are meta-actions |
