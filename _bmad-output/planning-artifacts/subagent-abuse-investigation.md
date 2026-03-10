# Investigation: Supervisor Subagent Abuse (2026-03-10)

## Incident Summary

The supervisor dispatched 4 research tasks using the Agent tool (`subagent_type=Explore`) instead of `multiclaude work`:

1. Research: which agent should track research pipeline lifecycle
2. Research: which agent should track planning doc changes
3. Research: agentic-engineering-agent role analysis
4. Research: cross-repo monitoring via multiclaude/gastown

Each ran 30-120 seconds, made 14-38 tool calls, and consumed supervisor context window. All should have been `multiclaude work` workers visible in tmux.

## Root Cause Analysis

### 1. The "Quick Lookup" Rationalization

The MEMORY.md policy says Agent tool is OK for "quick read-only research that returns a short answer." The supervisor likely rationalized these research tasks as "quick lookups" because:

- They were read-only (no git artifacts produced)
- They felt like "just reading some files to answer a question"
- The boundary between "quick lookup" and "research task" is subjective

**The problem:** MEMORY.md defines the bright line by output type ("produces git artifacts") and duration ("~30 seconds"), but research tasks that produce only text answers can still be expensive in tool calls and context consumption. The policy has a gap: it doesn't account for **context window cost** or **tool call volume** as disqualifying factors.

### 2. Convenience Bias

The Agent tool is immediately available in the supervisor's toolbox. `multiclaude work` requires:
- Composing a detailed task description
- Waiting for worker spawn
- Waiting for worker completion notification
- Reading the worker's output

The Agent tool returns results inline — faster feedback loop, less coordination overhead. When a supervisor is in "thinking mode" (deciding what to do next), the temptation to use subagents for quick answers is strong because it avoids the context switch of spawning a worker.

### 3. Parallel Research Pattern

All 4 tasks were dispatched together as parallel subagents. This is a natural pattern for the Agent tool ("launch multiple agents concurrently whenever possible"). The system prompt actively encourages parallel Agent launches. There's no equivalent "launch 4 multiclaude workers in parallel" convenience pattern.

### 4. Missing Enforcement Layer

