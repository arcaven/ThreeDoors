# Data Storage Schema

## tasks.yaml Schema

**Location:** `~/.threedoors/tasks.yaml`

**Format:** YAML with strict schema validation

**Root Structure:**
```yaml
tasks:
  - # Array of Task objects
```

**Task Object Schema:**

```yaml
id: string                 # UUID v4, required
text: string              # 1-500 chars, required
status: string            # Enum: todo|blocked|in-progress|in-review|complete, required
notes:                    # Array of TaskNote objects, can be empty
  - timestamp: datetime   # RFC3339 format, required
    text: string          # 1-1000 chars, required
blocker: string           # Empty or 1-500 chars, required when status=blocked
created_at: datetime      # RFC3339 format, required
updated_at: datetime      # RFC3339 format, required, >= created_at
completed_at: datetime    # RFC3339 format, nullable, only when status=complete
```

**Example:**
```yaml
tasks:
  - id: a1b2c3d4-e5f6-7890-abcd-ef1234567890
    text: Write architecture document for ThreeDoors
    status: in-progress
    notes:
      - timestamp: 2025-11-07T14:15:00Z
        text: Started with high-level overview
      - timestamp: 2025-11-07T14:45:00Z
        text: Completed data models section
    blocker: ""
    created_at: 2025-11-07T10:00:00Z
    updated_at: 2025-11-07T14:45:00Z
    completed_at: null
```

**Validation Rules:**
1. All timestamps in UTC (RFC3339 format)
2. Task IDs must be unique across all tasks
3. Empty blocker field when status != blocked
4. completedAt must be null unless status == complete
5. notes array preserves chronological order (newest last)

## completed.txt Schema

**Location:** `~/.threedoors/completed.txt`

**Format:** Plain text, append-only log

**Line Format:**
```
[YYYY-MM-DD HH:MM:SS] task_id | task_text
```

**Example:**
```
[2025-11-07 14:32:15] a1b2c3d4-e5f6-7890-abcd-ef1234567890 | Write architecture document for ThreeDoors
[2025-11-07 14:45:03] b2c3d4e5-f6a7-8901-bcde-f12345678901 | Implement Story 1.1 - Project Setup
```

---
