package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type targetItem struct {
	id          string
	title       string
	description string
}

func (i targetItem) Title() string       { return i.title }
func (i targetItem) Description() string { return i.description }
func (i targetItem) FilterValue() string { return i.title }

type Sidebar struct {
	l list.Model
}

func NewSidebar(items []list.Item) Sidebar {
	delegate := list.NewDefaultDelegate()
	l := list.New(items, delegate, 28, 16)
	l.Title = "Targets"
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.SetShowPagination(false)
	return Sidebar{l: l}
}

func (s *Sidebar) SetSize(width, height int) {
	s.l.SetSize(width, height-2)
}

func (s Sidebar) Update(msg tea.Msg) (Sidebar, tea.Cmd) {
	var cmd tea.Cmd
	s.l, cmd = s.l.Update(msg)
	return s, cmd
}

func (s Sidebar) View() string {
	return s.l.View()
}

func (s Sidebar) SelectedID() string {
	if it, ok := s.l.SelectedItem().(targetItem); ok {
		return it.id
	}
	return ""
}

func (s Sidebar) SelectedTitle() string {
	if it, ok := s.l.SelectedItem().(targetItem); ok {
		return it.title
	}
	return ""
}
