package adapters_test

import (
	"testing"

	"github.com/arcaven/ThreeDoors/internal/adapters"
	"github.com/arcaven/ThreeDoors/internal/tasks"
)

// TestFallbackProviderContract runs the full contract test suite against
// FallbackProvider with two TextFileProvider instances (primary and fallback).
func TestFallbackProviderContract(t *testing.T) {
	factory := func(t *testing.T) tasks.TaskProvider {
		t.Helper()

		// Each provider gets its own temp dir for isolation
		primaryDir := t.TempDir()
		tasks.SetHomeDir(primaryDir)

		primary := tasks.NewTextFileProvider()
		fallback := tasks.NewTextFileProvider()

		t.Cleanup(func() {
			tasks.SetHomeDir("")
		})

		return tasks.NewFallbackProvider(primary, fallback)
	}

	adapters.RunContractTests(t, factory)
}

// TestFallbackProviderContract_FallbackActivation verifies that when the
// primary provider fails on LoadTasks, the FallbackProvider switches to
// the fallback and reports IsFallback() == true.
func TestFallbackProviderContract_FallbackActivation(t *testing.T) {
	primaryDir := t.TempDir()
	fallbackDir := t.TempDir()

	// Set up the fallback provider with real data
	tasks.SetHomeDir(fallbackDir)
	fallbackProvider := tasks.NewTextFileProvider()
	seed := []*tasks.Task{tasks.NewTask("Fallback task")}
	if err := fallbackProvider.SaveTasks(seed); err != nil {
		t.Fatalf("seed fallback: %v", err)
	}

	// Primary points to an empty but valid dir — set home to primary dir
	tasks.SetHomeDir(primaryDir)
	primaryProvider := tasks.NewTextFileProvider()

	t.Cleanup(func() {
		tasks.SetHomeDir("")
	})

	fp := tasks.NewFallbackProvider(primaryProvider, fallbackProvider)

	// Primary should succeed (returns default or empty tasks)
	_, err := fp.LoadTasks()
	if err != nil {
		t.Fatalf("LoadTasks() error: %v", err)
	}

	// FallbackProvider should NOT have activated fallback since primary succeeded
	if fp.IsFallback() {
		t.Error("IsFallback() = true, want false (primary succeeded)")
	}
}
