package model

import "github.com/hashicorp/terraform-plugin-framework/types"

type AggregationRecommendation struct {
	AggregationRule

	RecommendedAction  string `json:"recommended_action"`
	UsagesInRules      int64  `json:"usages_in_rules"`
	UsagesInQueries    int64  `json:"usages_in_queries"`
	UsagesInDashboards int64  `json:"usages_in_dashboards"`

	KeptLabels                   []string `json:"kept_labels,omitempty"`
	TotalSeriesAfterAggregation  int64    `json:"total_series_after_aggregation,omitempty"`
	TotalSeriesBeforeAggregation int64    `json:"total_series_before_aggregation,omitempty"`
}

func (r *AggregationRecommendation) ToTF() AggregationRecommendationTF {
	return AggregationRecommendationTF{
		Metric:    types.StringValue(r.Metric),
		MatchType: types.StringValue(r.MatchType),

		Drop:       types.BoolValue(r.Drop),
		KeepLabels: toTypesStringSlice(r.KeepLabels),
		DropLabels: toTypesStringSlice(r.DropLabels),

		Aggregations: toTypesStringSlice(r.Aggregations),

		AggregationInterval: types.StringValue(r.AggregationInterval),
		AggregationDelay:    types.StringValue(r.AggregationDelay),

		RecommendedAction:  types.StringValue(r.RecommendedAction),
		UsagesInRules:      types.Int64Value(r.UsagesInRules),
		UsagesInQueries:    types.Int64Value(r.UsagesInQueries),
		UsagesInDashboards: types.Int64Value(r.UsagesInDashboards),

		KeptLabels:                   toTypesStringSlice(r.KeptLabels),
		TotalSeriesAfterAggregation:  types.Int64Value(r.TotalSeriesAfterAggregation),
		TotalSeriesBeforeAggregation: types.Int64Value(r.TotalSeriesBeforeAggregation),
	}
}

type AggregationRecommendationListTF struct {
	Verbose         types.Bool                    `tfsdk:"verbose"`
	Action          []types.String                `tfsdk:"action"`
	Segment         types.String                  `tfsdk:"segment"`
	Recommendations []AggregationRecommendationTF `tfsdk:"recommendations"`
}

func (tf *AggregationRecommendationListTF) IsVerbose() bool {
	return tf.Verbose.ValueBool()
}

func (tf *AggregationRecommendationListTF) GetActionIn() []string {
	return toStringSlice(tf.Action)
}

type AggregationRecommendationTF struct {
	// Note: these fields are copied from RuleTF because tfsdk doesn't support struct embedding.
	Metric    types.String `tfsdk:"metric"`
	MatchType types.String `tfsdk:"match_type"`

	Drop       types.Bool     `tfsdk:"drop"`
	KeepLabels []types.String `tfsdk:"keep_labels"`
	DropLabels []types.String `tfsdk:"drop_labels"`

	Aggregations []types.String `tfsdk:"aggregations"`

	AggregationInterval types.String `tfsdk:"aggregation_interval"`
	AggregationDelay    types.String `tfsdk:"aggregation_delay"`

	RecommendedAction  types.String `tfsdk:"recommended_action"`
	UsagesInRules      types.Int64  `tfsdk:"usages_in_rules"`
	UsagesInQueries    types.Int64  `tfsdk:"usages_in_queries"`
	UsagesInDashboards types.Int64  `tfsdk:"usages_in_dashboards"`

	KeptLabels                   []types.String `tfsdk:"kept_labels"`
	TotalSeriesAfterAggregation  types.Int64    `tfsdk:"total_series_after_aggregation"`
	TotalSeriesBeforeAggregation types.Int64    `tfsdk:"total_series_before_aggregation"`
}
