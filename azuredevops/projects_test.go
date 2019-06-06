package azuredevops_test

import (
	"fmt"
	"net/http"
	"testing"
)

const (
	projectListURL      = "/_apis/projects"
	projectListResponse = `{
		"count": 3,
		"value": [
			{
				"id": "00000000-0000-0000-0000-000000000001",
				"name": "Project 1",
				"description": "Project Description 1",
				"url": "https://dev.azure.com/AZURE_DEVOPS_Account/_apis/projects/00000000-0000-0000-0000-000000000001",
				"state": "wellFormed",
				"revision": 1,
				"visibility": "private",
				"lastUpdateTime": "2019-01-01T01:23:45.678Z"
			},
			{
				"id": "00000000-0000-0000-0000-000000000002",
				"name": "Project 2",
				"description": "Project Description 2",
				"url": "https://dev.azure.com/AZURE_DEVOPS_Account/_apis/projects/00000000-0000-0000-0000-000000000002",
				"state": "wellFormed",
				"revision": 2,
				"visibility": "private",
				"lastUpdateTime": "2019-01-02T01:23:45.678Z"
			},
			{
				"id": "00000000-0000-0000-0000-000000000003",
				"name": "Project 3",
				"description": "Project Description 2",
				"url": "https://dev.azure.com/AZURE_DEVOPS_Account/_apis/projects/00000000-0000-0000-0000-000000000003",
				"state": "wellFormed",
				"revision": 3,
				"visibility": "private",
				"lastUpdateTime": "2019-01-03T01:23:45.678Z"
			}
		]
	}`
)

func TestProjectsService_List(t *testing.T) {
	tt := []struct {
		name           string
		URL            string
		response       string
		count          int
		index          int
		projectName    string
		result         string
		definitionName string
	}{
		{name: "return two boards", URL: projectListURL, response: projectListResponse, count: 3, index: 0, projectName: "Project 1"},
		{name: "can handle no projects returned", URL: projectListURL, response: "{}", count: 0, index: -1},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c, mux, _, teardown := setupDevOpsClient()
			defer teardown()

			mux.HandleFunc(tc.URL, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				json := tc.response
				fmt.Fprint(w, json)
			})

			projects, err := c.Projects.List()
			if err != nil {
				t.Fatalf("returned error: %v", err)
			}

			if tc.index > -1 {
				if projects[tc.index].Name != tc.projectName {
					t.Fatalf("expected board name: %s, got %s", tc.projectName, projects[tc.index].Name)
				}
			}

			if len(projects) != tc.count {
				t.Fatalf("expected length of builds to be %d; got %d", tc.count, len(projects))
			}
		})
	}
}

func TestProjectsService_List_ResponseDecodeFailure(t *testing.T) {
	c, mux, _, teardown := setupDevOpsClient()
	defer teardown()

	mux.HandleFunc(boardListURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		json := "sad"
		fmt.Fprint(w, json)
	})

	_, err := c.Projects.List()
	if err == nil {
		t.Fatalf("expected error decoding the response, did not get one")
	}
}
