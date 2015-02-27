package opsgenie

// AlertService handles communication with the alert related methods of the API.
//
// API docs: http://www.opsgenie.com/docs/web-api/alert-api

import (
	"fmt"
)

const (
	resource = "alert"
)

type AlertService struct {
	client *Client
}

// AlertRequest represents a request to create an alert.
type AlertRequest struct {
	ApiKey      string `json:"apiKey"`
	Id          string `json:"id,omitempty"`
	Message     string `json:"message,omitempty"`
	Description string `json:"description,omitempty"`
	Source      string `json:"source,omitempty"`
	AlertId     string `json:"alertId,omitempty"`
	Status      string `json:"status,omitempty"`
}

type Alert struct {
	Id           string `json:"id,omitempty"`
	AlertId      string `json:"alertId,omitempty"`
	Alias        string `json:"alias,omitempty"`
	TinyId       string `json:"tinyId,omitempty"`
	Message      string `json:"message,omitempty"`
	Description  string `json:"description,omitempty"`
	Source       string `json:"source,omitempty"`
	Status       string `json:"status,omitempty"`
	Acknowledged bool   `json:"acknowledged,omitempty"`
	IsSeen       bool   `json:"isSeen,omitempty"`
	createdAt    int64  `json:"createdAt,omitempty"`
	UpdatedAt    int64  `json:"updatedAt,omitempty"`
}

type AlertList struct {
	Alerts []Alert `json:"alerts"`
}

func (s *AlertService) create(alert *AlertRequest) (*Alert, *Response, error) {
	req, err := s.client.newRequest("POST", resource, nil, alert)
	if err != nil {
		return nil, nil, err
	}

	a := new(Alert)
	resp, err := s.client.do(req, a)
	if err != nil {
		return nil, resp, err
	}

	return a, resp, err
}

func (s *AlertService) acknowledge(alert *AlertRequest) (*Alert, *Response, error) {
	u := fmt.Sprintf("%v/acknowledge", resource)
	req, err := s.client.newRequest("POST", u, nil, alert)
	if err != nil {
		return nil, nil, err
	}

	a := new(Alert)
	resp, err := s.client.do(req, a)
	if err != nil {
		return nil, resp, err
	}

	return a, resp, err
}

func (s *AlertService) get(alert *AlertRequest) (*Alert, *Response, error) {

	// Build the query string manually with url.Values
	// data := &url.Values{}
	// data.Set("apiKey", alert.ApiKey)
	// data.Add("id", alert.Id)

	// or utils a struct and `query` package
	type QsGet struct {
		ApiKey string `url:"apiKey"`
		Id     string `url:"id"`
	}
	qs := QsGet{alert.ApiKey, alert.Id}

	req, err := s.client.newRequest("GET", resource, qs, nil)
	if err != nil {
		return nil, nil, err
	}

	a := new(Alert)
	resp, err := s.client.do(req, a)
	if err != nil {
		return nil, resp, err
	}

	return a, resp, err
}

func (s *AlertService) list(alert *AlertRequest) (*AlertList, *Response, error) {

	type QsList struct {
		ApiKey string `url:"apiKey"`
		Status string `url:"status"`
		Limit  int    `url:"limit"`
	}
	qs := QsList{alert.ApiKey, alert.Status, 100}

	req, err := s.client.newRequest("GET", resource, qs, nil)
	if err != nil {
		return nil, nil, err
	}

	alertList := new(AlertList)
	resp, err := s.client.do(req, alertList)
	if err != nil {
		return nil, resp, err
	}

	return alertList, resp, err
}
