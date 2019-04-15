package azuredevops

import (
	"fmt"
	"time"
)

// TestsService handles communication with the Tests methods on the API
// utilising https://docs.microsoft.com/en-gb/rest/api/vsts/test
type TestsService struct {
	client *Client
}

// TestListResponse is the wrapper around the main response for the List of Tests
type TestListResponse struct {
	Tests []Test `json:"value"`
}

// Test represents a test
type Test struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	IsAutomated bool   `json:"isAutomated"`
	Iteration   string `json:"iteration"`
	Owner       *struct {
		ID          string `json:"id"`
		DisplayName string `json:"displayName"`
	} `json:"owner,omitempty"`
	StartedDate   string `json:"startedDate"`
	CompletedDate string `json:"completedDate"`
	State         string `json:"state"`
	Plan          *struct {
		ID string `json:"id"`
	} `json:"plan,omitempty"`
	Revision int `json:"revision"`
}

// TestsListOptions describes what the request to the API should look like
type TestsListOptions struct {
	Count    int    `url:"$top,omitempty"`
	BuildURI string `url:"buildUri,omitempty"`
}

// List returns list of the tests
// utilising https://docs.microsoft.com/en-gb/rest/api/vsts/test/runs/list
func (s *TestsService) List(opts *TestsListOptions) ([]Test, error) {
	URL := fmt.Sprintf("_apis/test/runs?api-version=4.1")
	URL, err := addOptions(URL, opts)

	request, err := s.client.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}
	var response TestListResponse
	_, err = s.client.Execute(request, &response)

	return response.Tests, err
}

// TestResultsListResponse is the wrapper around the main response for the List of Tests
type TestResultsListResponse struct {
	Results []TestResult `json:"value"`
}

// TestResult represents a test result
type TestResult struct {
	ID      int `json:"id"`
	Project struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"project"`
	StartedDate   time.Time `json:"startedDate"`
	CompletedDate time.Time `json:"completedDate"`
	DurationInMs  float64   `json:"durationInMs"`
	Outcome       string    `json:"outcome"`
	Revision      int       `json:"revision"`
	RunBy         struct {
		ID          string `json:"id"`
		DisplayName string `json:"displayName"`
		UniqueName  string `json:"uniqueName"`
		URL         string `json:"url"`
		ImageURL    string `json:"imageUrl"`
	} `json:"runBy"`
	State    string `json:"state"`
	TestCase struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"testCase"`
	TestRun struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"testRun"`
	LastUpdatedDate time.Time `json:"lastUpdatedDate"`
	LastUpdatedBy   struct {
		ID          string `json:"id"`
		DisplayName string `json:"displayName"`
		UniqueName  string `json:"uniqueName"`
		URL         string `json:"url"`
		ImageURL    string `json:"imageUrl"`
	} `json:"lastUpdatedBy"`
	Priority     int    `json:"priority"`
	ComputerName string `json:"computerName"`
	Build        struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"build"`
	CreatedDate          time.Time `json:"createdDate"`
	URL                  string    `json:"url"`
	FailureType          string    `json:"failureType"`
	AutomatedTestStorage string    `json:"automatedTestStorage"`
	AutomatedTestType    string    `json:"automatedTestType"`
	AutomatedTestTypeID  string    `json:"automatedTestTypeId"`
	AutomatedTestID      string    `json:"automatedTestId"`
	Area                 struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"area"`
	TestCaseTitle     string        `json:"testCaseTitle"`
	CustomFields      []interface{} `json:"customFields"`
	AutomatedTestName string        `json:"automatedTestName"`
	StackTrace        string        `json:"stackTrace"`
}

// TestResultsListOptions describes what the request to the API should look like
type TestResultsListOptions struct {
	Count int    `url:"$top,omitempty"`
	RunID string `url:"runId,omitempty"`
}

// ResultsList returns list of the test results
// utilising https://docs.microsoft.com/en-gb/rest/api/vsts/test/runs/list
func (s *TestsService) ResultsList(opts *TestResultsListOptions) ([]TestResult, error) {
	URL := fmt.Sprintf("_apis/test/Runs/%s/results?api-version=4.1", opts.RunID)
	opts.RunID = ""
	URL, err := addOptions(URL, opts)

	request, err := s.client.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}
	var response TestResultsListResponse
	_, err = s.client.Execute(request, &response)

	return response.Results, err
}
