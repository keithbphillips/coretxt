package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/muesli/termenv"
)

// TestDimBelowCursor checks that in typewriter mode the document text below the
// cursor's line is rendered dimmed (muted foreground) while the cursor line and
// the text above it stay bright. It runs across colour profiles because the
// cursor-row detection must not depend on truecolor backgrounds.
func TestDimBelowCursor(t *testing.T) {
	for _, prof := range []termenv.Profile{termenv.TrueColor, termenv.ANSI256, termenv.ANSI} {
		lipgloss.SetColorProfile(prof)

		m := newModel("")
		m.width = 24
		m.height = 18 // totalDoc = 14, half = 7
		m.ta.SetWidth(m.width)
		var lines []string
		for i := 1; i <= 12; i++ {
			lines = append(lines, fmt.Sprintf("line-%02d", i))
		}
		m.ta.SetValue(strings.Join(lines, "\n"))
		m.typewriterMode = true
		m.syncTaHeight()

		// Put the cursor near the top so there is bright text below it that is
		// still above the fold — the case that originally rendered wrong.
		for i := 0; i < 30; i++ {
			m.ta.CursorUp()
		}
		for i := 0; i < 2; i++ {
			m.ta.CursorDown()
		}

		out := strings.Split(m.View(), "\n")
		cursor := cursorDisplayRow(out)
		if cursor < 0 {
			t.Fatalf("profile=%v: cursor row not found in view", prof)
		}

		muted := fgColorToken(themes[m.themeIdx].Muted)

		// The row above the cursor is real text and must not be dimmed.
		if strings.Contains(out[cursor-1], muted) {
			t.Errorf("profile=%v: row above cursor dimmed: %q", prof, ansi.Strip(out[cursor-1]))
		}
		// Rows below the cursor (still above the fold, then past it) must be dimmed.
		for d := 1; d <= 3; d++ {
			row := out[cursor+d]
			if !strings.Contains(row, muted) {
				t.Errorf("profile=%v: row %d below cursor not dimmed: %q", prof, d, ansi.Strip(row))
			}
			if rowHasCursorAttr(row) {
				t.Errorf("profile=%v: dimmed row carries cursor attribute: %q", prof, ansi.Strip(row))
			}
		}

		// The shared viewport height must be restored after the full-height render.
		if got, want := m.ta.Height(), (m.height-4)/2; got != want {
			t.Errorf("profile=%v: textarea height not restored: got %d want %d", prof, got, want)
		}
	}
}

// fgColorToken returns the SGR foreground parameter run lipgloss emits for a
// colour (e.g. "38;2;160;120;48"), matchable as a substring of a rendered line.
func fgColorToken(color string) string {
	probe := lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(" ")
	start := strings.IndexByte(probe, '[')
	if start < 0 {
		return ""
	}
	rest := probe[start+1:]
	end := strings.IndexByte(rest, 'm')
	if end < 0 {
		return ""
	}
	return rest[:end]
}
