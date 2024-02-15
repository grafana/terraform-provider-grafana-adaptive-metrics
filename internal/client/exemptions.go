package client

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-provider-adaptive-metrics/internal/model"
)

const (
	exemptionsEndpoint = "/v1/recommendations/exemptions"
	exemptionEndpoint  = "/v1/recommendations/exemptions/%s"
)

func (c *Client) CreateExemption(ex model.Exemption) (model.Exemption, error) {
	body, err := json.Marshal(ex)
	if err != nil {
		return model.Exemption{}, err
	}

	resp := exemptionResp{}

	err = c.request("POST", exemptionsEndpoint, body, &resp)
	if err != nil {
		return model.Exemption{}, err
	}

	return resp.Result, nil
}

func (c *Client) ReadExemption(exID string) (model.Exemption, error) {
	resp := exemptionResp{}
	endpoint := fmt.Sprintf(exemptionEndpoint, exID)

	err := c.request("GET", endpoint, nil, &resp)
	return resp.Result, err
}

func (c *Client) UpdateExemption(ex model.Exemption) error {
	body, err := json.Marshal(ex)
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf(exemptionEndpoint, ex.ID)
	return c.request("PUT", endpoint, body, nil)
}

func (c *Client) DeleteExemption(exID string) error {
	endpoint := fmt.Sprintf(exemptionEndpoint, exID)
	return c.request("DELETE", endpoint, nil, nil)
}

func (c *Client) ListExemptions() ([]model.Exemption, error) {
	resp := exemptionsResp{}

	err := c.request("GET", exemptionsEndpoint, nil, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Result, nil
}

type exemptionResp struct {
	Result model.Exemption `json:"result"`
}

type exemptionsResp struct {
	Result []model.Exemption `json:"result"`
}
