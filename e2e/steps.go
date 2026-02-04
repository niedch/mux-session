package e2e

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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

		err := executeMuxSessionWithConfig("list-sessions", "-f", testCtx.tempConfigFile, "-L", testCtx.tmuxSessionName)(ctx, docString)
		if err != nil {
			return err
		}

		return nil
	})

	ctx.Step(`^I should see the following items in output:$`, func(ctx context.Context, table *godog.Table) error {
		testCtx := ctx.Value("testCtx").(*testContext)

		for _, row := range table.Rows[1:] {
			item := row.Cells[0].Value
			assert.Contains(godog.T(ctx), testCtx.lastOutput, item, "Expected output to contain: %s", item)
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

	ctx.Step(`^window "([^"]*)" in session "([^"]*)" contains following panels:$`, func(ctx context.Context, windowName, sessionName string, table *godog.Table) error {
		testCtx := ctx.Value("testCtx").(*testContext)

		// Map headers to tmux format specifiers
		headerMap := map[string]string{
			"command": "#{pane_current_command}",
			"cwd":     "#{pane_current_path}",
		}

		var formats []string
		var headers []string

		// Validate headers and build format string
		for _, cell := range table.Rows[0].Cells {
			header := cell.Value
			if format, ok := headerMap[header]; ok {
				formats = append(formats, format)
				headers = append(headers, header)
			} else {
				return fmt.Errorf("unknown column header: %s. Supported headers: command, cwd", header)
			}
		}

		if len(formats) == 0 {
			return fmt.Errorf("no valid headers provided in table")
		}

		formatString := strings.Join(formats, "|")
		prepended_args := []string{"-L", testCtx.tmuxSessionName, "list-panes", "-t", fmt.Sprintf("%s:%s", sessionName, windowName), "-F", formatString}

		// Wait for panes to initialize
		time.Sleep(3 * time.Second)

		output, err := executeCommand("tmux", prepended_args...)
		if err != nil {
			return fmt.Errorf("failed to list panes: %v", err)
		}

		outputLines := strings.Split(strings.TrimSpace(output), "\n")
		// Remove empty lines if any
		var actualLines []string
		for _, line := range outputLines {
			if line != "" {
				actualLines = append(actualLines, line)
			}
		}

		expectedRows := table.Rows[1:]
		if len(actualLines) != len(expectedRows) {
			return fmt.Errorf("expected %d panels, but found %d", len(expectedRows), len(actualLines))
		}

		for i, row := range expectedRows {
			actualParts := strings.Split(actualLines[i], "|")

			if len(actualParts) != len(headers) {
				return fmt.Errorf("output format mismatch for panel %d", i)
			}

			for j, cell := range row.Cells {
				expectedValue := cell.Value
				actualValue := actualParts[j]

				assert.Equal(godog.T(ctx), expectedValue, actualValue, "Panel %d: expected %s '%s', got '%s'", i, headers[j], expectedValue, actualValue)
			}
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
