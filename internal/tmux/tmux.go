package tmux

import (
	"fmt"
	"os/exec"
	"slices"
	"strings"

	"github.com/niedch/mux-session/internal/conf"
)

func NewTmux() (*Tmux, error) {
	return &Tmux{}, nil
}

type Tmux struct{}

func (t *Tmux) ListSessions() ([]string, error) {
	cmd := exec.Command("tmux", "list-sessions", "-F", "#S")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	result := make([]string, 0)
	for line := range strings.SplitSeq(strings.TrimSpace(string(output)), "\n") {
		if line != "" {
			result = append(result, line)
		}
	}

	return result, nil
}

func (t *Tmux) CurrentSession() (string, error) {
	cmd := exec.Command("tmux", "display-message", "-p", "#S")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

func (t *Tmux) SwitchSession(session_name string) error {
	current_session, err := t.CurrentSession()
	if err != nil {
		return err
	}

	// No need to switch Sessions
	if current_session == session_name {
		return nil
	}

	sessions, err := t.ListSessions()
	if err != nil {
		return err
	}

	if slices.Contains(sessions, session_name) {
		cmd := exec.Command("tmux", "switch-client", "-t", session_name)
		return cmd.Run()
	}

	return fmt.Errorf("Was not able to find session '%s'", session_name)
}

func (t *Tmux) CreateSession(dir_name string, projectConfig conf.ProjectConfig) error {
	sessionName := dir_name
	if projectConfig.Name != nil {
		sessionName = *projectConfig.Name
	}

	if len(projectConfig.WindowConfig) == 0 {
		return fmt.Errorf("no window configuration found for session '%s'", sessionName)
	}

	// Create new session with first window
	firstWindow := projectConfig.WindowConfig[0]
	cmd := exec.Command("tmux", "new-session", "-d", "-s", sessionName, "-n", firstWindow.WindowName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create session '%s': %w", sessionName, err)
	}

	// Create additional windows
	for _, window := range projectConfig.WindowConfig[1:] {
		cmd := exec.Command("tmux", "new-window", "-t", sessionName, "-n", window.WindowName)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to create window '%s' in session '%s': %w", window.WindowName, sessionName, err)
		}

		// Execute command if specified
		if window.Cmd != nil && *window.Cmd != "" {
			cmd := exec.Command("tmux", "send-keys", "-t", fmt.Sprintf("%s:%s", sessionName, window.WindowName), *window.Cmd, "C-m")
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to send command to window '%s': %w", window.WindowName, err)
			}
		}
	}

	// Setup first window command
	if firstWindow.Cmd != nil && *firstWindow.Cmd != "" {
		cmd := exec.Command("tmux", "send-keys", "-t", fmt.Sprintf("%s:%s", sessionName, firstWindow.WindowName), *firstWindow.Cmd, "C-m")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to send command to window '%s': %w", firstWindow.WindowName, err)
		}
	}

	// Switch to the new session
	return t.SwitchSession(sessionName)
}
