package themes

import "testing"

// Standard circular frames used by classic/autumn themes.
var standardFrames = HandleFrames{
	Rest:       "●",
	Turning:    "◐",
	Turned:     "○",
	SpringBack: "◑",
}

func TestHandleCharForEmphasis_ForwardSequence(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		emphasis float64
		want     string
	}{
		{"at rest", 0.0, "●"},
		{"just below turning threshold", 0.29, "●"},
		{"at turning threshold", 0.3, "◐"},
		{"mid turning", 0.45, "◐"},
		{"just below turned threshold", 0.59, "◐"},
		{"at turned threshold", 0.6, "○"},
		{"mid turned", 0.8, "○"},
		{"fully selected", 1.0, "○"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := HandleCharForEmphasis(standardFrames, tt.emphasis, false)
			if got != tt.want {
				t.Errorf("HandleCharForEmphasis(%.2f, forward) = %q, want %q", tt.emphasis, got, tt.want)
			}
		})
	}
}

func TestHandleCharForEmphasis_ReverseSequence(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		emphasis float64
		want     string
	}{
		{"fully turned", 1.0, "○"},
		{"still turned", 0.7, "○"},
		{"at turned threshold", 0.6, "○"},
		{"just below turned", 0.59, "◑"},
		{"mid springback", 0.45, "◑"},
		{"at springback threshold", 0.3, "◑"},
		{"just below springback", 0.29, "●"},
		{"near rest", 0.1, "●"},
		{"at rest", 0.0, "●"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := HandleCharForEmphasis(standardFrames, tt.emphasis, true)
			if got != tt.want {
				t.Errorf("HandleCharForEmphasis(%.2f, reverse) = %q, want %q", tt.emphasis, got, tt.want)
			}
		})
	}
}

func TestHandleCharForEmphasis_Deterministic(t *testing.T) {
	t.Parallel()

	// Same emphasis must always produce the same character.
	for _, emphasis := range []float64{0.0, 0.15, 0.3, 0.45, 0.6, 0.8, 1.0} {
		first := HandleCharForEmphasis(standardFrames, emphasis, false)
		for i := 0; i < 10; i++ {
			got := HandleCharForEmphasis(standardFrames, emphasis, false)
			if got != first {
				t.Errorf("emphasis %.2f: call %d returned %q, first was %q", emphasis, i, got, first)
			}
		}
	}
}

func TestHandleCharForEmphasis_ClampOvershoot(t *testing.T) {
	t.Parallel()

	// Spring physics can overshoot — values should be clamped.
	tests := []struct {
		name     string
		emphasis float64
		desel    bool
		want     string
	}{
		{"overshoot positive", 1.2, false, "○"},
		{"overshoot negative", -0.1, false, "●"},
		{"overshoot positive reverse", 1.3, true, "○"},
		{"overshoot negative reverse", -0.2, true, "●"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := HandleCharForEmphasis(standardFrames, tt.emphasis, tt.desel)
			if got != tt.want {
				t.Errorf("HandleCharForEmphasis(%.2f, desel=%v) = %q, want %q",
					tt.emphasis, tt.desel, got, tt.want)
			}
		})
	}
}

func TestHandleCharForEmphasis_EmptyFrames(t *testing.T) {
	t.Parallel()

	got := HandleCharForEmphasis(HandleFrames{}, 0.5, false)
	if got != "" {
		t.Errorf("empty frames should return empty string, got %q", got)
	}
}

func TestHandleCharForEmphasis_PerThemeFrames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		theme    string
		frames   HandleFrames
		emphasis float64
		desel    bool
		want     string
	}{
		// Winter diamond sequence
		{"winter forward rest", HandleFrames{"◆", "◇", "○", "◑"}, 0.0, false, "◆"},
		{"winter forward turning", HandleFrames{"◆", "◇", "○", "◑"}, 0.4, false, "◇"},
		{"winter forward turned", HandleFrames{"◆", "◇", "○", "◑"}, 0.8, false, "○"},
		{"winter reverse springback", HandleFrames{"◆", "◇", "○", "◑"}, 0.4, true, "◑"},

		// Summer square sequence
		{"summer forward rest", HandleFrames{"■", "◧", "□", "◧"}, 0.0, false, "■"},
		{"summer forward turning", HandleFrames{"■", "◧", "□", "◧"}, 0.4, false, "◧"},
		{"summer forward turned", HandleFrames{"■", "◧", "□", "◧"}, 0.8, false, "□"},
	}

	for _, tt := range tests {
		t.Run(tt.theme, func(t *testing.T) {
			t.Parallel()
			got := HandleCharForEmphasis(tt.frames, tt.emphasis, tt.desel)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
