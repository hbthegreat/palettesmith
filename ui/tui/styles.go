package tui

import "github.com/charmbracelet/lipgloss"

// Color constants for TUI styling
const (
	// Primary colors
	ColorAccent     = "86"  // Bright cyan - used for highlighted/selected text
	ColorDimmed     = "240" // Dark gray - used for unselected/dimmed text  
	ColorHelp       = "241" // Medium gray - used for help text
	ColorBackground = ""    // Default terminal background
	
	// Additional UI colors
	ColorBorder     = "#666666" // Border color for panels
	ColorTextNormal = "#b0b0b0" // Normal text color
	ColorTextFocus  = "#e6e6e6" // Focused text color
	ColorTextTag    = "#7d7d7d" // Tag text color (same as help)
	ColorTextHelp   = "#777777" // Help text color (slightly lighter)
	ColorError      = "#ff6b6b" // Error text color
	ColorStatus     = "#8ece6a" // Status/success color
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
		
	// Form styles
	FormRowNormalStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorTextNormal))
		
	FormRowFocusStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorTextFocus)).
		Bold(true)
		
	FormTagStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorTextTag))
		
	FormHelpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorTextHelp))
		
	FormErrorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorError))
		
	// App layout styles
	LeftPaneStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderRight(true).
		BorderForeground(lipgloss.Color(ColorBorder))
		
	RightPaneStyle = lipgloss.NewStyle().
		Padding(0, 2)
		
	TabActiveStyle = lipgloss.NewStyle().
		Bold(true)
		
	TabDimStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorTextTag))
		
	StatusStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorStatus))
)