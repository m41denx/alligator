package alligator

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/m41denx/alligator/options"
)

type AppServer struct {
	ID            int           `json:"id"`
	ExternalID    string        `json:"external_id"`
	UUID          string        `json:"uuid"`
	Identifier    string        `json:"identifier"`
	Name          string        `json:"name"`
	Description   string        `json:"description"`
	Status        string        `json:"status,omitempty"`
	Suspended     bool          `json:"suspended"`
	Limits        Limits        `json:"limits"`
	FeatureLimits FeatureLimits `json:"feature_limits"`
	UserID        int           `json:"user"`
	NodeID        int           `json:"node"`
	Allocation    int           `json:"allocation"`
	NestID        int           `json:"nest"`
	EggID         int           `json:"egg"`
	Container     struct {
		StartupCommand string                 `json:"startup_command"`
		Image          string                 `json:"image"`
		Installed      int                    `json:"installed"`
		Environment    map[string]interface{} `json:"environment"`
	} `json:"container"`
	CreatedAt   *time.Time     `json:"created_at"`
	UpdatedAt   *time.Time     `json:"updated_at,omitempty"`
	Allocations []*Allocation  `json:"-"`
	UserObject  *User          `json:"-"`
	Subusers    []*User        `json:"-"`
	Location    *Location      `json:"-"`
	NodeObject  *Node          `json:"-"`
	NestObject  *Nest          `json:"-"`
	EggObject   *Egg           `json:"-"`
	Variables   []*EggVariable `json:"-"`
	Databases   []*Database    `json:"-"` // Добавляем поле для баз данных
}

// BuildDescriptor creates a ServerBuildDescriptor for the AppServer.
// It extracts relevant server build details, such as allocation, OOM disabled,
// limits, and feature limits from the server object.
// This descriptor is used to initialize or modify server build parameters.
func (s *AppServer) BuildDescriptor() *ServerBuildDescriptor {
	return &ServerBuildDescriptor{
		Allocation:        s.Allocation,
		OOMDisabled:       s.Limits.OOMDisabled,
		Limits:            s.Limits,
		AddAllocations:    []int{},
		RemoveAllocations: []int{},
		FeatureLimits:     s.FeatureLimits,
	}
}

// DetailsDescriptor creates a ServerDetailsDescriptor for the AppServer.
// It extracts relevant server details, such as the external ID, name, user ID,
// and description from the server object.
// This descriptor is used to initialize or modify server details.
func (s *AppServer) DetailsDescriptor() *ServerDetailsDescriptor {
	return &ServerDetailsDescriptor{
		ExternalID: s.ExternalID,
		// Guess who's back?
		Name:        s.Name,
		User:        s.UserID,
		Description: s.Description,
	}
}

// StartupDescriptor creates a ServerStartupDescriptor for the AppServer.
// It extracts startup-related information, including the startup command,
// environment variables, egg ID, and image from the server's container configuration.
// This descriptor is used to initialize or modify server startup parameters.
func (s *AppServer) StartupDescriptor() *ServerStartupDescriptor {
	return &ServerStartupDescriptor{
		Startup:     s.Container.StartupCommand,
		Environment: s.Container.Environment,
		Egg:         s.EggID,
		Image:       s.Container.Image,
	}
}

