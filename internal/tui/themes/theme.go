package themes

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

// ThemeColors holds the color palette for a door theme.
// All fields use lipgloss.TerminalColor to support adaptive color profiles
// (TrueColor, ANSI256, ANSI 16-color) for graceful degradation on
// constrained terminals.
type ThemeColors struct {
	Frame    lipgloss.TerminalColor
	Fill     lipgloss.TerminalColor
	Accent   lipgloss.TerminalColor
	Selected lipgloss.TerminalColor

	// Stats dashboard colors (Story 40.9). Zero values are safe —
	// InsightsView falls back to the independent palette when these are empty.
	StatsAccent        string // panel borders, hero number (#RRGGBB)
	StatsGradientStart string // sparkline low end (#RRGGBB)
	StatsGradientEnd   string // sparkline high end (#RRGGBB)
}

// MonthDay represents a calendar day within any year (month + day).
type MonthDay struct {
	Month int
	Day   int
}

// DoorTheme defines the visual frame for a door.
type DoorTheme struct {
	Name        string
	Description string
	Render      func(content string, width int, height int, selected bool, hint string) string
	Colors      ThemeColors
	MinWidth    int
	MinHeight   int

	// Seasonal metadata. Zero-value Season ("") indicates a non-seasonal theme.
	Season      string
	SeasonStart MonthDay
	SeasonEnd   MonthDay

	// HandleFrames defines the 4-character rotation for the handle turn
	// micro-animation (Story 48.4). Zero value means no animation.
	HandleFrames HandleFrames

	// handleOverride is shared with the Render closure via pointer. When
	// non-empty, the render closure uses this character instead of its
	// static default. Set via SetHandleChar before each Render call.
	handleOverride *string
}

// SetHandleChar overrides the handle character for the next Render call.
// Pass "" to revert to the theme's static default.
func (dt *DoorTheme) SetHandleChar(char string) {
	if dt.handleOverride != nil {
		*dt.handleOverride = char
	}
}

// HandleCharForAnimation computes the animated handle character based on
// the current spring emphasis and selection direction, then sets it for
// the next Render call. Returns the computed character.
func (dt *DoorTheme) HandleCharForAnimation(emphasis float64, deselecting bool) string {
	if dt.HandleFrames.Rest == "" {
		return ""
	}
	char := HandleCharForEmphasis(dt.HandleFrames, emphasis, deselecting)
	dt.SetHandleChar(char)
	return char
}

// DefaultThemeName is the theme used when no theme is specified.
const DefaultThemeName = "modern"

// renderHandleWithHint builds a handle row line, placing hint text to the left
// of the handle symbol when hint is non-empty. When hint is empty, renders the
// handle in its standard position. innerWidth is the total interior width between
// the vertical border characters. knobPad is the default left padding for the
// handle symbol. handleSym is the handle character (e.g. "●", "○", "◈──┤").
func renderHandleWithHint(innerWidth, knobPad int, handleSym, hint string) string {
	if hint == "" {
		rightPad := innerWidth - knobPad - 1
		if rightPad < 0 {
			rightPad = 0
		}
		return strings.Repeat(" ", knobPad) + handleSym + strings.Repeat(" ", rightPad)
	}
	hintWidth := ansi.StringWidth(hint)
	handleWidth := ansi.StringWidth(handleSym)
	// Layout: [padding] hint [space] handle [rightPad]
	leftPad := innerWidth - hintWidth - 1 - handleWidth - 1
	if leftPad < 1 {
		leftPad = 1
	}
	rightPad := innerWidth - leftPad - hintWidth - 1 - handleWidth
	if rightPad < 0 {
		rightPad = 0
	}
	return strings.Repeat(" ", leftPad) + hint + " " + handleSym + strings.Repeat(" ", rightPad)
}
