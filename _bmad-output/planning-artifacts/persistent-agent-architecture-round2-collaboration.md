# Party Mode Artifact: Agent Collaboration & Propagation Chain (Round 2)

**Date:** 2026-03-08
**Topic:** How should persistent BMAD agents collaborate?
**Participants:** John (PM), Winston (Architect), Bob (SM), Murat (TEA)

## Adopted Approach

### Communication Model: Message-Driven Chain

Agents communicate via `multiclaude message send`, not shared file polling. Each agent has an independent monitoring loop AND message reactivity.

### Propagation Chain: PR Merge → Doc Cascade

```
PR Merged (detected by PM in polling loop)
  │
  ├── PM updates story status → Done (PR #NNN)
  ├── PM updates ROADMAP.md → Epic progress
  ├── PM checks PRD alignment
  │     └── If PRD drift detected:
  │           PM messages Architect: "PRD section X changed, verify architecture alignment"
  │           └── Architect reviews architecture docs
  │                 └── If architecture update needed:
  │                       Architect updates docs directly
  │                       Architect messages PM: "Architecture updated, stories may need tech note refresh"
  │                       └── PM flags affected stories for next worker
  │
  └── merge-queue handles the merge itself (existing)
```

### Architect Independent Loop

```
Architect polls recent merges (every 20-30 min)
  │
  ├── Checks code changes in internal/ against docs/architecture/
  │     └── If pattern divergence:
  │           Opens GitHub issue with details
  │           Messages PM: "Architecture drift detected, see issue #NNN"
  │
  └── Checks for undocumented new packages/interfaces
        └── If found: flags for documentation
```

### Anti-Patterns Avoided

1. **Circular notifications:** Each message includes a correlation ID (PR number). Agents skip already-processed PRs.
2. **Authority creep:** PM edits planning docs. Architect edits architecture docs. Neither creates stories or modifies code — those require worker spawning.
3. **Chatty agents:** Messages are sent only on state changes, never as status pings.

### Authority Boundaries

| Agent | Can Directly Edit | Must Spawn Worker | Must Escalate to Supervisor |
|-------|-------------------|-------------------|----------------------------|
| PM | Story files, ROADMAP.md | New story creation | Scope decisions, priority changes |
| Architect | Architecture docs | Code refactoring | Design decision overrides |
| SM (cron) | Sprint status doc | Nothing | Blocked items, risk alerts |
| QA (cron) | Coverage reports | Test improvements | Coverage policy changes |

### Cron-Based Agent Integration

- **SM cron (every 4 hours):** Queries PR status, worker status, story progress. Generates summary. Messages supervisor if risks detected.
- **QA cron (weekly):** Runs coverage analysis. Compares to baseline. Messages PM if regression detected.
- Both can be implemented via `multiclaude` `/loop` skill or external cron.

## Rejected Options

### Shared File Protocol (agents write to a shared state file)
- **Why rejected:** Race conditions in shared worktrees. Message passing is the established multiclaude pattern. File-based coordination adds complexity without benefit.

### Event-Driven Architecture (webhook-based)
- **Why rejected:** multiclaude doesn't support webhook triggers. Polling + messages is the available primitive. Over-engineering for the scale of this project.

### Dense Agent Mesh (every agent talks to every other agent)
- **Why rejected:** Combinatorial explosion. With N agents, N*(N-1) communication channels. The hub-and-spoke model (PM as hub) is simpler and sufficient.

### PM as Single Hub for All Communication
- **Why rejected (partially):** PM is the primary hub, but Architect needs an independent monitoring loop for code-level concerns that PM can't assess. Two independent loops with message bridges is better than one monolithic loop.
