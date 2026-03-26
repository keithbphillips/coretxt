package main

import (
	"github.com/charmbracelet/lipgloss"
)

func renderPrompt(m Model) string {
	t := themes[m.themeIdx]

	title := accentText(t).Render("  ⬡ Save Document As")
	hint := dimText(t).Render("\n  Enter to confirm  ·  Esc to cancel")

	inputBox := lipgloss.NewStyle().
		Background(lipgloss.Color(t.Background)).
		Foreground(lipgloss.Color(t.Foreground)).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(lipgloss.Color(t.Accent)).
		Width(46).
		Render(m.nameInput.View())

	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		inputBox,
		hint,
	)

	box := lipgloss.NewStyle().
		Background(lipgloss.Color(t.Background)).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(t.Accent)).
		Padding(1, 2).
		Render(content)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		box,
		lipgloss.WithWhitespaceBackground(lipgloss.Color(t.Background)),
	)
}
