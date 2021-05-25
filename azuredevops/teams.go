package azuredevops

import "fmt"

// TeamsService handles communication with the teams methods on the API
// utilising https://docs.microsoft.com/en-us/rest/api/vsts/core/teams/get%20all%20teams
type TeamsService struct {
	client *Client
}

// TeamsListResponse describes what the list API call returns for teams
type TeamsListResponse struct {
	Teams []Team `json:"value"`
	Count int    `json:"count,omitempty"`
}

// Team describes what a team looks like
type Team struct {
	ID          string `url:"id,omitempty"`
	Name        string `url:"name,omitempty"`
	URL         string `url:"url,omitempty"`
	Description string `url:"description,omitempty"`
}

// TeamsListOptions describes what the request to the API should look like
type TeamsListOptions struct {
	Mine bool `url:"$mine,omitempty"`
	Top  int  `url:"$top,omitempty"`
	Skip int  `url:"$skip,omitempty"`
}

// List returns list of the teams
func (s *TeamsService) List(opts *TeamsListOptions) ([]Team, int, error) {
	URL := fmt.Sprintf("/_apis/teams?api-version=6.1-preview.3")
	URL, err := addOptions(URL, opts)

	request, err := s.client.NewBaseRequest("GET", URL, nil)
	if err != nil {
		return nil, 0, err
	}
	var response TeamsListResponse
	_, err = s.client.Execute(request, &response)

	return response.Teams, response.Count, err
}
