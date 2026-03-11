package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/arcaven/ThreeDoors/internal/core"
	"github.com/arcaven/ThreeDoors/internal/core/connection"
	"github.com/spf13/cobra"
)

// testCredentialStore is an in-memory credential store for testing.
type testCredentialStore struct {
	creds map[string]string
}

func newTestCredentialStore() *testCredentialStore {
	return &testCredentialStore{creds: make(map[string]string)}
}

func (s *testCredentialStore) Get(connID string) (string, error) {
	v, ok := s.creds[connID]
	if !ok {
		return "", connection.ErrCredentialNotFound
	}
	return v, nil
}

func (s *testCredentialStore) Set(connID, value string) error {
	s.creds[connID] = value
	return nil
}

func (s *testCredentialStore) Delete(connID string) error {
	delete(s.creds, connID)
	return nil
}

// testConnectService creates a ConnectionService backed by in-memory stores.
func testConnectService(t *testing.T) (*connection.ConnectionService, *testCredentialStore) {
	t.Helper()
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	cfg := &core.ProviderConfig{SchemaVersion: 1, Provider: "textfile"}
	if err := core.SaveProviderConfig(configPath, cfg); err != nil {
		t.Fatalf("save test config: %v", err)
	}

	manager := connection.NewConnectionManager(nil)
	creds := newTestCredentialStore()
	svc, err := connection.NewConnectionService(connection.ServiceConfig{
		Manager:    manager,
		Creds:      creds,
		ConfigPath: configPath,
	})
	if err != nil {
		t.Fatalf("create test service: %v", err)
	}
	return svc, creds
}

// withTestService sets up a test ConnectionService and overrides the global
// builder. Must not be used in parallel tests.
func withTestService(t *testing.T) (*connection.ConnectionService, *testCredentialStore) {
	t.Helper()
	svc, creds := testConnectService(t)
	origBuild := buildConnectionServiceFn
	buildConnectionServiceFn = func() (*connection.ConnectionService, error) {
		return svc, nil
	}
	t.Cleanup(func() { buildConnectionServiceFn = origBuild })
	return svc, creds
}

// executeConnect builds a root command, wires the connect subcommand,
// sets args, and captures output.
func executeConnect(t *testing.T, args []string) (string, error) {
	t.Helper()
	var buf bytes.Buffer

	root := &cobra.Command{Use: "threedoors"}
	root.PersistentFlags().Bool("json", false, "output in JSON format")
	root.AddCommand(newConnectCmd())
	root.SetOut(&buf)
	root.SetErr(&buf)
	root.SetArgs(args)

	err := root.Execute()
	return buf.String(), err
}

func TestConnectTodoist(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		env        map[string]string
		wantErr    bool
		wantSubstr string
	}{
		{
			name:       "valid with token flag",
			args:       []string{"connect", "todoist", "--label", "Personal", "--token", "test-token-123"},
			wantSubstr: "Connection created successfully",
		},
		{
			name:       "valid with env token",
			args:       []string{"connect", "todoist", "--label", "Work"},
			env:        map[string]string{"TODOIST_API_TOKEN": "env-token"},
			wantSubstr: "Connection created successfully",
		},
		{
			name:       "with project-ids",
			args:       []string{"connect", "todoist", "--label", "Filtered", "--token", "tok", "--project-ids", "123,456"},
			wantSubstr: "Connection created successfully",
		},
		{
			name:       "with filter",
			args:       []string{"connect", "todoist", "--label", "Filtered", "--token", "tok", "--filter", "today"},
			wantSubstr: "Connection created successfully",
		},
		{
			name:       "missing label",
			args:       []string{"connect", "todoist", "--token", "tok"},
			wantErr:    true,
			wantSubstr: "--label",
		},
		{
			name:       "missing token",
			args:       []string{"connect", "todoist", "--label", "Test"},
			wantErr:    true,
			wantSubstr: "--token",
		},
		{
			name:       "no flags interactive fallback",
			args:       []string{"connect", "todoist"},
			wantErr:    true,
			wantSubstr: "interactive mode is not yet available",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withTestService(t)
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			out, err := executeConnect(t, tt.args)

			if tt.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v\noutput: %s", err, out)
			}
			if tt.wantSubstr != "" {
				combined := out + errStr(err)
				if !strings.Contains(combined, tt.wantSubstr) {
					t.Errorf("output %q does not contain %q", combined, tt.wantSubstr)
				}
			}
		})
	}
}

