# Requirements

## Technical Demo & Validation Phase Requirements

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

## Full MVP Requirements (Post-Validation - Deferred)

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

## Non-Functional Requirements

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
