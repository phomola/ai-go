package main

import (
	"context"
	"fmt"
	"log"

	"github.com/phomola/ai-go/gemini/ai"
)

// Actor contains information about an actor.
type Actor struct {
	Name        string
	YearOfBirth int
	Movies      []struct {
		Name string
		Year int
	}
}

// Print prints the information represented by [Actor].
func (a *Actor) Print() {
	fmt.Println("name:", a.Name)
	fmt.Println("year of birth:", a.YearOfBirth)
	fmt.Println("movies:")
	for _, m := range a.Movies {
		fmt.Println("-", m.Name, m.Year)
	}
}

func main() {
	ctx := context.Background()

	cl, err := ai.NewClient(ctx, ai.Gemini3FlashPreview)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := ai.Generate[Actor](ctx, cl, ai.NewText("Provide information about Brad Pitt with a list of all his movies."))
	if err != nil {
		log.Fatal(err)
	}

	resp.Print()
}
