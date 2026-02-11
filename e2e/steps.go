package e2e

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cucumber/godog"
)

type testContext struct {
	lastOutput      string
	tmuxSessionName string
	tempDir         string
	tempConfigFile  string
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		testCtx := &testContext{lastOutput: "", tmuxSessionName: fmt.Sprintf("test-%s", sc.Id)}

		return context.WithValue(ctx, "testCtx", testCtx), nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		testCtx := ctx.Value("testCtx").(*testContext)
		executeTmuxCommand("tmux", "kill-server")(ctx)

		socketPath := filepath.Join("/tmp", "tmux-1000", testCtx.tmuxSessionName)
		if err := os.Remove(socketPath); err != nil && !os.IsNotExist(err) {
			godog.Logf(ctx, "Failed to remove socket file %s: %v", socketPath, err)
		}

		if testCtx.tempConfigFile != "" {
			if err := os.Remove(testCtx.tempConfigFile); err != nil && !os.IsNotExist(err) {
				godog.Logf(ctx, "Failed to remove temp config file %s: %v", testCtx.tempConfigFile, err)
			}
		}

		if testCtx.tempDir != "" {
			if err := os.RemoveAll(testCtx.tempDir); err != nil {
				godog.Logf(ctx, "Failed to remove temp directory %s: %v", testCtx.tempDir, err)
			}
		}

		return ctx, nil
	})

	RegisterMuxSessionSteps(ctx)
	RegisterTmuxSteps(ctx)
	RegisterGitSteps(ctx)
}

func executeCommandStep(name string, args ...string) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		testCtx := ctx.Value("testCtx").(*testContext)

		comb_args := strings.Join(args, " ")
		godog.Logf(ctx, "Executing Command '%s %s'\n", name, comb_args)

		output, err := executeCommand(name, args...)
		if err != nil {
			return err
		}

		testCtx.lastOutput = output
		return nil
	}
}

func executeCommand(name string, args ...string) (string, error) {
	comb_args := strings.Join(args, " ")

	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("failed to run %s %s, %s output: %s", name, comb_args, err, string(output))
	}
	return string(output), nil
}
