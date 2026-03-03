package tui

import (
	"strings"
	"testing"

	"github.com/arcaven/ThreeDoors/internal/core"
)

// --- SourceBadgeLabel ---

func TestSourceBadgeLabel(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		provider string
		want     string
	}{
		{"textfile", "textfile", "TXT"},
		{"obsidian", "obsidian", "OBS"},
		{"applenotes", "applenotes", "NOTES"},
		{"unknown short", "jira", "JIRA"},
		{"unknown long", "todoist", "TODO"},
		{"empty", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := SourceBadgeLabel(tt.provider)
			if got != tt.want {
				t.Errorf("SourceBadgeLabel(%q) = %q, want %q", tt.provider, got, tt.want)
			}
		})
	}
}

// --- SourceBadge (rendered) ---

func TestSourceBadge_ContainsLabel(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		provider string
		contains string
	}{
		{"textfile", "textfile", "TXT"},
		{"obsidian", "obsidian", "OBS"},
		{"applenotes", "applenotes", "NOTES"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := SourceBadge(tt.provider)
			if !strings.Contains(got, tt.contains) {
				t.Errorf("SourceBadge(%q) = %q, want to contain %q", tt.provider, got, tt.contains)
			}
		})
	}
}

func TestSourceBadge_EmptyProviderReturnsEmpty(t *testing.T) {
	t.Parallel()
	got := SourceBadge("")
	if got != "" {
		t.Errorf("SourceBadge(\"\") = %q, want empty", got)
	}
}

// --- Integration: DoorsView shows source badge ---

func TestDoorsView_ShowsSourceBadge(t *testing.T) {
	t.Parallel()
	pool := core.NewTaskPool()
	task := core.NewTask("test task from obsidian")
	task.SourceProvider = "obsidian"
	pool.AddTask(task)

	// Add more tasks to fill doors
	for i := 0; i < 3; i++ {
		t2 := core.NewTask("filler task")
		t2.SourceProvider = "textfile"
		pool.AddTask(t2)
	}

	tracker := core.NewSessionTracker()
	dv := NewDoorsView(pool, tracker)
	view := dv.View()

	// The rendered view should contain source badge text
	if !strings.Contains(view, "OBS") && !strings.Contains(view, "TXT") {
		t.Error("expected door view to contain source badge (OBS or TXT)")
	}
}

// --- Integration: DetailView shows source provider ---

func TestDetailView_ShowsSourceProvider(t *testing.T) {
	t.Parallel()
	task := core.NewTask("test task from obsidian")
	task.SourceProvider = "obsidian"
	pool := core.NewTaskPool()
	pool.AddTask(task)
	tracker := core.NewSessionTracker()

	dv := NewDetailView(task, tracker, nil, pool)
	view := dv.View()

	if !strings.Contains(view, "obsidian") && !strings.Contains(view, "OBS") {
		t.Error("expected detail view to contain source provider info")
	}
}

// --- Integration: SearchView shows source badge ---

func TestSearchView_ShowsSourceBadge(t *testing.T) {
	t.Parallel()
	pool := core.NewTaskPool()
	task := core.NewTask("searchable task from obsidian")
	task.SourceProvider = "obsidian"
	pool.AddTask(task)

	tracker := core.NewSessionTracker()
	sv := NewSearchView(pool, tracker, nil, nil, nil)

	// Simulate search
	sv.results = []*core.Task{task}
	view := sv.View()

	if !strings.Contains(view, "OBS") {
		t.Error("expected search view to contain source badge OBS")
	}
}

// --- DuplicateIndicator ---

func TestDuplicateIndicator(t *testing.T) {
	t.Parallel()
	indicator := DuplicateIndicator()
	if !strings.Contains(indicator, "Possible duplicate") {
		t.Errorf("DuplicateIndicator() = %q, want to contain 'Possible duplicate'", indicator)
	}
}
