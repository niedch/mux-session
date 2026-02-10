package tree

import (
	"fmt"
	"os"
	"strings"

	"github.com/niedch/mux-session/internal/dataproviders"
)

const (
	TreeBranch   = " ├── "
	TreeLast     = " └── "
	TreeVertical = " │   "
	TreeEmpty    = "     "
)

type Renderer struct {
	items []dataproviders.Item
}

func NewRenderer(items []dataproviders.Item) *Renderer {
	return &Renderer{items: items}
}

func (r *Renderer) Render() []string {
	var output []string
	for _, item := range r.items {
		output = append(output, r.renderItem(item, 0, true, true))
	}
	return output
}

func (r *Renderer) renderItem(item dataproviders.Item, level int, isLast bool, isRoot bool) string {
	var prefix string

	if !isRoot {
		for i := 0; i < level-1; i++ {
			prefix += TreeVertical
		}
		if isLast {
			prefix += TreeLast
		} else {
			prefix += TreeBranch
		}
	}

	return prefix + item.Display
}

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

func GetSubdirectories(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var dirs []string
	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			dirs = append(dirs, entry.Name())
		}
	}

	return dirs, nil
}

func FormatSubdirectory(name string) string {
	return fmt.Sprintf("[ ] %s/", name)
}
