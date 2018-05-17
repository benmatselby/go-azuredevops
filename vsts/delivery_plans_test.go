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
