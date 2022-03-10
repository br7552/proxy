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
		url := "http://" + r.Host
		resp, err := client.Get(url)
		if err != nil {
			fmt.Fprintf(w, "proxy:%v\n", err)
			return
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Fprintf(w, "proxy:%v\n", err)
			return
		}

		fmt.Fprintf(w, "%s\n", body)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
