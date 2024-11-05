package alligator

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/m41denx/alligator/options"
)

type Nest struct {
	ID          int          `json:"id,omitempty"`
	UUID        string       `json:"uuid,omitempty"`
	Author      string       `json:"author,omitempty"`
	Name        string       `json:"name,omitempty"`
	Description string       `json:"description,omitempty"`
	CreatedAt   time.Time    `json:"created_at,omitempty"`
	UpdatedAt   time.Time    `json:"updated_at,omitempty"`
	Eggs        []*Egg       `json:"-"`
	Servers     []*AppServer `json:"-"`
}

type Egg struct {
	ID          int            `json:"id,omitempty"`
	UUID        string         `json:"uuid,omitempty"`
	Name        string         `json:"name,omitempty"`
	NestID      int            `json:"nest,omitempty"`
	Author      string         `json:"author,omitempty"`
	Description string         `json:"description,omitempty"`
	DockerImage string         `json:"docker_image,omitempty"`
	Config      EggConfig      `json:"config,omitempty"`
	Startup     string         `json:"startup,omitempty"`
	Script      EggScript      `json:"script,omitempty"`
	CreatedAt   time.Time      `json:"created_at,omitempty"`
	UpdatedAt   time.Time      `json:"updated_at,omitempty"`
	NestObject  *Nest          `json:"-"`
	Servers     []*AppServer   `json:"-"`
	Variables   []*EggVariable `json:"-"`
}

type EggConfig struct {
	Files   map[string]EggFileConfig `json:"files,omitempty"`
	Startup EggStartup               `json:"startup,omitempty"`
	Stop    string                   `json:"stop,omitempty"`
	Logs    EggLogs                  `json:"logs,omitempty"`
	Extends interface{}              `json:"extends,omitempty"`
}

type EggFileConfig struct {
	Parser string            `json:"parser,omitempty"`
	Find   map[string]string `json:"find,omitempty"`
}

type EggStartup struct {
	Done            string   `json:"done,omitempty"`
	UserInteraction []string `json:"userInteraction,omitempty"`
}

// UnmarshalJSON is a custom unmarshaller for EggStartup to handle the situation where "userInteraction" is not present in the JSON.
// If "userInteraction" is not present, it is assumed to be an empty array.
// This custom unmarshaller is necessary because the default unmarshaller for slices in Go will return a nil slice if the slice is not present in the JSON.
// This is a problem because the nil slice is not the same as an empty slice, and the nil slice will not be deserialized properly by the server.
// By using a custom unmarshaller, we can ensure that the slice is always deserialized as an empty slice if it is not present in the JSON.
func (e *EggStartup) UnmarshalJSON(b []byte) error {
	var startup struct {
		Done            string   `json:"done,omitempty"`
		UserInteraction []string `json:"userInteraction,omitempty"`
	}

	if err := json.Unmarshal([]byte(strings.Replace(string(b), "\"userInteraction\":{}", "\"userInteraction\":[]", -1)), &startup); err != nil {
		return err
	}

	e.UserInteraction = startup.UserInteraction
	e.Done = startup.Done

	return nil
}

type EggLogs struct {
	Custom   bool   `json:"custom,omitempty"`
	Location string `json:"location,omitempty"`
}

type EggScript struct {
	Privileged bool        `json:"privileged,omitempty"`
	Install    string      `json:"install,omitempty"`
	Entry      string      `json:"entry,omitempty"`
	Container  string      `json:"container,omitempty"`
	Extends    interface{} `json:"extends,omitempty"`
}

