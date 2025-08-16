package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSetupModel(t *testing.T) {
	t.Run("should_initialize_with_generic_as_default_preset", func(t *testing.T) {
		model := NewSetupModel()

		assert.Equal(t, "generic", model.selectedPreset)
		assert.False(t, model.showConfirm)
		assert.Equal(t, "yes", model.confirmChoice)
		assert.False(t, model.confirmed)
	})
}

func TestSetupModel_Init(t *testing.T) {
	t.Run("should_return_nil_command", func(t *testing.T) {
		model := NewSetupModel()

		cmd := model.Init()

		assert.Nil(t, cmd)
	})
}

func TestSetupModel_Update_KeyPresses(t *testing.T) {
	t.Run("should_show_confirmation_on_g_key", func(t *testing.T) {
		model := NewSetupModel()

		updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}})

		setupModel := updatedModel.(SetupModel)
		assert.Equal(t, "generic", setupModel.selectedPreset)
		assert.True(t, setupModel.showConfirm)
		assert.False(t, setupModel.confirmed)
		assert.Nil(t, cmd) // No message sent yet, just showing confirmation
	})

	t.Run("should_show_confirmation_on_o_key", func(t *testing.T) {
		model := NewSetupModel()

		updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'o'}})

		setupModel := updatedModel.(SetupModel)
		assert.Equal(t, "omarchy", setupModel.selectedPreset)
		assert.True(t, setupModel.showConfirm)
		assert.False(t, setupModel.confirmed)
		assert.Nil(t, cmd) // No message sent yet, just showing confirmation
	})

	t.Run("should_select_yes_on_y_key_during_confirmation", func(t *testing.T) {
		model := NewSetupModel()
		model.selectedPreset = "omarchy"
		model.showConfirm = true
		model.confirmChoice = "no" // start with no

		updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

		setupModel := updatedModel.(SetupModel)
		assert.Equal(t, "omarchy", setupModel.selectedPreset)
		assert.True(t, setupModel.showConfirm)
		assert.Equal(t, "yes", setupModel.confirmChoice)
		assert.False(t, setupModel.confirmed) // not confirmed until Enter
		assert.Nil(t, cmd) // no command yet
	})

	t.Run("should_confirm_and_send_message_on_enter_when_yes_selected", func(t *testing.T) {
		model := NewSetupModel()
		model.selectedPreset = "omarchy"
		model.showConfirm = true
		model.confirmChoice = "yes"

		updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})

		setupModel := updatedModel.(SetupModel)
		assert.Equal(t, "omarchy", setupModel.selectedPreset)
		assert.True(t, setupModel.showConfirm)
		assert.True(t, setupModel.confirmed)

		require.NotNil(t, cmd)
		msg := cmd()
		chosenMsg, ok := msg.(PresetChosenMsg)
		require.True(t, ok)
		assert.Equal(t, "omarchy", chosenMsg.Preset)
	})

	t.Run("should_select_no_on_n_key_during_confirmation", func(t *testing.T) {
		model := NewSetupModel()
		model.selectedPreset = "omarchy"
		model.showConfirm = true
		model.confirmChoice = "yes" // start with yes

		updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

		setupModel := updatedModel.(SetupModel)
		assert.Equal(t, "omarchy", setupModel.selectedPreset)
		assert.True(t, setupModel.showConfirm)
		assert.Equal(t, "no", setupModel.confirmChoice)
		assert.False(t, setupModel.confirmed)
		assert.Nil(t, cmd)
	})

	t.Run("should_go_back_to_selection_on_enter_when_no_selected", func(t *testing.T) {
		model := NewSetupModel()
		model.selectedPreset = "omarchy"
		model.showConfirm = true
		model.confirmChoice = "no"

		updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})

		setupModel := updatedModel.(SetupModel)
		assert.Equal(t, "omarchy", setupModel.selectedPreset)
		assert.False(t, setupModel.showConfirm)
		assert.False(t, setupModel.confirmed)
		assert.Nil(t, cmd)
	})

	t.Run("should_toggle_choice_with_arrow_keys_on_confirmation", func(t *testing.T) {
		model := NewSetupModel()
		model.showConfirm = true
		model.confirmChoice = "yes"

		// Arrow down should go to "no"
		updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyDown})
		setupModel := updatedModel.(SetupModel)
		assert.Equal(t, "no", setupModel.confirmChoice)
		assert.Nil(t, cmd)

		// Arrow up should go back to "yes"
		updatedModel, cmd = setupModel.Update(tea.KeyMsg{Type: tea.KeyUp})
		setupModel = updatedModel.(SetupModel)
		assert.Equal(t, "yes", setupModel.confirmChoice)
		assert.Nil(t, cmd)
	})

	t.Run("should_handle_uppercase_keys", func(t *testing.T) {
		model := NewSetupModel()

		// Test uppercase G shows confirmation
		updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}})
		setupModel := updatedModel.(SetupModel)
		assert.Equal(t, "generic", setupModel.selectedPreset)
		assert.True(t, setupModel.showConfirm)
		assert.False(t, setupModel.confirmed)

		// Reset model for uppercase O test
		model = NewSetupModel()
		updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'O'}})
		setupModel = updatedModel.(SetupModel)
		assert.Equal(t, "omarchy", setupModel.selectedPreset)
		assert.True(t, setupModel.showConfirm)
		assert.False(t, setupModel.confirmed)
	})

	t.Run("should_quit_on_q_key", func(t *testing.T) {
		model := NewSetupModel()

		updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

		setupModel := updatedModel.(SetupModel)
		assert.False(t, setupModel.confirmed) // Should not be confirmed

		// Test that quit command is returned (tea.Quit is a function)
		assert.NotNil(t, cmd)
		assert.Equal(t, tea.Quit(), cmd())
	})

	t.Run("should_quit_on_ctrl_c", func(t *testing.T) {
		model := NewSetupModel()

		updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})

		setupModel := updatedModel.(SetupModel)
		assert.False(t, setupModel.confirmed)

		// Test that quit command is returned
		assert.NotNil(t, cmd)
		assert.Equal(t, tea.Quit(), cmd())
	})
}

