package utils

import (
	"log"
	"net/http"
	"time"
)

type Yconfig struct {
    Mysql struct {
        Host string
	    User string
	    Pw string
	    Port int
	    Schema string
    }
}

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}
