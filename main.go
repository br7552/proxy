package main

import (
	"log"
	"net/http"
	"time"
)

type proxy struct {
	client http.Client
	cache  *RequestCache
}

func main() {
	p := &proxy{
		cache: NewRequestCache(),
	}

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      http.HandlerFunc(p.handler),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