func TestSetupModel_Update_ArrowKeyNavigation(t *testing.T) {
	t.Run("should_toggle_to_omarchy_when_up_key_pressed_from_generic", func(t *testing.T) {
		model := NewSetupModel() // starts with generic

		updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyUp})

		setupModel := updatedModel.(SetupModel)
		assert.Equal(t, "omarchy", setupModel.selectedPreset)
		assert.False(t, setupModel.confirmed) // Should not confirm yet
		assert.Nil(t, cmd)
	})

	t.Run("should_toggle_to_generic_when_down_key_pressed_from_omarchy", func(t *testing.T) {
		model := NewSetupModel()
		model.selectedPreset = "omarchy"

		updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyDown})

		setupModel := updatedModel.(SetupModel)
		assert.Equal(t, "generic", setupModel.selectedPreset)
		assert.False(t, setupModel.confirmed)
		assert.Nil(t, cmd)
	})

	t.Run("should_show_confirmation_on_enter", func(t *testing.T) {
		model := NewSetupModel() // starts with generic

		updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})

		setupModel := updatedModel.(SetupModel)
		assert.Equal(t, "generic", setupModel.selectedPreset)
		assert.True(t, setupModel.showConfirm)
		assert.False(t, setupModel.confirmed)
		assert.Nil(t, cmd) // No message sent yet, just showing confirmation
	})

	t.Run("should_show_confirmation_for_omarchy_when_navigated_to_and_enter_pressed", func(t *testing.T) {
		model := NewSetupModel()
		// Navigate to omarchy
		updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyUp})
		model = updatedModel.(SetupModel)

		updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})

		setupModel := updatedModel.(SetupModel)
		assert.Equal(t, "omarchy", setupModel.selectedPreset)
		assert.True(t, setupModel.showConfirm)
		assert.False(t, setupModel.confirmed)
		assert.Nil(t, cmd) // No message sent yet, just showing confirmation
	})
}

