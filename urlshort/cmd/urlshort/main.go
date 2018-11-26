package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/urfave/cli"
	bolt "go.etcd.io/bbolt"
)

var (
	port     string
	yamlFile string
	jsonFile string
	dbFile   string
	key      string
	value    string
	fallback http.Handler
)

func main() {
	flag.Parse()

	fallback = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "port",
			Value:       "3927",
			Usage:       "port to serve on",
			Destination: &port,
		},
	}
	app.Commands = []cli.Command{
		yamlCommand(),
		jsonCommand(),
		dbCommand(),
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func dbConnectOrDie(path string) *bolt.DB {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		log.Fatalf("Error opening bolt DB: %v", err)
	}
	return db
}

func readFileOrDie(path string) []byte {
	body, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}
	return body
}

func serve(handler http.Handler) {
	http.Handle("/", loggingMiddleware(handler))
	log.Printf("Serving on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
