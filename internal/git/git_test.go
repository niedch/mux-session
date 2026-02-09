package git

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewGit(t *testing.T) {
	g := NewGit("/some/path")
	if g == nil {
		t.Fatal("expected non-nil Git instance")
	}
	if g.path != "/some/path" {
		t.Errorf("expected path '/some/path', got '%s'", g.path)
	}
}

func TestIsGitRepository(t *testing.T) {
	tmpDir := t.TempDir()

	// Test non-git directory
	g := NewGit(tmpDir)
	if g.IsGitRepository() {
		t.Error("expected non-git directory to return false")
	}

	// Initialize a git repo
	cmd := NewCommand("init", WithPath(tmpDir))
	if _, err := cmd.Exec(); err != nil {
		t.Fatalf("failed to init git repo: %v", err)
	}

	g = NewGit(tmpDir)
	if !g.IsGitRepository() {
		t.Error("expected git directory to return true")
	}
}

func TestIsWorktree(t *testing.T) {
	// Create main repo
	mainDir := t.TempDir()
	cmd := NewCommand("init", WithPath(mainDir))
	if _, err := cmd.Exec(); err != nil {
		t.Fatalf("failed to init main git repo: %v", err)
	}

	// Configure git user for commits
	configCmd := NewCommand("config", WithArgs("user.email", "test@test.com"), WithPath(mainDir))
	configCmd.Exec()
	configCmd = NewCommand("config", WithArgs("user.name", "Test"), WithPath(mainDir))
	configCmd.Exec()

	// Create initial commit
	testFile := filepath.Join(mainDir, "test.txt")
	os.WriteFile(testFile, []byte("test"), 0644)
	addCmd := NewCommand("add", WithArgs("test.txt"), WithPath(mainDir))
	addCmd.Exec()
	commitCmd := NewCommand("commit", WithArgs("-m", "initial"), WithPath(mainDir))
	commitCmd.Exec()

	// Test main repo is not a worktree
	g := NewGit(mainDir)
	isWorktree, err := g.IsWorktree()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if isWorktree {
		t.Error("expected main repo to not be a worktree")
	}

	// Create worktree
	worktreeDir := filepath.Join(t.TempDir(), "worktree")
	worktreeCmd := NewCommand("worktree", WithArgs("add", worktreeDir), WithPath(mainDir))
	if _, err := worktreeCmd.Exec(); err != nil {
		t.Fatalf("failed to create worktree: %v", err)
	}

	// Test worktree detection
	g = NewGit(worktreeDir)
	isWorktree, err = g.IsWorktree()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !isWorktree {
		t.Error("expected worktree to be detected")
	}
}

func TestGetMainRepositoryPath(t *testing.T) {
	// Create main repo
	mainDir := t.TempDir()
	cmd := NewCommand("init", WithPath(mainDir))
	if _, err := cmd.Exec(); err != nil {
		t.Fatalf("failed to init main git repo: %v", err)
	}

	// Configure git user
	configCmd := NewCommand("config", WithArgs("user.email", "test@test.com"), WithPath(mainDir))
	configCmd.Exec()
	configCmd = NewCommand("config", WithArgs("user.name", "Test"), WithPath(mainDir))
	configCmd.Exec()

	// Create initial commit
	testFile := filepath.Join(mainDir, "test.txt")
	os.WriteFile(testFile, []byte("test"), 0644)
	addCmd := NewCommand("add", WithArgs("test.txt"), WithPath(mainDir))
	addCmd.Exec()
	commitCmd := NewCommand("commit", WithArgs("-m", "initial"), WithPath(mainDir))
	commitCmd.Exec()

	// Create worktree
	worktreeDir := filepath.Join(t.TempDir(), "worktree")
	worktreeCmd := NewCommand("worktree", WithArgs("add", worktreeDir), WithPath(mainDir))
	if _, err := worktreeCmd.Exec(); err != nil {
		t.Fatalf("failed to create worktree: %v", err)
	}

	// Test getting main repository path
	g := NewGit(worktreeDir)
	mainPath, err := g.GetMainRepositoryPath()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Normalize paths for comparison
	mainDirAbs, _ := filepath.Abs(mainDir)
	mainPathAbs, _ := filepath.Abs(mainPath)
	if mainDirAbs != mainPathAbs {
		t.Errorf("expected main path '%s', got '%s'", mainDirAbs, mainPathAbs)
	}
}

func TestListWorktrees(t *testing.T) {
	// Create main repo
	mainDir := t.TempDir()
	cmd := NewCommand("init", WithPath(mainDir))
	if _, err := cmd.Exec(); err != nil {
		t.Fatalf("failed to init main git repo: %v", err)
	}

	// Configure git user
	configCmd := NewCommand("config", WithArgs("user.email", "test@test.com"), WithPath(mainDir))
	configCmd.Exec()
	configCmd = NewCommand("config", WithArgs("user.name", "Test"), WithPath(mainDir))
	configCmd.Exec()

	// Create initial commit
	testFile := filepath.Join(mainDir, "test.txt")
	os.WriteFile(testFile, []byte("test"), 0644)
	addCmd := NewCommand("add", WithArgs("test.txt"), WithPath(mainDir))
	addCmd.Exec()
	commitCmd := NewCommand("commit", WithArgs("-m", "initial"), WithPath(mainDir))
	commitCmd.Exec()

	// Create worktrees
	worktree1Dir := filepath.Join(t.TempDir(), "worktree1")
	worktreeCmd := NewCommand("worktree", WithArgs("add", worktree1Dir), WithPath(mainDir))
	if _, err := worktreeCmd.Exec(); err != nil {
		t.Fatalf("failed to create worktree: %v", err)
	}

	worktree2Dir := filepath.Join(t.TempDir(), "worktree2")
	worktreeCmd = NewCommand("worktree", WithArgs("add", "-b", "feature", worktree2Dir), WithPath(mainDir))
	if _, err := worktreeCmd.Exec(); err != nil {
		t.Fatalf("failed to create worktree: %v", err)
	}

	// Test listing worktrees
	g := NewGit(mainDir)
	worktrees, err := g.ListWorktrees()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(worktrees) < 3 {
		t.Errorf("expected at least 3 worktrees (main + 2), got %d", len(worktrees))
	}

	// Check that main directory is in the list
	foundMain := false
	for _, wt := range worktrees {
		if wt.Path == mainDir {
			foundMain = true
			break
		}
	}
	if !foundMain {
		t.Error("expected main directory to be in worktree list")
	}
}

func TestGetCurrentBranch(t *testing.T) {
	mainDir := t.TempDir()
	cmd := NewCommand("init", WithPath(mainDir))
	if _, err := cmd.Exec(); err != nil {
		t.Fatalf("failed to init git repo: %v", err)
	}

	g := NewGit(mainDir)
	branch, err := g.GetCurrentBranch()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if branch == "" {
		t.Error("expected non-empty branch name")
	}
}