// Database represents a database instance for a server
type Database struct {
	ID        int        `json:"id"`
	ServerID  int        `json:"server"`
	HostID    int        `json:"host"`
	Database  string     `json:"database"`
	Username  string     `json:"username"`
	Remote    string     `json:"remote"`
	MaxSize   int        `json:"max_size"`
	Port      int        `json:"port,omitempty"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	Host      *Host      `json:"-"`
	Password  string     `json:"password,omitempty"`
}

// Host represents a database host
type Host struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Host      string     `json:"host"`
	Port      int        `json:"port"`
	Username  string     `json:"username"`
	Node      int        `json:"node"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// ResponseDatabase represents the API response structure for database queries
type ResponseDatabase struct {
	*Database
	Relationships struct {
		Host struct {
			Attributes *Host `json:"attributes"`
		} `json:"host"`
	} `json:"relationships"`
}

// getDatabase returns a Database object extracted from the ResponseDatabase, which
// is a API response structure for database queries. It creates a new Database
// object and copies the fields from the ResponseDatabase, and also sets the Host
// field.
func (r *ResponseDatabase) getDatabase() *Database {
	db := r.Database
	db.Host = r.Relationships.Host.Attributes
	return db
}

// ListDatabases retrieves a list of Database objects associated with a specified server ID
// from the Pterodactyl API. The function supports customization of the API request through
// optional ListDatabasesOptions, allowing the inclusion of related fields and filtering.
//
// Parameters:
//   - serverID: The ID of the server for which databases are to be listed.
//   - opts: Variadic ListDatabasesOptions for customizing the request.
//
// Returns:
//   - A slice of pointers to Database objects representing the databases associated
//     with the specified server.
//   - An error if the request fails or the response cannot be parsed.
func (a *Application) ListDatabases(serverID int, opts ...options.ListDatabasesOptions) ([]*Database, error) {
	var o string
	if opts != nil && len(opts) > 0 {
		o = options.ParseRequestOptions(&opts[0]) // было getOptions
	}
	req := a.newRequest("GET", fmt.Sprintf("/servers/%d/databases?%s", serverID, o), nil)
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
			Attributes *ResponseDatabase `json:"attributes"`
		} `json:"data"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	databases := make([]*Database, 0, len(model.Data))
	for _, d := range model.Data {
		databases = append(databases, d.Attributes.getDatabase())
	}

	return databases, nil
}

// CreateDatabaseOptions represents the options for creating a new database
type CreateDatabaseOptions struct {
	Database string `json:"database"`
	Remote   string `json:"remote"`
	HostID   int    `json:"host,omitempty"`
}

// CreateDatabase creates a new database for a server
func (a *Application) CreateDatabase(serverID int, opts CreateDatabaseOptions) (*Database, error) {
	data, _ := json.Marshal(opts)
	body := bytes.Buffer{}
	body.Write(data)

	req := a.newRequest("POST", fmt.Sprintf("/servers/%d/databases", serverID), &body)
	res, err := a.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model struct {
		Attributes Database `json:"attributes"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	return &model.Attributes, nil
}

