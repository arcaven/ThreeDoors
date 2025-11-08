package tasks

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestLoadTasks_NoFileExists(t *testing.T) {
	tempDir := t.TempDir()
	SetHomeDir(tempDir)  // Redirect home directory for testing
	defer SetHomeDir("") // Reset after test

	tasks, err := LoadTasks()
	if err != nil {
		t.Fatalf("LoadTasks() failed: %v", err)
	}

	if len(tasks) != len(defaultTasks) {
		t.Errorf("Expected %d default tasks, got %d", len(defaultTasks), len(tasks))
	}

	for i, task := range tasks {
		if task.Text != defaultTasks[i] {
			t.Errorf("Expected task %d to be '%s', got '%s'", i, defaultTasks[i], task.Text)
		}
	}

	// Verify file and directory were created
	configPath := filepath.Join(tempDir, configDir)
	tasksFilePath := filepath.Join(configPath, tasksFileName)

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Errorf("Config directory was not created at %s", configPath)
	}
	if _, err := os.Stat(tasksFilePath); os.IsNotExist(err) {
		t.Errorf("Tasks file was not created at %s", tasksFilePath)
	}

	// Verify content of the created file
	content, err := os.ReadFile(tasksFilePath)
	if err != nil {
		t.Fatalf("Failed to read created tasks file: %v", err)
	}
	expectedContent := strings.Join(defaultTasks, "\n") + "\n"
	if string(content) != expectedContent {
		t.Errorf("Created file content mismatch.\nExpected:\n%s\nGot:\n%s", expectedContent, string(content))
	}
}

func TestLoadTasks_FileExistsWithContent(t *testing.T) {
	tempDir := t.TempDir()
	SetHomeDir(tempDir)
	defer SetHomeDir("")

	configPath := filepath.Join(tempDir, configDir)
	tasksFilePath := filepath.Join(configPath, tasksFileName)

	// Create directory and file with custom content
	os.MkdirAll(configPath, 0755)
	customTasks := []string{"Custom Task A", "Custom Task B", "Custom Task C"}
	err := os.WriteFile(tasksFilePath, []byte(strings.Join(customTasks, "\n")),
		0644)
	if err != nil {
		t.Fatalf("Failed to write custom tasks file: %v", err)
	}

	tasks, err := LoadTasks()
	if err != nil {
		t.Fatalf("LoadTasks() failed: %v", err)
	}

	expectedTasks := []Task{
		{Text: "Custom Task A"},
		{Text: "Custom Task B"},
		{Text: "Custom Task C"},
	}

	if !reflect.DeepEqual(tasks, expectedTasks) {
		t.Errorf("Loaded tasks mismatch.\nExpected: %+v\nGot: %+v", expectedTasks, tasks)
	}
}

func TestLoadTasks_EmptyFile(t *testing.T) {
	tempDir := t.TempDir()
	SetHomeDir(tempDir)
	defer SetHomeDir("")

	configPath := filepath.Join(tempDir, configDir)
	tasksFilePath := filepath.Join(configPath, tasksFileName)

	// Create directory and an empty file
	os.MkdirAll(configPath, 0755)
	err := os.WriteFile(tasksFilePath, []byte(""), 0644)
	if err != nil {
		t.Fatalf("Failed to write empty tasks file: %v", err)
	}

	tasks, err := LoadTasks()
	if err != nil {
		t.Fatalf("LoadTasks() failed: %v", err)
	}

	if len(tasks) != 0 {
		t.Errorf("Expected 0 tasks from empty file, got %d", len(tasks))
	}
}

func TestLoadTasks_FileWithBlankLines(t *testing.T) {
	tempDir := t.TempDir()
	SetHomeDir(tempDir)
	defer SetHomeDir("")

	configPath := filepath.Join(tempDir, configDir)
	tasksFilePath := filepath.Join(configPath, tasksFileName)

	// Create directory and file with blank lines
	os.MkdirAll(configPath, 0755)
	content := "Task 1\n\nTask 2\n   \nTask 3"
	err := os.WriteFile(tasksFilePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write tasks file with blank lines: %v", err)
	}

	tasks, err := LoadTasks()
	if err != nil {
		t.Fatalf("LoadTasks() failed: %v", err)
	}

	expectedTasks := []Task{
		{Text: "Task 1"},
		{Text: "Task 2"},
		{Text: "Task 3"},
	}

	if !reflect.DeepEqual(tasks, expectedTasks) {
		t.Errorf("Loaded tasks mismatch.\nExpected: %+v\nGot: %+v", expectedTasks, tasks)
	}
}
