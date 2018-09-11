package azuredevops_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/benmatselby/go-azuredevops/azuredevops"
)

const (
	gitRefsListURL      = "/AZURE_DEVOPS_Project/_apis/git/repositories/vscode/refs/heads"
	gitRefsListResponse = `{
		"count": 6,
		"value": [
		  {
			"name": "refs/heads/develop",
			"objectId": "67cae2b029dff7eb3dc062b49403aaedca5bad8d",
			"url": "https://fabrikam.visualstudio.com/_apis/git/repositories/278d5cd2-584d-4b63-824a-2ba458937249/refs/heads/develop"
		  },
		  {
			"name": "refs/heads/master",
			"objectId": "23d0bc5b128a10056dc68afece360d8a0fabb014",
			"url": "https://fabrikam.visualstudio.com/_apis/git/repositories/278d5cd2-584d-4b63-824a-2ba458937249/refs/heads/master"
		  },
		  {
			"name": "refs/heads/npaulk/feature",
			"objectId": "23d0bc5b128a10056dc68afece360d8a0fabb014",
			"url": "https://fabrikam.visualstudio.com/_apis/git/repositories/278d5cd2-584d-4b63-824a-2ba458937249/refs/heads/npaulk/feature"
		  },
		  {
			"name": "refs/tags/v1.0",
			"objectId": "23d0bc5b128a10056dc68afece360d8a0fabb014",
			"url": "https://fabrikam.visualstudio.com/_apis/git/repositories/278d5cd2-584d-4b63-824a-2ba458937249/refs/tags/v1.0"
		  },
		  {
			"name": "refs/tags/v1.1",
			"objectId": "23d0bc5b128a10056dc68afece360d8a0fabb014",
			"url": "https://fabrikam.visualstudio.com/_apis/git/repositories/278d5cd2-584d-4b63-824a-2ba458937249/refs/tags/v1.1"
		  },
		  {
			"name": "refs/tags/v2.0",
			"objectId": "23d0bc5b128a10056dc68afece360d8a0fabb014",
			"url": "https://fabrikam.visualstudio.com/_apis/git/repositories/278d5cd2-584d-4b63-824a-2ba458937249/refs/tags/v2.0"
		  }
		]
	  }`
)

func TestGitService_ListRefs(t *testing.T) {
	tt := []struct {
		name     string
		URL      string
		response string
		count    int
		index    int
		refName  string
		refID    string
	}{
		{name: "return 6 refs", URL: gitRefsListURL, response: gitRefsListResponse, count: 6, index: 0, refName: "refs/heads/develop", refID: "67cae2b029dff7eb3dc062b49403aaedca5bad8d"},
		{name: "can handle no refs returned", URL: gitRefsListURL, response: "{}", count: 0, index: -1},
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

			opts := azuredevops.GitRefListOptions{}
			refs, count, err := c.Git.ListRefs("vscode", "heads", &opts)
			if err != nil {
				t.Fatalf("returned error: %v", err)
			}

			if tc.index > -1 {
				if refs[tc.index].Name != tc.refName {
					t.Fatalf("expected git ref name %s, got %s", tc.refName, refs[tc.index].Name)
				}
				if refs[tc.index].ObjectID != tc.refID {
					t.Fatalf("expected git ref object id %s, got %s", tc.refID, refs[tc.index].ObjectID)
				}
			}

			if len(refs) != tc.count {
				t.Fatalf("expected length of git refs to be %d; got %d", tc.count, len(refs))
			}

			if count != tc.count {
				t.Fatalf("expected git ref count to be %d; got %d", tc.count, count)
			}
		})
	}
}
