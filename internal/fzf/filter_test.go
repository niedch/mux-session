package fzf

import (
	"testing"

	"github.com/niedch/mux-session/internal/dataproviders"
	"github.com/stretchr/testify/assert"
)

func TestFilterTree_Normal(t *testing.T) {
	items := []dataproviders.Item{
		{Display: "apple"},
		{Display: "banana"},
		{Display: "cherry"},
	}

	filtered := FilterTree(items, "ap")
	assert.Len(t, filtered, 1)
	assert.Equal(t, "apple", filtered[0].Display)

	filtered = FilterTree(items, "nan")
	assert.Len(t, filtered, 1)
	assert.Equal(t, "banana", filtered[0].Display)

	filtered = FilterTree(items, "z")
	assert.Empty(t, filtered)
}

func TestFilterTree_WithSubItems(t *testing.T) {
	items := []dataproviders.Item{
		{
			Display: "fruit",
			SubItems: []dataproviders.Item{
				{Display: "apple"},
				{Display: "banana"},
			},
		},
		{
			Display: "vegetable",
			SubItems: []dataproviders.Item{
				{Display: "carrot"},
			},
		},
	}

	// Filter matches parent (fruit) -> should return parent with all children
	filtered := FilterTree(items, "fruit")
	assert.Len(t, filtered, 1)
	assert.Equal(t, "fruit", filtered[0].Display)
	assert.Len(t, filtered[0].SubItems, 2)

	// Filter matches child (apple) -> should return parent with only matching child
	filtered = FilterTree(items, "apple")
	assert.Len(t, filtered, 1)
	assert.Equal(t, "fruit", filtered[0].Display)
	assert.Len(t, filtered[0].SubItems, 1)
	assert.Equal(t, "apple", filtered[0].SubItems[0].Display)

	// Filter matches child (carrot) -> should return parent (vegetable)
	filtered = FilterTree(items, "carrot")
	assert.Len(t, filtered, 1)
	assert.Equal(t, "vegetable", filtered[0].Display)
	assert.Len(t, filtered[0].SubItems, 1)
	assert.Equal(t, "carrot", filtered[0].SubItems[0].Display)
}
