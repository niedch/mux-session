package dataproviders

import (
	"log"
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

				display := "[ ] " + fullPath

				containsWorktrees, _ := HasWorktrees(fullPath)
				if containsWorktrees {
					display = "[w] " + fullPath
				}

				item := Item{
					Id:         filepath.Base(fullPath),
					Display:    display,
					Path:       fullPath,
					IsWorktree: containsWorktrees,
				}

				// If this is a worktree, scan for subdirectories
				if containsWorktrees {
					subItems := GetSubdirectories(fullPath)
					if len(subItems) > 0 {
						log.Printf("Adding SubItems %d to %s", len(subItems), item.Display)
						item.SubItems = subItems
					}
				}

				dirs = append(dirs, item)
			}
		}
	}

	return dirs, nil
}

