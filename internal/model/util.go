package model

import "github.com/hashicorp/terraform-plugin-framework/types"

func toTypesStringSlice(in []string) []types.String {
	out := make([]types.String, len(in))
	for i, s := range in {
		out[i] = types.StringValue(s)
	}
	return out
}

func toStringSlice(in []types.String) []string {
	out := make([]string, len(in))
	for i, s := range in {
		out[i] = s.ValueString()
	}
	return out
}
