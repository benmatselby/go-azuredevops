package azuredevops

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

// Build represents a build
type Build struct {
	Status      string          `json:"status"`
	Result      string          `json:"result"`
	Definition  BuildDefinition `json:"definition"`
	BuildNumber string          `json:"buildNumber"`
	FinishTime  string          `json:"finishTime"`
	Branch      string          `json:"sourceBranch"`
}

// BuildsListOptions describes what the request to the API should look like
type BuildsListOptions struct {
	Definitions      string `url:"definitions,omitempty"`
	Branch           string `url:"branchName,omitempty"`
	Count            int    `url:"$top,omitempty"`
	Repository       string `url:"repositoryId,omitempty"`
	BuildIDs         string `url:"buildIds,omitempty"`
	Order            string `url:"queryOrder,omitempty"`
	Deleted          string `url:"deletedFilter,omitempty"`
	MaxPerDefinition string `url:"maxBuildsPerDefinition,omitempty"`
	Token            string `url:"continuationToken,omitempty"`
	Props            string `url:"properties,omitempty"`
	Tags             string `url:"tagFilters,omitempty"`
	Result           string `url:"resultFilter,omitempty"`
	Status           string `url:"statusFilter,omitempty"`
	Reason           string `url:"reasonFilter,omitempty"`
	UserID           string `url:"requestedFor,omitempty"`
	MaxTime          string `url:"maxTime,omitempty"`
	MinTime          string `url:"minTime,omitempty"`
	BuildNumber      string `url:"buildNumber,omitempty"`
	Queues           string `url:"queues,omitempty"`
	RepoType         string `url:"repositoryType,omitempty"`
}

// List returns list of the builds
// utilising https://docs.microsoft.com/en-gb/rest/api/vsts/build/builds/list
func (s *BuildsService) List(opts *BuildsListOptions) ([]Build, error) {
	URL := fmt.Sprintf("/_apis/build/builds?api-version=4.1")
	URL, err := addOptions(URL, opts)

	request, err := s.client.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}
	var response BuildsListResponse
	_, err = s.client.Execute(request, &response)

	return response.Builds, err
}
