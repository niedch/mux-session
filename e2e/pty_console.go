package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/creack/pty"
)

// PtyConsole wraps a PTY for TUI testing
type PtyConsole struct {
	pty        *os.File
	cmd        *exec.Cmd
	socket     string
	configPath string
	done       chan error
}

// NewPtyConsole creates a new PTY console for testing
func NewPtyConsole() (*PtyConsole, error) {
	return &PtyConsole{
		done: make(chan error, 1),
	}, nil
}

// Spawn starts the mux-session binary in the PTY
func (p *PtyConsole) Spawn(socket string, configPath string, searchPaths ...string) error {
	p.socket = socket

	args := []string{"-L", socket}
	if configPath != "" {
		// Use provided config file
		p.configPath = configPath
		args = append(args, "-f", p.configPath)
	} else if len(searchPaths) > 0 {
		// Use a temp config file with specific search paths
		p.configPath = createTempConfig(searchPaths)
		if p.configPath == "" {
			return fmt.Errorf("failed to create temp config")
		}
		args = append(args, "-f", p.configPath)
	}

	p.cmd = exec.Command("../bin/mux-session", args...)
	p.cmd.Env = os.Environ()

	// Start PTY
	ptmx, err := pty.Start(p.cmd)
	if err != nil {
		return fmt.Errorf("failed to start PTY: %w", err)
	}
	p.pty = ptmx

	// Set PTY size
	pty.Setsize(ptmx, &pty.Winsize{
		Rows: 24,
		Cols: 80,
	})

	// Set TERM for proper terminal handling
	p.cmd.Env = append(p.cmd.Env, "TERM=xterm-256color")

	// Monitor process completion in background
	go func() {
		p.done <- p.cmd.Wait()
	}()

	return nil
}

// Send sends raw bytes to the PTY
func (p *PtyConsole) Send(data []byte) error {
	if p.pty == nil {
		return fmt.Errorf("PTY not initialized")
	}
	_, err := p.pty.Write(data)
	return err
}

// SendString sends a string to the PTY
func (p *PtyConsole) SendString(s string) error {
	return p.Send([]byte(s))
}

// SendLine sends a string followed by Enter (newline)
func (p *PtyConsole) SendLine(s string) error {
	return p.SendString(s + "\n")
}

// SendEnter sends just the Enter key
func (p *PtyConsole) SendEnter() error {
	return p.SendString("\n")
}

// SendEscape sends the Escape key
func (p *PtyConsole) SendEscape() error {
	return p.Send([]byte{0x1b}) // Escape character
}

// SendCtrlC sends Ctrl+C
func (p *PtyConsole) SendCtrlC() error {
	return p.Send([]byte{0x03}) // Ctrl+C = ASCII 3
}

// SendArrowDown sends the Down arrow key
func (p *PtyConsole) SendArrowDown() error {
	// Arrow keys are escape sequences: ESC [ B
	return p.Send([]byte{0x1b, '[', 'B'})
}

// SendArrowUp sends the Up arrow key
func (p *PtyConsole) SendArrowUp() error {
	// Arrow keys are escape sequences: ESC [ A
	return p.Send([]byte{0x1b, '[', 'A'})
}

// Wait waits for the process to finish with timeout
func (p *PtyConsole) Wait() error {
	if p.cmd == nil {
		return nil
	}

	select {
	case err := <-p.done:
		return err
	case <-time.After(5 * time.Second):
		// Timeout - kill the process
		if p.cmd.Process != nil {
			p.cmd.Process.Kill()
		}
		return fmt.Errorf("timeout waiting for process")
	}
}

// ProcessState returns the process state if it has exited
func (p *PtyConsole) ProcessState() *os.ProcessState {
	if p.cmd == nil {
		return nil
	}
	return p.cmd.ProcessState
}

// Close closes the console and kills the process
func (p *PtyConsole) Close() error {
	// Kill process if still running
	if p.cmd != nil && p.cmd.Process != nil {
		p.cmd.Process.Kill()
		select {
		case <-p.done:
		case <-time.After(100 * time.Millisecond):
		}
	}

	// Close PTY
	if p.pty != nil {
		p.pty.Close()
	}

	// Clean up temp config file
	if p.configPath != "" {
		os.Remove(p.configPath)
	}

	return nil
}

// createTempConfig creates a temporary config file with specific search paths
func createTempConfig(searchPaths []string) string {
	// Create temp file
	tmpFile, err := os.CreateTemp("", "mux-session-config-*.toml")
	if err != nil {
		return ""
	}

	// Write minimal config with search paths
	tmpFile.WriteString("search_paths = [\n")
	for _, path := range searchPaths {
		tmpFile.WriteString(fmt.Sprintf("  %q,\n", path))
	}
	tmpFile.WriteString("]\n")
	tmpFile.Close()

	return tmpFile.Name()
}
