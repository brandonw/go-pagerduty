package pagerduty

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
)

const (
	defaultBaseURL = "https://api.pagerduty.com"
)

type service struct {
	client *Client
}

// Config represents the configuration for a PagerDuty client
type Config struct {
	BaseURL    string
	HTTPClient *http.Client
	Token      string
	UserAgent  string
}

// Client manages the communication with the PagerDuty API
type Client struct {
	baseURL            *url.URL
	client             *http.Client
	Config             *Config
	Abilities          *AbilityService
	Addons             *AddonService
	EscalationPolicies *EscalationPolicyService
	Schedules          *ScheduleService
	Services           *ServicesService
	Teams              *TeamService
	Users              *UserService
	Vendors            *VendorService
}

// Response is a wrapper around http.Response
type Response struct {
	*http.Response
}

// Pagination contains pagination information
type Pagination struct {
	Limit  int  `url:"limit,omitempty"`
	More   bool `url:"more,omitempty"`
	Offset int  `url:"offset,omitempty"`
	Total  int  `url:"total,omitempty"`
}

// NewClient returns a new PagerDuty API client.
func NewClient(config *Config) (*Client, error) {
	if config.HTTPClient == nil {
		config.HTTPClient = http.DefaultClient
	}

	if config.BaseURL == "" {
		config.BaseURL = defaultBaseURL
	}

	baseURL, err := url.Parse(config.BaseURL)
	if err != nil {
		return nil, err
	}

	c := &Client{
		baseURL: baseURL,
		client:  config.HTTPClient,
		Config:  config,
	}

	c.Abilities = &AbilityService{c}
	c.Addons = &AddonService{c}
	c.EscalationPolicies = &EscalationPolicyService{c}
	c.Schedules = &ScheduleService{c}
	c.Services = &ServicesService{c}
	c.Teams = &TeamService{c}
	c.Users = &UserService{c}
	c.Vendors = &VendorService{c}

	return c, nil
}

func (c *Client) newRequest(method, url string, body interface{}) (*http.Request, error) {
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	u := c.baseURL.String() + url

	req, err := http.NewRequest(method, u, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/vnd.pagerduty+json;version=2")
	req.Header.Add("Authorization", fmt.Sprintf("Token token=%s", c.Config.Token))
	req.Header.Add("Content-Type", "application/json")

	if c.Config.UserAgent != "" {
		req.Header.Add("User-Agent", c.Config.UserAgent)
	}

	return req, nil
}

func (c *Client) newRequestDo(method, url string, options, body, v interface{}) (*Response, error) {
	if options != nil {
		values, err := query.Values(options)
		if err != nil {
			return nil, err
		}

		url = fmt.Sprintf("%s?%s", url, values.Encode())
	}

	req, err := c.newRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	return c.do(req, v)
}

func (c *Client) do(req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response := &Response{resp}

	if err := checkResponse(response); err != nil {
		return response, err
	}

	if v != nil {
		if err := decodeJSON(response, v); err != nil {
			return response, err
		}
	}

	return response, nil
}

// ValidateAuth validates a token against the PagerDuty API
func (c *Client) ValidateAuth() error {
	_, _, err := c.Abilities.List()
	return err
}

func decodeJSON(r *Response, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func checkResponse(r *Response) error {
	if c := r.StatusCode; http.StatusOK <= c && c <= 299 {
		return nil
	}

	return decodeErrorResponse(r)
}

func decodeErrorResponse(r *Response) error {
	// Try to decode error response or fallback with standard error
	v := new(ErrorResponse)
	if err := decodeJSON(r, v); err != nil {
		return fmt.Errorf("%s API call to %s failed: %v", r.Request.Method, r.Request.URL.String(), r.Status)
	}

	return fmt.Errorf("%s API call to %s failed: %s : %v", r.Request.Method, r.Request.URL.String(), r.Status, v.ErrorResponse)
}
