package previewproviders

import (
	"github.com/niedch/mux-session/internal/dataproviders"
)

type AsyncProviderWrapper struct {
	inner PreviewProvider
	cache *AsyncCache
}

func NewAsyncProviderWrapper(inner PreviewProvider) *AsyncProviderWrapper {
	return &AsyncProviderWrapper{
		inner: inner,
		cache: NewAsyncCache(),
	}
}

func (w *AsyncProviderWrapper) Render(item any) (string, error) {
	dpItem, ok := item.(*dataproviders.Item)
	if !ok {

		return w.inner.Render(item)
	}

	val, inProgress, generation := w.cache.GetOrStart(dpItem.Path)
	if inProgress {
		return val, nil
	}

	go func() {
		res, err := w.inner.Render(item)
		if err != nil {
			res = err.Error()
		}
		w.cache.Finish(dpItem.Path, res, generation)
	}()

	return "Loading...", nil
}

func (w *AsyncProviderWrapper) Name() string {
	return w.inner.Name()
}

func (w *AsyncProviderWrapper) SetWidth(width int) error {
	w.cache.Clear()
	return w.inner.SetWidth(width)
}

func (w *AsyncProviderWrapper) SetUpdateChan(ch chan<- struct{}) {
	w.cache.SetUpdateChan(ch)
}