func TestConnectGithub(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		env        map[string]string
		wantErr    bool
		wantSubstr string
	}{
		{
			name:       "valid with repos and token",
			args:       []string{"connect", "github", "--label", "OSS", "--repos", "owner/repo1,owner/repo2", "--token", "gh-tok"},
			wantSubstr: "Connection created successfully",
		},
		{
			name:       "valid with GH_TOKEN env",
			args:       []string{"connect", "github", "--label", "OSS", "--repos", "owner/repo1"},
			env:        map[string]string{"GH_TOKEN": "env-gh-tok"},
			wantSubstr: "Connection created successfully",
		},
		{
			name:       "valid with GITHUB_TOKEN env",
			args:       []string{"connect", "github", "--label", "OSS", "--repos", "owner/repo1"},
			env:        map[string]string{"GITHUB_TOKEN": "env-github-tok"},
			wantSubstr: "Connection created successfully",
		},
		{
			name:       "missing repos",
			args:       []string{"connect", "github", "--label", "Test", "--token", "tok"},
			wantErr:    true,
			wantSubstr: "--repos",
		},
		{
			name:       "missing label",
			args:       []string{"connect", "github", "--repos", "owner/repo"},
			wantErr:    true,
			wantSubstr: "--label",
		},
		{
			name:       "no flags interactive fallback",
			args:       []string{"connect", "github"},
			wantErr:    true,
			wantSubstr: "interactive mode is not yet available",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withTestService(t)
			for k, v := range tt.env {
				t.Setenv(k, v)
			}
			out, err := executeConnect(t, tt.args)
			if tt.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v\noutput: %s", err, out)
			}
			if tt.wantSubstr != "" {
				combined := out + errStr(err)
				if !strings.Contains(combined, tt.wantSubstr) {
					t.Errorf("output %q does not contain %q", combined, tt.wantSubstr)
				}
			}
		})
	}
}

func TestConnectJira(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		env        map[string]string
		wantErr    bool
		wantSubstr string
	}{
		{
			name:       "valid with all flags",
			args:       []string{"connect", "jira", "--label", "Work", "--server", "https://jira.company.com", "--token", "jira-tok"},
			wantSubstr: "Connection created successfully",
		},
		{
			name:       "valid with env token",
			args:       []string{"connect", "jira", "--label", "Work", "--server", "https://jira.company.com"},
			env:        map[string]string{"JIRA_API_TOKEN": "env-jira-tok"},
			wantSubstr: "Connection created successfully",
		},
		{
			name:       "with optional jql and email",
			args:       []string{"connect", "jira", "--label", "Work", "--server", "https://jira.co", "--token", "tok", "--jql", "project = FOO", "--email", "me@co.com"},
			wantSubstr: "Connection created successfully",
		},
		{
			name:       "missing server",
			args:       []string{"connect", "jira", "--label", "Work", "--token", "tok"},
			wantErr:    true,
			wantSubstr: "--server",
		},
		{
			name:       "missing token",
			args:       []string{"connect", "jira", "--label", "Work", "--server", "https://jira.co"},
			wantErr:    true,
			wantSubstr: "--token",
		},
		{
			name:       "no flags interactive fallback",
			args:       []string{"connect", "jira"},
			wantErr:    true,
			wantSubstr: "interactive mode is not yet available",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withTestService(t)
			for k, v := range tt.env {
				t.Setenv(k, v)
			}
			out, err := executeConnect(t, tt.args)
			if tt.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v\noutput: %s", err, out)
			}
			if tt.wantSubstr != "" {
				combined := out + errStr(err)
				if !strings.Contains(combined, tt.wantSubstr) {
					t.Errorf("output %q does not contain %q", combined, tt.wantSubstr)
				}
			}
		})
	}
}

func TestConnectTextfile(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantErr    bool
		wantSubstr string
	}{
		{
			name:       "valid with label and path",
			args:       []string{"connect", "textfile", "--label", "Notes", "--path", "~/tasks.yaml"},
			wantSubstr: "Connection created successfully",
		},
		{
			name:       "missing path",
			args:       []string{"connect", "textfile", "--label", "Notes"},
			wantErr:    true,
			wantSubstr: "--path",
		},
		{
			name:       "missing label",
			args:       []string{"connect", "textfile", "--path", "~/tasks.yaml"},
			wantErr:    true,
			wantSubstr: "--label",
		},
		{
			name:       "no flags interactive fallback",
			args:       []string{"connect", "textfile"},
			wantErr:    true,
			wantSubstr: "interactive mode is not yet available",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withTestService(t)
			out, err := executeConnect(t, tt.args)
			if tt.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v\noutput: %s", err, out)
			}
			if tt.wantSubstr != "" {
				combined := out + errStr(err)
				if !strings.Contains(combined, tt.wantSubstr) {
					t.Errorf("output %q does not contain %q", combined, tt.wantSubstr)
				}
			}
		})
	}
}

