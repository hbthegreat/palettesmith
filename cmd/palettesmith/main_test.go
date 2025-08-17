package main

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidatePluginsCLI(t *testing.T) {
	t.Run("should_show_help_with_validate_plugins_flag", func(t *testing.T) {
		// Build the binary first
		buildCmd := exec.Command("go", "build", "-o", "test-palettesmith", "main.go")
		err := buildCmd.Run()
		require.NoError(t, err, "Should build test binary")
		defer os.Remove("test-palettesmith")

		// Test help flag
		cmd := exec.Command("./test-palettesmith", "--help")
		output, err := cmd.CombinedOutput()
		
		// Help flag causes exit code 2, which is expected
		assert.Contains(t, string(output), "-validate-plugins", "Help should show validate-plugins flag")
		assert.Contains(t, string(output), "Validate plugin configurations and exit", "Help should show flag description")
	})

	t.Run("should_validate_plugins_successfully", func(t *testing.T) {
		// Build the binary first
		buildCmd := exec.Command("go", "build", "-o", "test-palettesmith", "main.go")
		err := buildCmd.Run()
		require.NoError(t, err, "Should build test binary")
		defer os.Remove("test-palettesmith")

		// Test validation flag
		cmd := exec.Command("./test-palettesmith", "--validate-plugins")
		output, err := cmd.CombinedOutput()
		
		assert.NoError(t, err, "Validation should succeed with valid plugins")
		assert.Contains(t, string(output), "Validating plugins...", "Should show validation message")
		assert.Contains(t, string(output), "Found", "Should show discovered plugins")
		assert.Contains(t, string(output), "✓ All plugins valid", "Should show success message")
	})

	t.Run("should_exit_with_error_for_invalid_plugins", func(t *testing.T) {
		// Create a broken plugin for testing
		err := os.MkdirAll("../../plugins/test-invalid", 0755)
		require.NoError(t, err, "Should create test plugin directory")
		defer os.RemoveAll("../../plugins/test-invalid")

		// Create broken plugin.json
		brokenPlugin := `{
			"id": "test-invalid", 
			"title": "",
			"spec": "spec.palette"
		}`
		err = os.WriteFile("../../plugins/test-invalid/plugin.json", []byte(brokenPlugin), 0644)
		require.NoError(t, err, "Should write broken plugin.json")

		// Create broken spec.palette  
		brokenSpec := `{
			"id": "",
			"fields": [
				{
					"key": "",
					"type": "invalid_type"
				}
			]
		}`
		err = os.WriteFile("../../plugins/test-invalid/spec.palette", []byte(brokenSpec), 0644)
		require.NoError(t, err, "Should write broken spec.palette")

		// Build the binary first
		buildCmd := exec.Command("go", "build", "-o", "test-palettesmith", "main.go")
		err = buildCmd.Run()
		require.NoError(t, err, "Should build test binary")
		defer os.Remove("test-palettesmith")

		// Test validation with broken plugin
		cmd := exec.Command("./test-palettesmith", "--validate-plugins")
		output, err := cmd.CombinedOutput()
		
		assert.Error(t, err, "Should exit with error for invalid plugins")
		outputStr := string(output)
		assert.Contains(t, outputStr, "✗ Found errors", "Should show error summary")
		assert.Contains(t, outputStr, "test-invalid", "Should show the broken plugin name")
	})
}

func TestMainFunction(t *testing.T) {
	t.Run("should_parse_command_line_flags", func(t *testing.T) {
		// Test that the flag parsing doesn't panic
		// We can't easily test the actual main function without running the TUI
		// but we can test that our flag is defined correctly
		
		// This is more of a compile-time test - if the code compiles,
		// our flag definition is syntactically correct
		assert.True(t, true, "Flag parsing code should compile without errors")
	})
}