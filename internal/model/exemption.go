package model

import (
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Exemption struct {
	ID         string    `json:"id"`
	Metric     string    `json:"metric,omitempty"`
	KeepLabels []string  `json:"keep_labels,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
	ManagedBy  string    `json:"managed_by,omitempty"`
}

func (e Exemption) ToTF() ExemptionTF {
	return ExemptionTF{
		ID:         types.StringValue(e.ID),
		Metric:     types.StringValue(e.Metric),
		KeepLabels: toTypesStringSlice(e.KeepLabels),
		CreatedAt:  types.Int64Value(e.CreatedAt.UnixMilli()),
		UpdatedAt:  types.Int64Value(e.UpdatedAt.UnixMilli()),
	}
}

type ExemptionTF struct {
	ID         types.String   `tfsdk:"id"`
	Metric     types.String   `tfsdk:"metric"`
	KeepLabels []types.String `tfsdk:"keep_labels"`
	CreatedAt  types.Int64    `tfsdk:"created_at"`
	UpdatedAt  types.Int64    `tfsdk:"updated_at"`

	LastUpdated types.String `tfsdk:"-"`
}

func (e ExemptionTF) ToAPIReq() Exemption {
	return Exemption{
		ID:         e.ID.ValueString(),
		Metric:     e.Metric.ValueString(),
		KeepLabels: toStringSlice(e.KeepLabels),
		CreatedAt:  time.UnixMilli(e.CreatedAt.ValueInt64()),
		UpdatedAt:  time.UnixMilli(e.UpdatedAt.ValueInt64()),
		ManagedBy:  managedByTF,
	}
}
