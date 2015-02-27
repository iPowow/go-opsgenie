package opsgenie

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
)

const (
	defaultBaseURL = "https://api.opsgenie.com"
	defaultBaseURI = "/v1/json"
)

// A Client manages communication with the OpsGenie API.
type Client struct {
	// HTTP client used to communicate with the API.
	httpClient *http.Client

	// Base URL for API requests.
	baseURL *url.URL

	// Config object
	config *Config

	// Services used for talking to different parts of the OpsGenie API.
	alert *AlertService
}

type Config struct {
	apiKey string
}

func New(apiKey string) (client *Client) {
	httpClient := http.DefaultClient
	baseURL, _ := url.Parse(defaultBaseURL)

	config := Config{apiKey: apiKey}
	c := &Client{httpClient: httpClient, baseURL: baseURL, config: &config}
	c.alert = &AlertService{client: c}
	return c
}

///
/// Inplementation details
///

// Response wraps the standard http.Response
type Response struct {
	*http.Response
}

// newResponse creats a new Response for the provided http.Response.
func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	return response
}

// This will be used as error type since it implements Error() below
type ErrorResponse struct {
	Response *http.Response // HTTP response that caused this error
	Code     string         `json:"code"`  // error message
	Message  string         `json:"error"` // more detail on individual errors
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v %+v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Code, r.Message)
}

// Used to craft an http request for OpsGenie
// Entry point for any remote request
func (c *Client) newRequest(method string, resource string /* qs *url.Values */, qsOpt interface{}, body interface{}) (*http.Request, error) {

	u, _ := url.ParseRequestURI(c.baseURL.String())
	u.Path = fmt.Sprintf("%v/%v", defaultBaseURI, resource)
	// Add Query String if it's provided
	if qsOpt != nil {
		qs, err := query.Values(qsOpt)
		if err != nil {
			return nil, err
		}
		u.RawQuery = qs.Encode()
	}

	urlStr := fmt.Sprintf("%v", u)
	//fmt.Println("Req:", urlStr)

	buf, err := json.Marshal(body)
	if err != nil {
		fmt.Println("ERROR: invalid format (%s)\n", body)
		return nil, err
	}

	//fmt.Println("Data:", string(buf))

	r, err := http.NewRequest(method, urlStr, bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return r, nil
}

// Used to make and process an http request for OpsGenie
// Execution point for any remote request
func (c *Client) do(req *http.Request, result interface{}) (*Response, error) {

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response := newResponse(resp)

	err = checkResponse(resp)
	if err != nil {
		return response, err
	}

	// Inject the response body into the result struct
	if result != nil {
		if w, ok := result.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(result)
		}
	}

	return response, err
}

// Utility function used to check the return code of the response
// and parse the error.
func checkResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 210 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	//fmt.Println("response Body:", string(data))
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}
	return errorResponse
}

// addOptions adds the parameters in opt as URL query parameters to s.  opt
// must be a struct whose fields may contain "url" tags.
func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}
