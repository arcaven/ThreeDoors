# Sprint Change Proposal: Envoy Agent Implementation Infrastructure

**Date:** 2026-03-08
**Change Trigger:** Envoy agent definition merged (PRs #227, #232) with party mode research complete, but supporting infrastructure not yet built
**Change Scope:** Minor — Direct Adjustment (add stories to Epic 0)
**Approved By:** Supervisor (task assignment)

---

## Section 1: Issue Summary

The envoy agent definition (`agents/envoy.md`) and its comprehensive rules of behavior (`_bmad-output/planning-artifacts/envoy-rules-of-behavior-party-mode.md`) were produced and merged in PRs #227 and #232. However, the envoy cannot operate without supporting infrastructure:

- **No local issue tracker file** — The envoy's primary state file (`docs/issue-tracker.md`) doesn't exist yet
- **No authority tier configuration** — Owner/contributor tiers need to be embedded in the tracker header
- **No SOUL.md alignment patterns** — The envoy needs documented classification logic for issue triage
- **No integration documentation** — How the envoy interacts with merge-queue, pr-shepherd, and supervisor isn't codified beyond the agent prompt

This is operationalizing existing, fully-researched work — not new feature development.

## Section 2: Impact Analysis

### Epic Impact
- **Epic 0 (Infrastructure):** Add 2 new stories (0.28, 0.29). Currently 19/22 complete → becomes 19/24.
- **No other epics affected.** The envoy is process infrastructure, not Go application code.

### Story Impact
- No existing stories modified.
- No dependencies on in-progress stories.
- Stories 0.28 and 0.29 are sequential (0.29 depends on 0.28).

### Artifact Conflicts
- **PRD:** No changes. The envoy is development infrastructure, not a product requirement.
- **Architecture:** No changes. The envoy is a multiclaude agent, not a Go component.
- **UI/UX:** No changes. No user-facing impact.

### Technical Impact
- **docs/issue-tracker.md** — New file with tracker structure, authority tiers, metadata format
- **docs/envoy-operations.md** — New file documenting envoy operational patterns (SOUL.md alignment, staleness thresholds, cross-agent protocols)
- **No code changes.** No Go files modified.

## Section 3: Recommended Approach

**Selected: Direct Adjustment — Add stories to Epic 0**

**Rationale:**
- The research is 100% complete (5-round party mode with full team consensus)
- All decisions are documented in the party mode artifact
- The work is purely documentation/infrastructure creation
- Low effort (creating files from existing specifications), low risk
- No rollback or scope changes needed

**Effort:** Low (2 stories, docs-only)
**Risk:** Low (no code changes, no breaking changes)
**Timeline Impact:** None on existing work

## Section 4: Detailed Change Proposals

### Epic 0 Status Update

```
OLD:
**Status:** 19 of 22 stories complete. Stories 0.20 (CI Churn Reduction), 0.21 (Homebrew Public Distribution), and 0.24 (Renovate + Dependabot) not started.

NEW:
**Status:** 19 of 24 stories complete. Stories 0.20 (CI Churn Reduction), 0.21 (Homebrew Public Distribution), 0.24 (Renovate + Dependabot), 0.28 (Issue Tracker & Authority Config), and 0.29 (Envoy Operations Guide) not started.
```

### New Stories

- **Story 0.28:** Issue Tracker File Structure, Authority Configuration & Initial Content
- **Story 0.29:** Envoy Operations Guide & Integration Documentation

### ROADMAP.md Update

Add envoy stories to Infrastructure Backlog section.

## Section 5: Implementation Handoff

**Change Scope:** Minor — Direct implementation by workers via `/implement-story`

**Handoff:**
- Story 0.28 → Worker agent (creates `docs/issue-tracker.md`)
- Story 0.29 → Worker agent (creates `docs/envoy-operations.md`)
- Stories are sequential: 0.29 depends on 0.28

**Success Criteria:**
- `docs/issue-tracker.md` exists with correct format per party mode consensus
- Authority tiers configured with `arcaven` as owner
- SOUL.md alignment patterns documented for envoy use
- Staleness thresholds, cross-agent protocols, and patrol workflow documented
- Envoy agent can be spawned and begin patrol cycles

---

## References

- Envoy agent definition: `agents/envoy.md` (PRs #227, #232)
- Party mode research: `_bmad-output/planning-artifacts/envoy-rules-of-behavior-party-mode.md`
- SOUL.md: Project values referenced for alignment classification
