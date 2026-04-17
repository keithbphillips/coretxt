package main

import (
	"reflect"
	"unsafe"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// textareaViewportOffset reads the vertical scroll offset from the textarea's
// unexported viewport field. This lets us convert a click's terminal row into
// the correct document visual row.
//
// Rationale: bubbles/textarea v1.0.0 has no public mouse API. The viewport
// field is a *viewport.Model whose YOffset is exported. We access it via
// reflect.Value.Pointer() (which does not enforce exportedness) and unsafe,
// both of which are standard Go mechanisms for this pattern.
func textareaViewportOffset(ta textarea.Model) int {
	rv := reflect.ValueOf(ta)
	vpField := rv.FieldByName("viewport")
	if !vpField.IsValid() {
		return 0
	}
	ptr := vpField.Pointer()
	if ptr == 0 {
		return 0
	}
	vp := (*viewport.Model)(unsafe.Pointer(ptr)) //nolint:govet
	return vp.YOffset
}

// handleMouseClick moves the textarea cursor to the character at terminal
// position (x, y). Only valid when in modeEdit; silently ignored otherwise.
//
// Strategy:
//   - Convert click (x, y) to a document visual row using the viewport offset.
//   - Navigate with Ctrl+Home → Down×docRow → Right×x.
//     Down is O(visual rows) rather than O(characters), which is much faster
//     for typical prose documents.
func (m *Model) handleMouseClick(x, y int) {
	// Layout: row 0 = header, row 1 = spacer, rows 2..2+H-1 = textarea.
	const taStartRow = 2
	clickRow := y - taStartRow
	if clickRow < 0 || clickRow >= m.ta.Height() {
		return
	}

	viewportOffset := textareaViewportOffset(m.ta)
	targetDocRow := viewportOffset + clickRow

	// Ctrl+Home resets the cursor to row 0 and clears lastCharOffset to 0,
	// so subsequent Down keys navigate to the beginning of each visual row.
	m.ta, _ = m.ta.Update(tea.KeyMsg{Type: tea.KeyCtrlHome})
	for i := 0; i < targetDocRow; i++ {
		m.ta, _ = m.ta.Update(tea.KeyMsg{Type: tea.KeyDown})
	}
	// Right×x positions within the visual row. For typical ASCII prose each
	// character is one cell wide, so this lands exactly on the clicked cell.
	for i := 0; i < x; i++ {
		m.ta, _ = m.ta.Update(tea.KeyMsg{Type: tea.KeyRight})
	}

	m.typewriterMode = false
	m.syncTaHeight()
}
