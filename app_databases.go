package alligator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/m41denx/alligator/options"
)

type DatabaseManager struct {
	app *Application
}

// NewDatabaseManager returns a new DatabaseManager instance that allows to manage databases on a Pterodactyl panel.
func NewDatabaseManager(app *Application) *DatabaseManager {
	return &DatabaseManager{app: app}
}

type DatabaseCredentials struct {
	Password string `json:"password"`
}

type DatabaseUsage struct {
	Current int64 `json:"current"`
	Max     int64 `json:"max"`
}

type DatabaseRotatePasswordResponse struct {
	Password string `json:"password"`
}

// HELL YEAH 100 DATABSES IN ONE SERVER
func (dm *DatabaseManager) ListDatabases(serverID int, opts ...options.ListDatabasesOptions) ([]*Database, error) {
	var o string
	if opts != nil && len(opts) > 0 {
		o = options.ParseRequestOptions(opts[0].GetOptions())
	}

	req := dm.app.newRequest("GET", fmt.Sprintf("/servers/%d/databases?%s", serverID, o), nil)
	res, err := dm.app.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model struct {
		Data []struct {
			Attributes *Database `json:"attributes"`
		} `json:"data"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	databases := make([]*Database, 0, len(model.Data))
	for _, d := range model.Data {
		databases = append(databases, d.Attributes)
	}

	return databases, nil
}

// Getting Database
func (dm *DatabaseManager) GetDatabase(serverID, databaseID int, opts ...options.GetDatabaseOptions) (*Database, error) {
	var o string
	if opts != nil && len(opts) > 0 {
		o = options.ParseRequestOptions(opts[0].GetOptions())
	}

	req := dm.app.newRequest("GET", fmt.Sprintf("/servers/%d/databases/%d?%s", serverID, databaseID, o), nil)
	res, err := dm.app.Http.Do(req)
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

// HELL YEAH DATABASES
func (dm *DatabaseManager) CreateDatabase(serverID int, opts CreateDatabaseOptions) (*Database, error) {
	data, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal database options: %w", err)
	}

	body := bytes.Buffer{}
	body.Write(data)

	req := dm.app.newRequest("POST", fmt.Sprintf("/servers/%d/databases", serverID), &body)
	res, err := dm.app.Http.Do(req)
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

// Changing database password
func (dm *DatabaseManager) RotateDatabasePassword(serverID, databaseID int) (*DatabaseRotatePasswordResponse, error) {
	req := dm.app.newRequest("POST", fmt.Sprintf("/servers/%d/databases/%d/rotate-password", serverID, databaseID), nil)
	res, err := dm.app.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model struct {
		Data struct {
			Attributes DatabaseRotatePasswordResponse `json:"attributes"`
		} `json:"data"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	return &model.Data.Attributes, nil
}

// KABOOM DATABASE💥
func (dm *DatabaseManager) DeleteDatabase(serverID, databaseID int) error {
	req := dm.app.newRequest("DELETE", fmt.Sprintf("/servers/%d/databases/%d", serverID, databaseID), nil)
	res, err := dm.app.Http.Do(req)
	if err != nil {
		return err
	}

	_, err = validate(res)
	return err
}

// Getting database usage
func (dm *DatabaseManager) GetDatabaseUsage(serverID, databaseID int) (*DatabaseUsage, error) {
	req := dm.app.newRequest("GET", fmt.Sprintf("/servers/%d/databases/%d/usage", serverID, databaseID), nil)
	res, err := dm.app.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model struct {
		Data struct {
			Attributes DatabaseUsage `json:"attributes"`
		} `json:"data"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	return &model.Data.Attributes, nil
}

// Goofy ahh backup model
type DatabaseBackup struct {
	ID          string     `json:"uuid"`
	Name        string     `json:"name"`
	Size        int64      `json:"bytes"`
	Successful  bool       `json:"successful"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// Creating goofy ahh backups
func (dm *DatabaseManager) CreateDatabaseBackup(serverID, databaseID int, name string) (*DatabaseBackup, error) {
	data := map[string]string{"name": name}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal backup options: %w", err)
	}

	body := bytes.Buffer{}
	body.Write(jsonData)

	req := dm.app.newRequest("POST", fmt.Sprintf("/servers/%d/databases/%d/backup", serverID, databaseID), &body)
	res, err := dm.app.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model struct {
		Data struct {
			Attributes DatabaseBackup `json:"attributes"`
		} `json:"data"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	return &model.Data.Attributes, nil
}

// List of goofy ahh backups
func (dm *DatabaseManager) ListDatabaseBackups(serverID, databaseID int) ([]*DatabaseBackup, error) {
	req := dm.app.newRequest("GET", fmt.Sprintf("/servers/%d/databases/%d/backups", serverID, databaseID), nil)
	res, err := dm.app.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model struct {
		Data []struct {
			Attributes DatabaseBackup `json:"attributes"`
		} `json:"data"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	backups := make([]*DatabaseBackup, 0, len(model.Data))
	for _, b := range model.Data {
		backup := b.Attributes
		backups = append(backups, &backup)
	}

	return backups, nil
}

// Getting Database Backup
func (dm *DatabaseManager) GetDatabaseBackup(serverID, databaseID int, backupID string) (*DatabaseBackup, error) {
	req := dm.app.newRequest("GET", fmt.Sprintf("/servers/%d/databases/%d/backups/%s", serverID, databaseID, backupID), nil)
	res, err := dm.app.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model struct {
		Data struct {
			Attributes DatabaseBackup `json:"attributes"`
		} `json:"data"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	return &model.Data.Attributes, nil
}

// Deleting DatabaseBackup
func (dm *DatabaseManager) DeleteDatabaseBackup(serverID, databaseID int, backupID string) error {
	req := dm.app.newRequest("DELETE", fmt.Sprintf("/servers/%d/databases/%d/backups/%s", serverID, databaseID, backupID), nil)
	res, err := dm.app.Http.Do(req)
	if err != nil {
		return err
	}

	_, err = validate(res)
	return err
}
