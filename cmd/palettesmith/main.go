package main

import (
	"fmt"
	"os"
	"palettesmith/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(tui.New(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
