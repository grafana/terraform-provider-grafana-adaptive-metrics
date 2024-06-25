package client

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-provider-grafana-adaptive-metrics/internal/model"
)

const (
	exemptionsEndpoint = "/v1/recommendations/exemptions"
	exemptionEndpoint  = "/v1/recommendations/exemptions/%s"
)

func (c *Client) CreateExemption(segmentID string, ex model.Exemption) (model.Exemption, error) {
	body, err := json.Marshal(ex)
	if err != nil {
		return model.Exemption{}, err
	}

	resp := exemptionResp{}

	params := url.Values{
		"segment": {segmentID},
	}

	err = c.request("POST", exemptionsEndpoint, params, body, &resp)
	if err != nil {
		return model.Exemption{}, err
	}

	return resp.Result, nil
}

func (c *Client) ReadExemption(segmentID string, exID string) (model.Exemption, error) {
	resp := exemptionResp{}
	endpoint := fmt.Sprintf(exemptionEndpoint, exID)
	params := url.Values{
		"segment": {segmentID},
	}

	err := c.request("GET", endpoint, params, nil, &resp)
	return resp.Result, err
}

func (c *Client) UpdateExemption(segmentID string, ex model.Exemption) error {
	body, err := json.Marshal(ex)
	if err != nil {
		return err
	}
	params := url.Values{
		"segment": {segmentID},
	}

	endpoint := fmt.Sprintf(exemptionEndpoint, ex.ID)
	return c.request("PUT", endpoint, params, body, nil)
}

func (c *Client) DeleteExemption(segmentID string, exID string) error {
	endpoint := fmt.Sprintf(exemptionEndpoint, exID)
	params := url.Values{
		"segment": {segmentID},
	}

	return c.request("DELETE", endpoint, params, nil, nil)
}

type exemptionResp struct {
	Result model.Exemption `json:"result"`
}
