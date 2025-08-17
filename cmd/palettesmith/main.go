package main

import (
	"flag"
	"fmt"
	"os"
	"palettesmith/internal/config"
	"palettesmith/internal/plugin"
	"palettesmith/internal/validation"
	"palettesmith/ui/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Parse command line flags
	validatePlugins := flag.Bool("validate-plugins", false, "Validate plugin configurations and exit")
	flag.Parse()

	// Handle plugin validation flag
	if *validatePlugins {
		validateAndExit()
		return
	}

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

// validateAndExit validates all plugins and exits with appropriate status code
func validateAndExit() {
	fmt.Println("Validating plugins...")
	
	// Discover all plugins
	result := plugin.Discover()
	
	// Report discovered plugins
	plugins := result.Store.List()
	fmt.Printf("Found %d plugins:\n", len(plugins))
	
	validator := validation.NewPluginValidator()
	hasErrors := false
	
	// Validate each plugin
	for _, p := range plugins {
		fmt.Printf("\n▸ %s (%s)\n", p.Manifest.ID, p.Manifest.Title)
		
		// Check for loading errors
		if result.Store.HasErrors(p.Manifest.ID) {
			hasErrors = true
			fmt.Printf("  ✗ Plugin loading errors:\n")
			for _, err := range result.Store.Errors()[p.Manifest.ID] {
				fmt.Printf("    - %v\n", err)
			}
		}
		
		// Validate plugin structure
		pluginErrors := validator.ValidatePlugin(p)
		if len(pluginErrors) > 0 {
			hasErrors = true
			fmt.Printf("  ✗ Plugin validation errors:\n")
			for _, err := range pluginErrors {
				fmt.Printf("    - %v\n", err)
			}
		}
		
		// Validate field defaults
		fieldErrorCount := 0
		for _, field := range p.Spec.Fields {
			if field.Default != "" {
				errors := validator.ValidateField(field, field.Default)
				if len(errors) > 0 {
					if fieldErrorCount == 0 {
						hasErrors = true
						fmt.Printf("  ✗ Field validation errors:\n")
					}
					fieldErrorCount++
					fmt.Printf("    - Field '%s': %s\n", field.Key, errors[0].Message)
				}
			}
		}
		
		// Show success if no errors
		if !result.Store.HasErrors(p.Manifest.ID) && len(pluginErrors) == 0 && fieldErrorCount == 0 {
			fmt.Printf("  ✓ Plugin is valid\n")
		}
	}
	
	// Report global errors
	if len(result.Errors) > 0 {
		hasErrors = true
		fmt.Printf("\nGlobal errors:\n")
		for id, errs := range result.Errors {
			if id == "system" {
				for _, err := range errs {
					fmt.Printf("  ✗ %v\n", err)
				}
			}
		}
	}
	
	// Summary
	fmt.Printf("\nValidation complete: ")
	if hasErrors {
		fmt.Printf("✗ Found errors\n")
		os.Exit(1)
	} else {
		fmt.Printf("✓ All plugins valid\n")
		os.Exit(0)
	}
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
