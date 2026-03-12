package previewproviders

// PreviewProvider defines the interface for rendering content in the preview panel
type PreviewProvider interface {
	// Render generates the content to display in the preview panel
	// Returns the rendered string and any error that occurred
	Render(item interface{}) (string, error)

	// Name returns the identifier name of this provider
	Name() string

	// SetWidth updates the width for word wrapping
	SetWidth(width int) error

	// SetUpdateChan sets the channel to notify the UI of updates
	SetUpdateChan(ch chan<- struct{})
}
