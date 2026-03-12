package themes

// HandleFrames defines the 4-character rotation sequence for a door handle's
// turn micro-animation. Characters cycle based on the spring emphasis value
// during selection (forward) and deselection (reverse).
type HandleFrames struct {
	Rest       string // emphasis 0.0 — at rest (default handle symbol)
	Turning    string // emphasis ~0.3 — beginning to turn
	Turned     string // emphasis ~0.6+ — fully turned (door ajar)
	SpringBack string // reverse ~0.5 — springing back to rest
}

// HandleCharForEmphasis maps a spring emphasis value to the appropriate handle
// character for the current animation state. When deselecting is true, uses the
// reverse sequence (Turned → SpringBack → Rest); otherwise uses the forward
// sequence (Rest → Turning → Turned).
//
// Returns the Rest character when frames are empty (no animation configured).
func HandleCharForEmphasis(frames HandleFrames, emphasis float64, deselecting bool) string {
	if frames.Rest == "" {
		return ""
	}

	// Clamp to [0, 1] — spring physics can overshoot
	if emphasis < 0 {
		emphasis = 0
	}
	if emphasis > 1 {
		emphasis = 1
	}

	if deselecting {
		// Reverse: emphasis decreasing from 1.0 → 0.0
		// 1.0–0.6: Turned (still open)
		// 0.6–0.3: SpringBack (snapping shut)
		// 0.3–0.0: Rest (closed)
		switch {
		case emphasis >= 0.6:
			return frames.Turned
		case emphasis >= 0.3:
			return frames.SpringBack
		default:
			return frames.Rest
		}
	}

	// Forward: emphasis increasing from 0.0 → 1.0
	// 0.0–0.3: Rest (closed)
	// 0.3–0.6: Turning (beginning to turn)
	// 0.6–1.0: Turned (fully turned)
	switch {
	case emphasis >= 0.6:
		return frames.Turned
	case emphasis >= 0.3:
		return frames.Turning
	default:
		return frames.Rest
	}
}
