package azuredevops

import (
	"fmt"
	"time"
)

// GitService handles communication with the git methods on the API
// See: https://docs.microsoft.com/en-us/rest/api/vsts/git/
type GitService struct {
	client *ProjectClient
}

// RepositoriesListResponse describes git list response
type RepositoriesListResponse struct {
	Repositories []Repository `json:"value"`
	Count        int          `json:"count,omitempty"`
}

// Repository describes a git repository
type Repository struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	URL           string `json:"url"`
	DefaultBranch string `json:"defaultBranch"`
	Size          int    `json:"size"`
	RemoteURL     string `json:"remoteUrl"`
	SSHURL        string `json:"sshUrl"`
	WebURL        string `json:"webUrl"`
}

// GitListRefsResponse describes the git refs list response
type GitListRefsResponse struct {
	Count int   `json:"count"`
	Refs  []Ref `json:"value"`
}

// Ref describes what the git reference looks like
type Ref struct {
	Name     string `json:"name,omitempty"`
	ObjectID string `json:"objectId,omitempty"`
	URL      string `json:"url,omitempty"`
	Statuses []struct {
		ID          int    `json:"id,omitempty"`
		State       string `json:"state,omitempty"`
		Description string `json:"description,omitempty"`
		Context     struct {
			Name  string `json:"name,omitempty"`
			Genre string `json:"genre,omitempty"`
		} `json:"context,omitempty"`
		CreationDate time.Time `json:"creationDate,omitempty"`
		CreatedBy    struct {
			ID          string `json:"id,omitempty"`
			DisplayName string `json:"displayName,omitempty"`
			UniqueName  string `json:"uniqueName,omitempty"`
			URL         string `json:"url,omitempty"`
			ImageURL    string `json:"imageUrl,omitempty"`
		} `json:"createdBy,omitempty"`
		TargetURL string `json:"targetUrl,omitempty"`
	} `json:"statuses,omitempty"`
}

// GitRefListOptions describes what the request to the API should look like
type GitRefListOptions struct {
	Filter             string `url:"filter,omitempty"`
	IncludeStatuses    bool   `url:"includeStatuses,omitempty"`
	LatestStatusesOnly bool   `url:"latestStatusesOnly,omitempty"`
}

// List returns a list of the repositories
func (s *GitService) List() ([]Repository, error) {
	URL := "/_apis/git/repositories?api-version=5.0"

	request, err := s.client.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}
	var response RepositoriesListResponse
	_, err = s.client.Execute(request, &response)

	return response.Repositories, err
}

// ListRefs returns a list of the references for a git repo
func (s *GitService) ListRefs(repo, refType string, opts *GitRefListOptions) ([]Ref, int, error) {
	URL := fmt.Sprintf(
		"/_apis/git/repositories/%s/refs/%s?api-version=4.1",
		repo,
		refType,
	)

	URL, err := addOptions(URL, opts)

	request, err := s.client.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, 0, err
	}
	var response GitListRefsResponse
	_, err = s.client.Execute(request, &response)

	return response.Refs, response.Count, err
}
