package app

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Bold(true)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(2)
)

type SearchPort struct {
	viewport    viewport.Model
	keymap      *keymap
	help        help.Model
	searchInput textinput.Model
	list        *searchList
}

func newSearchPort(width, height int) *SearchPort {
	viewport := viewport.New(width, height)
	keymap := newKeymap()
	help := help.New()

	ti := textinput.New()
	ti.Placeholder = "search..."
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = 20

	items := []string{
		"Ramen",
		"Tomato Soup",
		"Hamburgers",
		"Cheeseburgers",
		"Currywurst",
		"Okonomiyaki",
		"Pasta",
		"Fillet Mignon",
		"Caviar",
		"Just Wine",
	}

	s := &SearchPort{
		viewport:    viewport,
		keymap:      keymap,
		help:        help,
		searchInput: ti,
		list:        newSearchList(items),
	}

	return s
}

func (s *SearchPort) Width() int {
	return s.viewport.Width
}

func (s *SearchPort) Height() int {
	return s.viewport.Height
}

func (s *SearchPort) SetContent(content string) {
	s.viewport.SetContent(content)
}

func (s *SearchPort) UpdateLayout(width, height int) {
	s.viewport.Width = width
	s.viewport.Height = height
}

func (s *SearchPort) Update(msg tea.Msg) (*SearchPort, tea.Cmd) {
	var cmd tea.Cmd

	s.viewport, cmd = s.viewport.Update(msg)
	s.searchInput, cmd = s.searchInput.Update(msg)
	s.list.SetSearch(s.searchInput.Value())
	s.list.Update(msg)

	return s, cmd
}

func (s *SearchPort) View() string {
	footer := lipgloss.NewStyle().
		Width(s.Width()).
		Render("Search: " + s.searchInput.View())

	helpView := s.help.View(s.keymap)
	footerHeight := lipgloss.Height(footer)
	helpHeight := lipgloss.Height(helpView)

	availableHeight := max(s.Height()-footerHeight-helpHeight, 0)

	s.list.SetWidth(s.Width())
	listArea := s.list.View(availableHeight)

	return lipgloss.JoinVertical(lipgloss.Top,
		listArea,
		footer,
		helpView,
	)
}
