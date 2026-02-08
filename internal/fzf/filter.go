package fzf

import (
	"strings"

	"github.com/niedch/mux-session/internal/dataproviders"
)

func filterItems(items []dataproviders.Item, query string) []listItem {
	var result []listItem
	queryRunes := []rune(strings.ToLower(query))

	for i, dpItem := range items {
		itemText := dpItem.Display
		targetRunes := []rune(strings.ToLower(itemText))
		queryIndex, targetIndex := 0, 0
		matches := make([]int, 0, len(queryRunes))

		for queryIndex < len(queryRunes) && targetIndex < len(targetRunes) {
			if queryRunes[queryIndex] == targetRunes[targetIndex] {
				matches = append(matches, targetIndex)
				queryIndex++
			}
			targetIndex++
		}

		if queryIndex == len(queryRunes) {
			result = append(result, listItem{text: itemText, index: i, matches: matches})
		}
	}
	return result
}
