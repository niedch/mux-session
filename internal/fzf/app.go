package fzf

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/niedch/mux-session/internal/conf"
	"github.com/niedch/mux-session/internal/dataproviders"
	"github.com/niedch/mux-session/internal/previewproviders"
	"golang.org/x/term"
)

const (
	inputHeight = 1
	helpHeight  = 1
)

var (
	rightBorderStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, true, false, false)
)

func Run(dataProvider dataproviders.DataProvider, config *conf.Config) (*dataproviders.Item, error) {
	items, err := dataProvider.GetItems()
	if err != nil {
		return nil, err
	}

	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		log.Fatal(err)
	}

	leftVpWidth, rightVpWidth := calculateLayout(w)

	previewProvider, err := previewproviders.CreatePreviewProvider(config, rightVpWidth)
	if err != nil {
		return nil, err
	}

	p := tea.NewProgram(initialModel(items, previewProvider, leftVpWidth, rightVpWidth, h), tea.WithAltScreen())
	m, err := p.Run()
	if err != nil {
		return nil, err
	}

	if model, ok := m.(model); ok {
		return model.selected, nil
	}

	return nil, nil
}

type model struct {
	searchPort    *searchPort
	previewPort   *previewPort
	selected      *dataproviders.Item
	lastSelection *dataproviders.Item
	width         int
	height        int
}

func initialModel(items []dataproviders.Item, provider previewproviders.PreviewProvider, leftVpWidth, rightVpWidth, h int) model {
	return model{
		searchPort:  newSearchPort(items, leftVpWidth, h),
		previewPort: newPreviewPort(provider, rightVpWidth, h),
	}
}

func (m model) Init() tea.Cmd {
	// Load initial README if items exist
	if item := m.searchPort.GetSelected(); item != nil {
		m.lastSelection = item
		m.previewPort.LoadItem(item)
	}
	return m.searchPort.textInput.Focus()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			m.selected = m.searchPort.GetSelected()
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		leftVpWidth, rightVpWidth := calculateLayout(m.width)

		m.searchPort.SetSize(leftVpWidth, msg.Height)
		m.previewPort.SetSize(rightVpWidth, msg.Height)
	}

	// Update searchPort first (this moves the cursor)
	cmd := m.searchPort.Update(msg)
	cmds = append(cmds, cmd)

	// Now check if selection changed and update preview
	currentSelection := m.searchPort.GetSelected()
	if currentSelection != nil && (m.lastSelection == nil || currentSelection.Id != m.lastSelection.Id) {
		m.lastSelection = currentSelection
		m.previewPort.LoadItem(currentSelection)
	}

	cmd = m.previewPort.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	searchView := rightBorderStyle.Width(m.searchPort.width + 1).Render(m.searchPort.View())
	previewView := m.previewPort.View()

	return lipgloss.JoinHorizontal(lipgloss.Top, searchView, previewView)
}

func calculateLayout(width int) (leftVpWidth, rightVpWidth int) {
	sepWidth := 1
	availableWidth := width - sepWidth
	leftVpWidth = availableWidth / 2
	rightVpWidth = availableWidth - leftVpWidth
	return
}
