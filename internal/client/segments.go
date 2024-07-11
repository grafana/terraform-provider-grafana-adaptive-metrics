package client

import (
	"encoding/json"
	"net/url"

	"github.com/hashicorp/terraform-provider-grafana-adaptive-metrics/internal/model"
)

const (
	segmentsEndpoint = "/aggregations/rules/segments"
)

func (c *Client) CreateSegment(s model.Segment) (model.Segment, error) {
	body, err := json.Marshal(s)
	if err != nil {
		return model.Segment{}, err
	}

	c.segmentMutex.Lock()
	defer c.segmentMutex.Unlock()

	var resp model.Segment
	err = c.request("POST", segmentsEndpoint, nil, body, &resp)
	if err != nil {
		return model.Segment{}, err
	}

	return resp, nil
}

func (c *Client) ListSegments() ([]model.Segment, error) {
	c.segmentMutex.Lock()
	defer c.segmentMutex.Unlock()

	resp := []model.Segment{}
	err := c.request("GET", segmentsEndpoint, nil, nil, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) ReadSegment(id string) (model.Segment, error) {
	c.segmentMutex.Lock()
	defer c.segmentMutex.Unlock()

	resp := []model.Segment{}
	err := c.request("GET", segmentsEndpoint, nil, nil, &resp)
	if err != nil {
		return model.Segment{}, err
	}

	for _, segment := range resp {
		if segment.ID == id {
			return segment, nil
		}
	}

	return model.Segment{}, ErrNotFound{
		BodyContents: []byte("segment not found"),
	}
}

func (c *Client) UpdateSegment(s model.Segment) error {
	body, err := json.Marshal(s)
	if err != nil {
		return err
	}

	c.segmentMutex.Lock()
	defer c.segmentMutex.Unlock()

	params := url.Values{
		"segment": []string{s.ID},
	}
	return c.request("PUT", segmentsEndpoint, params, body, nil)
}

func (c *Client) DeleteSegment(id string) error {
	c.segmentMutex.Lock()
	defer c.segmentMutex.Unlock()

	params := url.Values{
		"segment": []string{id},
	}
	return c.request("DELETE", segmentsEndpoint, params, nil, nil)
}
