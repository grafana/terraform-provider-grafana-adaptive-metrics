package client

import (
	"encoding/json"
	"net/url"

	"github.com/hashicorp/terraform-provider-grafana-adaptive-metrics/internal/model"
)

const (
	recommendationsEndpoint       = "/aggregations/recommendations"
	recommendationsConfigEndpoint = "/aggregations/recommendations/config"
)

func (c *Client) AggregationRecommendations(segmentID string, verbose bool, action []string) ([]model.AggregationRecommendation, error) {
	var recs []model.AggregationRecommendation
	params := url.Values{}
	if segmentID != "" {
		params.Add("segment", segmentID)
	}
	if verbose {
		params.Add("verbose", "true")
	}
	for _, a := range action {
		params.Add("action", a)
	}
	err := c.request("GET", recommendationsEndpoint, params, nil, &recs)
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
