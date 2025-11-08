# ThreeDoors Product Requirements Document (PRD)

**Document Version:** 1.1 (Technical Demo & Validation Phase)
**Last Updated:** 2025-11-07
**Project Repository:** github.com/arcaven/ThreeDoors.git

---

## Goals and Background Context

### Goals

**Technical Demo & Validation Phase (Pre-MVP):**
- Validate the Three Doors UX concept in 1 week (4-8 hours of development)
- Prove the core hypothesis: "Presenting three diverse tasks is better than presenting a list"
- Build working TUI with Bubbletea to demonstrate feasibility
- Use simple local text file for rapid task population and testing
- Gather real usage feedback before investing in complex integrations

**Full MVP Goals (Post-Validation):**
- Master BMAD methodology through authentic, real-world application
- Create a todo app that reduces friction and actually helps with organization
- Build a personal achievement partner that works with human psychology, not against it
- Enable seamless cross-context navigation across multiple platforms and tools
- Capture the full story (what AND why) to improve stakeholder communication
- Achieve measurably better personal organization than current scattered approach
- Demonstrate progress-over-perfection philosophy in both product design and development process

### Background Context

Traditional todo apps work well for already-organized people, but they're fundamentally rudimentary tools that haven't evolved alongside modern technology capabilities. While they help those who are naturally organized stay organized, they offer little support for adapting to the dynamic reality of modern life—where the same person occupies multiple roles (employee, parent, partner, learner), experiences varying moods and energy states, and faces constantly shifting priorities.

ThreeDoors recognizes that as technology has advanced, we can offer substantially more support. We can organize our organization tools themselves, bringing together tasks scattered across multiple systems. More importantly, we can adapt technology support dynamically: responding to the user's current context, role, mood, and circumstances, re-routing based on changing conditions and priorities. This PRD defines the MVP: a CLI/TUI application with Apple Notes integration that begins this journey, embodying "progress over perfection" philosophy while serving as a practical demonstration of the BMAD methodology.

### Change Log

| Date | Version | Description | Author |
|------|---------|-------------|--------|
| 2025-11-07 | 1.0 | Initial PRD creation from project brief | John (PM Agent) |
| 2025-11-07 | 1.1 | Pivoted to Technical Demo & Validation approach (Option C): Simplified to text file storage, 1-week validation timeline, deferred Apple Notes and learning features to post-validation phases | John (PM Agent) |

---

## Requirements

### Technical Demo & Validation Phase Requirements

**Core Requirements (Week 1 - Build & Validate):**

**TD1:** The system shall provide a CLI/TUI interface optimized for terminal emulators (iTerm2 and similar)

**TD2:** The system shall read tasks from a simple local text file (e.g., `~/.threedoors/tasks.txt`)

**TD3:** The system shall display the Three Doors interface showing three tasks selected from the text file

**TD4:** The system shall allow users to select a door (press 1, 2, or 3) to start working on that task

**TD5:** The system shall allow users to mark the selected task as complete

**TD6:** The system shall track and display task completion count for the current session

**TD7:** The system shall provide a refresh mechanism (press R) to generate a new set of three doors

**TD8:** The system shall embed "progress over perfection" messaging in the interface

**TD9:** The system shall write completed tasks to a separate file (e.g., `~/.threedoors/completed.txt`) with timestamp

**Success Criteria for Phase:**
- Built and running within 4-8 hours of development time
- Developer uses it daily for 1 week to validate UX concept
- Three Doors selection feels meaningfully different from a simple list
- Decisions made on whether to proceed to Full MVP based on real usage

---

### Full MVP Requirements (Post-Validation - Deferred)

**Phase 2 - Apple Notes Integration:**

**FR2:** The system shall integrate with Apple Notes as the primary task storage backend, enabling bidirectional sync

**FR4:** The system shall retrieve and display tasks from Apple Notes within the application interface

**FR5:** The system shall allow users to mark tasks as complete, updating both the application state and Apple Notes

**FR12:** The system shall allow updates to tasks from either the application or directly in Apple Notes on iPhone, with changes reflected bidirectionally

**FR15:** The system shall provide a health check command to verify Apple Notes connectivity and database integrity

