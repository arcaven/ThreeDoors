package core

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// doctorMockProvider implements TaskProvider for testing doctor provider checks.
type doctorMockProvider struct {
	name      string
	tasks     []*Task
	loadErr   error
	healthRes HealthCheckResult
}

func (m *doctorMockProvider) Name() string                   { return m.name }
func (m *doctorMockProvider) LoadTasks() ([]*Task, error)    { return m.tasks, m.loadErr }
func (m *doctorMockProvider) SaveTask(_ *Task) error         { return nil }
func (m *doctorMockProvider) SaveTasks(_ []*Task) error      { return nil }
func (m *doctorMockProvider) DeleteTask(_ string) error      { return nil }
func (m *doctorMockProvider) MarkComplete(_ string) error    { return nil }
func (m *doctorMockProvider) Watch() <-chan ChangeEvent      { return nil }
func (m *doctorMockProvider) HealthCheck() HealthCheckResult { return m.healthRes }

func newTestRegistry(providers map[string]*doctorMockProvider) *Registry {
	reg := NewRegistry()
	for name, mp := range providers {
		p := mp // capture
		if err := reg.Register(name, func(_ *ProviderConfig) (TaskProvider, error) {
			return p, nil
		}); err != nil {
			panic(fmt.Sprintf("failed to register %q: %v", name, err))
		}
	}
	return reg
}

func writeConfig(t *testing.T, dir, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, "config.yaml"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestCheckProviders_NoProviderConfigured(t *testing.T) {
	t.Parallel()

	// resolveProviderEntries returns nil when both Providers and Provider are empty.
	// LoadProviderConfig defaults Provider to "textfile", so this edge case can
	// only arise if config is manually constructed. We test it directly.
	dc := &DoctorChecker{configDir: t.TempDir(), registry: NewRegistry()}
	cfg := &ProviderConfig{SchemaVersion: CurrentSchemaVersion, Provider: ""}
	entries := dc.resolveProviderEntries(cfg)
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}

	// Now test the full checkProviders path with a config that has only
	// schema_version (the loader will default provider to textfile, but
	// the registry has no textfile registered, so it should fail as unregistered).
	tmpDir := t.TempDir()
	writeConfig(t, tmpDir, "schema_version: 2\n")

	dc2 := &DoctorChecker{configDir: tmpDir, registry: NewRegistry()}
	result := dc2.checkProviders()

	if len(result.Checks) != 1 {
		t.Fatalf("expected 1 check, got %d", len(result.Checks))
	}
	if result.Checks[0].Status != CheckFail {
		t.Errorf("status = %v, want %v", result.Checks[0].Status, CheckFail)
	}
}

func TestCheckProviders_SingleProviderOK(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	writeConfig(t, tmpDir, "schema_version: 2\nprovider: textfile\n")

	reg := newTestRegistry(map[string]*doctorMockProvider{
		"textfile": {
			name:  "textfile",
			tasks: []*Task{{ID: "1", Text: "test"}},
		},
	})

	dc := &DoctorChecker{configDir: tmpDir, registry: reg}
	result := dc.checkProviders()

	if len(result.Checks) != 1 {
		t.Fatalf("expected 1 check, got %d", len(result.Checks))
	}
	if result.Checks[0].Status != CheckOK {
		t.Errorf("status = %v, want %v (message: %s)", result.Checks[0].Status, CheckOK, result.Checks[0].Message)
	}
}

func TestCheckProviders_SingleProviderLoadFails(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	writeConfig(t, tmpDir, "schema_version: 2\nprovider: textfile\n")

	reg := newTestRegistry(map[string]*doctorMockProvider{
		"textfile": {
			name:    "textfile",
			loadErr: errors.New("file corrupt"),
		},
	})

	dc := &DoctorChecker{configDir: tmpDir, registry: reg}
	result := dc.checkProviders()

	if len(result.Checks) != 1 {
		t.Fatalf("expected 1 check, got %d", len(result.Checks))
	}
	if result.Checks[0].Status != CheckFail {
		t.Errorf("status = %v, want %v", result.Checks[0].Status, CheckFail)
	}
}

