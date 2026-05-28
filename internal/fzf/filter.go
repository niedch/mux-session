package fzf

import (
	"sort"
	"strings"

	"github.com/niedch/mux-session/internal/dataproviders"
)

const (
	prefixMatchBonus = 100
)

// getMatches finds the byte offsets of characters in lowercased text
// that match the query runes in sequence (subsequence/fuzzy match).
func getMatches(textLower string, queryRunes []rune) []int {
	if len(queryRunes) == 0 || len(queryRunes) > len(textLower) {
		return nil
	}
	matches := make([]int, 0, len(queryRunes))
	qi := 0
	for bi, r := range textLower {
		if queryRunes[qi] == r {
			matches = append(matches, bi)
			qi++
			if qi == len(queryRunes) {
				return matches
			}
		}
	}
	return nil
}

// getMatchQuality scores how well lowercased text matches the lowercased query
// as a contiguous substring. Higher is better.
func getMatchQuality(textLower string, queryLower string) int {
	matchStart := strings.Index(textLower, queryLower)
	if matchStart == -1 {
		return 0
	}
	suffixLength := len(textLower) - len(queryLower)
	if matchStart == 0 {
		return prefixMatchBonus + suffixLength
	}
	return suffixLength + 1
}

// filterNode recursively filters children of an item whose display did not match.
// It returns the (potentially modified) item and a boolean indicating if a match was found.
func filterNode(item dataproviders.Item, queryLower string, queryRunes []rune) (dataproviders.Item, bool) {
	if len(item.SubItems) == 0 {
		return dataproviders.Item{}, false
	}

	var matchingChildren []dataproviders.Item
	for _, child := range item.SubItems {
		textLower := strings.ToLower(child.Display)
		if getMatches(textLower, queryRunes) != nil {
			matchingChildren = append(matchingChildren, child)
		} else {
			filteredChild, childMatches := filterNode(child, queryLower, queryRunes)
			if childMatches {
				matchingChildren = append(matchingChildren, filteredChild)
			}
		}
	}

	if len(matchingChildren) > 0 {
		newItem := item
		newItem.SubItems = matchingChildren
		return newItem, true
	}

	return dataproviders.Item{}, false
}

// FilterTree filters a slice of items based on a query.
// It returns a new tree containing only items that match or have children that match.
// Items are sorted by match quality (longer items with the query as prefix come first).
func FilterTree(items []dataproviders.Item, query string) []dataproviders.Item {
	queryLower := strings.ToLower(query)
	if len(queryLower) == 0 {
		return items
	}
	if len(items) == 0 {
		return nil
	}
	queryRunes := []rune(queryLower)

	type itemWithQuality struct {
		item    dataproviders.Item
		quality int
	}

	itemsWithQuality := make([]itemWithQuality, 0, len(items))
	for _, item := range items {
		textLower := strings.ToLower(item.Display)
		if getMatches(textLower, queryRunes) != nil {
			quality := getMatchQuality(textLower, queryLower)
			itemsWithQuality = append(itemsWithQuality, itemWithQuality{item: item, quality: quality})
		} else {
			filteredItem, hasMatch := filterNode(item, queryLower, queryRunes)
			if hasMatch {
				itemsWithQuality = append(itemsWithQuality, itemWithQuality{item: filteredItem, quality: 0})
			}
		}
	}

	if len(itemsWithQuality) == 0 {
		return nil
	}

	sort.Slice(itemsWithQuality, func(i, j int) bool {
		return itemsWithQuality[i].quality < itemsWithQuality[j].quality
	})

	filteredItems := make([]dataproviders.Item, 0, len(itemsWithQuality))
	for _, iwq := range itemsWithQuality {
		filteredItems = append(filteredItems, iwq.item)
	}

	return filteredItems
}

// createListItems creates the list of listItem models for rendering.
func createListItems(items []dataproviders.Item, query string) []listItem {
	queryLower := strings.ToLower(query)
	if len(queryLower) == 0 || len(items) == 0 {
		return nil
	}
	queryRunes := []rune(queryLower)
	result := make([]listItem, 0, len(items))

	for i, item := range items {
		textLower := strings.ToLower(item.Display)
		matches := getMatches(textLower, queryRunes)
		result = append(result, listItem{
			text:    item.Display,
			index:   i,
			matches: matches,
		})
	}
	return result
}
