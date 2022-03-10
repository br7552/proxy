package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

var client http.Client

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		url := "http://" + r.Host + r.URL.Path
		req, err := http.NewRequest(r.Method, url, r.Body)
		if err != nil {
			fmt.Fprintf(w, "proxy:%v\n", err)
			c := http.StatusInternalServerError
			http.Error(w, http.StatusText(c), c)
			return
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Fprintf(w, "proxy:%v\n", err)
			c := http.StatusInternalServerError
			http.Error(w, http.StatusText(c), c)
			return
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Fprintf(w, "proxy:%v\n", err)
			c := http.StatusInternalServerError
			http.Error(w, http.StatusText(c), c)
			return
		}

		for k, v := range resp.Header {
			w.Header()[k] = v
		}
		w.WriteHeader(resp.StatusCode)
		w.Write(body)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
