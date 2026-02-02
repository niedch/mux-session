package fzf

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func StartApp(dataProvider DataProvider) (string, error) {
	items, err := dataProvider.GetItems()
	if err != nil {
		return "", err
	}

	p := tea.NewProgram(initialModel(items), tea.WithAltScreen())
	m, err := p.Run()
	if err != nil {
		return "", err
	}

	if model, ok := m.(model); ok {
		return model.selected, nil
	}
	return "", nil
}

type item struct {
	text  string
	index int
}

type model struct {
	textInput textinput.Model
	help      help.Model
	keymap    keymap
	items     []string
	filtered  []item
	cursor    int
	selected  string
	width     int
	height    int
}

type keymap struct{}

func (k keymap) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("up", "ctrl+p"), key.WithHelp("↑/ctrl+p", "up")),
		key.NewBinding(key.WithKeys("down", "ctrl+n"), key.WithHelp("↓/ctrl+n", "down")),
		key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
		key.NewBinding(key.WithKeys("esc", "ctrl+c"), key.WithHelp("esc", "quit")),
	}
}

func (k keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}

func initialModel(items []string) model {
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
				m.selected = m.filtered[m.cursor].text
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

func filterItems(items []string, query string) []item {
	var matches []item
	q := []rune(strings.ToLower(query))

	for i, t := range items {
		target := []rune(strings.ToLower(t))
		qi, ti := 0, 0
		for qi < len(q) && ti < len(target) {
			if q[qi] == target[ti] {
				qi++
			}
			ti++
		}
		if qi == len(q) {
			matches = append(matches, item{text: t, index: i})
		}
	}
	return matches
}

func (m model) View() string {
	// If height is not set (e.g. init), default to 10
	listHeight := 10
	if m.height > 0 {
		// listHeight is the space for the list AND the empty lines below it.
		// -1 for search input
		// -1 for help text
		listHeight = m.height - 2
		if listHeight < 1 {
			listHeight = 1
		}
	}

	var s strings.Builder

	// Render list
	start := 0
	end := len(m.filtered)

	if len(m.filtered) > listHeight {
		if m.cursor < listHeight/2 {
			start = 0
			end = listHeight
		} else if m.cursor > len(m.filtered)-listHeight/2 {
			start = len(m.filtered) - listHeight
			end = len(m.filtered)
		} else {
			start = m.cursor - listHeight/2
			end = m.cursor + listHeight/2
		}
	}

	var renderedItems []string
	for i := start; i < end; i++ {
		if i < 0 || i >= len(m.filtered) {
			continue
		}
		it := m.filtered[i]
		cursor := "  "
		style := lipgloss.NewStyle()
		if i == m.cursor {
			cursor = "> "
			style = style.Bold(true)
		}
		renderedItems = append(renderedItems, fmt.Sprintf("%s%s\n", cursor, style.Render(it.text)))
	}

	// Fill remaining empty lines to push the list to the bottom
	padding := listHeight - len(renderedItems)
	if padding > 0 {
		s.WriteString(strings.Repeat("\n", padding))
	}

	s.WriteString(strings.Join(renderedItems, ""))

	// Render search input and help at the bottom
	s.WriteString("Search: " + m.textInput.View() + "\n")
	s.WriteString(m.help.View(m.keymap))
	return s.String()
}
