package cache

import (
	"net/http"
	"sync"
	"time"
)

type cacheKey struct {
	url    string
	method string
}

type cacheItem struct {
	data  []byte
	added time.Time
}

type RequestCache struct {
	responses map[cacheKey]cacheItem
	mu        sync.Mutex
}

func New() *RequestCache {
	return &RequestCache{
		responses: make(map[cacheKey]cacheItem),
	}
}

func (rc *RequestCache) Get(req *http.Request, maxAge int) ([]byte, bool) {

	if maxAge == 0 {
		return nil, false
	}

	key := cacheKey{
		url:    req.URL.String(),
		method: req.Method,
	}

	rc.mu.Lock()
	defer rc.mu.Unlock()

	item, ok := rc.responses[key]
	if !ok {
		return nil, false
	}

	age := int(time.Since(item.added).Seconds())
	if maxAge > 0 && maxAge < age {
		return nil, false
	}

	return item.data, true
}

func (rc *RequestCache) Set(req *http.Request, data []byte) {
	key := cacheKey{
		url:    req.URL.String(),
		method: req.Method,
	}

	item := cacheItem{
		data:  data,
		added: time.Now(),
	}

	rc.mu.Lock()
	defer rc.mu.Unlock()

	rc.responses[key] = item
}
