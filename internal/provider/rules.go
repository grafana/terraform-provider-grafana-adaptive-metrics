package provider

import (
	"fmt"
	"sync"

	"github.com/hashicorp/terraform-provider-grafana-adaptive-metrics/internal/client"
	"github.com/hashicorp/terraform-provider-grafana-adaptive-metrics/internal/model"
)

type AggregationRules struct {
	client *client.Client
	mu     sync.RWMutex

	etag  string
	rules map[string]model.AggregationRule
}

func NewAggregationRules(c *client.Client) *AggregationRules {
	return &AggregationRules{client: c, mu: sync.RWMutex{}, rules: make(map[string]model.AggregationRule)}
}

func (r *AggregationRules) Init() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	rules, etag, err := r.client.AggregationRules(r.segmentID)
	if err != nil {
		return err
	}

	for _, rule := range rules {
		r.rules[rule.Metric] = rule
	}
	r.etag = etag
	return nil
}

func (r *AggregationRules) Create(rule model.AggregationRule) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	etag, err := r.client.CreateAggregationRule(r.segmentID, rule, r.etag)
	if err != nil {
		return err
	}

	r.etag = etag
	r.rules[rule.Metric] = rule
	return nil
}

func (r *AggregationRules) Read(metric string) (model.AggregationRule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	rule, ok := r.rules[metric]
	if !ok {
		return model.AggregationRule{}, fmt.Errorf("no rule for %s found", metric)
	}

	return rule, nil
}

func (r *AggregationRules) Update(rule model.AggregationRule) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	etag, err := r.client.UpdateAggregationRule(r.segmentID, rule, r.etag)
	if err != nil {
		return err
	}

	r.etag = etag
	r.rules[rule.Metric] = rule
	return nil
}

func (r *AggregationRules) Delete(rule model.AggregationRule) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	etag, err := r.client.DeleteAggregationRule(r.segmentID, rule.Metric, r.etag)
	if err != nil {
		return err
	}

	r.etag = etag
	delete(r.rules, rule.Metric)
	return nil
}
