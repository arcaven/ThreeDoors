# Source Tree

```
ThreeDoors/
├── cmd/
│   └── threedoors/
│       └── main.go                    # Application entry point, Bubbletea initialization
│
├── internal/                          # Private application code
│   ├── tui/                          # TUI Layer - Bubbletea components
│   │   ├── main_model.go            # Root Bubbletea model, view routing
│   │   ├── doors_view.go            # Three Doors display component
│   │   ├── task_detail_view.go      # Task detail and options component
│   │   ├── status_menu.go           # Status update menu subcomponent
│   │   ├── notes_input.go           # Notes text input subcomponent
│   │   ├── blocker_input.go         # Blocker input subcomponent
│   │   ├── styles.go                # Lipgloss style definitions
│   │   └── messages.go              # Bubbletea message types
│   │
│   └── tasks/                        # Domain Layer - Business logic
│       ├── task.go                  # Task model, methods, validation
│       ├── task_status.go           # TaskStatus enum, constants
│       ├── task_pool.go             # TaskPool collection manager
│       ├── door_selection.go        # DoorSelection model, algorithm
│       ├── door_selector.go         # Door selection logic
│       ├── status_manager.go        # Status transition validator
│       ├── file_manager.go          # YAML I/O, atomic writes
│       └── config.go                # Configuration model, defaults
│
├── docs/                             # Documentation
│   ├── prd.md                       # Product Requirements Document
│   ├── architecture.md              # This architecture document
│   └── stories/                     # Story breakdowns (from PRD)
│
├── .bmad-core/                       # BMAD methodology artifacts
│   ├── core-config.yaml
│   ├── agents/
│   ├── tasks/
│   ├── templates/
│   └── data/
│
├── .github/                          # GitHub configuration (Epic 2+)
│   └── workflows/                   # CI/CD pipelines (deferred)
│
├── bin/                              # Build output (gitignored)
│   └── threedoors                   # Compiled binary
│
├── go.mod                            # Go module definition
├── go.sum                            # Dependency checksums
├── Makefile                          # Build automation
├── .gitignore                        # Git ignore rules
└── README.md                         # Quick start guide

User Data Directory (created at runtime):
~/.threedoors/
├── tasks.yaml                        # Active tasks with metadata
└── completed.txt                     # Completed task log
```

**Key Organization Principles:**

1. **`cmd/` for entry points:** Single main.go bootstraps the application
2. **`internal/` for private code:** Cannot be imported by external projects
3. **`internal/tui/` for presentation:** All Bubbletea UI components
4. **`internal/tasks/` for domain:** Business logic, no UI dependencies
5. **Flat package structure:** No deep nesting (2 levels max)
6. **Clear separation:** TUI layer imports tasks, never vice versa

---
