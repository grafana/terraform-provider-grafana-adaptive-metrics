package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"sync"

	"github.com/hashicorp/go-cleanhttp"
)

// Client is a Grafana Cloud API client.
type Client struct {
	Cfg     *Config
	BaseURL url.URL
	client  *http.Client

	// This mutex is necessary for now because the API does not support
	// concurrent writes to the segments endpoint.
	segmentMutex *sync.Mutex
	// This mutex is necessary for now because the API does not support
	// concurrent writes to the policies endpoint.
	policyMutex *sync.Mutex
}

// Config contains client configuration.
type Config struct {
	// APIKey is an optional API key or service account token.
	APIKey string
	// HTTPHeaders are optional HTTP headers.
	HTTPHeaders map[string]string
	Debug       bool
	HttpClient  *http.Client

	UserAgent string
}

// New creates a new Grafana client.
func New(baseURL string, cfg *Config) (*Client, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	if cfg.HttpClient == nil {
		cfg.HttpClient = cleanhttp.DefaultClient()
	}

	return &Client{
		Cfg:     cfg,
		BaseURL: *u,
		client:  cfg.HttpClient,

		segmentMutex: &sync.Mutex{},
		policyMutex:  &sync.Mutex{},
	}, nil
}

func (c *Client) request(method, requestPath string, query url.Values, body []byte, responseStruct interface{}) error {
	_, err := c.requestWithHeaders(method, requestPath, query, nil, body, responseStruct)
	return err
}

func (c *Client) requestWithHeaders(method, requestPath string, query url.Values, header http.Header, body []byte, responseStruct interface{}) (http.Header, error) {
	log.Printf("request (%s) to %s with body data: %s", method, c.BaseURL.String(), string(body))
	req, err := c.newRequest(method, requestPath, query, header, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyContents, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if c.Cfg.Debug {
		log.Printf("response status %d with body %v", resp.StatusCode, string(bodyContents))
	}

	// check status code.
	switch {
	case resp.StatusCode == http.StatusNotFound:
		return nil, ErrNotFound{
			BodyContents: bodyContents,
		}
	case resp.StatusCode >= 400:
		// Try to parse the error response to extract meaningful error message
		errorMsg := c.extractErrorMessage(bodyContents, resp.StatusCode)
		return nil, fmt.Errorf("%s", errorMsg)
	}

	if responseStruct == nil {
		return resp.Header, nil
	}

	err = json.Unmarshal(bodyContents, responseStruct)
	if err != nil {
		return resp.Header, err
	}

	return resp.Header, nil
}

func (c *Client) newRequest(method, requestPath string, query url.Values, header http.Header, body io.Reader) (*http.Request, error) {
	u := c.BaseURL
	u.Path = path.Join(u.Path, requestPath)
	u.RawQuery = query.Encode()
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return req, err
	}

	if c.Cfg.APIKey != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Cfg.APIKey))
	}

	if c.Cfg.HTTPHeaders != nil {
		for k, v := range c.Cfg.HTTPHeaders {
			req.Header.Add(k, v)
		}
	}

	for k, vals := range header {
		for _, v := range vals {
			req.Header.Add(k, v)
		}
	}

	req.Header.Add("User-Agent", c.Cfg.UserAgent)

	if c.Cfg.Debug {
		if body == nil {
			log.Printf("request (%s) to %s with no body data", method, u.String())
		} else {
			reader, ok := body.(*bytes.Reader)
			if !ok {
				return nil, fmt.Errorf("unexpected request body type: %T", body)
			}

			if reader.Len() == 0 {
				log.Printf("request (%s) to %s with no body data", method, u.String())
			} else {
				contents := make([]byte, reader.Len())
				if _, err := reader.Read(contents); err != nil {
					return nil, fmt.Errorf("cannot read body contents for logging: %w", err)
				}
				if _, err := reader.Seek(0, io.SeekStart); err != nil {
					return nil, fmt.Errorf("failed to seek body reader to start after logging: %w", err)
				}
				log.Printf("request (%s) to %s with body data: %s", method, u.String(), string(contents))
			}
		}
	}

	req.Header.Add("Content-Type", "application/json")
	return req, err
}

func IsErrNotFound(err error) bool {
	var e ErrNotFound
	return errors.As(err, &e)
}

type ErrNotFound struct {
	BodyContents []byte
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("status: 404, body: %s", e.BodyContents)
}

func (c *Client) extractErrorMessage(bodyContents []byte, statusCode int) string {
	// Try to parse as JSON error response
	var apiError struct {
		Error string `json:"error"`
	}

	if err := json.Unmarshal(bodyContents, &apiError); err == nil && apiError.Error != "" {
		return fmt.Sprintf("API error (status %d): %s", statusCode, apiError.Error)
	}

	// Fallback to showing the raw response
	return fmt.Sprintf("status: %d, body: %s", statusCode, string(bodyContents))
}
