package alligator

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const Version = "1.2.0b"

type Application struct {
	PanelURL string
	ApiKey   string
	Http     *http.Client
}

type Client struct {
	PanelURL string
	ApiKey   string
	Http     *http.Client
}

// NewApp creates a new Application instance for interacting with the Pterodactyl API.
// It requires a valid panel URL and an application API key for authentication.
// Returns a pointer to the Application instance or an error if the URL or API key is invalid.
func NewApp(url, key string) (*Application, error) {
	if url == "" {
		return nil, errors.New("a valid panel url is required")
	}
	if key == "" {
		return nil, errors.New("a valid application api key is required")
	}

	app := &Application{
		PanelURL: url,
		ApiKey:   key,
		Http:     &http.Client{},
	}

	return app, nil
}

// newRequest creates a new HTTP request with the given method and path,
// and sets the appropriate headers for an application API request.
// The request is authenticated using the application's API key, and the
// request and response bodies are expected to be in JSON format.
func (a *Application) newRequest(method, path string, body io.Reader) *http.Request {
	req, _ := http.NewRequest(method, fmt.Sprintf("%s/api/application%s", a.PanelURL, path), body)

	req.Header.Set("User-Agent", "Alligator v"+Version)
	req.Header.Set("Authorization", "Bearer "+a.ApiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return req
}

// NewClient returns a new Client instance, given a valid panel URL and client API key.
// The returned Client instance is used to make API requests to the panel on behalf of the client.
// The client API key is used to authenticate the requests, and is required to be set.
// If either the panel URL or the client API key are blank, the function will return an error.
func NewClient(url, key string) (*Client, error) {
	if url == "" {
		return nil, errors.New("a valid panel url is required")
	}
	if key == "" {
		return nil, errors.New("a valid client api key is required")
	}

	client := &Client{
		PanelURL: url,
		ApiKey:   key,
		Http:     &http.Client{},
	}

	return client, nil
}

// newRequest creates a new HTTP request with the given method and path,
// and sets the appropriate headers for a client API request.
// The request is authenticated using the client's API key, and the
// request and response bodies are expected to be in JSON format.
func (c *Client) newRequest(method, path string, body io.Reader) *http.Request {
	req, _ := http.NewRequest(method, fmt.Sprintf("%s/api/client%s", c.PanelURL, path), body)

	req.Header.Set("User-Agent", "Alligator v"+Version)
	req.Header.Set("Authorization", "Bearer "+c.ApiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return req
}

// validate validates the response from the API request. If the status code is 200, 201, or 202, the response body is read and returned as a byte slice.
// If the status code is 204, the function returns (nil, nil). Otherwise, the function attempts to unmarshal the response body into an ApiError struct, and returns (nil, ApiError).
// If there is an error during unmarshalling, the function returns (nil, error).
func validate(res *http.Response) ([]byte, error) {
	switch res.StatusCode {
	case http.StatusOK:
		fallthrough

	case http.StatusCreated:
		fallthrough

	case http.StatusAccepted:
		defer res.Body.Close()
		buf, _ := io.ReadAll(res.Body)
		return buf, nil

	case http.StatusNoContent:
		return nil, nil

	default:
		defer res.Body.Close()
		buf, _ := io.ReadAll(res.Body)

		var errs *ApiError
		if err := json.Unmarshal(buf, &errs); err != nil {
			return nil, err
		}

		return nil, errs
	}
}
