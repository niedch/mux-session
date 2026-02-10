package dataproviders

import (
	"maps"
	"strings"
)

type DeduplicatorProvider struct {
	directoryProvider   DataProvider
	multiplexerProvider DataProvider
	markDuplicates      bool
}

func NewDeduplicatorProvider(directoryProvider DataProvider, multiplexerProvider DataProvider) *DeduplicatorProvider {
	return &DeduplicatorProvider{
		directoryProvider:   directoryProvider,
		multiplexerProvider: multiplexerProvider,
	}
}

func (dp *DeduplicatorProvider) WithMarkDuplicates(mark bool) *DeduplicatorProvider {
	dp.markDuplicates = mark
	return dp
}

func (dp *DeduplicatorProvider) GetItems() ([]Item, error) {
	if dp.markDuplicates {
		return dp.markAndFilterItems()
	}

	multiplexerItems, err := dp.multiplexerProvider.GetItems()
	if err != nil {
		return nil, err
	}

	directoryItems, err := dp.directoryProvider.GetItems()
	if err != nil {
		return nil, err
	}

	directoryIds := flattenItems(directoryItems)
	filteredMultiplexerItems := filterItems(multiplexerItems, directoryIds)

	return append(directoryItems, filteredMultiplexerItems...), nil
}

func (dp *DeduplicatorProvider) markAndFilterItems() ([]Item, error) {
	directoryItems, err := dp.directoryProvider.GetItems()
	if err != nil {
		return nil, err
	}
	multiplexerItems, err := dp.multiplexerProvider.GetItems()
	if err != nil {
		return nil, err
	}

	multiplexerIds := flattenItems(multiplexerItems)
	markDuplicatesInItems(&directoryItems, multiplexerIds)

	directoryIds := flattenItems(directoryItems)
	filteredMultiplexerItems := filterItems(multiplexerItems, directoryIds)

	return append(directoryItems, filteredMultiplexerItems...), nil
}

func flattenItems(items []Item) map[string]bool {
	itemMap := make(map[string]bool)
	for _, item := range items {
		itemMap[item.Id] = true
		if len(item.SubItems) > 0 {
			maps.Copy(itemMap, flattenItems(item.SubItems))
		}
	}
	return itemMap
}

func filterItems(items []Item, ids map[string]bool) []Item {
	var filtered []Item
	for _, item := range items {
		if _, found := ids[item.Id]; !found {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func markDuplicatesInItems(items *[]Item, ids map[string]bool) {
	for i := range *items {
		if _, found := ids[(*items)[i].Id]; found {
			(*items)[i].Display = strings.Replace((*items)[i].Display, "[ ]", "[x]", 1)
		}
		if len((*items)[i].SubItems) > 0 {
			markDuplicatesInItems(&(*items)[i].SubItems, ids)
		}
	}
}
