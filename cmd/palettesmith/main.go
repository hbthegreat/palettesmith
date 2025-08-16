package main

import (
	"fmt"
	"os"
	"palettesmith/internal/config"
	"palettesmith/ui/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Initialize config manager
	configManager, err := config.NewManager()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize configuration: %v\n", err)
		os.Exit(1)
	}

	// Check if this is the first run
	isFirstRun := configManager.IsFirstRun()

	var initialModel tea.Model
	if isFirstRun {
		// Show setup screen for first run
		initialModel = tui.NewSetupModel()
	} else {
		// Use existing config and go to main TUI
		initialModel = tui.New()
	}

	p := tea.NewProgram(initialModel, tea.WithAltScreen())

	// Handle setup completion
	for {
		finalModel, err := p.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Application error: %v\n", err)
			os.Exit(1)
		}

		// Check if we completed setup and need to transition to main app
		if setupModel, ok := finalModel.(tui.SetupModel); ok && setupModel.IsConfirmed() {
			// Get the chosen preset and configure it
			preset := setupModel.GetSelectedPreset()
			if err := configManager.SetPreset(preset); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to set preset '%s': %v\n", preset, err)
				os.Exit(1)
			}

			// Mark setup as complete and save the configuration
			configManager.MarkSetupComplete()
			if err := configManager.SaveConfig(); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to save configuration: %v\n", err)
				os.Exit(1)
			}

			// Now start the main application
			p = tea.NewProgram(tui.New(), tea.WithAltScreen())
			continue
		}

		// Normal exit
		break
	}
}
