package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"time"

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
	// NumRetries contains the number of attempted retries
	NumRetries int
	// RetryTimeout says how long to wait before retrying a request
	RetryTimeout time.Duration
	// RetryStatusCodes contains the list of status codes to retry, use "x" as a wildcard for a single digit (default: [429, 5xx])
	RetryStatusCodes []string
}

// New creates a new Grafana client.
func New(baseURL string, cfg *Config) (*Client, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	return &Client{
		Cfg:     cfg,
		BaseURL: *u,
		client:  cleanhttp.DefaultClient(),
	}, nil
}

func (c *Client) request(method, requestPath string, query url.Values, body []byte, responseStruct interface{}) error {
	_, err := c.requestWithHeaders(method, requestPath, query, nil, body, responseStruct)
	return err
}

func (c *Client) requestWithHeaders(method, requestPath string, query url.Values, header http.Header, body []byte, responseStruct interface{}) (http.Header, error) {
	var (
		req          *http.Request
		resp         *http.Response
		err          error
		bodyContents []byte
	)

	// retry logic
	for n := 0; n <= c.Cfg.NumRetries; n++ {
		req, err = c.newRequest(method, requestPath, query, header, bytes.NewReader(body))
		if err != nil {
			return nil, err
		}

		// Wait a bit if that's not the first request
		if n != 0 {
			if c.Cfg.RetryTimeout == 0 {
				c.Cfg.RetryTimeout = time.Second * 5
			}
			time.Sleep(c.Cfg.RetryTimeout)
		}

		resp, err = c.client.Do(req)

		// If err is not nil, retry again
		// That's either caused by client policy, or failure to speak HTTP (such as network connectivity problem). A
		// non-2xx status code doesn't cause an error.
		if err != nil {
			continue
		}

		// read the body (even on non-successful HTTP status codes), as that's what the unit tests expect
		bodyContents, err = io.ReadAll(resp.Body)
		_ = resp.Body.Close()

		// if there was an error reading the body, try again
		if err != nil {
			continue
		}

		shouldRetry, err := matchRetryCode(resp.StatusCode, c.Cfg.RetryStatusCodes)
		if err != nil {
			return nil, err
		}
		if !shouldRetry {
			break
		}
	}
	if err != nil {
		return nil, err
	}

	if os.Getenv("GF_LOG") != "" {
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

	if os.Getenv("GF_LOG") != "" {
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

// matchRetryCode checks if the status code matches any of the configured retry status codes.
func matchRetryCode(gottenCode int, retryCodes []string) (bool, error) {
	gottenCodeStr := strconv.Itoa(gottenCode)
	for _, retryCode := range retryCodes {
		if len(retryCode) != 3 {
			return false, fmt.Errorf("invalid retry status code: %s", retryCode)
		}
		matched := true
		for i := range retryCode {
			c := retryCode[i]
			if c == 'x' {
				continue
			}
			if gottenCodeStr[i] != c {
				matched = false
				break
			}
		}
		if matched {
			return true, nil
		}
	}

	return false, nil
}

type ErrNotFound struct {
	BodyContents []byte
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("status: 404, body: %s", e.BodyContents)
}
