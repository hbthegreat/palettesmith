// Package tui contains the Bubble Tea UI for Palettesmith.
package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	width, height int
	sidebar       Sidebar
}

func New() Model {
	return Model{
		sidebar: NewSidebar(),
	}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.sidebar.SetSize(28, m.height)
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.sidebar, cmd = m.sidebar.Update(msg)
	return m, cmd
}

var (
	sep        = lipgloss.NewStyle().Foreground(lipgloss.Color("#666"))
	rightStyle = lipgloss.NewStyle().Padding(1, 2)
	helpStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#7d7d7d"))
	titleStyle = lipgloss.NewStyle().Bold(true)
)

func (m Model) View() string {
	left := m.sidebar.View()

	separator := sep.Render(strings.Repeat("│", max(1, m.height-2)))

	target := m.sidebar.SelectedTitle()
	if target == "" {
		target = "Select a target"
	}
	body := fmt.Sprintf(
		"%s\n\n• Config path: (example) ~/.config/%s/...\n• Reload: (example) hyprctl reload\n• Notes: atomic writes + backups\n",
		titleStyle.Render(target),
		strings.ToLower(target),
	)

	rightWidth := max(40, m.width-32)
	right := rightStyle.Width(rightWidth).Render(body)

	footer := helpStyle.Render("↑/↓ navigate   / filter   q quit")

	return lipgloss.JoinHorizontal(lipgloss.Top, left, separator, right) + "\n" + footer + "\n"
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