**Phase 3 - Enhanced Interaction & Learning:**

**FR3:** The system shall allow users to capture new tasks with optional context (what and why) through the CLI/TUI

**FR6:** The system shall display user-defined values and goals persistently throughout task work sessions

**FR7:** The system shall provide a "choose-your-own-adventure" interactive navigation flow that presents options rather than demands

**FR8:** The system shall track daily task completion count and display comparison to previous day's count

**FR9:** The system shall prompt the user once per session with: "What's one thing you could improve about this list/task/goal right now?"

**FR10:** The system shall embed "progress over perfection" messaging throughout interaction patterns and interface copy (enhanced beyond Tech Demo)

**FR16:** The system shall support a "quick add" mode for capturing tasks with minimal interaction

**FR18:** The system shall allow users to provide feedback on why a specific door isn't suitable with options: Blocked, Not now, Needs breakdown, or Other comment

**FR19:** The system shall capture and store blocker information when a task is marked as blocked

**FR20:** The system shall use door selection and feedback patterns to inform future door selection (learning which task types suit which contexts)

**FR21:** The system shall categorize tasks by type, effort level, and context to enable diverse door selection

**Phase 4 - Data Layer & Enrichment:**

**FR11:** The system shall maintain a local enrichment layer (SQLite and/or vector database) for metadata, cross-references, and relationships that cannot be stored in source systems

### Non-Functional Requirements

**Technical Demo Phase:**

**TD-NFR1:** The system shall be built in Go 1.25.4+ using idiomatic patterns and gofumpt formatting standards

**TD-NFR2:** The system shall use the Bubbletea/Charm Bracelet ecosystem for TUI implementation

**TD-NFR3:** The system shall operate on macOS as the primary target platform

**TD-NFR4:** The system shall store all data in local text files (`~/.threedoors/` directory) with no external services or telemetry

**TD-NFR5:** The system shall respond to user interactions within the CLI/TUI with minimal latency (target: <100ms for typical operations given simple file I/O)

**TD-NFR6:** The system shall use Make as the build system with simple targets: `build`, `run`, `clean`

**TD-NFR7:** The system shall gracefully handle missing or corrupted task files by creating defaults

---

**Full MVP Phase (Post-Validation - Deferred):**

**NFR1:** The system shall maintain idiomatic Go patterns and gofumpt formatting standards

**NFR2:** The system shall continue using Bubbletea/Charm Bracelet ecosystem

**NFR3:** The system shall operate on macOS as primary platform

**NFR4:** The system shall store all user data locally or in user's iCloud (via Apple Notes), with no external telemetry or tracking

**NFR5:** The system shall store application state and enrichment data locally (cross-computer sync deferred to post-MVP)

**NFR6:** The system shall respond to user interactions within the CLI/TUI with minimal latency (target: <500ms for typical operations)

**NFR7:** The system shall provide graceful degradation when Apple Notes integration is unavailable, maintaining core functionality

**NFR8:** The system shall implement secure credential storage using OS keychain for any API keys or authentication tokens

**NFR9:** The system shall never log sensitive user data or credentials

**NFR10:** The system shall use Make as the build system

**NFR11:** The system shall maintain clear architectural separation between core engine, TUI layer, integration adapters, and enrichment storage

**NFR12:** The system shall maintain data integrity even when Apple Notes is modified externally while app is running

---

## User Interface Design Goals

### Overall UX Vision

ThreeDoors presents as a conversational partner rather than a demanding taskmaster. The central interface metaphor is literal: **three doors, three tasks, three different on-ramps to action**. At each session start, the user is presented with three carefully selected tasks that are very different from each other—different types of activities, different effort levels, different contexts—but all represent good starting points based on priorities. This design serves dual purposes: it gets the user in the habit of doing *something* (reducing inertia), and it teaches the tool about the user's current state by observing which types of tasks they gravitate toward or avoid.

The interface should feel like opening a dialogue, not confronting a backlog. Users are greeted with options that respect their current capacity—whether focused, overwhelmed, or stuck—and celebrate any choice as progress.

### Key Interaction Paradigms

