package main

import (
	"fmt"
	"github.com/m41denx/alligator/options"
	"os"

	gator "github.com/m41denx/alligator"
)

func main() {
	app, _ := gator.NewApp(os.Getenv("CROC_URL"), os.Getenv("CROC_KEY"))

	loc, err := app.CreateLocation("us", "United States")
	if err != nil {
		handleError(err)
		return
	}

	fmt.Printf("ID: %d - Short: %s - Long: %s\n", loc.ID, loc.Short, loc.Long)

	loc, err = app.UpdateLocation(loc.ID, "us", "United States of America")
	if err != nil {
		handleError(err)
		return
	}

	fmt.Printf("ID: %d - Short: %s - Long: %s\n", loc.ID, loc.Short, loc.Long)

	// List locations and nodes in them
	locations, err := app.ListLocations(options.ListLocationsOptions{
		Include: options.IncludeLocations{
			Nodes: true,
		}})
	if err != nil {
		fmt.Printf("%#v", err)
		return
	}

	for _, l := range locations {
		fmt.Printf("%d: %s (%s)\n", l.ID, l.Short, l.Long)
	}

	if err = app.DeleteLocation(loc.ID); err != nil {
		handleError(err)
	}
}
