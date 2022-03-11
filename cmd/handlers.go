package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

	// TODO: return bad request response if addr host is proxy server

	req, err := http.NewRequest(r.Method, addr, r.Body)
	if err != nil {
		serverErrorResponse(w, r, err)
		return
	}

	ccHeaders := make(map[string]string)
	cc := strings.Split(r.Header.Get("Cache-Control"), ",")
	for _, v := range cc {
		v := strings.TrimSpace(v)
		t := strings.Split(v, "=")
		switch len(t) {
		case 2:
			ccHeaders[t[0]] = t[1]
		case 1:
			ccHeaders[t[0]] = ""
		}
	}

	maxAge := -1
	if t, ok := ccHeaders["max-age"]; ok {
		if age, err := strconv.Atoi(t); nil == err {
			maxAge = age
		}
	}

	if body, ok := p.cache.Get(req, maxAge); ok {
		w.Write(body)
		return
	}

	resp, err := p.client.Do(req)
	if err != nil {
		serverErrorResponse(w, r, err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		serverErrorResponse(w, r, err)
		return
	}

	if _, ok := ccHeaders["no-store"]; !ok &&
		(req.Method == http.MethodGet || req.Method == http.MethodHead) {

		p.cache.Set(req, body)
	}

	w.Write(body)
}
