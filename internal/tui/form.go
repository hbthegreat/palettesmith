package tui

import (
	"fmt"
	"palettesmith/internal/plugin"
	"regexp"
	"strconv"
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
}

func newFormFromSpec(s plugin.Spec) formModel {
	makeInput := func(f plugin.Field) textinput.Model {
		ti := textinput.New()
		ti.Prompt = ""
		switch f.Type {
		case "color":
			ti.CharLimit = 7
			ti.Width = 8
			if f.Default != "" {
				ti.SetValue(f.Default)
			} else {
				ti.SetValue("#000000")
			}
		case "number":
			ti.CharLimit = 10
			ti.Width = 8
			if f.Default != "" {
				ti.SetValue(f.Default)
			} else {
				ti.SetValue("0")
			}
		case "select":
			// Render as text for now; later swap to a selector
			ti.CharLimit = 64
			ti.Width = 24
			if f.Default != "" {
				ti.SetValue(f.Default)
			}
		default: // "text"
			ti.CharLimit = 256
			ti.Width = 32
			ti.SetValue(f.Default)
		}
		return ti
	}

	out := make([]formField, 0, len(s.Fields))
	for _, f := range s.Fields {
		out = append(out, formField{spec: f, input: makeInput(f)})
	}
	fm := formModel{fields: out}
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
			f.fields[i].err = validateValue(f.fields[i].spec, f.fields[i].input.Value())
		}
	}
	return f, cmd
}

var colorRe = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)

func validateValue(spec plugin.Field, v string) string {
	switch spec.Type {
	case "color":
		if !colorRe.MatchString(v) {
			return "expect #RRGGBB"
		}
	case "number":
		n, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
		if err != nil {
			return "not a number"
		}
		if spec.Min != nil && n < *spec.Min {
			return fmt.Sprintf("< %.0f", *spec.Min)
		}
		if spec.Max != nil && n > *spec.Max {
			return fmt.Sprintf("> %.0f", *spec.Max)
		}
	}
	return ""
}

func (f formModel) View() string {
	var b strings.Builder
	for i, fld := range f.fields {
		cursor := "  "
		if i == f.focusIndex {
			cursor = "▸ "
		}

		label := lipgloss.NewStyle().Bold(true).Render(fld.spec.Label)
		swatch := ""
		if fld.spec.Type == "color" && colorRe.MatchString(fld.input.Value()) {
			swatch = " " + lipgloss.NewStyle().
				Background(lipgloss.Color(fld.input.Value())).
				Padding(0, 1).
				Render(" ")
		}
		help := ""
		if fld.spec.Help != "" {
			help = " " + lipgloss.NewStyle().Foreground(lipgloss.Color("#777777")).Render(fld.spec.Help)
		}
		err := ""
		if fld.err != "" {
			err = " " + lipgloss.NewStyle().Foreground(lipgloss.Color("#ff6b6b")).Render(fld.err)
		}

		fmt.Fprintf(&b, "%s%s %s%s%s%s\n", cursor, label, fld.input.View(), swatch, err, help)
	}
	b.WriteString("\nEnter to edit • ↑/↓ move • A Apply • Q quit")
	return b.String()
}
