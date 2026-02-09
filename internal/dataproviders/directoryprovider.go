package dataproviders

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/niedch/mux-session/internal/git"
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

				display := fullPath
				isWorktree := false
				g := git.NewGit(fullPath)
				if g.IsGitRepository() {
					var err error
					isWorktree, err = g.IsWorktree()
					if err == nil && isWorktree {
						display = "[w] " + fullPath
					}
				}

				item := Item{
					Id:         filepath.Base(fullPath),
					Display:    display,
					Path:       fullPath,
					IsWorktree: isWorktree,
				}

				// If this is a worktree, scan for subdirectories
				if isWorktree {
					subItems := dp.getSubdirectories(fullPath)
					if len(subItems) > 0 {
						item.SubItems = subItems
					}
				}

				dirs = append(dirs, item)
			}
		}
	}

	return dirs, nil
}

// getSubdirectories scans a directory for immediate subdirectories
func (dp *DirectoryProvider) getSubdirectories(parentPath string) []Item {
	var subItems []Item

	entries, err := os.ReadDir(parentPath)
	if err != nil {
		return subItems
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Skip hidden directories
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		subPath := filepath.Join(parentPath, entry.Name())

		// Check if subdirectory is a git repository
		g := git.NewGit(subPath)
		isGitRepo := g.IsGitRepository()

		display := "[ ] " + entry.Name() + "/"
		if isGitRepo {
			isWorktree, _ := g.IsWorktree()
			if isWorktree {
				display = "[w] " + entry.Name() + "/"
			}
		}

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
