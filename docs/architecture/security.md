# Security

## Input Validation

**Validation Location:** At data model layer (Task constructors)

**Required Rules:**
1. Length limits enforced (500 chars for task, 1000 for notes)
2. No newlines or tabs in task text
3. Trim whitespace before validation
4. Reject empty strings

**Implementation:**
```go
func (t *Task) Validate() error {
    trimmed := strings.TrimSpace(t.Text)
    if len(trimmed) == 0 {
        return errors.New("task text cannot be empty")
    }
    if len(trimmed) > 500 {
        return errors.New("task text exceeds 500 characters")
    }
    if strings.ContainsAny(trimmed, "\n\t") {
        return errors.New("task text contains invalid characters")
    }
    return nil
}
```

## Authentication & Authorization

**Status:** Not Applicable (local-only application)

## Secrets Management

**Status:** Not Applicable (no secrets in Tech Demo)

**Future (Epic 2+):**
- macOS Keychain for credentials
- Never log credentials

## Data Protection

**Encryption at Rest:** Not Implemented

**Rationale:**
- Tasks are not highly sensitive
- User's macOS FileVault provides disk-level encryption
- Plain text YAML allows manual editing

**PII Handling:**
- Task text may contain PII (user decides)
- Data stays local
- No telemetry

**Logging Restrictions:**
- Never log task text content
- Log task IDs and operations only

**Example:**
```go
// ✅ GOOD: Log without sensitive data
log.Printf("INFO: Saved task %s with status %s\n", task.ID, task.Status)

// ❌ BAD: Log task content
log.Printf("Saved task: %s\n", task.Text) // May contain PII!
```

## Dependency Security

**Update Policy:** Update quarterly or when vulnerabilities reported

**Current Dependencies:**
- `github.com/charmbracelet/bubbletea` - Well-maintained
- `github.com/charmbracelet/lipgloss` - Same org
- `gopkg.in/yaml.v3` - Mature, stable

## Security Checklist

- ✅ No network connections
- ✅ File permissions: 0644 for data files
- ✅ No shell command execution
- ✅ Input validation for all user data
- ✅ Atomic writes prevent corruption
- ✅ No logging of sensitive data

---
