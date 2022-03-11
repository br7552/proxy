package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/br7552/proxy/internal/router"
)

func (p *proxy) proxyHandler(w http.ResponseWriter, r *http.Request) {
	addr := getDestination(r)
	if _, err := url.ParseRequestURI(addr); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest)
		return
	}

	req, err := newForwardedRequest(r, addr)
	if err != nil {
		serverErrorResponse(w, r, err)
		return
	}

	if resp, ok := p.cache.Get(req); ok {
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

	ccHeaders := getCacheControlHeaders(resp.Header)

	maxAge := -1
	if t, ok := ccHeaders["max-age"]; ok {
		if age, err := strconv.Atoi(t); nil == err {
			maxAge = age
		}
	}

	if _, ok := ccHeaders["no-store"]; !ok &&
		(req.Method == http.MethodGet || req.Method == http.MethodHead) {

		p.cache.Set(req, resp, maxAge)
	}

	err = writeResponse(w, resp)
	if err != nil {
		serverErrorResponse(w, r, err)
	}
}

func (p *proxy) noCacheHandler(w http.ResponseWriter, r *http.Request) {
	addr := getDestination(r)
	if _, err := url.ParseRequestURI(addr); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest)
		return
	}

	req, err := newForwardedRequest(r, addr)
	if err != nil {
		serverErrorResponse(w, r, err)
		return
	}

	resp, err := p.client.Do(req)
	if err != nil {
		serverErrorResponse(w, r, err)
		return
	}

	err = writeResponse(w, resp)
	if err != nil {
		serverErrorResponse(w, r, err)
	}
}

func getDestination(r *http.Request) string {
	path := router.Param(r, "path")

	var addr string
	switch {
	case r.Host != r.URL.Host:
		addr = "http://" + r.Host + "/" + path
	default:
		addr = "http://" + path
	}

	return addr
}

func newForwardedRequest(r *http.Request, addr string) (*http.Request, error) {
	req, err := http.NewRequest(r.Method, addr, r.Body)
	if err != nil {
		return nil, err
	}

	for k, v := range r.Header {
		req.Header[k] = v
	}
	return req, nil
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
