package previewproviders

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/niedch/mux-session/internal/dataproviders"
	"github.com/niedch/mux-session/internal/tree"
)

type TreePreviewProvider struct {
	width int
}

func NewTreePreviewProvider(width int) (*TreePreviewProvider, error) {
	return &TreePreviewProvider{
		width: width,
	}, nil
}

func (r *TreePreviewProvider) Render(item any) (string, error) {
	dpItem, ok := item.(*dataproviders.Item)
	if !ok {
		return "", fmt.Errorf("expected *dataproviders.Item, got %T", item)
	}

	var builder strings.Builder
	fmt.Fprintf(&builder, "%s\n", filepath.Base(dpItem.Path))
	err := buildTree(dpItem.Path, "", 0, 2, &builder)
	if err != nil {
		return "", err
	}

	return builder.String(), nil
}

func buildTree(dir string, prefix string, currentDepth int, maxDepth int, builder *strings.Builder) error {
	if currentDepth >= maxDepth {
		return nil
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	var filtered []os.DirEntry
	for _, e := range entries {

		if e.Name() == ".git" {
			continue
		}
		filtered = append(filtered, e)
	}

	for i, entry := range filtered {
		isLast := i == len(filtered)-1

		connector := tree.TreeBranch
		if isLast {
			connector = tree.TreeLast
		}

		fmt.Fprintf(builder, "%s%s%s\n", prefix, connector, entry.Name())

		if entry.IsDir() {
			newPrefix := prefix + tree.TreeVertical
			if isLast {
				newPrefix = prefix + tree.TreeEmpty
			}
			if err := buildTree(filepath.Join(dir, entry.Name()), newPrefix, currentDepth+1, maxDepth, builder); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *TreePreviewProvider) Name() string {
	return "tree"
}

func (r *TreePreviewProvider) SetWidth(width int) error {
	r.width = width
	return nil
}
