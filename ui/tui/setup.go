package tui

import (
	"fmt"
	
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type PresetChosenMsg struct {
	Preset string
}

type SetupModel struct {
	selectedPreset string
	showConfirm    bool  // true when showing confirmation screen
	confirmChoice  string // "yes" or "no" - which option is selected on confirm screen
	confirmed      bool  // true when user confirms and we're ready to proceed
}

func NewSetupModel() SetupModel {
	return SetupModel{
		selectedPreset: "generic",
		showConfirm:    false,
		confirmChoice:  "yes", // default to "yes"
		confirmed:      false,
	}
}

func (m SetupModel) Init() tea.Cmd {
	return nil
}

func (m SetupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case PresetChosenMsg:
		// We sent this message ourselves - now quit so main can handle the transition
		return m, tea.Quit
		
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}

		if m.showConfirm {
			// On confirmation screen
			switch msg.String() {
			case "y", "Y":
				m.confirmChoice = "yes"
			case "n", "N":
				m.confirmChoice = "no"
			case "up", "down":
				// Toggle between yes/no
				if m.confirmChoice == "yes" {
					m.confirmChoice = "no"
				} else {
					m.confirmChoice = "yes"
				}
			case "enter":
				if m.confirmChoice == "yes" {
					// User confirms - proceed to main app
					m.confirmed = true
					return m, func() tea.Msg {
						return PresetChosenMsg{Preset: m.selectedPreset}
					}
				} else {
					// User wants to go back to selection
					m.showConfirm = false
				}
			case "esc":
				// Always go back on escape
				m.showConfirm = false
			}
		} else {
			// On selection screen
			switch msg.String() {
			case "g", "G":
				m.selectedPreset = "generic"
				m.showConfirm = true
			case "o", "O":
				m.selectedPreset = "omarchy"
				m.showConfirm = true
			case "up", "down":
				if m.selectedPreset == "generic" {
					m.selectedPreset = "omarchy"
				} else {
					m.selectedPreset = "generic"
				}
			case "enter":
				// Enter shows confirmation for current selection
				m.showConfirm = true
			}
		}
	}
	return m, nil
}

func (m SetupModel) View() string {
	title := TitleStyle.Render("Welcome to Palettesmith")

	if m.showConfirm {
		// Show confirmation screen
		presetName := m.selectedPreset
		if presetName == "omarchy" {
			presetName = "Omarchy (~/.config/omarchy/)"
		} else {
			presetName = "Generic (~/.config/palettesmith/)"
		}

		confirmMsg := fmt.Sprintf("Configure Palettesmith with %s preset?", presetName)
		
		var yesOption, noOption string
		
		if m.confirmChoice == "yes" {
			yesOption = SelectedStyle.Render("▸ [Y] Yes, continue to theming")
			noOption = UnselectedStyle.Render("  [N] No, change setup")
		} else {
			yesOption = UnselectedStyle.Render("  [Y] Yes, continue to theming")
			noOption = SelectedStyle.Render("▸ [N] No, change setup")
		}
		
		help := HelpStyle.Render("Y/N to select, ↑/↓ to navigate, Enter to confirm, Esc to go back, Q to quit")

		return lipgloss.JoinVertical(
			lipgloss.Left,
			"",
			" "+title,
			"",
			" "+confirmMsg,
			"",
			" "+yesOption,
			" "+noOption,
			"",
			" "+help,
			"",
		)
	}

	// Show selection screen
	genericStyle := UnselectedStyle
	omarchyStyle := UnselectedStyle

	if m.selectedPreset == "generic" {
		genericStyle = SelectedStyle
	} else {
		omarchyStyle = SelectedStyle
	}

	genericOption := genericStyle.Render("▸ [G] Generic - Use ~/.config/palettesmith/")
	omarchyOption := omarchyStyle.Render("  [O] Omarchy - Use ~/.config/omarchy/")

	if m.selectedPreset == "omarchy" {
		omarchyOption = omarchyStyle.Render("▸ [O] Omarchy - Use ~/.config/omarchy/")
		genericOption = genericStyle.Render("  [G] Generic - Use ~/.config/palettesmith/")
	}

	help := HelpStyle.Render("Press O/G to select, ↑/↓ to navigate, Enter to confirm, Q to quit")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		" "+title,
		"",
		" Choose your configuration:",
		"",
		" "+genericOption,
		" "+omarchyOption,
		"",
		" "+help,
		"",
	)
}

// IsConfirmed returns whether the user has confirmed their preset choice
func (m SetupModel) IsConfirmed() bool {
	return m.confirmed
}

// GetSelectedPreset returns the currently selected preset
func (m SetupModel) GetSelectedPreset() string {
	return m.selectedPreset
}
