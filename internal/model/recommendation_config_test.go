package model

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestRecommendationConfig_ToTF(t *testing.T) {
	tests := []struct {
		name     string
		input    AggregationRecommendationConfiguration
		expected AggregationRecommendationConfigurationTF
	}{
		{
			name: "basic recommendation config without auto_apply",
			input: AggregationRecommendationConfiguration{
				KeepLabels: []string{"namespace", "namespace2"},
				AutoApply:  nil,
			},
			expected: AggregationRecommendationConfigurationTF{
				KeepLabels: []types.String{types.StringValue("namespace"), types.StringValue("namespace2")},
				AutoApply:  types.ObjectNull(map[string]attr.Type{"enabled": types.BoolType}),
			},
		},
		{
			name: "recommendation config with auto_apply enabled",
			input: AggregationRecommendationConfiguration{
				KeepLabels: []string{"namespace", "namespace2"},
				AutoApply: &AutoApplyConfig{
					Enabled: true,
				},
			},
			expected: AggregationRecommendationConfigurationTF{
				KeepLabels: []types.String{types.StringValue("namespace"), types.StringValue("namespace2")},
				AutoApply: func() types.Object {
					obj, _ := types.ObjectValue(
						map[string]attr.Type{"enabled": types.BoolType},
						map[string]attr.Value{"enabled": types.BoolValue(true)},
					)
					return obj
				}(),
			},
		},
		{
			name: "recommendation config with auto_apply disabled",
			input: AggregationRecommendationConfiguration{
				KeepLabels: []string{"namespace", "namespace2"},
				AutoApply: &AutoApplyConfig{
					Enabled: false,
				},
			},
			expected: AggregationRecommendationConfigurationTF{
				KeepLabels: []types.String{types.StringValue("namespace"), types.StringValue("namespace2")},
				AutoApply: func() types.Object {
					obj, _ := types.ObjectValue(
						map[string]attr.Type{"enabled": types.BoolType},
						map[string]attr.Value{"enabled": types.BoolValue(false)},
					)
					return obj
				}(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.ToTF()
			assert.Equal(t, tt.expected.KeepLabels, result.KeepLabels)
			assert.Equal(t, tt.expected.AutoApply, result.AutoApply)
		})
	}
}

func TestRecommendationConfigTF_ToAPIReq(t *testing.T) {
	tests := []struct {
		name     string
		input    AggregationRecommendationConfigurationTF
		expected AggregationRecommendationConfiguration
	}{
		{
			name: "basic recommendation config without auto_apply",
			input: AggregationRecommendationConfigurationTF{
				KeepLabels: []types.String{types.StringValue("namespace"), types.StringValue("namespace2")},
				AutoApply:  types.ObjectNull(map[string]attr.Type{"enabled": types.BoolType}),
			},
			expected: AggregationRecommendationConfiguration{
				KeepLabels: []string{"namespace", "namespace2"},
				AutoApply:  nil,
			},
		},
		{
			name: "recommendation config with auto_apply enabled",
			input: AggregationRecommendationConfigurationTF{
				KeepLabels: []types.String{types.StringValue("namespace"), types.StringValue("namespace2")},
				AutoApply: func() types.Object {
					obj, _ := types.ObjectValue(
						map[string]attr.Type{"enabled": types.BoolType},
						map[string]attr.Value{"enabled": types.BoolValue(true)},
					)
					return obj
				}(),
			},
			expected: AggregationRecommendationConfiguration{
				KeepLabels: []string{"namespace", "namespace2"},
				AutoApply: &AutoApplyConfig{
					Enabled: true,
				},
			},
		},
		{
			name: "recommendation config with auto_apply disabled",
			input: AggregationRecommendationConfigurationTF{
				KeepLabels: []types.String{types.StringValue("namespace"), types.StringValue("namespace2")},
				AutoApply: func() types.Object {
					obj, _ := types.ObjectValue(
						map[string]attr.Type{"enabled": types.BoolType},
						map[string]attr.Value{"enabled": types.BoolValue(false)},
					)
					return obj
				}(),
			},
			expected: AggregationRecommendationConfiguration{
				KeepLabels: []string{"namespace", "namespace2"},
				AutoApply: &AutoApplyConfig{
					Enabled: false,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.ToAPIReq()
			assert.Equal(t, tt.expected.KeepLabels, result.KeepLabels)
			if tt.expected.AutoApply == nil {
				assert.Nil(t, result.AutoApply)
			} else {
				assert.NotNil(t, result.AutoApply)
				assert.Equal(t, tt.expected.AutoApply.Enabled, result.AutoApply.Enabled)
			}
		})
	}
}
