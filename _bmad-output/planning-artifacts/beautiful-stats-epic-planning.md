# Beautiful Stats Epic Planning — Party Mode Artifact

**Date:** 2026-03-09
**Type:** Epic planning
**Epic:** 39 — Beautiful Stats Display
**Source Research:** `_bmad-output/planning-artifacts/beautiful-stats-research.md`

---

## Decisions Summary

| ID | Decision | Rationale | Alternatives Rejected |
|----|----------|-----------|----------------------|
| 1 | Use Lipgloss gradients for sparkline coloring | Already in project dependencies; `Blend1D()` available in current version; zero new dependencies | ntcharts sparkline (unnecessary dependency for something simple), raw ANSI codes (non-idiomatic) |
| 2 | Lipgloss bordered panels for dashboard layout | `NewStyle().Border()` + `JoinHorizontal()`/`JoinVertical()` already available; consistent with project style | termdash (incompatible event loop), custom box-drawing (reinvents Lipgloss) |
| 3 | Fun facts from existing PatternAnalyzer data | Data already captured in sessions.jsonl; no new tracking needed | New analytics engine (overengineering), hardcoded facts (defeats purpose) |
| 4 | Custom in-tree horizontal bar charts | Simple proportional fill with Lipgloss styling; ~30 lines of code | ntcharts bar chart (overkill for proportional bars), asciigraph (line graphs, not bars) |
| 5 | Evaluate ntcharts for heatmap only | Heatmap is the one chart type complex enough to justify a dependency; sparklines and bars are simpler in-tree | Build heatmap from scratch (significant effort for grid + color math), skip heatmap entirely (loses most impactful visualization) |
| 6 | `tea.Tick`-based counter animation | Native Bubbletea pattern; already used in project for worker polling | CSS-style transitions (not available in TUI), instant display (misses "button feel" opportunity) |
| 7 | Milestone celebrations are one-time observations only | SOUL.md anti-gamification principle; celebrate what happened, never prescribe what should happen | Badge/achievement system (gamification), daily/weekly targets (creates pressure), leaderboards (competitive) |
| 8 | Extend existing ThemeColors for stats palette | Natural extension of Epic 17 theme infrastructure; each theme already defines colors | Separate stats theme system (fragmented config), hardcoded colors (ignores user theme choice) |
| 9 | Stats as "Trophy Room" door metaphor | Integrates stats into the door concept; makes viewing stats feel like a reward | Separate app section (disconnected), floating overlay (breaks Bubbletea patterns) |
| 10 | Three-phase rollout (polish → visualizations → integration) | Quick wins first to validate approach; each phase is independently valuable | Big-bang release (high risk), single story (too large) |

## Rejected Approaches

| Approach | Why Rejected |
|----------|--------------|
| termdash as visualization framework | Incompatible event loop; conflicts with Bubbletea architecture |
| Full charting library (go-echarts) | Web-oriented, renders HTML/SVG, not suitable for terminal |
| Achievement/badge system | SOUL.md: "no gamification, no guilt" |
| Productivity scores/grades | "Not a productivity report" — assigns judgment |
| Daily/weekly goal targets | Creates pressure; violates "progress over perfection" |
| Comparative analytics (less productive than...) | Negative framing violates encouraging tone |

## SOUL.md Alignment

All stories in this epic were evaluated against SOUL.md principles:

- **Green light:** Sparklines, panels, fun facts, heatmap, bar charts, animated counters, theme colors — pure visual delight and neutral observation
- **Yellow light (careful framing required):** Milestones — must be one-time celebration only, no "next milestone" messaging, no badges/levels/unlocks
- **Red light (excluded):** Leaderboards, productivity scores, overdue counts, achievement badges, "you should..." recommendations, completion targets

## Phase Structure

- **Phase 1 (Stories 40.1–40.3):** Visual polish of existing data — no new data infrastructure
- **Phase 2 (Stories 40.4–40.7):** New chart types and hidden metric surfacing — ntcharts evaluation before 40.5
- **Phase 3 (Stories 40.8–40.10):** Thematic integration with door metaphor and theme system
