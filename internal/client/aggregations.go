package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-provider-grafana-adaptive-metrics/internal/model"
)

const (
	aggregationRulesEndpoint = "/aggregations/rules"
	aggregationRuleEndpoint  = "/aggregations/rule/%s"
)

func (c *Client) AggregationRules() ([]model.AggregationRule, string, error) {
	var rules []model.AggregationRule
	header, err := c.requestWithHeaders("GET", aggregationRulesEndpoint, nil, nil, &rules)
	if err != nil {
		return rules, "", err
	}

	etag := header.Get("ETag")
	if etag == "" {
		return rules, "", fmt.Errorf("response from %s endpoint missing etag header", aggregationRulesEndpoint)
	}

	return rules, etag, err
}

func (c *Client) UpdateAggregationRules(rules []model.AggregationRule, etag string) (string, error) {
	body, err := json.Marshal(rules)
	if err != nil {
		return "", err
	}

	reqHeader := make(http.Header)
	reqHeader.Add("If-Match", etag)

	respHeader, err := c.requestWithHeaders("POST", aggregationRulesEndpoint, reqHeader, body, nil)
	if err != nil {
		return "", err
	}

	newEtag := respHeader.Get("ETag")
	if newEtag == "" {
		return "", fmt.Errorf("response from %s endpoint missing etag header", aggregationRulesEndpoint)
	}

	return newEtag, nil
}

func (c *Client) CreateAggregationRule(rule model.AggregationRule, etag string) (string, error) {
	body, err := json.Marshal(rule)
	if err != nil {
		return "", err
	}

	reqHeader := make(http.Header)
	reqHeader.Add("If-Match", etag)

	endpoint := fmt.Sprintf(aggregationRuleEndpoint, rule.Metric)

	respHeader, err := c.requestWithHeaders("POST", endpoint, reqHeader, body, nil)
	if err != nil {
		return "", err
	}

	newEtag := respHeader.Get("ETag")
	if newEtag == "" {
		return "", fmt.Errorf("response from %s endpoint missing etag header", endpoint)
	}

	return newEtag, nil
}

func (c *Client) ReadAggregationRule(metric string) (model.AggregationRule, string, error) {
	rule := model.AggregationRule{}
	endpoint := fmt.Sprintf(aggregationRuleEndpoint, metric)

	respHeader, err := c.requestWithHeaders("GET", endpoint, nil, nil, &rule)
	if err != nil {
		return rule, "", err
	}

	newEtag := respHeader.Get("ETag")
	if newEtag == "" {
		return rule, "", fmt.Errorf("response from %s endpoint missing etag header", endpoint)
	}

	return rule, newEtag, nil
}

func (c *Client) UpdateAggregationRule(rule model.AggregationRule, etag string) (string, error) {
	body, err := json.Marshal(rule)
	if err != nil {
		return "", err
	}

	reqHeader := make(http.Header)
	reqHeader.Add("If-Match", etag)

	endpoint := fmt.Sprintf(aggregationRuleEndpoint, rule.Metric)

	respHeader, err := c.requestWithHeaders("PUT", endpoint, reqHeader, body, nil)
	if err != nil {
		return "", err
	}

	newEtag := respHeader.Get("ETag")
	if newEtag == "" {
		return "", fmt.Errorf("response from %s endpoint missing etag header", endpoint)
	}

	return newEtag, nil
}

func (c *Client) DeleteAggregationRule(metric, etag string) (string, error) {
	reqHeader := make(http.Header)
	reqHeader.Add("If-Match", etag)

	endpoint := fmt.Sprintf(aggregationRuleEndpoint, metric)

	respHeader, err := c.requestWithHeaders("DELETE", endpoint, reqHeader, nil, nil)
	if err != nil {
		return "", err
	}

	newEtag := respHeader.Get("ETag")
	if newEtag == "" {
		return "", fmt.Errorf("response from %s endpoint missing etag header", endpoint)
	}

	return newEtag, nil
}
