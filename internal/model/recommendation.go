package model

type AggregationRecommendation struct {
	AggregationRule

	RecommendedAction  string `json:"recommended_action"`
	UsagesInRules      int    `json:"usages_in_rules"`
	UsagesInQueries    int    `json:"usages_in_queries"`
	UsagesInDashboards int    `json:"usages_in_dashboards"`

	KeptLabels                   []string `json:"kept_labels,omitempty"`
	TotalSeriesAfterAggregation  int      `json:"total_series_after_aggregation,omitempty"`
	TotalSeriesBeforeAggregation int      `json:"total_series_before_aggregation,omitempty"`
}
