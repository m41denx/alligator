package alligator

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/m41denx/alligator/options"
)

type Node struct {
	ID                 int           `json:"id"`
	Name               string        `json:"name"`
	Description        string        `json:"description"`
	LocationID         int           `json:"location_id"`
	Public             bool          `json:"public"`
	FQDN               string        `json:"fqdn"`
	Scheme             string        `json:"scheme"`
	BehindProxy        bool          `json:"behind_proxy"`
	Memory             int64         `json:"memory"`
	MemoryOverallocate int64         `json:"memory_overallocate"`
	Disk               int64         `json:"disk"`
	DiskOverallocate   int64         `json:"disk_overallocate"`
	DaemonBase         string        `json:"daemon_base"`
	DaemonSftp         int32         `json:"daemon_sftp"`
	DaemonListen       int32         `json:"daemon_listen"`
	MaintenanceMode    bool          `json:"maintenance_mode"`
	UploadSize         int64         `json:"upload_size"`
	CreatedAt          *time.Time    `json:"created_at"`
	UpdatedAt          *time.Time    `json:"updated_at,omitempty"`
	Location           *Location     `json:"-"`
	Allocations        []*Allocation `json:"-"`
	Servers            []*AppServer  `json:"-"`
}

// UpdateDescriptor returns a descriptor that can be used to update the current
// node. All of the fields on the descriptor are optional and will be ignored
// if they are not provided.
func (n *Node) UpdateDescriptor() *UpdateNodeDescriptor {
	return &UpdateNodeDescriptor{
		Name:               n.Name,
		Description:        n.Description,
		LocationID:         n.LocationID,
		Public:             n.Public,
		FQDN:               n.FQDN,
		Scheme:             n.Scheme,
		BehindProxy:        n.BehindProxy,
		Memory:             n.Memory,
		MemoryOverallocate: n.MemoryOverallocate,
		Disk:               n.Disk,
		DiskOverallocate:   n.DiskOverallocate,
		DaemonBase:         n.DaemonBase,
		DaemonSftp:         n.DaemonSftp,
		DaemonListen:       n.DaemonListen,
		UploadSize:         n.UploadSize,
	}
}

type ResponseNode struct {
	*Node
	Relationships struct {
		Allocations struct {
			Data []struct {
				Attributes *Allocation `json:"attributes"`
			} `json:"data"`
		} `json:"allocations"`
		Location struct {
			Attributes *Location `json:"attributes"`
		} `json:"location"`
		Servers struct {
			Data []struct {
				Attributes *AppServer `json:"attributes"`
			} `json:"data"`
		} `json:"servers"`
	} `json:"relationships"`
}

// getNode resolves and returns the Node object from the ResponseNode structure.
// It sets up the relationships for Allocations, Location, and Servers by extracting
// the respective attributes from the API response. This function ensures that
// the Node object is fully populated with its related entities.
func (r *ResponseNode) getNode() *Node {
	node := r.Node
	node.Allocations = make([]*Allocation, 0)
	for _, a := range r.Relationships.Allocations.Data {
		node.Allocations = append(node.Allocations, a.Attributes)
	}
	node.Location = r.Relationships.Location.Attributes
	node.Servers = make([]*AppServer, 0)
	for _, s := range r.Relationships.Servers.Data {
		node.Servers = append(node.Servers, s.Attributes)
	}
	return node
}

