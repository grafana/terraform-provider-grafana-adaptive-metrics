package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-provider-grafana-adaptive-metrics/internal/model"
)

func (c *Client) ReadAggregationRuleSet(segmentID string) ([]model.AggregationRule, string, error) {
	rules := []model.AggregationRule{}

	var params url.Values
	if segmentID != "" {
		params = url.Values{
			"segment": []string{segmentID},
		}
	}

	respHeader, err := c.requestWithHeaders("GET", aggregationRulesEndpoint, params, nil, nil, &rules)
	if err != nil {
		return rules, "", err
	}

	newEtag := respHeader.Get("ETag")
	if newEtag == "" {
		return rules, "", fmt.Errorf("response from %s endpoint missing etag header", aggregationRulesEndpoint)
	}

	return rules, newEtag, nil
}

func (c *Client) UpdateAggregationRuleSet(segmentID string, rules []model.AggregationRule, etag string) (string, error) {
	// We don't want to send null to the server
	if rules == nil {
		rules = []model.AggregationRule{}
	}

	body, err := json.Marshal(rules)
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

	respHeader, err := c.requestWithHeaders("POST", aggregationRulesEndpoint, params, reqHeader, body, nil)
	if err != nil {
		return "", err
	}

	newEtag := respHeader.Get("ETag")
	if newEtag == "" {
		return "", fmt.Errorf("response from %s endpoint missing etag header", aggregationRulesEndpoint)
	}

	return newEtag, nil
}
