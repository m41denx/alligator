package options

import (
	"testing"
)

func TestOptions(t *testing.T) {
	userOpts := ListUsersOptions{
		Include: IncludeUsers{Servers: true},
		Filters: FiltersUsers{
			Username:   "foo",
			ExternalId: "bar",
		},
		SortBy: ListUsersSort_ID_DESC,
	}

	expected := "filter%5Bexternal_id%5D=bar&filter%5Busername%5D=foo&include=servers&sort=-id"
	if ParseRequestOptions(&userOpts) != expected {
		t.Errorf("expected:\n\t%s,\ngot:\n\t%s", expected, ParseRequestOptions(&userOpts))
	}
}
