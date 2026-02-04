package tmux

import (
	"fmt"
	"path/filepath"
	"slices"

	"github.com/niedch/mux-session/internal/conf"
)

type SessionCreator struct {
	sessionName string
	dirPath     string
	config      conf.ProjectConfig
}

func NewSessionCreator(dirPath string, config conf.ProjectConfig) *SessionCreator {
	sessionName := filepath.Base(dirPath)
	if config.Name != nil {
		sessionName = *config.Name
	}
	return &SessionCreator{
		sessionName: sessionName,
		dirPath:     dirPath,
		config:      config,
	}
}

func (sc *SessionCreator) Create() error {
	if len(sc.config.WindowConfig) == 0 {
		return fmt.Errorf("no window configuration found for session %s", sc.sessionName)
	}

	firstWindow := sc.config.WindowConfig[0]

	if err := sc.createInitialSession(firstWindow); err != nil {
		return err
	}

	if err := sc.setupPanels(firstWindow); err != nil {
		return err
	}

	for _, window := range sc.config.WindowConfig[1:] {
		if err := sc.createAdditionalWindow(window); err != nil {
			return err
		}
	}

	if err := sc.executeFirstWindowCommand(firstWindow); err != nil {
		return err
	}

	if err := sc.selectPrimaryWindow(); err != nil {
		return err
	}

	return sc.switchToSession()
}

func (sc *SessionCreator) createInitialSession(firstWindow conf.WindowConfig) error {
	return NewSession(
		WithDetached(),
		WithSession(sc.sessionName),
		WithWindowName(firstWindow.WindowName),
		WithWorkingDir(sc.dirPath),
	)
}

func (sc *SessionCreator) setupPanels(window conf.WindowConfig) error {
	if len(window.PanelConfig) == 0 {
		return nil
	}

	target := fmt.Sprintf("%s:%s", sc.sessionName, window.WindowName)

	if len(window.PanelConfig) == 1 {
		if window.PanelConfig[0].Cmd != "" {
			return SendKeys(
				WithTarget(target),
				WithKey(window.PanelConfig[0].Cmd),
				WithEnter(),
			)
		}
		return nil
	}

	if window.PanelConfig[0].Cmd != "" {
		if err := SendKeys(
			WithTarget(target),
			WithKey(window.PanelConfig[0].Cmd),
			WithEnter(),
		); err != nil {
			return fmt.Errorf("failed to send command to first panel: %w", err)
		}
	}

	for i, panel := range window.PanelConfig[1:] {
		if err := sc.splitPane(target, panel); err != nil {
			return fmt.Errorf("failed to create split for panel %d: %w", i+1, err)
		}

		if panel.Cmd != "" {
			if err := SendKeys(
				WithTarget(target),
				WithKey(panel.Cmd),
				WithEnter(),
			); err != nil {
				return fmt.Errorf("failed to send command to panel %d: %w", i+1, err)
			}
		}
	}

	return nil
}

func (sc *SessionCreator) splitPane(target string, panel conf.PanelConfig) error {
	var splitOpt OptFunc
	switch panel.PanelDirection {
	case "v":
		splitOpt = WithVertical()
	case "h":
		splitOpt = WithHorizontal()
	default:
		return fmt.Errorf("invalid panel direction %s", panel.PanelDirection)
	}

	return SplitWindow(
		WithTarget(target),
		splitOpt,
		WithWorkingDir(sc.dirPath),
	)
}

func (sc *SessionCreator) createAdditionalWindow(window conf.WindowConfig) error {
	target := fmt.Sprintf("%s:", sc.sessionName)

	if err := NewWindow(
		WithTarget(target),
		WithWindowName(window.WindowName),
		WithWorkingDir(sc.dirPath),
	); err != nil {
		return fmt.Errorf("failed to create window %s: %w", window.WindowName, err)
	}

	if err := sc.setupPanelsForWindow(window); err != nil {
		return err
	}

	if window.Cmd != nil && *window.Cmd != "" {
		target := fmt.Sprintf("%s:%s", sc.sessionName, window.WindowName)
		if err := SendKeys(
			WithTarget(target),
			WithKey(*window.Cmd),
			WithEnter(),
		); err != nil {
			return fmt.Errorf("failed to send command to window %s: %w", window.WindowName, err)
		}
	}

	return nil
}

func (sc *SessionCreator) setupPanelsForWindow(window conf.WindowConfig) error {
	if len(window.PanelConfig) == 0 {
		return nil
	}

	target := fmt.Sprintf("%s:%s", sc.sessionName, window.WindowName)

	if len(window.PanelConfig) == 1 {
		if window.PanelConfig[0].Cmd != "" {
			return SendKeys(
				WithTarget(target),
				WithKey(window.PanelConfig[0].Cmd),
				WithEnter(),
			)
		}
		return nil
	}

	if window.PanelConfig[0].Cmd != "" {
		if err := SendKeys(
			WithTarget(target),
			WithKey(window.PanelConfig[0].Cmd),
			WithEnter(),
		); err != nil {
			return fmt.Errorf("failed to send command to first panel: %w", err)
		}
	}

	for i, panel := range window.PanelConfig[1:] {
		if err := sc.splitPane(target, panel); err != nil {
			return fmt.Errorf("failed to create split for panel %d: %w", i+1, err)
		}

		if panel.Cmd != "" {
			if err := SendKeys(
				WithTarget(target),
				WithKey(panel.Cmd),
				WithEnter(),
			); err != nil {
				return fmt.Errorf("failed to send command to panel %d: %w", i+1, err)
			}
		}
	}

	return nil
}

func (sc *SessionCreator) executeFirstWindowCommand(firstWindow conf.WindowConfig) error {
	if firstWindow.Cmd == nil || *firstWindow.Cmd == "" {
		return nil
	}

	target := fmt.Sprintf("%s:%s", sc.sessionName, firstWindow.WindowName)
	return SendKeys(
		WithTarget(target),
		WithKey(*firstWindow.Cmd),
		WithEnter(),
	)
}

func (sc *SessionCreator) selectPrimaryWindow() error {
	primaryWindow := sc.findPrimaryWindow()
	if primaryWindow == "" {
		return nil
	}

	target := fmt.Sprintf("%s:%s", sc.sessionName, primaryWindow)
	return SelectWindow(WithTarget(target))
}

func (sc *SessionCreator) findPrimaryWindow() string {
	for _, window := range sc.config.WindowConfig {
		if window.Primary != nil && *window.Primary {
			return window.WindowName
		}
	}
	return ""
}

func (sc *SessionCreator) switchToSession() error {
	currentSession, err := CurrentSession()
	if err != nil {
		return err
	}

	if currentSession == sc.sessionName {
		return nil
	}

	sessions, err := ListSessions(WithFormat("#S"))
	if err != nil {
		return err
	}

	if slices.Contains(sessions, sc.sessionName) {
		return SwitchClient(WithTarget(sc.sessionName))
	}

	return fmt.Errorf("session %s not found", sc.sessionName)
}
