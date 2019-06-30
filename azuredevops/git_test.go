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
	gitListURL      = "/AZURE_DEVOPS_Project/_apis/git/repositories"
	gitListResponse = `{
			"value": [{
					"id": "9ae84a75-c776-4c08-83b8-f9b19c520c99",
					"name": "Repo1",
					"url": "https://fabrikam.visualstudio.com/7a5df255-147e-4e18-b9d7-68fca6102808/_apis/git/repositories/9ae84a75-c776-4c08-83b8-f9b19c520c99",
					"project": {
							"id": "7a5df255-147e-4e18-b9d7-68fca6102808",
							"name": "AZURE_DEVOPS_Project",
							"description": "Project description.",
							"url": "https://fabrikam.visualstudio.com/_apis/projects/7a5df255-147e-4e18-b9d7-68fca6102808",
							"state": "wellFormed",
							"revision": 412,
							"visibility": "private",
							"lastUpdateTime": "2019-06-10T14:49:21.857Z"
					},
					"defaultBranch": "refs/heads/master",
					"size": 114417,
					"remoteUrl": "https://fabrikam.visualstudio.com/AZURE_DEVOPS_Project/_git/Repo1",
					"sshUrl": "fabrikam@vs-ssh.visualstudio.com:v3/fabrikam/AZURE_DEVOPS_Project/Repo1",
					"webUrl": "https://fabrikam.visualstudio.com/AZURE_DEVOPS_Project/_git/Repo1"
			},
			{
				"id": "9ae84a75-c776-4c08-83b8-f9b19c520c99",
				"name": "Repo 2",
				"url": "https://fabrikam.visualstudio.com/7a5df255-147e-4e18-b9d7-68fca6102808/_apis/git/repositories/9ae84a75-c776-4c08-83b8-f9b19c520c99",
				"project": {
						"id": "7a5df255-147e-4e18-b9d7-68fca6102808",
						"name": "AZURE_DEVOPS_Project",
						"description": "Project description.",
						"url": "https://fabrikam.visualstudio.com/_apis/projects/7a5df255-147e-4e18-b9d7-68fca6102808",
						"state": "wellFormed",
						"revision": 412,
						"visibility": "private",
						"lastUpdateTime": "2019-06-10T14:49:21.857Z"
				},
				"defaultBranch": "refs/heads/master",
				"size": 114417,
				"remoteUrl": "https://fabrikam.visualstudio.com/AZURE_DEVOPS_Project/_git/Repo2",
				"sshUrl": "fabrikam@vs-ssh.visualstudio.com:v3/fabrikam/AZURE_DEVOPS_Project/Repo2",
				"webUrl": "https://fabrikam.visualstudio.com/AZURE_DEVOPS_Project/_git/Repo2"
			}],
			"count": 2
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

func TestGitService_List(t *testing.T) {
	tt := []struct {
		name           string
		URL            string
		response       string
		count          int
		index          int
		repositoryName string
	}{
		{name: "return two repositories", URL: gitListURL, response: gitListResponse, count: 2, index: 0, repositoryName: "Repo1"},
		{name: "can handle no repositories returned", URL: gitListURL, response: "{}", count: 0, index: -1},
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

			repositories, err := c.Git.List()
			if err != nil {
				t.Fatalf("returned error: %v", err)
			}

			if tc.index > -1 {
				if repositories[tc.index].Name != tc.repositoryName {
					t.Fatalf("expected result %s, got %s", tc.repositoryName, repositories[tc.index].Name)
				}
			}

			if len(repositories) != tc.count {
				t.Fatalf("expected length of repositories to be %d; got %d", tc.count, len(repositories))
			}
		})
	}
}