func TestSetupModel_View(t *testing.T) {
	t.Run("should_show_confirmation_screen_with_yes_selected", func(t *testing.T) {
		model := NewSetupModel()
		model.showConfirm = true
		model.selectedPreset = "generic"
		model.confirmChoice = "yes"

		view := model.View()

		assert.Contains(t, view, "Configure Palettesmith with Generic")
		assert.Contains(t, view, "▸ [Y] Yes, continue to theming") // should be highlighted
		assert.Contains(t, view, "  [N] No, change setup")        // should not be highlighted
		assert.Contains(t, view, "↑/↓ to navigate")
	})

	t.Run("should_show_confirmation_screen_with_no_selected", func(t *testing.T) {
		model := NewSetupModel()
		model.showConfirm = true
		model.selectedPreset = "omarchy"
		model.confirmChoice = "no"

		view := model.View()

		assert.Contains(t, view, "Configure Palettesmith with Omarchy")
		assert.Contains(t, view, "  [Y] Yes, continue to theming") // should not be highlighted
		assert.Contains(t, view, "▸ [N] No, change setup")        // should be highlighted
		assert.Contains(t, view, "↑/↓ to navigate")
	})

	t.Run("should_show_setup_screen_when_not_confirmed", func(t *testing.T) {
		model := NewSetupModel()

		view := model.View()

		assert.Contains(t, view, "Welcome to Palettesmith")
		assert.Contains(t, view, "Choose your configuration")
		assert.Contains(t, view, "[G] Generic")
		assert.Contains(t, view, "[O] Omarchy")
		assert.Contains(t, view, "Press O/G to select")
	})

	t.Run("should_highlight_generic_option_when_selected", func(t *testing.T) {
		model := NewSetupModel() // defaults to generic

		view := model.View()

		// Generic should have the selection indicator
		assert.Contains(t, view, "▸ [G] Generic")
		// Omarchy should not have the indicator
		assert.Contains(t, view, "  [O] Omarchy")
	})

	t.Run("should_highlight_omarchy_option_when_selected", func(t *testing.T) {
		model := NewSetupModel()
		model.selectedPreset = "omarchy"

		view := model.View()

		// Omarchy should have the selection indicator
		assert.Contains(t, view, "▸ [O] Omarchy")
		// Generic should not have the indicator
		assert.Contains(t, view, "  [G] Generic")
	})
}

func TestSetupModel_Integration(t *testing.T) {
	t.Run("should_handle_complete_user_flow_with_navigation", func(t *testing.T) {
		model := NewSetupModel()

		// User navigates up to omarchy
		updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyUp})
		model = updatedModel.(SetupModel)
		assert.Equal(t, "omarchy", model.selectedPreset)
		assert.False(t, model.confirmed)
		assert.False(t, model.showConfirm)

		// User navigates back down to generic
		updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
		model = updatedModel.(SetupModel)
		assert.Equal(t, "generic", model.selectedPreset)
		assert.False(t, model.confirmed)
		assert.False(t, model.showConfirm)

		// User presses enter to show confirmation
		updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
		setupModel := updatedModel.(SetupModel)
		assert.Equal(t, "generic", setupModel.selectedPreset)
		assert.True(t, setupModel.showConfirm)
		assert.False(t, setupModel.confirmed)
		assert.Nil(t, cmd) // No message yet

		// User selects "yes" with Y
		updatedModel, cmd = setupModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
		setupModel = updatedModel.(SetupModel)
		assert.Equal(t, "yes", setupModel.confirmChoice)
		assert.False(t, setupModel.confirmed) // not confirmed until Enter
		assert.Nil(t, cmd) // no command yet

		// User confirms with Enter
		updatedModel, cmd = setupModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
		setupModel = updatedModel.(SetupModel)
		assert.True(t, setupModel.confirmed)

		// Should send correct message
		require.NotNil(t, cmd)
		msg := cmd()
		chosenMsg, ok := msg.(PresetChosenMsg)
		require.True(t, ok)
		assert.Equal(t, "generic", chosenMsg.Preset)
	})
}

func TestPresetChosenMsg(t *testing.T) {
	t.Run("should_create_message_with_correct_preset", func(t *testing.T) {
		msg := PresetChosenMsg{Preset: "omarchy"}
		assert.Equal(t, "omarchy", msg.Preset)

		msg = PresetChosenMsg{Preset: "generic"}
		assert.Equal(t, "generic", msg.Preset)
	})
}

func TestSetupModel_IsConfirmed(t *testing.T) {
	t.Run("should_return_false_when_not_confirmed", func(t *testing.T) {
		model := NewSetupModel()
		assert.False(t, model.IsConfirmed())
	})

	t.Run("should_return_true_when_confirmed", func(t *testing.T) {
		model := NewSetupModel()
		model.confirmed = true
		assert.True(t, model.IsConfirmed())
	})
}

func TestSetupModel_GetSelectedPreset(t *testing.T) {
	t.Run("should_return_default_generic_preset", func(t *testing.T) {
		model := NewSetupModel()
		assert.Equal(t, "generic", model.GetSelectedPreset())
	})

	t.Run("should_return_updated_preset_after_navigation", func(t *testing.T) {
		model := NewSetupModel()

		// Navigate to omarchy
		updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyUp})
		model = updatedModel.(SetupModel)

		assert.Equal(t, "omarchy", model.GetSelectedPreset())
	})
}
