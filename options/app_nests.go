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

type IncludeEggs struct {
	Nest      bool `param:"nest"`      // Information about the nest that owns the egg
	Servers   bool `param:"servers"`   // List of servers using the egg
	Variables bool `param:"variables"` // List of egg variables
}

type ListEggsOptions struct {
	requestOptions
	Include IncludeEggs
}

func (o *ListEggsOptions) getOptions() *requestOptions {
	return &requestOptions{
		Include: o.Include,
	}
}

type GetEggOptions ListEggsOptions

func (o *GetEggOptions) getOptions() *requestOptions {
	return &requestOptions{
		Include: o.Include,
	}
}
