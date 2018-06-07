package vsts

import (
	"fmt"
)

// BuildDefinitionsService handles communication with the build definitions methods on the API
// utilising https://docs.microsoft.com/en-gb/rest/api/vsts/build/definitions
type BuildDefinitionsService struct {
	client *Client
}

// BuildDefinitionsListResponse describes the build definitions list response
type BuildDefinitionsListResponse struct {
	BuildDefinitions []BuildDefinition `json:"value"`
	Count            int               `json:"count"`
}

// BuildDefinition represents a build definition in VSTS
type BuildDefinition struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Repository struct {
		ID   string `json:"id,omitempty"`
		Type string `json:"type,omitempty"`
		Name string `json:"name,omitempty"`
	} `json:"repository,omitempty"`
}

// BuildDefinitionsListOptions describes what the request to the API should look like
type BuildDefinitionsListOptions struct {
	Path                 string `url:"path,omitempty"`
	IncludeAllProperties bool   `url:"includeAllProperties,omitempty"`
}

// List returns a list of build definitions in VSTS
// utilising https://docs.microsoft.com/en-gb/rest/api/vsts/build/definitions/list
func (s *BuildDefinitionsService) List(opts *BuildDefinitionsListOptions) ([]BuildDefinition, error) {
	URL := fmt.Sprintf("_apis/build/definitions?api-version=5.0-preview.6")
	URL, err := addOptions(URL, opts)

	request, err := s.client.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}
	var response BuildDefinitionsListResponse
	_, err = s.client.Execute(request, &response)

	return response.BuildDefinitions, err
}
