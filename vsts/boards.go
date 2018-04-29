package vsts

import (
	"fmt"
	"net/url"
)

// BoardsService handles communication with the boards methods on the API
// utilising https://docs.microsoft.com/en-gb/rest/api/vsts/work/boards
type BoardsService struct {
	client *Client
}

// ListBoardsResponse describes the boards response
type ListBoardsResponse struct {
	Boards []Board `json:"value"`
}

// Board describes a board
type Board struct {
	ID      string        `json:"id"`
	Name    string        `json:"name"`
	Columns []BoardColumn `json:"columns,omitempty"`
}

// BoardColumn describes a column on the board
type BoardColumn struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// List returns list of the boards in VSTS
// utilising https://docs.microsoft.com/en-gb/rest/api/vsts/work/boards/list
func (s *BoardsService) List(team string) ([]Board, error) {
	URL := fmt.Sprintf(
		"/%s/_apis/work/boards?api-version=4.1-preview",
		url.PathEscape(team),
	)

	request, err := s.client.NewRequest("GET", URL)
	if err != nil {
		return nil, err
	}
	var response ListBoardsResponse
	_, err = s.client.Execute(request, &response)

	return response.Boards, err
}

// Get returns a single board utilising https://docs.microsoft.com/en-gb/rest/api/vsts/work/boards/get
func (s *BoardsService) Get(team string, id string) (*Board, error) {
	URL := fmt.Sprintf(
		"/%s/_apis/work/boards/%s?api-version=4.1-preview",
		url.PathEscape(team),
		id,
	)

	request, err := s.client.NewRequest("GET", URL)
	if err != nil {
		return nil, err
	}
	var response Board
	_, err = s.client.Execute(request, &response)

	return &response, err
}
