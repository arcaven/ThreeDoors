package core

import (
	"testing"
)

// --- Levenshtein Distance ---

func TestLevenshteinDistance(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		a, b string
		want int
	}{
		{"identical", "hello", "hello", 0},
		{"empty both", "", "", 0},
		{"empty a", "", "abc", 3},
		{"empty b", "abc", "", 3},
		{"one edit", "kitten", "sitten", 1},
		{"classic", "kitten", "sitting", 3},
		{"completely different", "abc", "xyz", 3},
		{"case sensitive", "Hello", "hello", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := LevenshteinDistance(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("LevenshteinDistance(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

// --- TextSimilarity ---

func TestTextSimilarity(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		a, b    string
		wantMin float64
		wantMax float64
	}{
		{"identical", "buy groceries", "buy groceries", 1.0, 1.0},
		{"empty both", "", "", 1.0, 1.0},
		{"completely different", "abc", "xyz", 0.0, 0.01},
		{"similar", "buy groceries", "buy grocery", 0.7, 1.0},
		{"case insensitive", "Buy Groceries", "buy groceries", 1.0, 1.0},
		{"extra whitespace", "  buy  groceries  ", "buy groceries", 1.0, 1.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := TextSimilarity(tt.a, tt.b)
			if got < tt.wantMin || got > tt.wantMax {
				t.Errorf("TextSimilarity(%q, %q) = %f, want [%f, %f]", tt.a, tt.b, got, tt.wantMin, tt.wantMax)
			}
		})
	}
}

// --- DuplicatePair ---

func TestDuplicatePair_Fields(t *testing.T) {
	t.Parallel()
	pair := DuplicatePair{
		TaskA:      NewTask("buy groceries"),
		TaskB:      NewTask("buy grocery"),
		Similarity: 0.85,
	}
	if pair.TaskA == nil || pair.TaskB == nil {
		t.Fatal("expected non-nil tasks")
	}
	if pair.Similarity < 0.8 || pair.Similarity > 1.0 {
		t.Errorf("unexpected similarity: %f", pair.Similarity)
	}
}

// --- DetectDuplicates ---

func TestDetectDuplicates_FindsCrossProviderDuplicates(t *testing.T) {
	t.Parallel()
	taskA := NewTask("buy groceries from the store")
	taskA.SourceProvider = "textfile"
	taskB := NewTask("buy groceries from the store")
	taskB.SourceProvider = "obsidian"
	taskC := NewTask("completely different task")
	taskC.SourceProvider = "textfile"

	tasks := []*Task{taskA, taskB, taskC}
	pairs := DetectDuplicates(tasks, 0.8)

	if len(pairs) != 1 {
		t.Fatalf("expected 1 duplicate pair, got %d", len(pairs))
	}
	if pairs[0].Similarity < 0.8 {
		t.Errorf("expected similarity >= 0.8, got %f", pairs[0].Similarity)
	}
}

func TestDetectDuplicates_IgnoresSameProvider(t *testing.T) {
	t.Parallel()
	taskA := NewTask("buy groceries")
	taskA.SourceProvider = "textfile"
	taskB := NewTask("buy groceries")
	taskB.SourceProvider = "textfile"

	tasks := []*Task{taskA, taskB}
	pairs := DetectDuplicates(tasks, 0.8)

	if len(pairs) != 0 {
		t.Errorf("expected 0 pairs for same-provider tasks, got %d", len(pairs))
	}
}

func TestDetectDuplicates_BelowThreshold(t *testing.T) {
	t.Parallel()
	taskA := NewTask("buy groceries from the store")
	taskA.SourceProvider = "textfile"
	taskB := NewTask("send email to boss about meeting")
	taskB.SourceProvider = "obsidian"

	tasks := []*Task{taskA, taskB}
	pairs := DetectDuplicates(tasks, 0.8)

	if len(pairs) != 0 {
		t.Errorf("expected 0 pairs below threshold, got %d", len(pairs))
	}
}

func TestDetectDuplicates_DeterministicOrdering(t *testing.T) {
	t.Parallel()
	tasks := make([]*Task, 0, 6)
	// Create 3 cross-provider pairs at different similarities
	for _, text := range []string{"task alpha one", "task alpha two", "task beta one"} {
		ta := NewTask(text)
		ta.SourceProvider = "textfile"
		tb := NewTask(text + " extra")
		tb.SourceProvider = "obsidian"
		tasks = append(tasks, ta, tb)
	}

	pairs1 := DetectDuplicates(tasks, 0.5)
	pairs2 := DetectDuplicates(tasks, 0.5)

	if len(pairs1) != len(pairs2) {
		t.Fatalf("pair count differs: %d vs %d", len(pairs1), len(pairs2))
	}
	for i := range pairs1 {
		if pairs1[i].TaskA.ID != pairs2[i].TaskA.ID || pairs1[i].TaskB.ID != pairs2[i].TaskB.ID {
			t.Errorf("pair %d differs between runs", i)
		}
		if pairs1[i].Similarity != pairs2[i].Similarity {
			t.Errorf("pair %d similarity differs: %f vs %f", i, pairs1[i].Similarity, pairs2[i].Similarity)
		}
	}
}

func TestDetectDuplicates_SortedByScoreDescending(t *testing.T) {
	t.Parallel()
	// Create pairs with different similarity levels
	taskA := NewTask("buy groceries from the store")
	taskA.SourceProvider = "textfile"
	taskB := NewTask("buy groceries from the store") // identical
	taskB.SourceProvider = "obsidian"

	taskC := NewTask("walk the dog today")
	taskC.SourceProvider = "textfile"
	taskD := NewTask("walk the dog today please") // slightly different
	taskD.SourceProvider = "obsidian"

	tasks := []*Task{taskA, taskB, taskC, taskD}
	pairs := DetectDuplicates(tasks, 0.5)

	for i := 1; i < len(pairs); i++ {
		if pairs[i].Similarity > pairs[i-1].Similarity {
			t.Errorf("pairs not sorted by score descending: pair[%d]=%f > pair[%d]=%f",
				i, pairs[i].Similarity, i-1, pairs[i-1].Similarity)
		}
	}
}

func TestDetectDuplicates_EmptyTasks(t *testing.T) {
	t.Parallel()
	pairs := DetectDuplicates(nil, 0.8)
	if len(pairs) != 0 {
		t.Errorf("expected 0 pairs for nil tasks, got %d", len(pairs))
	}
}

func TestDetectDuplicates_SingleTask(t *testing.T) {
	t.Parallel()
	task := NewTask("only task")
	task.SourceProvider = "textfile"
	pairs := DetectDuplicates([]*Task{task}, 0.8)
	if len(pairs) != 0 {
		t.Errorf("expected 0 pairs for single task, got %d", len(pairs))
	}
}

// --- DedupDecision ---

func TestDedupDecisionConstants(t *testing.T) {
	t.Parallel()
	if DecisionDuplicate != "duplicate" {
		t.Errorf("expected %q, got %q", "duplicate", DecisionDuplicate)
	}
	if DecisionDistinct != "distinct" {
		t.Errorf("expected %q, got %q", "distinct", DecisionDistinct)
	}
}
