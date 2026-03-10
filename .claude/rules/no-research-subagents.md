# Rule: No Research via Subagents

NEVER use the Agent tool for research tasks. The word "research", "investigate",
"analyze", or "evaluate" in a task description means it MUST be dispatched via
`multiclaude work`, not the Agent tool.

The Agent tool is ONLY for single-question codebase lookups that return 1-3
sentences (e.g., "find where TaskProvider is defined", "what does this function
return?"). If the answer requires reading more than 5 files or synthesizing
information across multiple sources, use `multiclaude work`.

**Decision heuristic:** "What is X?" → Agent OK. "What should we do about X?" → multiclaude work.

**Why this matters:** Subagent research tasks consume supervisor context window,
are invisible to the user in tmux, and violate the multiclaude architecture
where all substantive work should be visible and trackable.
