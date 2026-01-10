package fzf

import (
	"github.com/koki-develop/go-fzf"
)

// StartFzf starts fzf with the given data provider and returns selected indices
func StartFzf(provider DataProvider) ([]int, error) {
	items, err := provider.GetItems()
	if err != nil {
		return nil, err
	}

	// Create fzf instance
	f, err := fzf.New(
		fzf.WithInputPosition(fzf.InputPositionBottom),
	)
	if err != nil {
		return nil, err
	}

	// Run fzf to select items
	return f.Find(items, func(i int) string {
		return provider.GetDisplayString(i)
	}, fzf.WithPreviewWindow(func(i, width, height int) string {
		return provider.GetPreview(i, width, height)
	}))
}
