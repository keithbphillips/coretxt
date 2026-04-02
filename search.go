package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ─── Match finding ────────────────────────────────────────────────────────────

// findAllMatches returns rune offsets of all case-insensitive occurrences of
// query in text.
func findAllMatches(text, query string) []int {
	if query == "" {
		return nil
	}
	textRunes := []rune(strings.ToLower(text))
	queryRunes := []rune(strings.ToLower(query))
	queryLen := len(queryRunes)
	if queryLen > len(textRunes) {
		return nil
	}
	var offsets []int
	for i := 0; i <= len(textRunes)-queryLen; {
		match := true
		for j := 0; j < queryLen; j++ {
			if textRunes[i+j] != queryRunes[j] {
				match = false
				break
			}
		}
		if match {
			offsets = append(offsets, i)
			i += queryLen
		} else {
			i++
		}
	}
	return offsets
}

// replaceAllOccurrences replaces all case-insensitive occurrences of query in
// text with replacement. Returns the updated text and the replacement count.
func replaceAllOccurrences(text, query, replacement string) (string, int) {
	if query == "" {
		return text, 0
	}
	textRunes := []rune(text)
	lowerRunes := []rune(strings.ToLower(text))
	queryRunes := []rune(strings.ToLower(query))
	replRunes := []rune(replacement)
	queryLen := len(queryRunes)
	if queryLen > len(textRunes) {
		return text, 0
	}
	var result []rune
	count := 0
	i := 0
	for i <= len(textRunes)-queryLen {
		match := true
		for j := 0; j < queryLen; j++ {
			if lowerRunes[i+j] != queryRunes[j] {
				match = false
				break
			}
		}
		if match {
			result = append(result, replRunes...)
			i += queryLen
			count++
		} else {
			result = append(result, textRunes[i])
			i++
		}
	}
	result = append(result, textRunes[i:]...)
	return string(result), count
}

// ─── Rendering ────────────────────────────────────────────────────────────────

// renderSearchBar renders the search or search+replace bar at the bottom of
// the screen, replacing the normal key hints row.
func renderSearchBar(m Model) string {
	t := themes[m.themeIdx]

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Accent)).
		Background(lipgloss.Color(t.BarBg)).
		Bold(true)
	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Muted)).
		Background(lipgloss.Color(t.BarBg))
	countStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Foreground)).
		Background(lipgloss.Color(t.BarBg))
	barStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(t.BarBg)).
		Width(m.width)

	matchInfo := func() string {
		if m.searchInput.Value() == "" {
			return ""
		}
		if len(m.searchMatches) == 0 {
			return countStyle.Render("  no matches")
		}
		return countStyle.Render(fmt.Sprintf("  %d/%d", m.searchCurrent+1, len(m.searchMatches)))
	}

	if m.mode == modeSearch {
		label := labelStyle.Render(" ⌕ ")
		hints := hintStyle.Render("   ↑/↓:prev/next  Esc:close")
		return barStyle.Render(label + m.searchInput.View() + matchInfo() + hints)
	}

	// modeReplace
	findLabel := labelStyle.Render(" ⌕ Find: ")
	replLabel := labelStyle.Render("  → ")
	var hints string
	if m.searchReplaceFocus == 1 {
		hints = hintStyle.Render("   Enter:replace  ^A:all  Esc:close")
	} else {
		hints = hintStyle.Render("   Tab:switch field  ^A:all  Esc:close")
	}
	return barStyle.Render(findLabel + m.searchInput.View() + replLabel + m.replaceInput.View() + matchInfo() + hints)
}
