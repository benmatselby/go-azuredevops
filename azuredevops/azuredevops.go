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

// DevOpsClient represents methods for interacting with Core level of the Azure DevOps API
// More info: https://docs.microsoft.com/en-us/rest/api/azure/devops/core/
type DevOpsClient struct {
	client *http.Client

	BaseURL   string
	UserAgent string

	Account   string
	AuthToken string

	// Services used to proxy to other API endpoints
	Favourites *FavouritesService
	Projects   *ProjectsService
	Teams      *TeamsService
}

// ProjectClient represents methods for interacting with Project level of the Azure DevOps API
type ProjectClient struct {
	devOpsClient *DevOpsClient

	Project string

	// Services used to proxy to other API endpoints
	Boards           *BoardsService
	BuildDefinitions *BuildDefinitionsService
	Builds           *BuildsService
	DeliveryPlans    *DeliveryPlansService
	Git              *GitService
	Iterations       *IterationsService
	PullRequests     *PullRequestsService
	Tests            *TestsService
	WorkItems        *WorkItemsService
}

// NewProjectClient gets a new Project Client for Azure DevOps API
func NewProjectClient(account string, project string, token string) *ProjectClient {
	devOpsClient := NewDevOpsClient(account, token)

	return devOpsClient.NewProjectClient(project)
}

// NewDevOpsClient gets a new Azure DevOps Client
func NewDevOpsClient(account string, token string) *DevOpsClient {
	c := &DevOpsClient{
		Account:   account,
		AuthToken: token,
		UserAgent: userAgent,
	}

	c.BaseURL = fmt.Sprintf(baseURL, account)
	c.Favourites = &FavouritesService{client: c}
	c.Projects = &ProjectsService{client: c}
	c.Teams = &TeamsService{client: c}

	return c
}

// NewProjectClient gets a new Project Client for Azure DevOps API
func (c *DevOpsClient) NewProjectClient(project string) *ProjectClient {
	client := &ProjectClient{
		devOpsClient: c,
		Project:      project,
	}

	client.Boards = &BoardsService{client: client}
	client.BuildDefinitions = &BuildDefinitionsService{client: client}
	client.Builds = &BuildsService{client: client}
	client.Git = &GitService{client: client}
	client.Iterations = &IterationsService{client: client}
	client.PullRequests = &PullRequestsService{client: client}
	client.WorkItems = &WorkItemsService{client: client}
	client.Tests = &TestsService{client: client}
	client.DeliveryPlans = &DeliveryPlansService{client: client}

	return client
}

// NewRequest creates an API request where the URL is relative from https://dev.azure.com/%s.
// Basically this includes the project which is most requests to the API
func (c *ProjectClient) NewRequest(method, URL string, body interface{}) (*http.Request, error) {
	request, err := c.devOpsClient.NewRequest(
		method,
		fmt.Sprintf("/%s/%s", url.PathEscape(c.Project), URL),
		body,
	)
	return request, err
}

// NewRequest creates an API request where the URL is relative from https://dev.azure.com/%s.
// and simply uses the base https://dev.azure.com base URL
func (c *DevOpsClient) NewRequest(method, URL string, body interface{}) (*http.Request, error) {
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
func (c *ProjectClient) Execute(request *http.Request, r interface{}) (*http.Response, error) {
	return c.devOpsClient.Execute(request, r)
}

// Execute runs all the http requests on the API
func (c *DevOpsClient) Execute(request *http.Request, r interface{}) (*http.Response, error) {
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
