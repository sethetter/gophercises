package main

import (
	"log"

	"github.com/sethetter/gophercises/urlshort"
	"github.com/urfave/cli"
)

func jsonCommand() cli.Command {
	return cli.Command{
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
	}
}
