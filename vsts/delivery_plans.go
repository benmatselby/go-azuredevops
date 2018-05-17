package vsts

import "fmt"

// DeliveryPlansService handles communication with the deliverytimeline methods on the API
// utilising https://docs.microsoft.com/en-us/rest/api/vsts/work/deliverytimeline
type DeliveryPlansService struct {
	client *Client
}

// DeliveryPlansListResponse describes the delivery plans list response
type DeliveryPlansListResponse struct {
	DeliveryPlans []DeliveryPlan `json:"value"`
	Count         int            `json:"count"`
}

// DeliveryPlanTimeLine describes the delivery plan get response
type DeliveryPlanTimeLine struct {
	StartDate string         `json:"startDate"`
	EndDate   string         `json:"endDate"`
	ID        string         `json:"id"`
	Revision  int            `json:"revision"`
	Teams     []DeliveryTeam `json:"teams"`
}

// DeliveryPlan describes an plan
type DeliveryPlan struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Created string `json:"createdDate"`
	URL     string `json:"url"`
}

// DeliveryTeam describes the teams in a specific plan
type DeliveryTeam struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	Iterations []Iteration `json:"iterations"`
}

const (
	// DeliveryPlanWorkItemIDKey is the key for which part of the workItems[] slice has the ID
	DeliveryPlanWorkItemIDKey = 0
	// DeliveryPlanWorkItemIterationKey is the key for which part of the workItems[] slice has the Iteration
	DeliveryPlanWorkItemIterationKey = 1
	// DeliveryPlanWorkItemTypeKey is the key for which part of the workItems[] slice has the Type
	DeliveryPlanWorkItemTypeKey = 2
	// DeliveryPlanWorkItemNameKey is the key for which part of the workItems[] slice has the Name
	DeliveryPlanWorkItemNameKey = 4
	// DeliveryPlanWorkItemStatusKey is the key for which part of the workItems[] slice has the Status
	DeliveryPlanWorkItemStatusKey = 5
	// DeliveryPlanWorkItemTagKey is the key for which part of the workItems[] slice has the Tag
	DeliveryPlanWorkItemTagKey = 6
)

// DeliveryPlansListOptions describes what the request to the API should look like
type DeliveryPlansListOptions struct {
}

// List returns a list of delivery plans in VSTS
func (s *DeliveryPlansService) List(opts *DeliveryPlansListOptions) ([]DeliveryPlan, int, error) {
	URL := fmt.Sprintf("_apis/work/plans?api-version=5.0-preview.6")
	URL, err := addOptions(URL, opts)

	request, err := s.client.NewRequest("GET", URL)
	if err != nil {
		return nil, 0, err
	}
	var response DeliveryPlansListResponse
	_, err = s.client.Execute(request, &response)

	return response.DeliveryPlans, response.Count, err
}

// GetTimeLine will fetch the details about a specific delivery plan
func (s *DeliveryPlansService) GetTimeLine(ID string) (*DeliveryPlanTimeLine, error) {
	URL := fmt.Sprintf(
		"_apis/work/plans/%s/deliverytimeline?api-version=5.0-preview.1",
		ID,
	)

	request, err := s.client.NewRequest("GET", URL)
	if err != nil {
		return nil, err
	}
	var response DeliveryPlanTimeLine
	_, err = s.client.Execute(request, &response)

	return &response, err
}
