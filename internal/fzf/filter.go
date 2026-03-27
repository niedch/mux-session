package fzf

import (
	"sort"
	"strings"

	"github.com/niedch/mux-session/internal/dataproviders"
)

const (
	prefixMatchBonus = 100
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

// getMatchQuality returns how well text matches the query.
// Higher values mean better matches. It considers:
// - Items where the query appears as a prefix are best (prefixMatchBonus bonus)
// - Among those, shorter text (longer suffix after query) is better
// - Items where the query appears elsewhere are sorted by where they start
func getMatchQuality(text string, queryRunes []rune) int {
	textLower := strings.ToLower(text)
	matchStart := strings.Index(textLower, string(queryRunes))
	if matchStart == -1 {
		return 0
	}

	textLength := len(textLower)
	queryLength := len(queryRunes)
	suffixLength := textLength - queryLength

	if matchStart == 0 {
		return prefixMatchBonus + suffixLength
	}

	return suffixLength + 1
}

// FilterTree filters a slice of items based on a query.
// It returns a new tree containing only items that match or have children that match.
// Items are sorted by match quality (longer items with the query as prefix come first).
func FilterTree(items []dataproviders.Item, query string) []dataproviders.Item {
	queryRunes := []rune(strings.ToLower(query))
	if len(queryRunes) == 0 {
		return items
	}

	type itemWithQuality struct {
		item    dataproviders.Item
		quality int
	}

	var itemsWithQuality []itemWithQuality
	for _, item := range items {
		filteredItem, matches := filterNode(item, queryRunes)
		if matches {
			quality := getMatchQuality(item.Display, queryRunes)
			itemsWithQuality = append(itemsWithQuality, itemWithQuality{item: filteredItem, quality: quality})
		}
	}

	sort.Slice(itemsWithQuality, func(i, j int) bool {
		return itemsWithQuality[i].quality > itemsWithQuality[j].quality
	})

	var filteredItems []dataproviders.Item
	for _, iwq := range itemsWithQuality {
		filteredItems = append(filteredItems, iwq.item)
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
