package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/sethetter/gophercises/urlshort"
)

var (
	port = flag.String("port", "3927", "port to serve on")
)

func main() {
	urlmap := map[string]string{
		"ow": "https://openwichita.org",
		"se": "https://seth.computer",
	}
	fallback := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	urlshortHandler, err := urlshort.MapHandler(urlmap, fallback)
	if err != nil {
		log.Fatalf("Failed to create handler: %v", err)
	}
	http.Handle("/", urlshortHandler)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
