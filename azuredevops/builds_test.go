package azuredevops_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/benmatselby/go-azuredevops/azuredevops"
)

const (
	buildListURL      = "/AZURE_DEVOPS_Project/_apis/build/builds"
	buildListResponse = `{
		"value": [
			{
				"status": "completed",
				"result": "succeeded",
				"definition": {
					"name": "build-one"
				}
			},
			{
				"status": "completed",
				"result": "failed",
				"definition": {
					"name": "build-two"
				}
			}
		]
	}`
)

func TestBuildsService_List(t *testing.T) {
	tt := []struct {
		name           string
		URL            string
		response       string
		count          int
		index          int
		status         string
		result         string
		definitionName string
	}{
		{name: "return two builds", URL: buildListURL, response: buildListResponse, count: 2, index: 0, status: "completed", result: "succeeded", definitionName: "build-one"},
		{name: "can handle no builds returned", URL: buildListURL, response: "{}", count: 0, index: -1},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c, mux, _, teardown := setup()
			defer teardown()

			mux.HandleFunc(tc.URL, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				json := tc.response
				fmt.Fprint(w, json)
			})

			options := &azuredevops.BuildsListOptions{}
			builds, err := c.Builds.List(options)
			if err != nil {
				t.Fatalf("returned error: %v", err)
			}

			if tc.index > -1 {
				if builds[tc.index].Result != tc.result {
					t.Fatalf("expected result %s, got %s", tc.result, builds[tc.index].Result)
				}

				if builds[tc.index].Status != tc.status {
					t.Fatalf("expected status %s, got %s", tc.status, builds[tc.index].Status)
				}

				if builds[tc.index].Definition.Name != tc.definitionName {
					t.Fatalf("expected definition name %s, got %s", tc.definitionName, builds[tc.index].Definition.Name)
				}
			}

			if len(builds) != tc.count {
				t.Fatalf("expected length of builds to be %d; got %d", tc.count, len(builds))
			}
		})
	}
}
