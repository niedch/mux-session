package e2e

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"
)

func RegisterMuxSessionSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I build the mux-session$`, executeCommandStep("go", "build", "-o", "mux-session", ".."))
	ctx.Step(`^I run mux-session with help flag$`, executeCommandStep("./mux-session", "--help"))

	ctx.Step(`^I should see help output$`, func(ctx context.Context) error {
		testCtx := ctx.Value("testCtx").(*testContext)
		if testCtx.lastOutput == "" {
			return fmt.Errorf("no output captured")
		}

		godog.Logf(ctx, "Verifying output: %s\n", testCtx.lastOutput)
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
		err := executeMuxSessionWithConfig("list-sessions")(ctx, docString)
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

	ctx.Step(`^I run mux-session config-validate with config:$`, func(ctx context.Context, docString *godog.DocString) error {
		testCtx := ctx.Value("testCtx").(*testContext)
		err := executeMuxSessionWithConfig("config-validate")(ctx, docString)
		if err != nil {
			godog.Logf(ctx, "Error when executing '%s'\n", testCtx.lastOutput)
		}

		return nil
	})

	ctx.Step(`^I should see the following (lines?|items?) in output:$`, func(ctx context.Context, _ string, table *godog.Table) error {
		testCtx := ctx.Value("testCtx").(*testContext)

		for _, row := range table.Rows[1:] {
			pattern := row.Cells[0].Value
			re, err := regexp.Compile(pattern)
			if err != nil {
				return fmt.Errorf("invalid regex pattern '%s': %w", pattern, err)
			}

			assert.True(godog.T(ctx), re.MatchString(testCtx.lastOutput), "Expected output '%s' to match regex: '%s'", testCtx.lastOutput, pattern)
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
		if cmd != "config-validate" {
			args = append(args, "-L", testCtx.tmuxSessionName)
		}

		output, err := executeCommand("./mux-session", append([]string{cmd}, args...)...)
		if err != nil {
			testCtx.lastOutput = output
			godog.Logf(ctx, "Error when '%s' '%s'\n", cmd, testCtx.lastOutput)
			return err
		}

		testCtx.lastOutput = output
		return nil
	}
}
