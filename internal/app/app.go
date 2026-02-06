package app

import (
	"os"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

var (
	// Right border for the left viewport to create a separator.
	// We use NormalBorder for a simple vertical line.
	rightBorderStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, true, false, false)
)

type Model struct {
	leftViewport  viewport.Model
	searchPort    *SearchPort
	rightViewport viewport.Model
	keymap        keymap
	ready         bool
	height        int
	width         int
}

func New() Model {
	m := Model{
		ready: false,
	}

	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err == nil {
		m.layout(w, h)
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) layout(width, height int) {
	// Use full height (no top/bottom borders)
	m.height = height
	m.width = width

	// Width calculation
	// We allocate 1 unit for the vertical separator (the right border of the left viewport)
	sepWidth := 1
	availableWidth := width - sepWidth

	// Split available space
	leftVpWidth := availableWidth / 2
	rightVpWidth := availableWidth - leftVpWidth

	if !m.ready {
		m.searchPort = newSearchPort(width, height)
		m.searchPort.SetContent("searchPort")

		m.rightViewport = viewport.New(rightVpWidth, m.height)
		m.rightViewport.SetContent("Right Viewport")

		m.ready = true
	} else {
		m.searchPort.UpdateLayout(leftVpWidth, m.height)

		m.rightViewport.Width = rightVpWidth
		m.rightViewport.Height = m.height
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.layout(msg.Width, msg.Height)
	}

	// Update both viewports
	m.searchPort, cmd = m.searchPort.Update(msg)
	cmds = append(cmds, cmd)

	m.rightViewport, cmd = m.rightViewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	searchView := rightBorderStyle.Width(m.searchPort.Width() + 1).Render(m.searchPort.View())

	rightView := m.rightViewport.View()

	return lipgloss.JoinHorizontal(lipgloss.Top, searchView, rightView)
}

func Run() error {
	p := tea.NewProgram(New(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
