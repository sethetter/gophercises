package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func RunCLI(adventure Adventure) error {
	current := "intro"
	stdin := bufio.NewReader(os.Stdin)

	for {
		if arc, ok := adventure.arcs[current]; ok {
			// Print the current story arc.
			for _, p := range arc.Story {
				fmt.Println(p + "\n")
			}

			if len(arc.Options) > 0 {
				// Print the options.
				for i, o := range arc.Options {
					fmt.Printf("%d: %s\n", i, o.Text)
				}

				fmt.Print("\nWhat's your choice? ")

				// Collect the response.
				choiceStr, err := stdin.ReadString('\n')
				if err != nil {
					return err
				}

				choice, err := strconv.Atoi(choiceStr[:len(choiceStr)-1])
				if err != nil {
					return err
				}

				if len(arc.Options) > choice {
					c := arc.Options[choice]
					current = c.Arc
				} else {
					fmt.Println("Invalid choice!")
				}

				fmt.Print("\n------------\n\n")
			} else {
				return nil
			}
		}
	}
}
