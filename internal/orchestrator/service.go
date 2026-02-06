package orchestrator

import (
	"fmt"
	"log"
	"path/filepath"
	"slices"

	"github.com/niedch/mux-session/internal/conf"
	"github.com/niedch/mux-session/internal/dataproviders"
	"github.com/niedch/mux-session/internal/tmux"
)

type OrchestratorService struct {
	tmux *tmux.Tmux
}

func New(tmux *tmux.Tmux) *OrchestratorService {
	return &OrchestratorService{
		tmux: tmux,
	}
}

func (m *OrchestratorService) CreateSession(item *dataproviders.Item, projectConfig conf.ProjectConfig) error {
	dirPath := item.Path
	sessionName := filepath.Base(dirPath)

	if projectConfig.Name != nil {
		sessionName = *projectConfig.Name
	}

	if len(projectConfig.WindowConfig) == 0 {
		return fmt.Errorf("no window configuration found for session %s", sessionName)
	}

	firstWindow := projectConfig.WindowConfig[0]

	log.Printf("Creating Session %s\n", sessionName)
	if err := m.tmux.NewSession(sessionName, firstWindow.WindowName, dirPath); err != nil {
		return fmt.Errorf("Failed to create Session %s", sessionName)
	}

	if len(projectConfig.Env) > 0 {
		if err := m.tmux.SetEnvironment(sessionName, projectConfig.Env); err != nil {
			return fmt.Errorf("failed to set environment %s: %w", sessionName, err)
		}
	}

	// Setup panels for first window if configured
	if len(firstWindow.PanelConfig) > 0 {
		if err := m.setupPanels(sessionName, firstWindow.WindowName, dirPath, firstWindow.PanelConfig); err != nil {
			return fmt.Errorf("failed to setup panels for first window: %w", err)
		}
	}

	// Create additional windows
	for _, window := range projectConfig.WindowConfig[1:] {
		if err := m.createWindowWithPanels(sessionName, dirPath, window); err != nil {
			return err
		}
	}

	// Execute window command if specified (after panels are created)
	if firstWindow.Cmd != nil && *firstWindow.Cmd != "" {
		target := fmt.Sprintf("%s:%s", sessionName, firstWindow.WindowName)
		if err := m.tmux.SendKeys(target, *firstWindow.Cmd); err != nil {
			return fmt.Errorf("failed to send command to window %s: %w", firstWindow.WindowName, err)
		}
	}

	// Select primary window if configured, otherwise select first window
	primaryWindow := m.findPrimaryWindow(projectConfig.WindowConfig)
	if primaryWindow != "" {
		target := fmt.Sprintf("%s:%s", sessionName, primaryWindow)
		if err := m.tmux.FocusWindow(target); err != nil {
			return fmt.Errorf("Failed to focus primary window %s: %w", primaryWindow, err)

		}
	}

	// Switch to the new session
	return m.tmux.SwitchSession(sessionName)
}

func (m *OrchestratorService) SwitchSession(selected *dataproviders.Item) (bool, error) {
	sessions, err := m.tmux.ListSessions()
	if err != nil {
		return false, err
	}

	if slices.Contains(sessions, selected.Id) {
		if err := m.tmux.SwitchSession(selected.Id); err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}

func (m *OrchestratorService) createWindowWithPanels(sessionName string, dirPath string, window conf.WindowConfig) error {
	target := fmt.Sprintf("%s:", sessionName)

	if err := m.tmux.CreateWindow(target, window.WindowName, dirPath); err != nil {
		return fmt.Errorf("failed to create window %s in session %s: %w", window.WindowName, sessionName, err)
	}

	if len(window.PanelConfig) > 0 {
		if err := m.setupPanels(sessionName, window.WindowName, dirPath, window.PanelConfig); err != nil {
			return fmt.Errorf("failed to setup panels for window %s: %w", window.WindowName, err)
		}
	}

	// Execute window command if specified (after panels are created)
	if window.Cmd != nil && *window.Cmd != "" {
		target := fmt.Sprintf("%s:%s", sessionName, window.WindowName)

		if err := m.tmux.SendKeys(target, *window.Cmd); err != nil {
			return fmt.Errorf("failed to send command to window %s: %w", window.WindowName, err)
		}
	}

	return nil
}

func (m *OrchestratorService) setupPanels(sessionName, windowName, dirPath string, panels []conf.PanelConfig) error {
	if len(panels) == 0 {
		return nil
	}

	target := fmt.Sprintf("%s:%s", sessionName, windowName)

	// First panel is already created with the window
	if len(panels) == 1 {
		// Just execute the command for the single panel
		if panels[0].Cmd != "" {
			if err := m.tmux.SendKeys(target, panels[0].Cmd); err != nil {
				return fmt.Errorf("failed to send command to panel: %w", err)
			}
		}

		return nil
	}

	// Execute command for first panel if specified
	if panels[0].Cmd != "" {
		if err := m.tmux.SendKeys(target, panels[0].Cmd); err != nil {
			return fmt.Errorf("failed to send command to panel: %w", err)
		}
	}

	for i, panel := range panels[1:] {
		if err := m.tmux.SplitWindow(target, panel.PanelDirection, dirPath); err != nil {
			return fmt.Errorf("failed to create split for panel %d: %w", i+1, err)
		}

		if err := m.tmux.SendKeys(target, panel.Cmd); err != nil {
			return fmt.Errorf("failed to send command to panel %d: %w", i+1, err)
		}
	}

	return nil
}

func (m *OrchestratorService) findPrimaryWindow(windows []conf.WindowConfig) string {
	for _, window := range windows {
		if window.Primary != nil && *window.Primary {
			return window.WindowName
		}
	}

	return ""
}
