package main

import (
	"os"
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
