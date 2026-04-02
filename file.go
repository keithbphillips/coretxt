package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const autosaveInterval = 30 * time.Second
const backupWordThreshold = 200
const maxBackups = 5

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

// ─── Backups ──────────────────────────────────────────────────────────────────

// saveBackup writes a timestamped copy of content into a .coretxt_backups
// directory next to mainPath, then prunes old backups to keep at most maxBackups.
func saveBackup(mainPath string, content string) error {
	dir := filepath.Join(filepath.Dir(mainPath), ".coretxt_backups")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	base := strings.TrimSuffix(filepath.Base(mainPath), filepath.Ext(mainPath))
	ext := filepath.Ext(mainPath)
	timestamp := time.Now().Format("20060102_150405")
	backupPath := filepath.Join(dir, base+"_"+timestamp+ext)
	if err := os.WriteFile(backupPath, []byte(content), 0644); err != nil {
		return err
	}
	pruneBackups(dir, base, ext)
	return nil
}

// pruneBackups removes the oldest backup files so that at most maxBackups remain.
func pruneBackups(dir, base, ext string) {
	matches, err := filepath.Glob(filepath.Join(dir, base+"_*"+ext))
	if err != nil || len(matches) <= maxBackups {
		return
	}
	sort.Strings(matches) // timestamp names sort oldest-first
	for _, old := range matches[:len(matches)-maxBackups] {
		os.Remove(old)
	}
}
