package main

import (
	"net/http"
	"sync"
)

type cacheKey struct {
	url    string
	method string
}

type RequestCache struct {
	responses map[cacheKey]*http.Response
	mu        sync.Mutex
}

func NewRequestCache() *RequestCache {
	return &RequestCache{
		responses: make(map[cacheKey]*http.Response),
	}
}

func (rc *RequestCache) Get(req *http.Request) (*http.Response, bool) {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	key := cacheKey{
		url:    req.URL.String(),
		method: req.Method,
	}

	resp, ok := rc.responses[key]
	return resp, ok
}

func (rc *RequestCache) Set(req *http.Request, resp *http.Response) {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	key := cacheKey{
		url:    req.URL.String(),
		method: req.Method,
	}

	rc.responses[key] = resp
}
