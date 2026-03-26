package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	var filename string
	if len(os.Args) > 1 {
		filename = os.Args[1]
	}

	m := newModel(filename)
	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "coretxt:", err)
		os.Exit(1)
	}
}
