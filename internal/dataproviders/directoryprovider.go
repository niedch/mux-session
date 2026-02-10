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

				containsWorktrees, _ := dp.hasWorktrees(fullPath)
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
					subItems := dp.getSubdirectories(fullPath)
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

func (dp *DirectoryProvider) hasWorktrees(parentPath string) (bool, error) {
	worktreePath := filepath.Join(parentPath, "worktrees")
	fileStats, err := os.Stat(worktreePath)
	if err != nil {
		return false, err
	}

	if fileStats.IsDir() {
		return true, nil
	}

	return false, nil
}

// getSubdirectories scans a directory for immediate subdirectories
func (dp *DirectoryProvider) getSubdirectories(parentPath string) []Item {
	var subItems []Item
	worktreePath := filepath.Join(parentPath, "worktrees")

	entries, err := os.ReadDir(worktreePath)
	if err != nil {
		return subItems
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		subPath := filepath.Join(worktreePath, entry.Name())
		display := "[ ] " + entry.Name()

		subItems = append(subItems, Item{
			Id:         entry.Name(),
			Display:    display,
			Path:       subPath,
			IsWorktree: false,
			TreeLevel:  1,
		})
	}

	return subItems
}
