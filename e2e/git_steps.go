package e2e

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/cucumber/godog"
)

func RegisterGitSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^the directory "([^"]*)" is a git worktree$`, func(ctx context.Context, worktreeDirName string) error {
		testCtx := ctx.Value("testCtx").(*testContext)

		// Create a hidden main repository that will own the worktree
		mainRepoDir := filepath.Join(testCtx.tempDir, ".main-repo")
		if err := os.MkdirAll(mainRepoDir, 0755); err != nil {
			return fmt.Errorf("failed to create main repo directory: %v", err)
		}

		// The worktree directory that will be discovered
		worktreeDir := filepath.Join(testCtx.tempDir, worktreeDirName)

		// Initialize git repo in the hidden main directory
		if err := execGitCommand(mainRepoDir, "init"); err != nil {
			return fmt.Errorf("failed to init git repo: %v", err)
		}

		// Configure git user
		if err := execGitCommand(mainRepoDir, "config", "user.email", "test@test.com"); err != nil {
			return fmt.Errorf("failed to configure git email: %v", err)
		}
		if err := execGitCommand(mainRepoDir, "config", "user.name", "Test"); err != nil {
			return fmt.Errorf("failed to configure git name: %v", err)
		}

		// Create initial commit
		testFile := filepath.Join(mainRepoDir, "test.txt")
		if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
			return fmt.Errorf("failed to create test file: %v", err)
		}
		if err := execGitCommand(mainRepoDir, "add", "test.txt"); err != nil {
			return fmt.Errorf("failed to add test file: %v", err)
		}
		if err := execGitCommand(mainRepoDir, "commit", "-m", "initial"); err != nil {
			return fmt.Errorf("failed to commit: %v", err)
		}

		// Create the worktree at the target directory
		if err := execGitCommand(mainRepoDir, "worktree", "add", worktreeDir); err != nil {
			return fmt.Errorf("failed to create worktree: %v", err)
		}

		return nil
	})
}

func execGitCommand(dir string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git command failed: %v, output: %s", err, string(output))
	}
	return nil
}
