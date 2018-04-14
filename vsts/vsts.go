package vsts

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	baseURL = "https://%s.visualstudio.com/%s/"
)

// Client for interacting with VSTS
type Client struct {
	client *http.Client

	BaseURL   string
	UserAgent string

	// Services used to proxy to other API endpoints
	Iterations *IterationsService
	WorkItems  *WorkItemsService

	Account   string
	Project   string
	AuthToken string
}

// NewClient gets the VSTS Client
func NewClient(account string, project string, token string) *Client {
	c := &Client{
		Account:   account,
		Project:   project,
		AuthToken: token,
	}
	c.BaseURL = fmt.Sprintf(baseURL, account, url.PathEscape(project))
	c.Iterations = &IterationsService{client: c}
	c.WorkItems = &WorkItemsService{client: c}

	return c
}

// NewRequest creates an API request where the URL is relative from https://%s.visualstudio.com/%s
func (c *Client) NewRequest(method, URL string) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.BaseURL)
	}

	var buf io.ReadWriter

	request, err := http.NewRequest(method, c.BaseURL+URL, buf)
	return request, err
}

// Execute runs all the http requests to VSTS
func (c *Client) Execute(request *http.Request, r interface{}) (*http.Response, error) {
	request.SetBasicAuth("", c.AuthToken)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Request to %s responded with status %d", request.URL, response.StatusCode)
	}

	if err := json.NewDecoder(response.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("Decoding json response from %s failed: %v", request.URL, err)
	}

	return response, nil
}
