package previewproviders

import "github.com/niedch/mux-session/internal/conf"

func CreatePreviewProvider(config *conf.Config, width int) (PreviewProvider, error) {
	providerName := "readme"
	if config.PreviewProvider != nil {
		providerName = *config.PreviewProvider
	}

	var provider PreviewProvider
	var err error

	switch providerName {
	case "readme":
		provider, err = NewReadmePreviewProvider(width)
	case "github":
		provider, err = NewGithubPreviewProvider(width)
	case "tree":
		provider, err = NewTreePreviewProvider(width)
	default:
		provider, err = NewReadmePreviewProvider(width)
	}

	if err != nil {
		return nil, err
	}

	return NewAsyncProviderWrapper(provider), nil
}
