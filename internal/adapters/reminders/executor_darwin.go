//go:build darwin

package reminders

import (
	"context"
	"os/exec"
	"strings"
)

// OSAScriptExecutor implements CommandExecutor using osascript -l JavaScript.
type OSAScriptExecutor struct{}

// Execute runs a JXA script via osascript and returns the trimmed output.
func (e *OSAScriptExecutor) Execute(ctx context.Context, script string) (string, error) {
	cmd := exec.CommandContext(ctx, "osascript", "-l", "JavaScript", "-e", script)
	out, err := cmd.Output()
	return strings.TrimSpace(string(out)), err
}
