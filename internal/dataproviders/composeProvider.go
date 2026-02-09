package dataproviders

import "strings"

type ComposeProvider struct {
	dataproviders  []DataProvider
	markDuplicates bool
}

func NewComposeProvider(dataproviders ...DataProvider) *ComposeProvider {
	return &ComposeProvider{
		dataproviders: dataproviders,
	}
}

func (dp *ComposeProvider) WithMarkDuplicates(mark bool) *ComposeProvider {
	dp.markDuplicates = mark
	return dp
}

func (dp *ComposeProvider) GetItems() ([]Item, error) {
	var items []Item
	seen := make(map[string]int)

	for _, dataprovider := range dp.dataproviders {
		dp_items, err := dataprovider.GetItems()
		if err != nil {
			return nil, err
		}

		for _, item := range dp_items {
			if idx, ok := seen[item.Id]; ok {
				if dp.markDuplicates {
					if strings.HasPrefix(items[idx].Display, "[ ] ") {
						items[idx].Display = strings.Replace(items[idx].Display, "[ ] ", "[x] ", 1)
					}
				}
				continue
			}

			if dp.markDuplicates && !strings.HasPrefix(item.Display, "[TMUX]") && !strings.HasPrefix(item.Display, "[w] ") {
				item.Display = "[ ] " + item.Display
			}
			items = append(items, item)
			seen[item.Id] = len(items) - 1
		}
	}
	return items, nil
}
