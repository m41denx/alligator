package alligator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/m41denx/alligator/options"
	"time"
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

func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

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

func (a *Application) DeleteUser(id int) error {
	req := a.newRequest("DELETE", fmt.Sprintf("/users/%d", id), nil)
	res, err := a.Http.Do(req)
	if err != nil {
		return err
	}

	_, err = validate(res)
	return err
}
