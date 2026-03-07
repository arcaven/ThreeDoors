package dispatch

import (
	"context"
	"fmt"
	"log"
	"strings"
)

// DispatchEngine orchestrates the dispatch pipeline: optional story generation
// followed by worker creation. It coordinates between StoryGenerator and Dispatcher.
type DispatchEngine struct {
	dispatcher Dispatcher
	stories    StoryGenerator // nil when no LLM backend is configured
	config     DispatchConfig
}

// NewDispatchEngine creates a DispatchEngine with the given dispatcher, optional
// story generator, and config. The stories parameter may be nil; if require_story
// is true and stories is nil, dispatch proceeds with a warning.
func NewDispatchEngine(dispatcher Dispatcher, stories StoryGenerator, config DispatchConfig) *DispatchEngine {
	return &DispatchEngine{
		dispatcher: dispatcher,
		stories:    stories,
		config:     config,
	}
}

// DispatchResult holds the outcome of a dispatch operation.
type DispatchResult struct {
	WorkerName  string
	StoryResult *StoryResult
	StoryErr    string // non-empty if story generation failed (non-fatal)
}

// Dispatch executes the full dispatch pipeline for a queue item:
// 1. Optionally generate story files (if require_story is true)
// 2. Build the task description (with story references if available)
// 3. Create a worker via the Dispatcher
func (e *DispatchEngine) Dispatch(ctx context.Context, item QueueItem) (*DispatchResult, error) {
	result := &DispatchResult{}

	// Phase 1: Optional story generation.
	if e.config.RequireStory {
		storyResult, err := e.generateStories(ctx, item)
		if err != nil {
			// Story generation failure is non-fatal (AC5).
			log.Printf("story generation failed for item %s: %v", item.ID, err)
			result.StoryErr = fmt.Sprintf("story generation failed: %v", err)
		} else {
			result.StoryResult = storyResult
		}
	}

	// Phase 2: Build task description with optional story references.
	taskDesc := BuildTaskDescriptionWithStories(item, result.StoryResult)

	// Phase 3: Create worker.
	workerName, err := e.dispatcher.CreateWorker(ctx, taskDesc)
	if err != nil {
		return nil, fmt.Errorf("dispatch item %s: %w", item.ID, err)
	}

	result.WorkerName = workerName
	return result, nil
}

// generateStories attempts to generate story files for the given queue item.
func (e *DispatchEngine) generateStories(ctx context.Context, item QueueItem) (*StoryResult, error) {
	if e.stories == nil {
		return nil, fmt.Errorf("story generator not configured (no LLM backend)")
	}

	taskDesc := item.TaskText
	if item.Context != "" {
		taskDesc = taskDesc + "\n\nContext: " + item.Context
	}

	return e.stories.GenerateStories(ctx, taskDesc)
}

// BuildTaskDescriptionWithStories constructs a worker prompt from a QueueItem,
// optionally including references to generated story files.
func BuildTaskDescriptionWithStories(item QueueItem, stories *StoryResult) string {
	if stories == nil {
		return BuildTaskDescription(item)
	}

	var b strings.Builder
	fmt.Fprintf(&b, "Task: %s\n", item.TaskText)

	if item.Context != "" {
		fmt.Fprintf(&b, "\nContext: %s\n", item.Context)
	}

	// Add story file references.
	fmt.Fprintf(&b, "\nGenerated Story Files (branch: %s):\n", stories.Branch)
	for _, fp := range stories.FilePaths {
		fmt.Fprintf(&b, "- %s\n", fp)
	}
	fmt.Fprintf(&b, "\nPlease implement the stories listed above. ")
	fmt.Fprintf(&b, "Check out branch %q and read each story file for detailed acceptance criteria and tasks.\n", stories.Branch)

	if len(item.AcceptanceCriteria) > 0 {
		fmt.Fprintf(&b, "\nAcceptance Criteria:\n")
		for _, ac := range item.AcceptanceCriteria {
			fmt.Fprintf(&b, "- %s\n", ac)
		}
	}

	if item.Scope != "" {
		fmt.Fprintf(&b, "\nScope: %s\n", item.Scope)
	}

	fmt.Fprintf(&b, "%s", taskSuffix)
	return b.String()
}
