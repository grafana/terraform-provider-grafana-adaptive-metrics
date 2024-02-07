package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type defaultBoolFalse struct{}

var _ defaults.Bool = defaultBoolFalse{}

func (d defaultBoolFalse) Description(_ context.Context) string {
	return "value defaults to false"
}

func (d defaultBoolFalse) MarkdownDescription(_ context.Context) string {
	return "value defaults to false"
}

func (d defaultBoolFalse) DefaultBool(_ context.Context, _ defaults.BoolRequest, resp *defaults.BoolResponse) {
	resp.PlanValue = types.BoolValue(false)
}

type defaultEmptyList struct{}

var _ defaults.List = defaultEmptyList{}

func (d defaultEmptyList) Description(_ context.Context) string {
	return "value defaults to []"
}

func (d defaultEmptyList) MarkdownDescription(_ context.Context) string {
	return "value defaults to []"
}

func (d defaultEmptyList) DefaultList(_ context.Context, _ defaults.ListRequest, resp *defaults.ListResponse) {
	resp.PlanValue = types.ListValueMust(types.StringType, []attr.Value{})
}
