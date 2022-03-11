package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

/* TODO:
- return bad request response if addr host is proxy server
- consider forwarding request headers and handling cache control
	in a different way
- consider caching original response headers and including them
	in proxy server response
*/

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

	for k, v := range r.Header {
		req.Header[k] = v
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

	if resp, ok := p.cache.Get(req, maxAge); ok {
		fmt.Println("CACHE HIT")
		err = writeResponse(w, resp)
		if err != nil {
			serverErrorResponse(w, r, err)
		}
		return
	}

	resp, err := p.client.Do(req)
	if err != nil {
		serverErrorResponse(w, r, err)
		return
	}

	if _, ok := ccHeaders["no-store"]; !ok &&
		(req.Method == http.MethodGet || req.Method == http.MethodHead) {

		p.cache.Set(req, resp)
	}

	err = writeResponse(w, resp)
	if err != nil {
		serverErrorResponse(w, r, err)
	}
}

func writeResponse(w http.ResponseWriter, resp *http.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	for k, v := range resp.Header {
		w.Header()[k] = v
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(body)

	return nil
}
