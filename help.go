package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func renderHelp(m Model) string {
	t := themes[m.themeIdx]

	accent := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Accent)).Bold(true)
	normal := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Foreground))
	dim := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Muted))
	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Accent)).
		Bold(true).
		Underline(true)

	type entry struct{ key, desc string }
	sections := []struct {
		heading string
		entries []entry
	}{
		{
			"FILE",
			[]entry{
				{"Ctrl+N", "New document"},
				{"Ctrl+S", "Save file"},
				{"Ctrl+Q", "Quit (confirms if unsaved)"},
				{"Ctrl+C", "Force quit"},
			},
		},
		{
			"NAVIGATION",
			[]entry{
				{"Ctrl+A / Ctrl+E", "Start / end of paragraph"},
				{"Ctrl+← / →", "Jump word"},
				{"Ctrl+Home/End", "Beginning / end of document"},
				{"PgUp / PgDn", "Scroll page"},
			},
		},
		{
			"EDITING",
			[]entry{
				{"Enter", "New line"},
				{"Backspace", "Delete back"},
				{"Ctrl+W", "Delete word back"},
				{"Ctrl+K", "Delete to end of line"},
			},
		},
		{
			"INTERFACE",
			[]entry{
				{"F1", "Toggle this help"},
				{"F2", "Cycle color theme"},
			},
		},
	}

	var lines []string

	banner := accent.Render("  ⬡ " + appName + "  ")
	lines = append(lines, banner)
	lines = append(lines, dim.Render("  A neon novel writing environment"))
	lines = append(lines, "")

	for _, sec := range sections {
		lines = append(lines, title.Render("  "+sec.heading))
		for _, e := range sec.entries {
			key := accent.Render("  " + padRight(e.key, 20))
			desc := normal.Render(e.desc)
			lines = append(lines, key+desc)
		}
		lines = append(lines, "")
	}

	lines = append(lines, dim.Render("  Press F1 or Esc to close"))

	boxContent := strings.Join(lines, "\n")

	boxWidth := 52
	boxStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(t.Background)).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(t.Accent)).
		Width(boxWidth).
		Padding(1, 0)

	box := boxStyle.Render(boxContent)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		box,
		lipgloss.WithWhitespaceBackground(lipgloss.Color(t.Background)),
	)
}

func padRight(s string, n int) string {
	if len(s) >= n {
		return s
	}
	return s + strings.Repeat(" ", n-len(s))
}
