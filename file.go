package main

import (
	"os"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const autosaveInterval = 30 * time.Second

// ─── Tick messages ────────────────────────────────────────────────────────────

type tickMsg time.Time
type statusClearMsg struct{}

func doTick() tea.Cmd {
	return tea.Tick(autosaveInterval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func clearStatus(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return statusClearMsg{}
	})
}

// ─── Paths ────────────────────────────────────────────────────────────────────

// docsDir returns ~/Documents, creating it if necessary.
func docsDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	dir := filepath.Join(home, "Documents")
	_ = os.MkdirAll(dir, 0755)
	return dir
}

// resolvePath returns the path unchanged if it is absolute, otherwise joins it
// with docsDir so bare filenames always map to ~/Documents.
func resolvePath(name string) string {
	if filepath.IsAbs(name) {
		return name
	}
	return filepath.Join(docsDir(), name)
}

// ─── File I/O ─────────────────────────────────────────────────────────────────

func loadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func saveFile(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}