func TestCheckProviders_MultipleProviders_OneFailsWarnCategory(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	writeConfig(t, tmpDir, `schema_version: 2
providers:
  - name: textfile
  - name: obsidian
    settings:
      vault_path: /nonexistent/vault/path
`)

	reg := newTestRegistry(map[string]*doctorMockProvider{
		"textfile": {
			name:  "textfile",
			tasks: []*Task{{ID: "1", Text: "test"}},
		},
		"obsidian": {
			name: "obsidian",
		},
	})

	dc := &DoctorChecker{configDir: tmpDir, registry: reg}
	result := dc.checkProviders()
	result.Status = worstCheckStatus(result.Checks)

	// Should have checks for both providers
	if len(result.Checks) < 2 {
		t.Fatalf("expected at least 2 checks, got %d", len(result.Checks))
	}

	// Category status should be WARN or FAIL (worst of all checks)
	hasOK := false
	hasFail := false
	for _, check := range result.Checks {
		if check.Status == CheckOK {
			hasOK = true
		}
		if check.Status == CheckFail {
			hasFail = true
		}
	}
	if !hasOK {
		t.Error("expected at least one OK check (textfile)")
	}
	if !hasFail {
		t.Error("expected at least one FAIL check (obsidian with bad vault path)")
	}

	// When one fails and one succeeds, worst status should be FAIL
	// The category aggregation (via worstCheckStatus) in DoctorChecker.Run()
	// will set this to WARN or FAIL depending on the worst check
	if result.Status != CheckFail {
		t.Errorf("category status = %v, want %v", result.Status, CheckFail)
	}
}

func TestCheckProviders_ObsidianVaultNotFound(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	badPath := "/nonexistent/vault/path/for/test"
	writeConfig(t, tmpDir, fmt.Sprintf(`schema_version: 2
providers:
  - name: obsidian
    settings:
      vault_path: %s
`, badPath))

	reg := newTestRegistry(map[string]*doctorMockProvider{
		"obsidian": {name: "obsidian"},
	})

	dc := &DoctorChecker{configDir: tmpDir, registry: reg}
	result := dc.checkProviders()

	if len(result.Checks) != 1 {
		t.Fatalf("expected 1 check, got %d", len(result.Checks))
	}
	check := result.Checks[0]
	if check.Status != CheckFail {
		t.Errorf("status = %v, want %v", check.Status, CheckFail)
	}
	wantMsg := fmt.Sprintf("Obsidian vault path not found: %s", badPath)
	if check.Message != wantMsg {
		t.Errorf("message = %q, want %q", check.Message, wantMsg)
	}
}

func TestCheckProviders_ObsidianVaultExists(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	vaultDir := filepath.Join(tmpDir, "my-vault")
	if err := os.MkdirAll(vaultDir, 0o755); err != nil {
		t.Fatal(err)
	}

	writeConfig(t, tmpDir, fmt.Sprintf(`schema_version: 2
providers:
  - name: obsidian
    settings:
      vault_path: %s
`, vaultDir))

	reg := newTestRegistry(map[string]*doctorMockProvider{
		"obsidian": {
			name:  "obsidian",
			tasks: []*Task{},
		},
	})

	dc := &DoctorChecker{configDir: tmpDir, registry: reg}
	result := dc.checkProviders()

	if len(result.Checks) != 1 {
		t.Fatalf("expected 1 check, got %d", len(result.Checks))
	}
	if result.Checks[0].Status != CheckOK {
		t.Errorf("status = %v, want %v (message: %s)", result.Checks[0].Status, CheckOK, result.Checks[0].Message)
	}
}

func TestCheckProviders_UnregisteredProvider(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	writeConfig(t, tmpDir, "schema_version: 2\nprovider: fakeprovider\n")

	dc := &DoctorChecker{configDir: tmpDir, registry: NewRegistry()}
	result := dc.checkProviders()

	if len(result.Checks) != 1 {
		t.Fatalf("expected 1 check, got %d", len(result.Checks))
	}
	if result.Checks[0].Status != CheckFail {
		t.Errorf("status = %v, want %v", result.Checks[0].Status, CheckFail)
	}
}

