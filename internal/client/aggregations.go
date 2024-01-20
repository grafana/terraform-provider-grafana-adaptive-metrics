package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-provider-adaptive-metrics/internal/model"
)

const (
	recommendationsEndpoint       = "/aggregations/recommendations"
	recommendationsConfigEndpoint = "/aggregations/recommendations/config"
	aggregationRulesEndpoint      = "/aggregations/rules"
)

func (c *Client) AggregationRecommendations() ([]model.AggregationRecommendation, error) {
	var recs []model.AggregationRecommendation
	err := c.request("GET", recommendationsEndpoint, nil, nil, &recs)
	return recs, err
}

func (c *Client) AggregationRecommendationsConfig() (model.AggregationRecommendationConfiguration, error) {
	config := model.AggregationRecommendationConfiguration{}
	err := c.request("GET", recommendationsConfigEndpoint, nil, nil, &config)
	return config, err
}

func (c *Client) UpdateAggregationRecommendationsConfig(config model.AggregationRecommendationConfiguration) error {
	body, err := json.Marshal(config)
	if err != nil {
		return err
	}

	return c.request("POST", recommendationsConfigEndpoint, nil, body, nil)
}

func (c *Client) AggregationRules() ([]model.AggregationRule, string, error) {
	var rules []model.AggregationRule
	header, err := c.requestWithHeaders("GET", aggregationRulesEndpoint, nil, nil, nil, &rules)
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

	respHeader, err := c.requestWithHeaders("POST", aggregationRulesEndpoint, nil, reqHeader, body, nil)
	if err != nil {
		return "", err
	}

	newEtag := respHeader.Get("ETag")
	if newEtag == "" {
		return "", fmt.Errorf("response from %s endpoint missing etag header", aggregationRulesEndpoint)
	}

	return newEtag, nil
}
