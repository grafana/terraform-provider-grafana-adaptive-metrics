package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"

	"github.com/hashicorp/go-cleanhttp"
)

// Client is a Grafana Cloud API client.
type Client struct {
	Cfg     *Config
	BaseURL url.URL
	client  *http.Client
}

// Config contains client configuration.
type Config struct {
	// APIKey is an optional API key or service account token.
	APIKey string
	// HTTPHeaders are optional HTTP headers.
	HTTPHeaders map[string]string
	Debug       bool
	HttpClient  *http.Client
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
	}, nil
}

func (c *Client) request(method, requestPath string, body []byte, responseStruct interface{}) error {
	_, err := c.requestWithHeaders(method, requestPath, nil, body, responseStruct)
	return err
}

func (c *Client) requestWithHeaders(method, requestPath string, header http.Header, body []byte, responseStruct interface{}) (http.Header, error) {
	req, err := c.newRequest(method, requestPath, header, bytes.NewReader(body))
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
		return nil, fmt.Errorf("status: %d, body: %v", resp.StatusCode, string(bodyContents))
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

func (c *Client) newRequest(method, requestPath string, header http.Header, body io.Reader) (*http.Request, error) {
	u := c.BaseURL
	u.Path = path.Join(u.Path, requestPath)
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

type ErrNotFound struct {
	BodyContents []byte
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("status: 404, body: %s", e.BodyContents)
}
