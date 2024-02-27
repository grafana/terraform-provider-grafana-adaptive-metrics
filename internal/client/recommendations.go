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

func (c *Client) AggregationRecommendations(verbose bool) ([]model.AggregationRecommendation, error) {
	var recs []model.AggregationRecommendation
	params := url.Values{}
	if verbose {
		params.Add("verbose", "true")
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
