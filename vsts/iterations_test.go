package vsts_test

import (
	"fmt"
	"net/http"
	"testing"
)

const (
	listURL      = "/VSTS_Project/VSTS_TEAM/_apis/work/teamsettings/iterations"
	listResponse = `{
		"value": [
		{
			"id": "a589a806-bf11-4d4f-a031-c19813331553",
			"name": "Sprint 1",
			"path": "",
			"url": "https://example.com/1"
		},
		{
			"id": "a589a806-bf11-4d4f-a031-c19813331554",
			"name": "Sprint 2",
			"path": "",
			"url": "https://example.com/2"
		}
		]
	}`
)

func TestIterationService_GetByName(t *testing.T) {
	tt := []struct {
		name      string
		iteration string
		found     bool
	}{
		{name: "iteration found", iteration: "Sprint 2", found: true},
		{name: "iteration not found", iteration: "Sprint Unkown", found: false},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c, mux, _, teardown := setup()
			defer teardown()

			mux.HandleFunc(listURL, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				json := listResponse
				fmt.Fprint(w, json)
			})

			iteration, err := c.Iterations.GetByName("VSTS_TEAM", tc.iteration)
			if err != nil {
				t.Fatalf("returned error: %v", err)
			}

			if tc.found {
				if iteration.Name != tc.iteration {
					t.Fatalf("expected name iteration name %s; got %s", tc.iteration, iteration.Name)
				}
			}

			if !tc.found {
				if iteration != nil {
					t.Fatalf("We should not have matched an iteration")
				}
			}
		})
	}
}

func TestIterationService_List(t *testing.T) {
	tt := []struct {
		name           string
		URL            string
		response       string
		iterationName  string
		iterationIndex int
	}{
		{name: "the first item is sprint 1", URL: listURL, iterationName: "Sprint 1", iterationIndex: 0, response: listResponse},
		{name: "there are no items in the response", URL: listURL, iterationIndex: -1, response: "{}"},
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

			iterations, err := c.Iterations.List("VSTS_TEAM")
			if err != nil {
				t.Fatalf("returned error: %v", err)
			}

			// We know we got some data back
			if tc.iterationIndex >= 0 {
				if iterations[tc.iterationIndex].Name != tc.iterationName {
					t.Fatalf("expected iteration name %s; got %s", tc.iterationName, iterations[tc.iterationIndex].Name)
				}
			} else {
				// We are testing the fact we didn't get data back
				if len(iterations) > 0 {
					t.Fatalf("expected length of iterations to be 0; got %d", len(iterations))
				}
			}
		})
	}
}