func TestCheckProviders_InvalidConfig(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	writeConfig(t, tmpDir, "{{{{invalid yaml")

	dc := &DoctorChecker{configDir: tmpDir, registry: NewRegistry()}
	result := dc.checkProviders()

	if len(result.Checks) != 1 {
		t.Fatalf("expected 1 check, got %d", len(result.Checks))
	}
	if result.Checks[0].Status != CheckFail {
		t.Errorf("status = %v, want %v", result.Checks[0].Status, CheckFail)
	}
}

func TestCheckProviders_LegacyFlatConfig(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	writeConfig(t, tmpDir, "schema_version: 2\nprovider: textfile\n")

	reg := newTestRegistry(map[string]*doctorMockProvider{
		"textfile": {
			name:  "textfile",
			tasks: []*Task{{ID: "1", Text: "hello"}},
		},
	})

	dc := &DoctorChecker{configDir: tmpDir, registry: reg}
	result := dc.checkProviders()

	if len(result.Checks) != 1 {
		t.Fatalf("expected 1 check, got %d", len(result.Checks))
	}
	if result.Checks[0].Status != CheckOK {
		t.Errorf("status = %v, want %v", result.Checks[0].Status, CheckOK)
	}
}

func TestCheckProviders_ProvidersListConfig(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	writeConfig(t, tmpDir, `schema_version: 2
providers:
  - name: textfile
  - name: applenotes
`)

	reg := newTestRegistry(map[string]*doctorMockProvider{
		"textfile":   {name: "textfile", tasks: []*Task{{ID: "1"}}},
		"applenotes": {name: "applenotes", tasks: []*Task{{ID: "2"}}},
	})

	dc := &DoctorChecker{configDir: tmpDir, registry: reg}
	result := dc.checkProviders()

	if len(result.Checks) != 2 {
		t.Fatalf("expected 2 checks, got %d", len(result.Checks))
	}
	for _, check := range result.Checks {
		if check.Status != CheckOK {
			t.Errorf("check %q: status = %v, want %v", check.Name, check.Status, CheckOK)
		}
	}
}

func TestCheckProviders_IntegrationWithRun(t *testing.T) {
	tmpDir := t.TempDir()
	writeConfig(t, tmpDir, "schema_version: 2\nprovider: textfile\n")

	reg := newTestRegistry(map[string]*doctorMockProvider{
		"textfile": {
			name:  "textfile",
			tasks: []*Task{{ID: "1", Text: "test task"}},
		},
	})

	dc := NewDoctorChecker(tmpDir)
	dc.SetRegistry(reg)
	result := dc.Run()

	// Should have Environment + Providers categories
	if len(result.Categories) != 2 {
		t.Fatalf("expected 2 categories, got %d", len(result.Categories))
	}

	providers := result.Categories[1]
	if providers.Name != "Providers" {
		t.Errorf("second category name = %q, want %q", providers.Name, "Providers")
	}
	if providers.Status != CheckOK {
		t.Errorf("providers status = %v, want %v", providers.Status, CheckOK)
	}
}

func TestCheckProviders_ObsidianNoVaultPath(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	writeConfig(t, tmpDir, `schema_version: 2
providers:
  - name: obsidian
`)

	reg := newTestRegistry(map[string]*doctorMockProvider{
		"obsidian": {name: "obsidian"},
	})

	dc := &DoctorChecker{configDir: tmpDir, registry: reg}
	result := dc.checkProviders()

	if len(result.Checks) != 1 {
		t.Fatalf("expected 1 check, got %d", len(result.Checks))
	}
	if result.Checks[0].Status != CheckFail {
		t.Errorf("status = %v, want %v", result.Checks[0].Status, CheckFail)
	}
	if result.Checks[0].Message != "Obsidian vault path not configured" {
		t.Errorf("message = %q", result.Checks[0].Message)
	}
}
