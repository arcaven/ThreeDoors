package reminders

import "context"

// CommandExecutor abstracts osascript execution for testability.
type CommandExecutor interface {
	Execute(ctx context.Context, script string) (string, error)
}