**The Three Doors (Primary Interaction):**
The main interface presents three tasks simultaneously as entry points. These tasks should be:
- **Intentionally diverse** - Different types of activities (e.g., creative vs. administrative vs. physical, or high-focus vs. low-friction vs. context-switching)
- **Small at first** - Especially in early usage, doors should present approachable tasks to build momentum
- **All viable next steps** - Each represents a legitimate priority, not filler options
- **Learning opportunities** - User's choice (or avoidance) informs the system about current mood, energy, and capacity

Over time, the system learns: "On Tuesday mornings, user picks Door 1 (focused work). On Friday afternoons, user picks Door 3 (quick wins). User never picks administrative tasks before 10am."

**Door Refresh & Feedback (MVP Core):**
- **Refresh/New Doors** - Simple keystroke (e.g., 'R' or 'N') to generate a new set of three doors if current options don't appeal. No judgment, no friction—just new options.
- **Door Feedback** - Option to indicate why a door isn't suitable (basic MVP options):
  - "Blocked" - Task cannot proceed (captures blocker)
  - "Not now" - Task is valid but doesn't fit current mood/context (teaches system about state)
  - "Needs breakdown" - Task is too big/unclear (MVP: flag for later attention; Post-MVP: may trigger breakdown assistance)
  - "Other comment" - Freeform note about the task (refactoring, context, etc.)

These interactions serve dual purposes: give users control (preventing feeling trapped) and provide rich learning signal to the system about task suitability, blockers, and user state.

**Choose-Your-Own-Adventure Navigation:**
Beyond the three doors, other decision points present 3-5 contextual options rather than requiring command memorization. Options adapt based on state and history.

**Progressive Disclosure:**
Start simple, reveal complexity only when needed. Quick add mode for speed, expanded capture for context when desired. Don't force decisions upfront.

**Persistent Context:**
Values/goals remain visible (but unobtrusive) throughout the session—likely as a subtle header or footer—reminding users of the "why" while working on the "what."

**Encouraging Tone:**
All messaging embodies "progress over perfection." Copy celebrates any action ("You picked a door and started. That's what matters.").

### Core Screens and Views

From a product perspective, these are the critical views necessary to deliver MVP value:

1. **Three Doors Dashboard (Primary Interface)** - Session entry point presenting three diverse tasks as "doors" to choose, with minimal surrounding context. Core question: "Which door feels right today?" Includes refresh option and per-door feedback mechanism.

2. **Task List View** - Full task display when user wants to see beyond the three doors, with filtering and status

3. **Quick Add Flow** - Minimal-friction task capture (possibly single input field)

4. **Extended Capture Flow** - Optional deeper capture including "why" context and task metadata (effort, type, context)

5. **Values/Goals Setup** - Initial and ongoing management of user-defined values that guide prioritization

6. **Progress View** - Visualization showing "better than yesterday" metrics and door choice patterns over time (e.g., "You've opened 5 doors this week, up from 3 last week" and "You tend to pick Door 1 in mornings, Door 3 in afternoons")

7. **Health Check View** - Diagnostic display showing Apple Notes connectivity and sync status

8. **Improvement Prompt** - End-of-session single question asking for one improvement

### Accessibility

**None** - MVP focuses on terminal interface for single user (developer). Accessibility requirements deferred to future phases when/if user base expands beyond CLI-comfortable users.

### Branding

**Terminal Aesthetic with Warmth:**
Leverage Charm Bracelet/Bubbletea's capabilities for styled terminal UI—think clean, readable typography with subtle use of color for status indication (green for progress, yellow for prompts, red sparingly for errors).

**Three Doors Visual Metaphor:**
The main interface could literally render three visual "doors" in ASCII art or styled terminal boxes:
```
┌─────────────┐  ┌─────────────┐  ┌─────────────┐
│   DOOR 1    │  │   DOOR 2    │  │   DOOR 3    │
│             │  │             │  │             │
│  [Task A]   │  │  [Task B]   │  │  [Task C]   │
│  Quick win  │  │  Deep work  │  │  Creative   │
│  ~5min      │  │  ~30min     │  │  ~15min     │
└─────────────┘  └─────────────┘  └─────────────┘

Press 1, 2, or 3 to enter  |  R to refresh  |  B to mark blocked
```

