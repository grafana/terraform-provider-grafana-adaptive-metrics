package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func testAutoApplyObject(enabled bool, policy *string) types.Object {
	gateTypes := map[string]attr.Type{"policy": types.StringType}
	gate := types.ObjectNull(gateTypes)
	if policy != nil {
		gate = types.ObjectValueMust(
			gateTypes,
			map[string]attr.Value{"policy": types.StringValue(*policy)},
		)
	}

	return types.ObjectValueMust(
		map[string]attr.Type{
			"enabled": types.BoolType,
			"gate":    types.ObjectType{AttrTypes: gateTypes},
		},
		map[string]attr.Value{
			"enabled": types.BoolValue(enabled),
			"gate":    gate,
		},
	)
}

func testNullAutoApplyObject() types.Object {
	return types.ObjectNull(map[string]attr.Type{
		"enabled": types.BoolType,
		"gate": types.ObjectType{AttrTypes: map[string]attr.Type{
			"policy": types.StringType,
		}},
	})
}

func testStringPointer(value string) *string {
	return &value
}
