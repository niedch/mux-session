package tmux

import (
	"fmt"
	"slices"
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

func (t *Tmux) NewSession(sessionName string, firstWindowName string, workingDir string) error {
	// Create new session with first window
	opts := append(t.commandOpts(),
		WithDetached(),
		WithSession(sessionName),
		WithWindowName(firstWindowName),
		WithWorkingDir(workingDir),
	)
	if err := NewSession(opts...); err != nil {
		return fmt.Errorf("failed to create session %s: %w", sessionName, err)
	}
	return nil
}

func (t *Tmux) CreateWindow(target string, windowName string, workingDir string) error {
	opts := append(t.commandOpts(),
		WithTarget(target),
		WithWindowName(windowName),
		WithWorkingDir(workingDir),
	)

	if err := NewWindow(opts...); err != nil {
		return fmt.Errorf("failed to create window %s in session %s: %w", windowName, t.socket, err)
	}

	return nil
}

func (t *Tmux) SendKeys(target string, cmd string) error {
	opts := append(t.commandOpts(),
		WithTarget(target),
		WithKey(cmd),
		WithEnter(),
	)
	if err := SendKeys(opts...); err != nil {
		return fmt.Errorf("failed to send command to session %s: %w", t.socket, err)
	}
	return nil
}

func (t *Tmux) SplitWindow(target string, direction string, workingdir string) error {
	var splitOpt OptFunc
	switch direction {
	case "v":
		splitOpt = WithVertical()
	case "h":
		splitOpt = WithHorizontal()
	default:
		return fmt.Errorf("invalid panel direction %s", direction)
	}

	opts := append(t.commandOpts(),
		WithTarget(target),
		splitOpt,
		WithWorkingDir(workingdir),
	)
	if err := SplitWindow(opts...); err != nil {
		return fmt.Errorf("failed to create split for panel: %w", err)
	}

	return nil
}

func (t *Tmux) SetEnvironment(target string, env map[string]string) error {
	for k, v := range env {
		opts := append(t.commandOpts(),
			WithTarget(target),
			WithArgs(k, v),
		)

		if err := SetEnvironment(opts...); err != nil {
			return fmt.Errorf("failed to set environment variable %s for session '%s': %w", k, target, err)
		}
	}

	return nil
}

func (t *Tmux) FocusWindow(target string) error {
	opts := append(t.commandOpts(), WithTarget(target))
	if err := SelectWindow(opts...); err != nil {
		return fmt.Errorf("failed to focus target: %s: %w", target, err)
	}

	return nil
}

// func (t *Tmux) CreateWindowWithPanels(sessionName string, dirPath string, window conf.WindowConfig) error {
// 	target := fmt.Sprintf("%s:", sessionName)
//
// 	// Create new window
// 	opts := append(t.commandOpts(),
// 		WithTarget(target),
// 		WithWindowName(window.WindowName),
// 		WithWorkingDir(dirPath),
// 	)
// 	if err := NewWindow(opts...); err != nil {
// 		return fmt.Errorf("failed to create window %s in session %s: %w", window.WindowName, sessionName, err)
// 	}
//
// 	// Setup panels if configured
// 	if len(window.PanelConfig) > 0 {
// 		if err := t.SetupPanels(sessionName, window.WindowName, dirPath, window.PanelConfig); err != nil {
// 			return fmt.Errorf("failed to setup panels for window %s: %w", window.WindowName, err)
// 		}
// 	}
//
// 	// Execute window command if specified (after panels are created)
// 	if window.Cmd != nil && *window.Cmd != "" {
// 		target := fmt.Sprintf("%s:%s", sessionName, window.WindowName)
//
// 		opts := append(t.commandOpts(),
// 			WithTarget(target),
// 			WithKey(*window.Cmd),
// 			WithEnter(),
// 		)
//
// 		if err := SendKeys(opts...); err != nil {
// 			return fmt.Errorf("failed to send command to window %s: %w", window.WindowName, err)
// 		}
// 	}
//
// 	return nil
// }

// func (t *Tmux) SetupPanels(sessionName, windowName, dirPath string, panels []conf.PanelConfig) error {
// 	if len(panels) == 0 {
// 		return nil
// 	}
//
// 	target := fmt.Sprintf("%s:%s", sessionName, windowName)
//
// 	// First panel is already created with the window
// 	if len(panels) == 1 {
// 		// Just execute the command for the single panel
// 		if panels[0].Cmd != "" {
// 			opts := append(t.commandOpts(),
// 				WithTarget(target),
// 				WithKey(panels[0].Cmd),
// 				WithEnter(),
// 			)
// 			if err := SendKeys(opts...); err != nil {
// 				return fmt.Errorf("failed to send command to panel: %w", err)
// 			}
// 		}
// 		return nil
// 	}
//
// 	// Execute command for first panel if specified
// 	if panels[0].Cmd != "" {
// 		opts := append(t.commandOpts(),
// 			WithTarget(target),
// 			WithKey(panels[0].Cmd),
// 			WithEnter(),
// 		)
// 		if err := SendKeys(opts...); err != nil {
// 			return fmt.Errorf("failed to send command to first panel: %w", err)
// 		}
// 	}
//
// 	// Create additional panels
// 	for i, panel := range panels[1:] {
// 		var splitOpt OptFunc
// 		switch panel.PanelDirection {
// 		case "v":
// 			splitOpt = WithVertical()
// 		case "h":
// 			splitOpt = WithHorizontal()
// 		default:
// 			return fmt.Errorf("invalid panel direction %s for panel %d", panel.PanelDirection, i+1)
// 		}
//
// 		opts := append(t.commandOpts(),
// 			WithTarget(target),
// 			splitOpt,
// 			WithWorkingDir(dirPath),
// 		)
// 		if err := SplitWindow(opts...); err != nil {
// 			return fmt.Errorf("failed to create split for panel %d: %w", i+1, err)
// 		}
//
// 		// Execute command for the new panel if specified
// 		if panel.Cmd != "" {
// 			opts := append(t.commandOpts(),
// 				WithTarget(target),
// 				WithKey(panel.Cmd),
// 				WithEnter(),
// 			)
// 			if err := SendKeys(opts...); err != nil {
// 				return fmt.Errorf("failed to send command to panel %d: %w", i+1, err)
// 			}
// 		}
// 	}
//
// 	return nil
// }
