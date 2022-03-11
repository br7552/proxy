package main

import (
	"log"
	"net/http"
	"runtime/debug"
)

func serverErrorResponse(w http.ResponseWriter,
	r *http.Request, err error) {

	log.Printf("method: %s, url: %s\n"+
		"error: %s\n"+
		"trace: %s\n",
		r.Method, r.URL.String(), err.Error(), string(debug.Stack()),
	)

	http.Error(w,
		http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError,
	)
}
