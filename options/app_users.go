package options

const (
	ListUsersSort_ID_DESC   = "-id"
	ListUsersSort_ID_ASC    = "id"
	ListUsersSort_UUID_DESC = "-uuid"
	ListUsersSort_UUID_ASC  = "uuid"
)

type IncludeUsers struct {
	Servers bool `param:"servers"` // List of servers the user has access to
}

type FiltersUsers struct {
	Email      string `param:"email"`
	UUID       string `param:"uuid"`
	Username   string `param:"username"`
	ExternalId string `param:"external_id"`
}

type ListUsersOptions struct {
	requestOptions
	Include IncludeUsers
	Filters FiltersUsers
	SortBy  string // -id | id | -uuid | uuid
}

func (o *ListUsersOptions) getOptions() *requestOptions {
	return &requestOptions{
		Include: o.Include,
		Filters: o.Filters,
		SortBy:  o.SortBy,
	}
}

type GetUserOptions struct {
	requestOptions
	Include IncludeUsers
}

func (o *GetUserOptions) getOptions() *requestOptions {
	return &requestOptions{
		Include: o.Include,
	}
}
