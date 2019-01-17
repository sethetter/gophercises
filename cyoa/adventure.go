package main

import (
	"encoding/json"
	"io"
)

// Adventure is the top level data structure holding the individual story arcs.
type Adventure struct {
	arcs map[string]Arc
}

// Arc represents a single arc of the story.
type Arc struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Option `json:"options"`
}

// Option is a single option presented to a user at the end of a story arc.
type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

func parseAdventure(f io.Reader) (Adventure, error) {
	var adv Adventure
	if err := json.NewDecoder(f).Decode(&adv.arcs); err != nil {
		return adv, err
	}
	return adv, nil
}

// Runner is the interface for types that can run an adventure.
type Runner = func(Adventure) error
