package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/arcaven/ThreeDoors/internal/core"
	"github.com/arcaven/ThreeDoors/internal/core/connection"
	"github.com/spf13/cobra"
)

// connectResult holds the outcome of a connect operation for output rendering.
type connectResult struct {
	ID       string            `json:"id"`
	Provider string            `json:"provider"`
	Label    string            `json:"label"`
	State    string            `json:"state"`
	Health   *healthResult     `json:"health,omitempty"`
	Settings map[string]string `json:"settings,omitempty"`
}

// healthResult holds health check output.
type healthResult struct {
	Healthy      bool              `json:"healthy"`
	APIReachable bool              `json:"api_reachable"`
	TokenValid   bool              `json:"token_valid"`
	RateLimitOK  bool              `json:"rate_limit_ok"`
	TaskCount    int               `json:"task_count"`
	Details      map[string]string `json:"details,omitempty"`
}

func newConnectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connect <provider>",
		Short: "Connect a data source",
		Long: `Connect a data source provider to ThreeDoors.

Supported providers: todoist, github, jira, textfile

Use flags for non-interactive (scriptable) setup.
Run without flags to launch the interactive wizard (coming soon).`,
	}
	cmd.AddCommand(newConnectTodoistCmd())
	cmd.AddCommand(newConnectGithubCmd())
	cmd.AddCommand(newConnectJiraCmd())
	cmd.AddCommand(newConnectTextfileCmd())
	return cmd
}

func newConnectTodoistCmd() *cobra.Command {
	var label, token, projectIDs, filter string

	cmd := &cobra.Command{
		Use:   "todoist",
		Short: "Connect a Todoist account",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !hasAnyFlag(cmd, "label", "token") {
				return fmt.Errorf("interactive mode is not yet available; use --label and --token flags\n\nUsage: threedoors connect todoist --label <name> --token <token>")
			}
			missing := requireFlags(cmd, "label")
			if len(missing) > 0 {
				return newMissingFlagsError(missing)
			}

			// Token from flag or env
			if token == "" {
				token = os.Getenv("TODOIST_API_TOKEN")
			}
			if token == "" {
				return newMissingFlagsError([]string{"--token (or set TODOIST_API_TOKEN env var)"})
			}

			settings := make(map[string]string)
			if projectIDs != "" {
				settings["project_ids"] = projectIDs
			}
			if filter != "" {
				settings["filter"] = filter
			}

			return runConnect(cmd, "todoist", label, token, settings)
		},
	}
	cmd.Flags().StringVar(&label, "label", "", "connection label (e.g. \"Personal\")")
	cmd.Flags().StringVar(&token, "token", "", "Todoist API token (or set TODOIST_API_TOKEN)")
	cmd.Flags().StringVar(&projectIDs, "project-ids", "", "comma-separated project IDs to sync")
	cmd.Flags().StringVar(&filter, "filter", "", "Todoist filter expression")
	return cmd
}

func newConnectGithubCmd() *cobra.Command {
	var label, token, repos string

	cmd := &cobra.Command{
		Use:   "github",
		Short: "Connect GitHub Issues",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !hasAnyFlag(cmd, "label", "repos", "token") {
				return fmt.Errorf("interactive mode is not yet available; use --label and --repos flags\n\nUsage: threedoors connect github --label <name> --repos owner/repo1,owner/repo2")
			}
			missing := requireFlags(cmd, "label", "repos")
			if len(missing) > 0 {
				return newMissingFlagsError(missing)
			}

			// Token from flag, then env vars
			if token == "" {
				token = os.Getenv("GH_TOKEN")
			}
			if token == "" {
				token = os.Getenv("GITHUB_TOKEN")
			}

			settings := map[string]string{
				"repos": repos,
			}

			return runConnect(cmd, "github", label, token, settings)
		},
	}
	cmd.Flags().StringVar(&label, "label", "", "connection label (e.g. \"OSS\")")
	cmd.Flags().StringVar(&token, "token", "", "GitHub token (or set GH_TOKEN / GITHUB_TOKEN)")
	cmd.Flags().StringVar(&repos, "repos", "", "comma-separated repos in owner/repo format")
	return cmd
}

