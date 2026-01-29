package fzf

import (
	"github.com/koki-develop/go-fzf"
)

func StartFzf(provider DataProvider) (*string, error) {
	items, err := provider.GetItems()
	if err != nil {
		return nil, err
	}

	f, err := fzf.New(
		fzf.WithInputPosition(fzf.InputPositionBottom),
	)
	if err != nil {
		return nil, err
	}

	// Run fzf to select items
	selectedIndex, err := f.Find(items, func(i int) string {
		return provider.GetDisplayString(i)
	}, fzf.WithPreviewWindow(func(i, width, height int) string {
		return provider.GetPreview(i, width, height)
	}))

	if err != nil {
		return nil, err
	}

	// Assume only one is selected
	return &items[selectedIndex[0]], nil
}
