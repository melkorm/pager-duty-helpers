package pd

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func createHTTPMockServer(response []byte, code int) (*httptest.Server, *Client) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		w.Write(response)
	}))

	c := NewClient("testowy_token", s.Client(), s.URL)

	return s, c
}

func TestIfGetOncallsReturnsOncalls(t *testing.T) {
	r := []byte(`{"oncalls":[{"escalation_policy":{"id":"PQQD111","type":"escalation_policy_reference","summary":"summary one","self":"https://api.pagerduty.com/escalation_policies/PQQD111","html_url":"https://app.pagerduty.com/escalation_policies/PQQD112"},"escalation_level":1,"schedule":{"id":"PC78222","type":"schedule_reference","summary":"summary two","self":"https://api.pagerduty.com/schedules/PC78222","html_url":"https://app.pagerduty.com/schedules/PC78223"},"user":{"id":"PVKKAQ2","type":"user_reference","summary":"Foo Bar Baz","self":"https://api.pagerduty.com/users/PVKKAA1","html_url":"https://app.pagerduty.com/users/PVKKAA1"},"start":"2019-05-27T08:30:00Z","end":"2019-06-03T08:30:00Z"}],"limit":25,"offset":2,"more":false,"total":null}`)
	fmt.Printf("%s", r)
	s, c := createHTTPMockServer(r, http.StatusOK)
	defer s.Close()

	oncalls, err := c.GetOncalls(1, nil, nil)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(oncalls.Items) != 1 {
		t.Errorf("expected items length of 1, got %d", len(oncalls.Items))
		t.FailNow()
	}

	if oncalls.Items[0].EscalationPolicy.ID != "PQQD111" {
		t.Errorf("expected escalation id `PQQD111` got %s", oncalls.Items[0].EscalationPolicy.ID)
		t.FailNow()
	}
}
