package fzf

import (
	"strings"
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

	assert.Equal(t, "apple\n", renderItems(FilterTree(items, "ap")))
	assert.Equal(t, "banana\n", renderItems(FilterTree(items, "nan")))
	assert.Equal(t, "", renderItems(FilterTree(items, "z")))
	assert.Equal(t, "apple\nbanana\ncherry\n", renderItems(FilterTree(items, "")))
}

func TestFilterTree_PrioritizeCloserMatches(t *testing.T) {
	items := []dataproviders.Item{
		{Display: "/home/nic/nixos-dotfiles"},
		{Display: "/home/nic/dotfiles"},
		{Display: "/home/nic/mux-prompt"},
	}

	assert.Equal(t, "/home/nic/mux-prompt\n/home/nic/nixos-dotfiles\n", renderItems(FilterTree(items, "nix")))
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

	assert.Equal(t, "fruit\n  apple\n  banana\n", renderItems(FilterTree(items, "fruit")))
	assert.Equal(t, "fruit\n  apple\n", renderItems(FilterTree(items, "apple")))
	assert.Equal(t, "vegetable\n  carrot\n", renderItems(FilterTree(items, "carrot")))
}


func renderItems(items []dataproviders.Item) string {
	var b strings.Builder
	renderItemsRec(&b, items, "")
	return b.String()
}

func renderItemsRec(b *strings.Builder, items []dataproviders.Item, indent string) {
	for _, item := range items {
		b.WriteString(indent)
		b.WriteString(item.Display)
		b.WriteByte('\n')
		if len(item.SubItems) > 0 {
			renderItemsRec(b, item.SubItems, indent+"  ")
		}
	}
}
