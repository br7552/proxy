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
	resp   *http.Response
	maxAge int
	added  time.Time
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

func (rc *RequestCache) Get(req *http.Request) (*http.Response,
	bool) {

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
	if item.maxAge > 0 && item.maxAge < age {
		delete(rc.responses, key)
		return nil, false
	}

	return item.resp, true
}

func (rc *RequestCache) Set(req *http.Request, resp *http.Response,
	maxAge int) {

	if maxAge == 0 {
		return
	}

	key := cacheKey{
		url:    req.URL.String(),
		method: req.Method,
	}

	item := cacheItem{
		resp:   resp,
		added:  time.Now(),
		maxAge: maxAge,
	}

	rc.mu.Lock()
	defer rc.mu.Unlock()

	rc.responses[key] = item
}
