package model

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestSegment_ToTF(t *testing.T) {
	tests := []struct {
		name     string
		input    Segment
		expected SegmentTF
	}{
		{
			name: "basic segment without auto_apply",
			input: Segment{
				ID:                "test-id",
				Name:              "test-name",
				Selector:          "{namespace=\"test\"}",
				FallbackToDefault: true,
				AutoApply:         nil,
			},
			expected: SegmentTF{
				ID:                types.StringValue("test-id"),
				Name:              types.StringValue("test-name"),
				Selector:          types.StringValue("{namespace=\"test\"}"),
				FallbackToDefault: types.BoolValue(true),
				AutoApply:         testNullAutoApplyObject(),
			},
		},
		{
			name: "segment with auto_apply and no gate",
			input: Segment{
				ID:                "test-id-4",
				Name:              "test-name-4",
				Selector:          "{namespace=\"dev\"}",
				FallbackToDefault: true,
				AutoApply:         &AutoApplyConfig{Enabled: true},
			},
			expected: SegmentTF{
				ID:                types.StringValue("test-id-4"),
				Name:              types.StringValue("test-name-4"),
				Selector:          types.StringValue("{namespace=\"dev\"}"),
				FallbackToDefault: types.BoolValue(true),
				AutoApply:         testAutoApplyObject(true, nil),
			},
		},
		{
			name: "segment with no-increase gate",
			input: Segment{
				ID:                "test-id-2",
				Name:              "test-name-2",
				Selector:          "{namespace=\"prod\"}",
				FallbackToDefault: false,
				AutoApply: &AutoApplyConfig{
					Enabled: true,
					Gate:    &GateConfig{Policy: GatePolicyNoIncrease},
				},
			},
			expected: SegmentTF{
				ID:                types.StringValue("test-id-2"),
				Name:              types.StringValue("test-name-2"),
				Selector:          types.StringValue("{namespace=\"prod\"}"),
				FallbackToDefault: types.BoolValue(false),
				AutoApply:         testAutoApplyObject(true, testStringPointer(GatePolicyNoIncrease)),
			},
		},
		{
			name: "segment with unbounded gate",
			input: Segment{
				ID:                "test-id-3",
				Name:              "test-name-3",
				Selector:          "{namespace=\"staging\"}",
				FallbackToDefault: true,
				AutoApply: &AutoApplyConfig{
					Enabled: false,
					Gate:    &GateConfig{Policy: GatePolicyUnbounded},
				},
			},
			expected: SegmentTF{
				ID:                types.StringValue("test-id-3"),
				Name:              types.StringValue("test-name-3"),
				Selector:          types.StringValue("{namespace=\"staging\"}"),
				FallbackToDefault: types.BoolValue(true),
				AutoApply:         testAutoApplyObject(false, testStringPointer(GatePolicyUnbounded)),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.ToTF()
			assert.Equal(t, tt.expected.ID, result.ID)
			assert.Equal(t, tt.expected.Name, result.Name)
			assert.Equal(t, tt.expected.Selector, result.Selector)
			assert.Equal(t, tt.expected.FallbackToDefault, result.FallbackToDefault)
			assert.Equal(t, tt.expected.AutoApply, result.AutoApply)
		})
	}
}

func TestSegmentTF_ToAPIReq(t *testing.T) {
	tests := []struct {
		name     string
		input    SegmentTF
		expected Segment
	}{
		{
			name: "basic segment without auto_apply",
			input: SegmentTF{
				ID:                types.StringValue("test-id"),
				Name:              types.StringValue("test-name"),
				Selector:          types.StringValue("{namespace=\"test\"}"),
				FallbackToDefault: types.BoolValue(true),
				AutoApply:         testNullAutoApplyObject(),
			},
			expected: Segment{
				ID:                "test-id",
				Name:              "test-name",
				Selector:          "{namespace=\"test\"}",
				FallbackToDefault: true,
				AutoApply:         nil,
			},
		},
		{
			name: "segment with auto_apply and no gate",
			input: SegmentTF{
				ID:                types.StringValue("test-id-4"),
				Name:              types.StringValue("test-name-4"),
				Selector:          types.StringValue("{namespace=\"dev\"}"),
				FallbackToDefault: types.BoolValue(true),
				AutoApply:         testAutoApplyObject(true, nil),
			},
			expected: Segment{
				ID:                "test-id-4",
				Name:              "test-name-4",
				Selector:          "{namespace=\"dev\"}",
				FallbackToDefault: true,
				AutoApply:         &AutoApplyConfig{Enabled: true},
			},
		},
		{
			name: "segment with no-increase gate",
			input: SegmentTF{
				ID:                types.StringValue("test-id-2"),
				Name:              types.StringValue("test-name-2"),
				Selector:          types.StringValue("{namespace=\"prod\"}"),
				FallbackToDefault: types.BoolValue(false),
				AutoApply:         testAutoApplyObject(true, testStringPointer(GatePolicyNoIncrease)),
			},
			expected: Segment{
				ID:                "test-id-2",
				Name:              "test-name-2",
				Selector:          "{namespace=\"prod\"}",
				FallbackToDefault: false,
				AutoApply: &AutoApplyConfig{
					Enabled: true,
					Gate:    &GateConfig{Policy: GatePolicyNoIncrease},
				},
			},
		},
		{
			name: "segment with unbounded gate",
			input: SegmentTF{
				ID:                types.StringValue("test-id-3"),
				Name:              types.StringValue("test-name-3"),
				Selector:          types.StringValue("{namespace=\"staging\"}"),
				FallbackToDefault: types.BoolValue(true),
				AutoApply:         testAutoApplyObject(false, testStringPointer(GatePolicyUnbounded)),
			},
			expected: Segment{
				ID:                "test-id-3",
				Name:              "test-name-3",
				Selector:          "{namespace=\"staging\"}",
				FallbackToDefault: true,
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
			assert.Equal(t, tt.expected.ID, result.ID)
			assert.Equal(t, tt.expected.Name, result.Name)
			assert.Equal(t, tt.expected.Selector, result.Selector)
			assert.Equal(t, tt.expected.FallbackToDefault, result.FallbackToDefault)
			assert.Equal(t, tt.expected.AutoApply, result.AutoApply)
		})
	}
}
