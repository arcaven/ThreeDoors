# Party Mode Artifact: Persistent Agent Role Evaluation (Round 1)

**Date:** 2026-03-08
**Topic:** Which BMAD roles should become persistent multiclaude agents?
**Participants:** John (PM), Winston (Architect), Bob (SM), Quinn (QA), Murat (TEA), Mary (Analyst), Paige (Tech Writer), Sally (UX Designer)

## Adopted Approach

### Persistent Agents (Always-On with Polling)

1. **PM Agent ("project-watchdog")** — Strongest case
   - **Monitor:** Merged PRs, story completion status, PRD alignment, ROADMAP.md accuracy, story sequencing
   - **Trigger:** PR merge events (poll `gh pr list --state merged`), message from other agents
   - **Authority:** Direct updates to story files, ROADMAP.md. Flag PRD drift. Message architect for cascading changes. Cannot create stories or modify code.
   - **Polling interval:** Every 10-15 minutes
   - **Rationale:** Every merged PR potentially drifts planning docs. This is the highest-frequency governance gap.

2. **Architect Agent ("arch-watchdog")** — Second strongest case
   - **Monitor:** Code changes in `internal/` vs architecture docs, new patterns introduced, architectural debt
   - **Trigger:** Code changes detected in recent merges, messages from PM about PRD changes
   - **Authority:** Direct updates to architecture docs. Flag code divergence via issues. Cannot refactor code directly.
   - **Polling interval:** Every 20-30 minutes
   - **Rationale:** Code-to-doc divergence accumulates silently. 210+ PRs means significant drift potential.

### Cron/Periodic (Not Persistent)

3. **SM (Sprint Health)** — Every 4 hours
   - Sprint status summary: blocked stories, stale PRs, idle workers
   - Overlaps significantly with merge-queue and pr-shepherd functions
   - Value is in summarization, not continuous monitoring

4. **QA/TEA (Coverage Audit)** — Weekly
   - Run `go test -cover ./...`, compare to baseline
   - Flag regressions to PM
   - CI already catches per-PR issues; this catches trend drift

5. **Tech Writer (Doc Staleness)** — Weekly
   - Check last-modified dates vs code changes
   - Flag docs that may be stale
   - Low frequency of doc changes makes persistence wasteful

### Stay Ephemeral (On-Demand Only)

6. **Analyst** — Research sweeps are monthly at most; fold into PM's responsibilities
7. **UX Designer** — CLI/TUI changes are story-driven, always ephemeral
8. **Dev** — Always ephemeral, spawned per-story via `/implement-story`

## Rejected Options

### All Roles Persistent
- **Why rejected:** Resource overhead (API costs, tmux sessions, compute) for 8+ persistent agents is untenable. Most roles don't generate enough monitoring events to justify always-on status.

### SM as Persistent
- **Why rejected:** Merge-queue and pr-shepherd already handle the mechanical aspects of process health (PR status, CI failures, rebasing). The SM's value is in periodic summarization, not continuous monitoring. A cron job every 4 hours provides equivalent value at 1/10th the cost.

### QA as Persistent
- **Why rejected:** CI runs on every PR and catches quality issues at the PR level. Persistent QA would mostly be idle. Coverage trend monitoring is a weekly concern, not a minute-by-minute one.

### Analyst as Persistent
- **Why rejected:** Research findings accumulate slowly (days/weeks between new research docs). A monthly sweep is sufficient. The PM can absorb this function as part of its monitoring loop.

### Tech Writer as Persistent
- **Why rejected:** Documentation drift happens over weeks, not hours. A weekly cron audit achieves the same result as persistence.

### UX Designer as Persistent
- **Why rejected:** Zero monitoring surface for a CLI/TUI project. UX decisions are made during story planning, not discovered through continuous monitoring.
