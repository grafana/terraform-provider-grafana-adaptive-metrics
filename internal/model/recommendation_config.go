package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AggregationRecommendationConfiguration struct {
	KeepLabels []string         `json:"keep_labels,omitempty" tfsdk:"keep_labels"`
	AutoApply  *AutoApplyConfig `json:"auto_apply,omitempty" tfsdk:"auto_apply"`
}

const (
	GatePolicyUnbounded  = "unbounded"
	GatePolicyNoIncrease = "no-increase"
)

type AutoApplyConfig struct {
	Enabled bool        `json:"enabled" tfsdk:"enabled"`
	Gate    *GateConfig `json:"gate,omitempty" tfsdk:"gate"`
}

type GateConfig struct {
	Policy string `json:"policy,omitempty" tfsdk:"policy"`
}

func gateAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{"policy": types.StringType}
}

func autoApplyAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"enabled": types.BoolType,
		"gate":    types.ObjectType{AttrTypes: gateAttributeTypes()},
	}
}

func autoApplyConfigToTF(config *AutoApplyConfig) types.Object {
	if config == nil {
		return types.ObjectNull(autoApplyAttributeTypes())
	}

	gate := types.ObjectNull(gateAttributeTypes())
	if config.Gate != nil {
		gate = types.ObjectValueMust(
			gateAttributeTypes(),
			map[string]attr.Value{"policy": types.StringValue(config.Gate.Policy)},
		)
	}

	return types.ObjectValueMust(
		autoApplyAttributeTypes(),
		map[string]attr.Value{
			"enabled": types.BoolValue(config.Enabled),
			"gate":    gate,
		},
	)
}

func autoApplyConfigFromTF(value types.Object) *AutoApplyConfig {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	enabled, ok := value.Attributes()["enabled"].(types.Bool)
	if !ok {
		return nil
	}

	config := &AutoApplyConfig{Enabled: enabled.ValueBool()}
	gate, ok := value.Attributes()["gate"].(types.Object)
	if !ok || gate.IsNull() || gate.IsUnknown() {
		return config
	}

	policy, ok := gate.Attributes()["policy"].(types.String)
	if ok && !policy.IsNull() && !policy.IsUnknown() {
		config.Gate = &GateConfig{Policy: policy.ValueString()}
	}

	return config
}

func (c AggregationRecommendationConfiguration) ToTF() AggregationRecommendationConfigurationTF {
	cfg := AggregationRecommendationConfigurationTF{
		KeepLabels: toTypesStringSlice(c.KeepLabels),
	}

	cfg.AutoApply = autoApplyConfigToTF(c.AutoApply)

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

	cfg.AutoApply = autoApplyConfigFromTF(c.AutoApply)

	return cfg
}
