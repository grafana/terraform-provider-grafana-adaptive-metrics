package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const managedByTF = "terraform"

type SegmentedRuleSet struct {
	Etag    string            `json:"etag"`
	Segment Segment           `json:"segment"`
	Rules   []AggregationRule `json:"rules"`
}

type AggregationRule struct {
	Metric    string `json:"metric"`
	MatchType string `json:"match_type,omitempty"`

	Drop       bool     `json:"drop,omitempty"`
	KeepLabels []string `json:"keep_labels,omitempty"`
	DropLabels []string `json:"drop_labels,omitempty"`

	Aggregations []string `json:"aggregations,omitempty"`

	AggregationInterval string `json:"aggregation_interval,omitempty"`
	AggregationDelay    string `json:"aggregation_delay,omitempty"`

	ManagedBy string `json:"managed_by,omitempty"`

	Ingest bool `json:"ingest,omitempty"`
}

func (r AggregationRule) ToTF() RuleTF {
	return RuleTF{
		Metric:    types.StringValue(r.Metric),
		MatchType: types.StringValue(r.MatchType),

		Drop:       types.BoolValue(r.Drop),
		KeepLabels: toTypesStringSlice(r.KeepLabels),
		DropLabels: toTypesStringSlice(r.DropLabels),

		Aggregations: toTypesStringSlice(r.Aggregations),

		AggregationInterval: types.StringValue(r.AggregationInterval),
		AggregationDelay:    types.StringValue(r.AggregationDelay),
	}
}

func (r AggregationRule) ToRuleSetRuleTF() RuleSetRuleTF {
	return RuleSetRuleTF{
		Metric:    types.StringValue(r.Metric),
		MatchType: types.StringValue(r.MatchType),

		Drop:       types.BoolValue(r.Drop),
		KeepLabels: toTypesStringSlice(r.KeepLabels),
		DropLabels: toTypesStringSlice(r.DropLabels),

		Aggregations: toTypesStringSlice(r.Aggregations),

		AggregationInterval: types.StringValue(r.AggregationInterval),
		AggregationDelay:    types.StringValue(r.AggregationDelay),
	}
}

type RuleTF struct {
	Segment   types.String `tfsdk:"segment"`
	Metric    types.String `tfsdk:"metric"`
	MatchType types.String `tfsdk:"match_type"`

	Drop       types.Bool     `tfsdk:"drop"`
	KeepLabels []types.String `tfsdk:"keep_labels"`
	DropLabels []types.String `tfsdk:"drop_labels"`

	Aggregations []types.String `tfsdk:"aggregations"`

	AggregationInterval types.String `tfsdk:"aggregation_interval"`
	AggregationDelay    types.String `tfsdk:"aggregation_delay"`

	AutoImport types.Bool `tfsdk:"auto_import"`

	LastUpdated types.String `tfsdk:"-"`
}

func (r RuleTF) ToAPIReq() AggregationRule {
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
