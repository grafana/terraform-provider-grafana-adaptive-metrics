package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type exemptionMatchTypeValidator struct{}

func (exemptionMatchTypeValidator) Description(context.Context) string {
	return "value must be exact, prefix, or suffix"
}
func (v exemptionMatchTypeValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}
func (exemptionMatchTypeValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	switch req.ConfigValue.ValueString() {
	case "exact", "prefix", "suffix":
		return
	}
	resp.Diagnostics.AddAttributeError(req.Path, "Invalid exemption match type", "Expected exact, prefix, or suffix.")
}
