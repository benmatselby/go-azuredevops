package vsts

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// WorkItemsService handles communication with the work items methods on the API
//
// Docs: https://docs.microsoft.com/en-gb/rest/api/vsts/wit/work%20items
type WorkItemsService struct {
	client *Client
}

// WorkItemsResponse describes the relationships between work items in VSTS
type WorkItemsResponse struct {
	WorkItemRelationships []WorkItemRelationship `json:"workItemRelations"`
}

// WorkItemRelationship describes the workitem section of the response
type WorkItemRelationship struct {
	Target WorkItemRelation `json:"target"`
}

// WorkItemRelation describes an intermediary between iterations and work items
type WorkItemRelation struct {
	ID int `json:"id"`
}

// WorkItemListResponse describes the list reponse for work items
type WorkItemListResponse struct {
	WorkItems []WorkItem `json:"value"`
}

// WorkItem describes an individual work item in TFS
type WorkItem struct {
	ID     int            `json:"id"`
	Rev    int            `json:"rev"`
	Fields WorkItemFields `json:"fields"`
}

// WorkItemFields describes all the fields for a given work item
type WorkItemFields struct {
	ID    int    `json:"System.Id"`
	Title string `json:"System.Title"`
	State string `json:"System.State"`
	Type  string `json:"System.WorkItemType"`
}

// GetForIteration will get a list of work items based on an iteration name
func (s *WorkItemsService) GetForIteration(team string, iteration Iteration) ([]WorkItem, error) {
	queryIds, error := s.GetIdsForIteration(team, iteration)
	if error != nil {
		return nil, error
	}

	var workIds []string
	for index := 0; index < len(queryIds); index++ {
		workIds = append(workIds, strconv.Itoa(queryIds[index]))
	}

	// Now we want to pad out the fields for the work items
	URL := fmt.Sprintf(
		"_apis/wit/workitems?ids=%s&fields=%s&api-version=%s",
		strings.Join(workIds, ","),
		// https://docs.microsoft.com/en-us/rest/api/vsts/wit/work%20item%20types%20field/list
		"System.Id,System.Title,System.State,System.WorkItemType",
		"4.1-preview",
	)

	request, err := s.client.NewRequest("GET", URL)
	if err != nil {
		return nil, err
	}

	var response WorkItemListResponse
	_, err = s.client.Execute(request, &response)

	return response.WorkItems, nil
}

// GetIdsForIteration will return an array of ids for a given iteration
func (s *WorkItemsService) GetIdsForIteration(team string, iteration Iteration) ([]int, error) {
	URL := fmt.Sprintf(
		"%s/_apis/work/teamsettings/iterations/%s/workitems?api-version=%s",
		url.PathEscape(team),
		iteration.ID,
		"4.1-preview",
	)

	request, err := s.client.NewRequest("GET", URL)
	if err != nil {
		return nil, err
	}

	var response WorkItemsResponse
	_, err = s.client.Execute(request, &response)

	var queryIds []int
	for index := 0; index < len(response.WorkItemRelationships); index++ {
		relationship := response.WorkItemRelationships[index]
		queryIds = append(queryIds, relationship.Target.ID)
	}

	return queryIds, nil
}
