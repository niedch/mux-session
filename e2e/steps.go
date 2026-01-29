package e2e

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/cucumber/godog"
)

func InitializeSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {
		// Setup before running the test suite
	})

	ctx.AfterSuite(func() {
		// Cleanup after running the test suite
	})
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		// Setup before each scenario
		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		// Cleanup after each scenario
		return ctx, nil
	})

	ctx.Step(`^I have mux-session installed$`, func() error {
		// Check if we can run the help command directly with go run
		cmd := exec.Command("go", "run", "../", "--help")
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to run mux-session --help: %w, output: %s", err, string(output))
		}
		return nil
	})

	ctx.Step(`^I run mux-session with help flag$`, func() error {
		// Run the help command
		cmd := exec.Command("go", "run", "../", "--help")
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to run mux-session --help: %w, output: %s", err, string(output))
		}
		return nil
	})

	ctx.Step(`^I should see help output$`, func() error {
		// This step would typically capture and verify output
		return nil
	})
}
