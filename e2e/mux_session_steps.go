package e2e

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/cucumber/godog"
)

func RegisterMuxSessionSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I run mux-session with help flag$`, executeCommandStep("./mux-session", "--help"))

	ctx.Step(`^I run mux-session list-sessions with config:$`, func(ctx context.Context, docString *godog.DocString) error {
		testCtx := ctx.Value("testCtx").(*testContext)

		err := executeMuxSessionWithConfig("list-sessions", "-f", testCtx.tempConfigFile, "-L", testCtx.tmuxSessionName)(ctx, docString)
		if err != nil {
			return err
		}

		return nil
	})

	ctx.Step(`^I run mux-session switch "([^"]*)" with config:$`, func(ctx context.Context, dirName string, docString *godog.DocString) error {
		err := executeMuxSessionWithConfig("switch", dirName)(ctx, docString)
		if err != nil {
			return err
		}

		return nil
	})
}

func executeMuxSessionWithConfig(cmd string, args ...string) func(ctx context.Context, config *godog.DocString) error {
	return func(ctx context.Context, config *godog.DocString) error {
		testCtx := ctx.Value("testCtx").(*testContext)

		if testCtx.tempDir == "" {
			return fmt.Errorf("no temp directory created, ensure 'I have the following directories' step is called first")
		}

		configContent := strings.ReplaceAll(config.Content, "<search_path>", testCtx.tempDir)

		tempConfigFile, err := os.CreateTemp("", fmt.Sprintf("mux-session-config-%s-*.toml", testCtx.tmuxSessionName))
		if err != nil {
			return fmt.Errorf("failed to create temp config file: %v", err)
		}
		testCtx.tempConfigFile = tempConfigFile.Name()

		if _, err := tempConfigFile.WriteString(configContent); err != nil {
			return fmt.Errorf("failed to write config content: %v", err)
		}
		tempConfigFile.Close()

		args = append(args, "-f", testCtx.tempConfigFile)
		args = append(args, "-L", testCtx.tmuxSessionName)

		output, err := executeCommand("./mux-session", append([]string{cmd}, args...)...)
		if err != nil {
			return err
		}

		testCtx.lastOutput = output
		return nil
	}
}
