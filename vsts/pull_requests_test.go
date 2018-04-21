package vsts

import (
	"fmt"
	"net/http"
	"testing"
)

const (
	pullrequestsListURL = "/_apis/git/pullrequests"
	// https://docs.microsoft.com/en-us/rest/api/vsts/git/pull%20requests/get%20pull%20requests%20by%20project
	pullrequestsResponse = `{
		"value": [
		  {
			"repository": {
			  "id": "3411ebc1-d5aa-464f-9615-0b527bc66719",
			  "name": "2016_10_31",
			  "url": "https://fabrikam.visualstudio.com/_apis/git/repositories/3411ebc1-d5aa-464f-9615-0b527bc66719",
			  "project": {
				"id": "a7573007-bbb3-4341-b726-0c4148a07853",
				"name": "2016_10_31",
				"state": "unchanged"
			  }
			},
			"pullRequestId": 22,
			"codeReviewId": 22,
			"status": "active",
			"createdBy": {
			  "id": "d6245f20-2af8-44f4-9451-8107cb2767db",
			  "displayName": "Normal Paulk",
			  "uniqueName": "fabrikamfiber16@hotmail.com",
			  "url": "https://fabrikam.visualstudio.com/_apis/Identities/d6245f20-2af8-44f4-9451-8107cb2767db",
			  "imageUrl": "https://fabrikam.visualstudio.com/_api/_common/identityImage?id=d6245f20-2af8-44f4-9451-8107cb2767db"
			},
			"creationDate": "2016-11-01T16:30:31.6655471Z",
			"title": "A new feature",
			"description": "Adding a new feature",
			"sourceRefName": "refs/heads/npaulk/my_work",
			"targetRefName": "refs/heads/new_feature",
			"mergeStatus": "succeeded",
			"mergeId": "f5fc8381-3fb2-49fe-8a0d-27dcc2d6ef82",
			"lastMergeSourceCommit": {
			  "commitId": "b60280bc6e62e2f880f1b63c1e24987664d3bda3",
			  "url": "https://fabrikam.visualstudio.com/_apis/git/repositories/3411ebc1-d5aa-464f-9615-0b527bc66719/commits/b60280bc6e62e2f880f1b63c1e24987664d3bda3"
			},
			"lastMergeTargetCommit": {
			  "commitId": "f47bbc106853afe3c1b07a81754bce5f4b8dbf62",
			  "url": "https://fabrikam.visualstudio.com/_apis/git/repositories/3411ebc1-d5aa-464f-9615-0b527bc66719/commits/f47bbc106853afe3c1b07a81754bce5f4b8dbf62"
			},
			"lastMergeCommit": {
			  "commitId": "39f52d24533cc712fc845ed9fd1b6c06b3942588",
			  "url": "https://fabrikam.visualstudio.com/_apis/git/repositories/3411ebc1-d5aa-464f-9615-0b527bc66719/commits/39f52d24533cc712fc845ed9fd1b6c06b3942588"
			},
			"reviewers": [
			  {
				"reviewerUrl": "https://fabrikam.visualstudio.com/_apis/git/repositories/3411ebc1-d5aa-464f-9615-0b527bc66719/pullRequests/22/reviewers/d6245f20-2af8-44f4-9451-8107cb2767db",
				"vote": 0,
				"id": "d6245f20-2af8-44f4-9451-8107cb2767db",
				"displayName": "Normal Paulk",
				"uniqueName": "fabrikamfiber16@hotmail.com",
				"url": "https://fabrikam.visualstudio.com/_apis/Identities/d6245f20-2af8-44f4-9451-8107cb2767db",
				"imageUrl": "https://fabrikam.visualstudio.com/_api/_common/identityImage?id=d6245f20-2af8-44f4-9451-8107cb2767db"
			  }
			],
			"url": "https://fabrikam.visualstudio.com/_apis/git/repositories/3411ebc1-d5aa-464f-9615-0b527bc66719/pullRequests/22",
			"supportsIterations": true
		  }
		],
		"count": 1
	}`
)

func TestPullRequestsService_List(t *testing.T) {
	tt := []struct {
		name     string
		URL      string
		title    string
		response string
		count    int
	}{
		{name: "happy", URL: pullrequestsListURL, response: pullrequestsResponse, count: 1, title: "A new feature"},
		{name: "there are no items in the response", URL: pullrequestsListURL, response: "{}", count: 0},
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

			opt := &PullRequestListOptions{}
			response, count, err := c.PullRequests.List(opt)
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
				if tc.title != response[0].Title {
					t.Fatalf("expected title to be %s; got %s", tc.title, response[0].Title)
				}
			}
		})
	}
}