**"Progress Over Perfection" Visual Language:**
Use asymmetry, incomplete progress bars, and "good enough" indicators. The three doors might be slightly different sizes or styles, reinforcing that perfection isn't required—just pick one and start.

### Target Device and Platforms

**Primary: macOS Terminal Emulators (iTerm2, Terminal.app, Alacritty)**
- CLI/TUI optimized for 80x24 minimum, responsive to larger terminal sizes
- Assumes modern terminal with 256-color support minimum
- Keyboard-driven navigation (arrow keys, vim-style hjkl, number keys 1-3 for door selection)

**Secondary: Remote Terminal Access**
- Should function over SSH connections (for future Geodesic/remote environment access)
- ASCII fallback for constrained environments

**Mobile Access (Indirect):**
- No dedicated mobile UI in MVP
- Mobile interaction happens through Apple Notes app directly (view/edit tasks on iPhone)
- Sync bidirectionally when user returns to terminal interface

---

## Technical Assumptions

### Technical Demo Phase Architecture

**Decision:** Minimal monolithic application with simple text file I/O

**Rationale:**
- **Speed to validation**: Build and test in 4-8 hours
- **Simple is fast**: No database, no complex integrations, no abstractions until needed
- **Easy external task population**: Text files can be edited with any editor, populated from scripts, etc.
- **Prove the concept first**: Validate Three Doors UX before investing in infrastructure
- **Low risk**: Can throw away and rebuild if concept fails validation

**Tech Demo Structure:**
```
ThreeDoors/
├── cmd/
│   └── threedoors/        # Main application (single file initially)
├── internal/
│   ├── tui/              # Bubbletea Three Doors interface
│   └── tasks/            # Simple file I/O (read tasks.txt, write completed.txt)
├── docs/                  # Documentation (including this PRD)
├── .bmad-core/           # BMAD methodology artifacts
├── Makefile              # Simple build: build, run, clean
└── README.md             # Quick start guide
```

**Data Files (created at runtime in `~/.threedoors/`):**
```
~/.threedoors/
├── tasks.txt             # One task per line (user can edit directly)
├── completed.txt         # Completed tasks with timestamps
└── config.txt            # Optional: Simple key=value config (if needed)
```

---

### Full MVP Architecture (Post-Validation - Deferred)

**Structure evolves to:**
```
ThreeDoors/
├── cmd/                    # CLI entry points
│   └── threedoors/        # Main application
├── internal/              # Private application code
│   ├── core/             # Core domain logic
│   ├── tui/              # Bubbletea interface components
│   ├── integrations/     # Adapter implementations
│   │   ├── textfile/    # Text file backend (from Tech Demo)
│   │   └── applenotes/  # Apple Notes integration
│   ├── enrichment/       # Local enrichment storage
│   └── learning/         # Door selection & pattern tracking
├── pkg/                   # Public, reusable packages (if any)
├── docs/                  # Documentation (including this PRD)
├── .bmad-core/           # BMAD methodology artifacts
└── Makefile              # Build automation
```

### Service Architecture

**Technical Demo Phase:**

**Decision:** Single-layer CLI/TUI application with direct file I/O

**Rationale:**
- **No abstractions yet**: Build for one thing (text files), refactor when adding second thing
- **Validate UX first**: Door selection algorithm is the innovation, not the data layer
- **Fast iteration**: Change anything without navigating architecture layers

**Demo Architecture:**
- **TUI Layer (Bubbletea)** - Three Doors interface, keyboard handling, rendering
- **Direct File I/O** - Read tasks.txt, write completed.txt, no abstraction layer
- **Simple Door Selection** - Random selection of 3 tasks from available pool (no learning/categorization yet)

---

**Full MVP Phase (Post-Validation - Deferred):**

**Decision:** Layered architecture with pluggable integration adapters

