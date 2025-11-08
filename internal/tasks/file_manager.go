package tasks

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	tasksFileName = "tasks.txt"
	configDir     = ".threedoors"
)

var testHomeDir string // Used for testing purposes to override os.UserHomeDir

// SetHomeDir sets the home directory for testing purposes.
// Pass an empty string to reset to default os.UserHomeDir.
func SetHomeDir(dir string) {
	testHomeDir = dir
}

// GetConfigDirPath returns the full path to the application's configuration directory.
func GetConfigDirPath() (string, error) {
	var homeDir string
	if testHomeDir != "" {
		homeDir = testHomeDir
	} else {
		var err error
		homeDir, err = os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user home directory: %w", err)
		}
	}
	return filepath.Join(homeDir, configDir), nil
}

// GetTasksFilePath returns the full path to the tasks file.
func GetTasksFilePath() (string, error) {
	configPath, err := GetConfigDirPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(configPath, tasksFileName), nil
}

// defaultTasks provides a set of sample tasks for initial file creation.
var defaultTasks = []string{
	"Task 1: Learn Go",
	"Task 2: Build a TUI app",
	"Task 3: Explore Bubbletea",
	"Task 4: Write tests",
	"Task 5: Deploy application",
}

// LoadTasks reads tasks from the tasks file. If the file or directory does not exist,
// it creates them with default tasks.
func LoadTasks() ([]Task, error) {
	tasksFilePath, err := GetTasksFilePath()
	if err != nil {
		return nil, err
	}

	configDirPath, err := GetConfigDirPath()
	if err != nil {
		return nil, err
	}

	// Ensure the config directory exists
	if _, err := os.Stat(configDirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(configDirPath, 0755); err != nil {
			return nil, fmt.Errorf("failed to create config directory %s: %w", configDirPath, err)
		}
	}

	// Check if tasks file exists, if not, create it with default tasks
	if _, err := os.Stat(tasksFilePath); os.IsNotExist(err) {
		file, err := os.Create(tasksFilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to create tasks file %s: %w", tasksFilePath, err)
		}
		defer file.Close()

		writer := bufio.NewWriter(file)
		for _, task := range defaultTasks {
			_, err := writer.WriteString(task + "\n")
			if err != nil {
				return nil, fmt.Errorf("failed to write default task to file: %w", err)
			}
		}
		writer.Flush()
	} else if err != nil {
		return nil, fmt.Errorf("failed to stat tasks file %s: %w", tasksFilePath, err)
	}

	// Read tasks from the file
	file, err := os.Open(tasksFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open tasks file %s: %w", tasksFilePath, err)
	}
	defer file.Close()

	var tasks []Task
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text != "" {
			tasks = append(tasks, Task{Text: text})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading tasks file: %w", err)
	}

	return tasks, nil
}
