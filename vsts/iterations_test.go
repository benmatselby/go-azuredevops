package vsts

import (
	"fmt"
	"net/http"
	"testing"
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

			mux.HandleFunc("/VSTS_TEAM/_apis/work/teamsettings/iterations", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				json := `{
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
