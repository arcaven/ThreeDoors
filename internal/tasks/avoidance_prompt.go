package tasks

import "time"

// AvoidanceAction represents the user's choice when prompted about an avoided task.
type AvoidanceAction string

const (
	ActionReconsider AvoidanceAction = "reconsider"
	ActionBreakDown  AvoidanceAction = "breakdown"
	ActionDefer      AvoidanceAction = "defer"
	ActionArchive    AvoidanceAction = "archive"
)

// AvoidancePromptState tracks which tasks have been prompted in the current session.
type AvoidancePromptState struct {
	promptedThisSession map[string]bool
}

// NewAvoidancePromptState creates a new session-scoped prompt state.
func NewAvoidancePromptState() *AvoidancePromptState {
	return &AvoidancePromptState{
		promptedThisSession: make(map[string]bool),
	}
}

// ShouldPrompt returns true if a task should trigger a re-evaluation prompt.
// Requires bypass count >= 10 and not already prompted this session.
func (s *AvoidancePromptState) ShouldPrompt(taskText string, bypassCount int) bool {
	if bypassCount < 10 {
		return false
	}
	return !s.promptedThisSession[taskText]
}

// MarkPrompted records that a task has been prompted this session.
func (s *AvoidancePromptState) MarkPrompted(taskText string) {
	s.promptedThisSession[taskText] = true
}

// DeferTask sets DeferredUntil to 7 days from now on the given task.
// Note: caller should also call UpdateStatus(StatusDeferred) which handles UpdatedAt.
func DeferTask(task *Task) {
	deferUntil := time.Now().UTC().AddDate(0, 0, 7)
	task.DeferredUntil = &deferUntil
}
