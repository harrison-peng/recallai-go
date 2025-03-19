package recallaigo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
)

const apiVersion = "v1"

type Token string

func (t Token) String() string {
	return string(t)
}

type Region string

const (
	UsEast Region = "us-east-1"
	UsWest Region = "us-west-2"
	Eu     Region = "eu-central-1"
	Japan  Region = "ap-northeast-1"
)

func (r Region) String() string {
	return string(r)
}

func (r Region) BaseURL() string {
	return fmt.Sprintf("https://%s.recall.ai", r)
}

type errJsonDecodeFunc func(data []byte) error

// ClientOption to configure API client
type ClientOption func(*Client)

type Client struct {
	httpClient *http.Client
	baseUrl    *url.URL
	apiVersion string
	Region     Region
	Token      Token

	Bot BotService
}

func NewClient(token Token, opts ...ClientOption) *Client {
	client := &Client{
		httpClient: http.DefaultClient,
		Token:      token,
		apiVersion: apiVersion,
		Region:     UsEast,
	}

	client.Bot = &BotClient{client: client}

	if err := client.setBaseURL(client.Region); err != nil {
		panic(fmt.Errorf("failed to set base URL: %w", err))
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func (c *Client) setBaseURL(region Region) error {
	apiURL := region.BaseURL()
	u, err := url.Parse(apiURL)
	if err != nil {
		return fmt.Errorf("failed to parse base URL: %w", err)
	}
	c.baseUrl = u
	return nil
}

// WithHTTPClient overrides the default http.Client
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = client
	}
}

// WithRegion sets the region for the API client
func WithRegion(region Region) ClientOption {
	return func(c *Client) {
		c.Region = region
		if err := c.setBaseURL(region); err != nil {
			panic(fmt.Errorf("failed to set base URL: %w", err))
		}
	}
}

func (c *Client) request(ctx context.Context, method, urlStr string, queryParams map[string][]string, requestBody interface{}) (*http.Response, error) {
	return c.requestImpl(ctx, method, urlStr, queryParams, requestBody, decodeClientError)
}

func (c *Client) requestImpl(ctx context.Context, method, urlStr string, queryParams map[string][]string, requestBody interface{}, errDecoder errJsonDecodeFunc) (*http.Response, error) {
	// Construct the request URL
	u, err := c.baseUrl.Parse(fmt.Sprintf("api/%s/%s", c.apiVersion, urlStr))
	if err != nil {
		return nil, fmt.Errorf("failed to parse request URL: %w", err)
	}

	// Prepare the request body
	var buf io.ReadWriter
	if requestBody != nil && !reflect.ValueOf(requestBody).IsNil() {
		body, err := json.Marshal(requestBody)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		buf = bytes.NewBuffer(body)
	}

	// Add query parameters to the URL
	if len(queryParams) > 0 {
		q := u.Query()
		for k, values := range queryParams {
			for _, value := range values {
				q.Add(k, value)
			}
		}
		u.RawQuery = q.Encode()
	}

	// Create the HTTP request
	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create new HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", c.Token))

	// Execute the request
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	// Handle non-OK responses
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read error response body: %w", err)
		}

		return nil, errDecoder(data)
	}

	return res, nil
}

func decodeClientError(data []byte) error {
	var apiErr Error
	if err := json.Unmarshal(data, &apiErr); err != nil {
		return fmt.Errorf("failed to decode error response: %w", err)
	}
	return &apiErr
}
