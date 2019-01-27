package azuredevops

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"

	"github.com/google/go-querystring/query"
)

const (
	baseURL   = "https://dev.azure.com/%s"
	userAgent = "go-azuredevops"
)

// Client for interacting with the Azure DevOps API
type Client struct {
	client *http.Client

	BaseURL   string
	UserAgent string

	Account   string
	Project   string
	AuthToken string

	// Services used to proxy to other API endpoints
	Boards           *BoardsService
	BuildDefinitions *BuildDefinitionsService
	Builds           *BuildsService
	DeliveryPlans    *DeliveryPlansService
	Favourites       *FavouritesService
	Git              *GitService
	Iterations       *IterationsService
	PullRequests     *PullRequestsService
	Teams            *TeamsService
	WorkItems        *WorkItemsService
}

// NewClient gets a new Azure DevOps Client
func NewClient(account string, project string, token string) *Client {
	c := &Client{
		Account:   account,
		Project:   project,
		AuthToken: token,
	}
	c.BaseURL = fmt.Sprintf(baseURL, account)

	c.Boards = &BoardsService{client: c}
	c.BuildDefinitions = &BuildDefinitionsService{client: c}
	c.Builds = &BuildsService{client: c}
	c.Favourites = &FavouritesService{client: c}
	c.Git = &GitService{client: c}
	c.Iterations = &IterationsService{client: c}
	c.PullRequests = &PullRequestsService{client: c}
	c.WorkItems = &WorkItemsService{client: c}
	c.Teams = &TeamsService{client: c}
	c.DeliveryPlans = &DeliveryPlansService{client: c}

	return c
}

// NewRequest creates an API request where the URL is relative from https://%s.visualstudio.com/%s.
// Basically this includes the project which is most requests to the API
func (c *Client) NewRequest(method, URL string, body interface{}) (*http.Request, error) {
	request, err := c.NewBaseRequest(
		method,
		fmt.Sprintf("/%s/%s", url.PathEscape(c.Project), URL),
		body,
	)
	return request, err
}

// NewBaseRequest does not take into consideration the project
// and simply uses the base https://%s.visualstudio.com base URL
func (c *Client) NewBaseRequest(method, URL string, body interface{}) (*http.Request, error) {
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	request, err := http.NewRequest(method, c.BaseURL+URL, buf)

	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	if c.UserAgent == "" {
		c.UserAgent = userAgent
	}
	request.Header.Set("User-Agent", c.UserAgent)
	return request, err
}

// Execute runs all the http requests on the API
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

	for k, v := range u.Query() {
		qs[k] = v
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}
