package main

import (
	"log"
	"net/http"
	"time"

	"github.com/br7552/proxy/internal/cache"
	"github.com/br7552/proxy/internal/router"
)

type proxy struct {
	client http.Client
	cache  *cache.RequestCache
}

func main() {
	p := &proxy{
		cache: cache.New(),
	}

	mux := router.New()
	mux.HandleAllFunc("/:path", p.proxyHandler)
	mux.HandleAllFunc("/no-cache/:path", p.noCacheHandler)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
