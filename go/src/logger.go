package main

import (
	"net/http"
	"time"
)

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		rollingLog.Printf(
			"%s\t%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			r.Host,
			name,
			time.Since(start),
		)
	})
}
