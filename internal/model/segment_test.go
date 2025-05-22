package model

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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
				AutoApply:         types.ObjectNull(map[string]attr.Type{"enabled": types.BoolType}),
			},
		},
		{
			name: "segment with auto_apply enabled",
			input: Segment{
				ID:                "test-id-2",
				Name:              "test-name-2",
				Selector:          "{namespace=\"prod\"}",
				FallbackToDefault: false,
				AutoApply: &AutoApplyConfig{
					Enabled: true,
				},
			},
			expected: SegmentTF{
				ID:                types.StringValue("test-id-2"),
				Name:              types.StringValue("test-name-2"),
				Selector:          types.StringValue("{namespace=\"prod\"}"),
				FallbackToDefault: types.BoolValue(false),
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
			name: "segment with auto_apply disabled",
			input: Segment{
				ID:                "test-id-3",
				Name:              "test-name-3",
				Selector:          "{namespace=\"staging\"}",
				FallbackToDefault: true,
				AutoApply: &AutoApplyConfig{
					Enabled: false,
				},
			},
			expected: SegmentTF{
				ID:                types.StringValue("test-id-3"),
				Name:              types.StringValue("test-name-3"),
				Selector:          types.StringValue("{namespace=\"staging\"}"),
				FallbackToDefault: types.BoolValue(true),
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
				AutoApply:         types.ObjectNull(map[string]attr.Type{"enabled": types.BoolType}),
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
			name: "segment with auto_apply enabled",
			input: SegmentTF{
				ID:                types.StringValue("test-id-2"),
				Name:              types.StringValue("test-name-2"),
				Selector:          types.StringValue("{namespace=\"prod\"}"),
				FallbackToDefault: types.BoolValue(false),
				AutoApply: func() types.Object {
					obj, _ := types.ObjectValue(
						map[string]attr.Type{"enabled": types.BoolType},
						map[string]attr.Value{"enabled": types.BoolValue(true)},
					)
					return obj
				}(),
			},
			expected: Segment{
				ID:                "test-id-2",
				Name:              "test-name-2",
				Selector:          "{namespace=\"prod\"}",
				FallbackToDefault: false,
				AutoApply: &AutoApplyConfig{
					Enabled: true,
				},
			},
		},
		{
			name: "segment with auto_apply disabled",
			input: SegmentTF{
				ID:                types.StringValue("test-id-3"),
				Name:              types.StringValue("test-name-3"),
				Selector:          types.StringValue("{namespace=\"staging\"}"),
				FallbackToDefault: types.BoolValue(true),
				AutoApply: func() types.Object {
					obj, _ := types.ObjectValue(
						map[string]attr.Type{"enabled": types.BoolType},
						map[string]attr.Value{"enabled": types.BoolValue(false)},
					)
					return obj
				}(),
			},
			expected: Segment{
				ID:                "test-id-3",
				Name:              "test-name-3",
				Selector:          "{namespace=\"staging\"}",
				FallbackToDefault: true,
				AutoApply: &AutoApplyConfig{
					Enabled: false,
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
			if tt.expected.AutoApply == nil {
				assert.Nil(t, result.AutoApply)
			} else {
				assert.NotNil(t, result.AutoApply)
				assert.Equal(t, tt.expected.AutoApply.Enabled, result.AutoApply.Enabled)
			}
		})
	}
}
