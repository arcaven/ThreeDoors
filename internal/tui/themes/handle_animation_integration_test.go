package themes

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// TestHandleAnimation_ThemeIntegration verifies that SetHandleChar changes
// the rendered handle character for each theme's door-mode output.
func TestHandleAnimation_ThemeIntegration(t *testing.T) {
	lipgloss.SetColorProfile(termenv.Ascii)
	t.Cleanup(func() { lipgloss.SetColorProfile(termenv.TrueColor) })

	tests := []struct {
		theme       *DoorTheme
		defaultChar string // static default in door mode
		animChar    string // character to override with
	}{
		{NewClassicTheme(), "●", "◐"},
		{NewModernTheme(), "○", "◑"},
		{NewShojiTheme(), "○", "●"},
		{NewWinterTheme(), "◆", "◇"},
		{NewSpringTheme(), "○", "◑"},
		{NewSummerTheme(), "■", "□"},
		{NewAutumnTheme(), "●", "○"},
	}

	for _, tt := range tests {
		t.Run(tt.theme.Name+"_default", func(t *testing.T) {
			tt.theme.SetHandleChar("")
			output := tt.theme.Render("Test", 40, 16, false, "")
			if !strings.Contains(output, tt.defaultChar) {
				t.Errorf("%s: expected default handle %q in output", tt.theme.Name, tt.defaultChar)
			}
		})

		t.Run(tt.theme.Name+"_animated", func(t *testing.T) {
			tt.theme.SetHandleChar(tt.animChar)
			output := tt.theme.Render("Test", 40, 16, false, "")
			if !strings.Contains(output, tt.animChar) {
				t.Errorf("%s: expected animated handle %q in output", tt.theme.Name, tt.animChar)
			}
			// Reset for next test
			tt.theme.SetHandleChar("")
		})
	}
}

// TestHandleAnimation_HandleCharForAnimation verifies the convenience method
// sets the override correctly based on emphasis and direction.
func TestHandleAnimation_HandleCharForAnimation(t *testing.T) {
	t.Parallel()

	theme := NewClassicTheme()

	tests := []struct {
		name        string
		emphasis    float64
		deselecting bool
		want        string
	}{
		{"rest", 0.0, false, "●"},
		{"turning", 0.4, false, "◐"},
		{"turned", 0.8, false, "○"},
		{"desel_turned", 0.8, true, "○"},
		{"desel_springback", 0.4, true, "◑"},
		{"desel_rest", 0.1, true, "●"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := theme.HandleCharForAnimation(tt.emphasis, tt.deselecting)
			if got != tt.want {
				t.Errorf("HandleCharForAnimation(%.1f, %v) = %q, want %q",
					tt.emphasis, tt.deselecting, got, tt.want)
			}
		})
	}
}

// TestHandleAnimation_SciFiStatic verifies that SciFi theme (no HandleFrames)
// always renders its static multi-character handle regardless of override attempts.
func TestHandleAnimation_SciFiStatic(t *testing.T) {
	lipgloss.SetColorProfile(termenv.Ascii)
	t.Cleanup(func() { lipgloss.SetColorProfile(termenv.TrueColor) })

	theme := NewSciFiTheme()

	// SciFi has empty HandleFrames — HandleCharForAnimation returns ""
	got := theme.HandleCharForAnimation(0.5, false)
	if got != "" {
		t.Errorf("SciFi HandleCharForAnimation should return empty, got %q", got)
	}

	// Static handle should always be present
	output := theme.Render("Test", 40, 16, false, "")
	if !strings.Contains(output, "◈") {
		t.Errorf("SciFi should always render its static ◈──┤ handle")
	}
}
