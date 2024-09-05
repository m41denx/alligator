package options

type IncludeNests struct {
	Servers bool `param:"servers"` // List of eggs in the location
	Eggs    bool `param:"eggs"`    // List of servers in the location
}

type ListNestsOptions struct {
	requestOptions
	Include IncludeNests
}

func (o *ListNestsOptions) getOptions() *requestOptions {
	return &requestOptions{
		Include: o.Include,
	}
}

type GetNestOptions ListNestsOptions

func (o *GetNestOptions) getOptions() *requestOptions {
	return &requestOptions{
		Include: o.Include,
	}
}
