package main

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type appMode int

const (
	modeEdit appMode = iota
	modeHelp
	modePromptName
	modeFileBrowser
	modeSpellCheck
)

// Model is the central bubbletea application state.
type Model struct {
	ta          textarea.Model
	nameInput   textinput.Model
	fileBrowser list.Model
	themeIdx    int
	filename    string
	dirty       bool
	lastSaved   time.Time
	mode        appMode
	width       int
	height      int
	statusMsg        string
	quitConfirm      bool
	typewriterMode   bool   // true while typing, false while navigating
	browserDir       string // current directory shown in the file browser
	browserSaveMode  bool   // true when browser was opened for saving
	promptDir        string // directory used by name prompt when saving
	spellWord        string
	spellSuggestions []string
	spellWordLeft    int
	spellWordRight   int
}

func newModel(filename string) Model {
	ta := textarea.New()
	ta.Placeholder = "Begin your story here..."
	ta.ShowLineNumbers = false
	ta.Prompt = ""
	ta.CharLimit = 0
	ta.MaxHeight = 0

	ni := textinput.New()
	ni.Placeholder = "my-novel.txt"
	ni.CharLimit = 256
	ni.Width = 44
	ni.Prompt = "  Filename: "

	p := loadPrefs()

	browserDir := p.LastDir
	if browserDir == "" {
		browserDir = docsDir()
	}

	m := Model{
		ta:         ta,
		nameInput:  ni,
		themeIdx:   p.ThemeIdx,
		filename:   filename,
		browserDir: browserDir,
	}

	applyTheme(&m.ta, themes[m.themeIdx])
	applyThemeToInput(&m.nameInput, themes[m.themeIdx])

	// Add ctrl+arrow word movement on top of the default alt+arrow bindings.
	m.ta.KeyMap.WordForward.SetKeys("alt+right", "alt+f", "ctrl+right")
	m.ta.KeyMap.WordBackward.SetKeys("alt+left", "alt+b", "ctrl+left")

	m.ta.Focus()

	contentLoaded := false
	if filename != "" {
		filename = resolvePath(filename)
		if content, err := loadFile(filename); err == nil && content != "" {
			m.ta.SetValue(content)
			m.lastSaved = time.Now()
			contentLoaded = true
		}
	}

	if !contentLoaded {
		// Blank document: pre-indent the first paragraph, dirty stays false
		m.ta.InsertString("    ")
	}

	return m
}

// ─── Init ─────────────────────────────────────────────────────────────────────

func (m Model) Init() tea.Cmd {
	return doTick()
}

