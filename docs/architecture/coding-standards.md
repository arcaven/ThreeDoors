# Coding Standards

**⚠️ MANDATORY for AI Agents:** These standards directly control code generation behavior.

## Core Standards

**Languages & Runtimes:**
- Go 1.25.4+ strictly
- No external languages in codebase

**Style & Linting:**
- **Formatting:** `gofumpt` - run before every commit
- **Linting:** `golangci-lint run ./...` - must pass with zero warnings
- **Import ordering:** Standard library → external → internal (auto-formatted)

**Test Organization:**
- Test files: `*_test.go` alongside source files
- Table-driven tests preferred
- Test fixtures: `testdata/` directory

## Naming Conventions

| Element | Convention | Example |
|---------|-----------|---------|
| **Packages** | Lowercase, single word | `tui`, `tasks` |
| **Files** | Lowercase, snake_case | `task_pool.go`, `doors_view.go` |
| **Types (exported)** | PascalCase | `TaskPool`, `DoorSelection` |
| **Types (private)** | camelCase | `internalState` |
| **Functions (exported)** | PascalCase | `NewTaskPool`, `SelectDoors` |
| **Functions (private)** | camelCase | `validateTask`, `renderDoor` |
| **Constants** | PascalCase | `StatusTodo`, `MaxTasks` |

## Critical Rules

**MUST Follow:**

1. **Never use fmt.Println for user output in TUI code**
   - TUI output goes through Bubbletea View() methods only
   - Logging goes through log.Printf() to stderr

2. **All file writes must use atomic write pattern**
   - Write to `.tmp` file
   - Sync to disk
   - Atomic rename
   - Cleanup temp on error

3. **Always validate status transitions before applying**
   - Call StatusManager.ValidateTransition() first
   - Never allow direct Task.Status field assignment from UI

4. **Errors must be wrapped with context**
   - Use `%w` verb: `fmt.Errorf("operation failed: %w", err)`
   - Preserves error chain for errors.Is() and errors.As()

5. **No panics in user-facing code**
   - Bubbletea Update() and View() must never panic
   - Return error values, handle gracefully

6. **Task IDs are immutable**
   - UUID assigned at creation
   - Never modify Task.ID after creation

7. **Timestamps always stored in UTC**
   - Use `time.Now().UTC()` not `time.Now()`
   - Convert to local timezone only for display

8. **YAML field tags match schema exactly**
   - Use `yaml:"field_name"` tags
   - Use `omitempty` for nullable fields

## Atomic Write Pattern Checklist

**CRITICAL:** Every file write operation MUST follow this exact pattern to prevent data corruption:

```
✅ Step 1: Create temp path
   tempPath := targetPath + ".tmp"

✅ Step 2: Write to temp file
   if err := os.WriteFile(tempPath, data, 0644); err != nil {
       return fmt.Errorf("failed to write temp file: %w", err)
   }

✅ Step 3: Sync to disk (flush buffers)
   f, err := os.OpenFile(tempPath, os.O_RDWR, 0644)
   if err == nil {
       f.Sync()
       f.Close()
   }

✅ Step 4: Atomic rename
   if err := os.Rename(tempPath, targetPath); err != nil {
       os.Remove(tempPath)  // Cleanup on failure
       return fmt.Errorf("failed to commit changes: %w", err)
   }

✅ Step 5: Success - temp file now atomically replaces target
```

**Why This Matters:**
- Prevents partial writes (crash during write leaves original intact)
- Prevents corruption (temp file discarded if write fails)
- Atomic rename is OS-level operation (succeeds or fails completely)

**Reference Implementation:** See `FileManager.SaveTasks()` in Section 5 (Components)

---
