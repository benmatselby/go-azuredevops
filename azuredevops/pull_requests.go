package azuredevops

import "fmt"

// PullRequestsService handles communication with the pull requests methods on the API
// utilising https://docs.microsoft.com/en-us/rest/api/vsts/git/pull%20requests
type PullRequestsService struct {
	client *Client
}

// PullRequestsResponse describes the pull requests response
type PullRequestsResponse struct {
	PullRequests []PullRequest `json:"value"`
	Count        int           `json:"count"`
}

// PullRequest describes the pull request
type PullRequest struct {
	ID          int             `json:"pullRequestId,omitempty"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Status      string          `json:"status"`
	Created     string          `json:"creationDate"`
	Repo        PullRequestRepo `json:"repository"`
	URL         string          `json:"url"`
}

// PullRequestRepo describes the repo within the pull request
type PullRequestRepo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

// PullRequestListOptions describes what the request to the API should look like
type PullRequestListOptions struct {
	// https://docs.microsoft.com/en-us/rest/api/vsts/git/pull%20requests/get%20pull%20requests%20by%20project#pullrequeststatus
	State string `url:"searchCriteria.status,omitempty"`
}

// List returns list of the pull requests
// utilising https://docs.microsoft.com/en-us/rest/api/vsts/git/pull%20requests/get%20pull%20requests%20by%20project
func (s *PullRequestsService) List(opts *PullRequestListOptions) ([]PullRequest, int, error) {
	URL := fmt.Sprintf("/_apis/git/pullrequests?api-version=4.1")
	URL, err := addOptions(URL, opts)

	request, err := s.client.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, 0, err
	}
	var response PullRequestsResponse
	_, err = s.client.Execute(request, &response)

	return response.PullRequests, response.Count, err
}
