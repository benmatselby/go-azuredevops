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

// DeliveryPlan describes an plan
type DeliveryPlan struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Created string `json:"createdDate"`
	URL     string `json:"url"`
}

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
