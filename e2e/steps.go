package e2e

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"
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
			log.Printf("Failed to remove socket file %s: %v", socketPath, err)
		}

		if testCtx.tempConfigFile != "" {
			if err := os.Remove(testCtx.tempConfigFile); err != nil && !os.IsNotExist(err) {
				log.Printf("Failed to remove temp config file %s: %v", testCtx.tempConfigFile, err)
			}
		}

		if testCtx.tempDir != "" {
			if err := os.RemoveAll(testCtx.tempDir); err != nil {
				log.Printf("Failed to remove temp directory %s: %v", testCtx.tempDir, err)
			}
		}

		return ctx, nil
	})

	ctx.Step(`^a new tmux server$`, func(ctx context.Context) error {
		if err := executeTmuxCommand("tmux", "start-server")(ctx); err != nil {
			return err
		}

		if err := executeTmuxCommand("tmux", "new-session", "-d", "-s", "test-session")(ctx); err != nil {
			return err
		}

		return nil
	})

	ctx.Step(`^I build the mux-session$`, executeCommandStep("go", "build", "-o", "mux-session", ".."))
	ctx.Step(`^I run list-sessions$`, executeTmuxCommand("tmux", "list-sessions"))
	ctx.Step(`^I run mux-session with help flag$`, executeCommandStep("./mux-session", "--help"))

	ctx.Step(`^I expect following sessions:$`, func(ctx context.Context, docString *godog.DocString) error {
		testCtx := ctx.Value("testCtx").(*testContext)
		if err := executeTmuxCommand("tmux", "list-sessions")(ctx); err != nil {
			return err
		}

		for expected_session := range strings.SplitSeq(docString.Content, "\n") {
			assert.Contains(godog.T(ctx), testCtx.lastOutput, expected_session)
		}

		return nil
	})

	ctx.Step(`^I should see help output$`, func(ctx context.Context) error {
		testCtx := ctx.Value("testCtx").(*testContext)
		if testCtx.lastOutput == "" {
			return fmt.Errorf("no output captured")
		}

		fmt.Printf("Verifying output: %s\n", testCtx.lastOutput) // Print for debugging
		return nil
	})

	ctx.Step(`^I have the following directories:$`, func(ctx context.Context, table *godog.Table) error {
		testCtx := ctx.Value("testCtx").(*testContext)

		tempDir, err := os.MkdirTemp("", fmt.Sprintf("mux-session-test-%s-*", testCtx.tmuxSessionName))
		if err != nil {
			return fmt.Errorf("failed to create temp directory: %v", err)
		}
		testCtx.tempDir = tempDir

		for _, row := range table.Rows[1:] {
			dirName := row.Cells[0].Value
			dirPath := filepath.Join(tempDir, dirName)
			if err := os.MkdirAll(dirPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %v", dirPath, err)
			}
		}

		return nil
	})

	ctx.Step(`^I run mux-session list-sessions with config:$`, func(ctx context.Context, docString *godog.DocString) error {
		testCtx := ctx.Value("testCtx").(*testContext)

		if testCtx.tempDir == "" {
			return fmt.Errorf("no temp directory created, ensure 'I have the following directories' step is called first")
		}

		configContent := strings.ReplaceAll(docString.Content, "<search_path>", testCtx.tempDir)

		tempConfigFile, err := os.CreateTemp("", fmt.Sprintf("mux-session-config-%s-*.toml", testCtx.tmuxSessionName))
		if err != nil {
			return fmt.Errorf("failed to create temp config file: %v", err)
		}
		testCtx.tempConfigFile = tempConfigFile.Name()

		if _, err := tempConfigFile.WriteString(configContent); err != nil {
			return fmt.Errorf("failed to write config content: %v", err)
		}
		tempConfigFile.Close()

		output, err := executeCommand("./mux-session", "list-sessions", "-f", testCtx.tempConfigFile, "-L", testCtx.tmuxSessionName)
		if err != nil {
			return err
		}

		testCtx.lastOutput = output
		return nil
	})

	ctx.Step(`^I should see the following directories in output:$`, func(ctx context.Context, table *godog.Table) error {
		testCtx := ctx.Value("testCtx").(*testContext)

		for _, row := range table.Rows[1:] {
			dirName := row.Cells[0].Value
			assert.Contains(godog.T(ctx), testCtx.lastOutput, dirName, "Expected output to contain directory: %s", dirName)
		}

		return nil
	})

}

func executeTmuxCommand(name string, args ...string) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		testCtx := ctx.Value("testCtx").(*testContext)
		prepended_args := append([]string{"-L", testCtx.tmuxSessionName}, args...)

		output, err := executeCommand(name, prepended_args...)
		if err != nil {
			return err
		}

		testCtx.lastOutput = output
		return nil
	}
}

func executeCommandStep(name string, args ...string) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		testCtx := ctx.Value("testCtx").(*testContext)

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
	log.Printf("Executing Command '%s %s'\n", name, comb_args)

	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run %s %s, %s output: %s", name, comb_args, err, string(output))
	}
	return string(output), nil
}
