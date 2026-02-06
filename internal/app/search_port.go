package app

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Bold(true)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(2)
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type SearchPort struct {
	viewport    viewport.Model
	keymap      *keymap
	help        help.Model
	searchInput textinput.Model
	list        list.Model
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

	items := []list.Item{
		item("Ramen"),
		item("Tomato Soup"),
		item("Hamburgers"),
		item("Cheeseburgers"),
		item("Currywurst"),
		item("Okonomiyaki"),
		item("Pasta"),
		item("Fillet Mignon"),
		item("Caviar"),
		item("Just Wine"),
	}

	list := list.New(items, itemDelegate{}, width, height)
	list.SetShowHelp(false)
	list.SetShowTitle(false)

	s := &SearchPort{
		viewport:    viewport,
		keymap:      keymap,
		help:        help,
		searchInput: ti,
		list:        list,
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
	s.list, cmd = s.list.Update(msg)

	return s, cmd
}

func (s *SearchPort) View() string {
	footer := lipgloss.NewStyle().
		Width(s.Width()).
		Render("Search: " + s.searchInput.View())

	listStyle := lipgloss.NewStyle().
		Width(s.Width()).
		Height(s.Height() - 1 - 1).
		Render(s.list.View())
	
	listBlock := lipgloss.PlaceVertical(120, lipgloss.Bottom, listStyle)


	return lipgloss.JoinHorizontal(lipgloss.Bottom,
		listBlock,
		footer,
		s.help.View(s.keymap),
	)
}
