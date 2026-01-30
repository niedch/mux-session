package e2e

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"
)

type testContext struct {
	lastOutput      string
	tmuxSessionName string
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		testCtx := &testContext{lastOutput: "", tmuxSessionName: fmt.Sprintf("test-%s", sc.Id)}

		return context.WithValue(ctx, "testCtx", testCtx), nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		executeTmuxCommand("tmux", "kill-server")(ctx)
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

	ctx.Step(`^I run list-sessions$`, executeTmuxCommand("tmux", "list-sessions"))

	ctx.Step(`^I build the mux-session$`, executeCommandStep("go", "build", "../"))
	ctx.Step(`^I run mux-session with help flag$`, executeCommandStep("mux-session", "--help"))

	ctx.Step(`^I expect following Sessions:$`, func(ctx context.Context, docString *godog.DocString) error {
		testCtx := ctx.Value("testCtx").(*testContext)
		if err := executeTmuxCommand("tmux", "list-sessions")(ctx); err != nil {
			return err
		}

		expected_sessions := strings.Split(docString.Content, "\n")
		for _, expected_session := range expected_sessions {
			assert.Contains(godog.T(ctx), testCtx.lastOutput, expected_session)
		}

		return nil
	})

	ctx.Step(`^I expect following sessions:$`, func(ctx context.Context) error {
		testCtx := ctx.Value("testCtx").(*testContext)
		if err := executeTmuxCommand("tmux", "list-sessions")(ctx); err != nil {
			return err
		}

		log.Printf("Actual sessions: %s", testCtx.lastOutput)
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
