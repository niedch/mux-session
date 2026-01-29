package tmux

import (
	"fmt"
	"os/exec"
	"slices"
	"strings"
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
