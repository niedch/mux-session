package dataproviders

type ComposeProvider struct {
	dataproviders []DataProvider
}

func NewComposeProvider(dataproviders ...DataProvider) *ComposeProvider {
	return &ComposeProvider{
		dataproviders: dataproviders,
	}
}

// Returns Items from all providers and removes duplicates that have the same Id 
func (dp *ComposeProvider) GetItems() ([]Item, error) {
	var items []Item
	seen := make(map[string]bool)

	for _, dataprovider := range dp.dataproviders {
		dp_items, err := dataprovider.GetItems()
		if err != nil {
			return nil, err
		}

		for _, item := range dp_items {
			if _, ok := seen[item.Id]; !ok {
				items = append(items, item)
				seen[item.Id] = true
			}
		}
	}
	return items, nil
}