// ListNodes retrieves a list of Node objects from the Pterodactyl API,
// with the option to include related allocations, location, and servers.
// The opts argument is a variable length argument of options.ListNodesOptions
// structs, which are used to customize the API request and response.
// The function returns a slice of Node objects, with their related entities
// resolved, and an error return value to indicate any errors that occurred
// while executing the request.
func (a *Application) ListNodes(opts ...options.ListNodesOptions) ([]*Node, error) {
	var o string
	if opts != nil && len(opts) > 0 {
		o = options.ParseRequestOptions(&opts[0])
	}
	req := a.newRequest("GET", fmt.Sprintf("/nodes?%s", o), nil)
	res, err := a.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	// Holllllly shiii
	var model struct {
		Data []struct {
			Attributes *ResponseNode `json:"attributes"`
		} `json:"data"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	nodes := make([]*Node, 0, len(model.Data))
	for _, n := range model.Data {
		nodes = append(nodes, n.Attributes.getNode())
	}

	return nodes, nil
}

// GetNode retrieves a Node object by its ID, with the option to include related
// allocations, location, and servers. The opts argument is a variable length
// argument of options.GetNodeOptions structs, which are used to customize the API
// request and response. The function returns a single Node object, with its
// related entities resolved, and an error return value to indicate any errors
// that occurred while executing the request.
func (a *Application) GetNode(id int, opts ...options.GetNodeOptions) (*Node, error) {
	var o string
	if opts != nil && len(opts) > 0 {
		o = options.ParseRequestOptions(&opts[0])
	}
	req := a.newRequest("GET", fmt.Sprintf("/nodes/%d?%s", id, o), nil)
	res, err := a.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model struct {
		Attributes *ResponseNode `json:"attributes"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	return model.Attributes.getNode(), nil
}

type DeployableNodesDescriptor struct {
	Page         int   `json:"page,omitempty"`
	Memory       int64 `json:"memory"`
	Disk         int64 `json:"disk"`
	LocationsIDs []int `json:"location_ids,omitempty"`
}

type NodeConfiguration struct {
	Debug   bool   `json:"debug"`
	UUID    string `json:"uuid"`
	TokenID string `json:"token_id"`
	Token   string `json:"token"`
	API     struct {
		Host string `json:"host"`
		Port int32  `json:"port"`
		SSL  struct {
			Enabled bool   `json:"enabled"`
			Cert    string `json:"cert"`
			Key     string `json:"key"`
		} `json:"ssl"`
		UploadLimit int64 `json:"upload_limit"`
	} `json:"api"`
	System struct {
		Data string `json:"data"`
		SFTP struct {
			BindPort int32 `json:"bind_port"`
		} `json:"sftp"`
	} `json:"system"`
	AllowedMounts []string `json:"allowed_mounts"`
	Remote        string   `json:"remote"`
}

// GetNodeConfiguration returns the configuration of the node with the given ID.
// The function makes a GET request to the API, and returns the response as a
// NodeConfiguration object, with an error return value to indicate any errors
// that occurred while executing the request.
func (a *Application) GetNodeConfiguration(id int) (*NodeConfiguration, error) {
	req := a.newRequest("GET", fmt.Sprintf("/nodes/%d/configuration", id), nil)
	res, err := a.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model *NodeConfiguration
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	return model, nil
}

type CreateNodeDescriptor struct {
	Name               string `json:"name"`
	Description        string `json:"description"`
	LocationID         int    `json:"location_id"`
	Public             bool   `json:"public"`
	FQDN               string `json:"fqdn"`
	Scheme             string `json:"scheme"`
	BehindProxy        bool   `json:"behind_proxy"`
	Memory             int64  `json:"memory"`
	MemoryOverallocate int64  `json:"memory_overallocate"`
	Disk               int64  `json:"disk"`
	DiskOverallocate   int64  `json:"disk_overallocate"`
	DaemonBase         string `json:"daemon_base"`
	DaemonSftp         int32  `json:"daemon_sftp"`
	DaemonListen       int32  `json:"daemon_listen"`
	UploadSize         int64  `json:"upload_size"`
}

// CreateNode sends a POST request to create a new node with the specified fields.
// The fields parameter is a CreateNodeDescriptor containing details about the node.
// The function returns a pointer to the newly created Node object, or an error if
// the request fails.
func (a *Application) CreateNode(fields CreateNodeDescriptor) (*Node, error) {
	data, _ := json.Marshal(fields)
	body := bytes.Buffer{}
	body.Write(data)

	req := a.newRequest("POST", "/nodes", &body)
	res, err := a.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model struct {
		Attributes Node `json:"attributes"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	return &model.Attributes, nil
}

type UpdateNodeDescriptor struct {
	Name               string `json:"name"`
	Description        string `json:"description"`
	LocationID         int    `json:"location_id"`
	Public             bool   `json:"public"`
	FQDN               string `json:"fqdn"`
	Scheme             string `json:"scheme"`
	BehindProxy        bool   `json:"behind_proxy"`
	Memory             int64  `json:"memory"`
	MemoryOverallocate int64  `json:"memory_overallocate"`
	Disk               int64  `json:"disk"`
	DiskOverallocate   int64  `json:"disk_overallocate"`
	DaemonBase         string `json:"daemon_base"`
	DaemonSftp         int32  `json:"daemon_sftp"`
	DaemonListen       int32  `json:"daemon_listen"`
	UploadSize         int64  `json:"upload_size"`
}

// UpdateNode updates the node with the specified ID using the provided fields.
// It accepts an integer ID representing the node to be updated and a
// UpdateNodeDescriptor struct with the fields to be updated. If no fields are
// specified, it returns an error. The function makes a PATCH request to the
// API, and on success, returns a pointer to the updated Node object. If any
// errors occur during the request or response processing, an error is returned.
func (a *Application) UpdateNode(id int, fields UpdateNodeDescriptor) (*Node, error) {
	data, _ := json.Marshal(fields)
	if len(data) == 2 {
		return nil, errors.New("no update fields specified")
	}

	body := bytes.Buffer{}
	body.Write(data)

	req := a.newRequest("PATCH", fmt.Sprintf("/nodes/%d", id), &body)
	res, err := a.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model struct {
		Attributes Node `json:"attributes"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	return &model.Attributes, nil
}

// DeleteNode deletes the node with the given ID.
//
// The function takes an integer ID representing the node to be deleted.
// The function returns an error if the request fails.
func (a *Application) DeleteNode(id int) error {
	req := a.newRequest("DELETE", fmt.Sprintf("/nodes/%d", id), nil)
	res, err := a.Http.Do(req)
	if err != nil {
		return err
	}

	_, err = validate(res)
	return err
}

type Allocation struct {
	ID       int        `json:"id"`
	IP       string     `json:"ip"`
	Alias    string     `json:"alias,omitempty"`
	Port     int32      `json:"port"`
	Notes    string     `json:"notes,omitempty"`
	Assigned bool       `json:"assigned"`
	Node     *Node      `json:"-"`
	Server   *AppServer `json:"-"`
}

type ResponseAllocation struct {
	*Allocation
	Relationships struct {
		Node struct {
			Attributes *Node `json:"attributes"`
		} `json:"node"`
		Server struct {
			Attributes *AppServer `json:"attributes"`
		} `json:"server"`
	} `json:"relationships"`
}

// getAllocation returns the nested Allocation object, with it's relationships
// resolved from the API response.

func (r *ResponseAllocation) getAllocation() *Allocation {
	alloc := r.Allocation
	alloc.Node = r.Relationships.Node.Attributes
	alloc.Server = r.Relationships.Server.Attributes
	return alloc
}

// ListNodeAllocations retrieves a list of Allocation objects for a specified node
// from the Pterodactyl API. The function accepts a node ID and an optional
// variable-length argument of options.ListNodeAllocationsOptions structs. These
// options are used to customize the API request and response. The function
// returns a slice of Allocation objects with their relationships resolved, and an
// error return value to indicate any errors that occurred during the request.
func (a *Application) ListNodeAllocations(node int, opts ...options.ListNodeAllocationsOptions) ([]*Allocation, error) {
	var o string
	if opts != nil && len(opts) > 0 {
		o = options.ParseRequestOptions(&opts[0])
	}
	req := a.newRequest("GET", fmt.Sprintf("/nodes/%d/allocations?%s", node, o), nil)
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
			Attributes *ResponseAllocation `json:"attributes"`
		} `json:"data"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	allocs := make([]*Allocation, 0, len(model.Data))
	for _, alloc := range model.Data {
		allocs = append(allocs, alloc.Attributes.getAllocation())
	}

	return allocs, nil
}

type CreateAllocationsDescriptor struct {
	IP    string   `json:"ip"`
	Alias string   `json:"alias,omitempty"`
	Ports []string `json:"ports"`
}

// CreateNodeAllocations creates a new allocation on the specified node. The
// fields argument is a CreateAllocationsDescriptor containing details about the
// allocation. The function returns an error if the request fails.
func (a *Application) CreateNodeAllocations(node int, fields CreateAllocationsDescriptor) error {
	data, _ := json.Marshal(fields)
	body := bytes.Buffer{}
	body.Write(data)

	req := a.newRequest("POST", fmt.Sprintf("/nodes/%d/allocations", node), &body)
	res, err := a.Http.Do(req)
	if err != nil {
		return err
	}

	_, err = validate(res)
	return err
}

// DeleteNodeAllocation deletes the allocation with the given ID from the specified node.
//
// The function takes two integer arguments: the first is the node ID, and the second
// is the allocation ID. The function returns an error if the request fails.
func (a *Application) DeleteNodeAllocation(node, id int) error {
	req := a.newRequest("DELETE", fmt.Sprintf("/nodes/%d/allocations/%d", node, id), nil)
	res, err := a.Http.Do(req)
	if err != nil {
		return err
	}

	_, err = validate(res)
	return err
}
