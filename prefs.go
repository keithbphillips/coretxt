package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type prefs struct {
	ThemeIdx int    `json:"theme_idx"`
	LastDir  string `json:"last_dir"`
	LastFile string `json:"last_file"`
}

func prefsPath() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		dir = os.Getenv("HOME")
	}
	return filepath.Join(dir, "coretxt", "prefs.json")
}

func loadPrefs() prefs {
	data, err := os.ReadFile(prefsPath())
	if err != nil {
		return prefs{}
	}
	var p prefs
	if err := json.Unmarshal(data, &p); err != nil {
		return prefs{}
	}
	// Guard against out-of-range index if themes list ever shrinks
	if p.ThemeIdx < 0 || p.ThemeIdx >= len(themes) {
		p.ThemeIdx = 0
	}
	return p
}

func savePrefs(p prefs) {
	path := prefsPath()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return
	}
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(path, data, 0644)
}
