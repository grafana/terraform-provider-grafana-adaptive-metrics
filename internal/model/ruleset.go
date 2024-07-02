package model

import "github.com/hashicorp/terraform-plugin-framework/types"

type AggregationRuleSet []AggregationRule

func (a AggregationRuleSet) ToTF(segment types.String) RuleSetTF {
	output := make([]RuleSetRuleTF, len(a))
	for i, rule := range a {
		output[i] = rule.ToRuleSetRuleTF()
	}

	return RuleSetTF{
		Segment: segment,
		Rules:   output,
	}
}

type RuleSetTF struct {
	Segment types.String    `tfsdk:"segment"`
	Rules   []RuleSetRuleTF `tfsdk:"rules"`
}

func (r RuleSetTF) ToAPIReq() []AggregationRule {
	output := make([]AggregationRule, len(r.Rules))
	for i, rule := range r.Rules {
		output[i] = rule.ToAPIReq()
	}
	return output
}

// RuleSetRule is a subset of RuleTF that is used in the RuleSetTF struct
// This is necessary because the tfsdk doesn't support embedding structs.
type RuleSetRuleTF struct {
	Metric    types.String `tfsdk:"metric"`
	MatchType types.String `tfsdk:"match_type"`

	Drop       types.Bool     `tfsdk:"drop"`
	KeepLabels []types.String `tfsdk:"keep_labels"`
	DropLabels []types.String `tfsdk:"drop_labels"`

	Aggregations []types.String `tfsdk:"aggregations"`

	AggregationInterval types.String `tfsdk:"aggregation_interval"`
	AggregationDelay    types.String `tfsdk:"aggregation_delay"`
}

func (r AggregationRule) IsExactMatch() bool {
	return r.MatchType == "" || r.MatchType == "exact"
}

func (r RuleSetRuleTF) ToAPIReq() AggregationRule {
	return AggregationRule{
		Metric:    r.Metric.ValueString(),
		MatchType: r.MatchType.ValueString(),

		Drop:       r.Drop.ValueBool(),
		KeepLabels: toStringSlice(r.KeepLabels),
		DropLabels: toStringSlice(r.DropLabels),

		Aggregations: toStringSlice(r.Aggregations),

		AggregationInterval: r.AggregationInterval.ValueString(),
		AggregationDelay:    r.AggregationDelay.ValueString(),

		ManagedBy: managedByTF,
	}
}

// AlignUpstreamWithState reorders the rules from upstream to match the state.
// This must mantain the semantic equality of the rules, but may reorder them as
// long as all non-exact rules are present in the same order as they were
// originally.
func AlignUpstreamWithState(state AggregationRuleSet, upstream AggregationRuleSet) AggregationRuleSet {
	// If the two rulesets are semantically different, then we can't reorder them
	// without changing the meaning of the rules.
	if !semanticallyEqualOrdering(state, upstream) {
		return upstream
	}

	output := make(AggregationRuleSet, 0, max(len(state), len(upstream)))

	// we know there can be no duplicates by metric name, so we can use a map
	upstreamMap := make(map[string]AggregationRule, len(upstream))
	for _, rule := range upstream {
		upstreamMap[rule.Metric] = rule
	}

	// we iterate over the state rules, and if we find a matching rule in the
	// upstream, we use that one and remove it from the map.
	for _, rule := range state {
		if upstreamRule, ok := upstreamMap[rule.Metric]; ok {
			output = append(output, upstreamRule)
			delete(upstreamMap, rule.Metric)
		}
	}

	// at this point, all rules with an equivalent in the state have been added
	// to the output. We can now add the remaining rules from the upstream.
	for _, rule := range upstream {
		// we check if the rule still exists in the map, as it may have been
		// removed by the previous loop.
		if _, ok := upstreamMap[rule.Metric]; ok {
			output = append(output, rule)
		}
	}

	return output
}

// semanticallyEqualOrdering checks if the ordering of the rules in two
// AggregationRuleSets is semantically equal. This means that the non-exact rules
// are in the same order in both sets.
func semanticallyEqualOrdering(a, b AggregationRuleSet) bool {
	aIndex := 0
	bIndex := 0

	for aIndex < len(a) && bIndex < len(b) {
		// Advance past exact matches
		if a[aIndex].IsExactMatch() {
			aIndex++
			continue
		}
		if b[bIndex].IsExactMatch() {
			bIndex++
			continue
		}

		// Compare the non-exact rules
		if a[aIndex].Metric != b[bIndex].Metric {
			return false
		}

		aIndex++
		bIndex++
	}

	return true
}
