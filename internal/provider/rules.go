package provider

import (
	"fmt"
	"slices"
	"sync"

	"github.com/hashicorp/terraform-provider-adaptive-metrics/internal/client"
	"github.com/hashicorp/terraform-provider-adaptive-metrics/internal/model"
)

type AggregationRules struct {
	client *client.Client
	mu     sync.RWMutex // TODO: Does this need to be thread-safe?

	// Preserve the original order of rules so we can replay them in the same order when we submit updates.
	orderedRules []string
	etag         string
	rules        map[string]model.AggregationRule
}

func NewAggregationRules(c *client.Client) *AggregationRules {
	return &AggregationRules{client: c, mu: sync.RWMutex{}, rules: make(map[string]model.AggregationRule)}
}

func (r *AggregationRules) Init() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	rules, etag, err := r.client.AggregationRules()
	if err != nil {
		return err
	}

	for _, rule := range rules {
		r.rules[rule.Metric] = rule
		r.orderedRules = append(r.orderedRules, rule.Metric)
	}
	r.etag = etag
	return nil
}

func (r *AggregationRules) Create(rule model.AggregationRule) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.rules[rule.Metric]; ok {
		return fmt.Errorf("rule for %s already exists", rule.Metric)
	}
	r.rules[rule.Metric] = rule
	r.orderedRules = append(r.orderedRules, rule.Metric)

	return r.syncRules()
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

	if _, ok := r.rules[rule.Metric]; !ok {
		return fmt.Errorf("no rule for %s found", rule.Metric)
	}
	r.rules[rule.Metric] = rule

	return r.syncRules()
}

func (r *AggregationRules) Delete(rule model.AggregationRule) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.rules[rule.Metric]; !ok {
		return fmt.Errorf("no rule for %s found", rule.Metric)
	}

	for i, metric := range r.orderedRules {
		if metric == rule.Metric {
			r.orderedRules = slices.Delete(r.orderedRules, i, i+1)
		}
	}
	delete(r.rules, rule.Metric)

	return r.syncRules()
}

func (r *AggregationRules) syncRules() error {
	payload := make([]model.AggregationRule, len(r.rules))
	for i, metric := range r.orderedRules {
		payload[i] = r.rules[metric]
	}

	etag, err := r.client.UpdateAggregationRules(payload, r.etag)
	if err != nil {
		return fmt.Errorf("could not update aggregation rules: %w", err)
	}

	r.etag = etag
	return nil
}
