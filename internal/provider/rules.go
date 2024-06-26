package provider

import (
	"sync"

	"github.com/hashicorp/terraform-provider-grafana-adaptive-metrics/internal/client"
	"github.com/hashicorp/terraform-provider-grafana-adaptive-metrics/internal/model"
)

type AggregationRules struct {
	client *client.Client
	mu     sync.RWMutex

	segmentEtags map[string]string
}

func NewAggregationRules(c *client.Client) *AggregationRules {
	return &AggregationRules{client: c, mu: sync.RWMutex{}}
}

func (r *AggregationRules) Init() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	ruleSets, err := r.client.SegmentedAggregationRules()
	if err != nil {
		return err
	}

	r.segmentEtags = make(map[string]string, len(ruleSets))
	for _, ruleSet := range ruleSets {
		r.segmentEtags[ruleSet.Segment.ID] = ruleSet.Etag
	}

	return nil
}

func (r *AggregationRules) Create(segmentID string, rule model.AggregationRule) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	etag, err := r.client.CreateAggregationRule(segmentID, rule, r.segmentEtags[segmentID])
	if err != nil {
		return err
	}

	r.segmentEtags[segmentID] = etag
	return nil
}

func (r *AggregationRules) Read(segmentID string, metric string) (model.AggregationRule, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	rule, etag, err := r.client.ReadAggregationRule(segmentID, metric)
	if err != nil {
		return model.AggregationRule{}, err
	}

	r.segmentEtags[segmentID] = etag
	return rule, nil
}

func (r *AggregationRules) Update(segmentID string, rule model.AggregationRule) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	etag, err := r.client.UpdateAggregationRule(segmentID, rule, r.segmentEtags[segmentID])
	if err != nil {
		return err
	}

	r.segmentEtags[segmentID] = etag
	return nil
}

func (r *AggregationRules) Delete(segmentID string, rule model.AggregationRule) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	etag, err := r.client.DeleteAggregationRule(segmentID, rule.Metric, r.segmentEtags[segmentID])
	if err != nil {
		return err
	}

	r.segmentEtags[segmentID] = etag
	return nil
}

func (r *AggregationRules) ReadRuleSet(segmentID string) (model.AggregationRuleSet, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	rules, etag, err := r.client.ReadAggregationRuleSet(segmentID)
	if err != nil {
		return nil, err
	}

	r.segmentEtags[segmentID] = etag
	return rules, nil
}

func (r *AggregationRules) UpdateRuleSet(segmentID string, rules model.AggregationRuleSet) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	etag, err := r.client.UpdateAggregationRuleSet(segmentID, rules, r.segmentEtags[segmentID])
	if err != nil {
		return err
	}

	r.segmentEtags[segmentID] = etag
	return nil
}
