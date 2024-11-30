package options

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

// Базовая структура - просто хранит данные
type requestOptions struct {
	Include    interface{}
	Filters    interface{}
	Parameters interface{}
	SortBy     string
}

type options interface {
	getOptions() *requestOptions
}

func ParseRequestOptions(opts options) string {
	o := opts.getOptions()
	vals := url.Values{}

	// Includes
	{
		t := reflect.TypeOf(o.Include)
		v := reflect.ValueOf(o.Include)
		if t != nil && t.Kind() == reflect.Struct {
			for i := 0; i < t.NumField(); i++ {
				ft := t.Field(i)
				fv := v.Field(i)
				if fv.Interface() == reflect.Zero(ft.Type).Interface() {
					// Skip zero values
					continue
				}
				if fv.Kind() != reflect.Bool {
					// Skip non-bools
					continue
				}
				if fv.Bool() {
					vals.Add("include", ft.Tag.Get("param"))
				}
			}
		}
		if inc := vals.Get("include"); inc != "" {
			vals.Set("include", strings.Join(vals["include"], ","))
		}
	}

	// Filters
	{
		t := reflect.TypeOf(o.Filters)
		v := reflect.ValueOf(o.Filters)
		if t != nil && t.Kind() == reflect.Struct {
			for i := 0; i < t.NumField(); i++ {
				ft := t.Field(i)
				fv := v.Field(i)
				if fv.Interface() == reflect.Zero(ft.Type).Interface() {
					// Skip zero values
					continue
				}
				// If we declare fields ourselves then it's 100% string. In any other case it won't panic anyway
				vals.Set(fmt.Sprintf("filter[%s]", ft.Tag.Get("param")), fv.String())
			}
		}
	}

	//Sort
	if o.SortBy != "" {
		vals.Set("sort", o.SortBy)
	}

	// Other parameters
	{
		t := reflect.TypeOf(o.Parameters)
		v := reflect.ValueOf(o.Parameters)
		if t != nil && t.Kind() == reflect.Struct {
			for i := 0; i < t.NumField(); i++ {
				ft := t.Field(i)
				fv := v.Field(i)
				if fv.Interface() == reflect.Zero(ft.Type).Interface() {
					// Skip zero values
					continue
				}
				vals.Set(ft.Tag.Get("param"), fv.String())
			}
		}
	}

	return vals.Encode()
}
