package previewproviders

import (
	"fmt"

	"github.com/niedch/mux-session/internal/dataproviders"
	"github.com/niedch/mux-session/internal/previewproviders/github"
)

type GithubPreviewProvider struct {
	width int
}

func NewGithubPreviewProvider(width int) (*GithubPreviewProvider, error) {
	return &GithubPreviewProvider{
		width: width,
	}, nil
}

func (r *GithubPreviewProvider) Render(item any) (string, error) {
	dpItem, ok := item.(*dataproviders.Item)
	if !ok {
		return "", fmt.Errorf("expected *dataproviders.Item, got %T", item)
	}
	return r.fetch(dpItem)
}

func (r *GithubPreviewProvider) fetch(dpItem *dataproviders.Item) (string, error) {
	info, fallbackMsg, err := github.FetchRepoInfo(dpItem.Path)
	if err != nil {
		return "", err
	}
	if info == nil {
		return fallbackMsg, nil
	}

	return github.RenderUI(info, r.width), nil
}

func (r *GithubPreviewProvider) Name() string {
	return "github"
}

func (r *GithubPreviewProvider) SetWidth(width int) error {
	r.width = width
	return nil
}
