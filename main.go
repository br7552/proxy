package main

import (
	"log"
	"net/http"
	"net/url"
)

type cacheKey struct {
	addr   string
	method string
}

var cache = make(map[cacheKey]*http.Response)
var client http.Client

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		addr := "http://" + r.Host + r.URL.Path
		if _, err := url.ParseRequestURI(addr); err != nil {
			c := http.StatusBadRequest
			http.Error(w, http.StatusText(c), c)
			return
		}

		key := cacheKey{
			addr:   addr,
			method: r.Method,
		}

		if resp, ok := cache[key]; ok {
			resp.Write(w)
			return
		}

		req, err := http.NewRequest(r.Method, addr, r.Body)
		if err != nil {
			log.Printf("proxy:%v\n", err)
			c := http.StatusInternalServerError
			http.Error(w, http.StatusText(c), c)
			return
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("proxy:%v\n", err)
			c := http.StatusInternalServerError
			http.Error(w, http.StatusText(c), c)
			return
		}

		if r.Method == http.MethodGet || r.Method == http.MethodPost {
			cache[key] = resp
		}

		resp.Write(w)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
