<h1 align="center">
    <img src=".github/alligator.png" width="240" />
</h1>

# 🐊 Alligator Go Client for Pterodactyl
> This fork attempts to follow [Official Pterodactyl API Docs @Dashflo](https://dashflo.net/docs/api/pterodactyl/v1)
> as closely as possible, but is not maintained by Pterodactyl team.
---

## So, what's new?

### ⬆⬇ Support for Filters, Include and Endpoint-specific parameters
- **Example for [List Users](https://dashflo.net/docs/api/pterodactyl/v1/#req_5703791f721f4b50bb0318cf19e2262d) endpoint**
```go
package main
 import (
	 gator "github.com/m41denx/alligator"
	 "github.com/m41denx/alligator/options"
 )

 func main() {
	 app, _ := gator.NewApp("https://panel.pterodactyl.io", "ApplicationKeyYouCreated")
	 
	 // Fetch some users
	 users, err := app.ListUsers(options.ListUsersOptions{
		 Include: options.IncludeUsers{
			 Servers: true,
		 },
		 Filters: options.FiltersUsers{
			 Username: "example",
		 },
		 SortBy: options.ListUsersSort_UUID_DESC, // Same as "-uuid"
	 })
	 if err != nil {
		 fmt.Printf("%#v", err)
		 return
	 }
	 
	 // options are optional btw
	 // users, err := app.ListUsers() // <- is also valid

	 for _, u := range users {
		 fmt.Printf("%d: %s\n", u.ID, u.Username)
	 }
 }
```

**More examples at [📁 _examples](_examples)**