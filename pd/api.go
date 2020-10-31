package pd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Client struct {
	client *http.Client
	token  string
	apiURL string
}

func (c *Client) GetOncalls(offset int, since, until *time.Time) (*Oncalls, error) {
	request, _ := http.NewRequest("GET", c.apiURL+"/oncalls", nil)
	params := request.URL.Query()
	params.Add("time_zone", "UTC")
	params.Add("offset", strconv.Itoa(offset))
	if since != nil {
		params.Add("since", since.Format("2006-01-02"))
	}
	if until != nil {
		params.Add("until", until.Format("2006-01-02"))
	}
	request.URL.RawQuery = params.Encode()
	// &since=2019-05-20&until=2019-05-27

	request.Header.Set("Accept", "application/vnd.pagerduty+json;version=2")
	request.Header.Set("Authorization", fmt.Sprintf("Token token=%s", c.token))

	resp, err := c.client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	oncalls := &Oncalls{}
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, oncalls)
	if err != nil {
		return nil, err
	}

	return oncalls, nil
}

func NewClient(token string, client *http.Client, apiUrl string) *Client {
	if client == nil {
		client = http.DefaultClient
	}
	if apiUrl == "" {
		apiUrl = "https://api.pagerduty.com"
	}
	return &Client{token: token, client: client, apiURL: apiUrl}
}

type Oncalls struct {
	Items []struct {
		EscalationPolicy struct {
			ID      string `json:"id"`
			Type    string `json:"type"`
			Summary string `json:"summary"`
			Self    string `json:"self"`
			HTMLURL string `json:"html_url"`
		} `json:"escalation_policy"`
		EscalationLevel int `json:"escalation_level"`
		Schedule        struct {
			ID      string `json:"id"`
			Type    string `json:"type"`
			Summary string `json:"summary"`
			Self    string `json:"self"`
			HTMLURL string `json:"html_url"`
		} `json:"schedule"`
		User struct {
			ID      string `json:"id"`
			Type    string `json:"type"`
			Summary string `json:"summary"`
			Self    string `json:"self"`
			HTMLURL string `json:"html_url"`
		} `json:"user"`
		Start *time.Time `json:"start"`
		End   *time.Time `json:"end"`
	} `json:"oncalls"`
	Limit  int         `json:"limit"`
	Offset int         `json:"offset"`
	More   bool        `json:"more"`
	Total  interface{} `json:"total"`
}
