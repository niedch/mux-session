package e2e

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"
)

type testContext struct {
	lastOutput      string
	tmuxSessionName string
	ptyConsole      *PtyConsole
	testDirs        []string
	configPath      string
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
	ctx.Step(`^I run mux-session with help flag$`, executeCommandStep("./mux-session", "--help"))

	ctx.Step(`^I expect following Sessions:$`, func(ctx context.Context, docString *godog.DocString) error {
		testCtx := ctx.Value("testCtx").(*testContext)
		if err := executeTmuxCommand("tmux", "list-sessions")(ctx); err != nil {
			return err
		}

		for expected_session := range strings.SplitSeq(docString.Content, "\n") {
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

	// TUI test steps
	ctx.Step(`^test directories exist:$`, func(ctx context.Context, table *godog.Table) error {
		testCtx := ctx.Value("testCtx").(*testContext)

		// Create base test directory
		baseDir := filepath.Join("/tmp", fmt.Sprintf("mux-test-%s", testCtx.tmuxSessionName))
		if err := os.MkdirAll(baseDir, 0755); err != nil {
			return fmt.Errorf("failed to create test base dir: %w", err)
		}

		// Create subdirectories from table
		for _, row := range table.Rows[1:] { // Skip header
			dirName := row.Cells[0].Value
			dirPath := filepath.Join(baseDir, dirName)
			if err := os.MkdirAll(dirPath, 0755); err != nil {
				return fmt.Errorf("failed to create test dir %s: %w", dirName, err)
			}
			testCtx.testDirs = append(testCtx.testDirs, dirPath)
		}

		return nil
	})

	ctx.Step(`^I spawn mux-session with socket "([^"]*)" using config:$`, func(ctx context.Context, socket string, docString *godog.DocString) error {
		testCtx := ctx.Value("testCtx").(*testContext)

		// Resolve template variables
		configContent := docString.Content
		if len(testCtx.testDirs) > 0 {
			testDir := filepath.Dir(testCtx.testDirs[0])
			configContent = strings.ReplaceAll(configContent, "{{TEST_DIR}}", testDir)
		}

		// Create temp config file with resolved content
		tmpFile, err := os.CreateTemp("", "mux-session-config-*.toml")
		if err != nil {
			return fmt.Errorf("failed to create temp config file: %w", err)
		}
		defer tmpFile.Close()

		if _, err := tmpFile.WriteString(configContent); err != nil {
			return fmt.Errorf("failed to write config content: %w", err)
		}

		testCtx.configPath = tmpFile.Name()

		// Create PTY console
		ptyConsole, err := NewPtyConsole()
		if err != nil {
			return fmt.Errorf("failed to create PTY console: %w", err)
		}
		testCtx.ptyConsole = ptyConsole

		// Spawn mux-session with socket
		if err := ptyConsole.Spawn(socket, testCtx.configPath); err != nil {
			return fmt.Errorf("failed to spawn mux-session: %w", err)
		}

		// Wait for fzf to appear and render
		time.Sleep(1 * time.Second)

		// Send a "wakeup" keystroke to activate the terminal input
		// This is a workaround for bubbletea issue #1116 where first keypress is ignored
		if err := ptyConsole.SendString(" "); err != nil {
			return fmt.Errorf("failed to send wakeup keystroke: %w", err)
		}
		time.Sleep(100 * time.Millisecond)
		// Send backspace to remove the space
		if err := ptyConsole.Send([]byte{0x7f}); err != nil {
			return fmt.Errorf("failed to send backspace: %w", err)
		}
		time.Sleep(100 * time.Millisecond)

		return nil
	})

	ctx.Step(`^I type "([^"]*)" in fzf$`, func(ctx context.Context, query string) error {
		testCtx := ctx.Value("testCtx").(*testContext)
		if testCtx.ptyConsole == nil {
			return fmt.Errorf("PTY console not initialized")
		}

		if err := testCtx.ptyConsole.SendString(query); err != nil {
			return err
		}

		// Small delay to let fzf process the input
		time.Sleep(200 * time.Millisecond)
		return nil
	})

	ctx.Step(`^I press Enter$`, func(ctx context.Context) error {
		testCtx := ctx.Value("testCtx").(*testContext)
		if testCtx.ptyConsole == nil {
			return fmt.Errorf("PTY console not initialized")
		}

		return testCtx.ptyConsole.SendEnter()
	})

	ctx.Step(`^I press Escape$`, func(ctx context.Context) error {
		testCtx := ctx.Value("testCtx").(*testContext)
		if testCtx.ptyConsole == nil {
			return fmt.Errorf("PTY console not initialized")
		}

		return testCtx.ptyConsole.SendEscape()
	})

	ctx.Step(`^I press Ctrl\+C$`, func(ctx context.Context) error {
		testCtx := ctx.Value("testCtx").(*testContext)
		if testCtx.ptyConsole == nil {
			return fmt.Errorf("PTY console not initialized")
		}

		return testCtx.ptyConsole.SendCtrlC()
	})

	ctx.Step(`^I press Down$`, func(ctx context.Context) error {
		testCtx := ctx.Value("testCtx").(*testContext)
		if testCtx.ptyConsole == nil {
			return fmt.Errorf("PTY console not initialized")
		}

		return testCtx.ptyConsole.SendArrowDown()
	})

	ctx.Step(`^I manually select "([^"]*)" in fzf$`, func(ctx context.Context, selection string) error {
		return fmt.Errorf("manual interaction required: please select %q in the fzf UI", selection)
	})

	ctx.Step(`^a tmux session "([^"]*)" should exist on socket "([^"]*)"$`, func(ctx context.Context, sessionName, socket string) error {
		testCtx := ctx.Value("testCtx").(*testContext)

		// Wait for mux-session to complete and session to be created
		if testCtx.ptyConsole != nil {
			err := testCtx.ptyConsole.Wait()
			if err != nil {
				log.Printf("mux-session process exited with error: %v", err)
			} else {
				log.Printf("mux-session process completed successfully")
			}
		}

		time.Sleep(500 * time.Millisecond)

		// Check if session exists
		output, err := executeCommand("tmux", "-L", socket, "list-sessions")
		if err != nil {
			return fmt.Errorf("failed to list sessions: %w", err)
		}

		if !strings.Contains(output, sessionName) {
			return fmt.Errorf("session %q not found in output: %s", sessionName, output)
		}

		return nil
	})

	ctx.Step(`^session "([^"]*)" on socket "([^"]*)" should have windows:$`, func(ctx context.Context, sessionName, socket string, table *godog.Table) error {
		// Get list of windows in the session
		output, err := executeCommand("tmux", "-L", socket, "list-windows", "-t", sessionName, "-F", "#W")
		if err != nil {
			return fmt.Errorf("failed to list windows for session %s: %w", sessionName, err)
		}

		// Parse expected windows from table
		expectedWindows := make([]string, 0)
		for _, row := range table.Rows[1:] { // Skip header
			expectedWindows = append(expectedWindows, row.Cells[0].Value)
		}

		// Parse actual windows from output
		actualWindows := strings.Split(strings.TrimSpace(output), "\n")
		for i := range actualWindows {
			actualWindows[i] = strings.TrimSpace(actualWindows[i])
		}

		// Verify all expected windows exist
		for _, expected := range expectedWindows {
			found := false
			for _, actual := range actualWindows {
				if actual == expected {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("expected window %q not found in session %s. Actual windows: %v", expected, sessionName, actualWindows)
			}
		}

		return nil
	})

	ctx.Step(`^no tmux session should exist on socket "([^"]*)"$`, func(ctx context.Context, socket string) error {
		// Wait a moment
		time.Sleep(200 * time.Millisecond)

		// Try to list sessions - should fail if none exist
		_, err := executeCommand("tmux", "-L", socket, "list-sessions")
		if err == nil {
			return fmt.Errorf("expected no sessions but some exist")
		}

		return nil
	})

	// Cleanup PTY console
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		testCtx := ctx.Value("testCtx").(*testContext)

		if testCtx.ptyConsole != nil {
			testCtx.ptyConsole.Close()
		}

		// Clean up test directories
		if len(testCtx.testDirs) > 0 {
			baseDir := filepath.Dir(testCtx.testDirs[0])
			os.RemoveAll(baseDir)
		}

		return ctx, nil
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
