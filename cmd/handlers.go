package main

import (
	"io"
	"net/http"
	"net/url"
)

func (p *proxy) handler(w http.ResponseWriter, r *http.Request) {
	var addr string
	switch {
	case r.Host != r.URL.Host:
		addr = "http://" + r.Host + r.URL.Path
	default:
		addr = "http://" + r.URL.Path[1:]
	}

	if _, err := url.ParseRequestURI(addr); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest)
		return
	}

	req, err := http.NewRequest(r.Method, addr, r.Body)
	if err != nil {
		serverErrorResponse(w, r, err)
		return
	}

	// TODO: get maxAge from r cc header
	if resp, ok := p.cache.Get(req, 5); ok {
		writeResponse(w, resp)
		return
	}

	resp, err := p.client.Do(req)
	if err != nil {
		serverErrorResponse(w, r, err)
		return
	}

	// TODO: don't cache if no-store is in r cc header
	if req.Method == http.MethodGet ||
		req.Method == http.MethodHead {
		p.cache.Set(req, resp)
	}

	writeResponse(w, resp)
}

func writeResponse(w http.ResponseWriter, resp *http.Response) {
	for k, v := range resp.Header {
		w.Header()[k] = v
	}
	w.WriteHeader(resp.StatusCode)

	defer resp.Body.Close()
	io.Copy(w, resp.Body)
}
