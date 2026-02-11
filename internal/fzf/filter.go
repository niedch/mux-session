package fzf

import (
	"strings"

	"github.com/niedch/mux-session/internal/dataproviders"
)

// getMatches finds the indices of characters in text that match the query runes.
// This is used for highlighting.
func getMatches(text string, queryRunes []rune) []int {
	if len(queryRunes) == 0 {
		return nil
	}
	matches := make([]int, 0, len(queryRunes))
	targetRunes := []rune(strings.ToLower(text))
	queryIndex, targetIndex := 0, 0

	for queryIndex < len(queryRunes) && targetIndex < len(targetRunes) {
		if queryRunes[queryIndex] == targetRunes[targetIndex] {
			matches = append(matches, targetIndex)
			queryIndex++
		}
		targetIndex++
	}

	// Only return matches if the entire query was found
	if queryIndex == len(queryRunes) {
		return matches
	}

	return nil
}

// filterNode recursively filters a single item and its children.
// It returns the (potentially modified) item and a boolean indicating if a match was found.
func filterNode(item dataproviders.Item, queryRunes []rune) (dataproviders.Item, bool) {
	isMatch := getMatches(item.Display, queryRunes) != nil

	if isMatch {
		// If the parent matches, return it with all its original children.
		return item, true
	}

	// If the parent does not match, check the children.
	var matchingChildren []dataproviders.Item
	hasMatchingChild := false
	for _, child := range item.SubItems {
		filteredChild, childMatches := filterNode(child, queryRunes)
		if childMatches {
			hasMatchingChild = true
			matchingChildren = append(matchingChildren, filteredChild)
		}
	}

	// If any child matches, return a new item with only the matching children.
	if hasMatchingChild {
		newItem := item
		newItem.SubItems = matchingChildren
		return newItem, true
	}

	// If neither the item nor its children match, it's not included.
	return dataproviders.Item{}, false
}

// FilterTree filters a slice of items based on a query.
// It returns a new tree containing only items that match or have children that match.
func FilterTree(items []dataproviders.Item, query string) []dataproviders.Item {
	queryRunes := []rune(strings.ToLower(query))
	if len(queryRunes) == 0 {
		return items
	}

	var filteredItems []dataproviders.Item
	for _, item := range items {
		filteredItem, matches := filterNode(item, queryRunes)
		if matches {
			filteredItems = append(filteredItems, filteredItem)
		}
	}
	return filteredItems
}

// createListItems creates the list of listItem models for rendering.
func createListItems(items []dataproviders.Item, query string) []listItem {
	var result []listItem
	queryRunes := []rune(strings.ToLower(query))

	for i, item := range items {
		matches := getMatches(item.Display, queryRunes)
		result = append(result, listItem{
			text:    item.Display,
			index:   i,
			matches: matches,
		})
	}
	return result
}
