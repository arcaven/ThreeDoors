package adapters_test

import (
	"context"
	"testing"

	"github.com/arcaven/ThreeDoors/internal/adapters"
	"github.com/arcaven/ThreeDoors/internal/adapters/reminders"
	"github.com/arcaven/ThreeDoors/internal/core"
)

// stubExecutor returns canned JSON for ReadReminders calls.
type stubExecutor struct{}

func (s *stubExecutor) Execute(_ context.Context, _ string) (string, error) {
	return `[]`, nil
}

func TestRemindersProviderContract(t *testing.T) {
	t.Parallel()

	factory := func(t *testing.T) core.TaskProvider {
		t.Helper()
		return reminders.NewRemindersProvider(&stubExecutor{}, []string{"Test"})
	}

	adapters.RunContractTests(t, factory)
}
