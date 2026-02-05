package fzf

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/niedch/mux-session/internal/dataproviders"
)

const (
	inputHeight = 1
	helpHeight  = 1
)

func StartApp(dataProvider dataproviders.DataProvider) (*dataproviders.Item, error) {
	items, err := dataProvider.GetItems()
	if err != nil {
		return nil, err
	}

	p := tea.NewProgram(initialModel(items), tea.WithAltScreen())
	m, err := p.Run()
	if err != nil {
		return nil, err
	}

	if model, ok := m.(model); ok {
		return model.selected, nil
	}

	return nil, nil
}

type item struct {
	text    string
	index   int
	matches []int
}

type model struct {
	textInput textinput.Model
	help      help.Model
	keymap    keymap
	items     []dataproviders.Item
	filtered  []item
	cursor    int
	selected  *dataproviders.Item
	width     int
	height    int
}

func initialModel(items []dataproviders.Item) model {
	// reverse items
	for i, j := 0, len(items)-1; i < j; i, j = i+1, j-1 {
		items[i], items[j] = items[j], items[i]
	}

	ti := textinput.New()
	ti.Placeholder = "search..."
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = 20

	h := help.New()

	km := keymap{}

	m := model{
		textInput: ti,
		items:     items,
		help:      h,
		keymap:    km,
	}
	m.updateFiltered()
	if len(m.filtered) > 0 {
		m.cursor = len(m.filtered) - 1
	}
	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			if len(m.filtered) > 0 {
				m.selected = &m.items[m.filtered[m.cursor].index]
			}
			return m, tea.Quit
		case "up", "ctrl+p":
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil
		case "down", "ctrl+n":
			if m.cursor < len(m.filtered)-1 {
				m.cursor++
			}
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
		m.width = msg.Width
		m.height = msg.Height
	}

	var cmd tea.Cmd
	oldValue := m.textInput.Value()
	m.textInput, cmd = m.textInput.Update(msg)
	if m.textInput.Value() != oldValue {
		m.updateFiltered()
		if len(m.filtered) > 0 {
			m.cursor = len(m.filtered) - 1
		} else {
			m.cursor = 0
		}
	}

	return m, cmd
}

func (m *model) updateFiltered() {
	query := m.textInput.Value()
	m.filtered = filterItems(m.items, query)
}

func (m model) View() string {
	listHeight := max(m.height - inputHeight - helpHeight, 1)

	start, end := m.calculateVisibleRange(listHeight)

	var renderedItems []string
	for i := start; i < end; i++ {
		renderedItems = append(renderedItems, m.renderItem(i))
	}

	var s strings.Builder

	// Fill remaining empty lines to push the list to the bottom.
	padding := listHeight - len(renderedItems)
	if padding > 0 {
		s.WriteString(strings.Repeat("\n", padding))
	}

	s.WriteString(strings.Join(renderedItems, ""))

	// Render search input and help at the bottom.
	s.WriteString("Search: " + m.textInput.View() + "\n")
	s.WriteString(m.help.View(m.keymap))
	return s.String()
}

// calculateVisibleRange determines the start and end indices of the items to be displayed.
func (m model) calculateVisibleRange(listHeight int) (int, int) {
	if len(m.filtered) <= listHeight {
		return 0, len(m.filtered)
	}

	// When cursor is near the top.
	if m.cursor < listHeight/2 {
		return 0, listHeight
	}

	// When cursor is near the bottom.
	if m.cursor > len(m.filtered)-listHeight/2 {
		return len(m.filtered) - listHeight, len(m.filtered)
	}

	// When cursor is in the middle.
	start := m.cursor - listHeight/2
	end := m.cursor + listHeight/2
	return start, end
}

// renderItem renders a single item in the list.
func (m model) renderItem(i int) string {
	if i < 0 || i >= len(m.filtered) {
		return ""
	}

	it := m.filtered[i]
	isSelected := i == m.cursor

	cursor := "  "
	style := lipgloss.NewStyle()
	if isSelected {
		cursor = "> "
		style = style.Bold(true)
	}

	var highlightedText strings.Builder
	matchStyle := style.Copy().Bold(true)

	matchMap := make(map[int]struct{})
	for _, idx := range it.matches {
		matchMap[idx] = struct{}{}
	}

	for i, r := range it.text {
		if _, ok := matchMap[i]; ok {
			highlightedText.WriteString(matchStyle.Render(string(r)))
		} else {
			highlightedText.WriteString(style.Render(string(r)))
		}
	}

	return fmt.Sprintf("%s%s\n", cursor, highlightedText.String())
}
