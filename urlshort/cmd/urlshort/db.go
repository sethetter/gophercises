package main

import (
	"fmt"
	"log"

	"github.com/sethetter/gophercises/urlshort"
	"github.com/urfave/cli"
	bolt "go.etcd.io/bbolt"
)

func dbCommand() cli.Command {
	return cli.Command{
		Name: "db",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "path",
				Usage:       "path to the db file",
				Value:       "dev.db",
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
	}
}