**Architecture Layers:**
1. **TUI Layer (Bubbletea)** - User interaction, rendering, keyboard handling
2. **Core Domain Logic** - Task management, door selection algorithm, progress tracking
3. **Integration Adapters** - Abstract interface with concrete implementations (text file, Apple Notes, others later)
4. **Enrichment Storage** - Metadata, cross-references, learning patterns not stored in source systems
5. **Configuration & State** - User preferences, values/goals, application state

**Key Architectural Principles:**
- Core domain logic has NO dependencies on specific integrations (dependency inversion)
- Integrations implement common `TaskProvider` interface
- Enrichment layer wraps tasks from any source with additional metadata
- TUI layer depends only on core domain, not specific integrations

### Testing Requirements

**Technical Demo Phase:**

**Decision:** Manual testing only - validate UX through real use

**Demo Testing Approach:**
- **No automated tests for Tech Demo** - premature given throwaway prototype nature
- **Manual testing** via daily use for 1 week
- **Success measurement**: Does Three Doors feel better than a list? Yes/No decision point
- **Quality gate**: If it crashes or feels bad to use, iterate or abandon concept

**Rationale:**
- 4-8 hours to build entire demo - testing infrastructure would consume half that time
- Real usage is the test: if developer won't use it daily, concept fails regardless of test coverage
- Can add tests when/if proceeding to Full MVP

---

**Full MVP Phase (Post-Validation - Deferred):**

**Testing Scope:**
- **Unit tests** for core domain logic (door selection algorithm, categorization, progress tracking)
- **Integration tests** for backend adapters (text file, Apple Notes)
- **Manual testing** for TUI interactions (Bubbletea testing framework is immature)

**Test Coverage Goals:**
- Core domain logic: 70%+ coverage (pragmatic, not perfectionist)
- Integration adapters: Critical paths covered (read, write, sync scenarios)
- TUI layer: Manual testing via developer use

**Testing Strategy:**
- Table-driven tests (idiomatic Go pattern)
- Test fixtures for data structures
- Mock `TaskProvider` interface for testing core logic without real integrations
- CI/CD runs tests on every commit (GitHub Actions)

**Deferred for Post-MVP:**
- End-to-end testing framework
- Property-based testing for door selection algorithm
- Performance/load testing

### Additional Technical Assumptions and Requests

**Technical Demo Phase Assumptions:**

**Text File Format:**
- **Simple line-delimited format**: One task per line in `tasks.txt`
- **Completed format**: `[timestamp] task description` in `completed.txt`
- **No metadata yet**: Task is just text; no categories, priorities, or context for Tech Demo
- **Easy population**: User can edit files with any text editor, generate from scripts, copy-paste, etc.

**Door Selection Algorithm (Tech Demo):**
- **Random selection**: Pick 3 random tasks from available pool
- **Simple diversity**: Ensure no duplicates in the three doors
- **No intelligence yet**: No learning, no categorization, no context awareness
- **Validation goal**: Prove that having 3 options reduces friction vs. scrolling a full list

**File I/O:**
- **Go standard library**: Use `os`, `bufio`, `io/ioutil` - no external dependencies for file operations
- **Error handling**: Create files with defaults if missing; graceful degradation if corrupted
- **Concurrency**: Not a concern for single-user local files

---

**Full MVP Phase Assumptions (Post-Validation - Deferred):**

**Apple Notes Integration:**
- **Options Identified (2025):**
  1. **DarwinKit (github.com/progrium/darwinkit)** - Native macOS API access from Go; requires translating Objective-C patterns; full API control but higher complexity
  2. **Direct SQLite Database Access** - Apple Notes stores data in `~/Library/Group Containers/group.com.apple.notes/NoteStore.sqlite`; note content is gzipped protocol buffers in `ZICNOTEDATA.ZDATA` column; read-only safe, write risks corruption
  3. **AppleScript Bridge** - Use `os/exec` to invoke AppleScript; simpler than native APIs; proven approach (see `sballin/alfred-search-notes-app`)
  4. **Existing MCP Server** - `mcp-apple-notes` server exists for Apple Notes integration; could potentially leverage this instead of building from scratch
