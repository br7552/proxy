package main

import (
	"net/http"
	"strings"
)

func getCacheControlHeaders(headers http.Header) map[string]string {
	ccHeaders := make(map[string]string)
	cc := strings.Split(headers.Get("Cache-Control"), ",")
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

	return ccHeaders
}
