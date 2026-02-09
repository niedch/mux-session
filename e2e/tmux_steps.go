package e2e

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"
)

func RegisterTmuxSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^a new tmux server$`, func(ctx context.Context) error {
		if err := executeTmuxCommand("tmux", "start-server")(ctx); err != nil {
			return err
		}

		if err := executeTmuxCommand("tmux", "new-session", "-d", "-s", "test-session")(ctx); err != nil {
			return err
		}

		// Add sleep to wait for initialization of tmux server
		time.Sleep(100 * time.Millisecond)

		return nil
	})

	ctx.Step(`^window "([^"]*)" in session "([^"]*)" is split (vertically|horizontally)$`, func(ctx context.Context, windowName, sessionName, splitType string) error {
		testCtx := ctx.Value("testCtx").(*testContext)

		// Get pane coordinates: index, left, top
		args := []string{"-L", testCtx.tmuxSessionName, "list-panes", "-t", fmt.Sprintf("%s:%s", sessionName, windowName), "-F", "#{pane_index}:#{pane_left},#{pane_top}"}
		output, err := executeCommand("tmux", args...)
		if err != nil {
			return fmt.Errorf("failed to list panes: %v", err)
		}

		lines := strings.Split(strings.TrimSpace(output), "\n")
		var actualLines []string
		for _, line := range lines {
			if line != "" {
				actualLines = append(actualLines, line)
			}
		}

		if len(actualLines) < 2 {
			return fmt.Errorf("expected at least 2 panes to determine split, found %d", len(actualLines))
		}

		// Parse coordinates of first two panes
		parseCoords := func(line string) (int, int, error) {
			parts := strings.Split(line, ":")
			if len(parts) != 2 {
				return 0, 0, fmt.Errorf("invalid format")
			}
			coords := strings.Split(parts[1], ",")
			if len(coords) != 2 {
				return 0, 0, fmt.Errorf("invalid coords")
			}
			x, err := strconv.Atoi(coords[0])
			if err != nil {
				return 0, 0, err
			}
			y, err := strconv.Atoi(coords[1])
			if err != nil {
				return 0, 0, err
			}
			return x, y, nil
		}

		x0, y0, err := parseCoords(actualLines[0])
		if err != nil {
			return fmt.Errorf("failed to parse pane 0 coords: %v", err)
		}
		x1, y1, err := parseCoords(actualLines[1])
		if err != nil {
			return fmt.Errorf("failed to parse pane 1 coords: %v", err)
		}

		// Check orientation based on coordinates
		// Vertical Split (Side-by-Side): Same Top (Y), Different Left (X)
		// Horizontal Split (Top-Bottom): Same Left (X), Different Top (Y)

		isVertical := y0 == y1 && x0 != x1
		isHorizontal := x0 == x1 && y0 != y1

		switch splitType {
		case "vertically":
			if !isVertical {
				return fmt.Errorf("expected vertical split (side-by-side), but found panes at P0(%d,%d) and P1(%d,%d)", x0, y0, x1, y1)
			}
		case "horizontally":
			if !isHorizontal {
				return fmt.Errorf("expected horizontal split (top-bottom), but found panes at P0(%d,%d) and P1(%d,%d)", x0, y0, x1, y1)
			}
		default:
			return fmt.Errorf("unknown split type: %s", splitType)
		}

		return nil
	})

	ctx.Step(`^I run list-sessions$`, executeTmuxCommand("tmux", "list-sessions"))

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

	ctx.Step(`^session "([^"]*)" contains following windows:$`, func(ctx context.Context, sessionName string, table *godog.Table) error {
		testCtx := ctx.Value("testCtx").(*testContext)

		prepended_args := []string{"-L", testCtx.tmuxSessionName, "list-windows", "-t", sessionName}
		output, err := executeCommand("tmux", prepended_args...)
		if err != nil {
			return err
		}

		for _, row := range table.Rows[1:] {
			windowName := row.Cells[0].Value
			assert.Contains(godog.T(ctx), output, windowName, "Expected session %s to contain window: %s", sessionName, windowName)
		}

		return nil
	})

	ctx.Step(`^window "([^"]*)" in session "([^"]*)" has environment variable "([^"]*)" set to "([^"]*)"$`, func(ctx context.Context, windowName, sessionName, envVar, expectedValue string) error {
		testCtx := ctx.Value("testCtx").(*testContext)

		// Use display-message to retrieve the environment variable value
		cmd := executeTmuxCommand("tmux", "display-message", "-p", "-t", fmt.Sprintf("%s:%s", sessionName, windowName), fmt.Sprintf("#{%s}", envVar))
		if err := cmd(ctx); err != nil {
			return fmt.Errorf("failed to get environment variable: %v", err)
		}

		actualValue := strings.TrimSpace(testCtx.lastOutput)
		if actualValue != expectedValue {
			return fmt.Errorf("expected env var %s to be '%s', but got '%s'", envVar, expectedValue, actualValue)
		}

		return nil
	})

	ctx.Step(`^window "([^"]*)" in session "([^"]*)" is active$`, func(ctx context.Context, windowName, sessionName string) error {
		testCtx := ctx.Value("testCtx").(*testContext)

		cmd := executeTmuxCommand("tmux", "list-windows", "-t", sessionName, "-F", "#{window_name}:#{window_active}")
		if err := cmd(ctx); err != nil {
			return fmt.Errorf("failed to list windows: %v", err)
		}

		lines := strings.Split(strings.TrimSpace(testCtx.lastOutput), "\n")
		found := false
		isActive := false

		for _, line := range lines {
			parts := strings.Split(line, ":")
			if len(parts) != 2 {
				continue
			}
			if parts[0] == windowName {
				found = true
				if parts[1] == "1" {
					isActive = true
				}
				break
			}
		}

		if !found {
			return fmt.Errorf("window '%s' not found in session '%s'", windowName, sessionName)
		}
		if !isActive {
			return fmt.Errorf("window '%s' is not active", windowName)
		}

		return nil
	})

	ctx.Step(`^I execute following Command in Session "([^"]*)" on Window "([^"]*)":$`, func(ctx context.Context, sessionName, windowName string, command *godog.DocString) error {
		cmdStr := strings.TrimSpace(command.Content)
		target := fmt.Sprintf("%s:%s", sessionName, windowName)

		// Wait for Session to initialized
		time.Sleep(2 * time.Second)

		// Send the command and Enter
		cmd := executeTmuxCommand("tmux", "send-keys", "-t", target, cmdStr, "C-m")
		if err := cmd(ctx); err != nil {
			return fmt.Errorf("failed to send command: %v", err)
		}

		// Wait for command execution
		time.Sleep(2 * time.Second)

		// Capture pane content to allow verification
		captureCmd := executeTmuxCommand("tmux", "capture-pane", "-p", "-t", target)
		if err := captureCmd(ctx); err != nil {
			return fmt.Errorf("failed to capture pane output: %v", err)
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
