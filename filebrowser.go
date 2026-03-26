package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

// ─── List item ────────────────────────────────────────────────────────────────

type fileItem struct {
	path    string // full path
	name    string // display name (may include trailing "/" for dirs)
	size    int64
	modTime time.Time
	isDir   bool
}

func (f fileItem) Title() string       { return f.name }
func (f fileItem) FilterValue() string { return f.name }
func (f fileItem) Description() string {
	if f.isDir {
		return "folder"
	}
	return fmt.Sprintf("%s  ·  %s", formatFileSize(f.size), formatFileAge(f.modTime))
}

// ─── Directory scan ───────────────────────────────────────────────────────────

// scanDir returns list items for dir: a ".." entry (unless at fs root),
// then subdirectories, then .txt files.
func scanDir(dir string) []list.Item {
	var items []list.Item

	// Parent entry
	parent := filepath.Dir(dir)
	if parent != dir {
		items = append(items, fileItem{path: parent, name: "../", isDir: true})
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return items
	}

	// Dirs first, then .txt files
	for _, e := range entries {
		if e.IsDir() {
			items = append(items, fileItem{
				path:  filepath.Join(dir, e.Name()),
				name:  e.Name() + "/",
				isDir: true,
			})
		}
	}
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".txt" {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		items = append(items, fileItem{
			path:    filepath.Join(dir, e.Name()),
			name:    e.Name(),
			size:    info.Size(),
			modTime: info.ModTime(),
		})
	}
	return items
}

// ─── Constructor ──────────────────────────────────────────────────────────────

func newFileBrowser(t Theme, items []list.Item, w, h int, dir string) list.Model {
	d := list.NewDefaultDelegate()

	// Selected item styles
	d.Styles.SelectedTitle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.Color(t.Accent)).
		Foreground(lipgloss.Color(t.Accent)).
		Bold(true).
		Padding(0, 0, 0, 1)
	d.Styles.SelectedDesc = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.Color(t.Accent)).
		Foreground(lipgloss.Color(t.Foreground)).
		Padding(0, 0, 0, 1)

	// Normal item styles
	d.Styles.NormalTitle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Foreground)).
		Padding(0, 0, 0, 2)
	d.Styles.NormalDesc = lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Muted)).
		Padding(0, 0, 0, 2)

	// Dimmed (during filter) styles
	d.Styles.DimmedTitle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Muted)).
		Padding(0, 0, 0, 2)
	d.Styles.DimmedDesc = lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Muted)).
		Padding(0, 0, 0, 2)

	// Filter match highlight
	d.Styles.FilterMatch = lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Accent)).
		Bold(true)

	l := list.New(items, d, w, h)
	l.Title = "⬡ " + filepath.Base(dir)
	l.SetShowHelp(false)
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)

	// Title bar styles
	l.Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Accent)).
		Background(lipgloss.Color(t.BarBg)).
		Bold(true).
		Padding(0, 1)
	l.Styles.TitleBar = lipgloss.NewStyle().
		Background(lipgloss.Color(t.BarBg)).
		Padding(0, 0, 1, 0)

	// Filter prompt
	l.Styles.FilterPrompt = lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Accent))
	l.Styles.FilterCursor = lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Accent))

	// Status bar (shows item count)
	l.Styles.StatusBar = lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Muted)).
		Padding(0, 0, 0, 2)

	// Pagination dots
	l.Styles.ActivePaginationDot = lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Accent)).
		SetString("•")
	l.Styles.InactivePaginationDot = lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Muted)).
		SetString("•")

	if len(items) == 0 {
		l.SetShowStatusBar(false)
	}

	return l
}

// ─── Renderer ─────────────────────────────────────────────────────────────────

func renderFileBrowser(m Model) string {
	t := themes[m.themeIdx]

	hint := dimText(t).Render("  ↑↓ Navigate  Enter Open  ⌫ Up  / Filter  Esc Close  ")

	inner := lipgloss.JoinVertical(lipgloss.Left,
		m.fileBrowser.View(),
		hint,
	)

	box := lipgloss.NewStyle().
		Background(lipgloss.Color(t.Background)).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(t.Accent)).
		Render(inner)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		box,
		lipgloss.WithWhitespaceBackground(lipgloss.Color(t.Background)),
	)
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func formatFileSize(size int64) string {
	switch {
	case size < 1024:
		return fmt.Sprintf("%d B", size)
	case size < 1024*1024:
		return fmt.Sprintf("%d KB", size/1024)
	default:
		return fmt.Sprintf("%d MB", size/(1024*1024))
	}
}

func formatFileAge(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	default:
		return t.Format("Jan 2, 2006")
	}
}
