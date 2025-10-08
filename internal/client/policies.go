package client

import (
	"encoding/json"

	"github.com/hashicorp/terraform-provider-grafana-adaptive-metrics/internal/model"
)

const (
	policiesEndpoint = "/aggregations/policies"
)

func (c *Client) CreatePolicy(p model.Policy) (model.Policy, error) {
	body, err := json.Marshal(p)
	if err != nil {
		return model.Policy{}, err
	}

	c.policyMutex.Lock()
	defer c.policyMutex.Unlock()

	var resp model.Policy
	err = c.request("POST", policiesEndpoint, nil, body, &resp)
	if err != nil {
		return model.Policy{}, err
	}

	return resp, nil
}

func (c *Client) ReadPolicy(id string) (model.Policy, error) {
	c.policyMutex.Lock()
	defer c.policyMutex.Unlock()

	var resp model.Policy
	err := c.request("GET", policiesEndpoint+"/"+id, nil, nil, &resp)
	if err != nil {
		return model.Policy{}, err
	}

	return resp, nil
}

func (c *Client) UpdatePolicy(p model.Policy) error {
	body, err := json.Marshal(p)
	if err != nil {
		return err
	}

	c.policyMutex.Lock()
	defer c.policyMutex.Unlock()

	return c.request("PUT", policiesEndpoint+"/"+p.ID, nil, body, nil)
}

func (c *Client) DeletePolicy(id string) error {
	c.policyMutex.Lock()
	defer c.policyMutex.Unlock()

	return c.request("DELETE", policiesEndpoint+"/"+id, nil, nil, nil)
}
