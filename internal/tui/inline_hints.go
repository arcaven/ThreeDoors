package tui

import "github.com/charmbracelet/lipgloss"

// Inline hint ANSI colors per Story 39.9 / 39.10.
var (
	// hintStyleNormal renders key hints in dim gray (ANSI 245).
	hintStyleNormal = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))

	// hintStyleFade renders key hints in extra-dim gray (ANSI 240) for fade mode.
	hintStyleFade = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	// hintStyleBright renders key hints in bright white (ANSI 255) for selected doors.
	hintStyleBright = lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Bold(true)

	// hintStyleDim renders key hints in dark gray (ANSI 240) for unselected doors when a selection is active.
	hintStyleDim = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

// renderInlineHint returns a Lipgloss-styled "[key]" string when enabled,
// or an empty string when disabled. When fade is true, uses the extra-dim
// ANSI 240 style instead of the normal ANSI 245.
func renderInlineHint(key string, enabled bool, fade bool) string {
	if !enabled {
		return ""
	}
	style := hintStyleNormal
	if fade {
		style = hintStyleFade
	}
	return style.Render("[" + key + "]")
}

// renderDoorHint returns a selection-state-aware inline hint for a door.
// When selected, the hint brightens (ANSI 255 + bold). When another door is
// selected, unselected doors' hints dim (ANSI 240). When no door is selected,
// all hints use normal brightness. Fade mode overrides to extra-dim for all.
func renderDoorHint(key string, enabled bool, fade bool, isSelected bool, hasSelection bool) string {
	if !enabled {
		return ""
	}
	text := "[" + key + "]"
	if fade {
		return hintStyleFade.Render(text)
	}
	if hasSelection {
		if isSelected {
			return hintStyleBright.Render(text)
		}
		return hintStyleDim.Render(text)
	}
	return hintStyleNormal.Render(text)
}
