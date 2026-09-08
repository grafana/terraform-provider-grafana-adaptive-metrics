package model

import (
	"testing"

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
				AutoApply:  testNullAutoApplyObject(),
			},
		},
		{
			name: "recommendation config with auto_apply and no gate",
			input: AggregationRecommendationConfiguration{
				KeepLabels: []string{"namespace", "namespace2"},
				AutoApply:  &AutoApplyConfig{Enabled: true},
			},
			expected: AggregationRecommendationConfigurationTF{
				KeepLabels: []types.String{types.StringValue("namespace"), types.StringValue("namespace2")},
				AutoApply:  testAutoApplyObject(true, nil),
			},
		},
		{
			name: "recommendation config with no-increase gate",
			input: AggregationRecommendationConfiguration{
				KeepLabels: []string{"namespace", "namespace2"},
				AutoApply: &AutoApplyConfig{
					Enabled: true,
					Gate:    &GateConfig{Policy: GatePolicyNoIncrease},
				},
			},
			expected: AggregationRecommendationConfigurationTF{
				KeepLabels: []types.String{types.StringValue("namespace"), types.StringValue("namespace2")},
				AutoApply:  testAutoApplyObject(true, testStringPointer(GatePolicyNoIncrease)),
			},
		},
		{
			name: "recommendation config with unbounded gate",
			input: AggregationRecommendationConfiguration{
				KeepLabels: []string{"namespace", "namespace2"},
				AutoApply: &AutoApplyConfig{
					Enabled: false,
					Gate:    &GateConfig{Policy: GatePolicyUnbounded},
				},
			},
			expected: AggregationRecommendationConfigurationTF{
				KeepLabels: []types.String{types.StringValue("namespace"), types.StringValue("namespace2")},
				AutoApply:  testAutoApplyObject(false, testStringPointer(GatePolicyUnbounded)),
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
				AutoApply:  testNullAutoApplyObject(),
			},
			expected: AggregationRecommendationConfiguration{
				KeepLabels: []string{"namespace", "namespace2"},
				AutoApply:  nil,
			},
		},
		{
			name: "recommendation config with auto_apply and no gate",
			input: AggregationRecommendationConfigurationTF{
				KeepLabels: []types.String{types.StringValue("namespace"), types.StringValue("namespace2")},
				AutoApply:  testAutoApplyObject(true, nil),
			},
			expected: AggregationRecommendationConfiguration{
				KeepLabels: []string{"namespace", "namespace2"},
				AutoApply:  &AutoApplyConfig{Enabled: true},
			},
		},
		{
			name: "recommendation config with no-increase gate",
			input: AggregationRecommendationConfigurationTF{
				KeepLabels: []types.String{types.StringValue("namespace"), types.StringValue("namespace2")},
				AutoApply:  testAutoApplyObject(true, testStringPointer(GatePolicyNoIncrease)),
			},
			expected: AggregationRecommendationConfiguration{
				KeepLabels: []string{"namespace", "namespace2"},
				AutoApply: &AutoApplyConfig{
					Enabled: true,
					Gate:    &GateConfig{Policy: GatePolicyNoIncrease},
				},
			},
		},
		{
			name: "recommendation config with unbounded gate",
			input: AggregationRecommendationConfigurationTF{
				KeepLabels: []types.String{types.StringValue("namespace"), types.StringValue("namespace2")},
				AutoApply:  testAutoApplyObject(false, testStringPointer(GatePolicyUnbounded)),
			},
			expected: AggregationRecommendationConfiguration{
				KeepLabels: []string{"namespace", "namespace2"},
				AutoApply: &AutoApplyConfig{
					Enabled: false,
					Gate:    &GateConfig{Policy: GatePolicyUnbounded},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.ToAPIReq()
			assert.Equal(t, tt.expected.KeepLabels, result.KeepLabels)
			assert.Equal(t, tt.expected.AutoApply, result.AutoApply)
		})
	}
}
