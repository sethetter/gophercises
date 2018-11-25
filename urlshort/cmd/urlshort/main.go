package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/sethetter/gophercises/urlshort"
)

var (
	port     = flag.String("port", "3927", "port to serve on")
	yamlFile = flag.String("yaml", "", "path to yaml file with url mappings")
)

func main() {
	flag.Parse()

	fallback := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	var err error
	var urlshortHandler http.Handler
	if *yamlFile != "" {
		var yml []byte
		yml, err = ioutil.ReadFile(*yamlFile)
		if err != nil {
			log.Fatalf("Error reading YAML file: %v", err)
		}
		urlshortHandler, err = urlshort.YAMLHandler(yml, fallback)
	} else {
		urlmap := map[string]string{
			"ow": "https://openwichita.org",
			"se": "https://seth.computer",
		}
		urlshortHandler, err = urlshort.MapHandler(urlmap, fallback)
	}
	if err != nil {
		log.Fatalf("Failed to create default handler: %v", err)
	}

	http.Handle("/", urlshortHandler)
	log.Printf("Serving on port %s\n", *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