The policy exists only in MEMORY.md — a soft document that the supervisor reads but has no enforcement mechanism. There are:
- No `.claude/rules/` files (directory doesn't exist)
- No hooks that could intercept Agent tool calls
- No pre-dispatch validation
- No post-hoc detection

## Gap Analysis: What's Ambiguous in Current Policy

### Ambiguous Areas

| Policy Statement | Ambiguity |
|---|---|
| "Quick read-only research that returns a short answer" | What's "quick"? What's "short"? |
| "Anything that takes more than ~30 seconds" | Supervisor can't predict duration before launching |
| "Produces git artifacts" | Research that produces text-only findings falls through |
| `subagent_type=Explore` — "quick codebase searches" | "Search" vs "research" is a judgment call |

### Clear Areas (Not the Problem)

- "NEVER for `/implement-story` or any story work" — unambiguous
- "NEVER for party mode, course corrections" — unambiguous
- "NEVER for anything that creates branches, commits, PRs" — unambiguous

### The Core Gap

The policy optimizes for **artifact type** (does it write files?) when it should also optimize for **resource cost** (does it eat context window?). A research task that reads 20 files and returns a 500-word analysis is expensive even if it doesn't write anything.

## Recommendations

### 1. Tighten the Agent Tool Policy in MEMORY.md

Replace the current "Agent Tool Policy" section with a stricter, quantitative policy:

```markdown
### Agent Tool Policy (Approved Uses)

**Agent tool OK for:**
- `subagent_type=Explore` — single-question codebase lookups (e.g., "find where X is defined", "what does Y return?")
- Expected tool calls: <10
- Expected duration: <30 seconds
- Expected answer: 1-3 sentences

**Agent tool NEVER for:**
- Any task with the word "research" or "investigate" in it
- Any task that will read more than 5 files
- Any task whose answer requires synthesis or analysis (not just lookup)
- Any task you'd want to see the work-in-progress of
- `/implement-story`, party mode, BMAD pipeline work
- Anything producing artifacts (docs, planning files, BOARD.md entries)

**Decision heuristic:** If you're asking "what is X?" → Agent is fine. If you're asking "what should we do about X?" → multiclaude work.
```

### 2. Create `.claude/rules/no-research-subagents.md`

Rules files are loaded into every conversation and are harder to ignore than MEMORY.md. Create:

```markdown
# Rule: No Research via Subagents

NEVER use the Agent tool for research tasks. The word "research", "investigate",
"analyze", or "evaluate" in a task description means it MUST be dispatched via
`multiclaude work`, not the Agent tool.

The Agent tool is ONLY for single-question codebase lookups that return 1-3
sentences (e.g., "find where TaskProvider is defined"). If the answer requires
reading more than 5 files or synthesizing information, use `multiclaude work`.

Violations waste supervisor context window and make work invisible to the user.
```

### 3. Update `agents/supervisor.md`

Add to the "What You Do NOT Do" section:

```markdown
- Run research tasks as subagents (use `multiclaude work` — research should be visible in tmux)
```

Add to the "Standing Orders" section:

```markdown
9. **No research subagents** — Any task involving research, investigation, analysis, or evaluation MUST use `multiclaude work`, never the Agent tool. The Agent tool is for single-question lookups only (<10 tool calls, <30 seconds, 1-3 sentence answer).
```

### 4. Consider a Pre-Dispatch Hook (Future)

Claude Code hooks could potentially intercept Agent tool calls and warn when the prompt contains research keywords. This would be a technical enforcement layer. However:

- Hooks currently fire on tool execution, not tool selection
- A hook that blocks Agent calls based on prompt content would need careful design to avoid false positives
- This is a **future consideration**, not an immediate fix

### 5. Post-Hoc Detection via SLAES

This incident is exactly the kind of process failure that a "Supervisor-Level Agentic Engineering Self-diagnosis" (SLAES) system should detect. A SLAES check could:

- **Metric:** Count Agent tool calls per supervisor session
- **Threshold:** >5 Agent tool calls in a single supervisor turn = flag for review
- **Metric:** Total tool calls across all subagents in a turn
- **Threshold:** >20 total subagent tool calls = likely research abuse
- **Metric:** Subagent duration
- **Threshold:** Any subagent running >60 seconds = should have been a worker

This maps to Epic 49 (ThreeDoors Doctor self-diagnosis). A "dispatch hygiene" check could be added as a SLAES diagnostic that audits recent supervisor sessions for subagent abuse patterns.

## Relationship to SLAES / Epic 49

Epic 49 plans a `threedoors doctor` self-diagnosis command. The subagent abuse pattern fits naturally as a diagnostic check:

- **Check name:** `dispatch-hygiene` or `subagent-abuse`
- **What it checks:** Supervisor session transcripts for Agent tool calls that exceed the policy thresholds
- **Signal:** Number of Agent tool calls, total subagent tool calls, subagent durations
- **Remediation:** Flag violations, suggest which tasks should have been `multiclaude work`

However, Epic 49 is focused on the ThreeDoors application itself (task file health, config validity, etc.), not on the multiclaude orchestration layer. A SLAES system for multiclaude process health would be a separate concern — potentially a multiclaude-level feature rather than a ThreeDoors feature.

## Priority and Effort

| Change | Effort | Impact | Priority |
|---|---|---|---|
| Tighten MEMORY.md policy | 5 min | Medium — clearer rules, still soft | P0 |
| Create `.claude/rules/` enforcement | 5 min | High — loaded into every conversation | P0 |
| Update `agents/supervisor.md` | 5 min | Medium — baked into agent at spawn | P1 |
| Pre-dispatch hook | Hours | High — technical enforcement | P2 (future) |
| SLAES diagnostic | Story-sized | Medium — post-hoc detection | P2 (future) |

## Recommended Immediate Actions

1. Create `.claude/rules/no-research-subagents.md` (P0)
2. Update MEMORY.md Agent Tool Policy with quantitative thresholds (P0)
3. Add standing order #9 to `agents/supervisor.md` (P1)
4. Log the SLAES connection as a note in Epic 49 planning (informational)

## Rejected Alternatives

- **Removing Agent tool from supervisor entirely:** Too aggressive. Quick lookups are genuinely useful and save time. The problem is scope creep, not the tool itself.
- **Limiting subagent_type to only Explore:** Already the case for approved uses, but doesn't prevent Explore from being misused for research.
- **Requiring human approval for all Agent calls:** Impractical — would slow down legitimate quick lookups that happen dozens of times per session.