func newConnectJiraCmd() *cobra.Command {
	var label, token, server, jql, email, authType string

	cmd := &cobra.Command{
		Use:   "jira",
		Short: "Connect a Jira instance",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !hasAnyFlag(cmd, "label", "server", "token") {
				return fmt.Errorf("interactive mode is not yet available; use --label, --server, and --token flags\n\nUsage: threedoors connect jira --label <name> --server <url> --token <token>")
			}
			missing := requireFlags(cmd, "label", "server")
			if len(missing) > 0 {
				return newMissingFlagsError(missing)
			}

			// Token from flag or env
			if token == "" {
				token = os.Getenv("JIRA_API_TOKEN")
			}
			if token == "" {
				return newMissingFlagsError([]string{"--token (or set JIRA_API_TOKEN env var)"})
			}

			if authType == "" {
				authType = "pat"
			}

			settings := map[string]string{
				"url":       server,
				"auth_type": authType,
			}
			if jql != "" {
				settings["jql"] = jql
			}
			if email != "" {
				settings["email"] = email
			}

			return runConnect(cmd, "jira", label, token, settings)
		},
	}
	cmd.Flags().StringVar(&label, "label", "", "connection label (e.g. \"Work\")")
	cmd.Flags().StringVar(&token, "token", "", "Jira API token (or set JIRA_API_TOKEN)")
	cmd.Flags().StringVar(&server, "server", "", "Jira server URL (e.g. https://jira.company.com)")
	cmd.Flags().StringVar(&jql, "jql", "", "custom JQL filter")
	cmd.Flags().StringVar(&email, "email", "", "Jira email (for basic auth)")
	cmd.Flags().StringVar(&authType, "auth-type", "", "auth type: basic or pat (default: pat)")
	return cmd
}

func newConnectTextfileCmd() *cobra.Command {
	var label, path string

	cmd := &cobra.Command{
		Use:   "textfile",
		Short: "Connect a local text file",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !hasAnyFlag(cmd, "label", "path") {
				return fmt.Errorf("interactive mode is not yet available; use --label and --path flags\n\nUsage: threedoors connect textfile --label <name> --path <filepath>")
			}
			missing := requireFlags(cmd, "label", "path")
			if len(missing) > 0 {
				return newMissingFlagsError(missing)
			}

			settings := map[string]string{
				"path": path,
			}

			return runConnect(cmd, "textfile", label, "", settings)
		},
	}
	cmd.Flags().StringVar(&label, "label", "", "connection label (e.g. \"Notes\")")
	cmd.Flags().StringVar(&path, "path", "", "path to YAML task file")
	return cmd
}

// runConnect creates a connection, optionally tests it, and prints the result.
func runConnect(cmd *cobra.Command, providerName, label, credential string, settings map[string]string) error {
	isJSON := isJSONOutput(cmd)
	formatter := NewOutputFormatter(cmd.OutOrStdout(), isJSON)

	svc, err := buildConnectionServiceFn()
	if err != nil {
		if isJSON {
			_ = formatter.WriteJSONError("connect", ExitGeneralError, err.Error(), "")
		}
		return err
	}

	conn, err := svc.Add(providerName, label, settings, credential)
	if err != nil {
		if isJSON {
			_ = formatter.WriteJSONError("connect", ExitGeneralError, err.Error(), "")
		}
		return fmt.Errorf("create connection: %w", err)
	}

	result := connectResult{
		ID:       conn.ID,
		Provider: conn.ProviderName,
		Label:    conn.Label,
		State:    conn.State.String(),
		Settings: redactSettings(settings),
	}

	// Attempt health check (non-fatal if unavailable or fails)
	hr, testErr := svc.TestConnection(conn.ID)
	if testErr == nil {
		result.Health = &healthResult{
			Healthy:      hr.Healthy(),
			APIReachable: hr.APIReachable,
			TokenValid:   hr.TokenValid,
			RateLimitOK:  hr.RateLimitOK,
			TaskCount:    hr.TaskCount,
			Details:      hr.Details,
		}
	}

	if isJSON {
		return formatter.WriteJSON("connect", result, nil)
	}

	return writeConnectTable(formatter, result, testErr)
}

