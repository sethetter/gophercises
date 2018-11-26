package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/sethetter/gophercises/urlshort"
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
)

func main() {
	flag.Parse()

	fallback := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		{
			Name: "yaml",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "path",
					Usage:       "path to the yaml config",
					Destination: &yamlFile,
				},
			},
			Action: func(c *cli.Context) error {
				log.Printf("Loading urls from yaml file: %s", yamlFile)
				yml := readFileOrDie(yamlFile)
				urlshortHandler, err := urlshort.YAMLHandler(yml, fallback)
				if err != nil {
					log.Fatalf("Error creating handler: %v", err)
					return err
				}
				serve(urlshortHandler)
				return nil
			},
		},
		{
			Name: "json",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "path",
					Usage:       "path to the json config",
					Destination: &jsonFile,
				},
			},
			Action: func(c *cli.Context) error {
				log.Printf("Loading urls from json file: %s", jsonFile)
				json := readFileOrDie(jsonFile)
				urlshortHandler, err := urlshort.JSONHandler(json, fallback)
				if err != nil {
					log.Fatalf("Error creating handler: %v", err)
					return err
				}
				serve(urlshortHandler)
				return nil
			},
		},
		{
			Name: "db",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "path",
					Usage:       "path to the db file",
					Destination: &dbFile,
				},
			},
			Action: func(c *cli.Context) error {
				log.Printf("Loading urls from db: %s", dbFile)
				db := dbConnectOrDie(dbFile)
				defer db.Close()
				urlshortHandler, err := urlshort.DBHandler(db, fallback)
				if err != nil {
					log.Fatalf("Error creating handler: %v", err)
					return err
				}
				serve(urlshortHandler)
				return nil
			},
			Subcommands: []cli.Command{
				{
					Name:  "add",
					Usage: "add a shorturl to a bolt db",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:        "key",
							Usage:       "short url key",
							Destination: &key,
						},
						cli.StringFlag{
							Name:        "value",
							Usage:       "short url value",
							Destination: &value,
						},
					},
					Action: func(c *cli.Context) error {
						db := dbConnectOrDie(dbFile)
						defer db.Close()
						return db.Update(func(tx *bolt.Tx) error {
							bucket, err := tx.CreateBucketIfNotExists(urlshort.BucketName)
							if err != nil {
								return err
							}
							if err := bucket.Put([]byte(key), []byte(value)); err != nil {
								return err
							}
							return nil
						})
					},
				},
				{
					Name:  "list",
					Usage: "list all the url mappings in a bolt db",
					Action: func(c *cli.Context) error {
						db := dbConnectOrDie(dbFile)
						defer db.Close()
						return db.View(func(tx *bolt.Tx) error {
							bucket := tx.Bucket(urlshort.BucketName)
							if bucket == nil {
								fmt.Println("No urls found!")
								return nil
							}
							c := bucket.Cursor()
							for k, v := c.First(); k != nil; k, v = c.Next() {
								fmt.Printf("%s: %s", k, v)
							}
							return nil
						})
					},
				},
			},
		},
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
