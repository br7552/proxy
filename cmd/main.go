package main

import (
	"log"
	"net/http"
	"time"

	"github.com/br7552/proxy/internal/cache"
)

type proxy struct {
	client http.Client
	cache  *cache.RequestCache
}

func main() {
	p := &proxy{
		cache: cache.New(),
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