// ─── Update ───────────────────────────────────────────────────────────────────

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ta.SetWidth(msg.Width)
		m.syncTaHeight()
		return m, nil

	case tickMsg:
		cmds = append(cmds, doTick())
		if m.dirty && m.filename != "" && m.filename != "untitled.txt" {
			if err := m.performSave(); err == nil {
				m.statusMsg = "Auto-saved"
				cmds = append(cmds, clearStatus(2*time.Second))
			}
		}
		return m, tea.Batch(cmds...)

	case statusClearMsg:
		m.statusMsg = ""
		return m, nil

	case spellResultMsg:
		if msg.err != "" {
			m.statusMsg = msg.err
			return m, clearStatus(3 * time.Second)
		}
		if len(msg.suggestions) == 0 {
			m.statusMsg = "✓ \"" + msg.word + "\" is spelled correctly"
			return m, clearStatus(3 * time.Second)
		}
		m.spellWord = msg.word
		m.spellSuggestions = msg.suggestions
		m.spellWordLeft = msg.wordLeft
		m.spellWordRight = msg.wordRight
		m.mode = modeSpellCheck
		return m, nil

	case tea.KeyMsg:

		// ── Filename prompt mode ──────────────────────────────────────────
		if m.mode == modePromptName {
			switch msg.String() {
			case "enter":
				name := strings.TrimSpace(m.nameInput.Value())
				if name != "" {
					if filepath.Ext(name) == "" {
						name += ".txt"
					}
					if m.promptDir != "" && !filepath.IsAbs(name) {
						m.filename = filepath.Join(m.promptDir, name)
					} else {
						m.filename = resolvePath(name)
					}
				}
				m.promptDir = ""
				if m.filename != "" && m.filename != "untitled.txt" {
					if err := m.performSave(); err != nil {
						m.statusMsg = "Save error: " + err.Error()
					} else {
						m.statusMsg = "Saved as \"" + filepath.Base(m.filename) + "\" ✓"
					}
				}
				m.mode = modeEdit
				m.nameInput.SetValue("")
				return m, clearStatus(3 * time.Second)

			case "esc":
				m.mode = modeEdit
				m.nameInput.SetValue("")
				return m, nil

			default:
				var tiCmd tea.Cmd
				m.nameInput, tiCmd = m.nameInput.Update(msg)
				return m, tiCmd
			}
		}

		// ── File browser mode ─────────────────────────────────────────────
		if m.mode == modeFileBrowser {
			filtering := m.fileBrowser.FilterState() == list.Filtering
			switch msg.String() {
			case "esc":
				if filtering {
					var lCmd tea.Cmd
					m.fileBrowser, lCmd = m.fileBrowser.Update(msg)
					return m, lCmd
				}
				m.mode = modeEdit
				m.browserSaveMode = false
				return m, nil

			case "backspace":
				if filtering {
					var lCmd tea.Cmd
					m.fileBrowser, lCmd = m.fileBrowser.Update(msg)
					return m, lCmd
				}
				parent := filepath.Dir(m.browserDir)
				if parent != m.browserDir {
					m.browserDir = parent
					m.rebuildFileBrowser()
				}
				return m, nil

			case "n":
				if m.browserSaveMode && !filtering {
					m.promptDir = m.browserDir
					prefill := ""
					if m.filename != "" && filepath.Dir(m.filename) == m.browserDir {
						prefill = filepath.Base(m.filename)
					}
					return m, m.openNamePrompt(prefill)
				}
				var lCmd tea.Cmd
				m.fileBrowser, lCmd = m.fileBrowser.Update(msg)
				return m, lCmd

			case "enter":
				if filtering {
					var lCmd tea.Cmd
					m.fileBrowser, lCmd = m.fileBrowser.Update(msg)
					return m, lCmd
				}
				if item, ok := m.fileBrowser.SelectedItem().(fileItem); ok {
					if item.isDir {
						m.browserDir = item.path
						m.rebuildFileBrowser()
						return m, nil
					}
					if m.browserSaveMode {
						m.filename = item.path
						m.browserSaveMode = false
						m.mode = modeEdit
						if err := m.performSave(); err != nil {
							m.statusMsg = "Save error: " + err.Error()
						} else {
							m.statusMsg = "Saved as \"" + filepath.Base(item.path) + "\" ✓"
							p := loadPrefs()
							p.LastDir = m.browserDir
							savePrefs(p)
						}
						return m, clearStatus(3 * time.Second)
					}
					content, err := loadFile(item.path)
					if err == nil {
						m.ta.SetValue(content)
						m.filename = item.path
						m.dirty = false
						m.lastSaved = time.Now()
						m.statusMsg = "Opened \"" + filepath.Base(item.path) + "\""
						p := loadPrefs()
						p.LastDir = m.browserDir
						savePrefs(p)
					} else {
						m.statusMsg = "Error: " + err.Error()
					}
					m.mode = modeEdit
					return m, clearStatus(3 * time.Second)
				}
				return m, nil

			default:
				var lCmd tea.Cmd
				m.fileBrowser, lCmd = m.fileBrowser.Update(msg)
				return m, lCmd
			}
		}

		// ── Spell-check overlay ───────────────────────────────────────────
		if m.mode == modeSpellCheck {
			switch msg.String() {
			case "esc":
				m.mode = modeEdit
			case "1", "2", "3", "4", "5", "6", "7", "8", "9":
				idx := int(msg.String()[0]-'1')
				if idx < len(m.spellSuggestions) {
					suggestion := m.spellSuggestions[idx]
					for i := 0; i < m.spellWordLeft; i++ {
						m.ta, _ = m.ta.Update(tea.KeyMsg{Type: tea.KeyBackspace})
					}
					for i := 0; i < m.spellWordRight; i++ {
						m.ta, _ = m.ta.Update(tea.KeyMsg{Type: tea.KeyDelete})
					}
					m.ta.InsertString(suggestion)
					var taCmd tea.Cmd
					m.ta, taCmd = m.ta.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{}})
					m.dirty = true
					m.mode = modeEdit
					return m, taCmd
				}
			}
			return m, nil
		}

		// ── Escape clears overlay states ──────────────────────────────────
		if msg.Type == tea.KeyEsc {
			if m.mode == modeHelp {
				m.mode = modeEdit
				return m, nil
			}
			if m.quitConfirm {
				m.quitConfirm = false
				m.statusMsg = ""
				return m, nil
			}
		}

		// ── Help overlay swallows all keys except F1 ──────────────────────
		if m.mode == modeHelp {
			if msg.Type == tea.KeyF1 {
				m.mode = modeEdit
			}
			return m, nil
		}

		// ── Normal editing keys ───────────────────────────────────────────
		switch msg.String() {

		case "ctrl+@": // Ctrl+Space — spell check word at cursor
			word, left, right := wordAtCursor(m)
			if word == "" {
				m.statusMsg = "No word at cursor"
				return m, clearStatus(2 * time.Second)
			}
			m.statusMsg = "Checking \"" + word + "\"…"
			return m, tea.Batch(clearStatus(5*time.Second), checkSpelling(word, left, right))

		case "ctrl+s":
			if m.filename == "" || m.filename == "untitled.txt" {
				m.openSaveBrowser()
				return m, nil
			}
			if err := m.performSave(); err != nil {
				m.statusMsg = "Save error: " + err.Error()
			} else {
				m.statusMsg = "Saved ✓"
			}
			return m, clearStatus(2 * time.Second)

		case "ctrl+o": // Open file browser
			m.browserSaveMode = false
			m.rebuildFileBrowser()
			m.mode = modeFileBrowser
			return m, nil

		case "f3": // Save As
			m.openSaveBrowser()
			return m, nil

		case "ctrl+q":
			if m.dirty {
				if m.quitConfirm {
					return m, tea.Quit
				}
				m.quitConfirm = true
				m.statusMsg = "Unsaved changes — Ctrl+Q again to quit"
				return m, clearStatus(4 * time.Second)
			}
			return m, tea.Quit

		case "ctrl+c":
			if m.dirty {
				m.statusMsg = "Unsaved changes — use Ctrl+Q to quit"
				return m, clearStatus(3 * time.Second)
			}
			return m, tea.Quit

		case "f1":
			if m.mode == modeHelp {
				m.mode = modeEdit
			} else {
				m.mode = modeHelp
			}
			return m, nil

		case "f2":
			m.themeIdx = (m.themeIdx + 1) % len(themes)
			savePrefs(prefs{ThemeIdx: m.themeIdx})
			applyTheme(&m.ta, themes[m.themeIdx])
			applyThemeToInput(&m.nameInput, themes[m.themeIdx])
			focusCmd := m.ta.Focus() // resets internal style pointer to current copy
			if m.width > 0 {
				m.ta.SetWidth(m.width)
				m.ta.SetHeight(taHeight(m.height))
			}
			return m, focusCmd

		case "enter":
			m.typewriterMode = true
			m.quitConfirm = false
			m.dirty = true
			if isEndOfParagraph(m.ta) {
				m.ta.InsertString("\n\n    ")
			} else {
				m.ta.InsertString("\n    ")
			}
			// InsertString doesn't call repositionView internally, so the cursor
			// can land below the viewport and disappear. A no-op Update fixes it.
			var taCmd tea.Cmd
			m.ta, taCmd = m.ta.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{}})
			m.syncTaHeight()
			cmds = append(cmds, taCmd)
			return m, tea.Batch(cmds...)

		default:
			m.quitConfirm = false
			if !isNavigationKey(msg) {
				m.dirty = true
				m.typewriterMode = true
			} else {
				m.typewriterMode = false
			}
			var taCmd tea.Cmd
			m.ta, taCmd = m.ta.Update(msg)
			m.syncTaHeight()
			cmds = append(cmds, taCmd)
			return m, tea.Batch(cmds...)
		}
	}

	// All other messages (mouse, focus, blink) go to textarea
	var taCmd tea.Cmd
	m.ta, taCmd = m.ta.Update(msg)
	m.syncTaHeight()
	cmds = append(cmds, taCmd)
	return m, tea.Batch(cmds...)
}

