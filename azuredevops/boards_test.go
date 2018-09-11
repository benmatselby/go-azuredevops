package azuredevops_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

const (
	boardListURL      = "/AZURE_DEVOPS_Project/AZURE_DEVOPS_TEAM/_apis/work/boards"
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
	boardGetURL      = "/AZURE_DEVOPS_Project/AZURE_DEVOPS_TEAM/_apis/work/boards/de70b6e6-7cf3-4075-bbe8-8de651f37149"
	boardGetResponse = `{
		"id": "de70b6e6-7cf3-4075-bbe8-8de651f37149",
		"name": "Iteration x",
		"columns": [
			{
				"id": "7612dd0e-89a4-4439-8d31-d1ae6434fac9",
				"name": "Backlog"
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

			boards, err := c.Boards.List("AZURE_DEVOPS_TEAM")
			if err != nil {
				t.Fatalf("returned error: %v", err)
			}

			if tc.index > -1 {
				if boards[tc.index].Name != tc.boardName {
					t.Fatalf("expected board name: %s, got %s", tc.boardName, boards[tc.index].Name)
				}
			}

			if len(boards) != tc.count {
				t.Fatalf("expected length of builds to be %d; got %d", tc.count, len(boards))
			}
		})
	}
}

func TestBuildsService_List_ResponseDecodeFailure(t *testing.T) {
	c, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(boardListURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		json := "sad"
		fmt.Fprint(w, json)
	})

	_, err := c.Boards.List("AZURE_DEVOPS_TEAM")
	if err == nil {
		t.Fatalf("expected error decoding the response, did not get one")
	}
}

func TestBuildsService_List_CallFailureForBuildingURL(t *testing.T) {
	c, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(boardListURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		json := "{}"
		fmt.Fprint(w, json)
	})

	_, err := c.Boards.List("")
	if err != nil && !strings.Contains(err.Error(), "404") {
		t.Fatalf("expected 404 error, got %s", err.Error())
	}
}

func TestBuildsService_Get(t *testing.T) {
	tt := []struct {
		name        string
		URL         string
		response    string
		boardId     string
		boardName   string
		columnCount int
		columnId    string
		columnName  string
	}{
		{
			name:        "we get a build",
			boardId:     "de70b6e6-7cf3-4075-bbe8-8de651f37149",
			boardName:   "Iteration x",
			URL:         boardGetURL,
			response:    boardGetResponse,
			columnCount: 1,
			columnId:    "7612dd0e-89a4-4439-8d31-d1ae6434fac9",
			columnName:  "Backlog",
		},
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

			board, err := c.Boards.Get("AZURE_DEVOPS_TEAM", tc.boardId)
			if err != nil {
				t.Fatalf("returned error: %v", err)
			}

			if board.Name != tc.boardName {
				t.Fatalf("expected board name: %s, got %s", tc.boardName, board.Name)
			}

			if len(board.Columns) != tc.columnCount {
				t.Fatalf("expected board column count: %d, got %d", tc.columnCount, len(board.Columns))
			}

			if board.Columns[0].ID != tc.columnId {
				t.Fatalf("expected column id: %s, got %s", tc.columnId, board.Columns[0].ID)
			}

			if board.Columns[0].Name != tc.columnName {
				t.Fatalf("expected column name: %s, got %s", tc.columnName, board.Columns[0].Name)
			}
		})
	}
}

func TestBuildsService_Get_ResponseDecodeFailure(t *testing.T) {
	c, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(boardGetURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		json := "sad"
		fmt.Fprint(w, json)
	})

	_, err := c.Boards.Get("AZURE_DEVOPS_TEAM", "b5f5e386-fd86-4459-af9a-72f881bd1b23")
	if err == nil {
		t.Fatalf("expected error decoding the response, did not get one")
	}
}

func TestBuildsService_Get_CallFailureForBuildingURL(t *testing.T) {
	c, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(boardGetURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		json := "{}"
		fmt.Fprint(w, json)
	})

	_, err := c.Boards.Get("AZURE_DEVOPS_TEAM", "b5f5e386-fd86-4459-af9a-72f881bd1b23")
	if err != nil && !strings.Contains(err.Error(), "404") {
		t.Fatalf("expected 404 error, got %s", err.Error())
	}
}
