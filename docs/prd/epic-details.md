# Epic Details

## Epic 1: Three Doors Technical Demo

**Epic Goal:** Build and validate the Three Doors interface with minimal viable functionality to prove the UX concept reduces friction compared to traditional task lists.

**Scope:** CLI/TUI application that reads tasks from a text file, presents three random tasks as "doors," allows refresh and selection, marks tasks complete, and tracks progress.

**Story Sequence Optimization:** Stories reordered to validate refresh UX before completion (moved 1.6→1.4). Non-essential features simplified or made optional to focus on core validation.

---

### Story 1.1: Project Setup & Basic Bubbletea App

**As a** developer,
**I want** a working Go project with Bubbletea framework,
**so that** I have a foundation for building the Three Doors TUI.

**Acceptance Criteria:**
1. Go module initialized with `go mod init github.com/arcaven/ThreeDoors`
2. Bubbletea and Lipgloss dependencies added
3. Basic TUI application renders "ThreeDoors - Technical Demo" header
4. Application responds to 'q' keypress to quit
5. `Makefile` with `build`, `run`, and `clean` targets works
6. Application compiles and runs without errors

**Estimated Time:** 30-45 minutes

---

### Story 1.2: File I/O for Tasks (Simplified)

**As a** user,
**I want** the app to read my tasks from a simple text file,
**so that** I can populate tasks easily outside the app.

**Acceptance Criteria:**
1. Application creates `~/.threedoors/` directory on first run if it doesn't exist
2. Application reads tasks from `~/.threedoors/tasks.txt` (one task per line, simple parsing)
3. Application displays count of loaded tasks (e.g., "Loaded 12 tasks")
4. If `tasks.txt` missing or empty, show message: "Create ~/.threedoors/tasks.txt and add tasks (one per line) to get started"
5. Gracefully handles file read errors with helpful error message

**Deferred to MVP:**
- Auto-creation of sample tasks (you'll populate real tasks manually)
- Comment support (`#` lines) - simpler parsing for Tech Demo
- Empty line handling - just skip them or treat as blank tasks

**Estimated Time:** 20-30 minutes

---

### Story 1.3: Three Doors Display

**As a** user,
**I want** to see three tasks displayed as "doors,"
**so that** I can quickly choose what to work on without scanning a long list.

**Acceptance Criteria:**
1. Three tasks are randomly selected from loaded tasks
2. Tasks are displayed in three visual "boxes" (ASCII art or Lipgloss styled borders)
3. Each door is labeled: "Door 1", "Door 2", "Door 3"
4. Task text is displayed inside each door (truncated if too long, max ~40 chars)
5. Instructions displayed at bottom: "Press 1, 2, or 3 to select | R to refresh | Q to quit"
6. No duplicate tasks appear in the three doors simultaneously
7. If fewer than 3 tasks available, show what's available (handle edge case gracefully)

**Estimated Time:** 45-60 minutes

---

### Story 1.4: Door Refresh Mechanism (MOVED UP)

**As a** user,
**I want** to refresh the three doors if none appeal to me,
**so that** I have control over my options without feeling trapped.

**Acceptance Criteria:**
1. Pressing R generates a new set of three doors
2. New selection is different from current selection (no duplicates of currently shown tasks)
3. Random selection ensures variety over multiple refreshes
4. Edge case: If 3 or fewer tasks remain total, show message "All available tasks are already showing"

**Deferred to MVP:**
- Refresh count tracking/display (not essential for validation)

**Rationale for Moving Up:** Validates the refresh UX flow (display → refresh → refresh → select) before implementing completion. User flow is: see doors, refresh until one appeals, then select.

**Estimated Time:** 15-20 minutes

---

### Story 1.5: Door Selection & Task Completion (MERGED)

**As a** user,
**I want** to select a door and mark the task as complete,
**so that** I can make progress on my tasks.

**Acceptance Criteria:**
1. Pressing 1, 2, or 3 selects the corresponding door
2. Selected task is highlighted/indicated visually
3. Prompt appears: "Working on: [task text] - Press C to complete, B to go back"
4. Pressing C marks task as complete
5. Completed task is removed from available task pool (in-memory)
6. New set of three doors is displayed automatically after completion
7. Session completion count increments and displays (e.g., "✓ Completed this session: 3")
8. "Progress over perfection" message shown after completing a task (e.g., "Nice! Any progress is good progress.")
9. **OPTIONAL (implement if time allows):** Completed tasks appended to `~/.threedoors/completed.txt` with timestamp format `[YYYY-MM-DD HH:MM:SS] task description`

**Rationale for Merge:** Persistent storage (completed.txt) is nice-to-have but not essential for validating Three Doors UX. Session count in-memory is sufficient. If you have extra time, add file persistence.

**Estimated Time:** 45-60 minutes (without persistence), 60-75 minutes (with persistence)

---

### Story 1.6: Essential Polish (SIMPLIFIED)

**As a** user,
**I want** the app to feel polished enough to use daily,
**so that** I enjoy the validation experience.

**Acceptance Criteria:**
1. Lipgloss styling applied: distinct colors for doors, success messages (green), prompts (yellow/blue)
2. "Progress over perfection" message embedded in interface (startup greeting or post-completion)
3. Application feels responsive and smooth (no noticeable lag)

**Deferred to MVP:**
- README.md (you're the only user for validation)
- Extensive edge case handling (all tasks completed celebration, 1-2 tasks remaining display logic)
- Advanced error messaging

**Rationale for Simplification:** Focus on making the core experience pleasant. Edge cases are unlikely to be hit during 1-week validation. README isn't needed when you built it.

**Estimated Time:** 20-30 minutes

---

## Epic 2-5: Post-Validation Epics (Placeholder)

**Note:** These epics are placeholders for post-validation planning. Detailed stories will be created only if Epic 1 successfully validates the Three Doors concept.

**Epic 2: Foundation & Apple Notes Integration**
*Stories to be defined after Epic 1 validation and Apple Notes integration spike*

**Epic 3: Enhanced Interaction & Task Context**
*Stories to be defined based on learnings from Epic 2 usage patterns*

**Epic 4: Learning & Intelligent Door Selection**
*Stories to be defined based on sufficient usage data from Epic 3*

**Epic 5: Data Layer & Enrichment (Optional)**
*Stories to be defined only if clear need emerges from Epic 4*

---
