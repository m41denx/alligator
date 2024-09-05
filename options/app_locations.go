package options

type IncludeLocations struct {
	Nodes   bool `param:"nodes"`   // List of nodes assigned to the location
	Servers bool `param:"servers"` // List of servers in the location
}

type ListLocationsOptions struct {
	requestOptions
	Include IncludeLocations
}

func (o *ListLocationsOptions) getOptions() *requestOptions {
	return &requestOptions{
		Include: o.Include,
	}
}

type GetLocationOptions ListLocationsOptions

func (o *GetLocationOptions) getOptions() *requestOptions {
	return &requestOptions{
		Include: o.Include,
	}
}