// writeConnectTable renders the connect result in human-readable format.
func writeConnectTable(formatter *OutputFormatter, result connectResult, testErr error) error {
	if err := formatter.Writef("Connection created successfully.\n\n"); err != nil {
		return err
	}

	tw := formatter.TableWriter()
	_, _ = fmt.Fprintf(tw, "ID\t%s\n", result.ID)
	_, _ = fmt.Fprintf(tw, "Provider\t%s\n", result.Provider)
	_, _ = fmt.Fprintf(tw, "Label\t%s\n", result.Label)
	_, _ = fmt.Fprintf(tw, "State\t%s\n", result.State)
	if err := tw.Flush(); err != nil {
		return err
	}

	if result.Health != nil {
		if err := formatter.Writef("\nHealth Check:\n"); err != nil {
			return err
		}
		tw = formatter.TableWriter()
		_, _ = fmt.Fprintf(tw, "  Healthy\t%t\n", result.Health.Healthy)
		_, _ = fmt.Fprintf(tw, "  API Reachable\t%t\n", result.Health.APIReachable)
		_, _ = fmt.Fprintf(tw, "  Token Valid\t%t\n", result.Health.TokenValid)
		_, _ = fmt.Fprintf(tw, "  Rate Limit OK\t%t\n", result.Health.RateLimitOK)
		_, _ = fmt.Fprintf(tw, "  Task Count\t%d\n", result.Health.TaskCount)
		return tw.Flush()
	}

	if testErr != nil {
		return formatter.Writef("\nHealth check skipped: %v\n", testErr)
	}

	return nil
}

// buildConnectionServiceFn is the function used to create a ConnectionService.
// It can be replaced in tests to inject mock services.
var buildConnectionServiceFn = buildConnectionService

// buildConnectionService creates a ConnectionService for CLI use.
func buildConnectionService() (*connection.ConnectionService, error) {
	configDir, err := core.GetConfigDirPath()
	if err != nil {
		return nil, fmt.Errorf("config dir: %w", err)
	}
	configPath := filepath.Join(configDir, "config.yaml")

	manager := connection.NewConnectionManager(nil)
	creds := connection.NewEnvCredentialStore()

	// Load existing connections so the manager has full state for persistence.
	cfg, err := core.LoadProviderConfig(configPath)
	if err == nil {
		for _, cc := range cfg.Connections {
			addExistingConnection(manager, cc)
		}
	}

	svc, err := connection.NewConnectionService(connection.ServiceConfig{
		Manager:    manager,
		Creds:      creds,
		ConfigPath: configPath,
	})
	if err != nil {
		return nil, fmt.Errorf("create service: %w", err)
	}
	return svc, nil
}

// addExistingConnection registers an existing connection from config into the manager.
func addExistingConnection(manager *connection.ConnectionManager, cc core.ConnectionConfig) {
	// We don't use manager.Add() because we want to preserve the existing ID
	_, _ = manager.Add(cc.Provider, cc.Label, cc.Settings)
}

// redactSettings returns a copy of settings with sensitive keys masked.
func redactSettings(settings map[string]string) map[string]string {
	redacted := make(map[string]string, len(settings))
	for k, v := range settings {
		if strings.Contains(strings.ToLower(k), "token") || strings.Contains(strings.ToLower(k), "secret") {
			redacted[k] = connection.MaskCredential(v)
		} else {
			redacted[k] = v
		}
	}
	return redacted
}

// hasAnyFlag checks if any of the named flags were explicitly set by the user.
func hasAnyFlag(cmd *cobra.Command, names ...string) bool {
	for _, name := range names {
		if cmd.Flags().Changed(name) {
			return true
		}
	}
	return false
}

// requireFlags checks that the named flags were set and returns the missing ones.
func requireFlags(cmd *cobra.Command, names ...string) []string {
	var missing []string
	for _, name := range names {
		if !cmd.Flags().Changed(name) {
			missing = append(missing, "--"+name)
		}
	}
	return missing
}

// newMissingFlagsError returns an error listing missing required flags.
func newMissingFlagsError(missing []string) error {
	return fmt.Errorf("missing required flags: %s", strings.Join(missing, ", "))
}
