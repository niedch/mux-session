package tmux

import (
	"fmt"
	"path/filepath"
	"slices"

	"github.com/niedch/mux-session/internal/conf"
)

func NewTmux(socket ...string) (*Tmux, error) {
	t := &Tmux{}
	if len(socket) > 0 {
		t.socket = socket[0]
	}
	return t, nil
}

type Tmux struct {
	socket string
}

func (t *Tmux) commandOpts() []OptFunc {
	if t.socket != "" {
		return []OptFunc{WithSocket(t.socket)}
	}
	return nil
}

func (t *Tmux) SetEnvironment(env map[string]string) error {
	opts := append(t.commandOpts(),
		WithTarget(t.socket),
	)

	for k, v := range env {
		opts = append(opts, WithArgs(k, v))
	}

	if err := SetEnvironment(opts...); err != nil {
		return fmt.Errorf("failed to set environment variables for session '%s': %w", t.socket, err)
	}

	return nil
}

func (t *Tmux) ListSessions() ([]string, error) {
	opts := append(t.commandOpts(), WithFormat("#S"))
	return ListSessions(opts...)
}

func (t *Tmux) CurrentSession() (string, error) {
	opts := append(t.commandOpts(), WithPrint(), WithFormat("#S"))
	return DisplayMessage(opts...)
}

func (t *Tmux) SwitchSession(sessionName string) error {
	currentSession, err := t.CurrentSession()
	if err != nil {
		return err
	}

	if currentSession == sessionName {
		return nil
	}

	sessions, err := t.ListSessions()
	if err != nil {
		return err
	}

	if slices.Contains(sessions, sessionName) {
		opts := append(t.commandOpts(), WithTarget(sessionName))
		return SwitchClient(opts...)
	}

	return fmt.Errorf("session %s not found", sessionName)
}

func (t *Tmux) CreateSession(dirPath string, projectConfig conf.ProjectConfig) error {
	sessionName := filepath.Base(dirPath)
	if projectConfig.Name != nil {
		sessionName = *projectConfig.Name
	}

	if len(projectConfig.WindowConfig) == 0 {
		return fmt.Errorf("no window configuration found for session %s", sessionName)
	}

	firstWindow := projectConfig.WindowConfig[0]

	// Create new session with first window
	opts := append(t.commandOpts(),
		WithDetached(),
		WithSession(sessionName),
		WithWindowName(firstWindow.WindowName),
		WithWorkingDir(dirPath),
	)
	if err := NewSession(opts...); err != nil {
		return fmt.Errorf("failed to create session %s: %w", sessionName, err)
	}

	if len(projectConfig.Env) > 0 {
		if err := t.SetEnvironment(projectConfig.Env); err != nil {
			return fmt.Errorf("failed to set environment %s: %w", sessionName, err)
		}
	}

	// Setup panels for first window if configured
	if len(firstWindow.PanelConfig) > 0 {
		if err := t.setupPanels(sessionName, firstWindow.WindowName, dirPath, firstWindow.PanelConfig); err != nil {
			return fmt.Errorf("failed to setup panels for first window: %w", err)
		}
	}

	// Create additional windows
	for _, window := range projectConfig.WindowConfig[1:] {
		if err := t.createWindowWithPanels(sessionName, dirPath, window); err != nil {
			return err
		}
	}

	// Execute window command if specified (after panels are created)
	if firstWindow.Cmd != nil && *firstWindow.Cmd != "" {
		target := fmt.Sprintf("%s:%s", sessionName, firstWindow.WindowName)
		opts := append(t.commandOpts(),
			WithTarget(target),
			WithKey(*firstWindow.Cmd),
			WithEnter(),
		)
		if err := SendKeys(opts...); err != nil {
			return fmt.Errorf("failed to send command to window %s: %w", firstWindow.WindowName, err)
		}
	}

	// Select primary window if configured, otherwise select first window
	primaryWindow := t.findPrimaryWindow(projectConfig.WindowConfig)
	if primaryWindow != "" {
		target := fmt.Sprintf("%s:%s", sessionName, primaryWindow)
		opts := append(t.commandOpts(), WithTarget(target))
		if err := SelectWindow(opts...); err != nil {
			return fmt.Errorf("failed to select primary window %s: %w", primaryWindow, err)
		}
	}

	// Switch to the new session
	return t.SwitchSession(sessionName)
}

