package dataproviders

import (
	"os"
	"path/filepath"
	"strings"
)

// DirectoryProvider implements DataProvider for directory browsing
type DirectoryProvider struct {
	searchPaths []string
}

// NewDirectoryProvider creates a new directory provider
func NewDirectoryProvider(searchPaths []string) *DirectoryProvider {
	return &DirectoryProvider{
		searchPaths: searchPaths,
	}
}

// GetItems returns the directories to display
func (dp *DirectoryProvider) GetItems() ([]Item, error) {
	var dirs []Item
	for _, searchPath := range dp.searchPaths {
		entries, err := os.ReadDir(searchPath)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			if !strings.HasPrefix(entry.Name(), ".") {
				fullPath := filepath.Join(searchPath, entry.Name())

				dirs = append(dirs, Item{
					Id:      filepath.Base(fullPath),
					Display: fullPath,
				})
			}
		}
	}

	return dirs, nil
}
