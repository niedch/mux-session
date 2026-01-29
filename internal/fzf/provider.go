package fzf

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// DataProvider provides data for fzf selection
type DataProvider interface {
	GetItems() ([]string, error)

	GetPreview(index int, width, height int) string
	GetDisplayString(index int) string
}

// DirectoryProvider implements DataProvider for directory browsing
type DirectoryProvider struct {
	searchPaths      []string
	items            []string
}

// NewDirectoryProvider creates a new directory provider
func NewDirectoryProvider(searchPaths []string) *DirectoryProvider {
	return &DirectoryProvider{
		searchPaths:      searchPaths,
	}
}

// GetItems returns the directories to display
func (dp *DirectoryProvider) GetItems() ([]string, error) {
	if len(dp.items) > 0 {
		return dp.items, nil
	}

	var dirs []string
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
				dirs = append(dirs, fullPath)
			}
		}
	}

	dp.items = dirs
	return dp.items, nil
}

// GetPreview returns preview content for a directory
func (dp *DirectoryProvider) GetPreview(index int, width, height int) string {
	if index == -1 || index >= len(dp.items) {
		return ""
	}

	dirPath := dp.items[index]

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Sprintf("Error reading %s: %v", dirPath, err)
	}

	var preview strings.Builder
	fmt.Fprintf(&preview, "Directory: %s\n\n", dirPath)

	for _, entry := range entries {
		if entry.IsDir() {
			fmt.Fprintf(&preview, "📁 %s/\n", entry.Name())
		} else {
			fmt.Fprintf(&preview, "📄 %s\n", entry.Name())
		}
	}

	return preview.String()
}

// GetDisplayString returns the display string for the item
func (dp *DirectoryProvider) GetDisplayString(index int) string {
	if index >= len(dp.items) {
		return ""
	}
	return dp.items[index]
}
