package app

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type searchList struct {
	allItems    []searchListItem
	items       []searchListItem
	cursor      int
	width       int
	searchQuery string
}

type searchListItem string

func newSearchList(items []string) *searchList {
	slItems := make([]searchListItem, len(items))
	for i, item := range items {
		slItems[i] = searchListItem(item)
	}

	return &searchList{
		allItems: slItems,
		items:    slItems,
		cursor:   0,
		width:    0,
	}
}

func (sl *searchList) SetWidth(width int) {
	sl.width = width
}

func (sl *searchList) SetSearch(query string) {
	if sl.searchQuery == query {
		return
	}

	sl.searchQuery = query
	sl.cursor = 0

	if query == "" {
		sl.items = sl.allItems
		return
	}

	var filtered []searchListItem
	lowerQuery := strings.ToLower(query)
	for _, item := range sl.allItems {
		if strings.Contains(strings.ToLower(string(item)), lowerQuery) {
			filtered = append(filtered, item)
		}
	}
	sl.items = filtered
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
	if height <= 0 || len(sl.items) == 0 {
		return strings.Repeat("\n", height)
	}

	// Ensure cursor is within bounds
	if sl.cursor >= len(sl.items) {
		sl.cursor = len(sl.items) - 1
	}
	if sl.cursor < 0 {
		sl.cursor = 0
	}

	// Calculate which items to show based on cursor position
	// Show items around the cursor, centered if possible
	var startIdx int
	if len(sl.items) <= height {
		// All items fit, show from start
		startIdx = 0
	} else if sl.cursor < height/2 {
		// Cursor is near the top
		startIdx = 0
	} else if sl.cursor > len(sl.items)-height/2-1 {
		// Cursor is near the bottom
		startIdx = len(sl.items) - height
	} else {
		// Cursor is in the middle, center it
		startIdx = sl.cursor - height/2
	}

	endIdx := startIdx + height
	if endIdx > len(sl.items) {
		endIdx = len(sl.items)
	}

	// Render visible items
	var visibleItems []string
	for idx := startIdx; idx < endIdx; idx++ {
		if idx == sl.cursor {
			visibleItems = append(visibleItems, selectedItemStyle.Render("> "+string(sl.items[idx])))
		} else {
			visibleItems = append(visibleItems, itemStyle.Render("  "+string(sl.items[idx])))
		}
	}

	// If we have fewer items than available height, add padding at top
	// (to push items to the bottom)
	paddingLines := height - len(visibleItems)
	padding := strings.Repeat("\n", paddingLines)
	visibleItems = append([]string{padding}, visibleItems...)

	listContent := strings.Join(visibleItems, "\n")

	return lipgloss.NewStyle().
		Width(sl.width).
		Render(listContent)
}
