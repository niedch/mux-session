package tree

import (
	"github.com/niedch/mux-session/internal/dataproviders"
)

const (
	TreeBranch   = " ├── "
	TreeLast     = " └── "
	TreeVertical = " │   "
	TreeEmpty    = "     "
)

func GeneratePrefix(level int, isLast bool) string {
	if level == 0 {
		return ""
	}

	var prefix string
	for i := 0; i < level-1; i++ {
		prefix += TreeVertical
	}

	if isLast {
		prefix += TreeLast
	} else {
		prefix += TreeBranch
	}

	return prefix
}

func FlattenItems(items []dataproviders.Item) []dataproviders.Item {
	var result []dataproviders.Item

	for i, item := range items {
		isLast := i == len(items)-1

		if item.TreeLevel > 0 {
			prefix := GeneratePrefix(item.TreeLevel, isLast)
			item.Display = prefix + item.Display
		}

		result = append(result, item)

		if len(item.SubItems) > 0 {
			subItems := make([]dataproviders.Item, len(item.SubItems))
			copy(subItems, item.SubItems)

			for j := range subItems {
				subItems[j].TreeLevel = item.TreeLevel + 1
			}

			result = append(result, FlattenItems(subItems)...)
		}
	}

	return result
}
