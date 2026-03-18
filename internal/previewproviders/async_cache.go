package previewproviders

import "sync"

type AsyncCache struct {
	mu         sync.Mutex
	cache      map[string]string
	inProgress map[string]bool
	generation int
	updateChan chan<- struct{}
}

func NewAsyncCache() *AsyncCache {
	return &AsyncCache{
		cache:      make(map[string]string),
		inProgress: make(map[string]bool),
	}
}

func (ac *AsyncCache) SetUpdateChan(ch chan<- struct{}) {
	ac.updateChan = ch
}

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

func (ac *AsyncCache) Finish(key string, result string, generation int) {
	ac.mu.Lock()
	if generation == ac.generation {
		ac.cache[key] = result
		delete(ac.inProgress, key)
	}
	ac.mu.Unlock()

	if ac.updateChan != nil {

		select {
		case ac.updateChan <- struct{}{}:
		default:
		}
	}
}

func (ac *AsyncCache) Clear() {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	ac.cache = make(map[string]string)
	ac.inProgress = make(map[string]bool)
	ac.generation++
}
