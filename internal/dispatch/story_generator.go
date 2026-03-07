package dispatch

import "context"

// StoryResult holds the output of story generation for use during dispatch.
type StoryResult struct {
	// Branch is the git branch where story files were committed.
	Branch string
	// FilePaths lists the generated story file paths relative to the repo root.
	FilePaths []string
}

// StoryGenerator generates BMAD-style story files from a task description.
// Implementations bridge to AgentService or other story-generation backends.
type StoryGenerator interface {
	GenerateStories(ctx context.Context, taskDescription string) (*StoryResult, error)
}