func TestConnectJSONOutput(t *testing.T) {
	withTestService(t)

	out, err := executeConnect(t, []string{"connect", "todoist", "--label", "JSON-Test", "--token", "tok", "--json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var env JSONEnvelope
	if err := json.Unmarshal([]byte(out), &env); err != nil {
		t.Fatalf("invalid JSON output: %v\nraw: %s", err, out)
	}

	if env.SchemaVersion != 1 {
		t.Errorf("schema_version = %d, want 1", env.SchemaVersion)
	}
	if env.Command != "connect" {
		t.Errorf("command = %q, want %q", env.Command, "connect")
	}
	if env.Error != nil {
		t.Errorf("unexpected error in JSON: %+v", env.Error)
	}

	data, ok := env.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("data is not a map: %T", env.Data)
	}
	if data["provider"] != "todoist" {
		t.Errorf("provider = %v, want todoist", data["provider"])
	}
	if data["label"] != "JSON-Test" {
		t.Errorf("label = %v, want JSON-Test", data["label"])
	}
	if data["id"] == nil || data["id"] == "" {
		t.Errorf("id should not be empty")
	}
}

func TestConnectCredentialStored(t *testing.T) {
	_, creds := withTestService(t)

	_, err := executeConnect(t, []string{"connect", "todoist", "--label", "CredTest", "--token", "my-secret-token"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found := false
	for _, v := range creds.creds {
		if v == "my-secret-token" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("credential was not stored; stored keys: %v", creds.creds)
	}
}

func TestConnectSettings(t *testing.T) {
	withTestService(t)

	out, err := executeConnect(t, []string{"connect", "github", "--label", "SettingsTest", "--repos", "owner/repo1,owner/repo2", "--token", "tok", "--json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var env JSONEnvelope
	if err := json.Unmarshal([]byte(out), &env); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	data, ok := env.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("data is not a map")
	}
	settings, ok := data["settings"].(map[string]interface{})
	if !ok {
		t.Fatalf("settings is not a map: %T", data["settings"])
	}
	if settings["repos"] != "owner/repo1,owner/repo2" {
		t.Errorf("repos = %v, want owner/repo1,owner/repo2", settings["repos"])
	}
}

func TestRedactSettings(t *testing.T) {
	t.Parallel()

	settings := map[string]string{
		"repos":     "owner/repo",
		"api_token": "secret123",
		"url":       "https://example.com",
	}

	redacted := redactSettings(settings)

	if redacted["repos"] != "owner/repo" {
		t.Errorf("repos should not be redacted, got %q", redacted["repos"])
	}
	if redacted["api_token"] != "••••" {
		t.Errorf("api_token should be redacted, got %q", redacted["api_token"])
	}
	if redacted["url"] != "https://example.com" {
		t.Errorf("url should not be redacted, got %q", redacted["url"])
	}
}

func TestHasAnyFlag(t *testing.T) {
	t.Parallel()

	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("label", "", "")
	cmd.Flags().String("token", "", "")

	if hasAnyFlag(cmd, "label", "token") {
		t.Error("hasAnyFlag should return false when no flags set")
	}

	if err := cmd.Flags().Set("label", "test"); err != nil {
		t.Fatal(err)
	}
	if !hasAnyFlag(cmd, "label", "token") {
		t.Error("hasAnyFlag should return true when label is set")
	}
}

func TestRequireFlags(t *testing.T) {
	t.Parallel()

	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("label", "", "")
	cmd.Flags().String("token", "", "")

	missing := requireFlags(cmd, "label", "token")
	if len(missing) != 2 {
		t.Errorf("expected 2 missing, got %d: %v", len(missing), missing)
	}

	if err := cmd.Flags().Set("label", "test"); err != nil {
		t.Fatal(err)
	}
	missing = requireFlags(cmd, "label", "token")
	if len(missing) != 1 {
		t.Errorf("expected 1 missing, got %d: %v", len(missing), missing)
	}
	if missing[0] != "--token" {
		t.Errorf("expected --token missing, got %s", missing[0])
	}
}

func TestConnectKnownSubcommands(t *testing.T) {
	t.Parallel()

	subs := KnownSubcommands()
	found := false
	for _, s := range subs {
		if s == "connect" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("connect not in KnownSubcommands: %v", subs)
	}
}

// errStr returns the error string or empty string for nil.
func errStr(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}

// TestMain clears env vars that could pollute tests.
func TestMain(m *testing.M) {
	for _, key := range []string{"TODOIST_API_TOKEN", "GH_TOKEN", "GITHUB_TOKEN", "JIRA_API_TOKEN"} {
		_ = os.Unsetenv(key)
	}
	os.Exit(m.Run())
}