// ResetDatabasePassword resets the password for a specific database
func (a *Application) ResetDatabasePassword(serverID, databaseID int) (*Database, error) {
	req := a.newRequest("POST", fmt.Sprintf("/servers/%d/databases/%d/reset-password", serverID, databaseID), nil)
	res, err := a.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model struct {
		Attributes Database `json:"attributes"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	return &model.Attributes, nil
}

// DeleteDatabase deletes a specific database
func (a *Application) DeleteDatabase(serverID, databaseID int) error {
	req := a.newRequest("DELETE", fmt.Sprintf("/servers/%d/databases/%d", serverID, databaseID), nil)
	res, err := a.Http.Do(req)
	if err != nil {
		return err
	}

	_, err = validate(res)
	return err
}

type ResponseServer struct {
	*AppServer
	Relationships struct {
		Allocations struct {
			Data []struct {
				Attributes *Allocation `json:"attributes"`
			} `json:"data"`
		} `json:"allocations"`
		User struct {
			Attributes *User `json:"attributes"`
		} `json:"user"`
		Subusers struct {
			Data []struct {
				Attributes *User `json:"attributes"`
			} `json:"data"`
		} `json:"subusers"`
		Location struct {
			Attributes *Location `json:"attributes"`
		} `json:"location"`
		Node struct {
			Attributes *Node `json:"attributes"`
		} `json:"node"`
		Nest struct {
			Attributes *Nest `json:"attributes"`
		} `json:"nest"`
		Egg struct {
			Attributes *Egg `json:"attributes"`
		} `json:"egg"`
		Variables struct {
			Data []struct {
				Attributes *EggVariable `json:"attributes"`
			} `json:"data"`
		} `json:"variables"`
		Databases struct { // Добавляем структуру для баз данных
			Data []struct {
				Attributes *Database `json:"attributes"`
			} `json:"data"`
		} `json:"databases"`
	} `json:"relationships"`
}

// getServer resolves and returns the Server object from the ResponseServer structure.
// It sets up the relationships for Allocations, User, Subusers, Location, Node, Nest, Egg, Variables, and Databases by extracting
// the respective attributes from the API response. This function ensures that the Server object is fully populated with its related entities.
func (r *ResponseServer) getServer() *AppServer {
	server := r.AppServer
	server.Allocations = make([]*Allocation, 0)
	for _, a := range r.Relationships.Allocations.Data {
		server.Allocations = append(server.Allocations, a.Attributes)
	}
	server.UserObject = r.Relationships.User.Attributes
	server.Subusers = make([]*User, 0)
	for _, s := range r.Relationships.Subusers.Data {
		server.Subusers = append(server.Subusers, s.Attributes)
	}
	server.Location = r.Relationships.Location.Attributes
	server.NodeObject = r.Relationships.Node.Attributes
	server.NestObject = r.Relationships.Nest.Attributes
	server.EggObject = r.Relationships.Egg.Attributes
	server.Variables = make([]*EggVariable, 0)
	for _, v := range r.Relationships.Variables.Data {
		server.Variables = append(server.Variables, v.Attributes)
	}
	server.Databases = make([]*Database, 0) // Инициализируем слайс баз данных
	for _, d := range r.Relationships.Databases.Data {
		server.Databases = append(server.Databases, d.Attributes)
	}
	return server
}

// ListServers retrieves a list of Server objects from the Pterodactyl API, with
// the option to include related allocations, user, subusers, location, node, nest, egg, variables, and databases.
//
// The opts argument is a variable length argument of options.ListServersOptions
// structs. These options are used to customize the API request and response.
//
// The function returns a slice of Server objects, with their related allocations,
// user, subusers, location, node, nest, egg, variables, and databases resolved.
// The error return value is used to indicate any errors that occurred while
// executing the request.
func (a *Application) ListServers(opts ...options.ListServersOptions) ([]*AppServer, error) {
	var o string
	if opts != nil && len(opts) > 0 {
		o = options.ParseRequestOptions(&opts[0])
	}
	req := a.newRequest("GET", fmt.Sprintf("/servers?%s", o), nil)
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
			Attributes *ResponseServer `json:"attributes"`
		} `json:"data"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	servers := make([]*AppServer, 0, len(model.Data))
	for _, s := range model.Data {
		servers = append(servers, s.Attributes.getServer())
	}

	return servers, nil
}

// GetServer retrieves a Server object by its ID, with its related allocations,
// user, subusers, location, node, nest, egg, variables, and databases resolved.
// The function takes a variable number of options, which are used to customize
// the API request and response. The error return value is used to indicate any
// errors that occurred while executing the request.
func (a *Application) GetServer(id int, opts ...options.GetServerOptions) (*AppServer, error) {
	var o string
	if opts != nil && len(opts) > 0 {
		opt := opts[0]
		// Включаем загрузку баз данных
		opt.Include.Databases = true
		o = options.ParseRequestOptions(&opt)
	}
	req := a.newRequest("GET", fmt.Sprintf("/servers/%d?%s", id, o), nil)
	res, err := a.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model struct {
		Attributes ResponseServer `json:"attributes"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	return model.Attributes.getServer(), nil
}

// GetServerExternal retrieves a Server object by its external ID, with its related
// allocations, user, subusers, location, node, nest, egg, variables, and databases
// resolved. The function takes a variable number of options, which are used to
// customize the API request and response. The error return value is used to
// indicate any errors that occurred while executing the request.
func (a *Application) GetServerExternal(id string, opts ...options.GetServerOptions) (*AppServer, error) {
	var o string
	if opts != nil && len(opts) > 0 {
		o = options.ParseRequestOptions(&opts[0])
	}
	req := a.newRequest("GET", fmt.Sprintf("/servers/external/%s?%s", id, o), nil)
	res, err := a.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model struct {
		Attributes ResponseServer `json:"attributes"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	return model.Attributes.getServer(), nil
}

type AllocationDescriptor struct {
	Default    int   `json:"default"`
	Additional []int `json:"additional,omitempty"`
}

type DeployDescriptor struct {
	Locations   []int    `json:"locations"`
	DedicatedIP bool     `json:"dedicated_ip"`
	PortRange   []string `json:"port_range"`
}

type CreateServerDescriptor struct {
	ExternalID        string                 `json:"external_id,omitempty"`
	Name              string                 `json:"name"`
	Description       string                 `json:"description,omitempty"`
	User              int                    `json:"user"`
	Egg               int                    `json:"egg"`
	DockerImage       string                 `json:"docker_image"`
	Startup           string                 `json:"startup"`
	Environment       map[string]interface{} `json:"environment"`
	SkipScripts       bool                   `json:"skip_scripts,omitempty"`
	OOMDisabled       bool                   `json:"oom_disabled"`
	Limits            *Limits                `json:"limits"`
	FeatureLimits     FeatureLimits          `json:"feature_limits"`
	Allocation        *AllocationDescriptor  `json:"allocation,omitempty"`
	Deploy            *DeployDescriptor      `json:"deploy,omitempty"`
	StartOnCompletion bool                   `json:"start_on_completion,omitempty"`
}

// CreateServer creates a new server from the given descriptor. At least one of Allocation or Deploy must be specified.
func (a *Application) CreateServer(fields CreateServerDescriptor) (*AppServer, error) {
	if fields.Allocation == nil && fields.Deploy == nil {
		return nil, errors.New("the allocation object or deploy object must be specified")
	}

	data, _ := json.Marshal(fields)
	body := bytes.Buffer{}
	body.Write(data)

	req := a.newRequest("POST", "/servers", &body)
	res, err := a.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model struct {
		Attributes AppServer `json:"attributes"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	return &model.Attributes, nil
}

type ServerBuildDescriptor struct {
	Allocation        int           `json:"allocation,omitempty"`
	OOMDisabled       bool          `json:"oom_disabled,omitempty"`
	Limits            Limits        `json:"limits,omitempty"`
	AddAllocations    []int         `json:"add_allocations,omitempty"`
	RemoveAllocations []int         `json:"remove_allocations,omitempty"`
	FeatureLimits     FeatureLimits `json:"feature_limits,omitempty"`
}

// UpdateServerBuild updates the build parameters of a server. The fields parameter
// must include at least one of Allocation, OOMDisabled, Limits, AddAllocations,
// RemoveAllocations, or FeatureLimits. The function returns the updated server
// object if successful.
func (a *Application) UpdateServerBuild(id int, fields ServerBuildDescriptor) (*AppServer, error) {
	data, _ := json.Marshal(fields)
	if len(data) == 2 {
		return nil, errors.New("no build fields specified")
	}

	body := bytes.Buffer{}
	body.Write(data)

	req := a.newRequest("PATCH", fmt.Sprintf("/servers/%d/build", id), &body)
	res, err := a.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model struct {
		Attributes AppServer `json:"attributes"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	return &model.Attributes, nil
}

type ServerDetailsDescriptor struct {
	ExternalID  string `json:"external_id,omitempty"`
	Name        string `json:"name,omitempty"`
	User        int    `json:"user,omitempty"`
	Description string `json:"description,omitempty"`
}

// UpdateServerDetails updates the details of a server using the specified server ID and
// details descriptor. It sends a PATCH request to the Pterodactyl API, which updates
// the server's external ID, name, user, or description based on the provided fields.
//
// Parameters:
//   - id: The ID of the server to be updated.
//   - fields: A ServerDetailsDescriptor containing the fields to be updated.
//
// Returns:
//   - A pointer to the updated AppServer object.
//   - An error if the request fails, if the response cannot be parsed, or if no details fields are specified.
func (a *Application) UpdateServerDetails(id int, fields ServerDetailsDescriptor) (*AppServer, error) {
	data, _ := json.Marshal(fields)
	if len(data) == 2 {
		return nil, errors.New("no details fields specified")
	}

	body := bytes.Buffer{}
	body.Write(data)

	req := a.newRequest("PATCH", fmt.Sprintf("/servers/%d/details", id), &body)
	res, err := a.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model struct {
		Attributes AppServer `json:"attributes"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	return &model.Attributes, nil
}

type ServerStartupDescriptor struct {
	Startup     string                 `json:"startup"`
	Environment map[string]interface{} `json:"environment"`
	Egg         int                    `json:"egg,omitempty"`
	Image       string                 `json:"image"`
	SkipScripts bool                   `json:"skip_scripts"`
}

// UpdateServerStartup updates the startup configuration of the server with the specified ID using the provided fields.
//
//   - An integer ID representing the server to be updated.
//   - A ServerStartupDescriptor struct with the fields to be updated.
//   - If no fields are specified, it returns an error.
//   - The function makes a PATCH request to the API, and on success, returns a pointer to the updated AppServer object.
//   - If any errors occur during the request or response processing, an error is returned.
func (a *Application) UpdateServerStartup(id int, fields ServerStartupDescriptor) (*AppServer, error) {
	data, _ := json.Marshal(fields)
	if len(data) == 2 {
		return nil, errors.New("no startup fields specified")
	}

	body := bytes.Buffer{}
	body.Write(data)

	req := a.newRequest("PATCH", fmt.Sprintf("/servers/%d/startup", id), &body)
	res, err := a.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model struct {
		Attributes AppServer `json:"attributes"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	return &model.Attributes, nil
}

// SuspendServer suspends the server with the specified ID, effectively stopping
// it but not deallocating its resources. The function makes a POST request to the
// API, and on success, returns nil. If any errors occur during the request or
// response processing, an error is returned.
func (a *Application) SuspendServer(id int) error {
	req := a.newRequest("POST", fmt.Sprintf("/servers/%d/suspend", id), nil)
	res, err := a.Http.Do(req)
	if err != nil {
		return err
	}

	_, err = validate(res)
	return err
}

// UnsuspendServer resumes the operation of a previously suspended server with the specified ID.
// The function makes a POST request to the API to unsuspend the server, allowing it to run again.
// It returns nil if successful, or an error if any issues occur during the request or response processing.
func (a *Application) UnsuspendServer(id int) error {
	req := a.newRequest("POST", fmt.Sprintf("/servers/%d/unsuspend", id), nil)
	res, err := a.Http.Do(req)
	if err != nil {
		return err
	}

	_, err = validate(res)
	return err
}

// ReinstallServer reinstalls the server with the specified ID, resetting its configuration to the egg defaults.
// The function makes a POST request to the API, and on success, returns nil. If any errors occur during the request or
// response processing, an error is returned.
func (a *Application) ReinstallServer(id int) error {
	req := a.newRequest("POST", fmt.Sprintf("/servers/%d/reinstall", id), nil)
	res, err := a.Http.Do(req)
	if err != nil {
		return err
	}

	_, err = validate(res)
	return err
}

// DeleteServer deletes the server with the specified ID. If the force argument is true, it will delete the server even if it is not in a stopped state. The function makes a DELETE request to the API, and on success, returns nil. If any errors occur during the request or response processing, an error is returned.
func (a *Application) DeleteServer(id int, force bool) error {
	url := fmt.Sprintf("/servers/%d", id)
	if force {
		url += "/force"
	}

	req := a.newRequest("DELETE", url, nil)
	res, err := a.Http.Do(req)
	if err != nil {
		return err
	}

	_, err = validate(res)
	return err
}
