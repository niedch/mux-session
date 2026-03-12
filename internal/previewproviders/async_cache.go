package previewproviders

import "sync"

// AsyncCache handles caching and concurrent access for preview providers
type AsyncCache struct {
	mu         sync.Mutex
	cache      map[string]string
	inProgress map[string]bool
	generation int
	updateChan chan<- struct{}
}

// NewAsyncCache creates a new async cache
func NewAsyncCache() *AsyncCache {
	return &AsyncCache{
		cache:      make(map[string]string),
		inProgress: make(map[string]bool),
	}
}

// SetUpdateChan sets the channel to notify the UI of updates
func (ac *AsyncCache) SetUpdateChan(ch chan<- struct{}) {
	ac.updateChan = ch
}

// GetOrStart checks if a key is cached or currently being fetched.
// It returns (value, true, generation) if it's cached or fetching (where value is "Loading..." if fetching).
// It returns ("", false, generation) if it needs to be fetched, and marks it as in progress.
func (ac *AsyncCache) GetOrStart(key string) (string, bool, int) {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	if cached, found := ac.cache[key]; found {
		return cached, true, ac.generation
	}

	if ac.inProgress[key] {
		return "Loading...", true, ac.generation
	}

	ac.inProgress[key] = true
	return "", false, ac.generation
}

// Finish updates the cache with the fetched result if the generation matches, and clears the in-progress flag.
// It also triggers an update on the updateChan if one is set.
func (ac *AsyncCache) Finish(key string, result string, generation int) {
	ac.mu.Lock()
	if generation == ac.generation {
		ac.cache[key] = result
		delete(ac.inProgress, key)
	}
	ac.mu.Unlock()

	if ac.updateChan != nil {
		// Non-blocking send just in case
		select {
		case ac.updateChan <- struct{}{}:
		default:
		}
	}
}

// Clear empties the cache, in-progress map, and increments the generation counter.
func (ac *AsyncCache) Clear() {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	ac.cache = make(map[string]string)
	ac.inProgress = make(map[string]bool)
	ac.generation++
}
