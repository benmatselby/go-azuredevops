package vsts_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/benmatselby/go-vsts/vsts"
)

const (
	buildDefinitionListURL      = "/VSTS_Project/_apis/build/definitions"
	buildDefinitionListResponse = `{
		"value": [
			{
				"id": 1,
				"name": "build-death-star"
			},
			{
				"id": 2,
				"name": "build-ark"
			}
		],
		"count": 2
	}`
)

func TestBuildDefinitionsService_List(t *testing.T) {
	tt := []struct {
		name     string
		URL      string
		response string
		count    int
		index    int
		defName  string
		defId    int
	}{
		{name: "return two build definitions", URL: buildDefinitionListURL, response: buildDefinitionListResponse, count: 2, index: 0, defName: "build-death-star", defId: 1},
		{name: "can handle no build definitions returned", URL: buildDefinitionListURL, response: "{}", count: 0, index: -1},
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

			options := &vsts.BuildDefinitionsListOptions{}
			buildDefs, err := c.BuildDefinitions.List(options)
			if err != nil {
				t.Fatalf("returned error: %v", err)
			}

			if tc.index > -1 {
				if buildDefs[tc.index].ID != tc.defId {
					t.Fatalf("expected build definition id %d, got %d", tc.defId, buildDefs[tc.index].ID)
				}

				if buildDefs[tc.index].Name != tc.defName {
					t.Fatalf("expected build definition name %s, got %s", tc.defName, buildDefs[tc.index].Name)
				}
			}

			if len(buildDefs) != tc.count {
				t.Fatalf("expected length of build definitions to be %d; got %d", tc.count, len(buildDefs))
			}
		})
	}
}
