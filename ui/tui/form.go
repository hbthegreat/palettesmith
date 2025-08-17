package tui

import (
	"fmt"
	"palettesmith/internal/plugin"
	"palettesmith/internal/theme"
	"palettesmith/internal/validation"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type formField struct {
	spec  plugin.Field
	input textinput.Model
	err   string
}

type formModel struct {
	fields     []formField
	focusIndex int
	labelW     int

	pluginID  string
	theme     *theme.Store
	validator *validation.PluginValidator
}

func newFormFromSpec(s plugin.Spec, pluginID string, th *theme.Store) formModel {
	makeInput := func(f plugin.Field, val string) textinput.Model {
		ti := textinput.New()
		ti.Prompt = ""
		switch f.Type {
		case "color":
			ti.CharLimit = 12
			ti.Width = 12
		case "number":
			ti.CharLimit = 10
			ti.Width = 8
		case "select":
			// Render as text for now; later swap to a selector
			ti.CharLimit = 64
			ti.Width = 24
		default: // "text"
			ti.CharLimit = 256
			ti.Width = 32
		}
		ti.SetValue(val)
		return ti
	}

	out := make([]formField, 0, len(s.Fields))
	lw := 0

	for _, f := range s.Fields {
		// Resolve initial value: override > theme default > field default
		val := f.Default
		if th != nil {
			val = th.Resolve(pluginID, f.Key, f.Default)
		}
		ff := formField{spec: f, input: makeInput(f, val)}
		out = append(out, ff)
		if n := len(f.Label); n > lw {
			lw = n
		}
	}

	fm := formModel{
		fields:    out,
		labelW:    lw,
		pluginID:  pluginID,
		theme:     th,
		validator: validation.NewPluginValidator(),
	}
	if len(fm.fields) > 0 {
		fm.fields[0].input.Focus()
	}
	return fm
}

func (f formModel) Palette() map[string]string {
	m := make(map[string]string, len(f.fields))
	for _, ff := range f.fields {
		m[ff.spec.Key] = ff.input.Value()
	}
	return m
}

func (f formModel) Update(msg tea.Msg) (formModel, tea.Cmd) {
	switch m := msg.(type) {
	case tea.KeyMsg:
		switch m.String() {
		case "up":
			if f.focusIndex > 0 {
				f.fields[f.focusIndex].input.Blur()
				f.focusIndex--
				f.fields[f.focusIndex].input.Focus()
			}
			return f, nil
		case "down":
			if f.focusIndex < len(f.fields)-1 {
				f.fields[f.focusIndex].input.Blur()
				f.focusIndex++
				f.fields[f.focusIndex].input.Focus()
			}
			return f, nil
		}
	}

	var cmd tea.Cmd
	for i := range f.fields {
		if i == f.focusIndex {
			f.fields[i].input, cmd = f.fields[i].input.Update(msg)
			
			// Use validation service instead of old validateValue function
			errors := f.validator.ValidateField(f.fields[i].spec, f.fields[i].input.Value())
			if len(errors) > 0 {
				f.fields[i].err = errors[0].Message
			} else {
				f.fields[i].err = ""
			}
			
			if f.theme != nil {
				f.theme.SetOverride(f.pluginID, f.fields[i].spec.Key, f.fields[i].input.Value())
			}
		}
	}
	return f, cmd
}

// isValidColor checks if a color value is valid for color swatch display
func isValidColor(value string) bool {
	// Use validation service to check if it's a valid color
	validator := validation.NewPluginValidator()
	field := plugin.Field{Type: "color"}
	errors := validator.ValidateField(field, value)
	return len(errors) == 0
}

func (f formModel) View() string {
	var b strings.Builder

	for i, fld := range f.fields {
		cursor := "  "
		if i == f.focusIndex {
			cursor = "▸ "
		}

		padded := fmt.Sprintf("%-*s", f.labelW, fld.spec.Label)
		label := lipgloss.NewStyle().Bold(true).Render(padded)

		swatch := ""
		if fld.spec.Type == "color" && isValidColor(fld.input.Value()) {
			swatch = " " + lipgloss.NewStyle().
				Background(lipgloss.Color(fld.input.Value())).
				Padding(0, 1).
				Render(" ")
		}

		help := ""
		if fld.spec.Help != "" {
			help = " " + FormHelpStyle.Render(fld.spec.Help)
		}

		err := ""
		if fld.err != "" {
			err = " " + FormErrorStyle.Render(fld.err)
		}

		source := "default"
		if f.theme != nil {
			if f.theme.HasOverride(f.pluginID, fld.spec.Key) {
				source = "override"
			} else if f.theme.HasDefault(fld.spec.Key) {
				source = "theme"
			}
		}

		tag := FormTagStyle.Render(" [" + source + "]")

		row := fmt.Sprintf("%s%s  %s%s%s%s%s", cursor, label, fld.input.View(), swatch, err, help, tag)

		if i == f.focusIndex {
			fmt.Fprintln(&b, FormRowFocusStyle.Render(row))
		} else {
			fmt.Fprintln(&b, FormRowNormalStyle.Render(row))
		}
	}

	return b.String()
}
