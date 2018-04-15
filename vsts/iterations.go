package vsts

import (
	"fmt"
	"net/url"
)

// IterationsService handles communication with the work items methods on the API
// utilising https://docs.microsoft.com/en-gb/rest/api/vsts/work/iterations
type IterationsService struct {
	client *Client
}

// IterationsResponse describes the iterations response
type IterationsResponse struct {
	Iterations []Iteration `json:"value"`
}

// Iteration describes an iteration
type Iteration struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
	URL  string `json:"url"`
}

// List returns list of the iterations available to the user in VSTS
// utilising https://docs.microsoft.com/en-gb/rest/api/vsts/work/iterations/list
func (s *IterationsService) List(team string) ([]Iteration, error) {
	URL := fmt.Sprintf(
		"%s/_apis/work/teamsettings/iterations?api-version=4.1-preview",
		url.PathEscape(team),
	)

	request, err := s.client.NewRequest("GET", URL)
	if err != nil {
		return nil, err
	}
	var ir IterationsResponse
	_, err = s.client.Execute(request, &ir)

	return ir.Iterations, err
}

// GetByName will search the iterations for the account and project
// and return a single iteration if the names match
func (s *IterationsService) GetByName(team string, name string) (*Iteration, error) {
	iterations, error := s.List(team)
	if error != nil {
		return nil, error
	}

	for index := 0; index < len(iterations); index++ {
		if name == iterations[index].Name {
			iteration := iterations[index]
			return &iteration, nil
		}
	}

	return nil, nil
}
