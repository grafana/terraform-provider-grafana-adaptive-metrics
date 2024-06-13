package client

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-provider-grafana-adaptive-metrics/internal/model"
)

const (
	segmentsEndpoint = "/aggregations/rules/segments"
)

func (c *Client) CreateSegment(s model.Segment) (model.Segment, error) {
	fmt.Println("CreateSegment", s)
	body, err := json.Marshal(s)
	if err != nil {
		return model.Segment{}, err
	}

	fmt.Println("CreateSegment", string(body))

	err = c.request("POST", segmentsEndpoint, nil, body, nil)
	if err != nil {
		return model.Segment{}, err
	}

	// TODO: modify API to return created segment
	allSegments, err := c.ListSegments()
	if err != nil {
		return model.Segment{}, err
	}

	for _, segment := range allSegments {
		if segment.Selector == s.Selector {
			return segment, nil
		}
	}

	return model.Segment{}, fmt.Errorf("segment not found after creation")
}

func (c *Client) ReadSegment(id string) (model.Segment, error) {
	allSegments, err := c.ListSegments()
	if err != nil {
		return model.Segment{}, err
	}

	for _, segment := range allSegments {
		if segment.ID == id {
			return segment, nil
		}
	}

	return model.Segment{}, fmt.Errorf("segment not found")
}

func (c *Client) UpdateSegment(s model.Segment) error {
	body, err := json.Marshal(s)
	if err != nil {
		return err
	}

	return c.request("PUT", segmentsEndpoint, nil, body, nil)
}

func (c *Client) DeleteSegment(selector string) error {
	params := url.Values{
		"segment": []string{selector},
	}
	return c.request("DELETE", segmentsEndpoint, params, nil, nil)
}

func (c *Client) ListSegments() ([]model.Segment, error) {
	resp := []model.Segment{}

	err := c.request("GET", segmentsEndpoint, nil, nil, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type SegmentResp struct {
	Result model.Segment `json:"result"`
}
