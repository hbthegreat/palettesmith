package tui

import "github.com/charmbracelet/lipgloss"

// Color constants for TUI styling
const (
	// Primary colors
	ColorAccent     = "86"  // Bright cyan - used for highlighted/selected text
	ColorDimmed     = "240" // Dark gray - used for unselected/dimmed text  
	ColorHelp       = "241" // Medium gray - used for help text
	ColorBackground = ""    // Default terminal background
)

// Style presets for common UI elements
var (
	// TitleStyle for main headings
	TitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(ColorAccent))

	// SelectedStyle for highlighted/selected items
	SelectedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorAccent)).
		Bold(true)

	// UnselectedStyle for normal/dimmed items
	UnselectedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorDimmed))

	// HelpStyle for help text
	HelpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorHelp))
)