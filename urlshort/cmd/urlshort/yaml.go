package main

import (
	"log"

	"github.com/sethetter/gophercises/urlshort"
	"github.com/urfave/cli"
)

func yamlCommand() cli.Command {
	return cli.Command{
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
	}
}
