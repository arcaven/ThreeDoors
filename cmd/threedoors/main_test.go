package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestUpdate(t *testing.T) {
	m := initialModel()

	// Test 'q' key press
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	newModel, cmd := m.Update(msg)

	if cmd == nil {
		t.Errorf("'q' key should trigger tea.Quit command, but got nil")
	}

	quitCmd, ok := cmd().(tea.QuitMsg)
	if !ok {
		t.Errorf("'q' key should return a tea.QuitMsg, but it didn't")
	}
	_ = quitCmd

	// Test 'ctrl+c' key press
	msgCtrlC := tea.KeyMsg{Type: tea.KeyCtrlC}
	newModel, cmd = newModel.Update(msgCtrlC)

	if cmd == nil {
		t.Errorf("'ctrl+c' key should trigger tea.Quit command, but got nil")
	}

	quitCmd, ok = cmd().(tea.QuitMsg)
	if !ok {
		t.Errorf("'ctrl+c' key should return a tea.QuitMsg, but it didn't")
	}
	_ = quitCmd

	// Test other key press
	msgOther := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
	_, cmd = newModel.Update(msgOther)

	if cmd != nil {
		t.Errorf("Another key should not trigger any command, but it did")
	}
}

func TestView(t *testing.T) {
	m := initialModel()
	expectedView := "ThreeDoors - Technical Demo"
	if view := m.View(); view != expectedView {
		t.Errorf("View() = %q, want %q", view, expectedView)
	}
}
