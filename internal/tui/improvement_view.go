package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// ImprovementView displays the session improvement prompt on exit.
type ImprovementView struct {
	input string
	width int
}

// NewImprovementView creates a new improvement prompt view.
func NewImprovementView() *ImprovementView {
	return &ImprovementView{}
}

// SetWidth sets the terminal width.
func (iv *ImprovementView) SetWidth(w int) {
	iv.width = w
}

// Update handles key input for the improvement prompt.
func (iv *ImprovementView) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			text := strings.TrimSpace(iv.input)
			if text != "" {
				return func() tea.Msg { return ImprovementSubmittedMsg{Text: text} }
			}
			// Empty input on Enter — skip
			return func() tea.Msg { return ImprovementSkippedMsg{} }
		case "esc", "ctrl+c":
			return func() tea.Msg { return ImprovementSkippedMsg{} }
		case "backspace":
			if len(iv.input) > 0 {
				iv.input = iv.input[:len(iv.input)-1]
			}
		default:
			if len(msg.String()) == 1 && len(iv.input) < 500 {
				iv.input += msg.String()
			}
		}
	}
	return nil
}

// View renders the improvement prompt dialog.
func (iv *ImprovementView) View() string {
	s := strings.Builder{}
	s.WriteString(improvementHeaderStyle.Render("Session Reflection"))
	s.WriteString("\n\n")

	s.WriteString("What's one thing you could improve about\n")
	s.WriteString("this list/task/goal right now?\n\n")

	s.WriteString("> " + iv.input + "_\n\n")

	s.WriteString(helpStyle.Render("Enter to save | Esc to skip"))

	w := iv.width - 6
	if w < 40 {
		w = 40
	}
	return detailBorder.Width(w).Render(s.String())
}
