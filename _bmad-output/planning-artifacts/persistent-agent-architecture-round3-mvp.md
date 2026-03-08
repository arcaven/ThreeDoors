# Party Mode Artifact: MVP Persistent Agent Selection (Round 3)

**Date:** 2026-03-08
**Topic:** MVP — Top 2-3 persistent agents for maximum autonomous governance
**Participants:** John (PM), Winston (Architect), Bob (SM), Murat (TEA), Mary (Analyst), Barry (Quick Flow)

## Adopted Approach: Two New Persistent Agents

### MVP: PM + Architect (Total: 5 Persistent Agents)

Current persistent agents (3):
1. merge-queue — merges PRs
2. pr-shepherd — rebases branches
3. envoy — issue triage

**Add (2):**

4. **project-watchdog (PM role)** — Planning-side governance
   - Watches merged PRs, updates story status and ROADMAP.md
   - Detects PRD drift
   - Validates story sequencing
   - Monthly research doc sweep (absorbs analyst function)
   - Polling: every 10-15 minutes

5. **arch-watchdog (Architect role)** — Implementation-side governance
   - Watches code changes vs architecture docs
   - Detects undocumented patterns
   - Flags architectural debt
   - Polling: every 20-30 minutes

### Supporting Cron Jobs (Not Persistent)

- **Sprint health (SM):** Every 4 hours via `/loop`
- **Coverage audit (QA/TEA):** Weekly via cron or scheduled worker

### Why Two, Not Three

- 5 total persistent agents is manageable for compute/API cost
- 6+ starts creating coordination overhead that outweighs governance value
- The two selected agents cover the two most critical governance gaps: planning-side drift and implementation-side drift
- Start with 2, observe for 2 weeks, add a third only if a clear gap emerges

### Projected Impact

| Problem | Solved By | Coverage |
|---------|-----------|----------|
| PRD drift unnoticed | project-watchdog | ~90% |
| Architecture divergence | arch-watchdog | ~80% |
| Stories out of sequence | project-watchdog | ~70% |
| Research findings unactioned | project-watchdog (monthly sweep) | ~50% |
| Test coverage regression | QA cron (weekly) | ~60% |
| Doc staleness | Tech writer cron (weekly) | ~50% |
| Sprint health blindness | SM cron (4-hourly) | ~70% |

## Rejected Options

### Three Persistent Agents (PM + Architect + SM)
- **Why rejected:** SM's monitoring overlaps with merge-queue and pr-shepherd. The marginal value of persistent SM doesn't justify the resource cost. A 4-hourly cron achieves 70%+ of the value at minimal cost.

### Three Persistent Agents (PM + Architect + QA)
- **Why rejected:** QA monitoring is episodic (per-PR via CI, weekly via audit). No continuous monitoring surface justifies persistence. Weekly cron is sufficient and much cheaper.

### One Persistent Agent (PM only)
- **Why rejected:** Missing the code-to-docs feedback loop. The PM can't assess whether code changes comply with architecture docs — that requires architect domain knowledge. Without the architect agent, architectural drift continues unchecked.

### Zero New Persistent Agents (All Cron-Based)
- **Why rejected:** Cron jobs lack the contextual awareness and message-reactivity of persistent agents. A cron PM would miss the cascade: "PR merged → story updated → ROADMAP updated → PRD checked → architect notified" because each step depends on the previous. Persistent agents maintain state across these cascading checks.

### Deferred for Future Evaluation
- **Tech Writer persistent agent:** Revisit after 1 month if doc staleness remains a problem
- **SM persistent agent:** Revisit if sprint health summaries from cron prove insufficient
- **Analyst persistent agent:** Revisit if research doc backlog grows beyond PM's monthly sweep capacity
