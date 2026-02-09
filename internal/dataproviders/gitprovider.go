package dataproviders

import (
	"github.com/niedch/mux-session/internal/git"
)

// GitProvider implements DataProvider for git repositories
type GitProvider struct {
	git *git.Git
}

// NewGitProvider creates a new git provider
func NewGitProvider(git *git.Git) *GitProvider {
	return &GitProvider{
		git: git,
	}
}

// GetItems returns information about the git repository
func (dp *GitProvider) GetItems() ([]Item, error) {
	var items []Item

	// Check if this is a git repository
	if !dp.git.IsGitRepository() {
		return items, nil
	}

	// Check if this is a worktree
	isWorktree, err := dp.git.IsWorktree()
	if err != nil {
		return items, nil
	}

	if isWorktree {
		mainPath, err := dp.git.GetMainRepositoryPath()
		if err == nil {
			items = append(items, Item{
				Id:      "main-repo",
				Display: "[WORKTREE] Main Repository: " + mainPath,
				Path:    mainPath,
			})
		}

		worktreeRoot, err := dp.git.GetWorktreeRoot()
		if err == nil {
			items = append(items, Item{
				Id:      "worktree-root",
				Display: "[WORKTREE] Worktree Root: " + worktreeRoot,
				Path:    worktreeRoot,
			})
		}
	}

	// Add current branch info
	branch, err := dp.git.GetCurrentBranch()
	if err == nil {
		items = append(items, Item{
			Id:      "current-branch",
			Display: "[GIT] Branch: " + branch,
			Path:    "",
		})
	}

	// List all worktrees
	worktrees, err := dp.git.ListWorktrees()
	if err == nil && len(worktrees) > 0 {
		for i, wt := range worktrees {
			display := "[GIT] Worktree: " + wt.Path
			if wt.Branch != "" {
				display += " (" + wt.Branch + ")"
			} else if wt.Detached {
				display += " (detached)"
			}
			items = append(items, Item{
				Id:      "worktree-" + string(rune('0'+i)),
				Display: display,
				Path:    wt.Path,
			})
		}
	}

	return items, nil
}
