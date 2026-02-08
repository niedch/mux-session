package fzf

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/niedch/mux-session/internal/dataproviders"
)

type listItem struct {
	text    string
	index   int
	matches []int
}

type list struct {
	items    []dataproviders.Item
	filtered []listItem
	cursor   int
}

func newList(items []dataproviders.Item) *list {
	// reverse items
	for i, j := 0, len(items)-1; i < j; i, j = i+1, j-1 {
		items[i], items[j] = items[j], items[i]
	}

	l := &list{
		items: items,
	}
	l.filter("")
	if len(l.filtered) > 0 {
		l.cursor = len(l.filtered) - 1
	}
	return l
}

func (l *list) filter(query string) {
	l.filtered = filterItems(l.items, query)
}

func (l *list) updateFilter(query string) {
	l.filter(query)
	if len(l.filtered) > 0 {
		l.cursor = len(l.filtered) - 1
	} else {
		l.cursor = 0
	}
}

func (l *list) moveUp() {
	if l.cursor > 0 {
		l.cursor--
	}
}

func (l *list) moveDown() {
	if l.cursor < len(l.filtered)-1 {
		l.cursor++
	}
}

func (l *list) getSelected() *dataproviders.Item {
	if len(l.filtered) > 0 {
		return &l.items[l.filtered[l.cursor].index]
	}
	return nil
}

func (l *list) calculateVisibleRange(listHeight int) (int, int) {
	if len(l.filtered) <= listHeight {
		return 0, len(l.filtered)
	}

	if l.cursor < listHeight/2 {
		return 0, listHeight
	}

	if l.cursor > len(l.filtered)-listHeight/2 {
		return len(l.filtered) - listHeight, len(l.filtered)
	}

	start := l.cursor - listHeight/2
	end := l.cursor + listHeight/2
	return start, end
}

func (l *list) renderItem(i int) string {
	if i < 0 || i >= len(l.filtered) {
		return ""
	}

	it := l.filtered[i]
	isSelected := i == l.cursor

	cursor := "  "
	style := lipgloss.NewStyle()
	if isSelected {
		cursor = "> "
		style = style.Bold(true)
	}

	var highlightedText strings.Builder
	matchStyle := style.Copy().Bold(true)

	matchMap := make(map[int]struct{})
	for _, idx := range it.matches {
		matchMap[idx] = struct{}{}
	}

	for i, r := range it.text {
		if _, ok := matchMap[i]; ok {
			highlightedText.WriteString(matchStyle.Render(string(r)))
		} else {
			highlightedText.WriteString(style.Render(string(r)))
		}
	}

	return fmt.Sprintf("%s%s\n", cursor, highlightedText.String())
}
