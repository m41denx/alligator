package options

type IncludeNodes struct {
	Allocations bool `param:"allocations"` // List of allocations added to the node
	Location    bool `param:"location"`    // Information about the location the node is assigned to
	Servers     bool `param:"servers"`     // List of servers on the node
}

type ListNodesOptions struct {
	requestOptions
	Include IncludeNodes
}

func (o *ListNodesOptions) getOptions() *requestOptions {
	return &requestOptions{
		Include: o.Include,
	}
}

type GetNodeOptions ListNodesOptions

func (o *GetNodeOptions) getOptions() *requestOptions {
	return &requestOptions{
		Include: o.Include,
	}
}

type IncludeAllocations struct {
	Node   bool `param:"node"`   // Information about the node the allocation belongs to
	Server bool `param:"server"` // Information about the server the allocation belongs to
}

type ListNodeAllocationsOptions struct {
	requestOptions
	Include IncludeAllocations
}

func (o *ListNodeAllocationsOptions) getOptions() *requestOptions {
	return &requestOptions{
		Include: o.Include,
	}
}
