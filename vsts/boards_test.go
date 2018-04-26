package vsts_test

import (
	"fmt"
	"net/http"
	"testing"
)

const (
	boardListURL      = "/VSTS_TEAM/_apis/work/boards"
	boardListResponse = `{
		"value": [
			{
				"id": "ac1760e7-5524-4d57-9596-fa8b9e859f89",
				"name": "Stories"
			},
			{
				"id": "a4dca894-65c0-4074-bf54-08f5c2639a5d",
				"name": "Epics"
			}
		]
	}`
)

func TestBoardsService_List(t *testing.T) {
	tt := []struct {
		name           string
		URL            string
		response       string
		count          int
		index          int
		boardName      string
		result         string
		definitionName string
	}{
		{name: "return two boards", URL: boardListURL, response: boardListResponse, count: 2, index: 0, boardName: "Stories"},
		{name: "can handle no boards returned", URL: boardListURL, response: "{}", count: 0, index: -1},
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

			boards, err := c.Boards.List("VSTS_TEAM")
			if err != nil {
				t.Fatalf("returned error: %v", err)
			}

			if tc.index > -1 {
				// if boards[tc.index]. != tc.result {
				// 	t.Fatalf("expected result %s, got %s", tc.result, boards[tc.index].Result)
				// }

				if boards[tc.index].Name != tc.boardName {
					t.Fatalf("expected board name %s, got %s", tc.boardName, boards[tc.index].Name)
				}
			}

			if len(boards) != tc.count {
				t.Fatalf("expected length of builds to be %d; got %d", tc.count, len(boards))
			}
		})
	}
}
