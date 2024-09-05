package main

import (
	"fmt"
	"github.com/m41denx/alligator/options"
	"os"

	gator "github.com/m41denx/alligator"
)

func main() {
	app, _ := gator.NewApp(os.Getenv("CROC_URL"), os.Getenv("CROC_KEY"))

	user, err := app.CreateUser(gator.CreateUserDescriptor{
		Email:     "example@example.com",
		Username:  "example",
		FirstName: "test",
		LastName:  "user",
	})
	if err != nil {
		fmt.Printf("%#v", err)
		return
	}

	fmt.Printf("ID: %d - Name: %s - RootAdmin: %v\n", user.ID, user.Username, user.RootAdmin)

	data := user.UpdateDescriptor()
	data.RootAdmin = true
	user, err = app.UpdateUser(user.ID, *data)
	if err != nil {
		fmt.Printf("%#v", err)
		return
	}

	fmt.Printf("ID: %d - Name: %s - RootAdmin: %v\n", user.ID, user.Username, user.RootAdmin)

	users, err := app.ListUsers(options.ListUsersOptions{
		Include: options.IncludeUsers{
			Servers: true,
		},
		Filters: options.FiltersUsers{
			Username: "example",
		},
		SortBy: options.ListUsersSort_UUID_DESC,
	})
	if err != nil {
		fmt.Printf("%#v", err)
		return
	}

	for _, u := range users {
		fmt.Printf("%d: %s\n", u.ID, u.Username)
	}

	if err = app.DeleteUser(user.ID); err != nil {
		fmt.Printf("%#v", err)
	}
}
