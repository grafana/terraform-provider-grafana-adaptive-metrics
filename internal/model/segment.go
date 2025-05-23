package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Segment struct {
	ID                string           `json:"id"`
	Name              string           `json:"name"`
	Selector          string           `json:"selector"`
	FallbackToDefault bool             `json:"fallback_to_default"`
	AutoApply         *AutoApplyConfig `json:"auto_apply,omitempty"`
}

type AutoApplyConfig struct {
	Enabled bool `json:"enabled" tfsdk:"enabled"`
}

func (e Segment) ToTF() SegmentTF {
	segment := SegmentTF{
		ID:                types.StringValue(e.ID),
		Name:              types.StringValue(e.Name),
		Selector:          types.StringValue(e.Selector),
		FallbackToDefault: types.BoolValue(e.FallbackToDefault),
	}

	if e.AutoApply != nil {
		segment.AutoApply, _ = types.ObjectValue(map[string]attr.Type{"enabled": types.BoolType}, map[string]attr.Value{"enabled": types.BoolValue(e.AutoApply.Enabled)})
	} else {
		segment.AutoApply = types.ObjectNull(map[string]attr.Type{"enabled": types.BoolType})
	}

	return segment
}

type SegmentTF struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	Selector          types.String `tfsdk:"selector"`
	FallbackToDefault types.Bool   `tfsdk:"fallback_to_default"`
	AutoApply         types.Object `tfsdk:"auto_apply"`
}

func (e SegmentTF) ToAPIReq() Segment {
	segment := Segment{
		ID:                e.ID.ValueString(),
		Name:              e.Name.ValueString(),
		Selector:          e.Selector.ValueString(),
		FallbackToDefault: e.FallbackToDefault.ValueBool(),
	}

	if !e.AutoApply.IsNull() {
		attrs := e.AutoApply.Attributes()
		if enabled, ok := attrs["enabled"]; ok {
			if boolVal, ok := enabled.(types.Bool); ok {
				segment.AutoApply = &AutoApplyConfig{
					Enabled: boolVal.ValueBool(),
				}
			}
		}
	}

	return segment
}
