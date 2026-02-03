package e2e

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/creack/pty"
)

type PtyConsole struct {
	pty        *os.File
	cmd        *exec.Cmd
	socket     string
	configPath string
	binaryPath string
	done       chan error
	output     bytes.Buffer
}

func NewPtyConsole() (*PtyConsole, error) {
	return &PtyConsole{
		done: make(chan error, 1),
	}, nil
}

func WithBinaryPath(binaryPath string) func(*PtyConsole) {
	return func(p *PtyConsole) {
		p.binaryPath = binaryPath
	}
}

func (p *PtyConsole) Spawn(socket string, configPath string) error {
	p.socket = socket

	args := []string{"-L", socket}
	if configPath != "" {
		p.configPath = configPath
		args = append(args, "-f", p.configPath)
	}

	binaryPath := "../mux-session"
	if p.binaryPath != "" {
		binaryPath = p.binaryPath
	}

	p.cmd = exec.Command(binaryPath, args...)
	p.cmd.Env = os.Environ()
	p.cmd.Env = append(p.cmd.Env, "TERM=xterm-256color")

	ptmx, err := pty.Start(p.cmd)
	if err != nil {
		return fmt.Errorf("failed to start PTY: %w", err)
	}
	p.pty = ptmx

	pty.Setsize(ptmx, &pty.Winsize{
		Rows: 24,
		Cols: 80,
	})

	p.cmd.Env = append(p.cmd.Env, "TERM=xterm-256color")

	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := p.pty.Read(buf)
			if n > 0 {
				data := buf[:n]
				p.output.Write(data)

				// Handle cursor position query
				if bytes.Contains(data, []byte("\x1b[6n")) {
					p.pty.Write([]byte("\x1b[1;1R"))
				}
			}
			if err != nil {
				if err != io.EOF {
					// Log error if needed, or just exit loop
				}
				break
			}
		}
	}()

	go func() {
		p.done <- p.cmd.Wait()
	}()

	return nil
}

func (p *PtyConsole) Send(data []byte) error {
	if p.pty == nil {
		return fmt.Errorf("PTY not initialized")
	}

	_, err := p.pty.Write(data)
	return err
}

func (p *PtyConsole) SendString(s string) error {
	return p.Send([]byte(s))
}

func (p *PtyConsole) SendLine(s string) error {
	return p.SendString(s + "\n")
}

func (p *PtyConsole) SendEnter() error {
	return p.SendString("\n")
}

func (p *PtyConsole) SendEscape() error {
	return p.Send([]byte{0x1b})
}

func (p *PtyConsole) SendCtrlC() error {
	return p.Send([]byte{0x03})
}

func (p *PtyConsole) SendArrowDown() error {
	return p.Send([]byte{0x1b, '[', 'B'})
}

func (p *PtyConsole) SendBackSpace() error {
	return p.Send([]byte{0x7f})
}

func (p *PtyConsole) SendArrowUp() error {
	return p.Send([]byte{0x1b, '[', 'A'})
}

func (p *PtyConsole) Wait() error {
	if p.cmd == nil {
		return nil
	}

	select {
	case err := <-p.done:
		return err
	case <-time.After(3 * time.Second):
		if p.cmd.Process != nil {
			p.cmd.Process.Kill()
		}

		return fmt.Errorf("timeout waiting for process.")
	}
}

func (p *PtyConsole) ProcessState() *os.ProcessState {
	if p.cmd == nil {
		return nil
	}
	return p.cmd.ProcessState
}

func (p *PtyConsole) Close() error {
	if p.cmd != nil && p.cmd.Process != nil {
		p.cmd.Process.Kill()
		select {
		case <-p.done:
		case <-time.After(100 * time.Millisecond):
		}
	}

	if p.pty != nil {
		p.pty.Close()
	}

	if p.configPath != "" {
		os.Remove(p.configPath)
	}

	return nil
}
