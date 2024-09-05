package main

import (
	"fmt"
	"os"

	gator "github.com/m41denx/alligator"
)

func main() {
	client, _ := gator.NewClient(os.Getenv("CROC_URL"), os.Getenv("CROC_KEY"))

	servers, err := client.GetServers()
	if err != nil {
		fmt.Printf("%#v", err)
		return
	}

	for _, s := range servers {
		fmt.Printf("%s (%d): %s\n", s.Identifier, s.InternalID, s.Name)
	}

	if len(servers) == 0 {
		fmt.Println("no servers to list")
		return
	}
	server := servers[0]

	if err = client.SetServerPowerState(server.Identifier, "restart"); err != nil {
		fmt.Printf("%#v", err)
		return
	}

	if err = client.SendServerCommand(server.Identifier, "say \"hello world\""); err != nil {
		fmt.Printf("%#v", err)
		return
	}
}
