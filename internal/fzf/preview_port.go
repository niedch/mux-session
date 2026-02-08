package fzf

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/niedch/mux-session/internal/previewproviders"
)

type previewPort struct {
	viewport viewport.Model
	provider previewproviders.PreviewProvider
	content  string
	width    int
	height   int
	lastItem interface{}
}

func newPreviewPort(provider previewproviders.PreviewProvider, width, height int) *previewPort {
	vp := viewport.New(width, height)

	return &previewPort{
		viewport: vp,
		provider: provider,
		width:    width,
		height:   height,
	}
}

func (p *previewPort) LoadItem(item interface{}) error {
	// Only reload if the item has changed
	if item == p.lastItem {
		return nil
	}
	p.lastItem = item

	if item == nil {
		p.content = ""
		p.viewport.SetContent(p.content)
		return nil
	}

	rendered, err := p.provider.Render(item)
	if err != nil {
		p.content = err.Error()
		p.viewport.SetContent(p.content)
		return err
	}

	p.content = rendered
	p.viewport.SetContent(p.content)
	return nil
}

func (p *previewPort) SetSize(width, height int) {
	p.width = width
	p.height = height
	p.viewport.Width = width
	p.viewport.Height = height

	// Update provider width if it supports it
	if provider, ok := p.provider.(*previewproviders.ReadmePreviewProvider); ok {
		provider.SetWidth(width)
	}
}

func (p *previewPort) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	p.viewport, cmd = p.viewport.Update(msg)
	return cmd
}

func (p *previewPort) View() string {
	return p.viewport.View()
}
