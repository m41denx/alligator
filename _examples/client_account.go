package main

import (
	"fmt"
	"os"

	gator "github.com/m41denx/alligator"
)

func main() {
	client, _ := gator.NewClient(os.Getenv("CROC_URL"), os.Getenv("CROC_KEY"))

	account, err := client.GetAccount()
	if err != nil {
		fmt.Printf("%#v", err)
		return
	}

	fmt.Printf("ID: %d - Name: %s\n", account.ID, account.FullName())

	apikeys, err := client.GetApiKeys()
	if err != nil {
		fmt.Printf("%#v", err)
		return
	}

	for _, k := range apikeys {
		fmt.Printf("%s: %s\n", k.Identifier, k.Description)
	}
}
