package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// UsageSource represents the sources of usage data to consider when generating recommendations.
type UsageSource string

const (
	UsageSourceDashboard UsageSource = "dashboard"
	UsageSourceRules     UsageSource = "rules"
	UsageSourceQueries   UsageSource = "queries"
)

// UnusedMetricsAction defines the action to take for metrics with no recorded usage.
type UnusedMetricsAction string

const (
	UnusedMetricsActionDropAllLabels    UnusedMetricsAction = "drop_all_labels"
	UnusedMetricsActionDropCustomLabels UnusedMetricsAction = "drop_custom_labels"
	UnusedMetricsActionKeepCustomLabels UnusedMetricsAction = "keep_custom_labels"
	UnusedMetricsActionBestGuess        UnusedMetricsAction = "best_guess"
)

// Policy defines the configuration for Adaptive Metrics recommendations generation.
// Policies can be assigned to segments to customize recommendation behavior.
type Policy struct {
	ID                                   string               `json:"id"`
	Name                                 string               `json:"name"`
	UsageSources                         *[]UsageSource       `json:"usage_sources,omitempty"`
	UnusedMetricsAction                  *UnusedMetricsAction `json:"unused_metrics_action,omitempty"`
	UnusedMetricsActionActOnCustomLabels *[]string            `json:"unused_metrics_action_act_on_custom_labels,omitempty"`
	MinQueryUsages                       *int32               `json:"min_query_usages,omitempty"`
}

func (p Policy) ToTF() PolicyTF {
	policy := PolicyTF{
		ID:   types.StringValue(p.ID),
		Name: types.StringValue(p.Name),
	}

	// Handle UsageSources
	if p.UsageSources != nil {
		usageSources := make([]attr.Value, len(*p.UsageSources))
		for i, source := range *p.UsageSources {
			usageSources[i] = types.StringValue(string(source))
		}
		policy.UsageSources, _ = types.ListValue(types.StringType, usageSources)
	} else {
		policy.UsageSources = types.ListNull(types.StringType)
	}

	// Handle UnusedMetricsAction
	if p.UnusedMetricsAction != nil {
		policy.UnusedMetricsAction = types.StringValue(string(*p.UnusedMetricsAction))
	} else {
		policy.UnusedMetricsAction = types.StringNull()
	}

	// Handle UnusedMetricsActionActOnCustomLabels
	if p.UnusedMetricsActionActOnCustomLabels != nil {
		customLabels := make([]attr.Value, len(*p.UnusedMetricsActionActOnCustomLabels))
		for i, label := range *p.UnusedMetricsActionActOnCustomLabels {
			customLabels[i] = types.StringValue(label)
		}
		policy.UnusedMetricsActionActOnCustomLabels, _ = types.ListValue(types.StringType, customLabels)
	} else {
		policy.UnusedMetricsActionActOnCustomLabels = types.ListNull(types.StringType)
	}

	// Handle MinQueryUsages
	if p.MinQueryUsages != nil {
		policy.MinQueryUsages = types.Int64Value(int64(*p.MinQueryUsages))
	} else {
		policy.MinQueryUsages = types.Int64Null()
	}

	return policy
}

type PolicyTF struct {
	ID                                   types.String `tfsdk:"id"`
	Name                                 types.String `tfsdk:"name"`
	UsageSources                         types.List   `tfsdk:"usage_sources"`
	UnusedMetricsAction                  types.String `tfsdk:"unused_metrics_action"`
	UnusedMetricsActionActOnCustomLabels types.List   `tfsdk:"unused_metrics_action_act_on_custom_labels"`
	MinQueryUsages                       types.Int64  `tfsdk:"min_query_usages"`
}

func (p PolicyTF) ToAPIReq() Policy {
	policy := Policy{
		ID:   p.ID.ValueString(),
		Name: p.Name.ValueString(),
	}

	// Handle UsageSources
	if !p.UsageSources.IsNull() && !p.UsageSources.IsUnknown() {
		usageSources := make([]UsageSource, len(p.UsageSources.Elements()))
		for i, elem := range p.UsageSources.Elements() {
			if strVal, ok := elem.(types.String); ok {
				usageSources[i] = UsageSource(strVal.ValueString())
			}
		}
		policy.UsageSources = &usageSources
	}

	// Handle UnusedMetricsAction
	if !p.UnusedMetricsAction.IsNull() && !p.UnusedMetricsAction.IsUnknown() {
		action := UnusedMetricsAction(p.UnusedMetricsAction.ValueString())
		policy.UnusedMetricsAction = &action
	}

	// Handle UnusedMetricsActionActOnCustomLabels
	if !p.UnusedMetricsActionActOnCustomLabels.IsNull() && !p.UnusedMetricsActionActOnCustomLabels.IsUnknown() {
		customLabels := make([]string, len(p.UnusedMetricsActionActOnCustomLabels.Elements()))
		for i, elem := range p.UnusedMetricsActionActOnCustomLabels.Elements() {
			if strVal, ok := elem.(types.String); ok {
				customLabels[i] = strVal.ValueString()
			}
		}
		policy.UnusedMetricsActionActOnCustomLabels = &customLabels
	}

	// Handle MinQueryUsages
	if !p.MinQueryUsages.IsNull() && !p.MinQueryUsages.IsUnknown() {
		minQueryUsages := int32(p.MinQueryUsages.ValueInt64())
		policy.MinQueryUsages = &minQueryUsages
	}

	return policy
}
