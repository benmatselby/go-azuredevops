package vsts

import (
	"fmt"
)

// BuildsService handles communication with the builds methods on the API
// utilising https://docs.microsoft.com/en-gb/rest/api/vsts/build/builds
type BuildsService struct {
	client *Client
}

// BuildsListResponse is the wrapper around the main response for the List of Builds
type BuildsListResponse struct {
	Builds []Build `json:"value"`
}

// Build represents a build in VSTS
type Build struct {
	Status      string          `json:"status"`
	Result      string          `json:"result"`
	Definition  BuildDefinition `json:"definition"`
	BuildNumber string          `json:"buildNumber"`
	FinishTime  string          `json:"finishTime"`
}

// BuildDefinition represents the `definition` aspect of the response
type BuildDefinition struct {
	Name string `json:"name"`
}

// List returns list of the builds in VSTS
// utilising https://docs.microsoft.com/en-gb/rest/api/vsts/build/builds/list#build
func (s *BuildsService) List() ([]Build, error) {
	URL := fmt.Sprintf("/_apis/build/builds?api-version=4.1")

	request, err := s.client.NewRequest("GET", URL)
	if err != nil {
		return nil, err
	}
	var builds BuildsListResponse
	_, err = s.client.Execute(request, &builds)

	return builds.Builds, err
}
