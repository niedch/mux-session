package previewproviders

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/glamour"
	"github.com/niedch/mux-session/internal/dataproviders"
)

// ReadmePreviewProvider renders README.md files from project directories
type ReadmePreviewProvider struct {
	renderer *glamour.TermRenderer
}

// NewReadmePreviewProvider creates a new README preview provider
func NewReadmePreviewProvider(width int) (*ReadmePreviewProvider, error) {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return nil, err
	}

	return &ReadmePreviewProvider{
		renderer: renderer,
	}, nil
}

// Render generates the README preview for the given item
func (r *ReadmePreviewProvider) Render(item interface{}) (string, error) {
	dpItem, ok := item.(*dataproviders.Item)
	if !ok {
		return "", fmt.Errorf("expected *dataproviders.Item, got %T", item)
	}

	readmePath := filepath.Join(dpItem.Path, "README.md")

	_, err := os.Stat(readmePath)
	if os.IsNotExist(err) {
		return fmt.Sprintf("No README.md found in %s", dpItem.Path), nil
	}

	data, err := os.ReadFile(readmePath)
	if err != nil {
		return "", fmt.Errorf("error reading README.md: %w", err)
	}

	rendered, err := r.renderer.Render(string(data))
	if err != nil {
		return "", fmt.Errorf("error rendering markdown: %w", err)
	}

	return rendered, nil
}

// Name returns the identifier name of this provider
func (r *ReadmePreviewProvider) Name() string {
	return "readme"
}

// SetWidth updates the renderer with a new width for word wrapping
func (r *ReadmePreviewProvider) SetWidth(width int) error {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return err
	}
	r.renderer = renderer
	return nil
}
