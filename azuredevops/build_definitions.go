package azuredevops

import (
	"fmt"
)

// BuildDefinitionsService handles communication with the build definitions methods on the API
// utilising https://docs.microsoft.com/en-gb/rest/api/vsts/build/definitions
type BuildDefinitionsService struct {
	client *ProjectClient
}

// BuildDefinitionsListResponse describes the build definitions list response
type BuildDefinitionsListResponse struct {
	BuildDefinitions []BuildDefinition `json:"value"`
	Count            int               `json:"count"`
}

// Repository represents a repository used by a build definition
type Repository struct {
	ID                 string                 `json:"id,omitempty"`
	Type               string                 `json:"type,omitempty"`
	Name               string                 `json:"name,omitempty"`
	URL                string                 `json:"url,omitempty"`
	RootFolder         string                 `json:"root_folder"`
	Properties         map[string]interface{} `json:"properties"`
	Clean              string                 `json:"clean"`
	DefaultBranch      string                 `json:"default_branch"`
	CheckoutSubmodules bool                   `json:"checkout_submodules"`
}

// BuildDefinition represents a build definition
type BuildDefinition struct {
	ID         int         `json:"id"`
	Name       string      `json:"name"`
	Repository *Repository `json:"repository,omitempty"`
}

// BuildDefinitionsListOptions describes what the request to the API should look like
type BuildDefinitionsListOptions struct {
	Path                 string `url:"path,omitempty"`
	IncludeAllProperties bool   `url:"includeAllProperties,omitempty"`
}

// List returns a list of build definitions
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
