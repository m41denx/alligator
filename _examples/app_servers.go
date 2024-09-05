package main

import (
	"fmt"
	"github.com/m41denx/alligator/options"
	"os"

	gator "github.com/m41denx/alligator"
)

func main() {
	app, _ := gator.NewApp(os.Getenv("CROC_URL"), os.Getenv("CROC_KEY"))

	server, err := app.CreateServer(gator.CreateServerDescriptor{
		Name:        "alligator server",
		Description: "test server",
		User:        5,
		Egg:         25,
		DockerImage: "quay.io/parkervcp/pterodactyl-images:base_debian",
		Startup:     "./${EXECUTABLE}",
		Environment: map[string]interface{}{
			"GO_PACKAGE": "github.com/m41denx/alligator",
			"EXECUTABLE": "alligator",
		},
		Limits:        &gator.Limits{Memory: 1024, Disk: 1024, IO: 10, CPU: 1, Threads: "0"},
		FeatureLimits: gator.FeatureLimits{Allocations: 1},
		Deploy:        &gator.DeployDescriptor{Locations: []int{1, 2}, PortRange: []string{}},
	})
	if err != nil {
		fmt.Printf("%#v", err)
		return
	}

	fmt.Printf("ID: %d - Name: %s - ExternalID: %s\n", server.ID, server.Name, server.ExternalID)

	data := server.DetailsDescriptor()
	data.ExternalID = "gator"
	server, err = app.UpdateServerDetails(server.ID, *data)
	if err != nil {
		fmt.Printf("%#v", err)
		return
	}

	fmt.Printf("ID: %d - Name: %s - ExternalID: %s\n", server.ID, server.Name, server.ExternalID)

	servers, err := app.ListServers()
	if err != nil {
		fmt.Printf("%#v", err)
		return
	}

	for _, s := range servers {
		fmt.Printf("%d: %s\n", s.ID, s.Name)
	}

	if err = app.DeleteServer(server.ID, false); err != nil {
		fmt.Printf("%#v", err)
		return
	}

	server2, err := app.GetServer(server.ID, options.GetServerOptions{
		Include: options.IncludeServers{
			User:     true,
			Subusers: true,
			Node:     true,
		},
	})

	if err != nil {
		fmt.Printf("%#v", err)
		return
	}

	fmt.Printf("ID: %d - Name: %s - ExternalID: %s\n", server2.ID, server2.Name, server2.ExternalID)
}
