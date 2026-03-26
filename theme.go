package main

// Theme defines a complete neon color palette for the UI.
type Theme struct {
	Name       string
	Background string // editor background
	Foreground string // document text
	BarFg      string // bar normal text (slightly dimmer than Foreground)
	Accent     string // neon highlight color (borders, active elements)
	Muted      string // line numbers, dim text
	CursorLine string // current line background
	BarBg      string // header / status bar background
	Modified   string // unsaved changes indicator
	Selection  string // selected text background
}

var themes = []Theme{
	{
		Name:       "CYBERPUNK",
		Background: "#0a0a12",
		Foreground: "#c8c8e8",
		BarFg:      "#8080a8",
		Accent:     "#00ffff",
		Muted:      "#6868a8",
		CursorLine: "#12122a",
		BarBg:      "#10102a",
		Modified:   "#ff2266",
		Selection:  "#1a1a40",
	},
	{
		Name:       "SYNTHWAVE",
		Background: "#150025",
		Foreground: "#f0d0ff",
		BarFg:      "#b080cc",
		Accent:     "#ff71ce",
		Muted:      "#a050c0",
		CursorLine: "#200038",
		BarBg:      "#1a0030",
		Modified:   "#ff0099",
		Selection:  "#2d0050",
	},
	{
		Name:       "MATRIX",
		Background: "#000d00",
		Foreground: "#00cc44",
		BarFg:      "#008830",
		Accent:     "#00ff41",
		Muted:      "#207840",
		CursorLine: "#001800",
		BarBg:      "#001200",
		Modified:   "#ffff00",
		Selection:  "#002200",
	},
	{
		Name:       "NEON AMBER",
		Background: "#0c0800",
		Foreground: "#ffcc66",
		BarFg:      "#c09040",
		Accent:     "#ffdd00",
		Muted:      "#a07830",
		CursorLine: "#1a1200",
		BarBg:      "#120a00",
		Modified:   "#ff3300",
		Selection:  "#2a1e00",
	},
}
