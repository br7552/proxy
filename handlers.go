package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

func (p *proxy) handler(w http.ResponseWriter, r *http.Request) {
	addr := "http://" + r.Host + r.URL.Path
	if _, err := url.ParseRequestURI(addr); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest)
		return
	}

	req, err := http.NewRequest(r.Method, addr, r.Body)
	if err != nil {
		log.Printf("proxy:%v\n", err)
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	resp, ok := p.cache.Get(req)
	if !ok {
		resp, err = p.client.Do(req)
		if err != nil {
			log.Printf("proxy:%v\n", err)
			http.Error(w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}

		if req.Method == http.MethodGet ||
			req.Method == http.MethodHead {
			p.cache.Set(req, resp)
		}
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
