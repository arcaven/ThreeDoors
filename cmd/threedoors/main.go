package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/arcaven/ThreeDoors/internal/tasks"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// model represents the application's state.
type model struct {
	allTasks     []tasks.Task
	displayedDoors []tasks.Task
}

// initialModel initializes the application's model.
func initialModel() model {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	allTasks, err := tasks.LoadTasks()
	if err != nil {
		// In a real application, you might want to handle this error more gracefully,
		// perhaps by showing an error screen or logging it. For now, we'll panic.
		panic(fmt.Sprintf("Failed to load tasks: %v", err))
	}

	if len(allTasks) < 3 {
		// Ensure we have at least 3 tasks to display.
		// If not, we'll just use all available tasks and duplicate if necessary.
		// For this story, we assume there will be enough tasks.
		panic("Not enough tasks loaded to display three doors.")
	}

	// Randomly select 3 tasks
	displayedDoors := make([]tasks.Task, 3)
	perm := rand.Perm(len(allTasks))
	for i := 0; i < 3; i++ {
		displayedDoors[i] = allTasks[perm[i]]
	}

	return model{
		allTasks:     allTasks,
		displayedDoors: displayedDoors,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	s := strings.Builder{}
	s.WriteString("ThreeDoors - Technical Demo\n\n")

	doorStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 2). // Using Padding(vertical, horizontal)
		Width(20)

	for i, task := range m.displayedDoors {
		doorContent := fmt.Sprintf("Door %d:\n%s", i+1, task.Text)
		s.WriteString(doorStyle.Render(doorContent))
		s.WriteString("  ") // Space between doors
	}
	s.WriteString("\n\nPress 'q' or 'ctrl+c' to quit.\n")
	return s.String()
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}