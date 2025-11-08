# Epic List

## Phase 1: Technical Demo & Validation (Immediate - Week 1)

**Epic 1: Three Doors Technical Demo**
- **Goal:** Build and validate the Three Doors interface with minimal viable functionality to prove the UX concept reduces friction compared to traditional task lists
- **Timeline:** 1 week (3-6 hours development time - optimized sequence)
- **Deliverables:** Working CLI/TUI showing Three Doors, reading from text file, door refresh, marking tasks complete, basic progress tracking
- **Success Criteria:**
  - Developer uses tool daily for 1 week
  - Three Doors selection feels meaningfully different from scrolling a list
  - Decision point reached: proceed to Full MVP or pivot/abandon
- **Tech Stack:** Go 1.25.4+, Bubbletea/Lipgloss, local text files
- **Risk:** UX concept might not feel better than simple list; easy to pivot if fails
- **Optimization:** Reordered stories to validate refresh UX before completion; merged/simplified non-essential features

---

## Phase 2: Post-Validation Roadmap (Conditional on Phase 1 Success)

**DECISION GATE:** Only proceed with these epics if Technical Demo validates the Three Doors concept through real usage.

**Epic 2: Foundation & Apple Notes Integration**
- **Goal:** Replace text file backend with Apple Notes integration, enabling mobile task editing while maintaining Three Doors UX
- **Prerequisites:** Epic 1 success; Apple Notes integration spike completed
- **Deliverables:**
  - Refactor to adapter pattern (text file + Apple Notes backends)
  - Bidirectional sync with Apple Notes
  - Health check command for Notes connectivity
  - Migration path from text files to Notes
- **Estimated Effort:** 3-4 weeks at 2-4 hrs/week (includes spike + implementation)
- **Risk:** Apple Notes integration complexity could exceed estimates; fallback to improved text file backend

**Epic 3: Enhanced Interaction & Task Context**
- **Goal:** Add task capture, values/goals display, and basic feedback mechanisms to improve task management workflow
- **Prerequisites:** Epic 2 complete (stable backend integration)
- **Deliverables:**
  - Quick add mode for task capture
  - Extended capture with "why" context
  - Values/goals setup and persistent display
  - Door feedback options (Blocked, Not now, Needs breakdown)
  - Blocker tracking
  - Improvement prompt at session end
- **Estimated Effort:** 2-3 weeks at 2-4 hrs/week
- **Risk:** Feature creep; maintain focus on minimal valuable additions

**Epic 4: Learning & Intelligent Door Selection**
- **Goal:** Implement pattern tracking and learning to make door selection context-aware and adaptive to user preferences
- **Prerequisites:** Epic 3 complete (enough usage data to learn from)
- **Deliverables:**
  - Task categorization (type, effort level, context)
  - Door selection pattern tracking
  - Learning algorithm that adapts based on user choices
  - Progress view showing door choice patterns over time
  - "Better than yesterday" multi-dimensional tracking
- **Estimated Effort:** 3-4 weeks at 2-4 hrs/week
- **Risk:** Algorithm complexity; may need to simplify learning approach

**Epic 5: Data Layer & Enrichment (Optional)**
- **Goal:** Add enrichment storage layer for metadata that cannot live in source systems
- **Prerequisites:** Epic 4 complete; proven need for enrichment beyond what backends support
- **Deliverables:**
  - SQLite enrichment database
  - Cross-reference tracking (tasks across multiple systems)
  - Metadata not supported by Apple Notes (categories, learning patterns, etc.)
  - Data migration and backup tooling
- **Estimated Effort:** 2-3 weeks at 2-4 hrs/week
- **Risk:** May be YAGNI; consider deferring indefinitely if not clearly needed

---

## Phase 3: Future Expansion (12+ months out)

**Epic 6+: Additional Integrations** (Jira, Linear, Google Calendar, Slack, etc.)
**Epic 7+: Cross-Computer Sync** (Implement alternative to monolithic SQLite on cloud storage)
**Epic 8+: LLM Integration** (Task breakdown assistance, assumption challenging, dependency collapse)
**Epic 9+: Advanced Features** (Voice interface, mobile app, web interface, trading mechanic, gamification)

**Guiding Principle:** Each epic must deliver tangible user value and be informed by real usage patterns from previous phases. No speculation-driven development.

---
