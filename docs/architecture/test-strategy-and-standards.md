# Test Strategy and Standards

## Testing Philosophy

**Approach:** Pragmatic Testing

- Focus on domain logic (tasks package)
- Minimal TUI testing (manual testing preferred)
- Table-driven tests for multiple scenarios
- No mocking frameworks - use interfaces and simple stubs

**Coverage Goals:**
- `internal/tasks/*`: 70%+ coverage
- `internal/tui/*`: 20%+ coverage
- Overall: 50%+ coverage

**Test Pyramid:**
- **70% Unit tests:** Fast, isolated, single functions
- **20% Integration tests:** Component interactions
- **10% Manual testing:** End-to-end TUI workflows

## Test Types

### Unit Tests

**Framework:** Go's built-in `testing` package

**File Convention:** `task_test.go` alongside `task.go`

**Coverage Requirement:** 70%+ for domain logic

**Example:**
```go
func TestTask_UpdateStatus(t *testing.T) {
    tests := []struct {
        name          string
        currentStatus TaskStatus
        newStatus     TaskStatus
        wantErr       bool
    }{
        {
            name:          "todo to in-progress valid",
            currentStatus: StatusTodo,
            newStatus:     StatusInProgress,
            wantErr:       false,
        },
        // ... more test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            task := &Task{Status: tt.currentStatus}
            err := task.UpdateStatus(tt.newStatus, "")
            if (err != nil) != tt.wantErr {
                t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Integration Tests

**Scope:** Multi-component interactions

**Infrastructure:**
- **File I/O:** Use `t.TempDir()` for test files
- **YAML parsing:** Test with real YAML strings

**Example:**
```go
func TestFileManager_SaveAndLoad(t *testing.T) {
    tempDir := t.TempDir()
    config := &Config{
        TasksPath: filepath.Join(tempDir, "tasks.yaml"),
    }
    fm := NewFileManager(config)

    // Save tasks
    originalPool := NewTaskPool(10)
    task := NewTask("Test task")
    originalPool.AddTask(task)
    fm.SaveTasks(originalPool)

    // Load tasks
    loadedPool, err := fm.LoadTasks()
    if err != nil {
        t.Fatalf("LoadTasks() error = %v", err)
    }

    // Verify
    if loadedPool.Count() != 1 {
        t.Errorf("count = %d, want 1", loadedPool.Count())
    }
}
```

### Manual Testing (TUI)

**Test Scenarios:**
1. First run - creates sample tasks
2. Door selection and navigation
3. Status updates and validation
4. Notes and blocker input
5. Task completion and removal
6. Edge cases (0-2 tasks, all completed)

## Test Data Management

**Strategy:** Inline test data and temp files

**Cleanup:** Automatic via `t.TempDir()` and `t.Cleanup()`

---
