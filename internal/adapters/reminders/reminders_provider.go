package reminders

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/arcaven/ThreeDoors/internal/core"
)

// RemindersProvider implements core.TaskProvider as a read-only adapter
// for Apple Reminders. Write operations return core.ErrReadOnly.
type RemindersProvider struct {
	executor CommandExecutor
	lists    []string
}

// NewRemindersProvider creates a RemindersProvider that reads from the given
// lists. An empty lists slice means all lists will be read.
func NewRemindersProvider(executor CommandExecutor, lists []string) *RemindersProvider {
	return &RemindersProvider{
		executor: executor,
		lists:    lists,
	}
}

// Name returns the provider identifier.
func (p *RemindersProvider) Name() string {
	return "reminders"
}

// LoadTasks reads incomplete reminders via JXA and maps them to core.Task values.
func (p *RemindersProvider) LoadTasks() ([]*core.Task, error) {
	ctx := context.Background()

	listNames, err := p.resolveListNames(ctx)
	if err != nil {
		return nil, fmt.Errorf("resolve reminder lists: %w", err)
	}

	var tasks []*core.Task
	for _, listName := range listNames {
		reminders, err := ReadReminders(ctx, p.executor, listName)
		if err != nil {
			return nil, fmt.Errorf("load reminders from %q: %w", listName, err)
		}
		for _, r := range reminders {
			tasks = append(tasks, mapReminderToTask(r, listName))
		}
	}

	return tasks, nil
}

// SaveTask returns ErrReadOnly.
func (p *RemindersProvider) SaveTask(_ *core.Task) error {
	return core.ErrReadOnly
}

// SaveTasks returns ErrReadOnly.
func (p *RemindersProvider) SaveTasks(_ []*core.Task) error {
	return core.ErrReadOnly
}

// DeleteTask returns ErrReadOnly.
func (p *RemindersProvider) DeleteTask(_ string) error {
	return core.ErrReadOnly
}

// MarkComplete returns ErrReadOnly.
func (p *RemindersProvider) MarkComplete(_ string) error {
	return core.ErrReadOnly
}

// Watch returns nil — the reminders provider does not support watching.
func (p *RemindersProvider) Watch() <-chan core.ChangeEvent {
	return nil
}

// HealthCheck attempts a lightweight read and reports TCC permission status.
func (p *RemindersProvider) HealthCheck() core.HealthCheckResult {
	start := time.Now().UTC()

	item := core.HealthCheckItem{Name: "Apple Reminders"}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := ReadLists(ctx, p.executor)
	if err != nil {
		item.Status = core.HealthFail
		errMsg := err.Error()
		if strings.Contains(errMsg, "not allowed") || strings.Contains(errMsg, "denied") || strings.Contains(errMsg, "1002") {
			item.Message = "Reminders access denied by macOS privacy settings"
			item.Suggestion = "Grant Reminders access in System Settings > Privacy & Security > Reminders"
		} else {
			item.Message = fmt.Sprintf("Cannot access Reminders: %s", errMsg)
			item.Suggestion = "Verify osascript is available and Reminders app is installed"
		}
	} else {
		item.Status = core.HealthOK
		item.Message = "Apple Reminders accessible"
	}

	return core.HealthCheckResult{
		Items:    []core.HealthCheckItem{item},
		Overall:  item.Status,
		Duration: time.Since(start),
	}
}

// resolveListNames returns the configured list names, or discovers all lists
// if none were configured.
func (p *RemindersProvider) resolveListNames(ctx context.Context) ([]string, error) {
	if len(p.lists) > 0 {
		return p.lists, nil
	}
	return ReadLists(ctx, p.executor)
}

// mapPriorityToEffort converts Apple Reminders priority (0-9) to core.TaskEffort.
// Reminders priority: 1 = highest, 5 = medium, 9 = lowest, 0 = none.
func mapPriorityToEffort(priority int) core.TaskEffort {
	switch {
	case priority >= 1 && priority <= 4:
		return core.EffortDeepWork
	case priority == 5:
		return core.EffortMedium
	case priority >= 6 && priority <= 9:
		return core.EffortQuickWin
	default:
		return ""
	}
}

// mapReminderToTask converts a ReminderJSON to a core.Task.
func mapReminderToTask(r ReminderJSON, listName string) *core.Task {
	status := core.StatusTodo
	if r.Completed {
		status = core.StatusComplete
	}

	createdAt := parseISOTime(r.CreationDate)
	updatedAt := parseISOTime(r.ModificationDate)

	var contextParts []string
	if r.DueDate != nil && *r.DueDate != "" {
		dueTime := parseISOTime(*r.DueDate)
		if !dueTime.IsZero() {
			contextParts = append(contextParts, fmt.Sprintf("due:%s", dueTime.Format("2006-01-02")))
		}
	}

	var notes []core.TaskNote
	if r.Body != "" {
		notes = append(notes, core.TaskNote{
			Timestamp: createdAt,
			Text:      r.Body,
		})
	}

	task := &core.Task{
		ID:             r.ID,
		Text:           r.Name,
		Context:        strings.Join(contextParts, " "),
		Status:         status,
		Effort:         mapPriorityToEffort(r.Priority),
		Notes:          notes,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
		SourceProvider: fmt.Sprintf("reminders:%s", listName),
		SourceRefs: []core.SourceRef{{
			Provider: fmt.Sprintf("reminders:%s", listName),
			NativeID: r.ID,
		}},
	}

	if r.Completed {
		task.CompletedAt = &updatedAt
	}

	return task
}

// parseISOTime parses an ISO 8601 timestamp, returning zero time on failure.
func parseISOTime(s string) time.Time {
	for _, layout := range []string{
		"2006-01-02T15:04:05.000Z",
		time.RFC3339,
		"2006-01-02T15:04:05Z",
	} {
		if t, err := time.Parse(layout, s); err == nil {
			return t.UTC()
		}
	}
	return time.Time{}
}
