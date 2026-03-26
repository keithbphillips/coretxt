package main

import (
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
)

// applyTheme re-styles the textarea model to match the given theme.
// Must be called after ShowLineNumbers and Prompt are set, before SetWidth.
func applyTheme(ta *textarea.Model, t Theme) {
	bg := lipgloss.Color(t.Background)
	fg := lipgloss.Color(t.Foreground)
	accent := lipgloss.Color(t.Accent)
	muted := lipgloss.Color(t.Muted)
	cursorLine := lipgloss.Color(t.CursorLine)

	focused := textarea.Style{
		Base:             lipgloss.NewStyle().Background(bg).Foreground(fg),
		CursorLine:       lipgloss.NewStyle().Background(cursorLine),
		CursorLineNumber: lipgloss.NewStyle().Foreground(accent).Background(cursorLine).Bold(true),
		EndOfBuffer:      lipgloss.NewStyle().Foreground(muted).Background(bg),
		LineNumber:       lipgloss.NewStyle().Foreground(muted).Background(bg),
		Placeholder:      lipgloss.NewStyle().Foreground(muted).Background(bg).Italic(true),
		Prompt:           lipgloss.NewStyle().Foreground(accent).Background(bg),
		Text:             lipgloss.NewStyle().Foreground(fg).Background(bg),
	}

	blurred := textarea.Style{
		Base:             lipgloss.NewStyle().Background(bg).Foreground(muted),
		CursorLine:       lipgloss.NewStyle().Background(bg),
		CursorLineNumber: lipgloss.NewStyle().Foreground(muted).Background(bg),
		EndOfBuffer:      lipgloss.NewStyle().Foreground(muted).Background(bg),
		LineNumber:       lipgloss.NewStyle().Foreground(muted).Background(bg),
		Placeholder:      lipgloss.NewStyle().Foreground(muted).Background(bg).Italic(true),
		Prompt:           lipgloss.NewStyle().Foreground(muted).Background(bg),
		Text:             lipgloss.NewStyle().Foreground(muted).Background(bg),
	}

	ta.FocusedStyle = focused
	ta.BlurredStyle = blurred

	// cursor.View() always applies Reverse(true) on top of this style.
	// Underline(true) gives the underline; Reverse inverts the char colors
	// so the cursor position is doubly visible.
	ta.Cursor.SetMode(cursor.CursorStatic)
	ta.Cursor.Style = lipgloss.NewStyle().Underline(true)
}

// ─── Bar styles ───────────────────────────────────────────────────────────────

func barBase(t Theme) lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(t.BarBg)).
		Foreground(lipgloss.Color(t.Muted))
}

func accentText(t Theme) lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(t.BarBg)).
		Foreground(lipgloss.Color(t.Accent)).
		Bold(true)
}

func normalText(t Theme) lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(t.BarBg)).
		Foreground(lipgloss.Color(t.BarFg))
}

func modifiedText(t Theme) lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(t.BarBg)).
		Foreground(lipgloss.Color(t.Modified)).
		Bold(true)
}

func dimText(t Theme) lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(t.BarBg)).
		Foreground(lipgloss.Color(t.Muted))
}
