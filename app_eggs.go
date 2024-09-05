package alligator

import (
	"encoding/json"
	"fmt"
	"github.com/m41denx/alligator/options"
	"strings"
	"time"
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
	ID            int          `json:"id,omitempty"`
	UUID          string       `json:"uuid,omitempty"`
	Name          string       `json:"name,omitempty"`
	Nest          int          `json:"nest,omitempty"`
	Author        string       `json:"author,omitempty"`
	Description   string       `json:"description,omitempty"`
	DockerImage   string       `json:"docker_image,omitempty"`
	Config        EggConfig    `json:"config,omitempty"`
	Startup       string       `json:"startup,omitempty"`
	Script        EggScript    `json:"script,omitempty"`
	CreatedAt     time.Time    `json:"created_at,omitempty"`
	UpdatedAt     time.Time    `json:"updated_at,omitempty"`
	Relationships EggRelations `json:"relationships,omitempty"`
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

type EggRelations struct {
	Variables EggVariables `json:"variables,omitempty"`
}

type EggVariables struct {
	Object string            `json:"object,omitempty"`
	Data   []EggRelationData `json:"data,omitempty"`
}

type EggRelationData struct {
	Object     string `json:"object,omitempty"`
	Attributes struct {
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
	} `json:"attributes,omitempty"`
}

func (e *EggVariables) UnmarshalJSON(b []byte) error {
	var eggVariables struct {
		Object string            `json:"object,omitempty"`
		Data   []EggRelationData `json:"data,omitempty"`
	}

	if err := json.Unmarshal(b, &eggVariables); err != nil {
		if eggVariables.Object == "list" {
			e.Data = []EggRelationData{}
		} else {
			return err
		}
	}

	e.Object = eggVariables.Object
	e.Data = eggVariables.Data

	return nil
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

func (a *Application) GetNestEggs(nestID int) ([]*Egg, error) {
	req := a.newRequest("GET", fmt.Sprintf("nests/%d/eggs", nestID), nil)
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
			Attributes *Egg `json:"attributes"`
		} `json:"data"`
	}

	err = json.Unmarshal(buf, &model)
	if err != nil {
		return nil, err
	}

	eggs := make([]*Egg, len(model.Data))
	for _, egg := range model.Data {
		eggs = append(eggs, egg.Attributes)
	}

	return eggs, nil
}

func (a *Application) GetEgg(nestID, eggID int) (*Egg, error) {
	req := a.newRequest("GET", fmt.Sprintf("nests/%d/eggs", nestID), nil)
	res, err := a.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model struct {
		Attributes *Egg `json:"attributes"`
	}

	err = json.Unmarshal(buf, &model)
	if err != nil {
		return nil, err
	}

	return model.Attributes, nil
}
