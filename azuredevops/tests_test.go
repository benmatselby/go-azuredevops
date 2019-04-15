package azuredevops_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/benmatselby/go-azuredevops/azuredevops"
)

const (
	testListURL      = "/AZURE_DEVOPS_Project/_apis/test/runs"
	testListResponse = `{
		"value": [
		  {
			"id": 1,
			"name": "NewTestRun2",
			"url": "https://dev.azure.com/fabrikam/fabrikam-fiber-tfvc/_apis/test/Runs/1",
			"isAutomated": false,
			"iteration": "Fabrikam-Fiber-TFVC\\Release 1\\Sprint 1",
			"owner": {
			  "id": "e5a5f7f8-6507-4c34-b397-6c4818e002f4",
			  "displayName": "Fabrikam Fiber"
			},
			"startedDate": "2014-05-04T12:50:33.17Z",
			"completedDate": "2014-05-04T12:50:31.953Z",
			"state": "Completed",
			"plan": {
			  "id": "1"
			},
			"revision": 4
		  },
		  {
			"id": 2,
			"name": "sprint1 (Manual)",
			"url": "https://dev.azure.com/fabrikam/fabrikam-fiber-tfvc/_apis/test/Runs/2",
			"isAutomated": false,
			"iteration": "Fabrikam-Fiber-TFVC\\Release 1\\Sprint 1",
			"owner": {
			  "id": "e5a5f7f8-6507-4c34-b397-6c4818e002f4",
			  "displayName": "Fabrikam Fiber"
			},
			"startedDate": "2014-05-04T12:58:36.907Z",
			"completedDate": "2014-05-04T12:58:36.47Z",
			"state": "Completed",
			"plan": {
			  "id": "1"
			},
			"revision": 3
		  }
		],
		"count": 2
		}`

	testResultsListURL      = "/AZURE_DEVOPS_Project/_apis/test/Runs/1/results"
	testResultsListResponse = `{
			"count": 1,
			"value": [
				{
					"id": 100000,
					"project": {
						"id": "5c3d39df-a0cb-49da-be01-42e53792c0e1",
						"name": "Fabrikam-Fiber-TFVC",
						"url": "https://dev.azure.com/fabrikam/_apis/projects/Fabrikam-Fiber-TFVC"
					},
					"startedDate": "2016-07-13T11:12:48.487Z",
					"completedDate": "2016-07-13T11:12:48.493Z",
					"durationInMs": 4,
					"outcome": "Passed",
					"revision": 1,
					"runBy": {
						"id": "a5cbf24d-799f-452e-82be-f049a85b5895",
						"displayName": "Fabrikam",
						"uniqueName": "fabrikamfiber.vsin@hotmail.com",
						"url": "https://dev.azure.com/fabrikam/_apis/Identities/a5cbf24d-799f-452e-82be-f049a85b5895",
						"imageUrl": "https://dev.azure.com/fabrikam/_api/_common/identityImage?id=a5cbf24d-799f-452e-82be-f049a85b5895"
					},
					"state": "Completed",
					"testCase": {
						"name": "Pass1"
					},
					"testRun": {
						"id": "16",
						"name": "VSTest Test Run release any cpu",
						"url": "https://dev.azure.com/fabrikam/Fabrikam-Fiber-TFVC/_apis/test/Runs/16"
					},
					"lastUpdatedDate": "2016-07-13T11:12:49.123Z",
					"lastUpdatedBy": {
						"id": "375baa5b-5148-4e89-a549-ec202b722d89",
						"displayName": "Project Collection Build Service (fabrikam)",
						"uniqueName": "Build\\78b5727d-4a24-4ec8-9caf-704685572174",
						"url": "https://vssps.dev.azure.com/fabrikam/_apis/Identities/375baa5b-5148-4e89-a549-ec202b722d89",
						"imageUrl": "https://dev.azure.com/fabrikam/_api/_common/identityImage?id=375baa5b-5148-4e89-a549-ec202b722d89"
					},
					"priority": 0,
					"computerName": "TASKAGENT5-0055",
					"build": {
						"id": "5",
						"name": "20160713.2",
						"url": "https://dev.azure.com/fabrikam/_apis/build/Builds/5"
					},
					"createdDate": "2016-07-13T11:12:49.123Z",
					"url": "https://dev.azure.com/fabrikam/Fabrikam-Fiber-TFVC/_apis/test/Runs/16/Results/100000",
					"failureType": "None",
					"automatedTestStorage": "unittestproject1.dll",
					"automatedTestType": "UnitTest",
					"automatedTestTypeId": "13cdc9d9-ddb5-4fa4-a97d-d965ccfc6d4b",
					"automatedTestId": "aefba017-ab06-be36-6b92-de4e29836f72",
					"area": {
						"id": "37528",
						"name": "Fabrikam-Fiber-TFVC",
						"url": "vstfs:///Classification/Node/ebe8ac79-8d9f-4a5b-8d0a-c3095c81e70e"
					},
					"testCaseTitle": "Pass1",
					"customFields": [],
					"automatedTestName": "UnitTestProject1.UnitTest1.Pass1"
				}
			]
		}`
)

func TestTestsService_List(t *testing.T) {
	tt := []struct {
		name     string
		URL      string
		response string
		count    int
		index    int
		state    string
		revision int
	}{
		{name: "return two tests", URL: testListURL, response: testListResponse, count: 2, index: 0, state: "Completed", revision: 4},
		{name: "can handle no builds returned", URL: testListURL, response: "{}", count: 0, index: -1},
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

			options := &azuredevops.TestsListOptions{}
			tests, err := c.Tests.List(options)
			if err != nil {
				t.Fatalf("returned error: %v", err)
			}

			if tc.index > -1 {
				if tests[tc.index].State != tc.state {
					t.Fatalf("expected state %s, got %s", tc.state, tests[tc.index].State)
				}

				if tests[tc.index].Revision != tc.revision {
					t.Fatalf("expected revision %d, got %d", tc.revision, tests[tc.index].Revision)
				}
			}

			if len(tests) != tc.count {
				t.Fatalf("expected length of tests to be %d; got %d", tc.count, len(tests))
			}
		})
	}
}

func TestTestsService_ResultsList(t *testing.T) {
	tt := []struct {
		name     string
		URL      string
		response string
		count    int
		index    int
		outcome  string
		testcase string
	}{
		{name: "return one result", URL: testResultsListURL, response: testResultsListResponse, count: 1, index: 0, outcome: "Passed", testcase: "Pass1"},
		{name: "can handle no results returned", URL: testResultsListURL, response: "{}", count: 0, index: -1},
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

			options := &azuredevops.TestResultsListOptions{RunID: "1"}
			tests, err := c.Tests.ResultsList(options)
			if err != nil {
				t.Fatalf("returned error: %v", err)
			}

			if tc.index > -1 {
				if tests[tc.index].Outcome != tc.outcome {
					t.Fatalf("expected outcome %s, got %s", tc.outcome, tests[tc.index].Outcome)
				}

				if tests[tc.index].TestCase.Name != tc.testcase {
					t.Fatalf("expected testcase %s got %s", tc.testcase, tests[tc.index].TestCase.Name)
				}
			}

			if len(tests) != tc.count {
				t.Fatalf("expected length of tests to be %d; got %d", tc.count, len(tests))
			}
		})
	}
}
