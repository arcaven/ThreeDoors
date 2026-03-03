package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Known provider badge labels.
var providerLabels = map[string]string{
	"textfile":   "TXT",
	"obsidian":   "OBS",
	"applenotes": "NOTES",
}

// Provider badge colors.
var providerColors = map[string]lipgloss.Color{
	"textfile":   lipgloss.Color("243"), // gray
	"obsidian":   lipgloss.Color("141"), // purple
	"applenotes": lipgloss.Color("220"), // yellow
}

// SourceBadgeLabel returns the short label for a provider name.
// Known providers get predefined labels; unknown providers get an
// uppercase abbreviation (up to 4 characters).
func SourceBadgeLabel(provider string) string {
	if provider == "" {
		return ""
	}
	if label, ok := providerLabels[provider]; ok {
		return label
	}
	upper := strings.ToUpper(provider)
	if len(upper) > 4 {
		return upper[:4]
	}
	return upper
}

// SourceBadge returns a styled badge string for display in TUI views.
func SourceBadge(provider string) string {
	label := SourceBadgeLabel(provider)
	if label == "" {
		return ""
	}

	color, ok := providerColors[provider]
	if !ok {
		color = lipgloss.Color("243") // gray fallback
	}

	style := lipgloss.NewStyle().Foreground(color)
	return style.Render("[" + label + "]")
}

// DuplicateIndicator returns a styled indicator for potential duplicates.
func DuplicateIndicator() string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("208")).Faint(true)
	return style.Render("Possible duplicate")
}
