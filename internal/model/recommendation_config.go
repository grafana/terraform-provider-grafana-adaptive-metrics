package model

type AggregationRecommendationConfiguration struct {
	KeepLabels []string `json:"keep_labels,omitempty" tfsdk:"keep_labels"`
}
