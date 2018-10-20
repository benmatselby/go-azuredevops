package azuredevops_test

import (
	"encoding/json"
	"fmt"
	"github.com/benmatselby/go-azuredevops/azuredevops"
	"net/http"
	"testing"
)

const (
	buildListURL      = "/AZURE_DEVOPS_Project/_apis/build/builds"
	queueBuildURL     = "/AZURE_DEVOPS_Project/_apis/build/builds"
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

func TestBuildsService_Queue(t *testing.T) {
	t.Run("Updates build struct with returned data", func(t *testing.T) {
		c, mux, _, teardown := setup()
		defer teardown()

		requestBuild := &azuredevops.Build{
			Status: "completed",
			Definition: azuredevops.BuildDefinition{
				Name: "build-one",
			},
		}

		responseBuild := requestBuild
		responseBuild.Result = "succeeded"

		mux.HandleFunc(queueBuildURL, func(w http.ResponseWriter, r *http.Request) {
  		b, err := json.Marshal(responseBuild)
			queueBuildResponse := string(b)

			if err != nil {
				t.Fatalf("returned error: %v", err)
			}

			testMethod(t, r, "POST")
			testBody(t, r, queueBuildResponse+"\n")

			fmt.Fprint(w, queueBuildResponse)
		})

		options := &azuredevops.QueueBuildOptions{}

		err := c.Builds.Queue(requestBuild, options)

		if err != nil {
			t.Fatalf("returned error: %v", err)
		}

		if requestBuild.Result != responseBuild.Result {
			t.Fatalf("expected result %s, got %s", responseBuild.Result, requestBuild.Result)
		}

		if requestBuild.Status != responseBuild.Status {
			t.Fatalf("expected status %s, got %s", responseBuild.Status, requestBuild.Status)
		}

		if requestBuild.Definition.Name != responseBuild.Definition.Name {
			t.Fatalf("expected definition name %s, got %s", responseBuild.Definition.Name, requestBuild.Definition.Name)
		}
	})
}
