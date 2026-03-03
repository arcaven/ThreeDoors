package tasks

import (
	"testing"
	"time"
)

func TestAvoidancePromptState_ShouldPrompt_Below10(t *testing.T) {
	state := NewAvoidancePromptState()
	if state.ShouldPrompt("task1", 9) {
		t.Error("ShouldPrompt() returned true for bypassCount=9, want false")
	}
}

func TestAvoidancePromptState_ShouldPrompt_AtThreshold(t *testing.T) {
	state := NewAvoidancePromptState()
	if !state.ShouldPrompt("task1", 10) {
		t.Error("ShouldPrompt() returned false for bypassCount=10, want true")
	}
}

func TestAvoidancePromptState_ShouldPrompt_Above10(t *testing.T) {
	state := NewAvoidancePromptState()
	if !state.ShouldPrompt("task1", 15) {
		t.Error("ShouldPrompt() returned false for bypassCount=15, want true")
	}
}

func TestAvoidancePromptState_ShouldPrompt_AlreadyPrompted(t *testing.T) {
	state := NewAvoidancePromptState()
	state.MarkPrompted("task1")
	if state.ShouldPrompt("task1", 10) {
		t.Error("ShouldPrompt() returned true after MarkPrompted, want false")
	}
}

func TestAvoidancePromptState_MarkPrompted(t *testing.T) {
	state := NewAvoidancePromptState()
	if state.ShouldPrompt("task1", 10) != true {
		t.Fatal("precondition: ShouldPrompt should be true before marking")
	}
	state.MarkPrompted("task1")
	if state.ShouldPrompt("task1", 10) {
		t.Error("ShouldPrompt() should be false after MarkPrompted")
	}
	// Different task should not be affected
	if !state.ShouldPrompt("task2", 10) {
		t.Error("ShouldPrompt() for different task should still be true")
	}
}

func TestAvoidancePromptState_ShouldPrompt_Zero(t *testing.T) {
	state := NewAvoidancePromptState()
	if state.ShouldPrompt("task1", 0) {
		t.Error("ShouldPrompt() returned true for bypassCount=0")
	}
}

func TestDeferTask(t *testing.T) {
	task := newTestTask("id1", "Test task", StatusTodo, baseTime)
	if task.DeferredUntil != nil {
		t.Fatal("precondition: DeferredUntil should be nil")
	}

	before := time.Now().UTC()
	DeferTask(task)
	after := time.Now().UTC()

	if task.DeferredUntil == nil {
		t.Fatal("DeferTask() did not set DeferredUntil")
	}

	expectedEarliest := before.AddDate(0, 0, 7)
	expectedLatest := after.AddDate(0, 0, 7)

	if task.DeferredUntil.Before(expectedEarliest) || task.DeferredUntil.After(expectedLatest) {
		t.Errorf("DeferredUntil = %v, expected between %v and %v", task.DeferredUntil, expectedEarliest, expectedLatest)
	}
}