func (t *Tmux) createWindowWithPanels(sessionName string, dirPath string, window conf.WindowConfig) error {
	target := fmt.Sprintf("%s:", sessionName)

	// Create new window
	opts := append(t.commandOpts(),
		WithTarget(target),
		WithWindowName(window.WindowName),
		WithWorkingDir(dirPath),
	)
	if err := NewWindow(opts...); err != nil {
		return fmt.Errorf("failed to create window %s in session %s: %w", window.WindowName, sessionName, err)
	}

	// Setup panels if configured
	if len(window.PanelConfig) > 0 {
		if err := t.setupPanels(sessionName, window.WindowName, dirPath, window.PanelConfig); err != nil {
			return fmt.Errorf("failed to setup panels for window %s: %w", window.WindowName, err)
		}
	}

	// Execute window command if specified (after panels are created)
	if window.Cmd != nil && *window.Cmd != "" {
		target := fmt.Sprintf("%s:%s", sessionName, window.WindowName)
		opts := append(t.commandOpts(),
			WithTarget(target),
			WithKey(*window.Cmd),
			WithEnter(),
		)
		if err := SendKeys(opts...); err != nil {
			return fmt.Errorf("failed to send command to window %s: %w", window.WindowName, err)
		}
	}

	return nil
}

func (t *Tmux) setupPanels(sessionName, windowName, dirPath string, panels []conf.PanelConfig) error {
	if len(panels) == 0 {
		return nil
	}

	target := fmt.Sprintf("%s:%s", sessionName, windowName)

	// First panel is already created with the window
	if len(panels) == 1 {
		// Just execute the command for the single panel
		if panels[0].Cmd != "" {
			opts := append(t.commandOpts(),
				WithTarget(target),
				WithKey(panels[0].Cmd),
				WithEnter(),
			)
			if err := SendKeys(opts...); err != nil {
				return fmt.Errorf("failed to send command to panel: %w", err)
			}
		}
		return nil
	}

	// Execute command for first panel if specified
	if panels[0].Cmd != "" {
		opts := append(t.commandOpts(),
			WithTarget(target),
			WithKey(panels[0].Cmd),
			WithEnter(),
		)
		if err := SendKeys(opts...); err != nil {
			return fmt.Errorf("failed to send command to first panel: %w", err)
		}
	}

	// Create additional panels
	for i, panel := range panels[1:] {
		var splitOpt OptFunc
		switch panel.PanelDirection {
		case "v":
			splitOpt = WithVertical()
		case "h":
			splitOpt = WithHorizontal()
		default:
			return fmt.Errorf("invalid panel direction %s for panel %d", panel.PanelDirection, i+1)
		}

		opts := append(t.commandOpts(),
			WithTarget(target),
			splitOpt,
			WithWorkingDir(dirPath),
		)
		if err := SplitWindow(opts...); err != nil {
			return fmt.Errorf("failed to create split for panel %d: %w", i+1, err)
		}

		// Execute command for the new panel if specified
		if panel.Cmd != "" {
			opts := append(t.commandOpts(),
				WithTarget(target),
				WithKey(panel.Cmd),
				WithEnter(),
			)
			if err := SendKeys(opts...); err != nil {
				return fmt.Errorf("failed to send command to panel %d: %w", i+1, err)
			}
		}
	}

	return nil
}

func (t *Tmux) findPrimaryWindow(windows []conf.WindowConfig) string {
	for _, window := range windows {
		if window.Primary != nil && *window.Primary {
			return window.WindowName
		}
	}
	return ""
}
