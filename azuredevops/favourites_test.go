package azuredevops_test

import (
	"fmt"
	"net/http"
	"testing"
)

const (
	favouritesListURL      = "/_apis/Favorite/Favorites"
	favouritesListResponse = `{
		"value": [
			{
				"id": "3597be53-16a3-4ba9-b383-c1d869503d6e",
				"artifactName": "build-death-star",
				"artifactType": "Microsoft.TeamFoundation.Build.Definition",
				"artifactId": "09747a31-8cfc-4c85-b41c-8ed6e03beb08"
			},
			{
				"id": "7d102f74-b5b6-47ef-82a0-e99e53e15bba",
				"artifactName": "death-star",
				"artifactType": "Microsoft.TeamFoundation.Git.Repository",
				"artifactId": "91382499-f7ab-445d-aee0-1202052ec059"
			}
		],
		"count": 2
	}`
)

func TestFavouritesService_List(t *testing.T) {
	tt := []struct {
		name         string
		URL          string
		response     string
		count        int
		index        int
		artifactName string
		artifactID   string
	}{
		{name: "return two favourites", URL: favouritesListURL, response: favouritesListResponse, count: 2, index: 0, artifactName: "build-death-star", artifactID: "09747a31-8cfc-4c85-b41c-8ed6e03beb08"},
		{name: "can handle no favourites returned", URL: favouritesListURL, response: "{}", count: 0, index: -1},
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

			favourites, count, err := c.Favourites.List()
			if err != nil {
				t.Fatalf("returned error: %v", err)
			}

			if tc.index > -1 {
				if favourites[tc.index].ArtifactName != tc.artifactName {
					t.Fatalf("expected favourite artifact name %s, got %s", tc.artifactName, favourites[tc.index].ArtifactName)
				}
				if favourites[tc.index].ArtifactID != tc.artifactID {
					t.Fatalf("expected favourite artifact id %s, got %s", tc.artifactID, favourites[tc.index].ArtifactID)
				}
			}

			if len(favourites) != tc.count {
				t.Fatalf("expected length of artifacts to be %d; got %d", tc.count, len(favourites))
			}

			if count != tc.count {
				t.Fatalf("expected artifact count to be %d; got %d", tc.count, count)
			}
		})
	}
}
