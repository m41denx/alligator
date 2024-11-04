package options

type IncludeServers struct {
	Allocations bool `param:"allocations"` // List of allocations assigned to the server
	User        bool `param:"user"`        // Information about the server owner
	Subusers    bool `param:"subusers"`    // List of users added to the server
	Pack        bool `param:"pack"`        // Information about the server pack
	Nest        bool `param:"nest"`        // Information about the server's egg nest
	Egg         bool `param:"egg"`         // Information about the server's egg
	Variables   bool `param:"variables"`   // List of server variables
	Location    bool `param:"location"`    // Information about server's node location
	Node        bool `param:"node"`        // Information about the server's node
	Databases   bool `param:"databases"`   // List of databases on the server
}

type ListDatabasesOptions struct {
	Include IncludeDatabases // Included Databases
	Fields  string           // Fields in Database
	Filter  string           // Filtering
	Sort    string           // Sorting
	PerPage int              // PerPage (goofy ahh limit)
	Page    int              // Page
}

type IncludeDatabases struct {
	Host bool `param:"host"`
}

func (o *ListDatabasesOptions) GetOptions() *requestOptions {
	return &requestOptions{
		Include: o.Include,
	}
}

type GetDatabaseOptions ListDatabasesOptions

func (o *GetDatabaseOptions) GetOptions() *requestOptions {
	return &requestOptions{
		Include: o.Include,
	}
}

type ListServersOptions struct {
	Include IncludeServers
}

func (o *ListServersOptions) getOptions() *requestOptions {
	return &requestOptions{
		Include: o.Include,
	}
}

type GetServerOptions ListServersOptions

func (o *GetServerOptions) getOptions() *requestOptions {
	return &requestOptions{
		Include: o.Include,
	}
}
