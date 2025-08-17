// Package tui contains the Bubble Tea UI for Palettesmith.
package tui

import (
	"fmt"
	"palettesmith/internal/plugin"
	"palettesmith/internal/theme"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type page int

const (
	pageExplainer page = iota
	pageForm
)

const sidebarW = 30

type Model struct {
	width, height int
	sidebar       Sidebar
	page          page
	store         *plugin.Store
	form          formModel
	specLoadedFor string
	status        string

	theme *theme.Store
}

func New() Model {
	result := plugin.Discover()
	st := result.Store

	items := []list.Item{}
	for _, p := range st.List() {
		items = append(items, targetItem{
			id:          p.Manifest.ID,
			title:       firstNonEmpty(p.Manifest.Title, p.Manifest.ID),
			description: "Themeable target",
		})
	}

	if len(items) == 0 {
		items = []list.Item{targetItem{id: "", title: "No plugins found", description: "Put plugins under ./plugins/<id>/"}}
	}

	th := theme.NewStore(theme.ThemeConfig{
		ThemeDefaults: map[string]string{
			"bg":     "#1e1e2e",
			"fg":     "#cdd6f4",
			"accent": "#89b4fa",
			"border": "#45475a",
		},
	})

	return Model{
		sidebar: NewSidebar(items),
		page:    pageExplainer,
		store:   st,
		theme:   th,
	}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) ensureFormFor(id string) Model {
	if id == "" || m.store == nil || m.specLoadedFor == id {
		return m
	}
	if plug, ok := m.store.Get(id); ok {
		m.form = newFormFromSpec(plug.Spec, id, m.theme)
		m.specLoadedFor = id
	}
	return m
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.sidebar.SetSize(sidebarW, m.height)

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab":
			if m.page == pageExplainer {
				m.page = pageForm
			} else {
				m.page = pageExplainer
			}
		case "a":
			sel := m.sidebar.SelectedID()
			if sel == "" {
				m.status = "No target selected"
				return m, clearAfter(2 * time.Second)
			}

			if plug, ok := m.store.Get(sel); ok {
				eff := map[string]string{}
				for _, f := range plug.Spec.Fields {
					eff[f.Key] = m.theme.Resolve(sel, f.Key, f.Default)
				}
				m.status = fmt.Sprintf("Apply (dry-run) %s: %v", sel, eff)
			} else {
				m.status = "Unknown target"
			}
			return m, clearAfter(2 * time.Second)
		}
	case statusClearMsg:
		m.status = ""
	}

	var cmd tea.Cmd
	m.sidebar, cmd = m.sidebar.Update(msg)
	if m.page == pageForm {
		m = m.ensureFormFor(m.sidebar.SelectedID())
		m.form, cmd = m.form.Update(msg)
	}
	return m, cmd
}


func (m Model) View() string {
	left := LeftPaneStyle.Width(sidebarW).Render(m.sidebar.View())
	rightWidth := max(40, m.width-sidebarW)

	title := m.sidebar.SelectedTitle()
	if title == "" {
		title = "Select a target"
	}

	tabs := lipgloss.JoinHorizontal(lipgloss.Top,
		boolStyle(m.page == pageExplainer, TabActiveStyle, TabDimStyle).Render("Explainer"),
		lipgloss.NewStyle().Padding(0, 1).Render("·"),
		boolStyle(m.page == pageForm, TabActiveStyle, TabDimStyle).Render("Form"),
	)

	selID := m.sidebar.SelectedID()
	var upaths, spaths, reload string
	if selID != "" && m.store != nil {
		if plug, ok := m.store.Get(selID); ok {
			upaths = strings.Join(plug.Manifest.UserPaths, ", ")
			spaths = strings.Join(plug.Manifest.SystemPaths, ", ")
			reload = strings.Join(plug.Manifest.Reload, " ")
		}
	}
	var body string
	switch m.page {
	case pageExplainer:
		body = fmt.Sprintf("%s\n\nThis target is provided by a plugin.\n• User paths: %s\n• System paths: %s\n• Reload: %s\n",
			TitleStyle.Render(title), nz(upaths, "—"), nz(spaths, "—"), nz(reload, "—"))
	case pageForm:
		body = TitleStyle.Render(title) + "\n\n" + m.form.View()
	}

	body = tabs + "\n\n" + body

	right := RightPaneStyle.Width(rightWidth).Render(body)

	statusLine := ""
	if m.status != "" {
		statusLine = StatusStyle.Render(m.status) + "\n"
	}
	var footerText string
	if m.page == pageForm {
		footerText = "Enter to edit • ↑/↓ Move • A Apply • Q Quit"
	} else {
		footerText = "Tab Explainer/Form • ↑/↓ Move • Q Quit • / Filter"
	}
	footer := HelpStyle.Render(footerText)
	return lipgloss.JoinHorizontal(lipgloss.Top, left, right) + "\n" + statusLine + footer + "\n"
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func firstNonEmpty(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

func nz(s, fallback string) string {
	if strings.TrimSpace(s) == "" {
		return fallback
	}
	return s
}

func clearAfter(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(time.Time) tea.Msg { return statusClearMsg{} })
}

type statusClearMsg struct{}

func boolStyle(ok bool, a, b lipgloss.Style) lipgloss.Style {
	if ok {
		return a
	}
	return b
}
