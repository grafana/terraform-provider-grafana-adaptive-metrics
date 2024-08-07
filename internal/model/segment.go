package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Segment struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Selector          string `json:"selector"`
	FallbackToDefault bool   `json:"fallback_to_default"`
}

func (e Segment) ToTF() SegmentTF {
	return SegmentTF{
		ID:                types.StringValue(e.ID),
		Name:              types.StringValue(e.Name),
		Selector:          types.StringValue(e.Selector),
		FallbackToDefault: types.BoolValue(e.FallbackToDefault),
	}
}

type SegmentTF struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	Selector          types.String `tfsdk:"selector"`
	FallbackToDefault types.Bool   `tfsdk:"fallback_to_default"`
}

func (e SegmentTF) ToAPIReq() Segment {
	return Segment{
		ID:                e.ID.ValueString(),
		Name:              e.Name.ValueString(),
		Selector:          e.Selector.ValueString(),
		FallbackToDefault: e.FallbackToDefault.ValueBool(),
	}
}
