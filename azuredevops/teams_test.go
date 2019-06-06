package azuredevops_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/benmatselby/go-azuredevops/azuredevops"
)

const (
	teamsListURL      = "/_apis/teams"
	teamsListResponse = `{
		"value": [
		  {
			"id": "564e8204-a90b-4432-883b-d4363c6125ca",
			"name": "Quality assurance",
			"url": "https://fabrikam.visualstudio.com/_apis/projects/eb6e4656-77fc-42a1-9181-4c6d8e9da5d1/teams/564e8204-a90b-4432-883b-d4363c6125ca",
			"description": "Testing staff",
			"identityUrl": "https://fabrikam.vssps.visualstudio.com/_apis/Identities/564e8204-a90b-4432-883b-d4363c6125ca"
		  },
		  {
			"id": "66df9be7-3586-467b-9c5f-425b29afedfd",
			"name": "Fabrikam-Fiber-TFVC Team",
			"url": "https://fabrikam.visualstudio.com/_apis/projects/eb6e4656-77fc-42a1-9181-4c6d8e9da5d1/teams/66df9be7-3586-467b-9c5f-425b29afedfd",
			"description": "The default project team.",
			"identityUrl": "https://fabrikam.vssps.visualstudio.com/_apis/Identities/66df9be7-3586-467b-9c5f-425b29afedfd"
		  }
		],
		"count": 2
	  }`
)

func TestTeamsService_List(t *testing.T) {
	tt := []struct {
		name     string
		URL      string
		title    string
		response string
		count    int
	}{
		{name: "can identify teams in the response", URL: teamsListURL, response: teamsListResponse, count: 2, title: "Quality assurance"},
		{name: "there are no items in the response", URL: teamsListURL, response: "{}", count: 0},
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

			opt := &azuredevops.TeamsListOptions{}
			response, count, err := c.Teams.List(opt)
			if err != nil {
				t.Fatalf("returned error: %v", err)
			}

			if count != tc.count {
				t.Fatalf("expected count in response to be %d; got %d", tc.count, count)
			}

			if len(response) != tc.count {
				t.Fatalf("expected length of response to be 0; got %d", len(response))
			}

			if count > 0 {
				if tc.title != response[0].Name {
					t.Fatalf("expected name to be %s; got %s", tc.title, response[0].Name)
				}
			}
		})
	}
}
