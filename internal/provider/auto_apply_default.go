package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type staticObjectDefault struct {
	value types.Object
}

func (d staticObjectDefault) Description(context.Context) string {
	return d.value.String()
}

func (d staticObjectDefault) MarkdownDescription(ctx context.Context) string {
	return d.Description(ctx)
}

func (d staticObjectDefault) DefaultObject(_ context.Context, _ defaults.ObjectRequest, response *defaults.ObjectResponse) {
	response.PlanValue = d.value
}

func disabledAutoApplyDefault() defaults.Object {
	gateTypes := map[string]attr.Type{"policy": types.StringType}

	return staticObjectDefault{
		value: types.ObjectValueMust(
			map[string]attr.Type{
				"enabled": types.BoolType,
				"gate":    types.ObjectType{AttrTypes: gateTypes},
			},
			map[string]attr.Value{
				"enabled": types.BoolValue(false),
				"gate":    types.ObjectNull(gateTypes),
			},
		),
	}
}
