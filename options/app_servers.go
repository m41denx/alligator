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

type DatabaseFilters struct {
	Name     string // Not "uwu_database_1"
	ServerID string // Please use real ID, not "1337"
}

type ServerFilters struct {
	Name       string // Server name (no, "mega_server_9000" is not okay)
	ExternalID string // External ID from... somewhere
	UUID       string // Unique ID (yes, more IDs!)
}

// Parameters - the fun continues
type DatabaseParameters struct {
	Host     string // Where the magic happens
	Username string // No, "admin123" is not secure
}

type ServerParameters struct {
	Owner bool // Is it yours? Really?
	Stats bool // How bad is CPU usage?
}

// ListDatabasesOptions holds all available params for listing databases
// Does not include admin's coffee preferences (yet)
type ListDatabasesOptions struct {
	requestOptions
	Include IncludeDatabases
}

// getOptions returns requestOptions from ListDatabasesOptions
// TLDR: yeet the options into standard format
func (o *ListDatabasesOptions) getOptions() *requestOptions {
	return &requestOptions{
		Include:    o.Include,
		Filters:    nil,
		Parameters: nil,
		SortBy:     "",
	}
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

func (o *GetDatabaseOptions) getOptions() *requestOptions {
	return &requestOptions{
		Include: o.Include,
	}
}

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
