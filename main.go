package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

var (
	client http.Client
	cache  = NewRequestCache()
)

func writeResponse(w http.ResponseWriter, resp *http.Response) {
	defer resp.Body.Close()

	for k, v := range resp.Header {
		w.Header()[k] = v
	}
	w.WriteHeader(resp.StatusCode)

	io.Copy(w, resp.Body)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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

		if resp, ok := cache.Get(req); ok {
			writeResponse(w, resp)
			return
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("proxy:%v\n", err)
			http.Error(w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}

		if r.Method == http.MethodGet || r.Method == http.MethodHead {
			cache.Set(req, resp)
		}

		writeResponse(w, resp)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
