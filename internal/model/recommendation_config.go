package model

import "github.com/hashicorp/terraform-plugin-framework/types"

type AggregationRecommendationConfiguration struct {
	KeepLabels []string `json:"keep_labels,omitempty" tfsdk:"keep_labels"`
}

func (c AggregationRecommendationConfiguration) ToTF() AggregationRecommendationConfigurationTF {
	return AggregationRecommendationConfigurationTF{
		KeepLabels: toTypesStringSlice(c.KeepLabels),
	}
}

type AggregationRecommendationConfigurationTF struct {
	KeepLabels  []types.String `tfsdk:"keep_labels"`
	LastUpdated types.String   `tfsdk:"-"`
}

func (c AggregationRecommendationConfigurationTF) ToAPIReq() AggregationRecommendationConfiguration {
	return AggregationRecommendationConfiguration{
		KeepLabels: toStringSlice(c.KeepLabels),
	}
}
