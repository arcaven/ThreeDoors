# Default Tooltips / Inline Keybinding Hints Mode — Party Mode Artifact

**Date:** 2026-03-09
**Participants:** Sally (UX Designer), John (PM), Winston (Architect)
**Facilitator:** Party Mode Orchestrator
**Related:** Epic 39 (Keybinding Display System), SOUL.md

---

## Problem Statement

New ThreeDoors users see three doors but don't know which keys to press. The existing help line at the bottom of the doors view (`a/left, w/up, d/right...`) is a wall of text that many users skip. The keybinding bar (39.2) and overlay (39.3) provide external reference, but nothing teaches users by showing key labels **directly on** the interactive elements.

## Adopted Approach

### Inline Hints as Epic 39 Stories (39.9–39.12)

Add inline keybinding hints rendered directly on/near interactive UI elements. Hints are ON by default for new users and auto-fade after N sessions.

**Key Design Decisions:**

1. **Doorknob Metaphor (Approach B — Frame Decoration)**
   - Key labels (`[a]`, `[w]`, `[d]`) rendered as part of the door frame/border, not inside the content area
   - Visually: the hint sits in the bottom-right of the door border, like a doorbell or doorknob label
   - Each theme's `Render()` function extended with a `hint string` parameter
   - When hints disabled, the border renders normally (no empty space)

   ```
   ┌──────────────────────┐
   │                      │
   │   Fix the login bug  │
   │   📋 S               │
   │                      │
   └─────────────── [a] ──┘
   ```

2. **Separate Config from Bar Toggle**
   - `show_inline_hints: true` — independent of `show_keybinding_bar`
   - No runtime toggle key — avoids collision with 39.4's `h` toggle
   - Re-enable/disable via `:hints on` / `:hints off` command
   - `:hints` command discoverable via overlay (39.3)

3. **Auto-Fade After N Sessions**
   - `inline_hints_session_count` incremented on each app launch
   - `inline_hints_fade_threshold` defaults to 5
   - Session N-1: hints render in dimmer style (graceful visual fade)
   - Session N: hints disappear, one-time flash message: "Key hints hidden — type :hints to bring them back"
   - User can re-enable anytime via `:hints on`

4. **Per-View Inline Hints (Non-Door Views)**
   - Detail view: `[esc] Back  [c] Complete  [b] Block` as subtle inline row
   - Search view: `[enter] search  [esc] cancel` near input field
   - Mood view: numbered labels next to emoji options
   - Add task view: `[enter] submit  [esc] cancel` near input
   - All hints sourced from keybinding registry (39.1)

5. **Self-Bootstrapping Discoverability**
   - Footer includes `hints: on` indicator when active
   - The hint system tells users how to dismiss itself
   - Overlay (39.3) lists `:hints` command

6. **Selected Door Hint Behavior**
   - When a door is selected, its hint brightens or changes to `[✓]`
   - Unselected door hints dim along with the rest of the door

### Config Model

```yaml
show_keybinding_bar: true        # 39.4 toggle (h key)
show_inline_hints: true          # tooltips, default on
inline_hints_session_count: 0    # auto-incremented
inline_hints_fade_threshold: 5   # auto-disable after N sessions
```

### SOUL.md Alignment

- "Hey, press me" energy — playful doorknobs, not patronizing tooltips
- Friction reduction — no need to read a help line or press `?`
- Physical metaphor — doorknob labels feel like part of the object
- Progressive disclosure — hints fade as user learns

## Rejected Approaches

### New Epic (Separate from Epic 39)
**Rationale for rejection:** Same keybinding registry data source (39.1), same config infrastructure, same toggle ecosystem. Creating a new epic would fragment the discoverability story unnecessarily and duplicate the data layer.

### Shared `h` Toggle (Tri-State Cycle)
**Rationale for rejection:** 39.4 explicitly defines `h` as bar-only toggle. Overloading it with inline hints creates a confusing tri-state cycle (bar+inline → bar only → clean → repeat). The bar is a reference tool (toggle on demand); inline hints are onboarding scaffolding (set and forget). Different lifecycle = different controls.

### Approach A — Content Injection
**Rationale for rejection:** Adding key labels inside the door's content area (mixed with task text, badges, status) breaks the door metaphor. The hint becomes "text inside a box" instead of "a label on the door frame." Approach B (frame decoration) preserves the physical object metaphor that SOUL.md demands.

### Per-Key Usage Tracking
**Rationale for rejection:** Tracking how many times each specific key has been pressed adds complex state management (per-key counters, cross-session persistence) for marginal UX benefit over simple session counting. Session count is a proxy that's 90% as effective with 10% of the complexity.

### Runtime Toggle Key for Inline Hints
**Rationale for rejection:** No good single-key binding available. `h` is taken by bar toggle. Adding a new key increases the keybinding surface that new users need to learn — defeating the purpose of discoverability. `:hints` command is sufficient for the rare case where an experienced user wants to re-enable hints.

## Deferred Ideas

- **Per-key progressive fade:** Track which specific keys the user has pressed and fade only those hints. Revisit if session-based fade feels too blunt in practice.
- **Animated hint reveal on first launch:** Hints could animate in (fade up) on the very first session for a "welcome" feel. Nice polish but unnecessary for v1.
- **Theme-specific hint styling:** Each door theme could style its hints differently (Gothic: ornate bracket, Modern: clean pill shape). Defer to a polish pass after base implementation.

## Story Breakdown

### Story 39.9: Inline Hint Rendering Infrastructure (P1, S)
- `renderInlineHint(key string, enabled bool) string` helper
- Theme `Render()` signature extended: `Render(content, width, height, selected, hint) string`
- Config model: `show_inline_hints`, `inline_hints_session_count`, `inline_hints_fade_threshold`
- `:hints on/off` command registration
- Unit tests

### Story 39.10: Door View Inline Hints (P1, M)
- `[a]`, `[w]`, `[d]` on door frames via Approach B
- Selected door hint changes (brighten or `[✓]`)
- `[s]` re-roll and `[n]` add hints in footer area
- Golden file tests for doors with/without hints
- All theme Render() implementations updated

### Story 39.11: Non-Door View Inline Hints (P1, M)
- Detail view, search view, mood view, add task view hints
- Hints sourced from keybinding registry (39.1)
- Hint rendering conditioned on `show_inline_hints` config
- Golden file tests per view

### Story 39.12: Auto-Fade After N Sessions (P2, S)
- Session counter increment on launch
- Graceful dim at session N-1
- Auto-disable at session N with flash message
- Config persistence of fade state
- Re-enable via `:hints on`

## Dependencies

- 39.1 (Keybinding Registry) — data source for all hints ✅ DONE
- 39.4 (Toggle + Config) — config persistence pattern, must be compatible
- Epic 17 (Theme System) — theme Render() signature change affects all themes
