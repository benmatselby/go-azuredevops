package vsts_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/benmatselby/go-vsts/vsts"
)

const (
	deliveryPlansListURL      = "/VSTS_Project/_apis/work/plans"
	deliveryPlansListResponse = `{
		"value": [
			{
				"id": "7154147c-43ca-44a9-9df0-2fa0a7f9d6b2",
				"name": "Plan One",
				"type": "deliveryTimelineView",
				"createdDate": "2017-12-14T16:54:06.74Z"

			},
			{
				"id": "643c57b0-ed96-45c4-b16b-77b150828eee",
				"name": "Plan Two",
				"type": "deliveryTimelineView",
				"createdDate": "2018-01-09T13:31:22.197Z"

			}
		],
		"count": 2
	}`
	deliveryPlanGetURL              = "/VSTS_Project/_apis/work/plans/7154147c-43ca-44a9-9df0-2fa0a7f9d6b2/deliverytimeline"
	deliveryPlanGetTimeLineResponse = `
	{
		"id": "7154147c-43ca-44a9-9df0-2fa0a7f9d6b2",
		"startDate": "2018-05-04T00:00:00+00:00",
		"endDate": "2018-07-06T00:00:00+00:00",
		"teams": [
			{
				"id": "c7d2dc3a-2d44-45e1-b1f1-ca2454ed368a",
				"name": "Team One",
				"iterations": [
					{
						"name": "Iteration One",
						"path": "Project\\Team\\1",
						"startDate": "2018-04-30T00:00:00Z",
						"finishDate": "2018-05-11T00:00:00Z",
						"workItems": [
							[
								1097,
								"Project\\Team\\1",
								"Feature",
								null,
								"Feaeture One",
								"New",
								"This is a tag",
								"Project",
								5,
								"Project\\Team",
								null
							]
						]
					}
				]
			}
		]
	}
	`
)

func TestDeliveryPlansService_List(t *testing.T) {
	tt := []struct {
		name     string
		URL      string
		response string
		count    int
		index    int
		planName string
		planID   string
	}{
		{name: "return two deliery plans", URL: deliveryPlansListURL, response: deliveryPlansListResponse, count: 2, index: 0, planName: "Plan One", planID: "7154147c-43ca-44a9-9df0-2fa0a7f9d6b2"},
		{name: "can handle no delivery plans returned", URL: deliveryPlansListURL, response: "{}", count: 0, index: -1},
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

			options := &vsts.DeliveryPlansListOptions{}
			plans, count, err := c.DeliveryPlans.List(options)
			if err != nil {
				t.Fatalf("returned error: %v", err)
			}

			if tc.index > -1 {
				if plans[tc.index].ID != tc.planID {
					t.Fatalf("expected delivery plan id %s, got %s", tc.planID, plans[tc.index].ID)
				}

				if plans[tc.index].Name != tc.planName {
					t.Fatalf("expected delivery plan name %s, got %s", tc.planName, plans[tc.index].Name)
				}
			}

			if len(plans) != tc.count {
				t.Fatalf("expected length of delivery plans to be %d; got %d", tc.count, len(plans))
			}

			if tc.count != count {
				t.Fatalf("expected delivery plan count to be %d; got %d", tc.count, count)
			}
		})
	}
}

func TestDeliveryPlansService_GetTimeLine(t *testing.T) {
	c, mux, _, teardown := setup()
	defer teardown()

	planID := "7154147c-43ca-44a9-9df0-2fa0a7f9d6b2"

	mux.HandleFunc(deliveryPlanGetURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		json := deliveryPlanGetTimeLineResponse
		fmt.Fprint(w, json)
	})

	timeline, err := c.DeliveryPlans.GetTimeLine(planID)
	if err != nil {
		t.Fatalf("returned error: %v", err)
	}

	if timeline.ID != "7154147c-43ca-44a9-9df0-2fa0a7f9d6b2" {
		t.Fatalf("expected delivery plan id %s, got %s", planID, timeline.ID)
	}

	if timeline.Teams[0].Name != "Team One" {
		t.Fatalf("expected delivery plan to have team[0].Name of %s, got %s", "Team One", timeline.Teams[0].Name)
	}

	if timeline.Teams[0].Iterations[0].Name != "Iteration One" {
		t.Fatalf(
			"expected delivery plan to have team[0].Iterations[0].Name of %s, got %s",
			"Iteration One",
			timeline.Teams[0].Name,
		)
	}

	if timeline.Teams[0].Iterations[0].WorkItems[0][vsts.DeliveryPlanWorkItemIDKey].(float64) != 1097 {
		t.Fatalf(
			"expected delivery plan to have teams[0].Iterations[0].WorkItems[0][vsts.DeliveryPlanWorkItemIDKey] of %v, got %v",
			1097,
			timeline.Teams[0].Iterations[0].WorkItems[0][vsts.DeliveryPlanWorkItemIDKey].(float64),
		)
	}
}
