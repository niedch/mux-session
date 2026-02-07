package app

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type searchList struct {
	items  []searchListItem
	cursor int
	width  int
}

type searchListItem string

func newSearchList(items []string) *searchList {
	slItems := make([]searchListItem, len(items))
	for i, item := range items {
		slItems[i] = searchListItem(item)
	}

	return &searchList{
		items:  slItems,
		cursor: 0,
		width:  0,
	}
}

func (sl *searchList) SetWidth(width int) {
	sl.width = width
}

func (sl *searchList) Update(msg tea.Msg) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "ctrl+p":
			if sl.cursor > 0 {
				sl.cursor--
			}
		case "down", "ctrl+n":
			if sl.cursor < len(sl.items)-1 {
				sl.cursor++
			}
		}
	}
}

func (sl *searchList) View(height int) string {
	if height <= 0 {
		return ""
	}

	// Render items from bottom to top (last items at bottom)
	var visibleItems []string
	for idx := len(sl.items) - 1; idx >= 0 && len(visibleItems) < height; idx-- {
		if idx == sl.cursor {
			visibleItems = append([]string{selectedItemStyle.Render("> " + string(sl.items[idx]))}, visibleItems...)
		} else {
			visibleItems = append([]string{itemStyle.Render("  " + string(sl.items[idx]))}, visibleItems...)
		}
	}

	// If we have fewer items than available height, add padding at top
	if len(visibleItems) < height {
		padding := strings.Repeat("\n", height-len(visibleItems))
		visibleItems = append([]string{padding}, visibleItems...)
	}

	listContent := strings.Join(visibleItems, "\n")

	return lipgloss.NewStyle().
		Width(sl.width).
		Render(listContent)
}
