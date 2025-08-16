package main

import (
	"fmt"
	"os"
	"palettesmith/internal/config"
	"palettesmith/ui/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	configManager := initializeConfig()
	runApplication(configManager)
}

// initializeConfig creates and returns a new config manager, exiting on failure
func initializeConfig() *config.Manager {
	configManager, err := config.NewManager()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize configuration: %v\n", err)
		os.Exit(1)
	}
	return configManager
}

// createInitialModel determines which model to start with based on config state
func createInitialModel(configManager *config.Manager) tea.Model {
	if configManager.IsFirstRun() {
		return tui.NewSetupModel()
	}
	return tui.New()
}

// handleSetupCompletion processes setup completion and configures the chosen preset
func handleSetupCompletion(setupModel tui.SetupModel, configManager *config.Manager) {
	preset := setupModel.GetSelectedPreset()
	if err := configManager.SetPreset(preset); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to set preset '%s': %v\n", preset, err)
		os.Exit(1)
	}

	configManager.MarkSetupComplete()
	if err := configManager.SaveConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to save configuration: %v\n", err)
		os.Exit(1)
	}
}

// runApplication executes the main application loop, handling setup transitions
func runApplication(configManager *config.Manager) {
	initialModel := createInitialModel(configManager)
	p := tea.NewProgram(initialModel, tea.WithAltScreen())

	for {
		finalModel, err := p.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Application error: %v\n", err)
			os.Exit(1)
		}

		// Check if we completed setup and need to transition to main app
		if setupModel, ok := finalModel.(tui.SetupModel); ok && setupModel.IsConfirmed() {
			handleSetupCompletion(setupModel, configManager)
			// Start the main application
			p = tea.NewProgram(tui.New(), tea.WithAltScreen())
			continue
		}

		// Normal exit
		break
	}
}
