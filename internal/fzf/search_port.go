package fzf

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/niedch/mux-session/internal/dataproviders"
)

type searchPort struct {
	textInput textinput.Model
	help      help.Model
	keymap    keymap
	list      *list
	width     int
	height    int
}

func newSearchPort(items []dataproviders.Item, width, height int) *searchPort {
	ti := textinput.New()
	ti.Placeholder = "search..."
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = 20

	h := help.New()
	km := keymap{}

	return &searchPort{
		textInput: ti,
		help:      h,
		keymap:    km,
		list:      newList(items),
		width:     width,
		height:    height,
	}
}

func (sp *searchPort) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "ctrl+p":
			sp.list.moveUp()
			return nil
		case "down", "ctrl+n":
			sp.list.moveDown()
			return nil
		}
	}

	var cmd tea.Cmd
	oldValue := sp.textInput.Value()
	sp.textInput, cmd = sp.textInput.Update(msg)
	if sp.textInput.Value() != oldValue {
		sp.list.updateFilter(sp.textInput.Value())
	}

	return cmd
}

func (sp *searchPort) SetSize(width, height int) {
	sp.width = width
	sp.height = height
}

func (sp *searchPort) GetSelected() *dataproviders.Item {
	return sp.list.getSelected()
}

func (sp *searchPort) View() string {
	listHeight := max(sp.height-inputHeight-helpHeight, 1)

	start, end := sp.list.calculateVisibleRange(listHeight)

	var renderedItems []string
	for i := start; i < end; i++ {
		renderedItems = append(renderedItems, sp.list.renderItem(i))
	}

	var s strings.Builder

	// Fill remaining empty lines to push the list to the bottom.
	padding := listHeight - len(renderedItems)
	if padding > 0 {
		s.WriteString(strings.Repeat("\n", padding))
	}

	s.WriteString(strings.Join(renderedItems, ""))

	// Render search input and help at the bottom.
	s.WriteString("Search: " + sp.textInput.View() + "\n")
	s.WriteString(sp.help.View(sp.keymap))
	return s.String()
}
