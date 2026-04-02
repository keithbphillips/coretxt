package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/charmbracelet/lipgloss"
)

const appName = "CORETXT"

// ─── Header ───────────────────────────────────────────────────────────────────

func renderHeader(m Model) string {
	t := themes[m.themeIdx]

	// Left: app logo
	left := accentText(t).Render(" ⬡ " + appName + " ")

	// Center: filename + modified dot
	fname := filepath.Base(m.filename)
	if m.filename == "" {
		fname = "untitled"
	}
	center := normalText(t).Render("  " + fname)
	if m.dirty {
		center += modifiedText(t).Render(" ●")
	} else {
		center += normalText(t).Render(" ◎")
	}

	// Right: theme name + word count
	words := wordCount(m.ta.Value())
	right := accentText(t).Render(t.Name)
	right += normalText(t).Render("  " + formatNum(words) + " words ")

	return buildBar(left, center, right, m.width, t)
}

// ─── Status Bar ───────────────────────────────────────────────────────────────

func renderStatusBar(m Model) string {
	t := themes[m.themeIdx]

	content := m.ta.Value()
	words := wordCount(content)

	sep := dimText(t).Render("  ◈  ")

	left := accentText(t).Render(" ◈") +
		normalText(t).Render(fmt.Sprintf(" %s words", formatNum(words))) +
		sep +
		dimText(t).Render(readingTime(words))

	// Right: status message or clock
	var right string
	if m.statusMsg != "" {
		right = accentText(t).Render(m.statusMsg + " ")
	} else {
		right = dimText(t).Render(time.Now().Format("15:04") + " ")
	}

	return buildBar(left, "", right, m.width, t)
}

// ─── Key Hints ────────────────────────────────────────────────────────────────

func renderKeyHints(m Model) string {
	t := themes[m.themeIdx]

	hints := []struct{ key, label string }{
		{"^N", "New"},
		{"^S", "Save"},
		{"^O", "Open"},
		{"^F", "Find"},
		{"^R", "Replace"},
		{"^Q", "Quit"},
		{"F1", "Help"},
		{"F2", "Theme"},
		{"F3", "Save As"},
		{"F7/^Spc", "Spell"},
	}

	var sb strings.Builder
	sb.WriteString("  ")
	for i, h := range hints {
		if i > 0 {
			sb.WriteString(dimText(t).Render("  "))
		}
		sb.WriteString(accentText(t).Render(h.key))
		sb.WriteString(dimText(t).Render(":" + h.label))
	}

	bar := sb.String()
	padded := lipgloss.NewStyle().
		Background(lipgloss.Color(t.BarBg)).
		Width(m.width).
		Render(bar)
	return padded
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

// buildBar assembles a full-width bar with left, optional center, and right sections.
func buildBar(left, center, right string, width int, t Theme) string {
	base := lipgloss.NewStyle().Background(lipgloss.Color(t.BarBg)).Width(width)

	used := lipgloss.Width(left) + lipgloss.Width(center) + lipgloss.Width(right)
	gap := width - used
	if gap < 0 {
		gap = 0
	}

	var row string
	if center == "" {
		// Push right section to the far right
		row = left + strings.Repeat(" ", gap) + right
	} else {
		// left | ~~~ center ~~~ | right
		leftGap := gap / 2
		rightGap := gap - leftGap
		row = left + strings.Repeat(" ", leftGap) + center + strings.Repeat(" ", rightGap) + right
	}

	return base.Render(row)
}

func wordCount(s string) int {
	return len(strings.FieldsFunc(s, func(r rune) bool {
		return unicode.IsSpace(r)
	}))
}

func charCount(s string) int {
	// Count non-whitespace characters
	n := 0
	for _, r := range s {
		if !unicode.IsSpace(r) {
			n++
		}
	}
	return n
}

func readingTime(words int) string {
	if words == 0 {
		return "0 min read"
	}
	mins := words / 238
	if mins < 1 {
		return "< 1 min read"
	}
	return fmt.Sprintf("%d min read", mins)
}

func timeSince(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "saved just now"
	case d < time.Hour:
		return fmt.Sprintf("saved %dm ago", int(d.Minutes()))
	default:
		return fmt.Sprintf("saved %dh ago", int(d.Hours()))
	}
}

func formatNum(n int) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}
	if n < 1_000_000 {
		return fmt.Sprintf("%d,%03d", n/1000, n%1000)
	}
	return fmt.Sprintf("%d,%03d,%03d", n/1_000_000, (n/1000)%1000, n%1000)
}
