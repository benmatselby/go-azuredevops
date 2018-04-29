package vsts

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"

	"github.com/google/go-querystring/query"
)

const (
	baseURL = "https://%s.visualstudio.com/%s/"
)

// Client for interacting with VSTS
type Client struct {
	client *http.Client

	BaseURL   string
	UserAgent string

	Account   string
	Project   string
	AuthToken string

	// Services used to proxy to other API endpoints
	Iterations   *IterationsService
	WorkItems    *WorkItemsService
	Builds       *BuildsService
	PullRequests *PullRequestsService
	Boards       *BoardsService
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
	c.Builds = &BuildsService{client: c}
	c.PullRequests = &PullRequestsService{client: c}
	c.Boards = &BoardsService{client: c}

	return c
}

// NewRequest creates an API request where the URL is relative from https://%s.visualstudio.com/%s
func (c *Client) NewRequest(method, URL string) (*http.Request, error) {
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

// addOptions adds the parameters in opt as URL query parameters to s. opt
// must be a struct whose fields may contain "url" tags.
// From: https://github.com/google/go-github/blob/master/github/github.go
func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}
