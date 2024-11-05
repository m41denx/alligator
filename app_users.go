package alligator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/m41denx/alligator/options"
)

type User struct {
	ID         int          `json:"id"`
	ExternalID string       `json:"external_id"`
	UUID       string       `json:"uuid"`
	Username   string       `json:"username"`
	Email      string       `json:"email"`
	FirstName  string       `json:"first_name"`
	LastName   string       `json:"last_name"`
	Language   string       `json:"language"`
	RootAdmin  bool         `json:"root_admin"`
	TwoFactor  bool         `json:"2fa"`
	CreatedAt  *time.Time   `json:"created_at"`
	UpdatedAt  *time.Time   `json:"updated_at,omitempty"`
	Servers    []*AppServer `json:"-"`
}

// FullName returns the full name of the user by concatenating the first name
// and last name with a space in between.
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

// UpdateDescriptor returns a descriptor that can be used to update the current
// user. All of the fields on the descriptor are optional and will be ignored
// if they are not provided.
func (u *User) UpdateDescriptor() *UpdateUserDescriptor {
	return &UpdateUserDescriptor{
		ExternalID: u.ExternalID,
		Email:      u.Email,
		Username:   u.Username,
		FirstName:  u.FirstName,
		LastName:   u.LastName,
		Language:   u.Language,
		RootAdmin:  u.RootAdmin,
	}
}

// ListUsers retrieves a list of User objects from the Pterodactyl API, with the
// option to include related servers. The opts argument is a variable length
// argument of options.ListUsersOptions structs, which are used to customize the
// API request and response. The function returns a slice of User objects, with
// their related servers resolved, and an error return value to indicate any
// errors that occurred while executing the request.
func (a *Application) ListUsers(opts ...options.ListUsersOptions) ([]*User, error) {
	var o string
	if opts != nil && len(opts) > 0 {
		o = options.ParseRequestOptions(&opts[0])
	}
	req := a.newRequest("GET", fmt.Sprintf("/users?%s", o), nil)
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
			Attributes struct {
				*User
				Relationships struct {
					Servers struct {
						Data []struct {
							Attributes *AppServer `json:"attributes"`
						} `json:"data"`
					} `json:"servers"`
				} `json:"relationships"`
			} `json:"attributes"`
		} `json:"data"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	users := make([]*User, 0)
	for _, u := range model.Data {
		user := u.Attributes.User
		user.Servers = make([]*AppServer, 0)
		for _, s := range u.Attributes.Relationships.Servers.Data {
			user.Servers = append(user.Servers, s.Attributes)
		}
		users = append(users, user)
	}

	return users, nil
}

// GetUser retrieves a User object by its ID, with its related servers resolved.
// The function takes a variable number of options, which are used to customize
// the API request and response. The error return value is used to indicate any
// errors that occurred while executing the request.
func (a *Application) GetUser(id int, opts ...options.GetUserOptions) (*User, error) {
	var o string
	if opts != nil && len(opts) > 0 {
		o = options.ParseRequestOptions(&opts[0])
	}
	req := a.newRequest("GET", fmt.Sprintf("/users/%d?%s", id, o), nil)
	res, err := a.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model struct {
		Attributes struct {
			*User
			Relationships struct {
				Servers struct {
					Data []struct {
						Attributes *AppServer `json:"attributes"`
					} `json:"data"`
				} `json:"servers"`
			} `json:"relationships"`
		} `json:"attributes"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	user := model.Attributes.User
	user.Servers = make([]*AppServer, 0)
	for _, s := range model.Attributes.Relationships.Servers.Data {
		user.Servers = append(user.Servers, s.Attributes)
	}

	return user, nil
}

// GetUserExternal retrieves a User object by its external ID, with its related
// servers resolved. The function takes a variable number of options, which are
// used to customize the API request and response. The error return value is used
// to indicate any errors that occurred while executing the request.
func (a *Application) GetUserExternal(id string, opts ...options.GetUserOptions) (*User, error) {
	var o string
	if opts != nil && len(opts) > 0 {
		o = options.ParseRequestOptions(&opts[0])
	}
	req := a.newRequest("GET", fmt.Sprintf("/users/external/%s?%s", id, o), nil)
	res, err := a.Http.Do(req)
	if err != nil {
		return nil, err
	}
	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model struct {
		Attributes struct {
			*User
			Relationships struct {
				Servers struct {
					Data []struct {
						Attributes *AppServer `json:"attributes"`
					} `json:"data"`
				} `json:"servers"`
			} `json:"relationships"`
		} `json:"attributes"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	user := model.Attributes.User
	user.Servers = make([]*AppServer, 0)
	for _, s := range model.Attributes.Relationships.Servers.Data {
		user.Servers = append(user.Servers, s.Attributes)
	}

	return user, nil
}

type CreateUserDescriptor struct {
	ExternalID string `json:"external_id,omitempty"`
	Email      string `json:"email"`
	Username   string `json:"username"`
	Password   string `json:"password,omitempty"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Language   string `json:"language,omitempty"`
	RootAdmin  bool   `json:"root_admin,omitempty"`
}

// CreateUser sends a POST request to create a new user with the specified fields.
// The fields parameter is a CreateUserDescriptor containing details about the
// user. The function returns a pointer to the newly created User object, or an
// error if the request fails.
func (a *Application) CreateUser(fields CreateUserDescriptor) (*User, error) {
	data, _ := json.Marshal(fields)
	body := bytes.Buffer{}
	body.Write(data)

	req := a.newRequest("POST", "/users", &body)
	res, err := a.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model struct {
		Attributes User `json:"attributes"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	return &model.Attributes, nil
}

type UpdateUserDescriptor struct {
	ExternalID string `json:"external_id,omitempty"`
	Email      string `json:"email,omitempty"`
	Username   string `json:"username,omitempty"`
	Password   string `json:"password,omitempty"`
	FirstName  string `json:"first_name,omitempty"`
	LastName   string `json:"last_name,omitempty"`
	Language   string `json:"language,omitempty"`
	RootAdmin  bool   `json:"root_admin,omitempty"`
}

// UpdateUser sends a PATCH request to update the user with the specified ID using the provided fields.
// The fields parameter is a UpdateUserDescriptor containing details about the user. The function returns a pointer
// to the updated User object, or an error if the request fails.
func (a *Application) UpdateUser(id int, fields UpdateUserDescriptor) (*User, error) {
	data, _ := json.Marshal(fields)
	body := bytes.Buffer{}
	body.Write(data)

	req := a.newRequest("PATCH", fmt.Sprintf("/users/%d", id), &body)
	res, err := a.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model struct {
		Attributes User `json:"attributes"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	return &model.Attributes, nil
}

// DeleteUser sends a DELETE request to remove the user with the specified ID.
// The function returns an error if the request fails.
func (a *Application) DeleteUser(id int) error {
	req := a.newRequest("DELETE", fmt.Sprintf("/users/%d", id), nil)
	res, err := a.Http.Do(req)
	if err != nil {
		return err
	}

	_, err = validate(res)
	return err
}
