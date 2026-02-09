package git

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func NewGit(path string) *Git {
	return &Git{
		path: path,
	}
}

type Git struct {
	path string
}

func (g *Git) commandOpts() []OptFunc {
	if g.path != "" {
		return []OptFunc{WithPath(g.path)}
	}
	return nil
}

func (g *Git) IsGitRepository() bool {
	_, err := g.GetGitDir()
	return err == nil
}

func (g *Git) IsWorktree() (bool, error) {
	_, err := g.GetGitDir()
	if err != nil {
		return false, err
	}

	gitPath := filepath.Join(g.path, ".git")
	info, err := os.Stat(gitPath)
	if err != nil {
		return false, fmt.Errorf("failed to stat .git: %w", err)
	}

	if info.IsDir() {
		return false, nil
	}

	content, err := os.ReadFile(gitPath)
	if err != nil {
		return false, fmt.Errorf("failed to read .git file: %w", err)
	}

	return strings.HasPrefix(string(content), "gitdir:"), nil
}

func (g *Git) GetGitDir() (string, error) {
	opts := append(g.commandOpts(), WithFlag("--git-dir"))
	return Output("rev-parse", opts...)
}

func (g *Git) GetWorktreeRoot() (string, error) {
	isWorktree, err := g.IsWorktree()
	if err != nil {
		return "", err
	}

	if !isWorktree {
		return "", fmt.Errorf("not a worktree")
	}

	gitPath := filepath.Join(g.path, ".git")
	content, err := os.ReadFile(gitPath)
	if err != nil {
		return "", fmt.Errorf("failed to read .git file: %w", err)
	}

	line := strings.TrimSpace(string(content))
	if !strings.HasPrefix(line, "gitdir:") {
		return "", fmt.Errorf("invalid .git file format")
	}

	gitDir := strings.TrimSpace(strings.TrimPrefix(line, "gitdir:"))

	// The gitDir points to something like /path/to/main/.git/worktrees/worktree-name
	// The commondir file is in that directory
	commonFile := filepath.Join(gitDir, "commondir")

	commonContent, err := os.ReadFile(commonFile)
	if err != nil {
		return "", fmt.Errorf("failed to read commondir file: %w", err)
	}

	commonDir := strings.TrimSpace(string(commonContent))
	if filepath.IsAbs(commonDir) {
		return filepath.Clean(commonDir), nil
	}

	return filepath.Clean(filepath.Join(gitDir, commonDir)), nil
}

func (g *Git) GetMainRepositoryPath() (string, error) {
	root, err := g.GetWorktreeRoot()
	if err != nil {
		return "", err
	}

	return filepath.Dir(root), nil
}

func (g *Git) ListWorktrees() ([]WorktreeInfo, error) {
	opts := append(g.commandOpts(), WithArgs("list", "--porcelain"))
	lines, err := OutputLines("worktree", opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to list worktrees: %w", err)
	}

	var worktrees []WorktreeInfo
	var current WorktreeInfo

	for _, line := range lines {
		if strings.HasPrefix(line, "worktree ") {
			if current.Path != "" {
				worktrees = append(worktrees, current)
			}
			current = WorktreeInfo{
				Path: strings.TrimPrefix(line, "worktree "),
			}
		} else if strings.HasPrefix(line, "HEAD ") {
			current.HEAD = strings.TrimPrefix(line, "HEAD ")
		} else if strings.HasPrefix(line, "branch ") {
			current.Branch = strings.TrimPrefix(line, "branch ")
		} else if strings.HasPrefix(line, "detached") {
			current.Detached = true
		}
	}

	if current.Path != "" {
		worktrees = append(worktrees, current)
	}

	return worktrees, nil
}

type WorktreeInfo struct {
	Path     string
	HEAD     string
	Branch   string
	Detached bool
}

func (g *Git) GetCurrentBranch() (string, error) {
	opts := append(g.commandOpts(), WithArgs("--short", "HEAD"))
	return Output("symbolic-ref", opts...)
}

func (g *Git) GetRemoteUrl(remote string) (string, error) {
	opts := append(g.commandOpts(), WithArgs("get-url", remote))
	return Output("remote", opts...)
}
