package main

import (
	"flag"
	"log"
)

var (
	storyFile *string = flag.String("story", "gopher.json", "Path to the file containing the adventure.")
	runner    *string = flag.String("runner", "cli", "Which way to run the adventure, 'web' or 'cli'.")
)

func main() {
	flag.Parse()

	// Parse the adventure.
	adventure, err := parseAdventureFile(*storyFile)
	if err != nil {
		log.Fatal(err)
	}

	var runnerF Runner

	switch *runner {
	case "cli":
		runnerF = RunCLI
	case "web":
		runnerF = RunWeb
	}

	if err = runnerF(adventure); err != nil {
		log.Fatal(err)
	}
}
