package dataproviders

import (
	"github.com/niedch/mux-session/internal/tmux"
)

// TmuxProvider implements DataProvider for directory browsing
type TmuxProvider struct {
	tmux *tmux.Tmux
}

// NewTmuxProvider creates a new directory provider
func NewTmuxProvider(tmux *tmux.Tmux) *TmuxProvider {
	return &TmuxProvider{
		tmux: tmux,
	}
}

// GetItems returns the directories to display
func (dp *TmuxProvider) GetItems() ([]Item, error) {
	var items []Item
	sessions, err := dp.tmux.ListSessions()
	if err != nil {
		// Return empty list if no tmux server is running
		return items, nil
	}

	for _, session := range sessions {
		items = append(items, Item{
			Id:      session,
			Display: "[TMUX] " + session,
			Path:    session,
		})
	}

	return items, nil
}
