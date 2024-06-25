package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-provider-grafana-adaptive-metrics/internal/model"
)

const (
	segmentedRulesEndpoint   = "/aggregations/segmented_rules"
	aggregationRulesEndpoint = "/aggregations/rules"
	aggregationRuleEndpoint  = "/aggregations/rule/%s"
)

func (c *Client) SegmentedAggregationRules() ([]model.SegmentedRuleSet, error) {
	var rules []model.SegmentedRuleSet
	err := c.request("GET", segmentedRulesEndpoint, nil, nil, &rules)
	if err != nil {
		return rules, err
	}

	return rules, err
}

func (c *Client) CreateAggregationRule(segmentID string, rule model.AggregationRule, etag string) (string, error) {
	body, err := json.Marshal(rule)
	if err != nil {
		return "", err
	}

	reqHeader := make(http.Header)
	reqHeader.Add("If-Match", etag)

	var params url.Values
	if segmentID != "" {
		params = url.Values{
			"segment": []string{segmentID},
		}
	}

	endpoint := fmt.Sprintf(aggregationRuleEndpoint, rule.Metric)

	respHeader, err := c.requestWithHeaders("POST", endpoint, params, reqHeader, body, nil)
	if err != nil {
		return "", err
	}

	newEtag := respHeader.Get("ETag")
	if newEtag == "" {
		return "", fmt.Errorf("response from %s endpoint missing etag header", endpoint)
	}

	return newEtag, nil
}

func (c *Client) ReadAggregationRule(segmentID string, metric string) (model.AggregationRule, string, error) {
	rule := model.AggregationRule{}
	endpoint := fmt.Sprintf(aggregationRuleEndpoint, metric)

	var params url.Values
	if segmentID != "" {
		params = url.Values{
			"segment": []string{segmentID},
		}
	}

	respHeader, err := c.requestWithHeaders("GET", endpoint, params, nil, nil, &rule)
	if err != nil {
		return rule, "", err
	}

	newEtag := respHeader.Get("ETag")
	if newEtag == "" {
		return rule, "", fmt.Errorf("response from %s endpoint missing etag header", endpoint)
	}

	return rule, newEtag, nil
}

func (c *Client) UpdateAggregationRule(segmentID string, rule model.AggregationRule, etag string) (string, error) {
	body, err := json.Marshal(rule)
	if err != nil {
		return "", err
	}

	reqHeader := make(http.Header)
	reqHeader.Add("If-Match", etag)

	var params url.Values
	if segmentID != "" {
		params = url.Values{
			"segment": []string{segmentID},
		}
	}

	endpoint := fmt.Sprintf(aggregationRuleEndpoint, rule.Metric)

	respHeader, err := c.requestWithHeaders("PUT", endpoint, params, reqHeader, body, nil)
	if err != nil {
		return "", err
	}

	newEtag := respHeader.Get("ETag")
	if newEtag == "" {
		return "", fmt.Errorf("response from %s endpoint missing etag header", endpoint)
	}

	return newEtag, nil
}

func (c *Client) DeleteAggregationRule(segmentID string, metric, etag string) (string, error) {
	reqHeader := make(http.Header)
	reqHeader.Add("If-Match", etag)

	endpoint := fmt.Sprintf(aggregationRuleEndpoint, metric)

	var params url.Values
	if segmentID != "" {
		params = url.Values{
			"segment": []string{segmentID},
		}
	}

	respHeader, err := c.requestWithHeaders("DELETE", endpoint, params, reqHeader, nil, nil)
	if err != nil {
		return "", err
	}

	newEtag := respHeader.Get("ETag")
	if newEtag == "" {
		return "", fmt.Errorf("response from %s endpoint missing etag header", endpoint)
	}

	return newEtag, nil
}
