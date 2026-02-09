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
	default:
		provider, err = NewReadmePreviewProvider(width)
	}

	return provider, err
}
