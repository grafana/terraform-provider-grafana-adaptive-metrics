package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AggregationRecommendationConfiguration struct {
	KeepLabels []string         `json:"keep_labels,omitempty" tfsdk:"keep_labels"`
	AutoApply  *AutoApplyConfig `json:"auto_apply,omitempty" tfsdk:"auto_apply"`
}

type AutoApplyConfig struct {
	Enabled bool `json:"enabled" tfsdk:"enabled"`
}

func (c AggregationRecommendationConfiguration) ToTF() AggregationRecommendationConfigurationTF {
	cfg := AggregationRecommendationConfigurationTF{
		KeepLabels: toTypesStringSlice(c.KeepLabels),
	}

	if c.AutoApply != nil {
		cfg.AutoApply, _ = types.ObjectValue(map[string]attr.Type{"enabled": types.BoolType}, map[string]attr.Value{"enabled": types.BoolValue(c.AutoApply.Enabled)})
	} else {
		cfg.AutoApply = types.ObjectNull(map[string]attr.Type{"enabled": types.BoolType})
	}

	return cfg
}

type AggregationRecommendationConfigurationTF struct {
	KeepLabels  []types.String `tfsdk:"keep_labels"`
	AutoApply   types.Object   `tfsdk:"auto_apply"`
	LastUpdated types.String   `tfsdk:"-"`
}

func (c AggregationRecommendationConfigurationTF) ToAPIReq() AggregationRecommendationConfiguration {
	cfg := AggregationRecommendationConfiguration{
		KeepLabels: toStringSlice(c.KeepLabels),
	}

	if !c.AutoApply.IsNull() {
		attrs := c.AutoApply.Attributes()
		if enabled, ok := attrs["enabled"]; ok {
			if boolVal, ok := enabled.(types.Bool); ok {
				cfg.AutoApply = &AutoApplyConfig{
					Enabled: boolVal.ValueBool(),
				}
			}
		}
	}

	return cfg
}
