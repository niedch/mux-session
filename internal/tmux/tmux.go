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
	if err := exec.Command("tmux", "new-session", "-d", "-s", sessionName, "-n", firstWindow.WindowName).Run(); err != nil {
		return fmt.Errorf("failed to create session '%s': %w", sessionName, err)
	}

	// Setup panels for first window if configured
	if len(firstWindow.PanelConfig) > 0 {
		if err := t.setupPanels(sessionName, firstWindow.WindowName, firstWindow.PanelConfig); err != nil {
			return fmt.Errorf("failed to setup panels for first window: %w", err)
		}
	}

	// Create additional windows
	for _, window := range projectConfig.WindowConfig[1:] {
		if err := t.createWindowWithPanels(sessionName, window); err != nil {
			return err
		}
	}

	// Execute window command if specified (after panels are created)
	if firstWindow.Cmd != nil && *firstWindow.Cmd != "" {
		target := fmt.Sprintf("%s:%s", sessionName, firstWindow.WindowName)

		if err := exec.Command("tmux", "send-keys", "-t", target, *firstWindow.Cmd, "C-m").Run(); err != nil {
			return fmt.Errorf("failed to send command to window '%s': %w", firstWindow.WindowName, err)
		}
	}

	// Switch to the new session
	return t.SwitchSession(sessionName)
}

func (t *Tmux) createWindowWithPanels(sessionName string, window conf.WindowConfig) error {
	// Create new window
	cmd := exec.Command("tmux", "new-window", "-t", sessionName, "-n", window.WindowName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create window '%s' in session '%s': %w", window.WindowName, sessionName, err)
	}

	// Setup panels if configured
	if len(window.PanelConfig) > 0 {
		if err := t.setupPanels(sessionName, window.WindowName, window.PanelConfig); err != nil {
			return fmt.Errorf("failed to setup panels for window '%s': %w", window.WindowName, err)
		}
	}

	// Execute window command if specified (after panels are created)
	if window.Cmd != nil && *window.Cmd != "" {
		target := fmt.Sprintf("%s:%s", sessionName, window.WindowName)
		cmd := exec.Command("tmux", "send-keys", "-t", target, *window.Cmd, "C-m")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to send command to window '%s': %w", window.WindowName, err)
		}
	}

	return nil
}

func (t *Tmux) setupPanels(sessionName, windowName string, panels []conf.PanelConfig) error {
	if len(panels) == 0 {
		return nil
	}

	// First panel is already created with the window
	if len(panels) == 1 {
		// Just execute the command for the single panel
		if panels[0].Cmd != "" {
			target := fmt.Sprintf("%s:%s", sessionName, windowName)
			if err := exec.Command("tmux", "send-keys", "-t", target, panels[0].Cmd, "C-m").Run(); err != nil {
				return fmt.Errorf("failed to send command to panel: %w", err)
			}
		}
		return nil
	}

	// Execute command for first panel if specified
	if panels[0].Cmd != "" {
		target := fmt.Sprintf("%s:%s", sessionName, windowName)
		cmd := exec.Command("tmux", "send-keys", "-t", target, panels[0].Cmd, "C-m")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to send command to first panel: %w", err)
		}
	}

	// Create additional panels
	for i, panel := range panels[1:] {
		target := fmt.Sprintf("%s:%s", sessionName, windowName)

		switch panel.PanelDirection {
		case "v":
			// Split vertically (top/bottom)
			cmd := exec.Command("tmux", "split-window", "-t", target, "-v")
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to create vertical split for panel %d: %w", i+1, err)
			}
		case "h":
			// Split horizontally (left/right)
			cmd := exec.Command("tmux", "split-window", "-t", target, "-h")
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to create horizontal split for panel %d: %w", i+1, err)
			}
		default:
			return fmt.Errorf("invalid panel direction '%s' for panel %d", panel.PanelDirection, i+1)
		}

		// Execute command for the new panel if specified
		if panel.Cmd != "" {
			// New panel becomes the active pane after split
			if err := exec.Command("tmux", "send-keys", "-t", target, panel.Cmd, "C-m").Run(); err != nil {
				return fmt.Errorf("failed to send command to panel %d: %w", i+1, err)
			}
		}
	}

	return nil
}