// ─── View ─────────────────────────────────────────────────────────────────────

func (m Model) View() string {
	if m.width == 0 {
		return "Initializing..."
	}

	switch m.mode {
	case modeHelp:
		return renderHelp(m)
	case modePromptName:
		return renderPrompt(m)
	case modeFileBrowser:
		return renderFileBrowser(m)
	case modeSpellCheck:
		return renderSpellCheck(m)
	}

	t := themes[m.themeIdx]
	spacer := lipgloss.NewStyle().
		Background(lipgloss.Color(t.Background)).
		Width(m.width).
		Render("")

	// Fill any remaining space below the textarea with styled blank lines so
	// the background colour is consistent across the whole screen.
	th := m.ta.Height()
	totalDoc := m.height - 4 // header + spacer + statusbar + keyhints
	padLines := totalDoc - th
	if padLines < 0 {
		padLines = 0
	}
	blankLine := lipgloss.NewStyle().
		Background(lipgloss.Color(t.Background)).
		Width(m.width).
		Render("")
	padding := strings.Repeat("\n"+blankLine, padLines)

	return lipgloss.JoinVertical(lipgloss.Left,
		renderHeader(m),
		spacer,
		m.ta.View()+padding,
		renderStatusBar(m),
		renderKeyHints(m),
	)
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

// rebuildFileBrowser recreates the file browser list for m.browserDir.
func (m *Model) rebuildFileBrowser() {
	bw := m.width - 8
	if bw > 64 {
		bw = 64
	}
	bh := m.height - 10
	if bh > 20 {
		bh = 20
	}
	items := scanDir(m.browserDir)
	m.fileBrowser = newFileBrowser(themes[m.themeIdx], items, bw, bh, m.browserDir)
	if m.browserSaveMode {
		m.fileBrowser.Title = "⬡ Save to: " + filepath.Base(m.browserDir)
	}
}

// openSaveBrowser switches to the file browser in save mode.
func (m *Model) openSaveBrowser() {
	m.browserSaveMode = true
	m.rebuildFileBrowser()
	m.mode = modeFileBrowser
}

// openNamePrompt switches to the filename prompt, pre-filling with prefill.
// Returns the Focus command for the text input.
func (m *Model) openNamePrompt(prefill string) tea.Cmd {
	m.mode = modePromptName
	m.nameInput.SetValue(prefill)
	// Position cursor at end of any pre-filled text
	m.nameInput.CursorEnd()
	return m.nameInput.Focus()
}

func (m *Model) performSave() error {
	if err := saveFile(m.filename, m.ta.Value()); err != nil {
		return err
	}
	m.dirty = false
	m.lastSaved = time.Now()
	m.quitConfirm = false
	return nil
}

func isNavigationKey(msg tea.KeyMsg) bool {
	switch msg.Type {
	case tea.KeyUp, tea.KeyDown, tea.KeyLeft, tea.KeyRight,
		tea.KeyPgUp, tea.KeyPgDown, tea.KeyHome, tea.KeyEnd,
		tea.KeyF1, tea.KeyF2, tea.KeyF3, tea.KeyF4,
		tea.KeyF5, tea.KeyF6, tea.KeyF7, tea.KeyF8:
		return true
	}
	switch msg.String() {
	case "ctrl+left", "ctrl+right", "ctrl+home", "ctrl+end",
		"alt+left", "alt+right", "alt+up", "alt+down":
		return true
	}
	return false
}

func isEndOfParagraph(ta textarea.Model) bool {
	lines := strings.Split(ta.Value(), "\n")
	row := ta.Line()
	if row >= len(lines) {
		return false
	}
	return strings.TrimSpace(lines[row]) != ""
}

// applyThemeToInput styles the textinput widget to match the current theme.
func applyThemeToInput(ni *textinput.Model, t Theme) {
	ni.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(t.Foreground))
	ni.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(t.Accent)).Bold(true)
	ni.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(t.Muted))
	ni.Cursor.SetMode(cursor.CursorStatic)
	ni.Cursor.Style = lipgloss.NewStyle().Underline(true)
}

// taHeight returns the textarea height for a given terminal height.
// The textarea occupies the upper half of the document area so the cursor
// (which the textarea always keeps at the bottom of its viewport) sits at
// mid-screen — typewriter scrolling.
func taHeight(termHeight int) int {
	available := termHeight - 4 // subtract header, spacer, statusbar, keyhints
	h := available / 2
	if h < 1 {
		h = 1
	}
	return h
}

// syncTaHeight switches between two modes:
//   - Typewriter mode (last key was a typing key): textarea occupies the top
//     half of the screen so the cursor sits at mid-screen while writing.
//   - Reading mode (last key was a navigation key): textarea fills the full
//     available height so scrolling uses the whole screen.
func (m *Model) syncTaHeight() {
	if m.height == 0 {
		return
	}
	available := m.height - 4
	if available < 1 {
		available = 1
	}
	if !m.typewriterMode {
		m.ta.SetHeight(available)
		return
	}
	half := available / 2
	if half < 1 {
		half = 1
	}
	m.ta.SetHeight(half)
}
