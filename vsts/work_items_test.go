package vsts_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/benmatselby/go-vsts/vsts"
)

const (
	getURL    = "/_apis/wit/workitems"
	getIdsURL = "/VSTS_TEAM/_apis/work/teamsettings/iterations/1/workitems"
	// Pulled from https://docs.microsoft.com/en-gb/rest/api/vsts/wit/work%20items/list
	getResponse = `{
		"count": 3,
		"value": [
		  {
			"id": 297,
			"rev": 1,
			"fields": {
			  "System.AreaPath": "Fabrikam-Fiber-Git",
			  "System.TeamProject": "Fabrikam-Fiber-Git",
			  "System.IterationPath": "Fabrikam-Fiber-Git",
			  "System.WorkItemType": "Product Backlog Item",
			  "System.State": "New",
			  "System.Reason": "New backlog item",
			  "System.CreatedDate": "2014-12-29T20:49:20.77Z",
			  "System.CreatedBy": "Jamal Hartnett ",
			  "System.ChangedDate": "2014-12-29T20:49:20.77Z",
			  "System.ChangedBy": "Jamal Hartnett ",
			  "System.Title": "Customer can sign in using their Microsoft Account",
			  "Microsoft.VSTS.Scheduling.Effort": 8,
			  "WEF_6CB513B6E70E43499D9FC94E5BBFB784_Kanban.Column": "New",
			  "System.Description": "Our authorization logic needs to allow for users with Microsoft accounts (formerly Live Ids) - http://msdn.microsoft.com/en-us/library/live/hh826547.aspx"
			},
			"url": "https://fabrikam.visualstudio.com/_apis/wit/workItems/297"
		  },
		  {
			"id": 299,
			"rev": 7,
			"fields": {
			  "System.AreaPath": "Fabrikam-Fiber-Git\\Website",
			  "System.TeamProject": "Fabrikam-Fiber-Git",
			  "System.IterationPath": "Fabrikam-Fiber-Git",
			  "System.WorkItemType": "Task",
			  "System.State": "To Do",
			  "System.Reason": "New task",
			  "System.AssignedTo": "Johnnie McLeod ",
			  "System.CreatedDate": "2014-12-29T20:49:21.617Z",
			  "System.CreatedBy": "Jamal Hartnett ",
			  "System.ChangedDate": "2014-12-29T20:49:28.74Z",
			  "System.ChangedBy": "Jamal Hartnett ",
			  "System.Title": "JavaScript implementation for Microsoft Account",
			  "Microsoft.VSTS.Scheduling.RemainingWork": 4,
			  "System.Description": "Follow the code samples from MSDN",
			  "System.Tags": "Tag1; Tag2"
			},
			"url": "https://fabrikam.visualstudio.com/_apis/wit/workItems/299"
		  },
		  {
			"id": 300,
			"rev": 1,
			"fields": {
			  "System.AreaPath": "Fabrikam-Fiber-Git",
			  "System.TeamProject": "Fabrikam-Fiber-Git",
			  "System.IterationPath": "Fabrikam-Fiber-Git",
			  "System.WorkItemType": "Task",
			  "System.State": "To Do",
			  "System.Reason": "New task",
			  "System.CreatedDate": "2014-12-29T20:49:22.103Z",
			  "System.CreatedBy": "Jamal Hartnett ",
			  "System.ChangedDate": "2014-12-29T20:49:22.103Z",
			  "System.ChangedBy": "Jamal Hartnett ",
			  "System.Title": "Unit Testing for MSA login",
			  "Microsoft.VSTS.Scheduling.RemainingWork": 3,
			  "System.Description": "We need to ensure we have coverage to prevent regressions"
			},
			"url": "https://fabrikam.visualstudio.com/_apis/wit/workItems/300"
		  }
		]
	  }
	`
	// Pulled from https://docs.microsoft.com/en-gb/rest/api/vsts/work/iterations/get%20iteration%20work%20items
	getIdsResponse = `{
		"workItemRelations": [
		  {
			"rel": null,
			"source": null,
			"target": {
			  "id": 1,
			  "url": "https://fabrikam.visualstudio.com/_apis/wit/workItems/1"
			}
		  },
		  {
			"rel": "System.LinkTypes.Hierarchy-Forward",
			"source": {
			  "id": 1,
			  "url": "https://fabrikam.visualstudio.com/_apis/wit/workItems/1"
			},
			"target": {
			  "id": 3,
			  "url": "https://fabrikam.visualstudio.com/_apis/wit/workItems/3"
			}
		  }
		],
		"url": "https://fabrikam.visualstudio.com/Fabrikam-Fiber/_apis/work/teamsettings/iterations/a589a806-bf11-4d4f-a031-c19813331553/workitems",
		"_links": {
		  "self": {
			"href": "https://fabrikam.visualstudio.com/Fabrikam-Fiber/_apis/work/teamsettings/iterations/a589a806-bf11-4d4f-a031-c19813331553/workitems"
		  },
		  "iteration": {
			"href": "https://fabrikam.visualstudio.com/Fabrikam-Fiber/_apis/work/teamsettings/iterations/a589a806-bf11-4d4f-a031-c19813331553"
		  }
		}
	  }
	`
)

func TestWorkItems_GetForIteration(t *testing.T) {
	tt := []struct {
		name              string
		idsBaseURL        string
		getBaseURL        string
		actualIdsURL      string
		actualGetURL      string
		idsResponse       string
		getResponse       string
		expectedWorkItems int
	}{
		{
			name:              "we get ids and we get iterations",
			idsBaseURL:        getIdsURL,
			actualIdsURL:      "/VSTS_TEAM/_apis/work/teamsettings/iterations/1/workitems?api-version=4.1-preview",
			getBaseURL:        getURL,
			actualGetURL:      "/_apis/wit/workitems?ids=1,3&fields=System.Id,System.Title,System.State,System.WorkItemType,Microsoft.VSTS.Scheduling.StoryPoints,System.BoardColumn,System.CreatedBy,System.AssignedTo&api-version=4.1-preview",
			idsResponse:       getIdsResponse,
			getResponse:       getResponse,
			expectedWorkItems: 3,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c, mux, _, teardown := setup()
			defer teardown()

			mux.HandleFunc(tc.idsBaseURL, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				testURL(t, r, tc.actualIdsURL)
				json := tc.idsResponse
				fmt.Fprint(w, json)
			})
			mux.HandleFunc(tc.getBaseURL, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				testURL(t, r, tc.actualGetURL)
				json := tc.getResponse
				fmt.Fprint(w, json)
			})

			iteration := vsts.Iteration{ID: "1"}
			workItems, err := c.WorkItems.GetForIteration("VSTS_TEAM", iteration)
			if err != nil {
				t.Fatalf("returned error: %v", err)
			}

			if len(workItems) != tc.expectedWorkItems {
				t.Fatalf("expected %d work items; got %d", tc.expectedWorkItems, len(workItems))
			}
		})
	}
}