type EggVariable struct {
	ID           int    `json:"id,omitempty"`
	EggID        int    `json:"egg_id,omitempty"`
	Name         string `json:"name,omitempty"`
	Description  string `json:"description,omitempty"`
	EnvVariable  string `json:"env_variable,omitempty"`
	DefaultValue string `json:"default_value,omitempty"`
	UserViewable int    `json:"user_viewable,omitempty"`
	UserEditable int    `json:"user_editable,omitempty"`
	Rules        string `json:"rules,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	UpdatedAt    string `json:"updated_at,omitempty"`
}

type ResponseNest struct {
	*Nest
	Relationships struct {
		Eggs struct {
			Data []struct {
				Attributes *Egg `json:"attributes"`
			} `json:"data"`
		} `json:"eggs"`
		Servers struct {
			Data []struct {
				Attributes *AppServer `json:"attributes"`
			} `json:"data"`
		} `json:"servers"`
	} `json:"relationships"`
}

// getNest returns the nested Nest object, with it's relationships
// resolved from the API response.
func (r *ResponseNest) getNest() *Nest {
	nest := r.Nest
	nest.Eggs = make([]*Egg, 0)
	for _, e := range r.Relationships.Eggs.Data {
		nest.Eggs = append(nest.Eggs, e.Attributes)
	}
	nest.Servers = make([]*AppServer, 0)
	for _, s := range r.Relationships.Servers.Data {
		nest.Servers = append(nest.Servers, s.Attributes)
	}
	return nest
}

// ListNests retrieves a list of Nest objects from the Pterodactyl API, with
// the option to include related servers and eggs.
//
// The opts argument is a variable length argument of options.ListNestsOptions
// structs. These options are used to customize the API request and response.
//
// The function returns a slice of Nest objects, with their related servers and
// eggs resolved. The error return value is used to indicate any errors that
// occurred while executing the request.
func (a *Application) ListNests(opts ...options.ListNestsOptions) ([]*Nest, error) {
	var o string
	if opts != nil && len(opts) > 0 {
		o = options.ParseRequestOptions(&opts[0])
	}
	req := a.newRequest("GET", fmt.Sprintf("/nests?%s", o), nil)
	res, err := a.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model struct {
		Data []struct {
			Attributes *ResponseNest `json:"attributes"`
		} `json:"data"`
	}

	err = json.Unmarshal(buf, &model)
	if err != nil {
		return nil, err
	}

	nests := make([]*Nest, len(model.Data))
	for _, nest := range model.Data {
		nests = append(nests, nest.Attributes.getNest())
	}

	return nests, nil
}

// GetNest returns a Nest object by its ID, with its related servers and eggs
// resolved. The function takes a variable number of options, which are used to
// customize the API request and response. The error return value is used to
// indicate any errors that occurred while executing the request.
func (a *Application) GetNest(nestID int, opts ...options.GetNestOptions) (*Nest, error) {
	var o string
	if opts != nil && len(opts) > 0 {
		o = options.ParseRequestOptions(&opts[0])
	}
	req := a.newRequest("GET", fmt.Sprintf("nests/%d?%s", nestID, o), nil)
	res, err := a.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model struct {
		Attributes *ResponseNest `json:"attributes"`
	}

	err = json.Unmarshal(buf, &model)
	if err != nil {
		return nil, err
	}

	return model.Attributes.getNest(), nil
}

type ResponseEgg struct {
	*Egg
	Relationships struct {
		Nest struct {
			Attributes *Nest `json:"attributes"`
		} `json:"nest"`
		Servers struct {
			Data []struct {
				Attributes *AppServer `json:"attributes"`
			} `json:"data"`
		} `json:"servers"`
		Variables struct {
			Data []struct {
				Attributes *EggVariable `json:"attributes"`
			} `json:"data"`
		} `json:"variables"`
	} `json:"relationships"`
}

// getEgg resolves and returns the Egg object from the ResponseEgg structure.
// It sets up the relationships for Nest, Servers, and Variables by extracting
// the respective attributes from the API response. This function ensures that
// the Egg object is fully populated with its related entities.
func (r *ResponseEgg) getEgg() *Egg {
	egg := r.Egg
	egg.NestObject = r.Relationships.Nest.Attributes
	egg.Servers = make([]*AppServer, 0)
	for _, s := range r.Relationships.Servers.Data {
		egg.Servers = append(egg.Servers, s.Attributes)
	}
	egg.Variables = make([]*EggVariable, 0)
	for _, v := range r.Relationships.Variables.Data {
		egg.Variables = append(egg.Variables, v.Attributes)
	}
	return egg
}

// ListNestEggs retrieves a list of Egg objects from a Nest, with the option to
// include related servers and variables.
//
// The opts argument is a variable length argument of options.ListEggsOptions
// structs. These options are used to customize the API request and response.
//
// The function returns a slice of Egg objects, with their related servers and
// variables resolved. The error return value is used to indicate any errors that
// occurred while executing the request.
func (a *Application) ListNestEggs(nestID int, opts ...options.ListEggsOptions) ([]*Egg, error) {
	var o string
	if opts != nil && len(opts) > 0 {
		o = options.ParseRequestOptions(&opts[0])
	}
	req := a.newRequest("GET", fmt.Sprintf("nests/%d/eggs?%s", nestID, o), nil)
	res, err := a.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model struct {
		Data []struct {
			Attributes *ResponseEgg `json:"attributes"`
		} `json:"data"`
	}

	err = json.Unmarshal(buf, &model)
	if err != nil {
		return nil, err
	}

	eggs := make([]*Egg, len(model.Data))
	for _, egg := range model.Data {
		eggs = append(eggs, egg.Attributes.getEgg())
	}

	return eggs, nil
}

// GetEgg returns a single Egg object by its ID, with its related servers and
// variables resolved. The opts argument is a variable length argument of
// options.GetEggOptions structs. These options are used to customize the API
// request and response.
//
// The function returns a single Egg object, with its related servers and
// variables resolved. The error return value is used to indicate any errors
// that occurred while executing the request.
func (a *Application) GetEgg(nestID, eggID int, opts ...options.GetEggOptions) (*Egg, error) {
	var o string
	if opts != nil && len(opts) > 0 {
		o = options.ParseRequestOptions(&opts[0])
	}
	req := a.newRequest("GET", fmt.Sprintf("nests/%d/eggs?%s", nestID, o), nil)
	res, err := a.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model struct {
		Attributes *ResponseEgg `json:"attributes"`
	}

	err = json.Unmarshal(buf, &model)
	if err != nil {
		return nil, err
	}

	return model.Attributes.getEgg(), nil
}
