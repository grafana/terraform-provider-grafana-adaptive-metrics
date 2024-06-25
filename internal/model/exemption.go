package model

import (
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Exemption struct {
	ID                     string    `json:"id"`
	Metric                 string    `json:"metric,omitempty"`
	KeepLabels             []string  `json:"keep_labels,omitempty"`
	DisableRecommendations bool      `json:"disable_recommendations,omitempty"`
	CreatedAt              time.Time `json:"created_at,omitempty"`
	UpdatedAt              time.Time `json:"updated_at,omitempty"`
	ManagedBy              string    `json:"managed_by,omitempty"`
	Reason                 string    `json:"reason,omitempty"`
}

func (e Exemption) ToTF() ExemptionTF {
	return ExemptionTF{
		ID:                     types.StringValue(e.ID),
		Metric:                 types.StringValue(e.Metric),
		KeepLabels:             toTypesStringSlice(e.KeepLabels),
		DisableRecommendations: types.BoolValue(e.DisableRecommendations),
		CreatedAt:              types.Int64Value(e.CreatedAt.UnixMilli()),
		UpdatedAt:              types.Int64Value(e.UpdatedAt.UnixMilli()),
		Reason:                 types.StringValue(e.Reason),
	}
}

type ExemptionTF struct {
	Segment                types.String   `tfsdk:"segment"`
	ID                     types.String   `tfsdk:"id"`
	Metric                 types.String   `tfsdk:"metric"`
	KeepLabels             []types.String `tfsdk:"keep_labels"`
	DisableRecommendations types.Bool     `tfsdk:"disable_recommendations"`
	Reason                 types.String   `tfsdk:"reason"`
	CreatedAt              types.Int64    `tfsdk:"created_at"`
	UpdatedAt              types.Int64    `tfsdk:"updated_at"`

	LastUpdated types.String `tfsdk:"-"`
}

func (e ExemptionTF) ToAPIReq() Exemption {
	return Exemption{
		ID:                     e.ID.ValueString(),
		Metric:                 e.Metric.ValueString(),
		KeepLabels:             toStringSlice(e.KeepLabels),
		DisableRecommendations: e.DisableRecommendations.ValueBool(),
		CreatedAt:              time.UnixMilli(e.CreatedAt.ValueInt64()),
		UpdatedAt:              time.UnixMilli(e.UpdatedAt.ValueInt64()),
		ManagedBy:              managedByTF,
		Reason:                 e.Reason.ValueString(),
	}
}
