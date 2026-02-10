package dataproviders

import (
	"os"
	"path/filepath"
)

func HasWorktrees(parentPath string) (bool, error) {
	worktreePath := filepath.Join(parentPath, ".git", "worktrees")
	fileStats, err := os.Stat(worktreePath)
	if err != nil {
		return false, err
	}

	if fileStats.IsDir() {
		return true, nil
	}

	return false, nil
}

func GetSubdirectories(parentPath string) []Item {
	var subItems []Item
	worktreeDefinitions := filepath.Join(parentPath, ".git", "worktrees")

	entries, err := os.ReadDir(worktreeDefinitions)
	if err != nil {
		return subItems
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		
		definitionFile := filepath.Join(worktreeDefinitions, entry.Name(), "gitdir")
		dirPointerBytes, err := os.ReadFile(definitionFile)
		if err != nil {
			continue
		}

		filePointer := string(dirPointerBytes)
		itemDir := filepath.Dir(filePointer)
		

		display := "[ ] " + entry.Name()

		subItems = append(subItems, Item{
			Id:         entry.Name(),
			Display:    display,
			Path:       itemDir,
			IsWorktree: false,
			TreeLevel:  1,
		})
	}

	return subItems
}
