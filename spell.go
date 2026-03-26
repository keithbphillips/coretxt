package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
	"unicode"

	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
)

// ─── Messages ─────────────────────────────────────────────────────────────────

type spellResultMsg struct {
	word        string
	suggestions []string // empty = correctly spelled; nil = aspell unavailable
	wordLeft    int      // runes from word-start to cursor
	wordRight   int      // runes from cursor to word-end
	err         string   // non-empty = error string to show
}

// ─── Word extraction ──────────────────────────────────────────────────────────

func isWordRune(r rune) bool {
	return unicode.IsLetter(r) || r == '\''
}

// wordAtCursor returns the word under the cursor plus how many runes lie to the
// left and right of the cursor within that word.
func wordAtCursor(m Model) (word string, left, right int) {
	lines := strings.Split(m.ta.Value(), "\n")
	if m.ta.Line() >= len(lines) {
		return "", 0, 0
	}
	runes := []rune(lines[m.ta.Line()])
	info := m.ta.LineInfo()
	col := info.StartColumn + info.ColumnOffset
	if col > len(runes) {
		col = len(runes)
	}

	// Expand left to find word start.
	start := col
	for start > 0 && isWordRune(runes[start-1]) {
		start--
	}
	// Expand right to find word end.
	end := col
	for end < len(runes) && isWordRune(runes[end]) {
		end++
	}
	if start == end {
		return "", 0, 0
	}
	return string(runes[start:end]), col - start, end - col
}

// ─── Aspell integration ───────────────────────────────────────────────────────

// checkSpelling returns a tea.Cmd that shells to aspell and sends back a
// spellResultMsg. Runs off the main goroutine so the UI never blocks.
func checkSpelling(word string, left, right int) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("aspell", "-a")
		cmd.Stdin = strings.NewReader(word + "\n")
		out, err := cmd.Output()
		if err != nil {
			return spellResultMsg{
				word: word,
				err:  "aspell not available",
			}
		}

		scanner := bufio.NewScanner(strings.NewReader(string(out)))
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "*") {
				// Correctly spelled.
				return spellResultMsg{word: word, wordLeft: left, wordRight: right}
			}
			if strings.HasPrefix(line, "&") {
				// Format: & word count offset: sug1, sug2, …
				colon := strings.Index(line, ":")
				if colon < 0 {
					break
				}
				parts := strings.Split(strings.TrimSpace(line[colon+1:]), ", ")
				if len(parts) > 9 {
					parts = parts[:9]
				}
				return spellResultMsg{
					word:        word,
					suggestions: parts,
					wordLeft:    left,
					wordRight:   right,
				}
			}
			if strings.HasPrefix(line, "#") {
				// No suggestions.
				return spellResultMsg{word: word, wordLeft: left, wordRight: right}
			}
		}
		return spellResultMsg{word: word, wordLeft: left, wordRight: right}
	}
}

// ─── Suggestion overlay ───────────────────────────────────────────────────────

func renderSpellCheck(m Model) string {
	t := themes[m.themeIdx]

	title := accentText(t).Render(fmt.Sprintf("  ⬡ Spelling: \"%s\"", m.spellWord))

	var body strings.Builder
	if len(m.spellSuggestions) == 0 {
		body.WriteString(normalText(t).Render("  ✓ Correctly spelled"))
	} else {
		for i, s := range m.spellSuggestions {
			num := accentText(t).Render(fmt.Sprintf("  %d", i+1))
			word := dimText(t).Render("  " + s)
			body.WriteString(num + word)
			if i < len(m.spellSuggestions)-1 {
				body.WriteByte('\n')
			}
		}
	}

	var hint string
	if len(m.spellSuggestions) > 0 {
		hint = dimText(t).Render("\n  1-9 to replace  ·  Esc to cancel")
	} else {
		hint = dimText(t).Render("\n  Esc to close")
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		body.String(),
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