- **Assumption:** Multiple viable paths exist; choice depends on read-only vs. read-write needs, complexity tolerance, and reliability requirements (WILL REQUIRE VALIDATION when implementing Phase 2)
- **Spike Required:** Evaluate options before implementing Apple Notes integration
- **Preferred Exploration Order:** Start with Option 4 (MCP server) or Option 2 (SQLite read-only), fall back to Option 3 (AppleScript) if bidirectional sync required, reserve Option 1 (DarwinKit) for complex scenarios

**Cloud Storage for Cross-Computer Sync (DEFERRED - Not MVP):**
- **Status:** Cross-computer sync is deferred post-MVP; single-computer local storage is sufficient for initial development and use
- **Future Exploration:** When implementing sync, explore alternatives to monolithic SQLite file:
  - Individual JSON/YAML files per task or per day (more granular, better suited for file-based cloud sync)
  - Conflict-free Replicated Data Types (CRDTs) for eventual consistency
  - Event sourcing with append-only logs
  - Cloud-native solutions (S3, Firebase, etc.) if local-first constraint relaxes
- **Awareness:** Monolithic SQLite on cloud storage (iCloud/Google Drive) is known problematic—corruption risk, locking issues, slow sync
- **MVP Decision:** Store enrichment data locally only; revisit sync architecture when/if multi-computer use becomes actual need

**Go Language & Ecosystem (Tech Demo):**
- **Language:** Go 1.25.4+ (current stable as of November 2025)
- **Formatting:** `gofumpt` (run before commits)
- **Linting:** Skip for Tech Demo (adds no validation value at this stage)
- **Dependency Management:** Go modules
- **TUI Framework:** Bubbletea + Lipgloss (styling) - minimal Bubbles components, only if needed

**Data Storage (Tech Demo):**
- **Storage:** Plain text files in `~/.threedoors/`
- **No database**: Not needed for line-delimited text
- **No configuration file initially**: Hardcode paths, add config only if needed

**Build & Development (Tech Demo):**
- **Build System:** Minimal Makefile
  ```makefile
  build:
      go build -o bin/threedoors cmd/threedoors/main.go

  run: build
      ./bin/threedoors

  clean:
      rm -rf bin/
  ```
- **Development Workflow:** Direct iteration on macOS
- **No CI/CD for Tech Demo**: Overkill for validation prototype

**Performance Expectations (Tech Demo):**
- **File I/O**: <10ms to read tasks.txt (even with 100+ tasks)
- **Door selection**: <1ms for random selection from array
- **TUI rendering**: Bubbletea handles 60fps, not a concern
- **Startup time**: <100ms total from launch to Three Doors display

**Security & Privacy (Tech Demo):**
- **Local files only**: No network, no external services
- **No logging**: Not even metadata for Tech Demo
- **File permissions**: Standard user file permissions on `~/.threedoors/`

---

**Full MVP Phase (Post-Validation - Deferred):**

**Go Language & Ecosystem:**
- **Language:** Go 1.25.4+
- **Formatting:** `gofumpt`
- **Linting:** `golangci-lint` with standard rule set
- **Dependency Management:** Go modules
- **TUI Framework:** Bubbletea + Lipgloss + Bubbles

**Data Storage:**
- **Primary:** Apple Notes (user-facing tasks) or text file backend
- **Enrichment:** SQLite for metadata (door feedback, blockers, categorization, learning patterns)
- **Configuration:** YAML or TOML for user preferences, values/goals
- **Location:** `~/.config/threedoors/` (XDG Base Directory spec on Linux, macOS equivalent)

**Build & Development:**
- **Build System:** Makefile with full targets (build, test, lint, install)
- **CI/CD:** GitHub Actions running tests on every commit
- **Development Workflow:** Direct iteration on macOS

**Performance Expectations:**
- Door selection algorithm: <100ms to choose 3 tasks from up to 1000 total tasks
- Backend sync: <2 seconds for typical data set
- TUI rendering: 60fps equivalent for smooth interaction

**Deferred Technical Decisions (Post-MVP):**
- Cross-computer sync architecture (see deferred section above)
- LLM provider integration architecture (local vs. cloud, which providers)
- Additional integration adapters (Jira, Linear, Google Calendar, etc.)
- Remote access agent for Geodesic environments
- Vector database for semantic task search
- Voice interface integration

---
