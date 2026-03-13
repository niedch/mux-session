package previewproviders

import (
	"github.com/niedch/mux-session/internal/dataproviders"
)

// AsyncProviderWrapper wraps a synchronous PreviewProvider to make it asynchronous with caching.
type AsyncProviderWrapper struct {
	inner PreviewProvider
	cache *AsyncCache
}

// NewAsyncProviderWrapper creates a new wrapper around an existing PreviewProvider.
func NewAsyncProviderWrapper(inner PreviewProvider) *AsyncProviderWrapper {
	return &AsyncProviderWrapper{
		inner: inner,
		cache: NewAsyncCache(),
	}
}

// Render returns cached content if available, or spawns a goroutine to fetch it via the inner provider.
func (w *AsyncProviderWrapper) Render(item any) (string, error) {
	dpItem, ok := item.(*dataproviders.Item)
	if !ok {
		// Fallback for non-dataproviders.Item
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

// Name delegates to the inner provider.
func (w *AsyncProviderWrapper) Name() string {
	return w.inner.Name()
}

// SetWidth delegates to the inner provider and clears the cache because word wrapping might change.
func (w *AsyncProviderWrapper) SetWidth(width int) error {
	w.cache.Clear()
	return w.inner.SetWidth(width)
}

// SetUpdateChan sets the channel on the wrapper's cache.
func (w *AsyncProviderWrapper) SetUpdateChan(ch chan<- struct{}) {
	w.cache.SetUpdateChan(ch)
}
