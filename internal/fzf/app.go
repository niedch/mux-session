package fzf

import (
	"log"
	"os"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/niedch/mux-session/internal/dataproviders"
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

func Run(dataProvider dataproviders.DataProvider) (*dataproviders.Item, error) {
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

type model struct {
	searchPort *searchPort
	preview    *viewport.Model
	selected   *dataproviders.Item
	width      int
	height     int
}

func initialModel(items []dataproviders.Item) model {
	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		log.Fatal(err)
	}

	sepWidth := 1
	availableWidth := w - sepWidth

	// Split available space
	leftVpWidth := availableWidth / 2
	rightVpWidth := availableWidth - leftVpWidth

	previewPort := viewport.New(rightVpWidth, h)

	return model{
		searchPort: newSearchPort(items, leftVpWidth, h),
		preview:    &previewPort,
	}
}

func (m model) Init() tea.Cmd {
	return m.searchPort.textInput.Focus()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

		sepWidth := 1
		availableWidth := msg.Width - sepWidth
		leftVpWidth := availableWidth / 2

		m.searchPort.SetSize(leftVpWidth, msg.Height)
		m.preview.Width = availableWidth - leftVpWidth
		m.preview.Height = msg.Height
	}

	cmd := m.searchPort.Update(msg)
	return m, cmd
}

func (m model) View() string {
	searchView := rightBorderStyle.Width(m.searchPort.width + 1).Render(m.searchPort.View())
	rightView := m.preview.View()

	return lipgloss.JoinHorizontal(lipgloss.Top, searchView, rightView)
}
