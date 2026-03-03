package adapters_test

import (
	"testing"

	"github.com/arcaven/ThreeDoors/internal/adapters"
	"github.com/arcaven/ThreeDoors/internal/tasks"
)

// TestWALProviderContract runs the full contract test suite against
// WALProvider wrapping a TextFileProvider. When the inner provider
// succeeds, WAL should be transparent — all operations pass through.
func TestWALProviderContract(t *testing.T) {
	factory := func(t *testing.T) tasks.TaskProvider {
		t.Helper()

		dir := t.TempDir()
		tasks.SetHomeDir(dir)

		inner := tasks.NewTextFileProvider()
		walProvider := tasks.NewWALProvider(inner, dir)

		t.Cleanup(func() {
			tasks.SetHomeDir("")
		})

		return walProvider
	}

	adapters.RunContractTests(t, factory)
}

// TestWALProviderContract_PendingCountZero verifies that when the inner
// provider is healthy, no WAL entries accumulate.
func TestWALProviderContract_PendingCountZero(t *testing.T) {
	dir := t.TempDir()
	tasks.SetHomeDir(dir)
	t.Cleanup(func() {
		tasks.SetHomeDir("")
	})

	inner := tasks.NewTextFileProvider()
	walProvider := tasks.NewWALProvider(inner, dir)

	task := tasks.NewTask("WAL test task")
	if err := walProvider.SaveTask(task); err != nil {
		t.Fatalf("SaveTask() error: %v", err)
	}

	if walProvider.PendingCount() != 0 {
		t.Errorf("PendingCount() = %d, want 0 (inner provider succeeded)", walProvider.PendingCount())
	}
}
